package syndicate

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// The site issues its own IndexNow key.
//
// The protocol proves domain ownership simply: the same key sits as a file on the
// site, and a submission points at that file. There is no secret in it — there is a
// requirement of permanence: a changed key invalidates the file already served, and
// every submission starts being rejected.
//
// So the key is neither taken from the config nor invented afresh at start-up: it is
// issued once and kept in the database. No manual setting is needed, which means it
// cannot be forgotten either — that is exactly why IndexNow had been silent.

// indexNowSetting is the row's name in app_settings.
const indexNowSetting = "indexnow_key"

// ensureIndexNowKey returns the site's key, issuing it on first use.
func (m *Module) ensureIndexNowKey(ctx context.Context) (string, error) {
	var key string
	err := m.db.QueryRow(ctx,
		`SELECT value FROM app_settings WHERE name = $1`, indexNowSetting).Scan(&key)
	if err == nil {
		if k, ok := normalizeIndexNowKey(key); ok && k != "" {
			return k, nil
		}
		// The row exists but is unusable — reissue: a silent IndexNow is worse than a key
		// change nobody had confirmed anyway.
	}

	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("ключ indexnow: %w", err)
	}
	key = hex.EncodeToString(buf)

	if _, err := m.db.Exec(ctx, `
		INSERT INTO app_settings (name, value) VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value, updated_at = now()`,
		indexNowSetting, key); err != nil {
		return "", fmt.Errorf("сохранение ключа indexnow: %w", err)
	}
	return key, nil
}
