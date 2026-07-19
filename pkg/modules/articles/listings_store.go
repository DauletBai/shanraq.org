package articles

import (
	"context"
	"encoding/json"
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
// maxListingPhotos caps how many photos a listing may carry.
const maxListingPhotos = 10

type ListingInput struct {
	DealType, PropertyType             string
	Country, Region, City, Village     string
	Price                              int64
	Area                               float64
	Rooms                              int
	Title, Description, Contact, Cover string
	Images                             []string
	LandArea                           float64
	Amenities                          []string
	RoomSpecs                          []RoomSpec
	NoFilters                          bool // author attested photos are not filter-distorted
	GeoNodeID                          *uuid.UUID
}

func (s *ListingStore) Create(ctx context.Context, authorID uuid.UUID, in ListingInput) (uuid.UUID, error) {
	rooms, err := json.Marshal(in.RoomSpecs)
	if err != nil || in.RoomSpecs == nil {
		rooms = []byte("[]")
	}
	// Coerce nil slices to empty so the NOT NULL array columns keep their default.
	if in.Images == nil {
		in.Images = []string{}
	}
	if in.Amenities == nil {
		in.Amenities = []string{}
	}
	var id uuid.UUID
	err = s.db.QueryRow(ctx, `
		INSERT INTO listings (author_id, deal_type, property_type, country, region, city, village,
		                      price, area, rooms, title, description, contact, cover_url, images, geo_node_id,
		                      land_area, amenities, room_specs, status, expires_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19::jsonb,'published', NOW() + INTERVAL '21 days')
		RETURNING id
	`, authorID, in.DealType, in.PropertyType, in.Country, in.Region, in.City, in.Village,
		in.Price, in.Area, in.Rooms, in.Title, in.Description, in.Contact, in.Cover, in.Images, in.GeoNodeID,
		in.LandArea, in.Amenities, string(rooms)).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create listing: %w", err)
	}
	return id, nil
}

const listingCols = `l.id, l.author_id, u.email, l.deal_type, l.property_type, l.country, l.region, l.city, l.village,
	l.price, l.area, l.rooms, l.title, l.description, l.contact, l.cover_url, l.images, l.status, l.created_at,
	l.expires_at, l.promoted_until, l.featured_until, l.views_count, l.contacts_count, l.land_area, l.amenities, l.room_specs`

func scanListing(row pgx.Row) (*Listing, error) {
	var l Listing
	var id, authorID uuid.UUID
	var roomsRaw []byte
	err := row.Scan(&id, &authorID, &l.AuthorEmail, &l.DealType, &l.PropertyType, &l.Country, &l.Region, &l.City, &l.Village,
		&l.Price, &l.Area, &l.Rooms, &l.Title, &l.Description, &l.Contact, &l.CoverURL, &l.Images, &l.Status, &l.CreatedAt,
		&l.ExpiresAt, &l.PromotedUntil, &l.FeaturedUntil, &l.ViewsCount, &l.ContactsCount, &l.LandArea, &l.Amenities, &roomsRaw)
	if err != nil {
		return nil, err
	}
	if len(roomsRaw) > 0 {
		_ = json.Unmarshal(roomsRaw, &l.RoomSpecs) // tolerate malformed JSON → empty
	}
	l.ID = id.String()
	l.AuthorID = authorID.String()
	return &l, nil
}

