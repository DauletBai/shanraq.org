package articles

import (
	"net"
	"net/http"
	"strings"

	maxminddb "github.com/oschwald/maxminddb-golang"
)

// geoIP resolves a visitor's IP to an ISO 3166-1 alpha-2 country code using a
// local MaxMind-format database (DB-IP Lite). It exists so the audience panel
// can tell domestic (KZ) readers from genuine foreign ones — the English content
// shared on LinkedIn draws both, and by-language view counts alone cannot say
// which.
//
// Privacy: the IP is looked up and immediately discarded; only the coarse
// country code is ever counted, exactly like the User-Agent in the device/OS
// panels. Nothing per-visitor is stored, so this keeps the aggregate-only,
// no-profiling promise.
//
// The feature is optional. With no database configured (or one that fails to
// open) the reader is nil and country() returns "", so country analytics simply
// stay empty rather than breaking the request path or blocking boot.
type geoIP struct {
	r *maxminddb.Reader
}

// openGeoIP opens the country database at path. An empty path or an open error
// yields a nil *geoIP (feature off) without an error — country enrichment is
// best-effort, never a boot blocker.
func openGeoIP(path string) *geoIP {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	r, err := maxminddb.Open(path)
	if err != nil {
		return nil
	}
	return &geoIP{r: r}
}

// close releases the mapped database. Safe on a nil receiver.
func (g *geoIP) close() {
	if g != nil && g.r != nil {
		_ = g.r.Close()
	}
}

// country returns the uppercase ISO country code for ip, or "" when the reader
// is off, the IP is nil, or the database has no record for it. Only the country
// field is decoded — nothing else about the record is read.
func (g *geoIP) country(ip net.IP) string {
	if g == nil || g.r == nil || ip == nil {
		return ""
	}
	var rec struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}
	if err := g.r.Lookup(ip, &rec); err != nil {
		return ""
	}
	return strings.ToUpper(strings.TrimSpace(rec.Country.ISOCode))
}

// clientIP extracts the origin client's IP from the proxy headers Caddy sets in
// front of the app (X-Forwarded-For, then X-Real-IP), falling back to the
// connection's own RemoteAddr. Loopback and private addresses return nil so
// internal traffic is never miscounted as a country.
func clientIP(r *http.Request) net.IP {
	// X-Forwarded-For is "client, proxy1, proxy2 …" — the first entry is the
	// origin client; take the first public address in the chain.
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		for _, part := range strings.Split(xff, ",") {
			if ip := net.ParseIP(strings.TrimSpace(part)); isPublicIP(ip) {
				return ip
			}
		}
	}
	if xr := strings.TrimSpace(r.Header.Get("X-Real-IP")); xr != "" {
		if ip := net.ParseIP(xr); isPublicIP(ip) {
			return ip
		}
	}
	host := r.RemoteAddr
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	if ip := net.ParseIP(strings.TrimSpace(host)); isPublicIP(ip) {
		return ip
	}
	return nil
}

// isPublicIP reports whether ip is a routable public address — not nil,
// loopback, unspecified, link-local, or an RFC1918/ULA private range. Only such
// addresses map to a meaningful visitor country.
func isPublicIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	return !(ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() || ip.IsUnspecified())
}
