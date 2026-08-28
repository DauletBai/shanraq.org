package articles

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Narration is a pre-rendered reading of one article in one language.
//
// The read-aloud button normally uses the browser's own speech synthesis, which
// is free and needs nothing from us. It also has no Kazakh voice -- no browser
// ships one -- so for Kazakh the player falls back to offering a Russian voice,
// and a Kazakh article gets read in Russian. Narration replaces that guess with
// an actual recording where one exists, and changes nothing where it does not.
type Narration struct {
	Lang        string
	URL         string
	StorageKey  string
	DurationSec int
	Bytes       int64
	Voice       string
	TextSHA256  string
	// Cues is the block-to-time map, already JSON, passed to the page as-is.
	// It is stored and served opaque because only the player reads it, and
	// parsing it here would buy nothing but a second place to keep its shape.
	Cues      []byte
	UpdatedAt time.Time
	// Stale reports that the article was edited after this audio was made. The
	// recording still plays; it just no longer matches the page. Saying so is
	// better than either hiding it or letting it quietly contradict the text.
	Stale bool
}

// AudioStore persists narration rows.
type AudioStore struct{ db *pgxpool.Pool }

// NewAudioStore builds an AudioStore over the shared pool.
func NewAudioStore(db *pgxpool.Pool) *AudioStore { return &AudioStore{db: db} }

// TextDigest is the fingerprint a narration is matched against.
//
// It hashes the text as the synthesiser received it, not the Markdown: the
// markup does not reach the ear, so a change from italics to bold must not mark
// perfectly good audio stale. Whitespace is collapsed for the same reason.
func TextDigest(title, body string) string {
	clean := func(s string) string { return strings.Join(strings.Fields(s), " ") }
	sum := sha256.Sum256([]byte(clean(title) + "\n" + clean(body)))
	return hex.EncodeToString(sum[:])
}

// Get returns the narration for one article and language, if there is one.
func (s *AudioStore) Get(ctx context.Context, articleID uuid.UUID, lang, digest string) (*Narration, error) {
	if s == nil || s.db == nil {
		return nil, nil
	}
	var n Narration
	err := s.db.QueryRow(ctx,
		`SELECT lang, url, storage_key, duration_sec, bytes, voice, text_sha256, cues, updated_at
		   FROM article_audio WHERE article_id=$1 AND lang=$2`,
		articleID, lang).Scan(&n.Lang, &n.URL, &n.StorageKey, &n.DurationSec,
		&n.Bytes, &n.Voice, &n.TextSHA256, &n.Cues, &n.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("article audio get: %w", err)
	}
	// An empty digest on either side means "cannot tell", and cannot-tell must
	// not read as stale: rows written before the column existed would all light
	// up as out of date on the first page view.
	n.Stale = digest != "" && n.TextSHA256 != "" && n.TextSHA256 != digest
	return &n, nil
}

// Upsert records narration for one article and language, replacing whatever was
// there. It returns the storage key it replaced, if any, so the caller can drop
// the old file once the row safely points at the new one.
func (s *AudioStore) Upsert(ctx context.Context, articleID uuid.UUID, n Narration) (replaced string, err error) {
	if s == nil || s.db == nil {
		return "", errors.New("article audio: no database")
	}
	err = s.db.QueryRow(ctx,
		`INSERT INTO article_audio
		     (article_id, lang, storage_key, url, duration_sec, bytes, voice, text_sha256, cues, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NOW())
		 ON CONFLICT (article_id, lang) DO UPDATE
		   SET storage_key=EXCLUDED.storage_key, url=EXCLUDED.url,
		       duration_sec=EXCLUDED.duration_sec, bytes=EXCLUDED.bytes,
		       voice=EXCLUDED.voice, text_sha256=EXCLUDED.text_sha256,
		       cues=EXCLUDED.cues, updated_at=NOW()
		 RETURNING (SELECT storage_key FROM article_audio
		             WHERE article_id=$1 AND lang=$2)`,
		articleID, n.Lang, n.StorageKey, n.URL, n.DurationSec, n.Bytes,
		n.Voice, n.TextSHA256, cuesOrEmpty(n.Cues)).Scan(&replaced)
	if err != nil {
		return "", fmt.Errorf("article audio upsert: %w", err)
	}
	if replaced == n.StorageKey {
		// Same key overwritten in place: deleting it would delete what we just
		// wrote.
		return "", nil
	}
	return replaced, nil
}

// Delete removes the narration row and returns the key that backed it.
func (s *AudioStore) Delete(ctx context.Context, articleID uuid.UUID, lang string) (string, error) {
	if s == nil || s.db == nil {
		return "", nil
	}
	var key string
	err := s.db.QueryRow(ctx,
		`DELETE FROM article_audio WHERE article_id=$1 AND lang=$2 RETURNING storage_key`,
		articleID, lang).Scan(&key)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("article audio delete: %w", err)
	}
	return key, nil
}

// cuesOrEmpty keeps the column valid JSON when a caller sends none. A narration
// without timings still plays; it just cannot move the highlight.
func cuesOrEmpty(b []byte) []byte {
	if len(b) == 0 {
		return []byte("[]")
	}
	return b
}
