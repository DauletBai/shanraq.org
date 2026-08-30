package articles

import (
	"testing"

	"go.uber.org/zap"
)

// A database that blinks should cost a delay, not a hole in the figures. The
// counts a failed batch never wrote go back into the buffer, and the next flush
// carries them.
func TestGiveBackReturnsUnwrittenCounts(t *testing.T) {
	mt := &Metrics{buf: map[metricKey]int64{}, log: zap.NewNop()}
	mt.buf[metricKey{"page", "article", true}] = 2

	mt.giveBack([]metricDelta{
		{metricKey{"page", "article", true}, 5},
		{metricKey{"page", "home", true}, 3},
	})

	if got := mt.buf[metricKey{"page", "article", true}]; got != 7 {
		t.Errorf("returned count did not add to what was already buffered: got %d, want 7", got)
	}
	if got := mt.buf[metricKey{"page", "home", true}]; got != 3 {
		t.Errorf("a count with no counterpart in the buffer was lost: got %d, want 3", got)
	}
	if mt.Dropped() != 0 {
		t.Errorf("nothing should have been dropped, Dropped() = %d", mt.Dropped())
	}
}

// A database down for hours must not be answered by growing a map until the
// process dies of it. Past the ceiling the counts are dropped and said to be.
func TestGiveBackStopsAtTheCeiling(t *testing.T) {
	mt := &Metrics{buf: map[metricKey]int64{}, log: zap.NewNop()}
	for i := 0; i < metricsBufMax; i++ {
		mt.buf[metricKey{"page", string(rune('a'+i%26)) + string(rune(i)), true}] = 1
	}
	mt.giveBack([]metricDelta{{metricKey{"page", "one-too-many", true}, 4}})

	if len(mt.buf) > metricsBufMax {
		t.Errorf("buffer grew past its ceiling: %d entries, max %d", len(mt.buf), metricsBufMax)
	}
	if mt.Dropped() == 0 {
		t.Error("counts were dropped without being counted as dropped")
	}
}
