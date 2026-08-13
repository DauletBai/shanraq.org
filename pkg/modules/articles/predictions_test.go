package articles

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// The score is the claim the whole page rests on, so its arithmetic is checked
// directly rather than through the handler.
func TestPredictionAccuracy(t *testing.T) {
	cases := []struct {
		name string
		s    PredictionScore
		want int
	}{
		// A ledger with nothing judged claims nothing: 0, not 100.
		{"empty", PredictionScore{}, 0},
		{"only open", PredictionScore{Open: 5}, 0},
		{"perfect", PredictionScore{Hit: 4}, 100},
		{"none right", PredictionScore{Miss: 3}, 0},
		{"half", PredictionScore{Hit: 2, Miss: 2}, 50},
		// A partial counts as half: calling it a hit flatters the number and
		// calling it a miss punishes admitting the middle ground exists.
		{"partial is half", PredictionScore{Partial: 2}, 50},
		{"mixed", PredictionScore{Hit: 3, Partial: 2, Miss: 5}, 40},
		// Open forecasts never move the number, in either direction.
		{"open is not counted", PredictionScore{Hit: 1, Miss: 1, Open: 98}, 50},
	}
	for _, c := range cases {
		if got := c.s.Accuracy(); got != c.want {
			t.Errorf("%s: accuracy = %d, want %d", c.name, got, c.want)
		}
	}
}

// The table forbids a settled forecast with no verdict date and an open one
// that has one. Validate has to fix both before the insert, or an operator gets
// a constraint violation instead of a form.
func TestPredictionInputValidate(t *testing.T) {
	base := func() PredictionInput {
		return PredictionInput{
			Status:    PredOpen,
			Statement: map[string]string{LangRU: "Курс упадёт"},
			Verdict:   map[string]string{},
		}
	}

	in := base()
	in.Status = PredHit
	if err := in.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
	if in.ResolvedOn == nil {
		t.Error("a settled forecast was left with no verdict date; the table would reject it")
	}

	in = base()
	when := time.Now()
	in.ResolvedOn = &when
	if err := in.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
	if in.ResolvedOn != nil {
		t.Error("reopening a forecast must clear its verdict date")
	}

	in = base()
	if in.MadeOn.IsZero() {
		if err := in.Validate(); err != nil {
			t.Fatalf("validate: %v", err)
		}
		if in.MadeOn.IsZero() {
			t.Error("a forecast with no date is undatable and therefore unfalsifiable")
		}
	}

	empty := base()
	empty.Statement = map[string]string{LangRU: "   "}
	if err := empty.Validate(); !errors.Is(err, ErrPredictionEmpty) {
		t.Errorf("a blank forecast was accepted: %v", err)
	}

	bad := base()
	bad.Status = "probably"
	if err := bad.Validate(); err == nil {
		t.Error("an unknown status was accepted")
	}
}

