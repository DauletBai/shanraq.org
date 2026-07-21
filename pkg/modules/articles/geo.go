package articles

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GeoStore reads the hierarchical location reference (country → region →
// city → district …) that powers the cascading location picker.
type GeoStore struct {
	db *pgxpool.Pool
}

func NewGeoStore(db *pgxpool.Pool) *GeoStore { return &GeoStore{db: db} }

// GeoNode is one location in the tree, localized to a requested language.
type GeoNode struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Level       int    `json:"level"`
	HasChildren bool   `json:"has_children"`
	// Lat/Lng let the listing form centre its map on the selected place.
	// Districts and countries have none, hence the pointers.
	Lat *float64 `json:"lat,omitempty"`
	Lng *float64 `json:"lng,omitempty"`
}

// geoNameCol maps a UI language to its name column; unknown falls back to ru.
func geoNameCol(lang string) string {
	switch lang {
	case LangKZ:
		return "name_kk"
	case LangEN:
		return "name_en"
	default:
		return "name_ru"
	}
}

func (s *GeoStore) query(ctx context.Context, lang, where string, args ...any) ([]GeoNode, error) {
	name := fmt.Sprintf("COALESCE(NULLIF(c.%s,''), c.name_ru)", geoNameCol(lang))
	q := fmt.Sprintf(`
		SELECT c.id, %s AS name, c.kind, c.level,
		       EXISTS(SELECT 1 FROM geo_nodes g WHERE g.parent_id = c.id) AS has_children,
		       c.lat, c.lng
		FROM geo_nodes c
		WHERE %s
		ORDER BY c.sort, c.population DESC NULLS LAST, name`, name, where)

	rows, err := s.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("geo query: %w", err)
	}
	defer rows.Close()

	out := []GeoNode{}
	for rows.Next() {
		var n GeoNode
		var id uuid.UUID
		if err := rows.Scan(&id, &n.Name, &n.Kind, &n.Level, &n.HasChildren, &n.Lat, &n.Lng); err != nil {
			return nil, err
		}
		n.ID = id.String()
		out = append(out, n)
	}
	return out, rows.Err()
}

// Roots returns the countries (top of the tree).
func (s *GeoStore) Roots(ctx context.Context, lang string) ([]GeoNode, error) {
	return s.query(ctx, lang, "c.parent_id IS NULL")
}

// Children returns the direct children of a node.
func (s *GeoStore) Children(ctx context.Context, parent uuid.UUID, lang string) ([]GeoNode, error) {
	return s.query(ctx, lang, "c.parent_id = $1", parent)
}

// Ancestry returns the path from the root down to node (inclusive), localized.
// Used to fill a listing's denormalized country/region/city/village fields.
func (s *GeoStore) Ancestry(ctx context.Context, node uuid.UUID, lang string) ([]GeoNode, error) {
	name := fmt.Sprintf("COALESCE(NULLIF(n.%s,''), n.name_ru)", geoNameCol(lang))
	q := fmt.Sprintf(`
		WITH RECURSIVE up AS (
			SELECT id, parent_id, level, kind, %s AS name FROM geo_nodes n WHERE id = $1
			UNION ALL
			SELECT n.id, n.parent_id, n.level, n.kind, %s FROM geo_nodes n
			JOIN up ON n.id = up.parent_id
		)
		SELECT id, name, kind, level FROM up ORDER BY level`, name, name)

	rows, err := s.db.Query(ctx, q, node)
	if err != nil {
		return nil, fmt.Errorf("geo ancestry: %w", err)
	}
	defer rows.Close()
	out := []GeoNode{}
	for rows.Next() {
		var n GeoNode
		var id uuid.UUID
		if err := rows.Scan(&id, &n.Name, &n.Kind, &n.Level); err != nil {
			return nil, err
		}
		n.ID = id.String()
		out = append(out, n)
	}
	return out, rows.Err()
}
