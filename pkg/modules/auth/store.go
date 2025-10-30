package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// ErrEmailExists indicates the email is already registered.
	ErrEmailExists = errors.New("email already registered")
	// ErrNotFound signals the requested user does not exist.
	ErrNotFound = errors.New("user not found")
)

// User represents a record stored in auth_users.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

// Store contains persistence helpers for auth module.
type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(ctx context.Context, email, hash string) (User, error) {
	userID := uuid.New()
	_, err := s.db.Exec(ctx, `
		INSERT INTO auth_users (id, email, password_hash)
		VALUES ($1, $2, $3)
	`, userID, email, hash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "auth_users_email_key" {
			return User{}, ErrEmailExists
		}
		return User{}, fmt.Errorf("insert user: %w", err)
	}

	return User{
		ID:           userID,
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}, nil
}

func (s *Store) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := s.db.QueryRow(ctx, `
		SELECT id, email, password_hash, created_at
		FROM auth_users
		WHERE email = $1
	`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("select by email: %w", err)
	}
	return u, nil
}

func (s *Store) GetByID(ctx context.Context, id string) (User, error) {
	var u User
	err := s.db.QueryRow(ctx, `
		SELECT id, email, password_hash, created_at
		FROM auth_users
		WHERE id = $1
	`, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("select by id: %w", err)
	}
	return u, nil
}
