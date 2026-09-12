package payments

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// The settings row is what stands between "an admin clicked something" and real
// money moving, so its validation is tested against the real table. Skipped
// unless SHANRAQ_TEST_DB names a test database.
func TestSettingsSaveValidation(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the payment settings test")
	}
	if !strings.Contains(dsn, "test") {
		t.Fatalf("SHANRAQ_TEST_DB must name a test database; refusing %q", dsn)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM payment_settings WHERE id = 1`) })

	st := NewSettingsStore(pool, Settings{})
	if err := st.Save(ctx, Settings{Enabled: true, Provider: "bogus"}, nil); err == nil {
		t.Error("enabled + unknown provider must be rejected")
	}
	if err := st.Save(ctx, Settings{Enabled: true, Provider: ProviderKaspi}, nil); err != nil {
		t.Errorf("enabled + kaspi must save: %v", err)
	}
	if got := st.Get(); !got.Enabled || got.Provider != ProviderKaspi {
		t.Errorf("cache not refreshed after save: %+v", got)
	}
	if err := st.Save(ctx, Settings{Enabled: false}, nil); err != nil {
		t.Errorf("disabled must save: %v", err)
	}
}
