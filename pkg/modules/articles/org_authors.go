package articles

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"shanraq.org/pkg/modules/auth"
)

// Publishing as an organisation: ЖКХ «Качарец», ТОО «Водоканал», an akimat.
//
// The organisation is a right attached to a person's account, not an account of
// its own. This site requires a real name so that a person stands behind every
// text; a shared departmental login cancels that rule — afterwards nobody can
// say who wrote it. Staff change, the account outlives the person, and a notice
// from an akimat has a responsible official behind it by law rather than a
// login. So the byline shows the organisation while the moderation ledger keeps
// pointing at Ivanov.
//
// Verification is not a refinement of this feature, it is the feature. A name
// anybody can type turns the platform into a machine for impersonating an
// akimat, and one forged notice costs more than a hundred genuine ones are
// worth. Until the status is verified the name appears nowhere.

// OrgAuthor is an organisation a person may publish on behalf of.
type OrgAuthor struct {
	UserID    uuid.UUID
	Name      string
	Kind      string
	BIN       string
	PlaceID   *uuid.UUID
	PlaceName string // filled for display
	About     string
	Contact   string
	Status    string
	Reason    string
	Email     string // owner's email, for the moderation queue
	OwnerName string // the person behind it, shown under the byline
}

// Verified reports whether a moderator has granted the badge. Everything that
// displays an organisation must ask this first.
func (o OrgAuthor) Verified() bool { return o.Status == orgVerified }

// orgKinds are the shapes an organisation can take here. The list is short on
// purpose: it exists so a reader can tell in half a second whether a notice is
// official, communal or commercial, and a longer list would blur exactly that.
var orgKinds = []string{"akimat", "utility", "company", "school", "clinic", "public"}

// IsOrgKind reports whether s is a known organisation kind.
func IsOrgKind(s string) bool {
	for _, k := range orgKinds {
		if k == s {
			return true
		}
	}
	return false
}

// NormalizeOrgKind falls back to the kind that claims least, so an unrecognised
// value can never quietly present itself as a state body.
func NormalizeOrgKind(s string) string {
	if IsOrgKind(strings.TrimSpace(s)) {
		return strings.TrimSpace(s)
	}
	return "public"
}

// KindLabelKey is the i18n key for the kind's label.
func (o OrgAuthor) KindLabelKey() string { return "org.kind_" + NormalizeOrgKind(o.Kind) }

// Official reports whether this is a state body, which the badge says louder:
// a message from an akimat is read differently from one by a company.
func (o OrgAuthor) Official() bool { return NormalizeOrgKind(o.Kind) == "akimat" }

const (
	orgPending  = "pending"
	orgVerified = "verified"
	orgRejected = "rejected"
)

func isOrgStatus(s string) bool {
	return s == orgPending || s == orgVerified || s == orgRejected
}

// OrgStore persists organisational author profiles.
type OrgStore struct{ db *pgxpool.Pool }

func NewOrgStore(db *pgxpool.Pool) *OrgStore { return &OrgStore{db: db} }

const orgCols = `o.user_id, o.name, o.kind, o.bin, o.geo_node_id, o.about, o.contact,
	o.status, o.reject_reason, COALESCE(u.email, ''),
	COALESCE(u.first_name, ''), COALESCE(u.last_name, ''), COALESCE(g.name_ru, '')`

func scanOrg(row pgx.Row) (*OrgAuthor, error) {
	var o OrgAuthor
	var first, last string
	if err := row.Scan(&o.UserID, &o.Name, &o.Kind, &o.BIN, &o.PlaceID, &o.About, &o.Contact,
		&o.Status, &o.Reason, &o.Email, &first, &last, &o.PlaceName); err != nil {
		return nil, err
	}
	o.OwnerName = composeName(first, last)
	return &o, nil
}

// ByUser returns the organisation a person applied for, or nil.
func (s *OrgStore) ByUser(ctx context.Context, userID uuid.UUID) (*OrgAuthor, error) {
	row := s.db.QueryRow(ctx, `
		SELECT `+orgCols+`
		FROM org_authors o
		JOIN auth_users u ON u.id = o.user_id
		LEFT JOIN geo_nodes g ON g.id = o.geo_node_id
		WHERE o.user_id = $1`, userID)
	o, err := scanOrg(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("org by user: %w", err)
	}
	return o, nil
}

