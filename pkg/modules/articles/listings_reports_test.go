package articles

import "testing"

func TestShouldAutoHide(t *testing.T) {
	cases := []struct {
		name    string
		reports int
		views   int
		want    bool
	}{
		{"below floor, tiny sample", 2, 2, false},               // 100% but only 2 reports
		{"floor met, high share", 3, 10, true},                  // 30% >= 15%
		{"floor met, low share on big audience", 3, 100, false}, // 3% < 15%
		{"exactly at percent", 3, 20, true},                     // 15%
		{"many reports huge audience below %", 10, 1000, false}, // 1%
		{"hard count overrides low share", 12, 1000, true},      // >=12 always
		{"zero views but enough reports", 3, 0, true},           // floored denom
		{"one report", 1, 1, false},
	}
	for _, c := range cases {
		if got := shouldAutoHide(c.reports, c.views); got != c.want {
			t.Errorf("%s: shouldAutoHide(%d,%d)=%v want %v", c.name, c.reports, c.views, got, c.want)
		}
	}
}
