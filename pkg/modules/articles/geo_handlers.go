package articles

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

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
