package articles

import (
	"net/http"
	"strings"
	"testing"
)

// Страница аналитики публичная: она существует ровно затем, чтобы её видел
// посторонний. Если она потребует входа, смысла в ней нет.
func TestTheAudiencePageIsOpenToEveryone(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	w := app.do(http.MethodGet, "/analytics", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("страница аналитики отдала %d без входа", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, T(LangRU, "stats.title")) {
		t.Error("на странице нет её собственного заголовка")
	}
	// Метод объявляется на самой странице: число без определения — это то, с
	// чем спорят наши же статьи.
	if !strings.Contains(body, T(LangRU, "stats.method")) {
		t.Error("на странице не сказано, как мы считаем")
	}
	// Краулеры показаны отдельно, а не подмешаны к людям: их больше, и
	// смешение завысило бы посещаемость вдвое.
	if !strings.Contains(body, T(LangRU, "stats.bots_note")) {
		t.Error("нет оговорки о том, что краулеры не входят в число людей")
	}
}

// Ни одна строка страницы не должна указывать на человека. Мы не пишем
// конкретные адреса — только тип страницы, — и страница не должна пересекать
// страну ни с чем.
func TestTheAudiencePageNamesNoReaderAndNoArticle(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	author := app.createUser("statspriv@example.com", "Parol123!")
	_, slug := app.seedArticle(author, "published")
	app.exec(`UPDATE articles SET published_at = NOW() WHERE slug = $1`, slug)

	body := app.do(http.MethodGet, "/analytics", nil).Body.String()
	if strings.Contains(body, "/read/"+slug) {
		t.Error("страница аналитики выдаёт адреса отдельных статей")
	}
	if strings.Contains(body, "statspriv@example.com") {
		t.Error("на странице оказался адрес электронной почты")
	}
}

// Страница обязана открываться на всех трёх языках: она часть издания, а не
// служебная выгрузка.
func TestTheAudiencePageSpeaksAllThreeLanguages(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	for _, lang := range []string{LangKZ, LangRU, LangEN} {
		w := app.do(http.MethodGet, "/analytics?lang="+lang, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("язык %s: %d", lang, w.Code)
		}
		if want := T(lang, "stats.lead"); !strings.Contains(w.Body.String(), want) {
			t.Errorf("язык %s: нет вводной строки на этом языке", lang)
		}
	}
}
