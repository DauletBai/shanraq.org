package articles

import (
	"context"
	"strings"
	"testing"
)

// Забракованное объявление должно объяснять себя целиком. «Скрыто: скрытая
// реклама» не говорит автору ничего, что можно исправить, — он остаётся гадать,
// какая фраза не понравилась, а догадка не есть исправление. Поэтому проверка
// объявлений отвечает тем же контрактом, что и проверка статей: список находок,
// в каждой — правило, точная цитата и что менять.
func TestListingPromptDemandsAnExplanationPerFinding(t *testing.T) {
	p := listingReviewPrompt("ru", "текст объявления")

	for _, code := range []string{
		"misleading:", "disguised_ad:", "scam:", "spam:",
		"prohibited_goods:", "personal_data:", "illegal:",
	} {
		if !strings.Contains(p, code) {
			t.Errorf("правило %q не описано", code)
		}
	}
	for _, must := range []string{
		`MUST quote the exact words`,
		`A finding the author cannot locate is useless`,
		`"findings"`,
		`"quote"`,
		`"note"`,
	} {
		if !strings.Contains(p, must) {
			t.Errorf("промпт не требует объяснения: нет %q", must)
		}
	}
	// Замечание пишется на языке автора, а не на языке промпта.
	if !strings.Contains(p, "Write \"note\" in ru") {
		t.Error("язык замечания не задан")
	}
	// И то, из-за чего проверка не должна ругаться на честное объявление.
	for _, guard := range []string{
		"must return an empty findings array",
		"badly written or overpriced listing breaks no rule",
		"own contact details are not personal_data",
	} {
		if !strings.Contains(p, guard) {
			t.Errorf("нет защиты от ложных срабатываний: %q", guard)
		}
	}
}

// Половина того, что делает объявление недостоверным, не видна в описании
// самом по себе — она в расхождении между описанием и полями. «Просторная
// квартира 120 м²» при заявленных 54 — это и есть самая частая жалоба, и
// увидеть её можно, только если обе половины лежат перед проверяющим.
func TestScreeningTextCarriesTheDeclaredFields(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	ownerID := app.createUser("listmod@example.com", "Parol123!")
	listingID := app.seedListing(ownerID)

	text, status, title, author, err := app.module().listingScreeningText(context.Background(), listingID)
	if err != nil {
		t.Fatalf("listingScreeningText: %v", err)
	}
	if status != "published" {
		t.Errorf("статус %q", status)
	}
	if author != ownerID {
		t.Errorf("автор потерян: %v против %v", author, ownerID)
	}
	if title == "" {
		t.Error("заголовок пуст")
	}
	for _, must := range []string{"area:", "rooms:", "price:", "Description:"} {
		if !strings.Contains(text, must) {
			t.Errorf("в тексте проверки нет %q:\n%s", must, text)
		}
	}
	if !strings.Contains(text, "60") {
		t.Errorf("заявленная площадь не попала в текст проверки:\n%s", text)
	}
}

// Проверка идёт в фоне, а не в обработчике формы: при трёх запросах в минуту
// синхронный вызов держал бы форму до двадцати секунд, а весь смысл доски
// объявлений в том, что подача занимает секунды. И если модель выключена,
// задача не ставится вовсе — подача работает ровно как раньше.
func TestScreeningIsQueuedNotAwaited(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	ownerID := app.createUser("listqueue@example.com", "Parol123!")
	listingID := app.seedListing(ownerID)

	app.module().enqueueListingScreening(context.Background(), listingID, ownerID)

	var queued int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM job_queue WHERE name = $1`, JobModerateListing).Scan(&queued); err != nil {
		t.Fatalf("очередь: %v", err)
	}
	if queued != 0 {
		t.Errorf("с выключенной моделью поставлено %d задач проверки", queued)
	}
}

// Уже снятое, истёкшее или пожалованное читателями объявление проверять нечего:
// решение по нему уже принято, и машине незачем его пересматривать.
func TestScreeningSkipsListingsThatAreNoLongerPublished(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	ownerID := app.createUser("listgone@example.com", "Parol123!")
	listingID := app.seedListing(ownerID)
	app.exec(`UPDATE listings SET status = 'flagged' WHERE id = $1`, listingID)

	_, status, _, _, err := app.module().listingScreeningText(context.Background(), listingID)
	if err != nil {
		t.Fatalf("listingScreeningText: %v", err)
	}
	if status == "published" {
		t.Error("статус не прочитан: проверка не сможет пропустить снятое объявление")
	}
}
