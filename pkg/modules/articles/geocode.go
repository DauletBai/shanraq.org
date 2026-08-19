package articles

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

// Geocoding proxy over OpenStreetMap Nominatim. The strict browser CSP blocks
// direct calls to external hosts, so the server forwards address↔coordinate
// lookups for the listing form (the same OSM project whose tiles we already
// use). A small in-memory cache keeps us within Nominatim's fair-use policy; at
// scale this seam can point at a self-hosted or paid geocoder without touching
// the form.
type geoCoder struct {
	http  *http.Client
	mu    sync.Mutex
	cache map[string]geoCacheEntry

	// paceMu is held across the wait, not just around the bookkeeping, so two
	// callers cannot both look at the same idle second and both decide it is
	// theirs. queued counts who is waiting behind it.
	paceMu sync.Mutex
	lastAt time.Time
	queued atomic.Int32
}

type geoCacheEntry struct {
	body []byte
	at   time.Time
}

const (
	geoCacheTTL = 6 * time.Hour

	// Nominatim's usage policy is one request per second for the whole
	// application, not per user: the limit belongs here, at the one place every
	// lookup leaves the process through, and not on the handlers.
	geoMinInterval = time.Second

	// Past this many callers waiting their turn, the queue is longer than
	// anyone will sit through — refuse quickly instead of promising a lookup
	// eight seconds from now.
	geoMaxQueued = 8

	// The cache was unbounded and never dropped an expired key, so a long
	// uptime with varied queries grew it without limit. Sweeping on write keeps
	// it a cache rather than a log of everything ever asked.
	geoCacheMax = 4096
)

// errGeoBusy is returned when the outbound queue is full: the caller should get
// a "try again", not a lookup that arrives after they have given up.
var errGeoBusy = errors.New("geocoder busy")

var geocoder = &geoCoder{
	http:  &http.Client{Timeout: 8 * time.Second},
	cache: map[string]geoCacheEntry{},
}

// wait spaces outbound calls at least geoMinInterval apart, across every
// handler and every user.
func (g *geoCoder) wait(ctx context.Context) error {
	if g.queued.Add(1) > geoMaxQueued {
		g.queued.Add(-1)
		return errGeoBusy
	}
	defer g.queued.Add(-1)

	g.paceMu.Lock()
	defer g.paceMu.Unlock()
	if d := time.Until(g.lastAt.Add(geoMinInterval)); d > 0 {
		timer := time.NewTimer(d)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
	g.lastAt = time.Now()
	return nil
}

// store caches one answer and keeps the cache bounded. Expired entries go
// first; if the map is still full after that, the oldest quarter goes too, so a
// burst of unique queries cannot pin the eviction cost on every later write.
func (g *geoCoder) store(endpoint string, body []byte) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.cache[endpoint] = geoCacheEntry{body: body, at: time.Now()}
	if len(g.cache) <= geoCacheMax {
		return
	}
	now := time.Now()
	for k, e := range g.cache {
		if now.Sub(e.at) >= geoCacheTTL {
			delete(g.cache, k)
		}
	}
	if len(g.cache) <= geoCacheMax {
		return
	}
	ages := make([]time.Time, 0, len(g.cache))
	for _, e := range g.cache {
		ages = append(ages, e.at)
	}
	sort.Slice(ages, func(i, j int) bool { return ages[i].Before(ages[j]) })
	cutoff := ages[len(ages)/4]
	for k, e := range g.cache {
		if e.at.Before(cutoff) {
			delete(g.cache, k)
		}
	}
}

func (g *geoCoder) get(ctx context.Context, endpoint string) ([]byte, error) {
	g.mu.Lock()
	if e, ok := g.cache[endpoint]; ok && time.Since(e.at) < geoCacheTTL {
		body := e.body
		g.mu.Unlock()
		return body, nil
	}
	g.mu.Unlock()

	if err := g.wait(ctx); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Shanraq/1.0 (+https://shanraq.org)")
	resp, err := g.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocode status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	g.store(endpoint, body)
	return body, nil
}

// streetRank is Nominatim's address rank for a street. Its scale runs from a
// country at 4 through a city at 16 and a suburb at 20 to a building at 30, so
// 26 is the coarsest answer that still names a street rather than a district.
//
// Nothing checked the granularity before, so any answer became the listing's
// own coordinates and therefore an "exact" pin. A query that degrades to a city
// or a suburb — an address with no street, or a street the data does not carry —
// would put a building-level claim on a district centroid.
//
// This is hardening, not the cause of a known incident: the misspelled
// "Проспект Будёного" that prompted the search returns nothing at all from
// Nominatim rather than falling back to the city. That listing got its city
// coordinates from the settlement picker instead (see place() in partials.html).
//
// Refusing a coarse answer costs nothing: the listing falls back to its
// settlement centre, which the map then labels as a settlement centre.
const streetRank = 26

// geocodeHit is one Nominatim jsonv2 result. jsonv2 is what carries place_rank;
// the plain json format does not, which is why nothing could be checked before.
type geocodeHit struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	Name        string `json:"display_name"`
	PlaceRank   int    `json:"place_rank"`
	AddressType string `json:"addresstype"`
}

