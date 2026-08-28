package articles

import "testing"

func TestSpeechTextDropsMarksAndKeepsWords(t *testing.T) {
	cases := [][2]string{
		{"Токаев (по данным ЦИК) получил 71,17%", "Токаев по данным ЦИК получил 71,17%"},
		{"Закон «О возврате активов» действует", "Закон О возврате активов действует"},
		{"Смотрите https://example.kz/page — там всё", "Смотрите , там всё"},
		{"Сумма 16 700 тенге", "Сумма 16700 тенге"},
		{"Итого 55 800 000 тенге", "Итого 55800000 тенге"},
		{"Источник · дата · автор", "Источник, дата, автор"},
		{"Это **важно** и `код`", "Это важно и код"},
		{"Кто-то из ИТ-отдела", "Кто-то из ИТ-отдела"},
	}
	for _, c := range cases {
		if got := speechText(c[0]); got != c[1] {
			t.Errorf("speechText(%q) = %q, want %q", c[0], got, c[1])
		}
	}
}

func TestNarrationBlocksFollowsTheDOM(t *testing.T) {
	got := NarrationBlocks(`<div class="prose">
	  <p>Первый абзац.</p>
	  <pre><code>not spoken</code></pre>
	  <blockquote><p>Внутренний абзац читается один раз.</p></blockquote>
	  <p aria-hidden="true">Скрытое</p>
	  <ul><li>Пункт списка</li></ul>
	</div>`)
	want := []string{"Первый абзац.", "Внутренний абзац читается один раз.", "Пункт списка"}
	if len(got) != len(want) {
		t.Fatalf("got %d blocks %q, want %d", len(got), got, len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("block %d = %q, want %q", i, got[i], want[i])
		}
	}
}
