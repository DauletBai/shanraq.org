package articles

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestParseVerdict(t *testing.T) {
	// A clean article returns nothing.
	fs, err := parseVerdict(`{"findings":[]}`)
	if err != nil || len(fs) != 0 {
		t.Fatalf("clean: %v %v", fs, err)
	}
	// Unknown rule codes are dropped: they have no translation to show.
	fs, err = parseVerdict(`{"findings":[{"rule":"made_up","severity":"block"},{"rule":"hatred","severity":"block","quote":"q","note":"n"}]}`)
	if err != nil || len(fs) != 1 || fs[0].RuleCode != "hatred" {
		t.Fatalf("filter: %+v %v", fs, err)
	}
	// A missing severity must default to blocking, not to permissive.
	fs, _ = parseVerdict(`{"findings":[{"rule":"illegal"}]}`)
	if len(fs) != 1 || !fs[0].Blocking() {
		t.Fatalf("default severity must block: %+v", fs)
	}
	// Prose around the JSON is tolerated; garbage is an error, never a pass.
	if _, err := parseVerdict("Here you go: {\"findings\":[]} thanks"); err != nil {
		t.Fatalf("wrapped json: %v", err)
	}
	if _, err := parseVerdict("I could not check this article."); err == nil {
		t.Fatal("unparseable output must be an error, not an empty pass")
	}
}

func TestHumanRejectionRequiresActionableReason(t *testing.T) {
	m := &Module{}
	for _, decision := range []string{"reject", "needs_work", "hide"} {
		err := m.DecideArticle(context.Background(), uuid.New(), decision, uuid.New(), "  ")
		if !errors.Is(err, ErrDecisionNoteRequired) {
			t.Errorf("%s: got %v, want ErrDecisionNoteRequired", decision, err)
		}
	}
}
