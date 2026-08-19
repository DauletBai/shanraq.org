package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// AdminUser is one account as the administrator sees it: who they are, how to
// reach them, what they have posted. Until this existed the panel could report
// how many accounts held each role and nothing else — not a name, not an
// address — which is no basis for acting on a complaint about a person.
type AdminUser struct {
	ID         uuid.UUID
	Email      string
	First      string
	Last       string
	Middle     string
	Role       string
	Country    string // ISO code captured at registration; empty if unknown
	Phone      string
	Verified   bool
	CreatedAt  time.Time
	Articles   int
	Listings   int
	Comments   int
	IsLastRoot bool // the only remaining admin/director — must not be removed
}

// FullName renders "Фамилия Имя Отчество", skipping what is missing.
func (u AdminUser) FullName() string {
	parts := []string{}
	for _, p := range []string{u.Last, u.First, u.Middle} {
		if p = strings.TrimSpace(p); p != "" {
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, " ")
}

// rootRoles are the roles that can administer other accounts. The last holder
// of one cannot be deleted or demoted, or the site locks everybody out of its
// own admin panel with no way back in short of a database console.
var rootRoles = []string{"admin", "director"}

func isRootRole(role string) bool {
	for _, r := range rootRoles {
		if r == role {
			return true
		}
	}
	return false
}

// ListUsers returns accounts matching a search over name and e-mail, newest
// first. An empty query returns everyone; limit is capped so a growing site
// cannot turn the panel into a slow page.
func (s *Store) ListUsers(ctx context.Context, query string, limit int) ([]AdminUser, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	q := "%" + strings.ToLower(strings.TrimSpace(query)) + "%"
	rows, err := s.db.Query(ctx, `
		SELECT u.id, u.email, u.first_name, u.last_name, u.middle_name, u.role,
		       COALESCE(u.signup_country, ''), COALESCE(u.phone, ''),
		       u.email_verified_at IS NOT NULL, u.created_at,
		       (SELECT count(*) FROM articles a WHERE a.author_id = u.id),
		       (SELECT count(*) FROM listings l WHERE l.author_id = u.id),
		       (SELECT count(*) FROM comments c WHERE c.user_id = u.id)
		  FROM auth_users u
		 WHERE $1 = '%%' OR lower(u.email) LIKE $1
		    OR lower(u.first_name) LIKE $1 OR lower(u.last_name) LIKE $1
		 ORDER BY u.created_at DESC
		 LIMIT $2`, q, limit)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	out := []AdminUser{}
	roots := 0
	for rows.Next() {
		var u AdminUser
		if err := rows.Scan(&u.ID, &u.Email, &u.First, &u.Last, &u.Middle, &u.Role,
			&u.Country, &u.Phone, &u.Verified, &u.CreatedAt,
			&u.Articles, &u.Listings, &u.Comments); err != nil {
			return nil, err
		}
		if isRootRole(u.Role) {
			roots++
		}
		out = append(out, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Counting within the page would be wrong on a filtered search, so ask the
	// table directly: a filter must never make the last admin look expendable.
	total, err := s.countRootAdmins(ctx)
	if err != nil {
		return nil, err
	}
	if total <= 1 {
		for i := range out {
			if isRootRole(out[i].Role) {
				out[i].IsLastRoot = true
			}
		}
	}
	return out, nil
}

func (s *Store) countRootAdmins(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRow(ctx,
		`SELECT count(*) FROM auth_users WHERE role = ANY($1)`, rootRoles).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("count admins: %w", err)
	}
	return n, nil
}

// GetAdminUser loads a single account for the edit form.
func (s *Store) GetAdminUser(ctx context.Context, id uuid.UUID) (AdminUser, error) {
	var u AdminUser
	err := s.db.QueryRow(ctx, `
		SELECT id, email, first_name, last_name, middle_name, role,
		       COALESCE(signup_country, ''), COALESCE(phone, ''),
		       email_verified_at IS NOT NULL, created_at
		  FROM auth_users WHERE id = $1`, id).
		Scan(&u.ID, &u.Email, &u.First, &u.Last, &u.Middle, &u.Role,
			&u.Country, &u.Phone, &u.Verified, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return u, ErrUserNotFound
	}
	if err != nil {
		return u, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

// ErrUserNotFound is returned when an admin action names an account that is not
// there — a stale tab, or a guessed id.
var ErrUserNotFound = errors.New("user not found")

// ErrLastAdmin refuses the change that would leave the site with no
// administrator at all.
var ErrLastAdmin = errors.New("the last administrator cannot be removed or demoted")

// UpdateUserByAdmin rewrites the fields an administrator may correct: the name
// carried on every byline, and the e-mail-verified flag.
//
// The e-mail address itself is deliberately not editable here. It is the login
// identity, so changing it silently hands the account to whoever the new address
// belongs to — that is a password reset with extra steps, and it should go
// through the owner, not around them.
func (s *Store) UpdateUserByAdmin(ctx context.Context, id uuid.UUID, first, last, middle string, verified bool) error {
	ct, err := s.db.Exec(ctx, `
		UPDATE auth_users
		   SET first_name = $2, last_name = $3, middle_name = $4,
		       email_verified_at = CASE
		           WHEN $5 AND email_verified_at IS NULL THEN now()
		           WHEN NOT $5 THEN NULL
		           ELSE email_verified_at END,
		       updated_at = now()
		 WHERE id = $1`, id, first, last, middle, verified)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

// adminGuardLock is an arbitrary but fixed key. Every transaction that could
// change the number of administrators takes it, so the count and the change
// that depends on it cannot interleave with another one.
const adminGuardLock = 4823710055

// lockAdminGuard serialises administrator-count changes for this transaction.
// The lock is released when the transaction ends, however it ends.
func lockAdminGuard(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock($1)`, int64(adminGuardLock)); err != nil {
		return fmt.Errorf("lock admin guard: %w", err)
	}
	return nil
}

// countRootAdminsTx counts administrators inside a transaction, so the answer
// still holds when it is acted upon.
func (s *Store) countRootAdminsTx(ctx context.Context, tx pgx.Tx) (int, error) {
	var n int
	if err := tx.QueryRow(ctx,
		`SELECT count(*) FROM auth_users WHERE role = ANY($1)`, rootRoles).Scan(&n); err != nil {
		return 0, fmt.Errorf("count admins: %w", err)
	}
	return n, nil
}

// SetRoleByID changes an account's role, refusing to demote the last
// administrator.
//
// The count and the demotion happen in one transaction behind an advisory lock.
// Read separately, as they were, two concurrent demotions each saw two
// administrators, each concluded it was safe to proceed, and the site was left
// with none.
//
// The lock is taken for every role change, not only for the ones that look like
// demotions, because whether this *is* a demotion is itself decided by reading
// the current role — and that read, taken before the lock as it once was, can be
// stale by the time it is acted on. A promotion racing a demotion could make the
// demoter believe it was touching an ordinary account. Inside the lock the role
// is re-read, so the decision and the change cannot come apart.
func (s *Store) SetRoleByID(ctx context.Context, id uuid.UUID, role string) error {
	primary, normalized := normalizeRoleSet(role)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin role tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if lockErr := lockAdminGuard(ctx, tx); lockErr != nil {
		return lockErr
	}
	var current string
	if err := tx.QueryRow(ctx, `SELECT role FROM auth_users WHERE id = $1`, id).Scan(&current); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("read role: %w", err)
	}

	if isRootRole(current) && !isRootRole(primary) {
		n, cerr := s.countRootAdminsTx(ctx, tx)
		if cerr != nil {
			return cerr
		}
		if n <= 1 {
			return ErrLastAdmin
		}
	}

	// A role change is exactly the event every issued token must stop surviving.
	if _, err := tx.Exec(ctx,
		`UPDATE auth_users SET role = $2, auth_version = auth_version + 1, updated_at = now() WHERE id = $1`,
		id, primary); err != nil {
		return fmt.Errorf("set role: %w", err)
	}
	// The join table is the other half of the truth; leaving it stale would let
	// a demoted account keep reaching role-gated routes.
	if _, err := tx.Exec(ctx, `DELETE FROM auth_user_roles WHERE user_id = $1`, id); err != nil {
		return fmt.Errorf("clear roles: %w", err)
	}
	for _, name := range normalized {
		roleID, ensureErr := s.ensureRoleTx(ctx, tx, name, "")
		if ensureErr != nil {
			return ensureErr
		}
		if assignErr := s.assignRoleTx(ctx, tx, id, roleID); assignErr != nil {
			return assignErr
		}
	}
	return tx.Commit(ctx)
}

// DeleteUserByAdmin removes an account and everything cascading from it.
// Refuses to delete the last administrator; the caller separately refuses to
// let an administrator delete themselves.
func (s *Store) DeleteUserByAdmin(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin delete tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	// Same guard as the demotion path, and for the same reason: two concurrent
	// deletions would otherwise each see a spare administrator that the other
	// was already removing. The role is read under the lock, not before it, so
	// a promotion landing in between cannot hide an administrator from the count.
	if lockErr := lockAdminGuard(ctx, tx); lockErr != nil {
		return lockErr
	}
	var role string
	if err := tx.QueryRow(ctx, `SELECT role FROM auth_users WHERE id = $1`, id).Scan(&role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("read role: %w", err)
	}
	if isRootRole(role) {
		n, cerr := s.countRootAdminsTx(ctx, tx)
		if cerr != nil {
			return cerr
		}
		if n <= 1 {
			return ErrLastAdmin
		}
	}
	ct, err := tx.Exec(ctx, `DELETE FROM auth_users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete account: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return tx.Commit(ctx)
}

// SetSignupCountry records where an account registered from. Best-effort: an
// unresolvable address simply leaves the column empty.
func (s *Store) SetSignupCountry(ctx context.Context, id uuid.UUID, cc string) error {
	cc = strings.ToUpper(strings.TrimSpace(cc))
	if len(cc) != 2 {
		return nil
	}
	_, err := s.db.Exec(ctx,
		`UPDATE auth_users SET signup_country = $2 WHERE id = $1 AND signup_country IS NULL`, id, cc)
	return err
}
