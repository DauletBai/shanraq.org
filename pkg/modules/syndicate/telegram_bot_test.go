package syndicate

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// The place line heads the message so a personal delivery answers the question
// it provokes: why am I being told this.
func TestPlacedMessageLeadsWithThePlace(t *testing.T) {
	msg := buildPlacedMessage("Качар & район", "Тепло дали", "Итоги", "https://shanraq.org/read/x")
	if !strings.HasPrefix(msg, "📍 <b>Качар &amp; район</b>") {
		t.Errorf("сообщение не начинается с места: %q", msg)
	}
	if !strings.Contains(msg, "Тепло дали") {
		t.Errorf("заголовок потерян: %q", msg)
	}
	// A place with no name must not leave an empty pin hanging over the title.
	plain := buildPlacedMessage("  ", "Заголовок", "", "https://shanraq.org/read/y")
	if strings.Contains(plain, "📍") {
		t.Errorf("пустое место оставило значок: %q", plain)
	}
}

func TestTelegramLanguageFallsBackToRussian(t *testing.T) {
	for code, want := range map[string]string{
		"kk": "kz", "kk-KZ": "kz", "en": "en", "en-GB": "en",
		"ru": "ru", "": "ru", "de": "ru",
	} {
		if got := tgLang(code); got != want {
			t.Errorf("tgLang(%q) = %q, ожидалось %q", code, got, want)
		}
	}
}

// Every string the bot says must exist in all three languages: a subscriber who
// set their phone to Kazakh and got a blank message would have no way to tell
// the bot is working at all.
func TestBotSpeaksAllThreeLanguages(t *testing.T) {
	for _, key := range []string{"welcome", "pick_region", "pick_inside", "pick_this", "back", "chosen", "stopped", "help"} {
		for _, lang := range []string{"kz", "ru", "en"} {
			if s := tgText(lang, key); strings.TrimSpace(s) == "" {
				t.Errorf("нет строки %q для языка %q", key, lang)
			}
		}
	}
}

// The delivery rule, against the real tree: a piece written for a region reaches
// everyone inside it, and a piece written for one village reaches that village
// alone. This is the whole point of moving placed articles off the channel, so
// it is asserted rather than assumed.
func TestSubscribersInsideFollowsTheTreeDownwards(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run this integration test")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()
	m := &Module{db: pool, log: zap.NewNop()}

	place := func(slug string) (id string) {
		if err := pool.QueryRow(ctx, `SELECT id::text FROM geo_nodes WHERE slug=$1`, slug).Scan(&id); err != nil {
			t.Skipf("в справочнике нет места %q", slug)
		}
		return id
	}
	oblast, city, village := place("kostanaiskaya-oblast"), place("kostanai"), place("kachar")

	// Three subscribers, one per level, with ids far outside any real range.
	subs := map[string]int64{oblast: -9001, city: -9002, village: -9003}
	for node, uid := range subs {
		if _, err := pool.Exec(ctx,
			`INSERT INTO telegram_subscribers (tg_user_id, geo_node_id, lang) VALUES ($1,$2::uuid,'ru')
			 ON CONFLICT (tg_user_id) DO UPDATE SET geo_node_id = EXCLUDED.geo_node_id, active = TRUE`,
			uid, node); err != nil {
			t.Fatalf("insert subscriber: %v", err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM telegram_subscribers WHERE tg_user_id IN (-9001,-9002,-9003)`)
	})

	got := func(node string) map[int64]bool {
		id := uuid.MustParse(node)
		ids, _, err := m.subscribersInside(ctx, id)
		if err != nil {
			t.Fatalf("subscribersInside: %v", err)
		}
		set := map[int64]bool{}
		for _, x := range ids {
			if x <= -9001 && x >= -9003 { // ignore anything a parallel test left
				set[x] = true
			}
		}
		return set
	}

	// A regional piece reaches all three.
	r := got(oblast)
	for _, want := range []int64{-9001, -9002, -9003} {
		if !r[want] {
			t.Errorf("областная статья не дошла до подписчика %d", want)
		}
	}
	// A village piece reaches the village alone.
	v := got(village)
	if !v[-9003] {
		t.Error("качарская статья не дошла до подписчика в Качаре")
	}
	if v[-9001] || v[-9002] {
		t.Errorf("качарская статья ушла за пределы Качара: %v", v)
	}
	// A city piece does not leak sideways into the village.
	c := got(city)
	if !c[-9002] {
		t.Error("костанайская статья не дошла до подписчика в Костанае")
	}
	if c[-9003] {
		t.Error("костанайская статья ушла в Качар")
	}
}
