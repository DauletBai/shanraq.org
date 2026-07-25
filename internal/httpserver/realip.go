package httpserver

import (
	"net"
	"net/http"
	"strings"
)

// parseTrustedProxies turns the configured list of CIDRs or bare IPs into
// networks. A bare IP becomes a /32 (or /128). Unparseable entries are skipped.
func parseTrustedProxies(entries []string) []*net.IPNet {
	var out []*net.IPNet
	for _, e := range entries {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		if _, n, err := net.ParseCIDR(e); err == nil {
			out = append(out, n)
			continue
		}
		if ip := net.ParseIP(e); ip != nil {
			bits := 32
			if ip.To4() == nil {
				bits = 128
			}
			out = append(out, &net.IPNet{IP: ip, Mask: net.CIDRMask(bits, bits)})
		}
	}
	return out
}

func ipInAny(ip net.IP, nets []*net.IPNet) bool {
	for _, n := range nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

// realClientIP resolves the client's IP, trusting forwarded headers only when
// the immediate TCP peer is a configured trusted proxy. It walks
// X-Forwarded-For from right (closest hop) to left and returns the first
// address that is NOT itself a trusted proxy — so a client-supplied, spoofed
// left-most entry is ignored. Returns "" when the peer is untrusted or no
// usable forwarded address is present, meaning the caller keeps the TCP peer.
func realClientIP(r *http.Request, trusted []*net.IPNet) string {
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		host = strings.TrimSpace(r.RemoteAddr)
	}
	peer := net.ParseIP(host)
	if peer == nil || !ipInAny(peer, trusted) {
		return "" // untrusted peer: never believe forwarded headers
	}
	if xr := strings.TrimSpace(r.Header.Get("X-Real-IP")); xr != "" {
		if ip := net.ParseIP(xr); ip != nil {
			return ip.String()
		}
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		for i := len(parts) - 1; i >= 0; i-- {
			ip := net.ParseIP(strings.TrimSpace(parts[i]))
			if ip != nil && !ipInAny(ip, trusted) {
				return ip.String()
			}
		}
	}
	return ""
}

// trustedRealIP rewrites r.RemoteAddr to the resolved client IP when it can be
// trusted. Unlike chi's middleware.RealIP, it does not believe forwarded
// headers from arbitrary clients. With no trusted proxies configured it is a
// no-op and the real TCP peer is preserved.
func trustedRealIP(entries []string) func(http.Handler) http.Handler {
	trusted := parseTrustedProxies(entries)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(trusted) > 0 {
				if ip := realClientIP(r, trusted); ip != "" {
					r.RemoteAddr = ip
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
