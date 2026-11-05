package aggregator

import (
	"context"
	"log"
	"math/rand/v2"
	"sync"
	"time"

	"clasp/internal/engine"
)

type cooldownError struct {
	engine string
	until  time.Time
}

func (e *cooldownError) Error() string {
	return e.engine + " is cooling down until " + e.until.UTC().Format(time.RFC3339)
}

type engineThrottle struct {
	mu sync.Mutex

	minInterval  time.Duration
	cooldownBase time.Duration
	cooldownMax  time.Duration

	next         time.Time
	blockedUntil time.Time
	strikes      int
}

func newEngineThrottle(minInterval, cooldownBase, cooldownMax time.Duration) *engineThrottle {
	return &engineThrottle{
		minInterval:  minInterval,
		cooldownBase: cooldownBase,
		cooldownMax:  cooldownMax,
	}
}

func (g *engineThrottle) reserve(ctx context.Context, name string) error {
	g.mu.Lock()
	now := time.Now()
	if now.Before(g.blockedUntil) {
		until := g.blockedUntil
		g.mu.Unlock()
		return &cooldownError{engine: name, until: until}
	}

	slot := g.next
	if slot.Before(now) {
		slot = now
	}
	g.next = slot.Add(g.minInterval + jitter(g.minInterval))
	g.mu.Unlock()

	if wait := time.Until(slot); wait > 0 {
		timer := time.NewTimer(wait)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
	return nil
}

func (g *engineThrottle) record(name string, err error) {
	if err == nil {
		g.mu.Lock()
		g.strikes = 0
		g.blockedUntil = time.Time{}
		g.mu.Unlock()
		return
	}
	if !engine.IsBlocked(err) {
		return
	}

	g.mu.Lock()
	g.strikes++
	cooldown := g.cooldownBase
	for i := 1; i < g.strikes && cooldown < g.cooldownMax; i++ {
		cooldown *= 2
	}
	if cooldown > g.cooldownMax {
		cooldown = g.cooldownMax
	}
	g.blockedUntil = time.Now().Add(cooldown)
	strikes := g.strikes
	g.mu.Unlock()

	log.Printf("aggregator: %s blocked, cooling down for %s (strike %d)", name, cooldown.Round(time.Second), strikes)
}

func jitter(d time.Duration) time.Duration {
	if d <= 0 {
		return 0
	}
	return time.Duration(rand.Int64N(int64(d/2) + 1))
}
