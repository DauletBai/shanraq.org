package articles

import (
	"net"
	"net/http"
	"strings"

	maxminddb "github.com/oschwald/maxminddb-golang"
)

// geoIP resolves a visitor's IP to an analytics label using local MaxMind-format
// databases (DB-IP Lite): a country database and an optional ASN database.
//
// It exists so the audience panel can tell domestic (KZ) readers from a genuine
// foreign audience — the English content shared on LinkedIn draws both, and
// by-language view counts alone cannot say which. The ASN database sharpens that
// further: a hit from a hosting/cloud network (a bot or a VPN exit) is bucketed
// as "datacenter" instead of a country, so the country rows reflect real eyeballs
// rather than US-hosted infrastructure — which otherwise dominates by IP.
//
// Privacy: the IP is looked up and immediately discarded; only the coarse label
// (country code or "datacenter") is ever counted, exactly like the User-Agent in
// the device/OS panels. Nothing per-visitor is stored.
//
// The feature is optional. With no country database configured (or one that
// fails to open) the reader is nil and geoLabel returns "", so country analytics
// simply stay empty rather than breaking the request path or blocking boot. The
// ASN database is independently optional: without it, isDatacenter is always
// false and every resolvable IP is bucketed by country.
type geoIP struct {
	cr *maxminddb.Reader // country
	ar *maxminddb.Reader // ASN (optional)
}

// datacenterLabel is the analytics bucket for hosting/cloud/VPN IPs — traffic
// that is a bot or a VPN exit and therefore not attributable to a reader's real
// country.
const datacenterLabel = "datacenter"

// datacenterOrgs are substrings of AS-organization names that mark a hosting /
// cloud / VPN network rather than a consumer ISP. The big clouds cover the bulk
// of automated and VPN traffic; the generic tokens ("hosting", "datacenter", …)
// catch the long tail. Matched case-insensitively against the AS org name.
var datacenterOrgs = []string{
	"amazon", "aws", "google", "microsoft", "azure", "cloudflare", "fastly",
	"akamai", "digitalocean", "digital ocean", "linode", "leaseweb", "hetzner",
	"ovh", "oracle", "alibaba", "aliyun", "tencent", "vultr", "choopa",
	"scaleway", "contabo", "m247", "datacamp", "gcore", "g-core", "netcup",
	"ionos", "kamatera", "upcloud", "psychz", "quadranet", "hivelocity",
	"limestone", "worldstream", "servers.com", "datacenter", "data center",
	"hosting", "colocation", "cloud", "dedicated server", "vpn",
}

// openGeoIP opens the country database at countryPath and, if given, the ASN
// database at asnPath. An empty country path or an open error yields a nil
// *geoIP (feature off) without an error — country enrichment is best-effort,
// never a boot blocker. A missing/broken ASN database is non-fatal: the reader
// stays nil and datacenter detection is simply skipped.
func openGeoIP(countryPath, asnPath string) *geoIP {
	if strings.TrimSpace(countryPath) == "" {
		return nil
	}
	cr, err := maxminddb.Open(countryPath)
	if err != nil {
		return nil
	}
	g := &geoIP{cr: cr}
	if p := strings.TrimSpace(asnPath); p != "" {
		if ar, err := maxminddb.Open(p); err == nil {
			g.ar = ar
		}
	}
	return g
}

// hasASN reports whether the datacenter/VPN filter is active.
func (g *geoIP) hasASN() bool { return g != nil && g.ar != nil }

// close releases the mapped databases. Safe on a nil receiver.
func (g *geoIP) close() {
	if g == nil {
		return
	}
	if g.cr != nil {
		_ = g.cr.Close()
	}
	if g.ar != nil {
		_ = g.ar.Close()
	}
}

// geoLabel returns the analytics label for ip: datacenterLabel when the IP
// belongs to a hosting/cloud/VPN network (a bot or VPN exit, not attributable to
// a reader's country), otherwise the ISO country code, or "" when the reader is
// off, the IP is nil, or nothing is known. The IP is read only for this label
// and never stored.
func (g *geoIP) geoLabel(ip net.IP) string {
	if g == nil || g.cr == nil || ip == nil {
		return ""
	}
	if g.isDatacenter(ip) {
		return datacenterLabel
	}
	return g.country(ip)
}

// country returns the uppercase ISO country code for ip, or "" when the reader
// is off or the database has no record. Only the country field is decoded.
func (g *geoIP) country(ip net.IP) string {
	if g == nil || g.cr == nil || ip == nil {
		return ""
	}
	var rec struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}
	if err := g.cr.Lookup(ip, &rec); err != nil {
		return ""
	}
	return strings.ToUpper(strings.TrimSpace(rec.Country.ISOCode))
}

// isDatacenter reports whether ip belongs to a hosting/cloud/VPN network, by
// matching its AS-organization name against datacenterOrgs. False when the ASN
// database is absent — the filter then simply does nothing.
func (g *geoIP) isDatacenter(ip net.IP) bool {
	if g == nil || g.ar == nil || ip == nil {
		return false
	}
	var rec struct {
		Org string `maxminddb:"autonomous_system_organization"`
	}
	if err := g.ar.Lookup(ip, &rec); err != nil {
		return false
	}
	org := strings.ToLower(rec.Org)
	for _, kw := range datacenterOrgs {
		if strings.Contains(org, kw) {
			return true
		}
	}
	return false
}

// clientIP returns the visitor's address, taken from RemoteAddr — which the
// server's trusted-proxy middleware has already rewritten wherever that could
// be done safely. Loopback and private addresses return nil, so internal
// traffic is never counted as a country.
//
// Reading the forwarded headers here as well was a hole, and a live one. This
// function used to take the FIRST public address out of X-Forwarded-For, and
// Caddy appends to that header rather than replacing it: a visitor who sent
// "X-Forwarded-For: 8.8.8.8" arrived as "8.8.8.8, <their real address>" and was
// counted as a reader in the United States. Anyone could colour the country
// panel and slip past the datacenter filter with one header.
//
// The middleware walks the same chain from the right and stops at the first
// address no trusted proxy vouches for. That end of the chain is written by our
// own proxy and cannot be forged, which is why the answer belongs there and not
// here. See internal/httpserver/realip.go.
func clientIP(r *http.Request) net.IP {
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
