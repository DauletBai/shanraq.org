package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultRoleName = "user"

var (
	// ErrEmailExists indicates the email is already registered.
	ErrEmailExists = errors.New("email already registered")
	// ErrNotFound signals the requested user does not exist.
	ErrNotFound = errors.New("user not found")
	// ErrRefreshNotFound indicates the refresh token does not exist.
	ErrRefreshNotFound = errors.New("refresh token not found")
	// ErrPasswordResetNotFound indicates the password reset token does not exist.
	ErrPasswordResetNotFound = errors.New("password reset token not found")
	// ErrMFATOTPNotFound indicates the user does not yet have a TOTP secret registered.
	ErrMFATOTPNotFound = errors.New("mfa totp secret not found")
)

// User represents a record stored in auth_users.
type User struct {
	ID                    uuid.UUID
	Email                 string
	PasswordHash          string
	Role                  string
	Roles                 []string
	PasswordResetRequired bool
	CreatedAt             time.Time
}

// MFATOTP persists user-specific TOTP secrets.
type MFATOTP struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Secret      string
	Confirmed   bool
	ConfirmedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	URI         string
}

// Store contains persistence helpers for auth module.
type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(ctx context.Context, email, hash string, roles ...string) (User, error) {
	var primaryInput string
	var extras []string
	if len(roles) > 0 {
		primaryInput = roles[0]
		if len(roles) > 1 {
			extras = roles[1:]
		}
	}

	primaryRole, normalizedRoles := normalizeRoleSet(primaryInput, extras...)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return User{}, fmt.Errorf("begin user tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	userID := uuid.New()
	var createdAt time.Time
	err = tx.QueryRow(ctx, `
		INSERT INTO auth_users (id, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`, userID, email, hash, primaryRole).Scan(&createdAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "auth_users_email_key" {
			return User{}, ErrEmailExists
		}
		return User{}, fmt.Errorf("insert user: %w", err)
	}

	roleIDs := make(map[string]uuid.UUID, len(normalizedRoles))
	for _, roleName := range normalizedRoles {
		roleID, ensureErr := s.ensureRoleTx(ctx, tx, roleName, "")
		if ensureErr != nil {
			return User{}, ensureErr
		}
		roleIDs[roleName] = roleID
	}

	for _, roleName := range normalizedRoles {
		roleID := roleIDs[roleName]
		if assignErr := s.assignRoleTx(ctx, tx, userID, roleID); assignErr != nil {
			return User{}, assignErr
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return User{}, fmt.Errorf("commit user: %w", err)
	}

	return User{
		ID:           userID,
		Email:        email,
		PasswordHash: hash,
		Role:         primaryRole,
		Roles:        normalizedRoles,
		CreatedAt:    createdAt,
	}, nil
}

func (s *Store) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := s.db.QueryRow(ctx, `
		SELECT id, email, password_hash, role, password_reset_required, created_at
		FROM auth_users
		WHERE email = $1
	`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.PasswordResetRequired, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("select by email: %w", err)
	}

	if err := s.hydrateUserRoles(ctx, &u); err != nil {
		return User{}, err
	}
	return u, nil
}

func (s *Store) GetByID(ctx context.Context, id string) (User, error) {
	var u User
	err := s.db.QueryRow(ctx, `
		SELECT id, email, password_hash, role, password_reset_required, created_at
		FROM auth_users
		WHERE id = $1
	`, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.PasswordResetRequired, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("select by id: %w", err)
	}

	if err := s.hydrateUserRoles(ctx, &u); err != nil {
		return User{}, err
	}
	return u, nil
}

// RefreshToken represents a persisted refresh token.
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time
}

func (s *Store) InsertRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) (RefreshToken, error) {
	var rt RefreshToken
	err := s.db.QueryRow(ctx, `
		INSERT INTO auth_refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, token_hash, expires_at, created_at, revoked_at
	`, userID, tokenHash, expiresAt).Scan(&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.CreatedAt, &rt.RevokedAt)
	if err != nil {
		return RefreshToken{}, fmt.Errorf("insert refresh token: %w", err)
	}
	return rt, nil
}

