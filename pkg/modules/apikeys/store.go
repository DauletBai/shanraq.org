package apikeys

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"shanraq.org/pkg/modules/auth"
)

type pool interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type poolWrapper struct {
	*pgxpool.Pool
}

func (p poolWrapper) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return p.Pool.Exec(ctx, sql, args...)
}

func (p poolWrapper) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.Pool.Query(ctx, sql, args...)
}

func (p poolWrapper) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.Pool.QueryRow(ctx, sql, args...)
}

type authLookup interface {
	GetByID(ctx context.Context, id string) (auth.User, error)
}

// Store coordinates persistence of API keys and related lookups.
type Store struct {
	db        pool
	authStore authLookup
}

func NewStore(db *pgxpool.Pool) *Store {
	return newStoreWithDeps(poolWrapper{db}, auth.NewStore(db))
}

func newStoreWithDeps(db pool, authStore authLookup) *Store {
	return &Store{db: db, authStore: authStore}
}

// Create issues a new API key for the given user, returning the plaintext secret once.
func (s *Store) Create(ctx context.Context, userID uuid.UUID, label string) (string, APIKey, error) {
	secret, prefix, hash, err := generateKey()
	if err != nil {
		return "", APIKey{}, err
	}

	var key APIKey
	err = s.db.QueryRow(ctx, `
		INSERT INTO auth_api_keys (user_id, key_hash, prefix, label)
		VALUES ($1, $2, $3, NULLIF($4, ''))
		RETURNING id, prefix, COALESCE(label, ''), created_at, revoked_at
	`, userID, hash, prefix, label).Scan(&key.ID, &key.Prefix, &key.Label, &key.CreatedAt, &key.RevokedAt)
	if err != nil {
		return "", APIKey{}, fmt.Errorf("insert api key: %w", err)
	}
	if key.Label == "" {
		key.Label = label
	}
	return secret, key, nil
}

// List returns non-revoked and revoked keys for the given user ordered by creation time.
func (s *Store) List(ctx context.Context, userID uuid.UUID) ([]APIKey, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, prefix, COALESCE(label, ''), created_at, revoked_at
		FROM auth_api_keys
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list api keys: %w", err)
	}
	defer rows.Close()

	var keys []APIKey
	for rows.Next() {
		var key APIKey
		if err := rows.Scan(&key.ID, &key.Prefix, &key.Label, &key.CreatedAt, &key.RevokedAt); err != nil {
			return nil, fmt.Errorf("scan api key: %w", err)
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

// Revoke marks the specified key as revoked if it belongs to the provided user.
func (s *Store) Revoke(ctx context.Context, userID, keyID uuid.UUID) error {
	tag, err := s.db.Exec(ctx, `
		UPDATE auth_api_keys
		SET revoked_at = NOW()
		WHERE id = $1 AND user_id = $2 AND revoked_at IS NULL
	`, keyID, userID)
	if err != nil {
		return fmt.Errorf("revoke api key: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// Validate returns the owning user for the provided secret if the key is active.
func (s *Store) Validate(ctx context.Context, token string) (auth.User, APIKey, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return auth.User{}, APIKey{}, fmt.Errorf("api key required")
	}
	hash := hashKey(token)

	var key APIKey
	var userID uuid.UUID
	err := s.db.QueryRow(ctx, `
		SELECT id, user_id, prefix, COALESCE(label, ''), created_at, revoked_at
		FROM auth_api_keys
		WHERE key_hash = $1
	`, hash).Scan(&key.ID, &userID, &key.Prefix, &key.Label, &key.CreatedAt, &key.RevokedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.User{}, APIKey{}, fmt.Errorf("api key not found")
		}
		return auth.User{}, APIKey{}, fmt.Errorf("lookup api key: %w", err)
	}
	if key.RevokedAt != nil {
		return auth.User{}, APIKey{}, fmt.Errorf("api key revoked")
	}

	user, err := s.authStore.GetByID(ctx, userID.String())
	if err != nil {
		return auth.User{}, APIKey{}, err
	}
	return user, key, nil
}

func generateKey() (plain string, prefix string, hash string, err error) {
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", "", fmt.Errorf("generate api key entropy: %w", err)
	}
	encoded := base64.RawURLEncoding.EncodeToString(buf)
	plain = "sk_" + encoded
	if len(plain) < 10 {
		return "", "", "", fmt.Errorf("api key generation produced short key")
	}
	prefix = plain[:10]
	hash = hashKey(plain)
	return plain, prefix, hash, nil
}

func hashKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", sum[:])
}
