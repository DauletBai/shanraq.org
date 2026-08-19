package media

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"shanraq.org/pkg/shanraq"
)

// ErrQuotaExceeded means this account already holds as much as it may.
var ErrQuotaExceeded = errors.New("storage quota exceeded")

// ErrStoreFull means the store as a whole is at its cap. It is deliberately a
// different error: one account being over its share and the disk being gone are
// different problems and want different answers.
var ErrStoreFull = errors.New("media store is full")

// ledger is the accounting half of storage: who uploaded what, how big it was,
// and — by the absence of any reference to it — what nobody is using.
//
// It is separate from Store on purpose. Store answers "where are the bytes",
// and gains an S3 backend one day without learning anything about accounts;
// ledger answers "whose bytes are they", which is a question about this site.
type ledger struct {
	db *pgxpool.Pool
}

// usage is how much this account holds. Keys are content hashes, so two people
// who uploaded the same photo each count it: the alternative is an account
// whose usage falls when a stranger deletes something.
func (l *ledger) usage(ctx context.Context, owner uuid.UUID) (int64, error) {
	var n int64
	err := l.db.QueryRow(ctx, `
		SELECT COALESCE(sum(o.bytes), 0)
		  FROM media_owners w JOIN media_objects o ON o.key = w.key
		 WHERE w.user_id = $1`, owner).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("media usage: %w", err)
	}
	return n, nil
}

// total is the whole store, counting each object once.
func (l *ledger) total(ctx context.Context) (int64, error) {
	var n int64
	if err := l.db.QueryRow(ctx, `SELECT COALESCE(sum(bytes), 0) FROM media_objects`).Scan(&n); err != nil {
		return 0, fmt.Errorf("media total: %w", err)
	}
	return n, nil
}

// held reports whether this account already owns this exact object, in which
// case storing it again costs nothing and must not be charged twice.
func (l *ledger) held(ctx context.Context, owner uuid.UUID, key string) (bool, error) {
	var ok bool
	err := l.db.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM media_owners WHERE user_id = $1 AND key = $2)`, owner, key).Scan(&ok)
	if err != nil {
		return false, fmt.Errorf("media held: %w", err)
	}
	return ok, nil
}

// known reports whether the object is already on disk under someone else's
// name. It costs no new bytes, so it does not count against the store cap.
func (l *ledger) known(ctx context.Context, key string) (bool, error) {
	var ok bool
	if err := l.db.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM media_objects WHERE key = $1)`, key).Scan(&ok); err != nil {
		return false, fmt.Errorf("media known: %w", err)
	}
	return ok, nil
}

