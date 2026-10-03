package limiter

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
)

var blackHole = netip.MustParseAddr("100::")

func parseAddr(raw string) (netip.Addr, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return netip.Addr{}, false
	}
	a, err := netip.ParseAddr(raw)
	if err != nil {
		return netip.Addr{}, false
	}
	if a.Zone() != "" {
		a = a.WithZone("")
	}
	if a.Is4In6() {
		a = netip.AddrFrom4(a.As4())
	}
	return a, true
}

func remoteAddr(raw string) (netip.Addr, bool) {
	if host, _, err := net.SplitHostPort(strings.TrimSpace(raw)); err == nil {
		return parseAddr(host)
	}
	return parseAddr(raw)
}

func (c Config) isTrustedProxy(a netip.Addr) bool {
	for _, p := range c.TrustedProxies {
		if p.Contains(a) {
			return true
		}
	}
	return false
}

func (c Config) clientAddr(r *http.Request) netip.Addr {
	peer, _ := remoteAddr(r.RemoteAddr)

	var chain []netip.Addr
	for _, part := range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		a, ok := parseAddr(part)
		if !ok {
			chain = nil
			break
		}
		chain = append(chain, a)
	}
	forwarded, hasForwarded := parseAddr(r.Header.Get("X-Real-IP"))

	if !c.isTrustedProxy(peer) {
		chain = nil
		hasForwarded = false
	}

	for i := len(chain) - 1; i >= 0; i-- {
		if !c.isTrustedProxy(chain[i]) {
			return chain[i]
		}
	}
	if hasForwarded {
		return forwarded
	}
	if peer.IsValid() {
		return peer
	}
	return blackHole
}

func (c Config) subnet(a netip.Addr) netip.Prefix {
	bits := c.IPv4Subnet
	if a.Is6() {
		bits = c.IPv6Subnet
	}
	return netip.PrefixFrom(a, bits).Masked()
}
