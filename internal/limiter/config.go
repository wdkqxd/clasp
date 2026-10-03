package limiter

import (
	"log"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"time"
)

type Rule struct {
	BurstWindow     time.Duration
	BurstMax        int
	SustainedWindow time.Duration
	SustainedMax    int
}

type Tier struct {
	Prefix string
	Rule   Rule
}

type Config struct {
	Enabled bool

	IPv4Subnet int
	IPv6Subnet int

	TrustedProxies []netip.Prefix

	Default Rule
	Tiers   []Tier

	ExemptPaths []string
}

func DefaultConfig() Config {
	return Config{
		Enabled:    true,
		IPv4Subnet: 24,
		IPv6Subnet: 56,
		Default: Rule{
			BurstWindow:     10 * time.Second,
			BurstMax:        20,
			SustainedWindow: time.Minute,
			SustainedMax:    120,
		},
		Tiers: []Tier{
			{Prefix: "/api/search", Rule: Rule{10 * time.Second, 10, time.Minute, 40}},
			{Prefix: "/api/images", Rule: Rule{10 * time.Second, 10, time.Minute, 60}},
			{Prefix: "/api/wikipedia", Rule: Rule{10 * time.Second, 10, time.Minute, 60}},
			{Prefix: "/api/suggest", Rule: Rule{2 * time.Second, 40, time.Minute, 300}},
			{Prefix: "/api/bangs", Rule: Rule{2 * time.Second, 30, time.Minute, 120}},
		},
		ExemptPaths:    []string{"/api/healthz"},
		TrustedProxies: parsePrefixes("127.0.0.0/8,::1"),
	}
}

func LoadConfig() Config {
	cfg := DefaultConfig()

	if v, ok := lookupBool("LIMITER_ENABLED"); ok {
		cfg.Enabled = v
	}
	if v, ok := lookupInt("LIMITER_IPV4_SUBNET"); ok {
		cfg.IPv4Subnet = clampInt(v, 0, 32)
	}
	if v, ok := lookupInt("LIMITER_IPV6_SUBNET"); ok {
		cfg.IPv6Subnet = clampInt(v, 0, 128)
	}
	if v, ok := lookupFloat("LIMITER_SCALE"); ok && v > 0 {
		scaleRules(&cfg, v)
	}
	if v, ok := os.LookupEnv("LIMITER_TRUSTED_PROXIES"); ok {
		cfg.TrustedProxies = parsePrefixes(v)
	}
	if v, ok := os.LookupEnv("LIMITER_EXEMPT_PATHS"); ok {
		cfg.ExemptPaths = splitList(v)
	}

	return cfg
}

func scaleRules(cfg *Config, scale float64) {
	scaleRule := func(r *Rule) {
		r.BurstMax = max(1, int(float64(r.BurstMax)*scale))
		r.SustainedMax = max(1, int(float64(r.SustainedMax)*scale))
	}
	scaleRule(&cfg.Default)
	for i := range cfg.Tiers {
		scaleRule(&cfg.Tiers[i].Rule)
	}
}

func parsePrefixes(raw string) []netip.Prefix {
	var out []netip.Prefix
	for _, item := range splitList(raw) {
		p, err := netip.ParsePrefix(item)
		if err != nil {
			if a, aerr := netip.ParseAddr(item); aerr == nil {
				out = append(out, netip.PrefixFrom(a, a.BitLen()))
				continue
			}
			log.Printf("limiter: invalid IP or CIDR %q, ignored", item)
			continue
		}
		out = append(out, p.Masked())
	}
	return out
}

func splitList(raw string) []string {
	var out []string
	for _, item := range strings.Split(raw, ",") {
		if item = strings.TrimSpace(item); item != "" {
			out = append(out, item)
		}
	}
	return out
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func lookupInt(key string) (int, bool) {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return 0, false
	}
	v, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		log.Printf("limiter: %s must be an integer, got %q", key, raw)
		return 0, false
	}
	return v, true
}

func lookupBool(key string) (bool, bool) {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return false, false
	}
	v, err := strconv.ParseBool(strings.TrimSpace(raw))
	if err != nil {
		log.Printf("limiter: %s must be a boolean, got %q", key, raw)
		return false, false
	}
	return v, true
}

func lookupFloat(key string) (float64, bool) {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return 0, false
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		log.Printf("limiter: %s must be a number, got %q", key, raw)
		return 0, false
	}
	return v, true
}
