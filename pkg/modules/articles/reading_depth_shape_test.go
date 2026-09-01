package articles

import "testing"

// A funnel that widens as it deepens is not a finding, it is a lost beacon.
// The panel showed three readers at three quarters and two at half, which no
// reader can produce: nobody passes 75% of a page without passing 50%.
func TestMonotonicDepth(t *testing.T) {
	cases := []struct {
		name string
		in   map[int]int64
		want map[int]int64
	}{
		{
			"дошедших до ¾ больше, чем до ½ — потерянный маяк",
			map[int]int64{25: 2, 50: 2, 75: 3, 100: 2},
			map[int]int64{25: 3, 50: 3, 75: 3, 100: 2},
		},
		{
			"нормальная воронка не меняется",
			map[int]int64{25: 10, 50: 9, 75: 9, 100: 5},
			map[int]int64{25: 10, 50: 9, 75: 9, 100: 5},
		},
		{
			"дочитали больше, чем начали",
			map[int]int64{25: 0, 50: 0, 75: 0, 100: 4},
			map[int]int64{25: 4, 50: 4, 75: 4, 100: 4},
		},
		{
			"пустая воронка остаётся пустой",
			map[int]int64{},
			map[int]int64{25: 0, 50: 0, 75: 0, 100: 0},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := monotonicDepth(c.in)
			for _, mark := range []int{25, 50, 75, 100} {
				if got[mark] != c.want[mark] {
					t.Errorf("на %d%%: получили %d, хотели %d", mark, got[mark], c.want[mark])
				}
			}
			// The property itself, not just the examples.
			if got[25] < got[50] || got[50] < got[75] || got[75] < got[100] {
				t.Errorf("воронка расширяется вглубь: %v", got)
			}
		})
	}
}
