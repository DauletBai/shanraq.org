package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNormalizeRoleSet(t *testing.T) {
	t.Parallel()
	primary, roles := normalizeRoleSet("Admin", "operator", "admin", "viewer")

	if primary != "admin" {
		t.Fatalf("expected primary role admin, got %s", primary)
	}
	expected := []string{"admin", "operator", "viewer", defaultRoleName}
	if len(roles) != len(expected) {
		t.Fatalf("expected %d roles, got %#v", len(expected), roles)
	}
	for i, want := range expected {
		if roles[i] != want {
			t.Fatalf("role[%d] = %s, want %s", i, roles[i], want)
		}
	}
}

func TestNormalizeRoleSetDefaultsToUser(t *testing.T) {
	t.Parallel()
	primary, roles := normalizeRoleSet("", "")

	if primary != defaultRoleName {
		t.Fatalf("expected default primary role %s, got %s", defaultRoleName, primary)
	}
	if len(roles) != 1 || roles[0] != defaultRoleName {
		t.Fatalf("expected only default role, got %#v", roles)
	}
}

func TestSanitizeRoles(t *testing.T) {
	input := []string{"Admin", "admin", "user", " ", "\tUSER"}
	got := sanitizeRoles(input)

	expected := []string{"admin", "user"}
	if len(got) != len(expected) {
		t.Fatalf("expected %d roles, got %#v", len(expected), got)
	}
	for i, want := range expected {
		if got[i] != want {
			t.Fatalf("role[%d] = %s, want %s", i, got[i], want)
		}
	}
}

func TestIsWeakSecret(t *testing.T) {
	cases := []struct {
		secret string
		weak   bool
	}{
		{"", true},
		{"replace-me-now", true},
		{"super-secret-key", true},
		{"changeme", true},
		{"short", false},
		{"this-is-a-secure-secret-string-with-entropy", false},
	}

	for _, tc := range cases {
		if isWeakSecret(tc.secret) != tc.weak {
			t.Fatalf("isWeakSecret(%q) = %t, want %t", tc.secret, !tc.weak, tc.weak)
		}
	}
}

func TestRequireRolesMiddleware(t *testing.T) {
	tokens := NewTokenService("super-secret", time.Minute)
	adminUser := User{
		ID:    uuid.New(),
		Email: "admin@example.com",
		Role:  "admin",
		Roles: []string{"admin", "operator"},
	}

	token, err := tokens.Generate(adminUser)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	module := &Module{tokens: tokens}
	calls := 0
	protected := module.RequireRoles("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if claims, ok := ClaimsFromContext(r.Context()); !ok || !claims.HasAnyRole("admin") {
			t.Fatalf("expected admin claims in context, got %#v", claims)
		}
		w.WriteHeader(http.StatusTeapot)
	}))

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	protected.ServeHTTP(rec, req)

	if rec.Code != http.StatusTeapot {
		t.Fatalf("expected handler to run with status 418, got %d", rec.Code)
	}
	if calls != 1 {
		t.Fatalf("expected handler to be invoked once, got %d", calls)
	}
}

func TestRequireRolesMiddlewareRejectsInsufficientRole(t *testing.T) {
	tokens := NewTokenService("another-secret", time.Minute)
	user := User{
		ID:    uuid.New(),
		Email: "user@example.com",
		Role:  "user",
		Roles: []string{"user"},
	}
	token, err := tokens.Generate(user)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	module := &Module{tokens: tokens}
	protected := module.RequireRoles("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	protected.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for insufficient role, got %d", rec.Code)
	}
}

func TestRequireRolesMiddlewareMissingToken(t *testing.T) {
	module := &Module{tokens: NewTokenService("secret", time.Minute)}
	protected := module.RequireRoles("user")(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rec := httptest.NewRecorder()
	protected.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing token, got %d", rec.Code)
	}
}
