package site

import "net/http"

// Supported interface languages. An author writes in one; the story is
// published in all three.
const (
	LangKZ = "kz"
	LangRU = "ru"
	LangEN = "en"
)

// Langs is the canonical, ordered list of supported languages.
var Langs = []string{LangKZ, LangRU, LangEN}

// LangLabels maps a language code to its short display label.
var LangLabels = map[string]string{
	LangKZ: "kz",
	LangRU: "ru",
	LangEN: "en",
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

// LangCookie remembers the reader's choice between visits.
const LangCookie = "shanraq_lang"

// ResolveLang picks the active language from ?lang=, then the cookie, then the
// default. A language given in the query is also remembered, so a link shared
// in Kazakh keeps the reader in Kazakh on the next page.
func ResolveLang(w http.ResponseWriter, r *http.Request) string {
	if q := r.URL.Query().Get("lang"); IsLang(q) {
		http.SetCookie(w, &http.Cookie{Name: LangCookie, Value: q, Path: "/", MaxAge: 31536000, SameSite: http.SameSiteLaxMode})
		return q
	}
	if c, err := r.Cookie(LangCookie); err == nil && IsLang(c.Value) {
		return c.Value
	}
	return LangRU
}

// HTMLLang maps our internal UI code to the BCP-47 language subtag used in
// HTML lang / hreflang / schema.org. Kazakh's ISO 639-1 code is "kk"; we keep
// "kz" internally (routing, ?lang=) but must present "kk" to browsers/crawlers.
func HTMLLang(lang string) string {
	if lang == LangKZ {
		return "kk"
	}
	return lang
}

// OGLocale maps a UI language to an Open Graph locale.
func OGLocale(lang string) string {
	switch lang {
	case LangKZ:
		return "kk_KZ"
	case LangEN:
		return "en_US"
	default:
		return "ru_RU"
	}
}

// CanonURL builds the canonical relative URL for a page in one language,
// preserving the whitelisted filters so /?cat=sport canonicalizes to itself
// (with its category), not to a bare "/".
func CanonURL(path, filters, lang string) string {
	q := "lang=" + lang
	if filters != "" {
		q = filters + "&" + q
	}
	return path + "?" + q
}

// LangLinks builds the per-language alternates for the current page, carrying
// the same whitelisted filters so switching language keeps the category/filter.
func LangLinks(base, filters string) map[string]string {
	out := make(map[string]string, len(Langs))
	for _, l := range Langs {
		out[l] = CanonURL(base, filters, l)
	}
	return out
}
