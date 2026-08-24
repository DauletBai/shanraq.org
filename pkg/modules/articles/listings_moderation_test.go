package articles

import (
	"context"
	"strings"
	"testing"
)

// A rejected listing has to explain itself in full. "Hidden: covert advertising"
// tells the author nothing they can act on — they are left guessing which phrase
// gave offence, and a guess is not a fix. So the listing check answers with the
// same contract as the article check: a list of findings, each carrying the rule,
// the exact quotation, and what to change.
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
	// The remark is written in the author's language, not the prompt's.
	if !strings.Contains(p, "Write \"note\" in ru") {
		t.Error("язык замечания не задан")
	}
	// And the thing the check must not complain about: an honest listing.
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

// Half of what makes a listing untrue is invisible in the description on its own —
// it lies in the gap between the description and the fields. "A spacious 120 m²
// flat" against a declared 54 is the single commonest complaint, and it can only
// be seen if both halves lie in front of the checker.
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

// The check runs in the background rather than in the form handler: at three
// requests a minute a synchronous call would hold the form for up to twenty
// seconds, and the whole point of a classifieds board is that posting takes
// seconds. And with the model switched off no job is queued at all — posting works
// exactly as it did before.
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

// A listing already withdrawn, expired or reported by readers is nothing to check:
// a decision on it has been taken, and the machine has no business revisiting it.
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
