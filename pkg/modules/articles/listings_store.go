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
// maxListingPhotos caps how many photos a listing may carry. Ten was a toy
// number next to the market: listings on Krisha, the competitor sellers here
// compare us with, routinely carry 15 to 34, and its own limit is higher still.
// Forty covers every room, the kitchen, the bathroom, the entrance, the view,
// the building and a floor plan with room to spare.
//
// The cost is bounded and measured, not guessed: after our resize, watermarking
// and re-encode a photo averages 165 KB in production, so forty of them is
// about 6.6 MB per listing and a thousand listings is under 7 GB.
const maxListingPhotos = 40
const maxListingDocs = 10

type ListingInput struct {
	DealType, PropertyType                      string
	Country, Region, City, District             string
	Microdistrict, Street, House                string
	Lat, Lng                                    *float64
	Price                                       int64
	Currency                                    string // "KZT" | "RUB"
	Area                                        float64
	Rooms                                       int
	Title, Description, Contact, Cover          string
	TitleKz, TitleRu, TitleEn                   string
	DescriptionKz, DescriptionRu, DescriptionEn string
	Images                                      []string
	Documents                                   []string
	ContractURL                                 string
	LandArea                                    float64
	BuildYear                                   int
	WallMaterial                                string
	CeilingHeight                               float64
	Amenities                                   []string
	RoomSpecs                                   []RoomSpec
	NoFilters                                   bool // author attested photos are not filter-distorted
	GeoNodeID                                   *uuid.UUID
}

// GeoNodeIDStr renders the chosen location for the form's hidden field: the
// bare uuid, or empty when nothing is chosen. Without it a re-rendered form
// drops the author's location and the price silently reverts to tenge.
func (in ListingInput) GeoNodeIDStr() string {
	if in.GeoNodeID == nil {
		return ""
	}
	return in.GeoNodeID.String()
}

// CurrencyOr returns the currency to preselect in the form, defaulting to the
// tenge so the price box is never labelled with a blank symbol.
func (in ListingInput) CurrencyOr() string {
	if listingCurrencies[in.Currency] {
		return in.Currency
	}
	return "KZT"
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
	if in.Documents == nil {
		in.Documents = []string{}
	}
	if in.Amenities == nil {
		in.Amenities = []string{}
	}
	if in.Currency == "" {
		in.Currency = "KZT"
	}
	var id uuid.UUID
	err = s.db.QueryRow(ctx, fmt.Sprintf(`
		INSERT INTO listings (author_id, deal_type, property_type, country, region, city, district,
		                      microdistrict, street, house, lat, lng,
		                      price, area, rooms, title, description, contact, cover_url, images, geo_node_id,
		                      land_area, amenities, room_specs, documents, contract_url,
		                      build_year, wall_material, ceiling_height,
		                      title_kz, title_ru, title_en, description_kz, description_ru, description_en, currency, status, expires_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24::jsonb,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,'published', NOW() + INTERVAL '%d days')
		RETURNING id
	`, freeDaysVal()), authorID, in.DealType, in.PropertyType, in.Country, in.Region, in.City, in.District,
		in.Microdistrict, in.Street, in.House, in.Lat, in.Lng,
		in.Price, in.Area, in.Rooms, in.Title, in.Description, in.Contact, in.Cover, in.Images, in.GeoNodeID,
		in.LandArea, in.Amenities, string(rooms), in.Documents, in.ContractURL,
		in.BuildYear, in.WallMaterial, in.CeilingHeight,
		in.TitleKz, in.TitleRu, in.TitleEn, in.DescriptionKz, in.DescriptionRu, in.DescriptionEn, in.Currency).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create listing: %w", err)
	}
	return id, nil
}

