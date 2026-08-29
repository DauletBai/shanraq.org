package articles

import "testing"

// A fifteen-minute article: 2700-odd words at the 180 a minute the page prints
// under the title.
const fifteenMin = 15 * 60

func TestReadCounts(t *testing.T) {
	cases := []struct {
		name   string
		depth  int
		secs   int
		expect int
		want   bool
	}{
		{"flicked to the bottom in four seconds", 100, 4, fifteenMin, false},
		{"read to the end in nine minutes", 100, 540, fifteenMin, true},
		{"nine minutes but stopped halfway", 50, 540, fifteenMin, false},
		{"exactly half the estimate, at the end", 100, fifteenMin / 2, fifteenMin, true},
		{"a second short of half", 100, fifteenMin/2 - 1, fifteenMin, false},
		{"no estimate to measure against", 100, 600, 0, false},
		{"a short article read quickly", 100, 40, 60, true},
	}
	for _, c := range cases {
		if got := readCounts(c.depth, c.secs, c.expect); got != c.want {
			t.Errorf("%s: readCounts(%d, %ds, %ds) = %v, want %v",
				c.name, c.depth, c.secs, c.expect, got, c.want)
		}
	}
}

// The estimate the threshold is measured against is the one the article shows
// its reader, so the two can never drift apart.
func TestReadingMinutesMatchesWhatThePageClaims(t *testing.T) {
	body := ""
	for i := 0; i < 2700; i++ {
		body += "слово "
	}
	if got, want := readingMinutes(body), 15; got != want {
		t.Errorf("readingMinutes(2700 words) = %d, want %d", got, want)
	}
}
