package articles

import (
	"strings"
	"testing"

	"shanraq.org/internal/config"
)

func TestApplyOperator_Configured(t *testing.T) {
	// Fake fixture values — the real operator BIN/name/address must never live
	// in the repo (they come from the environment on the live site).
	op := config.OperatorConfig{
		LegalName:   "ТОО «Тест Компани»",
		LegalNameEN: "Test Company LLP",
		BIN:         "123456789012",
		Address:     "г. Тестоград, микрорайон 1, дом 2",
		AddressEN:   "Testville, Microdistrict 1, bldg 2",
		Email:       "support@shanraq.org",
	}

	ru := applyOperator("Оператор — {{operator_block}}", op, LangRU)
	for _, want := range []string{"ТОО «Тест Компани»", "БИН 123456789012", "г. Тестоград", "support@shanraq.org", "Для обращений:"} {
		if !strings.Contains(ru, want) {
			t.Errorf("RU block missing %q in %q", want, ru)
		}
	}
	if strings.Contains(ru, "{{") {
		t.Errorf("RU still has an unsubstituted token: %q", ru)
	}

	// EN prefers the English legal name, address, and labels.
	en := applyOperator("Operator: {{operator_block}} {{op}}", op, LangEN)
	for _, want := range []string{"Test Company LLP", "BIN 123456789012", "Testville, Microdistrict 1", "Contact for enquiries:"} {
		if !strings.Contains(en, want) {
			t.Errorf("EN block missing %q in %q", want, en)
		}
	}
	if strings.Contains(en, "ТОО") || strings.Contains(en, "БИН") {
		t.Errorf("EN block leaked Russian text: %q", en)
	}

	// KZ uses the БСН label and the (shared) legal name.
	kz := applyOperator("{{operator_block}}", op, LangKZ)
	if !strings.Contains(kz, "БСН 123456789012") || !strings.Contains(kz, "Хабарласу үшін:") {
		t.Errorf("KZ block wrong: %q", kz)
	}
}

func TestApplyOperator_EmptyFallsBackWithoutLeakingLabels(t *testing.T) {
	// Nothing configured (the public-repo default): a generic owner phrase and
	// the default support email, with no empty "БИН"/"BIN" label dangling.
	out := applyOperator("Оператором данных является {{operator_block}}", config.OperatorConfig{}, LangRU)
	if !strings.Contains(out, "владелец платформы Shanraq") {
		t.Errorf("expected generic fallback name, got %q", out)
	}
	if !strings.Contains(out, "support@shanraq.org") {
		t.Errorf("expected default contact email, got %q", out)
	}
	if strings.Contains(out, "БИН") || strings.Contains(out, "{{") {
		t.Errorf("empty config must not print a BIN label or leave tokens: %q", out)
	}
}

func TestApplyOperator_PartialOmitsMissingFields(t *testing.T) {
	// Only a legal name + phone: no BIN and no address must appear.
	op := config.OperatorConfig{LegalName: "ТОО «Тест»", Phone: "+7 700 000 00 00"}
	out := applyOperator("{{operator_block}}", op, LangRU)
	if strings.Contains(out, "БИН") {
		t.Errorf("unset BIN should be omitted, got %q", out)
	}
	if !strings.Contains(out, "ТОО «Тест»") || !strings.Contains(out, "+7 700 000 00 00") {
		t.Errorf("name and phone should be present, got %q", out)
	}
}
