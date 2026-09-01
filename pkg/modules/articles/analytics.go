package articles

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/auth"
)

// Guest analytics is deliberately AGGREGATE ONLY: every hit is folded into a
// per-day counter keyed by a coarse page kind (or a named click event) and a
// single guest/registered flag. No visitor identifier, IP, or session is ever
// stored, so this honours the Privacy Policy's "minimal analytics, no
// behavioral profiling" promise while still telling us how many people read
// what, and what they click.

// pageKind maps a request path to the coarse bucket shown in the dashboard, or
// "" for paths that should not be counted (forms handled elsewhere, assets,
// API, staff areas). Kept pure for easy testing.
func pageKind(path string) string {
	switch path {
	case "/":
		return "home"
	case "/read":
		return "article"
	case "/listings", "/listings/new", "/listings/my":
		return "listings"
	case "/favorites":
		return "favorites"
	// Each standing page counts under its own name. They shared a "static"
	// bucket, which answered how many people read something and never which --
	// and the question the panel exists for is where the readers went.
	case "/about":
		return "about"
	case "/guide":
		return "guide"
	case "/formatting":
		return "formatting"
	case "/pricing":
		return "pricing"
	case "/support":
		return "support"
	case "/privacy":
		return "privacy"
	case "/terms":
		return "terms"
	case "/adam":
		return "adam"
	case "/framework":
		return "framework"
	case "/advertise":
		return "advertise"
	case "/predictions":
		return "predictions"
	case "/analytics":
		return "stats"
	case "/calculator":
		return "calculator"
	case "/rates":
		return "rates"
	case "/weather":
		return "weather"
	}
	switch {
	case strings.HasPrefix(path, "/read/"):
		return "article"
	case strings.HasPrefix(path, "/author/"):
		return "author"
	case strings.HasPrefix(path, "/agent/"):
		return "agent"
	case strings.HasPrefix(path, "/listings/"):
		return "listing"
	// The forecast pages are the largest part of the site by a wide margin and
	// counted for nothing: an unnamed path returns "", and trackTraffic skips
	// the request whole -- the view, the country, the device and the visitor
	// with it. A stranger arriving on the forecast for their own village was
	// invisible to every figure we publish.
	case strings.HasPrefix(path, "/weather/"):
		return "weather"
	case strings.HasPrefix(path, "/place/"):
		return "place"
	case strings.HasPrefix(path, "/archive/"):
		return "archive"
	}
	return ""
}

// trackedEvents is the closed set of click events the beacon may record. A
// closed set keeps arbitrary strings out of the counter table.
var trackedEvents = map[string]bool{
	"show_contact":  true, // revealed a listing's phone/contact
	"register_cta":  true, // clicked a sign-up call to action
	"login_cta":     true, // clicked a sign-in call to action
	"post_listing":  true, // clicked "post a listing"
	"write_article": true, // clicked "write"
	"follow_social": true, // clicked a social profile in the article-aside follow card
	"view_contract": true, // opened a listing's published lease contract
}

const (
	metricPage    = "page"
	metricClick   = "click"
	metricBot     = "bot"     // crawler/preview-bot page hit, by family
	metricSource  = "source"  // where a human visit came from, by referrer
	metricDevice  = "device"  // mobile / tablet / desktop
	metricOS      = "os"      // android / ios / windows / macos / linux
	metricBrowser = "browser" // chrome / safari / firefox / edge / …
	metricCountry = "country" // visitor country by IP (ISO code), IP not stored
	metricLang    = "lang"    // reading language of the served page (kz/ru/en)
	metricGeoLang = "geolang" // country|lang cross, e.g. "US|en", "datacenter|ru"
)

// readingLang reports the language a page is served in, for analytics. It mirrors
// resolveLang's precedence (query, then cookie, then the RU default) but has no
// cookie side effect, so it is safe to call from the tracking middleware.
func readingLang(r *http.Request) string {
	if q := r.URL.Query().Get("lang"); IsLang(q) {
		return q
	}
	if c, err := r.Cookie(langCookieName); err == nil && IsLang(c.Value) {
		return c.Value
	}
	return LangRU
}

// botLabel classifies a User-Agent as a crawler/bot and returns its family
// ("google", "yandex", "facebook", …), or "" for an ordinary human browser.
// The UA string is read and immediately discarded — nothing is stored, so this
// keeps the aggregate-only, no-profiling promise. Kept pure for testing.
func botLabel(ua string) string {
	u := strings.ToLower(ua)
	switch {
	case strings.TrimSpace(u) == "":
		return "other" // no UA at all → a script/bot, never a real browser
	case strings.Contains(u, "googlebot"), strings.Contains(u, "google-inspectiontool"), strings.Contains(u, "storebot-google"), strings.Contains(u, "apis-google"):
		return "google"
	case strings.Contains(u, "yandex"):
		return "yandex"
	case strings.Contains(u, "bingbot"), strings.Contains(u, "bingpreview"), strings.Contains(u, "msnbot"):
		return "bing"
	case strings.Contains(u, "facebookexternalhit"), strings.Contains(u, "facebot"), strings.Contains(u, "meta-external"):
		return "facebook"
	case strings.Contains(u, "telegrambot"), strings.Contains(u, "telegram"):
		return "telegram"
	case strings.Contains(u, "whatsapp"):
		return "whatsapp"
	case strings.Contains(u, "twitterbot"):
		return "twitter"
	case strings.Contains(u, "applebot"):
		return "apple"
	case strings.Contains(u, "ahrefsbot"), strings.Contains(u, "semrushbot"), strings.Contains(u, "mj12bot"), strings.Contains(u, "dotbot"), strings.Contains(u, "petalbot"), strings.Contains(u, "dataforseo"), strings.Contains(u, "blexbot"):
		return "seo"
	case strings.Contains(u, "gptbot"), strings.Contains(u, "oai-searchbot"), strings.Contains(u, "chatgpt"), strings.Contains(u, "claudebot"), strings.Contains(u, "anthropic"), strings.Contains(u, "ccbot"), strings.Contains(u, "perplexity"), strings.Contains(u, "bytespider"), strings.Contains(u, "amazonbot"):
		return "ai"
	case strings.Contains(u, "bot"), strings.Contains(u, "crawler"), strings.Contains(u, "spider"), strings.Contains(u, "slurp"),
		strings.Contains(u, "curl"), strings.Contains(u, "wget"), strings.Contains(u, "python"), strings.Contains(u, "go-http"),
		strings.Contains(u, "java/"), strings.Contains(u, "okhttp"), strings.Contains(u, "headless"), strings.Contains(u, "phantom"),
		strings.Contains(u, "scrapy"), strings.Contains(u, "axios"), strings.Contains(u, "node-fetch"), strings.Contains(u, "guzzle"),
		strings.Contains(u, "libwww"), strings.Contains(u, "httpclient"), strings.Contains(u, "http-client"), strings.Contains(u, "postman"),
		strings.Contains(u, "insomnia"), strings.Contains(u, "dart:io"), strings.Contains(u, "aiohttp"):
		return "other"
	}
	return ""
}

