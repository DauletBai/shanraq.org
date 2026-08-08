package articles

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"shanraq.org/internal/config"
	"shanraq.org/pkg/modules/ai"
	"shanraq.org/pkg/modules/auth"
	"shanraq.org/pkg/modules/media"
	"shanraq.org/pkg/modules/notifier"
	"shanraq.org/pkg/modules/syndicate"
	"shanraq.org/pkg/shanraq"
)

// testApp wires the real auth + articles modules against the test database and
// serves them through a chi router, so handlers are exercised over genuine HTTP
// requests. It is the harness for the handler-level coverage. Skipped unless
// SHANRAQ_TEST_DB names a test database.
type testApp struct {
	t      *testing.T
	router chi.Router
	pool   *pgxpool.Pool
	auth   *auth.Module
	origin string
	emails []string // created accounts, deleted on cleanup (cascade)
}

const testOrigin = "http://localhost:8080"

// newTestApp builds the app. Extra auth options let a test stand the module up
// in a configuration production does not currently use — MFA is the case that
// matters, since its whole risk is what happens the day it is switched on.
func newTestApp(t *testing.T, authOpts ...auth.Option) *testApp {
	t.Helper()
	dsn := requireTestDB(t)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}

	cfg := config.Config{
		Environment:   "test",
		PublicBaseURL: testOrigin,
		Auth: config.AuthConfig{
			TokenSecret: "test-token-secret-that-is-long-enough-1234567890",
			TokenTTL:    time.Hour,
		},
		Media: config.MediaConfig{
			Backend: "fs", Dir: t.TempDir(), PublicPrefix: "media",
			MaxDimension: 1600, MaxUploadBytes: 10 << 20,
		},
	}
	rt := &shanraq.Runtime{Config: cfg, Logger: zap.NewNop(), DB: pool, Router: chi.NewRouter()}

	mailer := notifier.New()
	authM := auth.New(append([]auth.Option{auth.WithMailer(mailer)}, authOpts...)...)
	aiM := ai.New()
	synM := syndicate.New(mailer)
	mediaM := media.New(authM)
	arts := New(authM, aiM, synM, mediaM, mailer)

	for _, m := range []shanraq.InitializerModule{authM, aiM, synM, mediaM, arts} {
		if err := m.Init(ctx, rt); err != nil {
			t.Fatalf("%s init: %v", m.Name(), err)
		}
	}
	authM.Routes(rt.Router)
	mediaM.Routes(rt.Router)
	arts.Routes(rt.Router)

	app := &testApp{t: t, router: rt.Router, pool: pool, auth: authM, origin: testOrigin}
	t.Cleanup(app.cleanup)
	return app
}

func (a *testApp) cleanup() {
	ctx := context.Background()
	for _, e := range a.emails {
		_, _ = a.pool.Exec(ctx, `DELETE FROM auth_users WHERE lower(email) = lower($1)`, e)
	}
	a.pool.Close()
}

// ---- request helpers ----

type reqOpt func(*http.Request)

func withCookie(c *http.Cookie) reqOpt {
	return func(r *http.Request) {
		if c != nil {
			r.AddCookie(c)
		}
	}
}

func withHeader(k, v string) reqOpt {
	return func(r *http.Request) { r.Header.Set(k, v) }
}

// do issues a request through the router and returns the recorder. A non-nil
// form is sent as urlencoded body; GET requests pass form as the query.
func (a *testApp) do(method, path string, form url.Values, opts ...reqOpt) *httptest.ResponseRecorder {
	a.t.Helper()
	var r *http.Request
	if form != nil && method != http.MethodGet {
		r = httptest.NewRequest(method, path, strings.NewReader(form.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		if form != nil {
			if strings.Contains(path, "?") {
				path += "&" + form.Encode()
			} else {
				path += "?" + form.Encode()
			}
		}
		r = httptest.NewRequest(method, path, nil)
	}
	// Same-origin so the CSRF guard on browser POSTs passes: the request Host
	// must equal the Origin's host (httptest defaults Host to example.com).
	r.Host = "localhost:8080"
	r.Header.Set("Origin", a.origin)
	for _, o := range opts {
		o(r)
	}
	w := httptest.NewRecorder()
	a.router.ServeHTTP(w, r)
	return w
}

// ---- account helpers ----

// createUser registers a verified-name account directly through the auth module
// (bypassing the studio gate, which is what setup wants) and returns its id.
func (a *testApp) createUser(email, password string) uuid.UUID {
	a.t.Helper()
	a.emails = append(a.emails, email)
	user, _, err := a.auth.RegisterPassword(context.Background(), email, password, "Тест", "Пользователь", "")
	if err != nil {
		a.t.Fatalf("createUser(%s): %v", email, err)
	}
	return user.ID
}

// login performs the studio login flow and returns the session cookie.
func (a *testApp) login(email, password string) *http.Cookie {
	a.t.Helper()
	w := a.do(http.MethodPost, "/studio/login", url.Values{"email": {email}, "password": {password}})
	if w.Code != http.StatusSeeOther {
		a.t.Fatalf("login want 303, got %d (%s)", w.Code, w.Body.String())
	}
	for _, c := range w.Result().Cookies() {
		if c.Name == auth.SessionCookieName {
			return c
		}
	}
	a.t.Fatal("login did not set a session cookie")
	return nil
}

// makeStaff promotes an account to admin (for admin-panel handlers).
func (a *testApp) makeStaff(email, role string) {
	a.t.Helper()
	if _, err := auth.NewStore(a.pool).SetPrimaryRole(context.Background(), email, role); err != nil {
		a.t.Fatalf("makeStaff: %v", err)
	}
}

// exec runs a raw statement against the test DB (for arranging fixtures).
func (a *testApp) exec(sql string, args ...any) {
	a.t.Helper()
	if _, err := a.pool.Exec(context.Background(), sql, args...); err != nil {
		a.t.Fatalf("exec %q: %v", sql, err)
	}
}

// seedArticle inserts an article by author with a Russian translation, at the
// given status (empty = draft). Returns its id and slug.
func (a *testApp) seedArticle(authorID uuid.UUID, status string) (uuid.UUID, string) {
	a.t.Helper()
	slug := "t-" + uuid.NewString()[:8]
	id, err := NewStore(a.pool).Create(context.Background(), authorID, slug, LangRU, "economy", "", "",
		[]TranslationInput{{Lang: LangRU, Title: "Тест заголовок", Summary: "Саммари", BodyMD: "## Тело\n\nтекст статьи", Source: "human"}})
	if err != nil {
		a.t.Fatalf("seedArticle: %v", err)
	}
	if status != "" && status != "draft" {
		a.exec(`UPDATE articles SET status = $2 WHERE id = $1`, id, status)
	}
	return id, slug
}

// seedListing inserts a published listing owned by author and returns its id.
func (a *testApp) seedListing(authorID uuid.UUID) uuid.UUID {
	a.t.Helper()
	id, err := NewListingStore(a.pool).Create(context.Background(), authorID, ListingInput{
		DealType: "sale", PropertyType: "apartment", Price: 12000000, Rooms: 2, Area: 60,
		Title: "Тестовая квартира", Description: "Описание", Contact: "+7 700 000 00 00",
		Country: "Казахстан", Region: "Алматы", City: "Алматы",
	})
	if err != nil {
		a.t.Fatalf("seedListing: %v", err)
	}
	return id
}
