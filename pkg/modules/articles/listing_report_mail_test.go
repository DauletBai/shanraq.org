package articles

import (
	"strings"
	"testing"

	"shanraq.org/internal/config"
	"shanraq.org/pkg/shanraq"
)

func reportMailModule() *Module {
	return &Module{rt: &shanraq.Runtime{Config: config.Config{PublicBaseURL: "https://shanraq.org"}}}
}

// A seller whose photos are reported has to learn three things: that it
// happened, what it will cost them, and where to fix it. The warning used to
// carry only the first, and it carried it in the REPORTER's interface language
// — so a Kazakh seller reported by an English-reading visitor got an English
// notice about their own listing. Accounts store no language preference, so all
// three are sent.
func TestReportEmailIsTrilingualAndActionable(t *testing.T) {
	m := reportMailModule()
	l := &Listing{ID: "11111111-1111-1111-1111-111111111111", TitleRu: "Сдается Студия"}

	subject, body := m.reportEmail(l, 1, false)
	for _, want := range []string{"жалоба", "шағым", "report"} {
		if !strings.Contains(strings.ToLower(subject), want) {
			t.Errorf("subject %q is missing the %q wording", subject, want)
		}
	}
	// Every language, so whoever opens it can read it.
	for lang, phrase := range map[string]string{
		"ru": "Убедитесь, что фото достоверны",
		"kz": "Фотоларыңыз",
		"en": "misrepresent reality",
	} {
		if !strings.Contains(body, phrase) {
			t.Errorf("the %s text is missing from the warning", lang)
		}
	}
	// The consequence, stated as a number rather than implied.
	if !strings.Contains(body, "После 3 жалоб") {
		t.Error("the warning does not say how many reports hide the listing")
	}
	// A way to act, not just a notice.
	if !strings.Contains(body, "https://shanraq.org/listings/"+l.ID+"/edit") {
		t.Error("the warning gives no link to fix the photos")
	}
	if !strings.Contains(body, l.TitleRu) {
		t.Error("the warning does not say which listing it is about")
	}
}

// Once it is hidden the message changes: the count is spent, what matters is
// that it is off the site and how to get it back.
func TestReportEmailSaysWhenHidden(t *testing.T) {
	m := reportMailModule()
	l := &Listing{ID: "22222222-2222-2222-2222-222222222222", TitleRu: "Сдается Студия"}

	subject, body := m.reportEmail(l, 3, true)
	if !strings.Contains(strings.ToLower(subject), "скрыто") {
		t.Errorf("subject %q does not say the listing is hidden", subject)
	}
	if !strings.Contains(body, "скрыто на проверку") {
		t.Error("the message does not say the listing is hidden for review")
	}
	if strings.Contains(body, "После 3 жалоб") {
		t.Error("a hidden listing should not still be counting down to being hidden")
	}
	if !strings.Contains(body, "/edit") {
		t.Error("the message gives no way to fix the photos")
	}
}
