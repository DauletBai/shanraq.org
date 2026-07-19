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
	// ErrEmailVerificationInvalid indicates the verification token is unknown, used, or expired.
	ErrEmailVerificationInvalid = errors.New("email verification token invalid or expired")
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

// CreateEmailVerification stores a verification token hash for a user.
func (s *Store) CreateEmailVerification(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO email_verification_tokens (user_id, token_hash, expires_at)
		 VALUES ($1,$2,$3)`,
		userID, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("insert email verification: %w", err)
	}
	return nil
}

// ConsumeEmailVerification validates an unused, unexpired token, marks the user
// verified and the token used (in one transaction), and returns the user id.
func (s *Store) ConsumeEmailVerification(ctx context.Context, tokenHash string) (uuid.UUID, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin verify tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var (
		id     uuid.UUID
		userID uuid.UUID
	)
	err = tx.QueryRow(ctx,
		`SELECT id, user_id FROM email_verification_tokens
		 WHERE token_hash = $1 AND used_at IS NULL AND expires_at > NOW()`,
		tokenHash).Scan(&id, &userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, ErrEmailVerificationInvalid
		}
		return uuid.Nil, fmt.Errorf("select email verification: %w", err)
	}
	if _, err = tx.Exec(ctx, `UPDATE email_verification_tokens SET used_at = NOW() WHERE id = $1`, id); err != nil {
		return uuid.Nil, fmt.Errorf("mark verification used: %w", err)
	}
	if _, err = tx.Exec(ctx, `UPDATE auth_users SET email_verified_at = NOW() WHERE id = $1 AND email_verified_at IS NULL`, userID); err != nil {
		return uuid.Nil, fmt.Errorf("mark email verified: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit verify: %w", err)
	}
	return userID, nil
}

// IsEmailVerified reports whether the user's email is verified.
func (s *Store) IsEmailVerified(ctx context.Context, userID uuid.UUID) (bool, error) {
	var verified bool
	err := s.db.QueryRow(ctx,
		`SELECT email_verified_at IS NOT NULL FROM auth_users WHERE id = $1`, userID).Scan(&verified)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("check email verified: %w", err)
	}
	return verified, nil
}

// SetName stores the author's real first and last name (the article byline).
func (s *Store) SetName(ctx context.Context, userID uuid.UUID, first, last string) error {
	_, err := s.db.Exec(ctx,
		`UPDATE auth_users SET first_name = $2, last_name = $3 WHERE id = $1`,
		userID, first, last)
	if err != nil {
		return fmt.Errorf("set name: %w", err)
	}
	return nil
}

// AuthorIdentity returns the user's name and phone-verification status.
func (s *Store) AuthorIdentity(ctx context.Context, userID uuid.UUID) (first, last string, phoneVerified bool, err error) {
	err = s.db.QueryRow(ctx,
		`SELECT COALESCE(first_name,''), COALESCE(last_name,''), phone_verified_at IS NOT NULL
		 FROM auth_users WHERE id = $1`, userID).Scan(&first, &last, &phoneVerified)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", false, nil
		}
		return "", "", false, fmt.Errorf("author identity: %w", err)
	}
	return first, last, phoneVerified, nil
}

// CreatePhoneCode stores a hashed one-time SMS code for a user+phone.
func (s *Store) CreatePhoneCode(ctx context.Context, userID uuid.UUID, phone, codeHash string, expiresAt time.Time) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO phone_verification_codes (user_id, phone, code_hash, expires_at)
		 VALUES ($1,$2,$3,$4)`,
		userID, phone, codeHash, expiresAt)
	if err != nil {
		return fmt.Errorf("create phone code: %w", err)
	}
	return nil
}

// VerifyPhoneCode checks the newest unused code for the user. On a correct,
// unexpired, under-attempt-limit match it marks the code used, stores the phone
// on the account and stamps phone_verified_at — all in one transaction. It
// returns true on success; false (with nil error) on a wrong/expired code.
func (s *Store) VerifyPhoneCode(ctx context.Context, userID uuid.UUID, codeHash string, maxAttempts int) (bool, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("begin phone verify: %w", err)
	}
	defer tx.Rollback(ctx)

	var (
		id       uuid.UUID
		phone    string
		attempts int
		expires  time.Time
		stored   string
	)
	err = tx.QueryRow(ctx,
		`SELECT id, phone, attempts, expires_at, code_hash
		 FROM phone_verification_codes
		 WHERE user_id = $1 AND used_at IS NULL
		 ORDER BY created_at DESC LIMIT 1`, userID).Scan(&id, &phone, &attempts, &expires, &stored)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("select phone code: %w", err)
	}
	if attempts >= maxAttempts || expires.Before(time.Now()) {
		return false, nil
	}
	if stored != codeHash {
		_, _ = tx.Exec(ctx, `UPDATE phone_verification_codes SET attempts = attempts + 1 WHERE id = $1`, id)
		if err = tx.Commit(ctx); err != nil {
			return false, fmt.Errorf("commit attempt: %w", err)
		}
		return false, nil
	}
	if _, err = tx.Exec(ctx, `UPDATE phone_verification_codes SET used_at = NOW() WHERE id = $1`, id); err != nil {
		return false, fmt.Errorf("mark code used: %w", err)
	}
	if _, err = tx.Exec(ctx,
		`UPDATE auth_users SET phone = $2, phone_verified_at = NOW() WHERE id = $1`, userID, phone); err != nil {
		return false, fmt.Errorf("mark phone verified: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("commit phone verify: %w", err)
	}
	return true, nil
}

// HasConsent reports whether the user accepted the given document at the given
// version. Checking the version (not just the document) means a materially
// changed document requires fresh consent.
func (s *Store) HasConsent(ctx context.Context, userID uuid.UUID, document, version string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM user_consents WHERE user_id=$1 AND document=$2 AND version=$3)`,
		userID, document, version).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("has consent: %w", err)
	}
	return exists, nil
}

// InsertConsent appends one consent record (append-only history).
func (s *Store) InsertConsent(ctx context.Context, userID uuid.UUID, document, version, source, ip string) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO user_consents (user_id, document, version, source, ip)
		 VALUES ($1,$2,$3,$4,$5)`,
		userID, document, version, source, ip)
	if err != nil {
		return fmt.Errorf("insert consent: %w", err)
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