// trafficSource classifies where a human visit came from, using only the
// Referer host (no query, no path, nothing stored). Empty and same-site
// referrers are "direct". Kept pure for testing.
func trafficSource(referer, host string) string {
	if strings.TrimSpace(referer) == "" {
		return "direct"
	}
	u, err := url.Parse(referer)
	if err != nil || u.Host == "" {
		return "direct"
	}
	h := strings.TrimPrefix(strings.ToLower(u.Host), "www.")
	host = strings.TrimPrefix(strings.ToLower(host), "www.")
	if h == host || strings.HasSuffix(h, "shanraq.org") {
		return "direct" // internal navigation, not an external source
	}
	switch {
	// Google is a dozen products, and only one of them is the search engine.
	// Lumping them together told us the site had 20 search visits in a month
	// while Search Console counted one click in three — the panel was quietly
	// filing Gmail and Translate as organic search, which is the difference
	// between "SEO is starting to work" and "it has not started".
	case strings.HasPrefix(h, "mail.google."):
		return "email"
	case strings.HasPrefix(h, "translate.google."), strings.HasSuffix(h, ".translate.goog"):
		return "translate"
	case h == "google" || strings.HasPrefix(h, "google."), strings.HasPrefix(h, "news.google."):
		return "google" // the search engine (and Google News), on any TLD
	case strings.Contains(h, "google"):
		return "other" // docs, drive, groups, googleusercontent — not search
	case strings.Contains(h, "yandex"), h == "ya.ru":
		return "yandex"
	case h == "t.me", strings.Contains(h, "telegram"):
		return "telegram"
	case strings.Contains(h, "facebook"), h == "fb.com", strings.Contains(h, "fb.me"), strings.Contains(h, "lm.facebook"):
		return "facebook"
	case strings.Contains(h, "instagram"):
		return "instagram"
	case strings.Contains(h, "linkedin"), h == "lnkd.in":
		return "linkedin"
	case strings.Contains(h, "twitter"), h == "x.com", h == "t.co":
		return "twitter"
	case strings.Contains(h, "youtube"), h == "youtu.be":
		return "youtube"
	case strings.Contains(h, "whatsapp"):
		return "whatsapp"
	case strings.Contains(h, "bing"):
		return "bing"
	case strings.Contains(h, "duckduckgo"):
		return "duckduckgo"
	}
	return "other"
}

// arrivalSource reports the channel a visit ARRIVED through, and whether this
// page view is an arrival at all.
//
// The distinction is the whole point. A source counter that fires on every page
// view measures reading depth, not acquisition: internal navigation sends a
// same-host referrer, trafficSource files that under "direct", and a reader who
// came from Facebook and then opened four more pages is recorded as Facebook 1,
// Direct 4. Do that across a month and the panel shows a direct-traffic
// majority that is really the site's own menu.
//
// So: an explicit utm_source is always an arrival (campaigns are tagged on the
// entry link). An external referrer is an arrival. No referrer at all is an
// arrival — with strict-origin-when-cross-origin, same-site navigation always
// carries one, so an empty Referer means a typed URL, a bookmark or an app that
// strips it. A same-host referrer is not an arrival, and is the case that was
// poisoning the numbers.
func arrivalSource(r *http.Request) (string, bool) {
	if src := utmSource(r.URL.Query().Get("utm_source")); src != "" {
		return src, true
	}
	ref := strings.TrimSpace(r.Header.Get("Referer"))
	if ref == "" {
		return "direct", true
	}
	u, err := url.Parse(ref)
	if err != nil || u.Host == "" {
		return "direct", true
	}
	h := strings.TrimPrefix(strings.ToLower(u.Host), "www.")
	host := strings.TrimPrefix(strings.ToLower(r.Host), "www.")
	if h == host || strings.HasSuffix(h, "shanraq.org") {
		return "", false // internal navigation — not an arrival
	}
	return trafficSource(ref, r.Host), true
}

// utmSource maps an explicit ?utm_source= tag to one of our known source
// labels, so a link shared with ?utm_source=telegram is attributed even when the
// browser strips the Referer (messengers, in-app browsers). Unknown values
// return "" so the caller falls back to referrer-based classification. The
// closed mapping keeps arbitrary strings out of the counter table.
func utmSource(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "instagram", "ig":
		return "instagram"
	case "telegram", "tg":
		return "telegram"
	case "facebook", "fb", "meta":
		return "facebook"
	case "youtube", "yt":
		return "youtube"
	case "google":
		return "google"
	case "yandex":
		return "yandex"
	case "twitter", "x":
		return "twitter"
	case "linkedin":
		return "linkedin"
	case "whatsapp", "wa":
		return "whatsapp"
	case "share", "copy":
		// Copied link or the OS share sheet: we know the reader passed it on,
		// not where to. Honest label beats guessing at "copy".
		return "share"
	case "qr":
		// A QR that is not one of the printed city campaigns — a business
		// card, a slide, a sticker. Worth separating from "direct": a scan is
		// someone standing in front of something we made, which is a different
		// event from a bookmark even though neither carries a referrer.
		return "qr"
	}
	// Printed campaigns own the rest of the QR labels. posterTargets is the
	// only place they can come from, so the set stays closed.
	return posterSource(strings.ToLower(strings.TrimSpace(v)))
}

