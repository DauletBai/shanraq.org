package articles

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ListingStore persists real-estate listings.
type ListingStore struct {
	db *pgxpool.Pool
}

func NewListingStore(db *pgxpool.Pool) *ListingStore { return &ListingStore{db: db} }

// ListingInput carries a submitted listing.
type ListingInput struct {
	DealType, PropertyType             string
	Country, Region, City, Village     string
	Price                              int64
	Area                               float64
	Rooms                              int
	Title, Description, Contact, Cover string
	GeoNodeID                          *uuid.UUID
}

func (s *ListingStore) Create(ctx context.Context, authorID uuid.UUID, in ListingInput) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.db.QueryRow(ctx, `
		INSERT INTO listings (author_id, deal_type, property_type, country, region, city, village,
		                      price, area, rooms, title, description, contact, cover_url, geo_node_id, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,'published')
		RETURNING id
	`, authorID, in.DealType, in.PropertyType, in.Country, in.Region, in.City, in.Village,
		in.Price, in.Area, in.Rooms, in.Title, in.Description, in.Contact, in.Cover, in.GeoNodeID).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create listing: %w", err)
	}
	return id, nil
}

const listingCols = `l.id, u.email, l.deal_type, l.property_type, l.country, l.region, l.city, l.village,
	l.price, l.area, l.rooms, l.title, l.description, l.contact, l.cover_url, l.images, l.status, l.created_at`

func scanListing(row pgx.Row) (*Listing, error) {
	var l Listing
	var id uuid.UUID
	err := row.Scan(&id, &l.AuthorEmail, &l.DealType, &l.PropertyType, &l.Country, &l.Region, &l.City, &l.Village,
		&l.Price, &l.Area, &l.Rooms, &l.Title, &l.Description, &l.Contact, &l.CoverURL, &l.Images, &l.Status, &l.CreatedAt)
	if err != nil {
		return nil, err
	}
	l.ID = id.String()
	return &l, nil
}

// ListingFilter captures the "find real estate" search criteria. Zero-valued
// fields are ignored.
type ListingFilter struct {
	Deal, PropertyType string
	GeoNodeID          *uuid.UUID // matches this node and its whole subtree
	RegionText         string     // plain-text region/city match (e.g. from a map click)
	PriceMin, PriceMax int64
	RoomsMin           int
	Query              string // free text over title/description
	Limit              int
}

// List returns published listings matching the filter, newest first.
func (s *ListingStore) List(ctx context.Context, f ListingFilter) ([]*Listing, error) {
	limit := f.Limit
	if limit <= 0 || limit > 60 {
		limit = 30
	}
	where := "l.status = 'published'"
	args := []any{}
	add := func(cond string, val any) {
		args = append(args, val)
		where += fmt.Sprintf(cond, len(args))
	}
	if isDealType(f.Deal) {
		add(" AND l.deal_type = $%d", f.Deal)
	}
	if isPropertyType(f.PropertyType) {
		add(" AND l.property_type = $%d", f.PropertyType)
	}
	if f.GeoNodeID != nil {
		args = append(args, *f.GeoNodeID)
		n := len(args)
		where += fmt.Sprintf(` AND (
			l.geo_node_id IN (
				WITH RECURSIVE sub AS (
					SELECT id FROM geo_nodes WHERE id = $%d
					UNION ALL SELECT g.id FROM geo_nodes g JOIN sub ON g.parent_id = sub.id
				) SELECT id FROM sub
			)
			OR l.region  = (SELECT name_ru FROM geo_nodes WHERE id = $%d)
			OR l.city    = (SELECT name_ru FROM geo_nodes WHERE id = $%d)
			OR l.country = (SELECT name_ru FROM geo_nodes WHERE id = $%d)
		)`, n, n, n, n)
	} else if txt := strings.TrimSpace(f.RegionText); txt != "" {
		args = append(args, txt)
		n := len(args)
		where += fmt.Sprintf(" AND (l.region = $%d OR l.city = $%d OR l.country = $%d)", n, n, n)
	}
	if f.PriceMin > 0 {
		add(" AND l.price >= $%d", f.PriceMin)
	}
	if f.PriceMax > 0 {
		add(" AND l.price <= $%d", f.PriceMax)
	}
	if f.RoomsMin > 0 {
		add(" AND l.rooms >= $%d", f.RoomsMin)
	}
	if q := strings.TrimSpace(f.Query); q != "" {
		args = append(args, "%"+q+"%")
		n := len(args)
		where += fmt.Sprintf(" AND (l.title ILIKE $%d OR l.description ILIKE $%d)", n, n)
	}
	args = append(args, limit)
	q := fmt.Sprintf(`SELECT %s FROM listings l JOIN auth_users u ON u.id = l.author_id
		WHERE %s ORDER BY l.created_at DESC LIMIT $%d`, listingCols, where, len(args))

	rows, err := s.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list listings: %w", err)
	}
	defer rows.Close()
	var out []*Listing
	for rows.Next() {
		l, err := scanListing(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

// GetByID loads a single published listing.
func (s *ListingStore) GetByID(ctx context.Context, id uuid.UUID) (*Listing, error) {
	row := s.db.QueryRow(ctx, fmt.Sprintf(`SELECT %s FROM listings l JOIN auth_users u ON u.id = l.author_id
		WHERE l.id = $1 AND l.status = 'published'`, listingCols), id)
	l, err := scanListing(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return l, nil
}
