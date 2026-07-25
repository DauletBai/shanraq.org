package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"shanraq.org/internal/config"
	"shanraq.org/pkg/shanraq"
)

// authTestApp wires the real auth module against the test database and serves
// its /auth JSON API through a chi router. Skipped unless SHANRAQ_TEST_DB names
// a test database.
type authTestApp struct {
	t      *testing.T
	router chi.Router
	pool   *pgxpool.Pool
	mod    *Module
	emails []string
}

func newAuthTestApp(t *testing.T) *authTestApp {
	t.Helper()
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run auth handler tests")
	}
	if !strings.Contains(dsn, "test") {
		t.Fatalf("SHANRAQ_TEST_DB must name a test database; refusing %q", dsn)
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	cfg := config.Config{
		Environment: "test",
		Auth:        config.AuthConfig{TokenSecret: "test-token-secret-that-is-long-enough-1234567890", TokenTTL: time.Hour},
	}
	rt := &shanraq.Runtime{Config: cfg, Logger: zap.NewNop(), DB: pool, Router: chi.NewRouter()}
	m := New()
	if err := m.Init(context.Background(), rt); err != nil {
		t.Fatalf("auth init: %v", err)
	}
	m.Routes(rt.Router)
	app := &authTestApp{t: t, router: rt.Router, pool: pool, mod: m}
	t.Cleanup(func() {
		for _, e := range app.emails {
			_, _ = pool.Exec(context.Background(), `DELETE FROM auth_users WHERE lower(email)=lower($1)`, e)
		}
		pool.Close()
	})
	return app
}

func (a *authTestApp) post(path, jsonBody string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodPost, path, strings.NewReader(jsonBody))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	a.router.ServeHTTP(w, r)
	return w
}

func TestAuthSignupFlow(t *testing.T) {
	app := newAuthTestApp(t)

	// Consent is mandatory.
	if w := app.post("/auth/signup", `{"email":"api1@t.test","password":"Parol12345","consent":false}`); w.Code != http.StatusBadRequest {
		t.Errorf("signup without consent = %d, want 400", w.Code)
	}
	// Invalid email is rejected.
	if w := app.post("/auth/signup", `{"email":"not-an-email","password":"Parol12345","consent":true}`); w.Code == http.StatusOK {
		t.Error("signup with bad email must not succeed")
	}
	// A valid signup succeeds and returns a token.
	app.emails = append(app.emails, "api-ok@t.test")
	w := app.post("/auth/signup", `{"email":"api-ok@t.test","password":"Parol12345","consent":true}`)
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Fatalf("valid signup = %d, want 2xx (%s)", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "token") {
		t.Error("signup response should carry a token")
	}
}

func TestAuthSigninFlow(t *testing.T) {
	app := newAuthTestApp(t)
	app.emails = append(app.emails, "signin@t.test")
	if _, _, err := app.mod.RegisterPassword(context.Background(), "signin@t.test", "Parol12345", "Тест", "Юзер", ""); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	// Correct credentials succeed.
	if w := app.post("/auth/signin", `{"email":"signin@t.test","password":"Parol12345"}`); w.Code != http.StatusOK {
		t.Errorf("valid signin = %d, want 200 (%s)", w.Code, w.Body.String())
	}
	// Wrong password is unauthorized.
	if w := app.post("/auth/signin", `{"email":"signin@t.test","password":"wrong-pass"}`); w.Code != http.StatusUnauthorized {
		t.Errorf("wrong password = %d, want 401", w.Code)
	}
	// Unknown account is unauthorized (no user enumeration).
	if w := app.post("/auth/signin", `{"email":"ghost@t.test","password":"whatever1"}`); w.Code != http.StatusUnauthorized {
		t.Errorf("unknown account = %d, want 401", w.Code)
	}
}

func TestAuthProfileRequiresAuth(t *testing.T) {
	app := newAuthTestApp(t)
	r := httptest.NewRequest(http.MethodGet, "/auth/profile", nil)
	w := httptest.NewRecorder()
	app.router.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("profile without auth = %d, want 401", w.Code)
	}
}