// deviceClass buckets a User-Agent into mobile / tablet / desktop. Coarse and
// pure — an aggregate device mix, never a per-visitor fingerprint. The UA is
// read and discarded, keeping the no-profiling promise.
func deviceClass(ua string) string {
	u := strings.ToLower(ua)
	switch {
	case strings.TrimSpace(u) == "":
		return "other"
	case strings.Contains(u, "ipad"), strings.Contains(u, "tablet"),
		strings.Contains(u, "android") && !strings.Contains(u, "mobile"):
		return "tablet"
	case strings.Contains(u, "mobile"), strings.Contains(u, "iphone"),
		strings.Contains(u, "ipod"), strings.Contains(u, "android"):
		return "mobile"
	default:
		return "desktop"
	}
}

// osFamily classifies the operating system. iOS is checked before macOS because
// the iPhone/iPad UA also contains "Mac OS X".
func osFamily(ua string) string {
	u := strings.ToLower(ua)
	switch {
	case strings.Contains(u, "android"):
		return "android"
	case strings.Contains(u, "iphone"), strings.Contains(u, "ipad"), strings.Contains(u, "ipod"):
		return "ios"
	case strings.Contains(u, "windows"):
		return "windows"
	case strings.Contains(u, "mac os x"), strings.Contains(u, "macintosh"):
		return "macos"
	case strings.Contains(u, "cros"):
		return "chromeos"
	case strings.Contains(u, "linux"):
		return "linux"
	default:
		return "other"
	}
}

// browserFamily classifies the browser. Order matters: Edge/Opera/Samsung/Yandex
// UAs all contain "Chrome", and Chrome's UA contains "Safari", so the more
// specific tokens are checked first.
func browserFamily(ua string) string {
	u := strings.ToLower(ua)
	switch {
	case strings.Contains(u, "edg"):
		return "edge"
	case strings.Contains(u, "opr"), strings.Contains(u, "opera"):
		return "opera"
	case strings.Contains(u, "samsungbrowser"):
		return "samsung"
	case strings.Contains(u, "yabrowser"):
		return "yandex"
	case strings.Contains(u, "firefox"), strings.Contains(u, "fxios"):
		return "firefox"
	case strings.Contains(u, "crios"), strings.Contains(u, "chrome"):
		return "chrome"
	case strings.Contains(u, "safari"):
		return "safari"
	default:
		return "other"
	}
}

// metricKey identifies one counter bucket.
type metricKey struct {
	kind  string
	label string
	guest bool
}

// Metrics buffers analytics increments in memory and flushes them to
// analytics_daily on a ticker, so a page view costs a map write under a mutex
// rather than a DB round-trip on the request path.
type Metrics struct {
	db  *pgxpool.Pool
	log *zap.Logger
	mu  sync.Mutex
	buf map[metricKey]int64
	// slots buffers visitor-slots on the same mutex and the same ticker as the
	// counters above; see analytics_slots.go for what a slot is and why the
	// identifiers in it cannot outlive the day that made them.
	slots map[slotKey]int64
	salts saltCache
	// dropped counts what the buffer had to throw away because the database
	// stayed unreachable long enough to fill it.
	dropped int64
}

// metricsBufMax caps the buffer. Each entry is a kind, a label and a count --
// a few dozen bytes -- so ten thousand is a comfortable ceiling for a site this
// size and still small enough that a database outage cannot exhaust memory.
const metricsBufMax = 10000

// metricDelta is one buffered count, ordered so a failed batch can say which of
// its entries never reached the database.
type metricDelta struct {
	k metricKey
	n int64
}

// NewMetrics returns a collector bound to the pool.
func NewMetrics(db *pgxpool.Pool, log *zap.Logger) *Metrics {
	return &Metrics{db: db, log: log, buf: map[metricKey]int64{}}
}

// inc buffers one hit. A nil collector or empty label is a no-op, so callers
// never have to guard.
func (mt *Metrics) inc(kind, label string, guest bool) {
	if mt == nil || label == "" {
		return
	}
	mt.mu.Lock()
	mt.buf[metricKey{kind, label, guest}]++
	mt.mu.Unlock()
}

// Flush writes the buffered deltas to analytics_daily and clears the buffer.
// It is best-effort: on a DB error the batch is dropped (aggregate counters,
// not billing) rather than retried forever.
func (mt *Metrics) Flush(ctx context.Context) {
	if mt == nil {
		return
	}
	mt.mu.Lock()
	if len(mt.buf) == 0 {
		mt.mu.Unlock()
		return
	}
	// Ordered, because what comes back on a failure has to be exactly what did
	// not go in. A batch's results arrive in the order it was queued, and a map
	// has no order to compare them against.
	batch := make([]metricDelta, 0, len(mt.buf))
	for k, n := range mt.buf {
		batch = append(batch, metricDelta{k, n})
	}
	mt.buf = map[metricKey]int64{}
	mt.mu.Unlock()

	b := &pgx.Batch{}
	for _, e := range batch {
		b.Queue(`
			INSERT INTO analytics_daily (day, kind, label, is_guest, n)
			VALUES (CURRENT_DATE, $1, $2, $3, $4)
			ON CONFLICT (day, kind, label, is_guest)
			DO UPDATE SET n = analytics_daily.n + EXCLUDED.n`,
			e.k.kind, e.k.label, e.k.guest, e.n)
	}
	res := mt.db.SendBatch(ctx, b)
	defer res.Close()
	for i := range batch {
		if _, err := res.Exec(); err != nil {
			mt.log.Warn("flush analytics", zap.Error(err), zap.Int("returned", len(batch)-i))
			mt.giveBack(batch[i:])
			return
		}
	}
}