func (s *Store) GetRefreshToken(ctx context.Context, tokenHash string) (RefreshToken, error) {
	var rt RefreshToken
	err := s.db.QueryRow(ctx, `
		SELECT id, user_id, token_hash, expires_at, created_at, revoked_at
		FROM auth_refresh_tokens
		WHERE token_hash = $1
	`, tokenHash).Scan(&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.CreatedAt, &rt.RevokedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RefreshToken{}, ErrRefreshNotFound
		}
		return RefreshToken{}, fmt.Errorf("select refresh token: %w", err)
	}
	return rt, nil
}

func (s *Store) RevokeRefreshToken(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `
		UPDATE auth_refresh_tokens
		SET revoked_at = NOW()
		WHERE id = $1 AND revoked_at IS NULL
	`, id)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

func (s *Store) RevokeUserTokens(ctx context.Context, userID uuid.UUID) error {
	_, err := s.db.Exec(ctx, `
		UPDATE auth_refresh_tokens
		SET revoked_at = NOW()
		WHERE user_id = $1 AND revoked_at IS NULL
	`, userID)
	if err != nil {
		return fmt.Errorf("revoke user tokens: %w", err)
	}
	return nil
}

func (s *Store) TrimActiveRefreshTokens(ctx context.Context, userID uuid.UUID, keep int) error {
	if keep <= 0 {
		return nil
	}
	_, err := s.db.Exec(ctx, `
		DELETE FROM auth_refresh_tokens
		WHERE id IN (
			SELECT id
			FROM auth_refresh_tokens
			WHERE user_id = $1 AND revoked_at IS NULL
			ORDER BY created_at DESC
			OFFSET $2
		)
	`, userID, keep)
	if err != nil {
		return fmt.Errorf("trim refresh tokens: %w", err)
	}
	return nil
}

func (s *Store) DeleteExpiredRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	_, err := s.db.Exec(ctx, `
		DELETE FROM auth_refresh_tokens
		WHERE user_id = $1 AND expires_at < NOW()
	`, userID)
	if err != nil {
		return fmt.Errorf("delete expired refresh tokens: %w", err)
	}
	return nil
}

// PasswordReset represents a pending password reset token.
type PasswordReset struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	UsedAt    *time.Time
}

func (s *Store) CreatePasswordReset(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) (PasswordReset, error) {
	var pr PasswordReset
	err := s.db.QueryRow(ctx, `
		INSERT INTO auth_password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, token_hash, expires_at, created_at, used_at
	`, userID, tokenHash, expiresAt).Scan(&pr.ID, &pr.UserID, &pr.TokenHash, &pr.ExpiresAt, &pr.CreatedAt, &pr.UsedAt)
	if err != nil {
		return PasswordReset{}, fmt.Errorf("insert password reset: %w", err)
	}
	return pr, nil
}

func (s *Store) GetPasswordReset(ctx context.Context, tokenHash string) (PasswordReset, error) {
	var pr PasswordReset
	err := s.db.QueryRow(ctx, `
		SELECT id, user_id, token_hash, expires_at, created_at, used_at
		FROM auth_password_resets
		WHERE token_hash = $1
	`, tokenHash).Scan(&pr.ID, &pr.UserID, &pr.TokenHash, &pr.ExpiresAt, &pr.CreatedAt, &pr.UsedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PasswordReset{}, ErrPasswordResetNotFound
		}
		return PasswordReset{}, fmt.Errorf("select password reset: %w", err)
	}
	return pr, nil
}

func (s *Store) MarkPasswordResetUsed(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `
		UPDATE auth_password_resets
		SET used_at = NOW()
		WHERE id = $1 AND used_at IS NULL
	`, id)
	if err != nil {
		return fmt.Errorf("mark password reset used: %w", err)
	}
	return nil
}

