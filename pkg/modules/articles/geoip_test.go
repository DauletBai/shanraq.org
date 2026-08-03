package articles

import (
	"net"
	"net/http"
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

func TestClientIP(t *testing.T) {
	tests := []struct {
		name       string
		xff, xreal string
		remote     string
		want       string // "" means nil expected
	}{
		{"xff first public", "203.0.113.9, 10.0.0.1", "", "10.0.0.1:1234", "203.0.113.9"},
		{"xff skips private", "10.0.0.1, 203.0.113.9", "", "10.0.0.1:1234", "203.0.113.9"},
		{"xreal fallback", "", "198.51.100.7", "10.0.0.1:1234", "198.51.100.7"},
		{"remoteaddr fallback", "", "", "8.8.4.4:5555", "8.8.4.4"},
		{"all private → nil", "10.0.0.1", "192.168.0.1", "127.0.0.1:80", ""},
		{"remoteaddr no port", "", "", "8.8.4.4", "8.8.4.4"},
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
