package notifier

import "testing"

func TestEnvelopeSender(t *testing.T) {
	cases := map[string]string{
		"Shanraq <no-reply@shanraq.org>": "no-reply@shanraq.org",
		"no-reply@shanraq.org":           "no-reply@shanraq.org",
		"  spaced@shanraq.org  ":         "spaced@shanraq.org",
		"not-an-address":                 "not-an-address", // left as-is (best effort)
	}
	for in, want := range cases {
		if got := envelopeSender(in); got != want {
			t.Errorf("envelopeSender(%q) = %q, want %q", in, got, want)
		}
	}
}
