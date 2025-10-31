package apikeys

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	pgxmock "github.com/pashagolub/pgxmock/v4"
	"shanraq.org/pkg/modules/auth"
)

func TestStoreCreate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock: %v", err)
	}
	defer mock.Close()

	userID := uuid.New()
	keyID := uuid.New()
	createdAt := time.Now()

	rows := pgxmock.NewRows([]string{"id", "prefix", "label", "created_at", "revoked_at"}).
		AddRow(keyID, "sk_prefix", "integration", createdAt, nil)

	mock.ExpectQuery("INSERT INTO auth_api_keys").
		WithArgs(userID, pgxmock.AnyArg(), pgxmock.AnyArg(), "integration").
		WillReturnRows(rows)

	store := newStoreWithDeps(mock, stubAuthStore{})
	secret, key, err := store.Create(context.Background(), userID, "integration")
	if err != nil {
		t.Fatalf("create api key: %v", err)
	}
	if secret == "" || key.ID != keyID || key.Prefix == "" {
		t.Fatalf("unexpected key data: secret=%q key=%+v", secret, key)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestStoreValidate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock: %v", err)
	}
	defer mock.Close()

	plain, prefix, hash, err := generateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	keyID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now()

	keyRows := pgxmock.NewRows([]string{"id", "user_id", "prefix", "label", "created_at", "revoked_at"}).
		AddRow(keyID, userID, prefix, "ci", createdAt, nil)

	mock.ExpectQuery("SELECT id, user_id, prefix, COALESCE\\(label").
		WithArgs(hash).
		WillReturnRows(keyRows)

	authUser := auth.User{ID: userID, Email: "owner@example.com", Role: "admin", Roles: []string{"admin"}}
	store := newStoreWithDeps(mock, stubAuthStore{user: authUser})
	user, key, err := store.Validate(context.Background(), plain)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if user.ID != userID {
		t.Fatalf("unexpected user id: %s", user.ID)
	}
	if key.ID != keyID {
		t.Fatalf("unexpected key id: %s", key.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

type stubAuthStore struct {
	user auth.User
	err  error
}

func (s stubAuthStore) GetByID(context.Context, string) (auth.User, error) {
	if s.err != nil {
		return auth.User{}, s.err
	}
	return s.user, nil
}
