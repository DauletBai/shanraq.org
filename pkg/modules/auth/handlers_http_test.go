package auth

import (
	"context"
	"errors"
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
	w := app.post("/auth/signup", `{"email":"api-ok@t.test","password":"Parol12345","consent":true,"first_name":"Даулет","last_name":"Баймурза"}`)
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

// The JSON signup used to be the one door with no lock on it: the browser form
// honours the registration service flag, while POST /auth/signup created an
// account and handed back tokens regardless. Closing registration in the admin
// panel therefore shut the site to visitors and left it open to anyone who
// could spell JSON.
func TestSignupObeysGate(t *testing.T) {
	app := newAuthTestApp(t)

	blocked := errors.New("registration is currently closed")
	app.mod.signupGate = func() error { return blocked }

	w := app.post("/auth/signup", `{"email":"gated@t.test","password":"Parol12345","consent":true,"first_name":"Даулет","last_name":"Баймурза"}`)
	if w.Code != http.StatusForbidden {
		t.Fatalf("signup with a closed gate = %d, want 403 (%s)", w.Code, w.Body.String())
	}

	// And the account must not exist: a refusal that still writes a row is worse
	// than no refusal at all, because it hides the breach.
	var n int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM auth_users WHERE lower(email)='gated@t.test'`).Scan(&n); err != nil {
		t.Fatalf("count: %v", err)
	}
	if n != 0 {
		t.Errorf("a blocked signup created %d user(s), want 0", n)
	}

	// With the gate open the endpoint works as before.
	app.mod.signupGate = func() error { return nil }
	app.emails = append(app.emails, "ungated@t.test")
	if w := app.post("/auth/signup", `{"email":"ungated@t.test","password":"Parol12345","consent":true,"first_name":"Даулет","last_name":"Баймурза"}`); w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Errorf("signup with an open gate = %d, want 2xx (%s)", w.Code, w.Body.String())
	}
}

// The browser form demands a real first and last name because every article,
// listing and comment on the site is signed by a person. The JSON endpoint used
// to skip that check, quietly making the rule optional for anyone who preferred
// JSON to the form.
//
// Signup is rate limited to a burst of three, so this stays at two requests and
// proves only that the endpoint enforces the rule and stores the result. The
// full table of what counts as a name lives in TestValidatePersonNameRejectsNonNames.
func TestSignupRequiresRealName(t *testing.T) {
	app := newAuthTestApp(t)

	if w := app.post("/auth/signup", `{"email":"noname@t.test","password":"Parol12345","consent":true}`); w.Code != http.StatusBadRequest {
		t.Errorf("a nameless signup = %d, want 400 (%s)", w.Code, w.Body.String())
	}
	var n int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM auth_users WHERE lower(email)='noname@t.test'`).Scan(&n); err != nil {
		t.Fatalf("count: %v", err)
	}
	if n != 0 {
		t.Errorf("a rejected signup created %d account(s), want 0", n)
	}

	// A named payload works, and the name actually lands in the row — a rule that
	// only rejects would pass the check above while storing nothing useful.
	app.emails = append(app.emails, "named@t.test")
	w := app.post("/auth/signup", `{"email":"named@t.test","password":"Parol12345","consent":true,"first_name":"Даулет","last_name":"Баймурза","middle_name":"Абаевич"}`)
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Fatalf("a named signup = %d, want 2xx (%s)", w.Code, w.Body.String())
	}
	var first, last string
	if err := app.pool.QueryRow(context.Background(),
		`SELECT COALESCE(first_name,''), COALESCE(last_name,'') FROM auth_users WHERE lower(email)='named@t.test'`).
		Scan(&first, &last); err != nil {
		t.Fatalf("read back: %v", err)
	}
	if first != "Даулет" || last != "Баймурза" {
		t.Errorf("name stored as %q %q, want Даулет Баймурза", first, last)
	}
}
