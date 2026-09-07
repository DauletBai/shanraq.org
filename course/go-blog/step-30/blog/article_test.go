package blog

import "testing"

func TestReadingTime(t *testing.T) {
	cases := []struct {
		name  string
		words int
		want  int
	}{
		{"нөл", 0, 0},
		{"бір сөз", 1, 1},
		{"дәл екі жүз", 200, 1},
		{"екі жүз бір", 201, 2},
		{"дәл төрт жүз", 400, 2},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := Article{Words: c.words}
			if got := a.ReadingTime(); got != c.want {
				t.Errorf("ReadingTime(%d) = %d, күткеніміз %d", c.words, got, c.want)
			}
		})
	}
}

func TestLetters(t *testing.T) {
	cases := []struct {
		name  string
		title string
		want  int
	}{
		{"латын", "Go", 2},
		{"қазақша", "Шаңырақ", 7},
		{"бос", "", 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := Article{Title: c.title}
			if got := a.Letters(); got != c.want {
				t.Errorf("Letters(%q) = %d, күткеніміз %d", c.title, got, c.want)
			}
		})
	}
}