// record files the object and the claim on it. Both halves are idempotent: the
// same upload twice is one object and one claim.
func (l *ledger) record(ctx context.Context, key string, bytes int64, contentType string, owner uuid.UUID) error {
	tx, err := l.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("media record: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if _, err := tx.Exec(ctx, `
		INSERT INTO media_objects (key, bytes, content_type) VALUES ($1, $2, $3)
		ON CONFLICT (key) DO NOTHING`, key, bytes, contentType); err != nil {
		return fmt.Errorf("media record object: %w", err)
	}
	if owner != uuid.Nil {
		if _, err := tx.Exec(ctx, `
			INSERT INTO media_owners (key, user_id) VALUES ($1, $2)
			ON CONFLICT DO NOTHING`, key, owner); err != nil {
			return fmt.Errorf("media record owner: %w", err)
		}
	}
	return tx.Commit(ctx)
}

// orphanQuery finds stored objects that no URL column points at. The columns
// listed here are every place a media URL is written as a whole value; free
// text that can embed one — article bodies, editable pages, ad copy, comments —
// is checked separately, because a substring match belongs nowhere near a
// cross join.
const orphanQuery = `
WITH refs AS (
	    SELECT cover_url AS u FROM articles
	UNION ALL SELECT avatar_url FROM auth_users
	UNION ALL SELECT cover_url FROM listings
	UNION ALL SELECT COALESCE(contract_url, '') FROM listings
	UNION ALL SELECT unnest(images) FROM listings
	UNION ALL SELECT unnest(documents) FROM listings
	UNION ALL SELECT image_url FROM ad_orders
)
SELECT o.key, o.bytes
  FROM media_objects o
 WHERE o.created_at < now() - make_interval(hours => $1::int)
   AND NOT EXISTS (SELECT 1 FROM refs WHERE refs.u <> '' AND refs.u LIKE '%' || o.key)
 ORDER BY o.created_at
 LIMIT $2`

// bodyQuery pulls the free text that could mention a media URL — and only the
// rows that do, so what comes back is a handful of strings rather than the
// site's prose.
const bodyQuery = `
	    SELECT COALESCE(body_md, '') FROM article_translations WHERE body_md LIKE $1
	UNION ALL SELECT COALESCE(body_md, '') FROM content_pages     WHERE body_md LIKE $1
	UNION ALL SELECT COALESCE(body, '')    FROM ad_orders         WHERE body    LIKE $1
	UNION ALL SELECT COALESCE(body, '')    FROM comments          WHERE body    LIKE $1`

type orphan struct {
	key   string
	bytes int64
}

// orphans lists objects older than grace that nothing on the site refers to.
//
// The grace period is the whole safety margin here: an upload becomes reachable
// only when the form it belongs to is saved, so anything younger than that is
// presumed to be a draft in progress rather than litter.
func (l *ledger) orphans(ctx context.Context, publicPrefix string, graceHours, limit int) ([]orphan, error) {
	rows, err := l.db.Query(ctx, orphanQuery, graceHours, limit)
	if err != nil {
		return nil, fmt.Errorf("media orphans: %w", err)
	}
	var candidates []orphan
	for rows.Next() {
		var o orphan
		if err := rows.Scan(&o.key, &o.bytes); err != nil {
			rows.Close()
			return nil, fmt.Errorf("media orphan scan: %w", err)
		}
		candidates = append(candidates, o)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("media orphans: %w", err)
	}
	if len(candidates) == 0 {
		return nil, nil
	}

	bodies, err := l.mentioningBodies(ctx, publicPrefix)
	if err != nil {
		// Not knowing what the prose refers to is not a licence to delete: a
		// sweep that cannot read half the references must not run at all.
		return nil, err
	}
	kept := candidates[:0]
	for _, c := range candidates {
		referenced := false
		for _, b := range bodies {
			if strings.Contains(b, c.key) {
				referenced = true
				break
			}
		}
		if !referenced {
			kept = append(kept, c)
		}
	}
	return kept, nil
}

func (l *ledger) mentioningBodies(ctx context.Context, publicPrefix string) ([]string, error) {
	pattern := "%" + strings.TrimSuffix(publicPrefix, "/") + "/%"
	rows, err := l.db.Query(ctx, bodyQuery, pattern)
	if err != nil {
		return nil, fmt.Errorf("media bodies: %w", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, fmt.Errorf("media body scan: %w", err)
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// forget drops the rows for objects whose bytes are gone.
func (l *ledger) forget(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	_, err := l.db.Exec(ctx, `DELETE FROM media_objects WHERE key = ANY($1)`, keys)
	if err != nil {
		return fmt.Errorf("media forget: %w", err)
	}
	return nil
}

// sweepInterval is how often the orphan sweep runs. Litter is not urgent; the
// point is that it is collected at all, not that it is collected promptly.
const sweepInterval = 6 * time.Hour

// orphanBatch bounds one sweep. A first run on a store that has never been
// swept should not turn into a single enormous transaction.
const orphanBatch = 500

// Start runs the orphan sweep in the background. Media is the one module whose
// storage grows from user action and shrinks from nothing: an upload that never
// made it into a saved listing is invisible to every screen on the site and
// still occupies the disk forever.
func (m *Module) Start(ctx context.Context, _ *shanraq.Runtime) error {
	if m.ledger == nil {
		return nil
	}
	go m.sweepLoop(ctx)
	return nil
}

func (m *Module) sweepLoop(ctx context.Context) {
	t := time.NewTicker(sweepInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			m.sweepOrphans(ctx)
		}
	}
}

// sweepOrphans deletes stored files nothing refers to any more, oldest first.
//
// It deliberately never deletes a file it cannot account for: only objects in
// the ledger are candidates, so anything uploaded before this existed — or
// written by hand into the media directory — is left where it is.
func (m *Module) sweepOrphans(ctx context.Context) {
	grace := m.cfg.OrphanGraceHours
	if grace <= 0 {
		return // sweeping disabled
	}
	found, err := m.ledger.orphans(ctx, m.cfg.PublicPrefix, grace, orphanBatch)
	if err != nil {
		m.logger.Warn("media sweep", zap.Error(err))
		return
	}
	if len(found) == 0 {
		return
	}
	keys := make([]string, 0, len(found))
	var freed int64
	for _, o := range found {
		if err := m.store.Delete(ctx, o.key); err != nil {
			// Leave the row: the file is still there, and the next sweep will
			// try again. Forgetting it here would make it invisible instead.
			m.logger.Warn("media sweep delete", zap.Error(err), zap.String("key", o.key))
			continue
		}
		keys = append(keys, o.key)
		freed += o.bytes
	}
	if err := m.ledger.forget(ctx, keys); err != nil {
		m.logger.Warn("media sweep forget", zap.Error(err))
		return
	}
	m.logger.Info("media sweep collected unreferenced uploads",
		zap.Int("files", len(keys)), zap.Int64("bytes", freed))
}