// giveBack returns counts a failed batch never wrote, so a database that is
// briefly unreachable costs a delay rather than a hole. It used to cost the
// hole: the buffer was emptied before the write and a failure dropped it, which
// under-reported traffic exactly during the trouble worth measuring.
//
// The buffer has a ceiling. Beyond it the returned counts are dropped and
// counted as dropped, because a database down for hours must not be answered by
// growing a map until the process dies of it.
func (mt *Metrics) giveBack(unsent []metricDelta) {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	for _, e := range unsent {
		if len(mt.buf) >= metricsBufMax {
			mt.dropped += int64(len(unsent))
			mt.log.Warn("analytics buffer full, counts dropped",
				zap.Int64("dropped_total", mt.dropped))
			return
		}
		mt.buf[e.k] += e.n
	}
}

// Dropped reports how many buffered counts have been thrown away for want of
// room. Zero unless the database has been unreachable for a long time.
func (mt *Metrics) Dropped() int64 {
	if mt == nil {
		return 0
	}
	mt.mu.Lock()
	defer mt.mu.Unlock()
	return mt.dropped
}

// analyticsOptOutCookie marks a browser whose traffic must never be counted —
// the team's own devices. It is a persistent flag, so it keeps excluding the
// device even after the owner logs out (e.g. switching from the admin account to
// a test account to like an article), which a role/identity check alone misses.
const analyticsOptOutCookie = "shanraq_notrack"

// excluded reports whether this request's traffic must be left out of analytics
// — the team's own browsing. It returns true when the opt-out cookie is present,
// or when the signed-in user is staff (admin/operator) or one of the configured
// excluded emails (e.g. the owner's test account); in the latter cases it also
// stamps the persistent cookie (when w != nil) so the device stays excluded once
// it logs out. Kept separate so both the page tracker and the click beacon share
// one rule.
func (m *Module) excluded(w http.ResponseWriter, r *http.Request) bool {
	if c, err := r.Cookie(analyticsOptOutCookie); err == nil && c.Value == "1" {
		return true
	}
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		return false
	}
	if claims.HasAnyRole("admin", "operator") || m.excludeEmails[strings.ToLower(strings.TrimSpace(claims.Email))] {
		if w != nil {
			http.SetCookie(w, &http.Cookie{
				Name: analyticsOptOutCookie, Value: "1", Path: "/",
				MaxAge: 60 * 60 * 24 * 365, SameSite: http.SameSiteLaxMode,
			})
		}
		return true
	}
	return false
}

// Audience buckets. Every countable hit lands in exactly one of them, and only
// one is the audience.
const (
	bucketAudience   = "audience"
	bucketBot        = "bot"
	bucketDatacenter = "datacenter"
	bucketDrop       = "drop"
)

// audienceBucket decides what a hit is, from the two signals worth trusting:
// what the client calls itself, and what kind of network it came from.
//
// The network verdict has to be taken here rather than further down. It used to
// be read inside the "this is a person" branch and kept as a dimension the panel
// could slice by — so an automated client sending an ordinary Chrome string from
// a cloud address was counted as a reader in every figure at once: views,
// visitors, hosts, sources (as "direct", having no referrer), devices, OS and
// browsers. One misplaced line, and nothing downstream could correct for it.
//
// What this still does not catch: automation behind a residential proxy, which
// arrives on a consumer ISP address and is indistinguishable from a reader by
// origin. The honest limit of the rule is "declared crawlers and hosting
// networks are out"; the scroll beacon, which needs a real browser and real
// time on the page, is what separates the rest.
func audienceBucket(bot, country string) string {
	switch {
	case bot == "seo":
		return bucketDrop
	case bot != "":
		return bucketBot
	case country == datacenterLabel:
		return bucketDatacenter
	default:
		return bucketAudience
	}
}