const listingCols = `l.id, l.author_id, u.email, l.deal_type, l.property_type, l.country, l.region, l.city, l.district,
	l.microdistrict, l.street, l.house, l.lat, l.lng,
	l.price, l.area, l.rooms, l.title, l.description, l.contact, l.cover_url, l.images, l.documents, l.contract_url,
	l.title_kz, l.title_ru, l.title_en, l.description_kz, l.description_ru, l.description_en, l.currency, l.status, l.created_at,
	l.expires_at, l.promoted_until, l.featured_until, l.banner_until, l.views_count, l.contacts_count, l.land_area, l.amenities, l.room_specs,
	l.build_year, l.wall_material, l.ceiling_height,
	l.geo_node_id, l.updated_at,
	(SELECT count(*) FROM listing_reports rp WHERE rp.listing_id = l.id),
	ra.user_id, ra.name`

func scanListing(row pgx.Row) (*Listing, error) {
	var l Listing
	var id, authorID uuid.UUID
	var roomsRaw []byte
	var agentID *uuid.UUID
	var agentName *string
	var geoNode *uuid.UUID
	err := row.Scan(&id, &authorID, &l.AuthorEmail, &l.DealType, &l.PropertyType, &l.Country, &l.Region, &l.City, &l.District,
		&l.Microdistrict, &l.Street, &l.House, &l.Lat, &l.Lng,
		&l.Price, &l.Area, &l.Rooms, &l.Title, &l.Description, &l.Contact, &l.CoverURL, &l.Images, &l.Documents, &l.ContractURL,
		&l.TitleKz, &l.TitleRu, &l.TitleEn, &l.DescriptionKz, &l.DescriptionRu, &l.DescriptionEn, &l.Currency, &l.Status, &l.CreatedAt,
		&l.ExpiresAt, &l.PromotedUntil, &l.FeaturedUntil, &l.BannerUntil, &l.ViewsCount, &l.ContactsCount, &l.LandArea, &l.Amenities, &roomsRaw,
		&l.BuildYear, &l.WallMaterial, &l.CeilingHeight,
		&geoNode, &l.UpdatedAt, &l.Reports,
		&agentID, &agentName)
	if err != nil {
		return nil, err
	}
	if geoNode != nil {
		l.GeoNodeID = geoNode.String()
	}
	if len(roomsRaw) > 0 {
		_ = json.Unmarshal(roomsRaw, &l.RoomSpecs) // tolerate malformed JSON → empty
	}
	l.ID = id.String()
	l.AuthorID = authorID.String()
	if agentID != nil {
		l.AgentID = agentID.String()
		if agentName != nil {
			l.AgentName = *agentName
		}
	}
	return &l, nil
}

