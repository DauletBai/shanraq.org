package articles

import (
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
)

// parseMaybe parses an IP string, returning nil for unparseable input.
func parseMaybe(s string) net.IP { return net.ParseIP(s) }

func TestIsPublicIP(t *testing.T) {
	cases := map[string]bool{
		"8.8.8.8":      true,  // public
		"2.56.100.1":   true,  // public
		"2a00:1450::1": true,  // public IPv6
		"127.0.0.1":    false, // loopback
		"10.0.0.5":     false, // RFC1918
		"192.168.1.10": false, // RFC1918
		"172.16.5.5":   false, // RFC1918
		"169.254.1.1":  false, // link-local
		"::1":          false, // IPv6 loopback
		"fc00::1":      false, // ULA private
		"not-an-ip":    false, // unparseable → nil
		"":             false,
	}
	for in, want := range cases {
		if got := isPublicIP(parseMaybe(in)); got != want {
			t.Errorf("isPublicIP(%q) = %v, want %v", in, got, want)
		}
	}
}

// clientIP reads only RemoteAddr now. The headers in these cases are exactly
// the ones a visitor can set on the way in, and none of them may move the
// answer: the trusted-proxy middleware upstream has already decided it.
func TestClientIP(t *testing.T) {
	tests := []struct {
		name       string
		xff, xreal string
		remote     string
		want       string // "" means nil expected
	}{
		{"resolved address is used", "", "", "203.0.113.9:1234", "203.0.113.9"},
		{"no port", "", "", "8.8.4.4", "8.8.4.4"},
		{"private peer → nil", "", "", "127.0.0.1:80", ""},
		{"docker peer → nil", "", "", "172.28.0.5:1234", ""},

		// The hole this replaced. Caddy APPENDS to X-Forwarded-For rather than
		// replacing it, so a visitor who sent "8.8.8.8" arrived as
		// "8.8.8.8, <their real address>" and the old code took the first entry.
		// One header and they were a reader in the United States, inside the
		// datacenter filter as well.
		{"forged xff is ignored", "8.8.8.8, 203.0.113.9", "", "203.0.113.9:1234", "203.0.113.9"},
		{"forged x-real-ip is ignored", "", "8.8.8.8", "203.0.113.9:1234", "203.0.113.9"},
		{"forged headers cannot invent a country", "8.8.8.8", "1.1.1.1", "127.0.0.1:80", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, "/", nil)
			if tc.xff != "" {
				r.Header.Set("X-Forwarded-For", tc.xff)
			}
			if tc.xreal != "" {
				r.Header.Set("X-Real-IP", tc.xreal)
			}
			r.RemoteAddr = tc.remote
			got := clientIP(r)
			if tc.want == "" {
				if got != nil {
					t.Errorf("clientIP = %v, want nil", got)
				}
				return
			}
			if got == nil || got.String() != tc.want {
				t.Errorf("clientIP = %v, want %s", got, tc.want)
			}
		})
	}
}

// The country panel and the datacenter filter must read the same address the
// rate limiter does. Two readers of the client address is how the two drifted
// apart in the first place.
func TestClientIPMatchesTheResolvedAddress(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("X-Forwarded-For", "8.8.8.8")
	r.RemoteAddr = "203.0.113.9:1234"
	if got := clientIP(r); got == nil || got.String() != "203.0.113.9" {
		t.Errorf("clientIP = %v, want the resolved 203.0.113.9", got)
	}
	if strings.Contains(clientIPSource(), "Header") {
		t.Error("clientIP reads a request header again; the resolution belongs upstream")
	}
}

// country()/geoLabel()/isDatacenter() on a nil *geoIP must be safe no-ops.
func TestNilGeoIP(t *testing.T) {
	var g *geoIP
	ip := parseMaybe("8.8.8.8")
	if got := g.country(ip); got != "" {
		t.Errorf("nil geoIP.country = %q, want empty", got)
	}
	if got := g.geoLabel(ip); got != "" {
		t.Errorf("nil geoIP.geoLabel = %q, want empty", got)
	}
	if g.isDatacenter(ip) {
		t.Error("nil geoIP.isDatacenter = true, want false")
	}
	if g.hasASN() {
		t.Error("nil geoIP.hasASN = true, want false")
	}
	g.close() // must not panic
}

func TestDatacenterOrgsNonEmpty(t *testing.T) {
	if len(datacenterOrgs) == 0 {
		t.Fatal("datacenterOrgs must not be empty")
	}
	for _, kw := range datacenterOrgs {
		if kw != strings.ToLower(kw) {
			t.Errorf("datacenter keyword %q must be lowercase (matched against a lowercased org)", kw)
		}
	}
}

// clientIPSource returns the body of clientIP as text, so the test above can
// assert that it reads no headers — a property no amount of table-driven cases
// can pin down, because a new header would simply not be in the table.
func clientIPSource() string {
	b, err := os.ReadFile("geoip.go")
	if err != nil {
		return ""
	}
	src := string(b)
	i := strings.Index(src, "func clientIP(r *http.Request) net.IP {")
	if i < 0 {
		return ""
	}
	rest := src[i:]
	if j := strings.Index(rest, "\n}\n"); j >= 0 {
		return rest[:j]
	}
	return rest
}
