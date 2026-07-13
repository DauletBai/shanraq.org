package ratings

import "testing"

func TestWeight(t *testing.T) {
	cases := map[int]int{
		-500: 1, // negative karma still counts as 1
		0:    1,
		50:   1,
		100:  2,
		250:  3,
		1000: 5,
		9999: 5, // capped
	}
	for karma, want := range cases {
		if got := Weight(karma); got != want {
			t.Errorf("Weight(%d) = %d, want %d", karma, got, want)
		}
	}
}
