package articles

import (
	"time"

	"github.com/google/uuid"
)

// Supported content languages. Author writes in one; the story is published in all three.
const (
	LangKZ = "kz"
	LangRU = "ru"
	LangEN = "en"
)

// Langs is the canonical, ordered list of supported languages.
var Langs = []string{LangKZ, LangRU, LangEN}

// LangLabels maps a language code to its short display label.
var LangLabels = map[string]string{
	LangKZ: "ҚАЗ",
	LangRU: "РУС",
	LangEN: "ENG",
}

// LangNames maps a language code to its full native name.
var LangNames = map[string]string{
	LangKZ: "Қазақша",
	LangRU: "Русский",
	LangEN: "English",
}

// IsLang reports whether code is a supported language.
func IsLang(code string) bool {
	switch code {
	case LangKZ, LangRU, LangEN:
		return true
	default:
		return false
	}
}

// Article is the language-independent story record.
type Article struct {
	ID           uuid.UUID
	AuthorID     uuid.UUID
	AuthorEmail  string
	Slug         string
	OriginalLang string
	Status       string
	Category     string
	CoverURL     string
	Score        int
	ViewsCount   int64
	PublishedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Translations keyed by language code, populated on demand.
	Translations map[string]*Translation
}

// Translation holds one language version of an article.
type Translation struct {
	Lang    string
	Title   string
	Summary string
	BodyMD  string
	Source  string // human | ai
	Status  string // draft | pending | ready
}

// Translation returns the version for lang, falling back to the original
// language, then to any available version. The second result is the language
// actually served.
func (a *Article) Translation(lang string) (*Translation, string) {
	if a.Translations == nil {
		return nil, ""
	}
	if tr, ok := a.Translations[lang]; ok && tr.hasContent() {
		return tr, lang
	}
	if tr, ok := a.Translations[a.OriginalLang]; ok && tr.hasContent() {
		return tr, a.OriginalLang
	}
	for _, code := range Langs {
		if tr, ok := a.Translations[code]; ok && tr.hasContent() {
			return tr, code
		}
	}
	return nil, ""
}

// AvailableLangs returns the languages that currently have readable content.
func (a *Article) AvailableLangs() []string {
	out := make([]string, 0, len(Langs))
	for _, code := range Langs {
		if tr, ok := a.Translations[code]; ok && tr.hasContent() {
			out = append(out, code)
		}
	}
	return out
}

func (t *Translation) hasContent() bool {
	return t != nil && t.Title != "" && t.BodyMD != ""
}

// AuthorName returns a human display name derived from the author's email.
func (a *Article) AuthorName() string {
	return displayName(a.AuthorEmail)
}

// AuthorStats aggregates a single author's publishing activity for the dashboard.
type AuthorStats struct {
	TotalArticles int
	Published     int
	Drafts        int
	TotalViews    int64
	ViewsByLang   map[string]int64
}