// MyListings returns all of an author's listings (active or expired), newest first.
func (s *ListingStore) MyListings(ctx context.Context, authorID uuid.UUID) ([]*Listing, error) {
	rows, err := s.db.Query(ctx, fmt.Sprintf(`SELECT %s FROM listings l JOIN auth_users u ON u.id = l.author_id
		WHERE l.author_id = $1 ORDER BY l.created_at DESC`, listingCols), authorID)
	if err != nil {
		return nil, fmt.Errorf("my listings: %w", err)
	}
	defer rows.Close()
	out := []*Listing{}
	for rows.Next() {
		l, err := scanListing(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

// Extend adds 14 days to a listing's life (owner-only). Returns ErrNotFound if
// the listing does not exist or is not owned by author.
func (s *ListingStore) Extend(ctx context.Context, id, author uuid.UUID) error {
	return s.touch(ctx, id, author, "expires_at = GREATEST(expires_at, NOW()) + INTERVAL '21 days', expiry_reminded = false")
}

// DueReminders returns active listings expiring within 2 days that have not yet
// been reminded, so the owner can be nudged to extend.
func (s *ListingStore) DueReminders(ctx context.Context) ([]*Listing, error) {
	rows, err := s.db.Query(ctx, fmt.Sprintf(`SELECT %s FROM listings l JOIN auth_users u ON u.id = l.author_id
		WHERE l.status = 'published' AND l.expiry_reminded = false
		  AND l.expires_at > NOW() AND l.expires_at <= NOW() + INTERVAL '2 days'
		ORDER BY l.expires_at`, listingCols))
	if err != nil {
		return nil, fmt.Errorf("due reminders: %w", err)
	}
	defer rows.Close()
	out := []*Listing{}
	for rows.Next() {
		l, err := scanListing(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

// MarkReminded records that the expiry reminder for id has been sent.
func (s *ListingStore) MarkReminded(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `UPDATE listings SET expiry_reminded = true WHERE id = $1`, id)
	return err
}

// Promote boosts a listing to the top of its section for 7 days (owner-only).
func (s *ListingStore) Promote(ctx context.Context, id, author uuid.UUID) error {
	return s.touch(ctx, id, author, "promoted_until = GREATEST(COALESCE(promoted_until, NOW()), NOW()) + INTERVAL '7 days'")
}

// Feature visually highlights a listing for 7 days (owner-only).
func (s *ListingStore) Feature(ctx context.Context, id, author uuid.UUID) error {
	return s.touch(ctx, id, author, "featured_until = GREATEST(COALESCE(featured_until, NOW()), NOW()) + INTERVAL '7 days'")
}

func (s *ListingStore) touch(ctx context.Context, id, author uuid.UUID, set string) error {
	ct, err := s.db.Exec(ctx, fmt.Sprintf(`UPDATE listings SET %s WHERE id = $1 AND author_id = $2`, set), id, author)
	if err != nil {
		return fmt.Errorf("update listing: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
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
	where := "l.status = 'published' AND l.expires_at > NOW()"
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
		WHERE %s ORDER BY COALESCE(l.promoted_until > NOW(), false) DESC, l.created_at DESC LIMIT $%d`, listingCols, where, len(args))

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

// ListingFacets holds active-listing counts per deal type and per property
// type, for the filter badges. Total is the overall active count.
type ListingFacets struct {
	Total int
	Deal  map[string]int
	Type  map[string]int
}

// Facets counts currently-active (published, unexpired) listings grouped by
// deal type and by property type, so the filter chips can show badge counts.
func (s *ListingStore) Facets(ctx context.Context) (ListingFacets, error) {
	fc := ListingFacets{Deal: map[string]int{}, Type: map[string]int{}}
	const active = "status = 'published' AND expires_at > NOW()"

	dealRows, err := s.db.Query(ctx, `SELECT deal_type, count(*) FROM listings WHERE `+active+` GROUP BY deal_type`)
	if err != nil {
		return fc, fmt.Errorf("facet deals: %w", err)
	}
	for dealRows.Next() {
		var k string
		var n int
		if err := dealRows.Scan(&k, &n); err != nil {
			dealRows.Close()
			return fc, err
		}
		fc.Deal[k] = n
		fc.Total += n // each listing has exactly one deal type
	}
	dealRows.Close()
	if err := dealRows.Err(); err != nil {
		return fc, err
	}

	typeRows, err := s.db.Query(ctx, `SELECT property_type, count(*) FROM listings WHERE `+active+` GROUP BY property_type`)
	if err != nil {
		return fc, fmt.Errorf("facet types: %w", err)
	}
	defer typeRows.Close()
	for typeRows.Next() {
		var k string
		var n int
		if err := typeRows.Scan(&k, &n); err != nil {
			return fc, err
		}
		fc.Type[k] = n
	}
	return fc, typeRows.Err()
}

// PurgeExpired permanently deletes listings whose 21-day free window has ended,
// together with their dependent rows (reports, favorites) — enforcing the
// "then all its data is deleted" policy. Runs from the background sweep, in one
// transaction so an extend that lands mid-sweep can't orphan dependents.
func (s *ListingStore) PurgeExpired(ctx context.Context) (int64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	const expired = "expires_at < NOW()"
	if _, err := tx.Exec(ctx, `DELETE FROM listing_reports WHERE listing_id IN (SELECT id FROM listings WHERE `+expired+`)`); err != nil {
		return 0, fmt.Errorf("purge reports: %w", err)
	}
	if _, err := tx.Exec(ctx, `DELETE FROM favorites WHERE item_type = 'listing' AND item_id IN (SELECT id FROM listings WHERE `+expired+`)`); err != nil {
		return 0, fmt.Errorf("purge favorites: %w", err)
	}
	ct, err := tx.Exec(ctx, `DELETE FROM listings WHERE `+expired)
	if err != nil {
		return 0, fmt.Errorf("purge listings: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return ct.RowsAffected(), nil
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
