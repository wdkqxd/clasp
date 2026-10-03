package limiter

import (
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

type budget struct {
	prefix    string
	burstMax  int
	longMax   int
	burst     *slidingWindow
	sustained *slidingWindow
}

type Limiter struct {
	cfg      Config
	tiers    []budget
	fallback budget
	windows  []*slidingWindow

	quit chan struct{}
}

func New(cfg Config) *Limiter {
	l := &Limiter{cfg: cfg, quit: make(chan struct{})}

	l.fallback = newBudget("", cfg.Default)
	l.tiers = make([]budget, 0, len(cfg.Tiers))
	for _, t := range cfg.Tiers {
		l.tiers = append(l.tiers, newBudget(t.Prefix, t.Rule))
	}

	l.windows = append(l.windows, l.fallback.burst, l.fallback.sustained)
	for i := range l.tiers {
		l.windows = append(l.windows, l.tiers[i].burst, l.tiers[i].sustained)
	}

	go l.sweepLoop()
	return l
}

func newBudget(prefix string, r Rule) budget {
	return budget{
		prefix:    prefix,
		burstMax:  r.BurstMax,
		longMax:   r.SustainedMax,
		burst:     newSlidingWindow(r.BurstWindow, r.BurstMax),
		sustained: newSlidingWindow(r.SustainedWindow, r.SustainedMax),
	}
}

func (l *Limiter) Close() {
	close(l.quit)
}

func (l *Limiter) sweepLoop() {
	longest := l.cfg.Default.SustainedWindow
	for _, t := range l.cfg.Tiers {
		longest = max(longest, t.Rule.SustainedWindow)
	}

	ticker := time.NewTicker(max(longest/4, time.Second))
	defer ticker.Stop()

	for {
		select {
		case <-l.quit:
			return
		case now := <-ticker.C:
			for _, w := range l.windows {
				w.sweep(now)
			}
		}
	}
}

func (l *Limiter) Middleware(next http.Handler) http.Handler {
	if !l.cfg.Enabled {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ok, retry := l.allow(r); !ok {
			tooManyRequests(w, retry)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (l *Limiter) allow(r *http.Request) (bool, time.Duration) {
	if !l.cfg.Enabled || r.Method == http.MethodOptions || slices.Contains(l.cfg.ExemptPaths, r.URL.Path) {
		return true, 0
	}

	b := l.budgetFor(r.URL.Path)
	key := l.cfg.subnet(l.cfg.clientAddr(r)).String()
	now := time.Now()

	for _, win := range []struct {
		name string
		win  *slidingWindow
		max  int
	}{
		{"burst", b.burst, b.burstMax},
		{"sustained", b.sustained, b.longMax},
	} {
		h := win.win.incr(key, now)
		if h.count <= win.max {
			continue
		}
		log.Printf("limiter: %s over %s limit (%d)", key, win.name, win.max)
		return false, h.retryIn
	}

	return true, 0
}

func (l *Limiter) budgetFor(path string) budget {
	best := l.fallback
	for _, t := range l.tiers {
		if strings.HasPrefix(path, t.prefix) && len(t.prefix) > len(best.prefix) {
			best = t
		}
	}
	return best
}

func tooManyRequests(w http.ResponseWriter, retry time.Duration) {
	if retry < time.Second {
		retry = time.Second
	}
	w.Header().Set("Retry-After", strconv.FormatInt(int64((retry+time.Second-1)/time.Second), 10))
	w.Header().Set("RateLimit-Reset", strconv.FormatInt(time.Now().Add(retry).Unix(), 10))
	w.Header().Set("Cache-Control", "no-store")
	http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
}
