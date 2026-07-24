package auth

import "testing"

func TestValidatePersonName(t *testing.T) {
	good := []string{
		"Даулет",      // Cyrillic
		"Баймурза",    //
		"Әсем",        // Kazakh letter
		"Ұлан",        //
		"Ivan",        // Latin
		"Аль-Фараби",  // one internal hyphen
		"Мухамед-Али", //
		"Ли",          // shortest allowed (2)
	}
	for _, s := range good {
		if err := ValidatePersonName(s); err != nil {
			t.Errorf("ValidatePersonName(%q) = %v, want nil", s, err)
		}
	}

	bad := map[string]string{
		"":            "empty",
		"д":           "single lowercase letter",
		"даулет":      "lowercase start",
		"Даулет1":     "digit",
		"Даулет Али":  "space",
		"Даулет.":     "punctuation",
		"Д":           "too short",
		"-Даулет":     "leading hyphen",
		"Даулет-":     "trailing hyphen",
		"Аль--Фараби": "doubled hyphen",
		"А-Б-В":       "more than one hyphen",
		"Иван@mail":   "symbol",
	}
	for s, why := range bad {
		if err := ValidatePersonName(s); err == nil {
			t.Errorf("ValidatePersonName(%q) = nil, want error (%s)", s, why)
		}
	}

	// The patronymic is optional but follows the same rule when present.
	if err := ValidateOptionalPersonName(""); err != nil {
		t.Errorf("empty patronymic should be allowed, got %v", err)
	}
	if err := ValidateOptionalPersonName("абаевич"); err == nil {
		t.Error("lowercase patronymic should be rejected")
	}
}

func TestValidatePassword(t *testing.T) {
	if err := ValidatePassword("parol123"); err != nil {
		t.Errorf("valid password rejected: %v", err)
	}
	for _, pw := range []string{"short1", "onlyletters", "12345678", ""} {
		if err := ValidatePassword(pw); err == nil {
			t.Errorf("ValidatePassword(%q) = nil, want error", pw)
		}
	}
}

func TestShortAndFullName(t *testing.T) {
	if got := ShortName("Даулет", "Баймурза", "Абаевич", "x@y.z"); got != "Баймурза Д.А." {
		t.Errorf("ShortName = %q, want %q", got, "Баймурза Д.А.")
	}
	if got := ShortName("Даулет", "Баймурза", "", "x@y.z"); got != "Баймурза Д." {
		t.Errorf("ShortName without patronymic = %q, want %q", got, "Баймурза Д.")
	}
	// Legacy account with no name falls back to the e-mail local part.
	if got := ShortName("", "", "", "someone@mail.kz"); got != "someone" {
		t.Errorf("ShortName fallback = %q, want %q", got, "someone")
	}
	if got := FullName("Даулет", "Баймурза", "x@y.z"); got != "Даулет Баймурза" {
		t.Errorf("FullName = %q, want %q", got, "Даулет Баймурза")
	}
	if got := FullName("", "", "someone@mail.kz"); got != "someone" {
		t.Errorf("FullName fallback = %q, want %q", got, "someone")
	}
}
