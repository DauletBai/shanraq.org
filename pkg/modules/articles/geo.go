package articles

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	Country     string `json:"country"` // ISO code (KZ/RU) — drives the listing currency
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
		SELECT c.id, %s AS name, c.kind, c.level, c.country,
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
		if err := rows.Scan(&id, &n.Name, &n.Kind, &n.Level, &n.Country, &n.HasChildren, &n.Lat, &n.Lng); err != nil {
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
			SELECT id, parent_id, level, kind, country, %s AS name FROM geo_nodes n WHERE id = $1
			UNION ALL
			SELECT n.id, n.parent_id, n.level, n.kind, n.country, %s FROM geo_nodes n
			JOIN up ON n.id = up.parent_id
		)
		SELECT id, name, kind, level, country FROM up ORDER BY level`, name, name)

	rows, err := s.db.Query(ctx, q, node)
	if err != nil {
		return nil, fmt.Errorf("geo ancestry: %w", err)
	}
	defer rows.Close()
	out := []GeoNode{}
	for rows.Next() {
		var n GeoNode
		var id uuid.UUID
		if err := rows.Scan(&id, &n.Name, &n.Kind, &n.Level, &n.Country); err != nil {
			return nil, err
		}
		n.ID = id.String()
		out = append(out, n)
	}
	return out, rows.Err()
}

// SetUserPlace remembers where a reader says they live, or forgets it when the
// node is nil.
//
// One node, not a set of columns: the ancestry of that node gives the district,
// the region and the country, and four separate fields would eventually
// disagree with one another — a city from one region stored beside another.
func (s *GeoStore) SetUserPlace(ctx context.Context, userID uuid.UUID, node *uuid.UUID) error {
	if node == nil {
		if _, err := s.db.Exec(ctx, `DELETE FROM user_places WHERE user_id = $1`, userID); err != nil {
			return fmt.Errorf("clear user place: %w", err)
		}
		return nil
	}
	_, err := s.db.Exec(ctx, `
		INSERT INTO user_places (user_id, geo_node_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET geo_node_id = EXCLUDED.geo_node_id, updated_at = NOW()
	`, userID, *node)
	if err != nil {
		return fmt.Errorf("set user place: %w", err)
	}
	return nil
}

// UserPlace returns where a reader said they live, or nil if they never said.
// A reader without a place is not a problem to be fixed: it means everything is
// shown to them, which is what the site did before places existed.
func (s *GeoStore) UserPlace(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error) {
	var id uuid.UUID
	err := s.db.QueryRow(ctx, `SELECT geo_node_id FROM user_places WHERE user_id = $1`, userID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("user place: %w", err)
	}
	return &id, nil
}

// PlaceLabel renders a reader's place the way it should be shown back to them:
// the place itself, then the region it sits in, so "Качар" reads as "Качар,
// Костанайская область" and cannot be confused with a namesake.
func (s *GeoStore) PlaceLabel(ctx context.Context, node uuid.UUID, lang string) (string, error) {
	chain, err := s.Ancestry(ctx, node, lang)
	if err != nil {
		return "", err
	}
	if len(chain) == 0 {
		return "", nil
	}
	// Ancestry runs country → … → node; the label reads the other way, and the
	// country is dropped: a reader in Kazakhstan does not need telling.
	self := chain[len(chain)-1].Name
	if len(chain) >= 3 {
		return self + ", " + chain[len(chain)-2].Name, nil
	}
	return self, nil
}
