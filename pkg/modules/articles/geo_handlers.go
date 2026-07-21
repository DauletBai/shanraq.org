package articles

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// MapStat is one region bubble on the listings map.
type MapStat struct {
	Name  string  `json:"name"`
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
	Count int     `json:"count"`
}

// handleListingsMap returns per-region listing counts with coordinates for the
// requested country, powering the self-hosted bubble map.
func (m *Module) handleListingsMap(w http.ResponseWriter, r *http.Request) {
	country := strings.ToUpper(r.URL.Query().Get("country"))
	if country != "RU" {
		country = "KZ"
	}
	counts, err := m.listings.RegionCounts(r.Context())
	if err != nil {
		m.rt.Logger.Error("region counts", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	stats := make([]MapStat, 0, len(regionCentroids))
	for name, c := range regionCentroids {
		if c.Country != country {
			continue
		}
		stats = append(stats, MapStat{Name: name, Lat: c.Lat, Lng: c.Lng, Count: counts[name]})
	}
	sort.Slice(stats, func(i, j int) bool { return stats[i].Name < stats[j].Name })

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=120")
	_ = json.NewEncoder(w).Encode(stats)
}

// handleListingPins serves the markers for the Leaflet map.
func (m *Module) handleListingPins(w http.ResponseWriter, r *http.Request) {
	pins, err := m.listings.ListingPins(r.Context(), 500)
	if err != nil {
		m.rt.Logger.Error("listing pins", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=60")
	_ = json.NewEncoder(w).Encode(pins)
}

// handleGeoRoots returns the countries (top of the location tree) as JSON.
func (m *Module) handleGeoRoots(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	nodes, err := m.geo.Roots(r.Context(), lang)
	if err != nil {
		m.rt.Logger.Error("geo roots", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeGeoJSON(w, nodes)
}

// handleGeoChildren returns the direct children of ?parent=<uuid> as JSON.
func (m *Module) handleGeoChildren(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	id, err := uuid.Parse(r.URL.Query().Get("parent"))
	if err != nil {
		http.Error(w, "bad parent id", http.StatusBadRequest)
		return
	}
	nodes, err := m.geo.Children(r.Context(), id, lang)
	if err != nil {
		m.rt.Logger.Error("geo children", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeGeoJSON(w, nodes)
}

func writeGeoJSON(w http.ResponseWriter, nodes []GeoNode) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_ = json.NewEncoder(w).Encode(nodes)
}
