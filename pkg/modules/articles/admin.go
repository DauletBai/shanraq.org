package articles

import (
	"context"
	"fmt"
	"net/http"
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
}
type CatViews struct {
	Slug  string
	N     int
	Views int64
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
	page.CanManageUsers = canManageUsers(claims)
	page.CanFinance = canViewFinance(claims)
	page.CanModerate = canModerate(claims)
	page.AssignRoles = assignableRoles
	page.Notice = r.URL.Query().Get("ok")
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
	http.Redirect(w, r, "/admin?ok=comment_hidden", http.StatusSeeOther)
}

func contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}
