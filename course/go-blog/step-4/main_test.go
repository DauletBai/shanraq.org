// Step 4 — after the lesson "Your first test".
//
// The same blog, now with the checks written down. They are the edges, not the
// middle: nothing, one, exactly on the boundary. The reading time in
// particular is the kind of quiet mistake nobody notices by looking.
package main

import "testing"

func TestReadingTime(t *testing.T) {
	cases := map[int]int{0: 0, 1: 1, 199: 1, 200: 1, 201: 2, 400: 2, 1000: 5}
	for words, want := range cases {
		if got := readingTime(words); got != want {
			t.Errorf("%d сөз: алдық %d, күттік %d", words, got, want)
		}
	}
}

func TestLonger(t *testing.T) {
	got := longer([]string{"Дала", "Шаңырақ", "Go", "Көш"}, 3)
	if len(got) != 2 {
		t.Errorf("алдық %d тақырып, күттік 2: %v", len(got), got)
	}
	if len(longer(nil, 3)) != 0 {
		t.Errorf("бос тізім бос емес нәтиже берді")
	}
}
