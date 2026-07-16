package auth

import "testing"

func TestNormalizeEmail(t *testing.T) {
	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		{"  User@Example.COM ", "user@example.com", true},
		{"a@b.co", "a@b.co", true},
		{"", "", false},
		{"   ", "", false},
		{"not-an-email", "", false},
		{"missing@domain", "", false},
		{"two@@at.com", "", false},
		{"name <name@x.com>", "", false}, // display-name form must be rejected
		{"sp ace@x.com", "", false},
	}
	for _, c := range cases {
		got, ok := NormalizeEmail(c.in)
		if ok != c.ok || got != c.want {
			t.Errorf("NormalizeEmail(%q) = (%q,%v), want (%q,%v)", c.in, got, ok, c.want, c.ok)
		}
	}
}
