package articles

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/auth"
)

// Staff roles that may open the admin dashboard.
var adminRoles = []string{"admin", "director", "manager", "editor"}

// Roles the dashboard can assign (admin is the seeded superuser, not assignable here).
var assignableRoles = []string{"director", "manager", "editor", "user"}

func canManageUsers(c *auth.Claims) bool { return c != nil && c.HasAnyRole("admin", "director") }
func canViewFinance(c *auth.Claims) bool {
	return c != nil && c.HasAnyRole("admin", "director", "manager")
}
func canModerate(c *auth.Claims) bool {
	return c != nil && c.HasAnyRole("admin", "director", "editor")
}

// AdminStore aggregates cross-cutting analytics for the control panel.
type AdminStore struct{ db *pgxpool.Pool }

func NewAdminStore(db *pgxpool.Pool) *AdminStore { return &AdminStore{db: db} }

type RoleCount struct {
	Role string
	N    int
	Pct  int
}
type CatViews struct {
	Slug  string
	N     int
	Views int64
	Pct   int
}
type TopArticle struct {
	Slug  string
	Title string
	Views int64
	Score int
}
type AdminComment struct {
	ID         string
	AuthorName string
	Body       string
	Slug       string
	CreatedAt  time.Time
}

// AdminStats is the full dashboard payload.
type AdminStats struct {
	Users           int
	UsersByRole     []RoleCount
	Articles        int
	Views           int64
	Comments        int
	ListingsActive  int
	ListingsExpired int
	ByCat           []CatViews
	TopArticles     []TopArticle
	RecentComments  []AdminComment
}

