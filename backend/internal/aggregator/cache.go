package aggregator

import (
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"clasp/internal/cachekey"
)

type cacheEntry[T any] struct {
	val     T
	expires time.Time
}

type inflightCall[T any] struct {
	done chan struct{}
	val  T
}

type ttlCache[T any] struct {
	mu      sync.Mutex
	ttl     time.Duration
	max     int
	entries map[string]cacheEntry[T]
	calls   map[string]*inflightCall[T]
}

func newTTLCache[T any](ttl time.Duration, max int) *ttlCache[T] {
	return &ttlCache[T]{
		ttl:     ttl,
		max:     max,
		entries: make(map[string]cacheEntry[T]),
		calls:   make(map[string]*inflightCall[T]),
	}
}

func (c *ttlCache[T]) do(key string, keep func(T) bool, fn func() T) T {
	c.mu.Lock()
	if entry, ok := c.entries[key]; ok {
		if time.Now().Before(entry.expires) {
			c.mu.Unlock()
			return entry.val
		}
		delete(c.entries, key)
	}
	if call, ok := c.calls[key]; ok {
		c.mu.Unlock()
		<-call.done
		return call.val
	}

	call := &inflightCall[T]{done: make(chan struct{})}
	c.calls[key] = call
	c.mu.Unlock()

	val := fn()

	c.mu.Lock()
	call.val = val
	if keep == nil || keep(val) {
		c.storeLocked(key, val)
	}
	delete(c.calls, key)
	close(call.done)
	c.mu.Unlock()

	return val
}

func (c *ttlCache[T]) storeLocked(key string, val T) {
	if c.max > 0 && len(c.entries) >= c.max {
		now := time.Now()
		for k, entry := range c.entries {
			if now.After(entry.expires) {
				delete(c.entries, k)
			}
		}
		if len(c.entries) >= c.max {
			clear(c.entries)
		}
	}
	c.entries[key] = cacheEntry[T]{val: val, expires: time.Now().Add(c.ttl)}
}

func cacheKey(kind, query string, limit, page int, engines []string) string {
	names := make([]string, len(engines))
	for i, name := range engines {
		names[i] = normalizeEngineName(name)
	}
	sort.Strings(names)

	normalizedQuery := strings.Join(strings.Fields(strings.ToLower(query)), " ")
	return cachekey.Hash(
		kind,
		normalizedQuery,
		strconv.Itoa(limit),
		strconv.Itoa(page),
		strings.Join(names, ","),
	)
}

func hasOKEngine(statuses map[string]Status) bool {
	for _, status := range statuses {
		if status.Status == "ok" {
			return true
		}
	}
	return false
}
