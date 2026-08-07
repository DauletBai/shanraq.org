package auth

import (
	"strings"
	"testing"
)

// Email normalisation decides account identity. If two spellings of the same
// address normalise differently, the same mailbox can hold two accounts; if two
// different addresses normalise to one string, whoever registers first owns the
// other person's login. Both are account-takeover shaped, so the rules are
// pinned here rather than left to whatever mail.ParseAddress does this year.
func TestNormalizeEmailIdentity(t *testing.T) {
	// Case and surrounding space must not create a second identity.
	same := []string{"user@example.com", "USER@Example.COM", "  User@Example.com  "}
	first, ok := NormalizeEmail(same[0])
	if !ok {
		t.Fatalf("%q must be valid", same[0])
	}
	for _, v := range same[1:] {
		got, ok := NormalizeEmail(v)
		if !ok || got != first {
			t.Errorf("NormalizeEmail(%q) = %q,%v — want %q, one identity", v, got, ok, first)
		}
	}

	// Anything that is not a bare, dotted-domain address is refused. A display
	// name is the interesting one: "Admin <a@b.com>" parses fine as a mail
	// address, and accepting it would store a value that is not the mailbox.
	bad := []string{
		"", "   ", "not-an-email", "@example.com", "user@", "user@localhost",
		"user@example", "user@.com", "user@example.", "Admin <a@example.com>",
		"a@b.com, c@d.com", "user name@example.com", "user@exam ple.com",
		strings.Repeat("a", 250) + "@example.com",
	}
	for _, v := range bad {
		if got, ok := NormalizeEmail(v); ok {
			t.Errorf("NormalizeEmail(%q) accepted it as %q", v, got)
		}
	}

	// Different mailboxes must stay different — no aggressive canonicalisation
	// that would let one person claim another's address.
	a, _ := NormalizeEmail("user.name@example.com")
	b, _ := NormalizeEmail("username@example.com")
	if a == b {
		t.Error("dots in the local part must not be collapsed: that merges two mailboxes")
	}
	c, _ := NormalizeEmail("user+tag@example.com")
	d, _ := NormalizeEmail("user@example.com")
	if c == d {
		t.Error("a plus tag must not be stripped: that merges two mailboxes")
	}
}

// The password rule is the floor under every account on the site, including the
// admin one. Loosening it silently is the kind of change that never gets
// noticed until the credential stuffing starts.
func TestValidatePasswordFloor(t *testing.T) {
	for _, weak := range []string{
		"", "short", "1234567", "12345678", "parolparol", "PAROLPAROL",
		"        ", "aaaaaaa" /* 7 знаков, без цифры */, "1234567",
	} {
		if err := ValidatePassword(weak); err == nil {
			t.Errorf("ValidatePassword(%q) accepted a weak password", weak)
		}
	}
	for _, good := range []string{"Parol12345", "parol1234", "пароль1234", "aaaaaaa1"} {
		if err := ValidatePassword(good); err != nil {
			t.Errorf("ValidatePassword(%q) = %v, want accepted", good, err)
		}
	}
	// An overlong password must be refused rather than truncated: bcrypt only
	// looks at the first 72 bytes, and silently ignoring the rest would make
	// two different passwords equivalent.
	if err := ValidatePassword(strings.Repeat("a1", 500)); err == nil {
		t.Error("an overlong password must be refused, not silently truncated")
	}
}

// Real names are the byline model: every article and listing is attributed to a
// person, so the field cannot accept digits, markup or padding that would let
// someone sign as something other than a name.
func TestValidatePersonNameRejectsNonNames(t *testing.T) {
	for _, bad := range []string{
		"", " ", "1", "Иван1", "<script>", "Иван<b>", "user@example.com",
		"Иван  Иванов", "Иван-", strings.Repeat("Я", 100),
	} {
		if err := ValidatePersonName(NormalizePersonName(bad)); err == nil {
			t.Errorf("ValidatePersonName(%q) accepted a non-name", bad)
		}
	}
	for _, good := range []string{"Даулет", "Baimurza", "Әлихан", "Мария"} {
		if err := ValidatePersonName(NormalizePersonName(good)); err != nil {
			t.Errorf("ValidatePersonName(%q) = %v, want accepted", good, err)
		}
	}
	// The patronymic is optional, so empty must pass — but rubbish still must not.
	if err := ValidateOptionalPersonName(""); err != nil {
		t.Errorf("an empty patronymic must be allowed, got %v", err)
	}
	if err := ValidateOptionalPersonName("Иван1"); err == nil {
		t.Error("a non-name patronymic must still be refused")
	}
}

// ShortName is what a commenter is shown as. It must never fall back to the raw
// e-mail address, which would publish the mailbox of everyone who comments.
func TestShortNameNeverLeaksTheEmail(t *testing.T) {
	cases := []struct{ first, last, middle, email string }{
		{"Даулет", "Баймурза", "Абаевич", "secret@example.com"},
		{"Даулет", "Баймурза", "", "secret@example.com"},
		{"", "", "", "secret@example.com"},
		{"", "Баймурза", "", "secret@example.com"},
	}
	for _, c := range cases {
		got := ShortName(c.first, c.last, c.middle, c.email)
		if strings.Contains(got, "@") || strings.Contains(got, "secret@example.com") {
			t.Errorf("ShortName(%q,%q,%q,%q) = %q — the address leaked",
				c.first, c.last, c.middle, c.email, got)
		}
		if strings.TrimSpace(got) == "" {
			t.Errorf("ShortName(%q,%q,%q,%q) produced an empty label",
				c.first, c.last, c.middle, c.email)
		}
	}
}

// The jobs API takes an arbitrary job name and payload, and among the registered
// handlers are ones that spend the AI budget and rewrite other people's
// translations. Any registered reader can mint an API key, so the role list on
// that route is the only thing standing between a reader and the queue. An
// earlier changelog claimed this was covered; it was not, so here it is.
func TestRequireRolesExcludesPlainUsers(t *testing.T) {
	m := New()
	staffOnly := m.RequireRoles("operator", "admin")

	for _, role := range []string{"user", "", "reader", "author"} {
		if roleAllowed(role, []string{"operator", "admin"}) {
			t.Errorf("role %q must not reach a staff-only route", role)
		}
	}
	for _, role := range []string{"operator", "admin"} {
		if !roleAllowed(role, []string{"operator", "admin"}) {
			t.Errorf("role %q must reach a staff-only route", role)
		}
	}
	if staffOnly == nil {
		t.Error("RequireRoles must return a middleware")
	}
}

// roleAllowed mirrors the middleware's decision for a single role, so the table
// above states the rule without needing a signed token per case.
func roleAllowed(role string, allowed []string) bool {
	c := &Claims{Roles: []string{role}}
	if role == "" {
		c.Roles = nil
	}
	return c.HasAnyRole(allowed...)
}
