package articles

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
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
}

type geoCacheEntry struct {
	body []byte
	at   time.Time
}

const geoCacheTTL = 6 * time.Hour

var geocoder = &geoCoder{
	http:  &http.Client{Timeout: 8 * time.Second},
	cache: map[string]geoCacheEntry{},
}

func (g *geoCoder) get(ctx context.Context, endpoint string) ([]byte, error) {
	g.mu.Lock()
	if e, ok := g.cache[endpoint]; ok && time.Since(e.at) < geoCacheTTL {
		body := e.body
		g.mu.Unlock()
		return body, nil
	}
	g.mu.Unlock()

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
	g.mu.Lock()
	g.cache[endpoint] = geoCacheEntry{body: body, at: time.Now()}
	g.mu.Unlock()
	return body, nil
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
	endpoint := "https://nominatim.openstreetmap.org/search?format=json&limit=1&accept-language=" +
		url.QueryEscape(lang) + "&q=" + url.QueryEscape(q)
	body, err := geocoder.get(r.Context(), endpoint)
	if err != nil {
		m.rt.Logger.Warn("geocode forward", zap.Error(err))
		http.Error(w, "geocode failed", http.StatusBadGateway)
		return
	}
	var res []struct {
		Lat  string `json:"lat"`
		Lon  string `json:"lon"`
		Name string `json:"display_name"`
	}
	if err := json.Unmarshal(body, &res); err != nil || len(res) == 0 {
		writeJSONObj(w, map[string]any{}) // empty → "not found", the JS keeps the map as-is
		return
	}
	lat, _ := strconv.ParseFloat(res[0].Lat, 64)
	lng, _ := strconv.ParseFloat(res[0].Lon, 64)
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
	endpoint := "https://nominatim.openstreetmap.org/search?format=json&limit=1&accept-language=" +
		url.QueryEscape(lang) + "&q=" + url.QueryEscape(query)
	body, err := geocoder.get(ctx, endpoint)
	if err != nil {
		return 0, 0, false
	}
	var res []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}
	if err := json.Unmarshal(body, &res); err != nil || len(res) == 0 {
		return 0, 0, false
	}
	lat, err1 := strconv.ParseFloat(res[0].Lat, 64)
	lng, err2 := strconv.ParseFloat(res[0].Lon, 64)
	if err1 != nil || err2 != nil || (lat == 0 && lng == 0) {
		return 0, 0, false
	}
	return lat, lng, true
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
