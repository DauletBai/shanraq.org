package ai

import (
	"strings"
	"testing"
)

// The first live translation came back with the link label still in the source
// language — `[ссылкой]` sitting inside a Kazakh sentence — because the prompt
// told the model to preserve the Markdown "exactly" and, with no room to think
// it over, it took that to cover the words inside the markup too.
//
// The distinction the prompt has to make: markup and URLs are structure, the
// words a reader sees are prose.
func TestTranslatePromptSeparatesMarkupFromProse(t *testing.T) {
	p := translateSystem("ru", "kz")

	for _, want := range []string{
		"link labels",   // the case that actually failed
		"[brackets]",    // named concretely, not left to inference
		"byte for byte", // what preservation applies to
		"(parentheses)", // URLs are not prose
		"code",          // nor is code
	} {
		if !strings.Contains(p, want) {
			t.Errorf("the prompt never mentions %q:\n%s", want, p)
		}
	}

	// The old wording is what caused it; it must not come back.
	if strings.Contains(p, "Markdown formatting exactly") {
		t.Error("the prompt still tells the model to preserve the formatting 'exactly', " +
			"which is what it read as 'do not touch the words inside it'")
	}

	// Both language names have to survive the format string — a %[2]s typo
	// would silently produce a prompt naming one language twice.
	if !strings.Contains(p, langFullName["ru"]) || !strings.Contains(p, langFullName["kz"]) {
		t.Errorf("the prompt does not name both languages:\n%s", p)
	}
	if strings.Count(p, langFullName["kz"]) < 2 {
		t.Error("the target language should be named where the prompt asks for idiomatic writing")
	}
}
