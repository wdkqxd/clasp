package aggregator

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type config struct {
	cacheTTL     time.Duration
	cacheMax     int
	minInterval  time.Duration
	cooldownBase time.Duration
	cooldownMax  time.Duration
}

func loadConfig() config {
	cfg := config{
		cacheTTL:     60 * time.Second,
		cacheMax:     1024,
		minInterval:  1200 * time.Millisecond,
		cooldownBase: 10 * time.Minute,
		cooldownMax:  time.Hour,
	}

	if v, ok := envSeconds("SEARCH_CACHE_TTL"); ok {
		cfg.cacheTTL = v
	}
	if v, ok := envInt("SEARCH_CACHE_MAX_ITEMS"); ok && v >= 0 {
		cfg.cacheMax = v
	}
	if v, ok := envMillis("ENGINE_MIN_INTERVAL"); ok && v >= 0 {
		cfg.minInterval = v
	}
	if v, ok := envSeconds("ENGINE_COOLDOWN"); ok && v > 0 {
		cfg.cooldownBase = v
	}
	if v, ok := envSeconds("ENGINE_COOLDOWN_MAX"); ok && v > 0 {
		cfg.cooldownMax = v
	}
	if cfg.cooldownMax < cfg.cooldownBase {
		cfg.cooldownMax = cfg.cooldownBase
	}
	return cfg
}

func envInt(key string) (int, bool) {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return 0, false
	}
	v, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		log.Printf("aggregator: %s must be an integer, got %q", key, raw)
		return 0, false
	}
	return v, true
}

func envSeconds(key string) (time.Duration, bool) {
	v, ok := envInt(key)
	if !ok {
		return 0, false
	}
	return time.Duration(v) * time.Second, true
}

func envMillis(key string) (time.Duration, bool) {
	v, ok := envInt(key)
	if !ok {
		return 0, false
	}
	return time.Duration(v) * time.Millisecond, true
}