func (s *AdminStore) Stats(ctx context.Context) (AdminStats, error) {
	var st AdminStats
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM auth_users`).Scan(&st.Users); err != nil {
		return st, err
	}
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM articles WHERE status='published'`).Scan(&st.Articles); err != nil {
		return st, err
	}
	if err := s.db.QueryRow(ctx, `SELECT COALESCE(SUM(views_count),0) FROM articles`).Scan(&st.Views); err != nil {
		return st, err
	}
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM comments WHERE status='published'`).Scan(&st.Comments); err != nil {
		return st, err
	}
	if err := s.db.QueryRow(ctx, `SELECT
			COUNT(*) FILTER (WHERE expires_at > NOW()),
			COUNT(*) FILTER (WHERE expires_at <= NOW())
		FROM listings WHERE status='published'`).Scan(&st.ListingsActive, &st.ListingsExpired); err != nil {
		return st, err
	}

	rows, err := s.db.Query(ctx, `SELECT role, COUNT(*) FROM auth_users GROUP BY role ORDER BY COUNT(*) DESC`)
	if err != nil {
		return st, err
	}
	for rows.Next() {
		var rc RoleCount
		if err := rows.Scan(&rc.Role, &rc.N); err != nil {
			rows.Close()
			return st, err
		}
		st.UsersByRole = append(st.UsersByRole, rc)
	}
	rows.Close()
	maxRole := 0
	for _, rc := range st.UsersByRole {
		if rc.N > maxRole {
			maxRole = rc.N
		}
	}
	for i := range st.UsersByRole {
		st.UsersByRole[i].Pct = barPct(int64(st.UsersByRole[i].N), int64(maxRole))
	}

	rows, err = s.db.Query(ctx, `SELECT category, COUNT(*), COALESCE(SUM(views_count),0)
		FROM articles WHERE status='published' GROUP BY category ORDER BY 3 DESC`)
	if err != nil {
		return st, err
	}
	for rows.Next() {
		var c CatViews
		if err := rows.Scan(&c.Slug, &c.N, &c.Views); err != nil {
			rows.Close()
			return st, err
		}
		st.ByCat = append(st.ByCat, c)
	}
	rows.Close()
	var maxViews int64
	for _, c := range st.ByCat {
		if c.Views > maxViews {
			maxViews = c.Views
		}
	}
	for i := range st.ByCat {
		st.ByCat[i].Pct = barPct(st.ByCat[i].Views, maxViews)
	}

	rows, err = s.db.Query(ctx, `SELECT a.slug, COALESCE(t.title,a.slug), a.views_count, a.score
		FROM articles a
		LEFT JOIN article_translations t ON t.article_id = a.id AND t.lang = 'ru'
		WHERE a.status='published' ORDER BY a.views_count DESC LIMIT 10`)
	if err != nil {
		return st, err
	}
	for rows.Next() {
		var t TopArticle
		if err := rows.Scan(&t.Slug, &t.Title, &t.Views, &t.Score); err != nil {
			rows.Close()
			return st, err
		}
		st.TopArticles = append(st.TopArticles, t)
	}
	rows.Close()

	rows, err = s.db.Query(ctx, `SELECT c.id, u.email, c.body, a.slug, c.created_at
		FROM comments c JOIN auth_users u ON u.id=c.user_id JOIN articles a ON a.id=c.article_id
		WHERE c.status='published' ORDER BY c.created_at DESC LIMIT 15`)
	if err != nil {
		return st, err
	}
	for rows.Next() {
		var c AdminComment
		var id uuid.UUID
		var email string
		if err := rows.Scan(&id, &email, &c.Body, &c.Slug, &c.CreatedAt); err != nil {
			rows.Close()
			return st, err
		}
		c.ID = id.String()
		c.AuthorName = displayName(email)
		st.RecentComments = append(st.RecentComments, c)
	}
	rows.Close()

	return st, nil
}

// AssignRole sets a user's role by email. Returns the number of rows updated.
func (s *AdminStore) AssignRole(ctx context.Context, email, role string) (int64, error) {
	ct, err := s.db.Exec(ctx, `UPDATE auth_users SET role=$2 WHERE lower(email)=lower($1)`, email, role)
	if err != nil {
		return 0, fmt.Errorf("assign role: %w", err)
	}
	return ct.RowsAffected(), nil
}

// HideComment marks a comment hidden (moderation).
func (s *AdminStore) HideComment(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `UPDATE comments SET status='hidden' WHERE id=$1`, id)
	return err
}

// ---- handlers ----

// AdminPage backs the dashboard template.
type AdminPage struct {
	Base
	Stats          AdminStats
	CanManageUsers bool
	CanFinance     bool
	CanModerate    bool
	AssignRoles    []string
	Notice         string
	Email          string
	Role           string
	// Moderation ledger: the work queue first, then the history.
	Appeals []ModAppeal
	ModLog  []ModAction
	Queue   []ReviewItem // articles awaiting a human decision
	// Growth analytics.
	Analytics AdminAnalytics
	// Operational service switches (maintenance mode per service).
	Services      []ServiceFlag
	ServiceStates []string    // selectable statuses: on | maintenance | off
	Site          ServiceFlag // the global site switch
	// Real-estate agents awaiting verification.
	PendingAgents []Agent
}

func (m *Module) handleAdmin(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	claims, _ := auth.ClaimsFromContext(r.Context())
	stats, err := m.admin.Stats(r.Context())
	if err != nil {
		m.rt.Logger.Error("admin stats", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	page := AdminPage{Base: m.base(r, T(lang, "admin.title"), lang)}
	page.Stats = stats
	if canModerate(claims) {
		if ap, err := m.mods.OpenAppeals(r.Context(), 50); err == nil {
			page.Appeals = ap
		} else {
			m.rt.Logger.Error("open appeals", zap.Error(err))
		}
		if lg, err := m.mods.Recent(r.Context(), 40); err == nil {
			page.ModLog = lg
		} else {
			m.rt.Logger.Error("moderation log", zap.Error(err))
		}
		if q, err := m.store.ReviewQueue(r.Context(), 100); err == nil {
			page.Queue = q
		} else {
			m.rt.Logger.Error("review queue", zap.Error(err))
		}
	}
	if an, err := m.adminAnalytics(r.Context()); err == nil {
		page.Analytics = an
	} else {
		m.rt.Logger.Error("admin analytics", zap.Error(err))
	}
	page.CanManageUsers = canManageUsers(claims)
	page.CanFinance = canViewFinance(claims)
	page.CanModerate = canModerate(claims)
	page.AssignRoles = assignableRoles
	if canManageUsers(claims) {
		page.Services = m.flags.All()
		page.ServiceStates = []string{svcOn, svcMaintenance, svcOff}
		page.Site = m.flags.SiteFlag()
		if pend, err := m.reagents.Pending(r.Context(), 100); err == nil {
			page.PendingAgents = pend
		} else {
			m.rt.Logger.Error("pending agents", zap.Error(err))
		}
	}
	page.Notice = r.URL.Query().Get("ok")
	if claims != nil {
		page.Email = claims.Email
		page.Role = claims.PrimaryRole
	}
	m.render(w, "admin", page)
}

func (m *Module) handleAdminAssignRole(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	_ = r.ParseForm()
	email := strings.TrimSpace(r.FormValue("email"))
	role := strings.TrimSpace(strings.ToLower(r.FormValue("role")))
	if email == "" || !contains(assignableRoles, role) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	n, err := m.admin.AssignRole(r.Context(), email, role)
	if err != nil {
		m.rt.Logger.Error("assign role", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	msg := "role_set"
	if n == 0 {
		msg = "user_missing"
	}
	http.Redirect(w, r, "/admin?ok="+msg, http.StatusSeeOther)
}

// handleAdminServiceFlag flips a service on / into maintenance / off and stores
// its localized notice. This is how a service is taken down for maintenance
// (or held back during the beta) without a redeploy.
func (m *Module) handleAdminServiceFlag(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	_ = r.ParseForm()
	code := strings.TrimSpace(r.FormValue("code"))
	status := strings.TrimSpace(r.FormValue("status"))
	if !validServiceCode(code) || !isServiceStatus(status) {
		http.Redirect(w, r, "/admin?ok=svc_bad", http.StatusSeeOther)
		return
	}
	var by *uuid.UUID
	if claims != nil {
		if id, err := uuid.Parse(claims.Subject); err == nil {
			by = &id
		}
	}
	// Optional maintenance window: minutes from now. 0 / empty = open-ended (no
	// countdown). Only meaningful when the service is not being turned on.
	var until *time.Time
	if status != svcOn {
		if mins, err := strconv.Atoi(strings.TrimSpace(r.FormValue("minutes"))); err == nil && mins > 0 {
			t := time.Now().Add(time.Duration(mins) * time.Minute)
			until = &t
		}
	}
	if err := m.flags.Set(r.Context(), code, status,
		strings.TrimSpace(r.FormValue("message_kz")),
		strings.TrimSpace(r.FormValue("message_ru")),
		strings.TrimSpace(r.FormValue("message_en")), until, by); err != nil {
		m.rt.Logger.Error("set service flag", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin?ok=svc_set", http.StatusSeeOther)
}

// handleAdminAgentDecide approves or rejects a real-estate agent profile. Only
// a verified agent gets the public "Agent" badge and page, so this is the trust
// gate for the whole feature.
func (m *Module) handleAdminAgentDecide(w http.ResponseWriter, r *http.Request) {
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
	status := agentRejected
	if r.FormValue("decision") == "verify" {
		status = agentVerified
	}
	reason := ""
	if status == agentRejected {
		reason = clip(strings.TrimSpace(r.FormValue("reason")), 300)
	}
	var reviewer *uuid.UUID
	if claims != nil {
		if id, perr := uuid.Parse(claims.Subject); perr == nil {
			reviewer = &id
		}
	}
	if err := m.reagents.SetStatus(r.Context(), uid, status, reason, reviewer); err != nil {
		m.rt.Logger.Error("agent decide", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin?ok=agent_set", http.StatusSeeOther)
}

func (m *Module) handleAdminHideComment(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canModerate(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := m.admin.HideComment(r.Context(), id); err != nil {
		m.rt.Logger.Error("hide comment", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	// A hidden comment used to disappear with no trace and no way for its
	// author to learn why. Record the decision so it can be read and contested.
	reason := strings.TrimSpace(r.FormValue("reason"))
	if !isModerationReason(reason) {
		reason = "off_topic"
	}
	var subject *uuid.UUID
	var title string
	if uid, body, err := m.admin.CommentAuthor(r.Context(), id); err == nil {
		subject, title = uid, clip(body, 90)
	}
	mid, _ := uuid.Parse(claims.Subject)
	if _, err := m.mods.Log(r.Context(), ModAction{
		TargetType: "comment", TargetID: id.String(), Title: title,
		Action: "hide", ReasonCode: reason,
		ReasonNote: clip(strings.TrimSpace(r.FormValue("note")), 500),
		ActorKind:  "human",
	}, subject, &mid); err != nil {
		m.rt.Logger.Error("log moderation", zap.Error(err))
	}
	http.Redirect(w, r, "/admin?ok=comment_hidden", http.StatusSeeOther)
}

// barPct scales v against max to a 0..100 width, with a small floor so
// non-zero bars stay visible.
func barPct(v, max int64) int {
	if max <= 0 || v <= 0 {
		return 0
	}
	p := int(v * 100 / max)
	if p < 4 {
		p = 4
	}
	return p
}

func contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}

// CommentAuthor returns who wrote a comment and its text, so a moderation
// entry can name the affected author and quote what was hidden.
func (s *AdminStore) CommentAuthor(ctx context.Context, id uuid.UUID) (*uuid.UUID, string, error) {
	var uid uuid.UUID
	var body string
	if err := s.db.QueryRow(ctx,
		`SELECT user_id, body FROM comments WHERE id = $1`, id).Scan(&uid, &body); err != nil {
		return nil, "", err
	}
	return &uid, body, nil
}
