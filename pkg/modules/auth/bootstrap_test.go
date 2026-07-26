package auth

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"shanraq.org/internal/config"
	"shanraq.org/pkg/shanraq"
)

func bootstrapTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run this integration test")
	}
	if !strings.Contains(dsn, "test") {
		t.Fatalf("SHANRAQ_TEST_DB must name a test database; refusing %q", dsn)
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func TestCreateVerifiedAdmin(t *testing.T) {
	pool := bootstrapTestPool(t)
	s := NewStore(pool)
	ctx := context.Background()
	email := "boot-" + uuid.NewString()[:8] + "@example.com"
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE email = $1`, email) })

	id, err := s.CreateVerifiedAdmin(ctx, email, "$2a$10$abcdefghijklmnopqrstuv", "Test", "Admin")
	if err != nil {
		t.Fatalf("CreateVerifiedAdmin: %v", err)
	}
	var role string
	var verified bool
	if err := pool.QueryRow(ctx,
		`SELECT role, email_verified_at IS NOT NULL FROM auth_users WHERE id = $1`, id).Scan(&role, &verified); err != nil {
		t.Fatalf("read back: %v", err)
	}
	if role != "admin" {
		t.Errorf("role = %q, want admin", role)
	}
	if !verified {
		t.Error("bootstrap admin email should be pre-verified")
	}

	if has, err := s.HasAnyStaffAdmin(ctx); err != nil || !has {
		t.Errorf("HasAnyStaffAdmin = %v (err %v), want true", has, err)
	}
}

// With an admin already present, the bootstrap must NOT create a second account.
func TestEnsureBootstrapAdminSkipsWhenAdminExists(t *testing.T) {
	pool := bootstrapTestPool(t)
	ctx := context.Background()
	s := NewStore(pool)

	seedEmail := "seed-admin-" + uuid.NewString()[:8] + "@example.com"
	if _, err := s.CreateVerifiedAdmin(ctx, seedEmail, "$2a$10$abcdefghijklmnopqrstuv", "Seed", "Admin"); err != nil {
		t.Fatalf("seed admin: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE email = $1`, seedEmail) })

	bootEmail := "wont-create-" + uuid.NewString()[:8] + "@example.com"
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE email = $1`, bootEmail) })
	m := &Module{
		rt: &shanraq.Runtime{
			Config: config.Config{Bootstrap: config.BootstrapConfig{AdminEmail: bootEmail, AdminPassword: "strong-enough"}},
			Logger: zap.NewNop(), DB: pool, Router: chi.NewRouter(),
		},
		store: s,
	}
	m.ensureBootstrapAdmin(ctx)

	var n int
	_ = pool.QueryRow(ctx, `SELECT count(*) FROM auth_users WHERE email = $1`, bootEmail).Scan(&n)
	if n != 0 {
		t.Errorf("bootstrap must skip when an admin already exists; created %d row(s) for %s", n, bootEmail)
	}
}
