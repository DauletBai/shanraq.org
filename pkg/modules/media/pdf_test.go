package media

import "testing"

// A PDF that runs on open is a program with a page count. These are the shapes
// that make one, written the way a real file writes them.
func TestActivePDFNamesWhatMakesADocumentDangerous(t *testing.T) {
	cases := []struct {
		name string
		body string
		want bool
	}{
		{"plain pages", "%PDF-1.4\n1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj\n%%EOF", false},
		{"script on open", "%PDF-1.4\n<</OpenAction<</S/JavaScript/JS(app.alert\\(1\\))>>>>\n%%EOF", true},
		{"script alone", "%PDF-1.4\n<</Names<</JavaScript 5 0 R>>>>\n%%EOF", true},
		{"hex-escaped name", "%PDF-1.4\n<</J#61vaScript 5 0 R>>\n%%EOF", true},
		{"additional actions", "%PDF-1.4\n<</AA<</O 6 0 R>>>>\n%%EOF", true},
		{"launches a program", "%PDF-1.4\n<</S/Launch/F(cmd.exe)>>\n%%EOF", true},
		{"carries another file", "%PDF-1.4\n<</Type/EmbeddedFile/Subtype/application#2Fzip>>\n%%EOF", true},
		{"submits a form", "%PDF-1.4\n<</S/SubmitForm/F(https://elsewhere.example/collect)>>\n%%EOF", true},
		{"fetches a remote page", "%PDF-1.4\n<</S/GoToR/F(https://elsewhere.example/x.pdf)>>\n%%EOF", true},
	}
	for _, c := range cases {
		got := activePDF([]byte(c.body)) != ""
		if got != c.want {
			t.Errorf("%s: refused=%v, want %v (reason %q)", c.name, got, c.want, activePDF([]byte(c.body)))
		}
	}
}

// The reason is shown to the author, so it has to name the thing rather than
// say "invalid".
func TestTheRefusalSaysWhatWasFound(t *testing.T) {
	got := activePDF([]byte("%PDF-1.4\n<</S/Launch/F(cmd.exe)>>"))
	if got != "a program launch" {
		t.Errorf("reason = %q, want it to name the launch", got)
	}
}

// An ordinary exported document must not be refused, or authors learn to work
// around the check instead of with it.
func TestAnOrdinaryDocumentIsAccepted(t *testing.T) {
	doc := "%PDF-1.7\n" +
		"1 0 obj<</Type/Catalog/Pages 2 0 R/Lang(ru-RU)>>endobj\n" +
		"2 0 obj<</Type/Pages/Kids[3 0 R]/Count 1>>endobj\n" +
		"3 0 obj<</Type/Page/Parent 2 0 R/Resources<</Font<</F1 4 0 R>>>>/Contents 5 0 R>>endobj\n" +
		"4 0 obj<</Type/Font/Subtype/Type1/BaseFont/Helvetica>>endobj\n" +
		"5 0 obj<</Length 44>>stream\nBT /F1 12 Tf 72 720 Td (Отчёт за август) Tj ET\nendstream endobj\n" +
		"trailer<</Root 1 0 R>>\n%%EOF"
	if reason := activePDF([]byte(doc)); reason != "" {
		t.Errorf("a plain document was refused for %q", reason)
	}
}