// trackTraffic counts one page view per GET of a countable page. It reads the
// (soft-loaded) session only to tell guests from signed-in users — nothing
// about who they are is recorded.
func (m *Module) trackTraffic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// The team's own devices never count as audience — skip early (and
			// stamp the opt-out cookie) so they don't inflate guest/Direct counts,
			// even while logged out on a test account.
			if m.excluded(w, r) {
				next.ServeHTTP(w, r)
				return
			}
			if kind := pageKind(r.URL.Path); kind != "" {
				ua := r.Header.Get("User-Agent")
				bot := botLabel(ua)
				// Where the visitor's network sits is decided BEFORE the audience
				// is, because it is a verdict and not a detail. It used to be read
				// one branch lower, inside "this is a person", and recorded as a
				// dimension the panel could slice by — so an automated client that
				// sends an ordinary Chrome string from a cloud address was counted
				// as a reader in every figure at once: views, visitors, hosts,
				// sources (as "direct", having no referrer), devices, OS, browsers.
				// That single misplacement is what made the panel unreadable, and
				// no amount of correcting the numbers downstream could fix it.
				country := m.geoip.geoLabel(clientIP(r))
				switch audienceBucket(bot, country) {
				case bucketDrop:
					// Commercial SEO scanners are turned away in robots.txt and
					// excluded from analytics entirely, so they neither count as
					// guests nor clutter the bot panel.
				case bucketBot:
					// Other crawlers are counted apart so they never inflate the
					// real human audience — that was the whole point.
					m.metrics.inc(metricBot, bot, true)
				case bucketDatacenter:
					// A hosting, cloud or VPN network. Some of this is a person
					// behind a VPN, but most of it is automation that declines to
					// say so, and the two cannot be told apart from the outside.
					// Counted and visible, like a crawler, and outside the
					// audience — for the same reason.
					m.metrics.inc(metricBot, datacenterLabel, true)
					m.metrics.inc(metricCountry, country, true)
					m.metrics.inc(metricGeoLang, country+"|"+readingLang(r), true)
				default:
					_, ok := auth.ClaimsFromContext(r.Context())
					guest := !ok
					m.metrics.inc(metricPage, kind, guest)
					// Prefer an explicit utm_source (survives the referrer being
					// stripped by messengers/apps); fall back to the Referer host.
					// Only an ARRIVAL counts as a source. Every page view used to
					// increment this, and internal navigation carries a same-host
					// referrer, which trafficSource classifies as "direct" — so a
					// reader who came from Facebook and opened four more pages
					// scored Facebook 1, Direct 4. The panel then reported a
					// direct-traffic majority that was really our own menu, and
					// the acquisition channels were unreadable.
					if src, ok := arrivalSource(r); ok {
						m.metrics.inc(metricSource, src, guest)
					}
					m.metrics.inc(metricDevice, deviceClass(ua), guest)
					m.metrics.inc(metricOS, osFamily(ua), guest)
					m.metrics.inc(metricBrowser, browserFamily(ua), guest)
					// Reading language of the served page — the signal a VPN cannot
					// mask, so it distinguishes real foreign readers (English) from
					// curious locals or VPN traffic from censored countries.
					lng := readingLang(r)
					m.metrics.inc(metricLang, lng, guest)
					// Visitor country (nil geoip → no-op). The IP was resolved to a
					// coarse label above — an ISO country code — and discarded.
					// Nothing per-visitor is stored.
					if country != "" {
						m.metrics.inc(metricCountry, country, guest)
						m.metrics.inc(metricGeoLang, country+"|"+lng, guest)
					}
					// The same hit against its visitor-slot, which is where
					// hosts, visitors and visits come from.
					m.metrics.noteSlot(r.Context(), r, country == "KZ", deviceClass(ua) == "mobile")
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

// handleTrack records a named click event sent via navigator.sendBeacon. It is
// fire-and-forget and always answers 204, even for unknown events.
func (m *Module) handleTrack(w http.ResponseWriter, r *http.Request) {
	e := strings.TrimSpace(r.URL.Query().Get("e"))
	// Skip the team's own clicks (opt-out cookie / staff / excluded email) so the
	// owner's own login/like round-trips at publish time don't inflate the panel.
	if trackedEvents[e] && botLabel(r.Header.Get("User-Agent")) == "" && !m.excluded(nil, r) {
		_, ok := auth.ClaimsFromContext(r.Context())
		m.metrics.inc(metricClick, e, !ok)
	}
	w.WriteHeader(http.StatusNoContent)
}

// metricsFlushLoop persists buffered counters periodically and once more on
// shutdown, so the last partial window is not lost.
func (m *Module) metricsFlushLoop(ctx context.Context) {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	// Expiring the old salts is not housekeeping: the privacy policy says
	// yesterday's key cannot be recovered, which holds only while something
	// deletes it. Hourly, and once at boot so a restarted process does not
	// leave a stale key sitting for an hour.
	purge := time.NewTicker(time.Hour)
	defer purge.Stop()
	go m.metrics.purge(ctx)
	for {
		select {
		case <-ctx.Done():
			// Detach from the cancelled request context for the final write.
			m.metrics.Flush(context.Background())
			m.metrics.flushSlots(context.Background())
			return
		case <-t.C:
			m.metrics.Flush(ctx)
			m.metrics.flushSlots(ctx)
		case <-purge.C:
			m.metrics.purge(ctx)
		}
	}
}

// ---- dashboard read side ----

// Audience splits a count between anonymous guests and signed-in users.
type Audience struct {
	Guest      int64
	Registered int64
}

// Total is the combined count.
func (a Audience) Total() int64 { return a.Guest + a.Registered }

// GuestPageRow is one page-kind's traffic over the reporting window.
type GuestPageRow struct {
	Kind  string
	Title string
	A     Audience
	Pct   int // guest share of the busiest kind, for the bar width
}

// GuestClickRow is one named click event's counts over the window.
type GuestClickRow struct {
	Name  string
	Title string
	A     Audience
}

// simpleRowsMax is how many named rows a panel shows before the remainder is
// folded into one. Twenty covers everything with a share worth acting on.
const simpleRowsMax = 20

// countryRowsMax is the exception. Browsers and devices come from a short list
// that repeats; countries do not, and where a reader is is the one of these a
// publisher acts on -- it decides what gets written and who the advertiser is
// sold. Thirty reaches well past the handful that dominate, and the tail is
// still folded so the total adds up.
const countryRowsMax = 30

// GuestSimpleRow is one labeled aggregate over the window — a traffic source or
// a bot family — carrying a single count (no guest/registered split).
type GuestSimpleRow struct {
	Name  string
	Title string
	N     int64
	Pct   int
}

// guestTrendDays is how far back the guest-views chart looks, and
// guestTrendLabels roughly how many dates it prints along the bottom.
const (
	guestTrendDays   = 30
	guestTrendLabels = 5
)

// GuestTrendDay is one day in the 30-day guest-views chart.
type GuestTrendDay struct {
	Label string
	N     int64
	Pct   int
	// Day is the day of the month. A full-width chart has room to number every
	// bar, which reads faster than a date and shows where the month turns over.
	Day int
	// Tick marks this day's label for the X axis. Thirty dates cannot be
	// printed side by side in a third of a row, so a narrow chart draws only
	// every few of them.
	Tick bool
}

// GuestAnalytics is the aggregate audience block shown at the bottom of the
// admin dashboard.
type GuestAnalytics struct {
	Day       Audience // views today
	Week      Audience // views in the last 7 days
	Month     Audience // views in the last 30 days
	Year      Audience // views in the last 365 days
	Pages     []GuestPageRow
	Clicks    []GuestClickRow
	Sources   []GuestSimpleRow // human visits by referrer (30 days)
	Bots      []GuestSimpleRow // crawler hits by family (30 days)
	Devices   []GuestSimpleRow // mobile / tablet / desktop (30 days)
	OS        []GuestSimpleRow // operating-system mix (30 days)
	Browsers  []GuestSimpleRow // browser mix (30 days)
	Countries []GuestSimpleRow // visitor country by IP (30 days)
	Langs     []GuestSimpleRow // reading-language mix kz/ru/en (30 days)
	EnglishBy []GuestSimpleRow // where English-reading visits come from (30 days)
	VPNLangs  []GuestSimpleRow // reading language within the ЦОД/VPN bucket (30 days)
	Trend     []GuestTrendDay
	TrendFrom string // first day label in the trend (oldest)
	TrendTo   string // last day label in the trend (today)
	// Y-axis ticks, top value first. Bar heights are a percentage of the first
	// one, so a bar's height can be read back as a number instead of only being
	// comparable to its neighbours.
	TrendTicks []AxisTick
	HasData    bool
}

// guestAnalytics assembles the audience dashboard. Titles are localized here so
// the template stays declarative. Errors on any part degrade to empty rather
// than failing the whole admin page.
func (m *Module) guestAnalytics(ctx context.Context, lang string) GuestAnalytics {
	var g GuestAnalytics
	db := m.rt.DB

	// Rolling windows (last N days), not calendar weeks/months, so the tiles
	// always nest: today ≤ 7 days ≤ 30 days ≤ 365 days. Calendar windows made the
	// current week exceed the month-to-date at the start of a month, which reads
	// as a bug even though it was arithmetically correct.
	_ = db.QueryRow(ctx, `
		SELECT
		  COALESCE(SUM(n) FILTER (WHERE day = CURRENT_DATE AND is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day = CURRENT_DATE AND NOT is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= CURRENT_DATE - INTERVAL '6 days' AND is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= CURRENT_DATE - INTERVAL '6 days' AND NOT is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= CURRENT_DATE - INTERVAL '29 days' AND is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= CURRENT_DATE - INTERVAL '29 days' AND NOT is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= CURRENT_DATE - INTERVAL '364 days' AND is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= CURRENT_DATE - INTERVAL '364 days' AND NOT is_guest), 0)
		FROM analytics_daily WHERE kind = 'page'`).
		Scan(&g.Day.Guest, &g.Day.Registered, &g.Week.Guest, &g.Week.Registered, &g.Month.Guest, &g.Month.Registered, &g.Year.Guest, &g.Year.Registered)

	// Pages, last 30 days.
	if rows, err := db.Query(ctx, `
		SELECT label,
		       COALESCE(SUM(n) FILTER (WHERE is_guest), 0),
		       COALESCE(SUM(n) FILTER (WHERE NOT is_guest), 0)
		FROM analytics_daily
		WHERE kind = 'page' AND day >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY label`); err == nil {
		var maxTotal int64
		for rows.Next() {
			var r GuestPageRow
			if err := rows.Scan(&r.Kind, &r.A.Guest, &r.A.Registered); err == nil {
				r.Title = T(lang, "ag.page."+r.Kind)
				if strings.HasPrefix(r.Title, "ag.page.") {
					r.Title = r.Kind
				}
				g.Pages = append(g.Pages, r)
				if r.A.Total() > maxTotal {
					maxTotal = r.A.Total()
				}
			}
		}
		rows.Close()
		sortPageRows(g.Pages)
		for i := range g.Pages {
			g.Pages[i].Pct = barPct(g.Pages[i].A.Total(), maxTotal)
		}
	} else {
		m.rt.Logger.Warn("guest analytics pages", zap.Error(err))
	}

	// Clicks, last 30 days.
	if rows, err := db.Query(ctx, `
		SELECT label,
		       COALESCE(SUM(n) FILTER (WHERE is_guest), 0),
		       COALESCE(SUM(n) FILTER (WHERE NOT is_guest), 0)
		FROM analytics_daily
		WHERE kind = 'click' AND day >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY label`); err == nil {
		for rows.Next() {
			var r GuestClickRow
			if err := rows.Scan(&r.Name, &r.A.Guest, &r.A.Registered); err == nil {
				r.Title = T(lang, "ag.click."+r.Name)
				if strings.HasPrefix(r.Title, "ag.click.") {
					r.Title = r.Name
				}
				g.Clicks = append(g.Clicks, r)
			}
		}
		rows.Close()
		sortClickRows(g.Clicks)
	} else {
		m.rt.Logger.Warn("guest analytics clicks", zap.Error(err))
	}

	// Traffic sources and bots, last 30 days.
	g.Sources = m.simpleRows(ctx, metricSource, "ag.source.", lang)
	g.Bots = m.simpleRows(ctx, metricBot, "ag.bot.", lang)
	g.Devices = m.simpleRows(ctx, metricDevice, "ag.device.", lang)
	g.OS = m.simpleRows(ctx, metricOS, "ag.os.", lang)
	g.Browsers = m.simpleRows(ctx, metricBrowser, "ag.browser.", lang)
	g.Countries = m.simpleRowsN(ctx, metricCountry, "ag.country.", lang, countryRowsMax)
	g.Langs = m.simpleRows(ctx, metricLang, "ag.lang.", lang)
	g.EnglishBy = m.englishByGeo(ctx, lang)
	g.VPNLangs = m.langOfGeo(ctx, datacenterLabel, lang)

	// 14-day guest-views sparkline, gaps filled with zero.
	g.Trend = m.guestTrend(ctx)
	if len(g.Trend) > 0 {
		g.TrendFrom = g.Trend[0].Label
		g.TrendTo = g.Trend[len(g.Trend)-1].Label
		var peak int64
		for _, d := range g.Trend {
			if d.N > peak {
				peak = d.N
			}
		}
		g.TrendTicks = axisTicks(peak)
	}

	g.HasData = g.Year.Total() > 0 || len(g.Pages) > 0 || len(g.Bots) > 0
	return g
}

// guestTrend returns the last 14 calendar days of guest page views, oldest
// first, with missing days present as zero so the sparkline is continuous.
func (m *Module) guestTrend(ctx context.Context) []GuestTrendDay {
	counts := map[string]int64{}
	if rows, err := m.rt.DB.Query(ctx, `
		SELECT day, COALESCE(SUM(n) FILTER (WHERE is_guest), 0)
		FROM analytics_daily
		WHERE kind = 'page' AND day >= CURRENT_DATE - make_interval(days => $1)
		GROUP BY day`, guestTrendDays-1); err == nil {
		for rows.Next() {
			var d time.Time
			var n int64
			if err := rows.Scan(&d, &n); err == nil {
				counts[d.Format("2006-01-02")] = n
			}
		}
		rows.Close()
	} else {
		m.rt.Logger.Warn("guest analytics trend", zap.Error(err))
	}

	today := time.Now()
	out := make([]GuestTrendDay, 0, guestTrendDays)
	var max int64
	for i := guestTrendDays - 1; i >= 0; i-- {
		d := today.AddDate(0, 0, -i)
		n := counts[d.Format("2006-01-02")]
		if n > max {
			max = n
		}
		out = append(out, GuestTrendDay{Label: d.Format("02.01"), Day: d.Day(), N: n})
	}
	// Scale to the axis top, not to the tallest bar: otherwise the peak always
	// touches the frame and the gridlines sit at meaningless values.
	top := max
	if ticks := axisTicks(max); len(ticks) > 0 {
		top = ticks[0].N
	}
	// About five dates along the axis, whatever the window: printing every
	// third day was readable over two weeks and would crowd the card over a
	// month.
	step := (len(out) + guestTrendLabels - 1) / guestTrendLabels
	if step < 1 {
		step = 1
	}
	for i := range out {
		out[i].Pct = barPct(out[i].N, top)
		// Anchor the ticks to the last day so "today" is always labelled, and
		// step back from there.
		out[i].Tick = i%step == (len(out)-1)%step
	}
	return out
}

// AxisTick is one Y-axis gridline: the value printed beside it and its height
// as a percentage of the axis, so the label and the line can be placed from the
// bottom at exactly the same offset.
type AxisTick struct {
	N   int64
	Pct int
}

// axisIntervals is how many bands the Y axis is cut into before rounding. Five
// is the usual compromise: enough lines to read a bar off the grid, few enough
// that the labels do not crowd.
const axisIntervals = 5

// niceStep rounds a raw interval up to a value people count in — 1, 2, 2.5 or 5
// times a power of ten. Dropping 2.5 makes the ladder jump straight from 2 to 5,
// which for a peak of 1040 would put the top of the axis at 1500 and leave the
// chart mostly empty.
func niceStep(raw float64) int64 {
	if raw <= 1 {
		return 1
	}
	pow := int64(1)
	for float64(pow)*10 <= raw {
		pow *= 10
	}
	for _, m := range []float64{1, 2, 2.5, 5, 10} {
		if step := float64(pow) * m; step >= raw {
			return int64(step)
		}
	}
	return pow * 10
}

// axisTicks builds the Y axis from the data: a step the readings actually need,
// and a top that is the first multiple of that step at or above the busiest day.
// The scale therefore moves with the traffic instead of sitting on fixed marks.
func axisTicks(max int64) []AxisTick {
	if max <= 0 {
		return []AxisTick{{N: 1, Pct: 100}, {N: 0, Pct: 0}}
	}
	step := niceStep(float64(max) / axisIntervals)
	top := ((max + step - 1) / step) * step
	out := make([]AxisTick, 0, top/step+1)
	for v := top; v >= 0; v -= step {
		out = append(out, AxisTick{N: v, Pct: int(v * 100 / top)})
	}
	return out
}

// simpleRows loads a 30-day aggregate for a single-count metric kind (bots or
// sources), grouped by label, busiest first, with bar percentages and localized
// titles. Unknown labels fall back to the raw name.
// Guest traffic only, like the trend chart beside these panels. Without the
// filter "Откуда пришли живые гости" counted signed-in traffic as well: one
// family iPhone testing listings over a VPN put Norway at 58 and Austria at 37,
// of which 50 and 34 were that phone — foreign readership that did not exist.
// The tiles above carry the signed-in share as their own line; these panels
// answer a different question, and it is about the audience.
func (m *Module) simpleRows(ctx context.Context, kind, i18nPrefix, lang string) []GuestSimpleRow {
	return m.simpleRowsN(ctx, kind, i18nPrefix, lang, simpleRowsMax)
}

// simpleRowsN is simpleRows with the row cap named, for the panel that wants a
// longer list than the rest.
func (m *Module) simpleRowsN(ctx context.Context, kind, i18nPrefix, lang string, keep int) []GuestSimpleRow {
	rows, err := m.rt.DB.Query(ctx, `
		SELECT label, COALESCE(SUM(n), 0)
		FROM analytics_daily
		WHERE kind = $1 AND is_guest AND day >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY label`, kind)
	if err != nil {
		m.rt.Logger.Warn("guest analytics "+kind, zap.Error(err))
		return nil
	}
	defer rows.Close()

	var out []GuestSimpleRow
	var max int64
	for rows.Next() {
		var r GuestSimpleRow
		if err := rows.Scan(&r.Name, &r.N); err != nil {
			continue
		}
		r.Title = T(lang, i18nPrefix+r.Name)
		if strings.HasPrefix(r.Title, i18nPrefix) {
			r.Title = r.Name
		}
		out = append(out, r)
		if r.N > max {
			max = r.N
		}
	}
	for i := 1; i < len(out); i++ { // busiest first (tiny list, insertion sort)
		for j := i; j > 0 && out[j].N > out[j-1].N; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	// A hundred countries, ninety of them with one hit apiece, is a list nobody
	// reads to the end -- and its tail is where one crawler's exit node looks
	// exactly like a country. Keep the rows with a share worth naming, fold the
	// rest into one so the total still adds up.
	if len(out) > keep {
		var tail int64
		for _, r := range out[keep:] {
			tail += r.N
		}
		out = out[:keep:keep]
		if tail > 0 {
			out = append(out, GuestSimpleRow{Name: "other", Title: T(lang, "ag.rest"), N: tail})
		}
	}
	for i := range out {
		out[i].Pct = barPct(out[i].N, max)
	}
	return out
}

// englishByGeo returns where English-reading visits came from over the last 30
// days — each origin bucket (a country code, or "datacenter" for hosting/VPN)
// with its English read count, busiest first. This is the sharpest answer to
// "are my English readers genuine foreigners?": English from a residential
// foreign country is a real foreign reader; English from KZ is a curious local;
// English from datacenter/VPN is masked traffic (e.g. a reader in a censored
// country on a VPN). It reads the country|lang cross counted in trackTraffic.
func (m *Module) englishByGeo(ctx context.Context, lang string) []GuestSimpleRow {
	rows, err := m.rt.DB.Query(ctx, `
		SELECT split_part(label, '|', 1) AS geo, COALESCE(SUM(n), 0)
		FROM analytics_daily
		WHERE kind = $1 AND is_guest AND label LIKE '%|en' AND day >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY geo`, metricGeoLang)
	if err != nil {
		m.rt.Logger.Warn("guest analytics english-by-geo", zap.Error(err))
		return nil
	}
	defer rows.Close()

	var out []GuestSimpleRow
	var max int64
	for rows.Next() {
		var r GuestSimpleRow
		if err := rows.Scan(&r.Name, &r.N); err != nil {
			continue
		}
		r.Title = T(lang, "ag.country."+r.Name)
		if strings.HasPrefix(r.Title, "ag.country.") {
			r.Title = r.Name
		}
		out = append(out, r)
		if r.N > max {
			max = r.N
		}
	}
	for i := 1; i < len(out); i++ { // busiest first (tiny list, insertion sort)
		for j := i; j > 0 && out[j].N > out[j-1].N; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	for i := range out {
		out[i].Pct = barPct(out[i].N, max)
	}
	return out
}

// langOfGeo returns the reading-language split (kk/ru/en) for a single origin
// bucket (e.g. "datacenter") over the last 30 days. For the masked VPN bucket
// this is the key discriminator: Russian implies Russia/CIS (Russian speakers on
// an always-on VPN), while English implies genuine international readers — China,
// Iran, the West — who bridge through English rather than Russian. It reads the
// country|lang cross counted in trackTraffic.
func (m *Module) langOfGeo(ctx context.Context, geo, lang string) []GuestSimpleRow {
	rows, err := m.rt.DB.Query(ctx, `
		SELECT split_part(label, '|', 2) AS lng, COALESCE(SUM(n), 0)
		FROM analytics_daily
		WHERE kind = $1 AND is_guest AND label LIKE $2 AND day >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY lng`, metricGeoLang, geo+"|%")
	if err != nil {
		m.rt.Logger.Warn("guest analytics lang-of-geo", zap.Error(err))
		return nil
	}
	defer rows.Close()

	var out []GuestSimpleRow
	var max int64
	for rows.Next() {
		var r GuestSimpleRow
		if err := rows.Scan(&r.Name, &r.N); err != nil {
			continue
		}
		r.Title = T(lang, "ag.lang."+r.Name)
		if strings.HasPrefix(r.Title, "ag.lang.") {
			r.Title = r.Name
		}
		out = append(out, r)
		if r.N > max {
			max = r.N
		}
	}
	for i := 1; i < len(out); i++ { // busiest first (tiny list, insertion sort)
		for j := i; j > 0 && out[j].N > out[j-1].N; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	for i := range out {
		out[i].Pct = barPct(out[i].N, max)
	}
	return out
}

// sortPageRows orders page kinds by total traffic, busiest first.
func sortPageRows(rs []GuestPageRow) {
	for i := 1; i < len(rs); i++ {
		for j := i; j > 0 && rs[j].A.Total() > rs[j-1].A.Total(); j-- {
			rs[j], rs[j-1] = rs[j-1], rs[j]
		}
	}
}

// sortClickRows orders click events by total count, most first.
func sortClickRows(rs []GuestClickRow) {
	for i := 1; i < len(rs); i++ {
		for j := i; j > 0 && rs[j].A.Total() > rs[j-1].A.Total(); j-- {
			rs[j], rs[j-1] = rs[j-1], rs[j]
		}
	}
}

// withUTM tags a share link with its channel, so a reader who passes an article
// on through WhatsApp arrives as WhatsApp rather than as "direct".
//
// This is not decoration. Messengers and in-app browsers strip the Referer, so
// without the tag every shared link lands in the direct bucket and the panel
// cannot tell a channel that works from one that does not — which is the whole
// question at this stage. utmSource maps the value back to a known label; a
// source it does not recognise is ignored rather than stored.
func withUTM(rawURL, source string) string {
	if rawURL == "" || source == "" {
		return rawURL
	}
	sep := "?"
	if strings.Contains(rawURL, "?") {
		sep = "&"
	}
	return rawURL + sep + "utm_source=" + url.QueryEscape(source) + "&utm_medium=share"
}