// Round trip through the real database: three languages in, three languages
// out, and the tally counts what was written.
func TestPredictionStoreRoundTrip(t *testing.T) {
	app := newTestApp(t)
	store := NewPredictionStore(app.pool)
	ctx := context.Background()

	made := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	id, err := store.Save(ctx, uuid.Nil, PredictionInput{
		MadeOn: made,
		Status: PredOpen,
		Statement: map[string]string{
			LangRU: "Ставка не опустится ниже 14% до конца года",
			LangKZ: "Мөлшерлеме жыл соңына дейін 14%-дан төмендемейді",
			LangEN: "The rate will not fall below 14% this year",
		},
		Verdict: map[string]string{},
	})
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	t.Cleanup(func() { _ = store.Delete(context.Background(), id) })

	got, err := store.Get(ctx, LangRU, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	for _, l := range Langs {
		if got.StatementIn(l) == "" {
			t.Errorf("%s statement was lost on the round trip", l)
		}
	}
	if got.Status != PredOpen || got.ResolvedOn != nil {
		t.Errorf("open forecast came back as %q / %v", got.Status, got.ResolvedOn)
	}

	// Settle it, and the tally must move.
	if _, err := store.Save(ctx, id, PredictionInput{
		MadeOn:    made,
		Status:    PredMiss,
		Statement: got.Statement,
		Verdict:   map[string]string{LangRU: "Опустилась в сентябре."},
	}); err != nil {
		t.Fatalf("resolve: %v", err)
	}
	got, err = store.Get(ctx, LangRU, id)
	if err != nil {
		t.Fatalf("get after resolve: %v", err)
	}
	if got.Status != PredMiss {
		t.Errorf("status = %q, want %q", got.Status, PredMiss)
	}
	if got.ResolvedOn == nil {
		t.Error("a settled forecast has no verdict date")
	}
	if !strings.Contains(got.VerdictIn(LangRU), "сентябре") {
		t.Errorf("verdict was not stored: %q", got.VerdictIn(LangRU))
	}
	// The statement must survive being resolved: rewriting history is the exact
	// failure this ledger exists to make impossible.
	if !strings.Contains(got.StatementIn(LangRU), "14%") {
		t.Errorf("the original forecast was lost when it was judged: %q", got.StatementIn(LangRU))
	}
}

// A miss has to be as visible as a hit, or the page is advertising rather than
// accounting.
func TestPredictionsPageShowsTheMisses(t *testing.T) {
	app := newTestApp(t)
	store := NewPredictionStore(app.pool)
	ctx := context.Background()

	miss, err := store.Save(ctx, uuid.Nil, PredictionInput{
		MadeOn:    time.Now().AddDate(0, -6, 0),
		Status:    PredMiss,
		Statement: map[string]string{LangRU: "Тенге укрепится до 400"},
		Verdict:   map[string]string{LangRU: "Не укрепился."},
	})
	if err != nil {
		t.Fatalf("save miss: %v", err)
	}
	t.Cleanup(func() { _ = store.Delete(context.Background(), miss) })

	open, err := store.Save(ctx, uuid.Nil, PredictionInput{
		MadeOn:    time.Now(),
		Status:    PredOpen,
		Statement: map[string]string{LangRU: "Ставка останется двузначной"},
		Verdict:   map[string]string{},
	})
	if err != nil {
		t.Fatalf("save open: %v", err)
	}
	t.Cleanup(func() { _ = store.Delete(context.Background(), open) })

	w := app.do(http.MethodGet, "/predictions", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("/predictions = %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Тенге укрепится до 400") {
		t.Error("the missed forecast is not on the public ledger")
	}
	if !strings.Contains(body, "Не укрепился.") {
		t.Error("the verdict is missing, so the reader cannot see what went wrong")
	}
	if !strings.Contains(body, "Ставка останется двузначной") {
		t.Error("the open forecast is missing")
	}
	if !strings.Contains(body, "pcard--miss") {
		t.Error("a miss is not marked as one")
	}
}

// Only leadership manages the ledger; a reader who could edit it could edit
// away a miss, which is the one thing that must be impossible.
func TestPredictionAdminNeedsLeadership(t *testing.T) {
	app := newTestApp(t)
	app.createUser("pred-reader@example.com", "Sup3r-Secret-Pass!")
	cookie := app.login("pred-reader@example.com", "Sup3r-Secret-Pass!")

	if w := app.do(http.MethodGet, "/admin/predictions", nil, withCookie(cookie)); w.Code == http.StatusOK {
		t.Error("an ordinary reader reached the prediction admin")
	}

	// Asserted on the ledger, not on the status code: a rejected write and a
	// redirect to the login page are both 303, so the code alone proves nothing.
	store := NewPredictionStore(app.pool)
	before, err := store.Score(context.Background())
	if err != nil {
		t.Fatalf("score: %v", err)
	}
	form := url.Values{"status": {PredHit}, "statement_ru": {"написано посторонним"}}
	app.do(http.MethodPost, "/admin/predictions", form, withCookie(cookie))
	after, err := store.Score(context.Background())
	if err != nil {
		t.Fatalf("score: %v", err)
	}
	if after.Total() != before.Total() {
		t.Errorf("an ordinary reader wrote to the ledger: %d entries became %d",
			before.Total(), after.Total())
	}
}

// An open forecast past its own deadline is how a ledger gets quietly gamed:
// never judge anything and the score never drops.
func TestOverdueOnlyAppliesToOpenForecasts(t *testing.T) {
	past := time.Now().AddDate(0, -1, 0)
	future := time.Now().AddDate(0, 1, 0)
	if p := (&Prediction{Status: PredOpen, Horizon: &past}); !p.Overdue() {
		t.Error("an open forecast past its deadline is not flagged")
	}
	if p := (&Prediction{Status: PredOpen, Horizon: &future}); p.Overdue() {
		t.Error("a forecast still inside its window was flagged as overdue")
	}
	if p := (&Prediction{Status: PredMiss, Horizon: &past}); p.Overdue() {
		t.Error("a judged forecast cannot be overdue — it has already been judged")
	}
	if p := (&Prediction{Status: PredOpen}); p.Overdue() {
		t.Error("a forecast with no stated deadline cannot be past it")
	}
}

// The management path end to end: leadership writes a forecast through the
// form, then settles it as a miss, and the public page shows both the original
// wording and the verdict. This is the workflow the whole feature is.
func TestPredictionAdminWritesAndSettles(t *testing.T) {
	app := newTestApp(t)
	app.createUser("pred-boss@example.com", "Sup3r-Secret-Pass!")
	app.makeStaff("pred-boss@example.com", "admin")
	cookie := app.login("pred-boss@example.com", "Sup3r-Secret-Pass!")

	store := NewPredictionStore(app.pool)
	form := url.Values{
		"status":       {PredOpen},
		"made_on":      {"2026-02-01"},
		"horizon":      {"2026-12-31"},
		"statement_ru": {"Инфляция не опустится ниже 8%"},
		"statement_en": {"Inflation will not fall below 8%"},
	}
	if w := app.do(http.MethodPost, "/admin/predictions", form, withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("save = %d, want 303", w.Code)
	}
	list, err := store.List(context.Background(), LangRU)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	var made *Prediction
	for _, p := range list {
		if strings.Contains(p.StatementIn(LangRU), "Инфляция не опустится") {
			made = p
		}
	}
	if made == nil {
		t.Fatal("the forecast was not written")
	}
	t.Cleanup(func() { _ = store.Delete(context.Background(), made.ID) })
	if made.Horizon == nil {
		t.Error("the deadline was dropped, so the forecast can never be overdue")
	}

	settle := url.Values{
		"id":           {made.ID.String()},
		"status":       {PredMiss},
		"made_on":      {"2026-02-01"},
		"statement_ru": {"Инфляция не опустится ниже 8%"},
		"verdict_ru":   {"Опустилась до 7,4% в ноябре."},
	}
	if w := app.do(http.MethodPost, "/admin/predictions", settle, withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("settle = %d, want 303", w.Code)
	}

	body := app.do(http.MethodGet, "/predictions", nil).Body.String()
	if !strings.Contains(body, "Инфляция не опустится ниже 8%") {
		t.Error("the original wording vanished once the forecast was judged")
	}
	if !strings.Contains(body, "Опустилась до 7,4% в ноябре.") {
		t.Error("the verdict is not shown to the reader")
	}
}

// An empty ledger showing "0% accuracy" would be the page libelling itself:
// nothing judged is not the same as never being right.
func TestEmptyLedgerClaimsNoScore(t *testing.T) {
	app := newTestApp(t)
	// The ledger is shared, so this only means anything on an empty one.
	if sc, err := NewPredictionStore(app.pool).Score(context.Background()); err != nil {
		t.Fatalf("score: %v", err)
	} else if sc.Resolved() > 0 {
		t.Skip("the ledger already has judged forecasts")
	}
	body := app.do(http.MethodGet, "/predictions", nil).Body.String()
	if strings.Contains(body, `pscore__pct">0%`) {
		t.Error("an empty ledger reports 0% accuracy, which reads as never being right")
	}
	if !strings.Contains(body, "pscore__pct--none") {
		t.Error("an empty ledger does not say that nothing has been judged")
	}
}

// The ledger's real entry point. A page linked only from the footer is a page
// nobody opens; the reader who has just finished the argument is the one who
// wants to know whether its author has been right before.
func TestArticleShowsItsOwnForecasts(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("pred-art@example.com", "Sup3r-Secret-Pass!")
	id, slug := app.seedCitable(author)

	store := NewPredictionStore(app.pool)
	pid, err := store.Save(context.Background(), uuid.Nil, PredictionInput{
		ArticleID: &id,
		MadeOn:    time.Now().AddDate(0, -3, 0),
		Status:    PredHit,
		Statement: map[string]string{LangRU: "Дефицит бюджета превысит план"},
		Verdict:   map[string]string{LangRU: "Превысил на 1,2 трлн."},
	})
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	t.Cleanup(func() { _ = store.Delete(context.Background(), pid) })

	body := app.do(http.MethodGet, "/read/"+slug, nil).Body.String()
	if !strings.Contains(body, "Дефицит бюджета превысит план") {
		t.Error("the article does not show the forecast it made")
	}
	if !strings.Contains(body, "Превысил на 1,2 трлн.") {
		t.Error("the outcome is missing, so the block only makes a claim")
	}
	if !strings.Contains(body, `href="/predictions`) {
		t.Error("the block does not lead to the ledger, which is the point of it")
	}

	// An article that forecast nothing must not grow an empty box.
	_, plain := app.seedCitable(author)
	if b := app.do(http.MethodGet, "/read/"+plain, nil).Body.String(); strings.Contains(b, `class="apred"`) {
		t.Error("an article with no forecasts rendered the forecasts block")
	}
}

// The rating row must be the same shape for everyone. A sentence printed in
// place of the buttons wrapped the row onto two lines for the author of the
// piece and for every signed-out reader — which is most of them.
func TestRatingRowKeepsItsShapeForEveryone(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("rate-shape@example.com", "Sup3r-Secret-Pass!")
	_, slug := app.seedCitable(author)

	for _, c := range []struct {
		who  string
		opts []reqOpt
	}{
		{"guest", nil},
		{"author", []reqOpt{withCookie(app.login("rate-shape@example.com", "Sup3r-Secret-Pass!"))}},
	} {
		body := app.do(http.MethodGet, "/read/"+slug, nil, c.opts...).Body.String()
		if n := strings.Count(body, "rating__btn"); n < 2 {
			t.Errorf("%s sees %d rating buttons, want the same two everyone else sees", c.who, n)
		}
		// The explanation belongs in the tooltip, not in the row.
		if strings.Contains(body, `class="hint" style="padding:0 6px"`) {
			t.Errorf("%s still gets the sentence that pushed the row onto two lines", c.who)
		}
		if !strings.Contains(body, "fhelp__tip") {
			t.Errorf("%s gets no explanation at all", c.who)
		}
	}
}
