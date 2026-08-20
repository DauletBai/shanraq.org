package articles

import "testing"

// Свёрнутый — не то же самое, что скрытый. Заминусованный комментарий остаётся
// на странице за одной строкой, и любой читатель может его раскрыть. Раньше на
// этом месте стояла модель, которая решала, увидит его кто-нибудь или нет.
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

// Порог не должен уехать в ноль незаметно: комментарий, свёрнутый первым же
// минусом, — это не читательская модерация, а глушилка.
func TestFoldThresholdLeavesRoomForDisagreement(t *testing.T) {
	if commentCollapseScore > -3 {
		t.Errorf("порог %d слишком мягкий: одиночное несогласие прячет комментарий", commentCollapseScore)
	}
}
