package syndicate

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Ключ должен переживать перезапуск. Протокол проверяет владение доменом по
// файлу с ключом, и новый ключ при каждом запуске делал бы уже отданный файл
// недействительным — все заявки отвергались бы.
func TestTheIndexNowKeySurvivesARestart(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("нужна SHANRAQ_TEST_DB")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("подключение: %v", err)
	}
	defer pool.Close()
	if _, err := pool.Exec(ctx, `DELETE FROM app_settings WHERE name = $1`, indexNowSetting); err != nil {
		t.Fatalf("очистка: %v", err)
	}

	m := &Module{db: pool, log: zap.NewNop(), baseURL: "https://shanraq.org"}
	first, err := m.ensureIndexNowKey(ctx)
	if err != nil {
		t.Fatalf("ключ не выписан: %v", err)
	}
	if len(first) < 8 {
		t.Fatalf("ключ длиной %d — протокол требует не меньше восьми знаков", len(first))
	}
	// Второй запуск: другой экземпляр модуля над той же базой.
	second, err := (&Module{db: pool, log: zap.NewNop()}).ensureIndexNowKey(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if second != first {
		t.Errorf("после перезапуска ключ сменился: %s → %s", first, second)
	}
	if _, ok := normalizeIndexNowKey(first); !ok {
		t.Errorf("выписанный ключ не проходит собственную проверку: %q", first)
	}
}

// Негодная строка в настройках не должна оставлять сайт без ключа: молчащий
// IndexNow хуже, чем ключ, которого ещё никто не подтверждал.
func TestABrokenStoredKeyIsReissued(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("нужна SHANRAQ_TEST_DB")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("подключение: %v", err)
	}
	defer pool.Close()
	if _, err := pool.Exec(ctx, `
		INSERT INTO app_settings (name, value) VALUES ($1, 'не ключ')
		ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value`, indexNowSetting); err != nil {
		t.Fatalf("подготовка: %v", err)
	}
	m := &Module{db: pool, log: zap.NewNop()}
	key, err := m.ensureIndexNowKey(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := normalizeIndexNowKey(key); !ok || key == "не ключ" {
		t.Errorf("негодный ключ остался в силе: %q", key)
	}
}

// Файл ключа обязан отдаваться: без него поисковик отвергает каждую заявку.
func TestTheKeyFileIsServed(t *testing.T) {
	m := &Module{indexNowKey: "0123456789abcdef", log: zap.NewNop()}
	w := httptest.NewRecorder()
	m.handleIndexNowKey(w, httptest.NewRequest(http.MethodGet, indexNowKeyPath, nil))
	if w.Code != http.StatusOK {
		t.Fatalf("файл ключа отдал %d", w.Code)
	}
	if got := strings.TrimSpace(w.Body.String()); got != "0123456789abcdef" {
		t.Errorf("в файле лежит %q, а не ключ", got)
	}
}
