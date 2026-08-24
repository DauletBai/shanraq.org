package articles

import "testing"

// Folded is not the same as hidden. A downvoted comment stays on the page behind a
// single row, and any reader can unfold it. What used to stand here was a model that
// decided whether anyone would see it.
func TestCommentFoldsOnlyWhenReadersBuryIt(t *testing.T) {
	cases := []struct {
		score int
		fold  bool
	}{
		{10, false},
		{0, false},
		{-1, false},
		{-4, false},
		{commentCollapseScore, true},
		{commentCollapseScore - 1, true},
	}
	for _, c := range cases {
		if got := (Comment{Score: c.score}).Collapsed(); got != c.fold {
			t.Errorf("счёт %d: свёрнут=%v, ожидалось %v", c.score, got, c.fold)
		}
	}
}

// The threshold must not slip to zero unnoticed: a comment folded by the very first
// downvote is not reader moderation, it is a silencer.
func TestFoldThresholdLeavesRoomForDisagreement(t *testing.T) {
	if commentCollapseScore > -3 {
		t.Errorf("порог %d слишком мягкий: одиночное несогласие прячет комментарий", commentCollapseScore)
	}
}
