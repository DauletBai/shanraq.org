package articles

import (
	"testing"
	"unicode/utf8"
)

func TestClipKeepsUTF8Valid(t *testing.T) {
	ru := "Kimi K3: крупнейшая открытая ИИ-модель в истории — и что гонка открытых моделей значит для нас и для рынка труда в Казахстане сегодня"
	for n := 1; n <= 200; n++ {
		got := clip(ru, n)
		if !utf8.ValidString(got) {
			t.Fatalf("n=%d produced invalid UTF-8: %q", n, got)
		}
		if utf8.RuneCountInString(got) > n+1 { // +1 for the ellipsis
			t.Fatalf("n=%d too long: %d runes", n, utf8.RuneCountInString(got))
		}
	}
	if clip("короткая строка", 100) != "короткая строка" {
		t.Fatal("short string must pass through unchanged")
	}
}