// VerifiedByUser returns the organisation only when it has been granted. Every
// byline goes through this, so an application in flight can never appear as a
// badge.
func (s *OrgStore) VerifiedByUser(ctx context.Context, userID uuid.UUID) (*OrgAuthor, error) {
	o, err := s.ByUser(ctx, userID)
	if err != nil || o == nil || !o.Verified() {
		return nil, err
	}
	return o, nil
}

// VerifiedNames returns the verified organisation name for each of the given
// authors, skipping those who have none.
//
// A feed asks about twenty-one authors at once. Asking one query per card would
// be twenty-one round trips to spell a name that is usually absent, so the
// lookup happens once and the caller reads a map.
func (s *OrgStore) VerifiedNames(ctx context.Context, users []uuid.UUID) (map[uuid.UUID]string, error) {
	out := map[uuid.UUID]string{}
	if len(users) == 0 {
		return out, nil
	}
	rows, err := s.db.Query(ctx, `
		SELECT user_id, name FROM org_authors
		WHERE status = $1 AND user_id = ANY($2)`, orgVerified, users)
	if err != nil {
		return nil, fmt.Errorf("org names: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		out[id] = name
	}
	return out, rows.Err()
}

// Pending lists applications waiting for a moderator.
func (s *OrgStore) Pending(ctx context.Context, limit int) ([]OrgAuthor, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	rows, err := s.db.Query(ctx, `
		SELECT `+orgCols+`
		FROM org_authors o
		JOIN auth_users u ON u.id = o.user_id
		LEFT JOIN geo_nodes g ON g.id = o.geo_node_id
		WHERE o.status = $1
		ORDER BY o.created_at
		LIMIT $2`, orgPending, limit)
	if err != nil {
		return nil, fmt.Errorf("pending orgs: %w", err)
	}
	defer rows.Close()
	var out []OrgAuthor
	for rows.Next() {
		o, err := scanOrg(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *o)
	}
	return out, rows.Err()
}

// Apply files or refiles an application. Refiling resets the status: an edited
// application has not been looked at yet, and a name changed after approval
// must not keep the badge it was granted under.
func (s *OrgStore) Apply(ctx context.Context, userID uuid.UUID, o OrgAuthor) error {
	name := strings.TrimSpace(o.Name)
	if name == "" {
		return errors.New("organisation name is required")
	}
	_, err := s.db.Exec(ctx, `
		INSERT INTO org_authors (user_id, name, kind, bin, geo_node_id, about, contact, status, reject_reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, '')
		ON CONFLICT (user_id) DO UPDATE SET
			name = EXCLUDED.name, kind = EXCLUDED.kind, bin = EXCLUDED.bin,
			geo_node_id = EXCLUDED.geo_node_id, about = EXCLUDED.about,
			contact = EXCLUDED.contact, status = $8, reject_reason = '',
			reviewed_at = NULL, reviewed_by = NULL, updated_at = NOW()
	`, userID, clip(name, 120), NormalizeOrgKind(o.Kind), strings.TrimSpace(o.BIN),
		o.PlaceID, clip(strings.TrimSpace(o.About), 600), clip(strings.TrimSpace(o.Contact), 200), orgPending)
	if err != nil {
		return fmt.Errorf("apply org: %w", err)
	}
	return nil
}

// SetStatus grants or refuses the badge, recording who decided and why.
func (s *OrgStore) SetStatus(ctx context.Context, userID uuid.UUID, status, reason string, reviewer *uuid.UUID) error {
	if !isOrgStatus(status) {
		return fmt.Errorf("unknown org status %q", status)
	}
	_, err := s.db.Exec(ctx, `
		UPDATE org_authors
		SET status = $2, reject_reason = $3, reviewed_at = NOW(), reviewed_by = $4, updated_at = NOW()
		WHERE user_id = $1`, userID, status, clip(strings.TrimSpace(reason), 300), reviewer)
	if err != nil {
		return fmt.Errorf("set org status: %w", err)
	}
	return nil
}

// MayPublishFor reports whether an organisation may publish for a place: its
// own territory, or anywhere inside it.
//
// An akimat of Kostanay publishing for Almaty is either a mistake or an abuse,
// and without this the verified badge becomes a pass to the whole country. An
// organisation that named no territory is unrestricted — that is what a
// nationwide body looks like, and a moderator saw the application.
func (s *OrgStore) MayPublishFor(ctx context.Context, org *OrgAuthor, place *uuid.UUID) (bool, error) {
	if org == nil || org.PlaceID == nil || place == nil {
		return true, nil
	}
	var ok bool
	err := s.db.QueryRow(ctx, `
		WITH RECURSIVE sub AS (
			SELECT id FROM geo_nodes WHERE id = $1
			UNION ALL
			SELECT g.id FROM geo_nodes g JOIN sub ON g.parent_id = sub.id
		)
		SELECT EXISTS(SELECT 1 FROM sub WHERE id = $2)`, *org.PlaceID, *place).Scan(&ok)
	if err != nil {
		return false, fmt.Errorf("org territory: %w", err)
	}
	return ok, nil
}

// ---- studio: applying to publish as an organisation ----

// OrgCabinetPage backs the studio screen where a person applies.
type OrgCabinetPage struct {
	Base
	Org     *OrgAuthor
	Kinds   []string
	Notice  string
	Error   string
	PlaceID string
}

func (m *Module) handleOrgCabinet(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	lang := m.resolveLang(w, r)
	page := OrgCabinetPage{Base: m.base(r, T(lang, "org.title"), lang), Kinds: orgKinds}
	page.Notice = orgNotice(lang, r.URL.Query().Get("ok"))
	if org, err := m.orgs.ByUser(r.Context(), userID); err != nil {
		m.rt.Logger.Error("org by user", zap.Error(err))
	} else if org != nil {
		page.Org = org
		if org.PlaceID != nil {
			page.PlaceID = org.PlaceID.String()
		}
	}
	m.render(w, "studio_org", page)
}

func orgNotice(lang, code string) string {
	if code == "" {
		return ""
	}
	msg := T(lang, "org."+code)
	if strings.HasPrefix(msg, "org.") {
		return ""
	}
	return msg
}

func (m *Module) handleOrgApply(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	lang := m.resolveLang(w, r)

	in := OrgAuthor{
		Name:    strings.TrimSpace(r.FormValue("name")),
		Kind:    NormalizeOrgKind(r.FormValue("kind")),
		BIN:     strings.TrimSpace(r.FormValue("bin")),
		About:   r.FormValue("about"),
		Contact: r.FormValue("contact"),
		PlaceID: formPlace(r),
	}
	fail := func(msg string) {
		page := OrgCabinetPage{Base: m.base(r, T(lang, "org.title"), lang), Kinds: orgKinds, Error: msg}
		page.Org = &in
		if in.PlaceID != nil {
			page.PlaceID = in.PlaceID.String()
		}
		m.render(w, "studio_org", page)
	}
	if in.Name == "" {
		fail(T(lang, "org.err_name"))
		return
	}
	// A БИН is what a moderator checks against the public register, so a
	// malformed one is refused here rather than wasting their time.
	if in.BIN != "" && !validBIN(in.BIN) {
		fail(T(lang, "org.err_bin"))
		return
	}
	if err := m.orgs.Apply(r.Context(), userID, in); err != nil {
		m.rt.Logger.Error("apply org", zap.Error(err))
		fail(T(lang, "org.err_save"))
		return
	}
	http.Redirect(w, r, "/studio/org?ok=applied", http.StatusSeeOther)
}

// ---- moderation ----

// handleOrgDecide grants or refuses the right to publish as an organisation.
//
// The whole feature rests on this handler: everything upstream is a form, and
// everything downstream reads a flag. Until somebody here says yes, a name is
// just a name in a table.
func (m *Module) handleOrgDecide(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	uid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	_ = r.ParseForm()
	status := orgRejected
	if r.FormValue("decision") == "verify" {
		status = orgVerified
	}
	reason := ""
	if status == orgRejected {
		reason = strings.TrimSpace(r.FormValue("reason"))
	}
	var reviewer *uuid.UUID
	if claims != nil {
		if id, perr := uuid.Parse(claims.Subject); perr == nil {
			reviewer = &id
		}
	}
	if err := m.orgs.SetStatus(r.Context(), uid, status, reason, reviewer); err != nil {
		m.rt.Logger.Error("org decide", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin?ok=org_set", http.StatusSeeOther)
}