// precise reports whether the hit is fine-grained enough to stand as a pin on a
// building, and parses it. A zero rank means the geocoder did not say, and an
// unlabelled answer is not evidence of precision.
func (h geocodeHit) precise() (float64, float64, bool) {
	if h.PlaceRank < streetRank {
		return 0, 0, false
	}
	lat, err1 := strconv.ParseFloat(h.Lat, 64)
	lng, err2 := strconv.ParseFloat(h.Lon, 64)
	if err1 != nil || err2 != nil || (lat == 0 && lng == 0) {
		return 0, 0, false
	}
	return lat, lng, true
}

// handleGeocode forwards an address string to coordinates (?q=...). Login-only,
// so the endpoint isn't a free public geocoder.
func (m *Module) handleGeocode(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.authorID(r); !ok {
		http.Error(w, "auth required", http.StatusUnauthorized)
		return
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if len(q) < 3 {
		http.Error(w, "query too short", http.StatusBadRequest)
		return
	}
	lang := m.resolveLang(w, r)
	endpoint := "https://nominatim.openstreetmap.org/search?format=jsonv2&limit=1&accept-language=" +
		url.QueryEscape(lang) + "&q=" + url.QueryEscape(q)
	body, err := geocoder.get(r.Context(), endpoint)
	if err != nil {
		if errors.Is(err, errGeoBusy) {
			// Our own queue is full, not the geocoder's fault: say "come back",
			// which is what the form's retry actually needs to hear.
			w.Header().Set("Retry-After", "2")
			http.Error(w, "geocoder busy", http.StatusServiceUnavailable)
			return
		}
		m.rt.Logger.Warn("geocode forward", zap.Error(err))
		http.Error(w, "geocode failed", http.StatusBadGateway)
		return
	}
	var res []geocodeHit
	if err := json.Unmarshal(body, &res); err != nil || len(res) == 0 {
		writeJSONObj(w, map[string]any{}) // empty → "not found", the JS keeps the map as-is
		return
	}
	lat, lng, ok := res[0].precise()
	if !ok {
		// A city-shaped answer to a street-shaped question is a miss, not a
		// result. Saying so leaves the map where the author put it.
		writeJSONObj(w, map[string]any{"coarse": res[0].AddressType})
		return
	}
	writeJSONObj(w, map[string]any{"lat": lat, "lng": lng, "label": res[0].Name})
}

// lookupAddress resolves a written address to coordinates. Used when a listing
// is saved: the author has already typed a street and a house number, so asking
// them to also open the map and drag a pin is asking twice for the same fact —
// and almost nobody does it, which is why the map sat empty while listings had
// full addresses. Best-effort by design: a miss leaves the listing to fall back
// on its settlement centre, never blocks the save.
func lookupAddress(ctx context.Context, lang string, parts ...string) (float64, float64, bool) {
	q := []string{}
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			q = append(q, p)
		}
	}
	query := strings.Join(q, ", ")
	if len(query) < 8 {
		return 0, 0, false
	}
	endpoint := "https://nominatim.openstreetmap.org/search?format=jsonv2&limit=1&accept-language=" +
		url.QueryEscape(lang) + "&q=" + url.QueryEscape(query)
	body, err := geocoder.get(ctx, endpoint)
	if err != nil {
		return 0, 0, false
	}
	var res []geocodeHit
	if err := json.Unmarshal(body, &res); err != nil || len(res) == 0 {
		return 0, 0, false
	}
	return res[0].precise()
}

// handleGeocodeReverse turns a pin's coordinates (?lat=&lng=) into the fine
// address fields (street, house, microdistrict) the author can keep or edit.
func (m *Module) handleGeocodeReverse(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.authorID(r); !ok {
		http.Error(w, "auth required", http.StatusUnauthorized)
		return
	}
	lat := strings.TrimSpace(r.URL.Query().Get("lat"))
	lng := strings.TrimSpace(r.URL.Query().Get("lng"))
	if lat == "" || lng == "" {
		http.Error(w, "lat/lng required", http.StatusBadRequest)
		return
	}
	lang := m.resolveLang(w, r)
	endpoint := "https://nominatim.openstreetmap.org/reverse?format=json&zoom=18&addressdetails=1&accept-language=" +
		url.QueryEscape(lang) + "&lat=" + url.QueryEscape(lat) + "&lon=" + url.QueryEscape(lng)
	body, err := geocoder.get(r.Context(), endpoint)
	if err != nil {
		if errors.Is(err, errGeoBusy) {
			// Our own queue is full, not the geocoder's fault: say "come back",
			// which is what the form's retry actually needs to hear.
			w.Header().Set("Retry-After", "2")
			http.Error(w, "geocoder busy", http.StatusServiceUnavailable)
			return
		}
		m.rt.Logger.Warn("geocode reverse", zap.Error(err))
		http.Error(w, "geocode failed", http.StatusBadGateway)
		return
	}
	var res struct {
		Address struct {
			Road          string `json:"road"`
			HouseNumber   string `json:"house_number"`
			Neighbourhood string `json:"neighbourhood"`
			Suburb        string `json:"suburb"`
			Quarter       string `json:"quarter"`
		} `json:"address"`
	}
	_ = json.Unmarshal(body, &res)
	micro := res.Address.Neighbourhood
	if micro == "" {
		micro = res.Address.Suburb
	}
	if micro == "" {
		micro = res.Address.Quarter
	}
	writeJSONObj(w, map[string]any{
		"street":        res.Address.Road,
		"house":         res.Address.HouseNumber,
		"microdistrict": micro,
	})
}

func writeJSONObj(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}
