package auth

import "testing"

func TestNormalizePhone(t *testing.T) {
	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		{"+7 701 222 33 44", "+77012223344", true},
		{"87012223344", "87012223344", true},
		{"+7 (701) 222-33-44", "+77012223344", true},
		{"", "", false},
		{"12345", "", false},             // too short
		{"+1234567890123456", "", false}, // too long (16 digits)
		{"7017abc4444", "", false},       // letters
		{"+7-701-ABC", "", false},
	}
	for _, c := range cases {
		got, ok := NormalizePhone(c.in)
		if ok != c.ok || got != c.want {
			t.Errorf("NormalizePhone(%q) = (%q,%v), want (%q,%v)", c.in, got, ok, c.want, c.ok)
		}
	}
}

func TestGenerateNumericCode(t *testing.T) {
	c, err := generateNumericCode(6)
	if err != nil {
		t.Fatalf("gen: %v", err)
	}
	if len(c) != 6 {
		t.Fatalf("len = %d, want 6", len(c))
	}
	for _, r := range c {
		if r < '0' || r > '9' {
			t.Fatalf("non-digit in code: %q", c)
		}
	}
}
