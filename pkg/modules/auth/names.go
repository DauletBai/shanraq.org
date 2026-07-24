package auth

import (
	"errors"
	"strings"
	"unicode"
)

// Errors for the person-name and password rules, so callers can show a precise
// message instead of a generic "invalid input".
var (
	ErrNameEmpty    = errors.New("name is required")
	ErrNameFormat   = errors.New("name must be letters only, starting with a capital")
	ErrNameLength   = errors.New("name length out of range")
	ErrPasswordWeak = errors.New("password must be at least 8 characters and contain a letter and a digit")
)

const (
	minNameLen = 2
	maxNameLen = 40
)

// NormalizePersonName trims surrounding whitespace. It deliberately does not
// change letter case: the capital-letter rule is validated, not silently
// "fixed", so the user sees and confirms exactly what will be shown publicly.
func NormalizePersonName(s string) string { return strings.TrimSpace(s) }

// ValidatePersonName checks a given name, family name or patronymic. The rule:
// letters only, starting with a capital, no digits, no spaces, no punctuation —
// with one exception, a single internal hyphen, because real Kazakh and Russian
// names use it ("Аль-Фараби", "Мухамед-Али"). Any script is accepted, so Kazakh
// letters (ә, ғ, қ, ң, ө, ұ, ү, і) and Latin both pass.
func ValidatePersonName(s string) error {
	if s == "" {
		return ErrNameEmpty
	}
	r := []rune(s)
	if len(r) < minNameLen || len(r) > maxNameLen {
		return ErrNameLength
	}
	if !unicode.IsUpper(r[0]) {
		return ErrNameFormat
	}
	hyphens := 0
	for i, c := range r {
		if unicode.IsLetter(c) {
			continue
		}
		// One hyphen, never leading or trailing, never doubled.
		if c == '-' && i > 0 && i < len(r)-1 && r[i-1] != '-' && hyphens == 0 {
			hyphens++
			continue
		}
		return ErrNameFormat
	}
	return nil
}

// ValidateOptionalPersonName applies the same rule but accepts an empty value,
// for the patronymic which not everyone has.
func ValidateOptionalPersonName(s string) error {
	if s == "" {
		return nil
	}
	return ValidatePersonName(s)
}

// ValidatePassword enforces the signup password rule: long enough, and mixing
// letters with digits so a password is not a bare word or a bare number.
func ValidatePassword(pw string) error {
	if len(pw) < minPasswordLen || len(pw) > maxPasswordLen {
		return ErrPasswordWeak
	}
	var hasLetter, hasDigit bool
	for _, c := range pw {
		switch {
		case unicode.IsLetter(c):
			hasLetter = true
		case unicode.IsDigit(c):
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return ErrPasswordWeak
	}
	return nil
}

// ShortName renders "Фамилия И.О." — the family name plus initials. Used where
// a compact, formal attribution is wanted (comments), and where showing the
// full given name of every commenter would be needlessly exposing.
func ShortName(first, last, middle, email string) string {
	first, last, middle = strings.TrimSpace(first), strings.TrimSpace(last), strings.TrimSpace(middle)
	initial := func(s string) string {
		for _, c := range s { // first rune, uppercased
			return string(unicode.ToUpper(c)) + "."
		}
		return ""
	}
	if last == "" {
		// No family name yet (legacy account): fall back to the long form.
		return FullName(first, last, email)
	}
	out := last
	if in := initial(first); in != "" {
		out += " " + in
		if mid := initial(middle); mid != "" {
			out += mid
		}
	}
	return out
}

// FullName renders a user's public display name from its parts, falling back to
// the local part of their e-mail for accounts created before names were
// required. Patronymic is not shown in the short form.
func FullName(first, last, email string) string {
	first, last = strings.TrimSpace(first), strings.TrimSpace(last)
	switch {
	case first != "" && last != "":
		return first + " " + last
	case first != "":
		return first
	case last != "":
		return last
	}
	if i := strings.IndexByte(email, '@'); i > 0 {
		return email[:i]
	}
	return email
}
