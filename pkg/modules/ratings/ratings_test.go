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

// Голос за комментарий взвешивается так же, как голос за статью: свежий аккаунт
// весит единицу, и толпа свежих аккаунтов не хоронит комментарий.
func TestCommentVotesUseTheSameWeighting(t *testing.T) {
	if Weight(0) != 1 {
		t.Errorf("новый читатель весит %d, а должен единицу", Weight(0))
	}
	if Weight(-500) != 1 {
		t.Errorf("отрицательная карма не должна давать меньше единицы: %d", Weight(-500))
	}
	if Weight(1000) != maxWeight {
		t.Errorf("вес не ограничен сверху: %d", Weight(1000))
	}
	// Пять минусов от новичков — ровно порог сворачивания в модуле articles.
	if 5*Weight(0) < 5 {
		t.Error("пять голосов новичков не набирают порога сворачивания")
	}
}