func (s *Store) UpdatePassword(ctx context.Context, userID uuid.UUID, hash string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE auth_users
		SET password_hash = $2,
		    password_reset_required = FALSE,
		    updated_at = NOW()
		WHERE id = $1
	`, userID, hash)
	if err != nil {
		return fmt.Errorf("update user password: %w", err)
	}
	return nil
}

func (s *Store) ensureRoleTx(ctx context.Context, tx pgx.Tx, name, description string) (uuid.UUID, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		name = defaultRoleName
	}

	var roleID uuid.UUID
	err := tx.QueryRow(ctx, `
		SELECT id
		FROM auth_roles
		WHERE name = $1
	`, name).Scan(&roleID)
	if err == nil {
		return roleID, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("lookup role %s: %w", name, err)
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO auth_roles (name, description)
		VALUES ($1, NULLIF($2, ''))
		RETURNING id
	`, name, description).Scan(&roleID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert role %s: %w", name, err)
	}
	return roleID, nil
}

func (s *Store) assignRoleTx(ctx context.Context, tx pgx.Tx, userID, roleID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO auth_user_roles (user_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`, userID, roleID)
	if err != nil {
		return fmt.Errorf("assign role: %w", err)
	}
	return nil
}

func (s *Store) hydrateUserRoles(ctx context.Context, u *User) error {
	roles, err := s.listUserRoles(ctx, u.ID)
	if err != nil {
		return err
	}

	primary, normalized := normalizeRoleSet(u.Role, roles...)
	u.Role = primary
	u.Roles = normalized
	return nil
}

func (s *Store) listUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	rows, err := s.db.Query(ctx, `
		SELECT r.name
		FROM auth_roles r
		INNER JOIN auth_user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
		ORDER BY r.name
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("query user roles: %w", err)
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan user role: %w", err)
		}
		roles = append(roles, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user roles: %w", err)
	}
	return roles, nil
}

func (s *Store) CreateTOTP(ctx context.Context, userID uuid.UUID, secret string) (MFATOTP, error) {
	id := uuid.New()
	if _, err := s.db.Exec(ctx, `
		INSERT INTO auth_mfa_totp (id, user_id, secret)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO NOTHING
	`, id, userID, secret); err != nil {
		return MFATOTP{}, fmt.Errorf("insert totp secret: %w", err)
	}
	return s.GetTOTPByUser(ctx, userID)
}

func (s *Store) GetTOTPByUser(ctx context.Context, userID uuid.UUID) (MFATOTP, error) {
	return scanMFATOTP(s.db.QueryRow(ctx, `
		SELECT id, user_id, secret, confirmed_at, created_at, updated_at
		FROM auth_mfa_totp
		WHERE user_id = $1
	`, userID))
}

func (s *Store) GetTOTPByID(ctx context.Context, id uuid.UUID) (MFATOTP, error) {
	return scanMFATOTP(s.db.QueryRow(ctx, `
		SELECT id, user_id, secret, confirmed_at, created_at, updated_at
		FROM auth_mfa_totp
		WHERE id = $1
	`, id))
}

func (s *Store) MarkTOTPConfirmed(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `
		UPDATE auth_mfa_totp
		SET confirmed_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("mark totp confirmed: %w", err)
	}
	return nil
}

func scanMFATOTP(row pgx.Row) (MFATOTP, error) {
	var (
		record      MFATOTP
		confirmedAt sql.NullTime
	)
	if err := row.Scan(&record.ID, &record.UserID, &record.Secret, &confirmedAt, &record.CreatedAt, &record.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return MFATOTP{}, ErrMFATOTPNotFound
		}
		return MFATOTP{}, fmt.Errorf("scan totp: %w", err)
	}
	if confirmedAt.Valid {
		record.Confirmed = true
		record.ConfirmedAt = &confirmedAt.Time
	}
	record.URI = ""
	return record, nil
}

func normalizeRoleSet(primary string, extras ...string) (string, []string) {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(extras)+2)

	add := func(role string) {
		role = strings.TrimSpace(strings.ToLower(role))
		if role == "" {
			return
		}
		if _, ok := seen[role]; ok {
			return
		}
		seen[role] = struct{}{}
		result = append(result, role)
	}

	add(primary)
	for _, role := range extras {
		add(role)
	}
	if _, ok := seen[defaultRoleName]; !ok {
		add(defaultRoleName)
	}
	if len(result) == 0 {
		add(defaultRoleName)
	}
	return result[0], result
}
