package syndicate

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// Ключ IndexNow сайт выписывает себе сам.
//
// Протокол доказывает владение доменом просто: тот же ключ лежит файлом на
// сайте, и заявка ссылается на этот файл. Секрета в нём нет — есть требование
// постоянства: сменившийся ключ делает недействительным уже отданный файл, и
// все заявки начинают отвергаться.
//
// Поэтому ключ не берётся из конфига и не выдумывается заново при запуске: он
// выписывается один раз и хранится в базе. Настройка руками не нужна, а значит
// нельзя и забыть её сделать — до сих пор IndexNow молчал ровно поэтому.

// indexNowSetting — имя строки в app_settings.
const indexNowSetting = "indexnow_key"

// ensureIndexNowKey возвращает ключ сайта, выписывая его при первом обращении.
func (m *Module) ensureIndexNowKey(ctx context.Context) (string, error) {
	var key string
	err := m.db.QueryRow(ctx,
		`SELECT value FROM app_settings WHERE name = $1`, indexNowSetting).Scan(&key)
	if err == nil {
		if k, ok := normalizeIndexNowKey(key); ok && k != "" {
			return k, nil
		}
		// Строка есть, но негодная — перевыписываем: молчащий IndexNow хуже
		// смены ключа, которую всё равно никто ещё не подтвердил.
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
