package articles

import (
	"errors"
	"strings"
	"testing"
)

func errorsAs(err error, target any) bool { return errors.As(err, target) }

// Three languages are the norm; a place-targeted piece may stop at two.
func TestRequiredLanguagesByReach(t *testing.T) {
	local := requiredLanguages(true)
	if len(local) != 2 || !hasLang(local, LangKZ) || !hasLang(local, LangRU) {
		t.Errorf("для местного материала ожидались казахский и русский, вышло %v", local)
	}
	if hasLang(local, LangEN) {
		t.Error("английский не должен быть обязательным для местного материала")
	}
	all := requiredLanguages(false)
	for _, l := range []string{LangKZ, LangRU, LangEN} {
		if !hasLang(all, l) {
			t.Errorf("материал на всю страну обязан выходить на %q", l)
		}
	}
}

// The refusal names the versions that are missing: "not all languages" sends an
// author hunting through their own piece for the gap.
func TestRefusalNamesTheMissingLanguages(t *testing.T) {
	e := &MissingLanguagesError{Missing: []string{LangEN}, Placed: false}
	msg := e.Reason(LangRU)
	if !strings.Contains(msg, "English") {
		t.Errorf("в отказе не назван недостающий язык: %q", msg)
	}
	if !strings.Contains(msg, "трёх") {
		t.Errorf("отказ для материала на всю страну должен говорить про три языка: %q", msg)
	}

	loc := &MissingLanguagesError{Missing: []string{LangKZ}, Placed: true}
	m2 := loc.Reason(LangRU)
	if !strings.Contains(m2, "Қазақша") {
		t.Errorf("не назван недостающий язык: %q", m2)
	}
	if strings.Contains(m2, "трёх") {
		t.Errorf("для местного материала три языка не требуются: %q", m2)
	}

	// A malformed list must still say something a person can act on.
	if s := (&MissingLanguagesError{Missing: []string{""}}).Reason(LangRU); strings.TrimSpace(s) == "" {
		t.Error("пустой список языков дал пустое сообщение")
	}
	// And the typed error must answer to the sentinel the handlers test for.
	var target *MissingLanguagesError
	if !errorsAs(e, &target) {
		t.Error("типизированная ошибка не опознаётся как ErrLocalNeedsBothLanguages")
	}
}

// Language quality is checked but never blocks: a misplaced comma is a reason
// to tell an author, not a reason to stop a piece that is otherwise sound.
func TestSpellingRulesExistAndDoNotBlock(t *testing.T) {
	for _, code := range []string{"spelling", "grammar", "punctuation", "local_languages"} {
		if !isReviewRule(code) {
			t.Errorf("правило %q не объявлено", code)
		}
		for _, lang := range []string{LangKZ, LangRU, LangEN} {
			if strings.TrimSpace(T(lang, "rule."+code)) == "" {
				t.Errorf("нет перевода rule.%s для %q", code, lang)
			}
		}
	}
	prompt := reviewPrompt(LangRU, "Заголовок", "Описание", "Текст")
	for _, want := range []string{"spelling", "grammar", "punctuation"} {
		if !strings.Contains(prompt, want) {
			t.Errorf("проверяющий не знает о правиле %q", want)
		}
	}
	// The checker must not be told to block on a typo, and must not be invited
	// to argue about style.
	if strings.Contains(prompt, `"block" only for defamation, hatred, personal_data, illegal, plagiarism, disguised_ad, spelling`) {
		t.Error("орфография попала в список блокирующих правил")
	}
	if !strings.Contains(prompt, "Do not report style") {
		t.Error("проверяющему не запрещено спорить о стиле")
	}
}

func hasLang(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}
