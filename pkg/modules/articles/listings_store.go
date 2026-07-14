package articles

import (
	"context"
	"fmt"

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
}

func (s *ListingStore) Create(ctx context.Context, authorID uuid.UUID, in ListingInput) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.db.QueryRow(ctx, `
		INSERT INTO listings (author_id, deal_type, property_type, country, region, city, village,
		                      price, area, rooms, title, description, contact, cover_url, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,'published')
		RETURNING id
	`, authorID, in.DealType, in.PropertyType, in.Country, in.Region, in.City, in.Village,
		in.Price, in.Area, in.Rooms, in.Title, in.Description, in.Contact, in.Cover).Scan(&id)
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

// List returns published listings, optionally filtered by deal and property type.
func (s *ListingStore) List(ctx context.Context, deal, propertyType string, limit int) ([]*Listing, error) {
	if limit <= 0 || limit > 60 {
		limit = 30
	}
	where := "l.status = 'published'"
	args := []any{}
	if isDealType(deal) {
		args = append(args, deal)
		where += fmt.Sprintf(" AND l.deal_type = $%d", len(args))
	}
	if isPropertyType(propertyType) {
		args = append(args, propertyType)
		where += fmt.Sprintf(" AND l.property_type = $%d", len(args))
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
