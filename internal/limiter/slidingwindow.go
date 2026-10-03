package limiter

import (
	"sync"
	"time"
)

type hit struct {
	count   int
	retryIn time.Duration
}

type slidingWindow struct {
	window time.Duration
	limit  int

	mu      sync.Mutex
	buckets map[string][]int64
}

func newSlidingWindow(window time.Duration, limit int) *slidingWindow {
	return &slidingWindow{
		window:  window,
		limit:   limit,
		buckets: make(map[string][]int64),
	}
}

func (s *slidingWindow) incr(key string, now time.Time) hit {
	cutoff := now.Add(-s.window).UnixNano()

	s.mu.Lock()
	defer s.mu.Unlock()

	events := s.buckets[key]
	kept := events[:0]
	for _, ts := range events {
		if ts > cutoff {
			kept = append(kept, ts)
		}
	}
	nowNano := now.UnixNano()
	events = append(kept, nowNano)
	s.buckets[key] = events

	h := hit{count: len(events)}
	if len(events) > s.limit {
		stale := events[len(events)-s.limit-1]
		if d := time.Duration(stale-nowNano) + s.window; d > 0 {
			h.retryIn = d
		}
	}
	return h
}

func (s *slidingWindow) sweep(now time.Time) int {
	cutoff := now.Add(-s.window).UnixNano()

	s.mu.Lock()
	defer s.mu.Unlock()

	removed := 0
	for key, events := range s.buckets {
		if n := len(events); n == 0 || events[n-1] <= cutoff {
			delete(s.buckets, key)
			removed++
		}
	}
	return removed
}
