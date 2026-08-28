package media

import (
	"context"
	"errors"
	"strings"
)

// SaveBlob stores bytes the system produced for itself and returns the URL a
// browser can fetch them from.
//
// It is deliberately not the upload path. Uploads are metered against an
// owner's quota and recorded in the ledger, because a person put them there and
// somebody has to answer for the space. This is for files the site generates:
// narration for an article, and whatever follows it. Nobody's quota should move
// because we rendered audio.
//
// Not recording them has a second, sharper reason. The sweeper decides a file
// is an orphan by searching article bodies for its key, and a generated file is
// referenced by a column rather than by the prose. Recorded, it would be swept
// on the next pass -- silently, and only for the articles nobody had opened
// recently enough to notice. Its lifetime is owned by whatever row points at it.
func (m *Module) SaveBlob(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	if m == nil || m.store == nil {
		return "", errors.New("media: storage not configured")
	}
	key = strings.TrimPrefix(strings.TrimSpace(key), "/")
	if key == "" {
		return "", errors.New("media: empty key")
	}
	if len(data) == 0 {
		return "", errors.New("media: empty blob")
	}
	if err := m.store.Put(ctx, key, data, contentType); err != nil {
		return "", err
	}
	return m.store.URL(key), nil
}

// DeleteBlob removes a generated file. A missing object is not an error, so a
// caller cleaning up after a half-finished write does not have to check first.
func (m *Module) DeleteBlob(ctx context.Context, key string) error {
	if m == nil || m.store == nil {
		return nil
	}
	key = strings.TrimPrefix(strings.TrimSpace(key), "/")
	if key == "" {
		return nil
	}
	return m.store.Delete(ctx, key)
}