// Update rewrites an author's own listing in place. Everything the form can set
// is rewritten, including the currency and the denormalized address, so a fixed
// listing looks exactly like a freshly posted one. Deliberately untouched: the
// id and the paid state (promotion, featuring, banner, expiry) — editing the
// text must never quietly buy or burn a paid slot, and must not reset the clock.
func (s *ListingStore) Update(ctx context.Context, id, authorID uuid.UUID, in ListingInput) error {
	rooms, err := json.Marshal(in.RoomSpecs)
	if err != nil || in.RoomSpecs == nil {
		rooms = []byte("[]")
	}
	if in.Images == nil {
		in.Images = []string{}
	}
	if in.Documents == nil {
		in.Documents = []string{}
	}
	if in.Amenities == nil {
		in.Amenities = []string{}
	}
	if in.Currency == "" {
		in.Currency = "KZT"
	}
	ct, err := s.db.Exec(ctx, `
		UPDATE listings SET
			deal_type = $3, property_type = $4,
			country = $5, region = $6, city = $7, district = $8,
			microdistrict = $9, street = $10, house = $11, lat = $12, lng = $13,
			price = $14, area = $15, rooms = $16, title = $17, description = $18,
			contact = $19, cover_url = $20, images = $21, geo_node_id = $22,
			land_area = $23, amenities = $24, room_specs = $25::jsonb,
			documents = $26, contract_url = $27,
			title_kz = $28, title_ru = $29, title_en = $30,
			description_kz = $31, description_ru = $32, description_en = $33,
			currency = $34,
			build_year = $35, wall_material = $36, ceiling_height = $37,
			updated_at = now()
		WHERE id = $1 AND author_id = $2
	`, id, authorID, in.DealType, in.PropertyType, in.Country, in.Region, in.City, in.District,
		in.Microdistrict, in.Street, in.House, in.Lat, in.Lng,
		in.Price, in.Area, in.Rooms, in.Title, in.Description, in.Contact, in.Cover, in.Images, in.GeoNodeID,
		in.LandArea, in.Amenities, string(rooms), in.Documents, in.ContractURL,
		in.TitleKz, in.TitleRu, in.TitleEn, in.DescriptionKz, in.DescriptionRu, in.DescriptionEn, in.Currency,
		in.BuildYear, in.WallMaterial, in.CeilingHeight)
	if err != nil {
		return fmt.Errorf("update listing: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete removes an author's own listing for good. Scoped to author_id in the
// statement itself, so a guessed id belonging to someone else deletes nothing
// rather than relying on a check the caller might forget.
func (s *ListingStore) Delete(ctx context.Context, id, authorID uuid.UUID) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("delete listing: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	ct, err := tx.Exec(ctx, `DELETE FROM listings WHERE id = $1 AND author_id = $2`, id, authorID)
	if err != nil {
		return fmt.Errorf("delete listing: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	// Neither table carries a foreign key to listings — favourites are
	// polymorphic over articles and listings, and reports were built without
	// one — so both have to be swept by hand or they linger pointing at nothing.
	if _, err := tx.Exec(ctx, `DELETE FROM favorites WHERE item_type = 'listing' AND item_id = $1`, id); err != nil {
		return fmt.Errorf("delete listing favorites: %w", err)
	}
	if _, err := tx.Exec(ctx, `DELETE FROM listing_reports WHERE listing_id = $1`, id); err != nil {
		return fmt.Errorf("delete listing reports: %w", err)
	}
	return tx.Commit(ctx)
}

// MyListings returns all of an author's listings (active or expired), newest first.
func (s *ListingStore) MyListings(ctx context.Context, authorID uuid.UUID) ([]*Listing, error) {
	rows, err := s.db.Query(ctx, fmt.Sprintf(`SELECT %s FROM listings l JOIN auth_users u ON u.id = l.author_id LEFT JOIN re_agents ra ON ra.user_id = l.author_id AND ra.status = 'verified'
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

// Extend adds another free window to a listing's life (owner-only). Returns
// ErrNotFound if the listing does not exist or is not owned by author.
func (s *ListingStore) Extend(ctx context.Context, id, author uuid.UUID) error {
	return s.touch(ctx, id, author,
		fmt.Sprintf("expires_at = GREATEST(expires_at, NOW()) + INTERVAL '%d days', expiry_reminded = false", freeDaysVal()))
}

// DueReminders returns active listings expiring within 2 days that have not yet
// been reminded, so the owner can be nudged to extend.
func (s *ListingStore) DueReminders(ctx context.Context) ([]*Listing, error) {
	rows, err := s.db.Query(ctx, fmt.Sprintf(`SELECT %s FROM listings l JOIN auth_users u ON u.id = l.author_id LEFT JOIN re_agents ra ON ra.user_id = l.author_id AND ra.status = 'verified'
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

// Promote boosts a listing to the top of its section for the promote window
// (owner-only).
func (s *ListingStore) Promote(ctx context.Context, id, author uuid.UUID) error {
	return s.touch(ctx, id, author,
		fmt.Sprintf("promoted_until = GREATEST(COALESCE(promoted_until, NOW()), NOW()) + INTERVAL '%d days'", promoteDaysVal()))
}

// Feature visually highlights a listing for the feature window (owner-only).
func (s *ListingStore) Feature(ctx context.Context, id, author uuid.UUID) error {
	return s.touch(ctx, id, author,
		fmt.Sprintf("featured_until = GREATEST(COALESCE(featured_until, NOW()), NOW()) + INTERVAL '%d days'", featureDaysVal()))
}

// Banner buys the sidebar banner slot on the real-estate page for days (1..7).
// The slot is always on screen while the visitor scrolls, so it costs more than
// Top/highlight. days is a validated small int, safe to inline.
func (s *ListingStore) Banner(ctx context.Context, id, author uuid.UUID, days int) error {
	if days < 1 || days > 7 {
		days = 1
	}
	return s.touch(ctx, id, author,
		fmt.Sprintf("banner_until = GREATEST(COALESCE(banner_until, NOW()), NOW()) + INTERVAL '%d days'", days))
}

// BannerListings returns listings currently holding a banner slot (most recently
// bought first), for the real-estate sidebar.
func (s *ListingStore) BannerListings(ctx context.Context, limit int) ([]*Listing, error) {
	if limit <= 0 || limit > 10 {
		limit = 2
	}
	rows, err := s.db.Query(ctx, fmt.Sprintf(`SELECT %s FROM listings l JOIN auth_users u ON u.id = l.author_id LEFT JOIN re_agents ra ON ra.user_id = l.author_id AND ra.status = 'verified'
		WHERE l.status = 'published' AND l.expires_at > NOW() AND l.banner_until > NOW()
		ORDER BY l.banner_until DESC LIMIT $1`, listingCols), limit)
	if err != nil {
		return nil, fmt.Errorf("banner listings: %w", err)
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
	// Currency scopes a price range to one currency. Follows the country the
	// searcher filtered by; the tenge when they named no country.
	Currency  string
	RoomsMin  int
	Amenities []string // listing must have ALL of these (garage, parking, …)
	Query     string   // free text over title/description
	Limit     int
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
	// A price range only means something inside one currency: "up to 20 000 000"
	// asked in tenge must not drag in a listing priced at 20 000 000 rubles,
	// which is a different amount of money by a factor of five. So a price bound
	// carries its currency with it.
	if f.PriceMin > 0 || f.PriceMax > 0 {
		cur := f.Currency
		if !listingCurrencies[cur] {
			cur = "KZT"
		}
		add(" AND l.currency = $%d", cur)
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
	if len(f.Amenities) > 0 {
		add(" AND l.amenities @> $%d::text[]", f.Amenities) // listing has all selected amenities
	}
	if q := strings.TrimSpace(f.Query); q != "" {
		args = append(args, "%"+q+"%")
		n := len(args)
		where += fmt.Sprintf(" AND (l.title ILIKE $%d OR l.description ILIKE $%d)", n, n)
	}
	args = append(args, limit)
	q := fmt.Sprintf(`SELECT %s FROM listings l JOIN auth_users u ON u.id = l.author_id LEFT JOIN re_agents ra ON ra.user_id = l.author_id AND ra.status = 'verified'
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
	// Expired listings 404 immediately (not only after the 6h purge sweep).
	row := s.db.QueryRow(ctx, fmt.Sprintf(`SELECT %s FROM listings l JOIN auth_users u ON u.id = l.author_id LEFT JOIN re_agents ra ON ra.user_id = l.author_id AND ra.status = 'verified'
		WHERE l.id = $1 AND l.status = 'published' AND l.expires_at > NOW()`, listingCols), id)
	l, err := scanListing(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return l, nil
}

// AgentListings returns an agent's active published listings, promoted first
// then newest — the feed for their public page.
func (s *ListingStore) AgentListings(ctx context.Context, agentID uuid.UUID) ([]*Listing, error) {
	rows, err := s.db.Query(ctx, fmt.Sprintf(`SELECT %s FROM listings l JOIN auth_users u ON u.id = l.author_id LEFT JOIN re_agents ra ON ra.user_id = l.author_id AND ra.status = 'verified'
		WHERE l.author_id = $1 AND l.status = 'published' AND l.expires_at > NOW()
		ORDER BY (l.promoted_until > NOW()) DESC, l.created_at DESC`, listingCols), agentID)
	if err != nil {
		return nil, fmt.Errorf("agent listings: %w", err)
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
