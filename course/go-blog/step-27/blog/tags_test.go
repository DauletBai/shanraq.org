package blog

import (
	"slices"
	"testing"
)

func TestParseTags(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"бос", "", nil},
		{"біреу", "дала", []string{"дала"}},
		{"шеттегі бос орындар", "  дала  ", []string{"дала"}},
		{"әртүрлі регистр", "Дала, дала", []string{"дала"}},
		{"үтірлер арасы бос", "дала,,блог", []string{"дала", "блог"}},
		{"рет сақталады", "блог, дала", []string{"блог", "дала"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ParseTags(c.in)
			if !slices.Equal(got, c.want) {
				t.Errorf("ParseTags(%q) = %v, күткеніміз %v", c.in, got, c.want)
			}
		})
	}
}
