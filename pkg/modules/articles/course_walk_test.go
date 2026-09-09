package articles

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// The reader's own path, end to end: the contents, a lesson, its exercise box,
// the tidy endpoint that box calls, and the link onwards.
//
// Every piece of this has been tested on its own and the path between them had
// not been: the contents numbered the announcement as lesson one for months
// while the strip beside a lesson counted correctly, and nothing failed,
// because nothing walked from one to the other.
func TestReaderWalksTheCourse(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	author := app.createUser("walk@t.test", "Parol12345")
	app.exec(`UPDATE auth_users SET email_verified_at = now() WHERE lower(email) = 'walk@t.test'`)
	cookie := app.login("walk@t.test", "Parol12345")

	series := uuid.New()
	app.exec(`INSERT INTO article_series (id, slug, cover_url, status, code_lang)
		VALUES ($1, $2, '', 'published', 'python')`, series, "walk-"+series.String()[:8])
	app.exec(`INSERT INTO article_series_i18n (series_id, lang, title, summary)
		VALUES ($1, 'ru', 'Курс для прохода', 'Проверка пути читателя')`, series)

	// The announcement stands in front of the lessons, at a position below ten.
	intro, _ := app.seedArticle(author, "published")
	app.exec(`INSERT INTO article_series_items (series_id, article_id, position) VALUES ($1, $2, 1)`, series, intro)

	// Two lessons, the first with an exercise the box can offer to check.
	first, firstSlug := app.seedArticle(author, "published")
	app.exec(`UPDATE article_translations SET body_md = $2 WHERE article_id = $1 AND lang = 'ru'`,
		first, "## Зачем это нужно\n\nТекст урока.\n\n## Задание\n\n**Обязательное.** Напишите функцию.\n")
	app.exec(`INSERT INTO article_series_items (series_id, article_id, position) VALUES ($1, $2, 10)`, series, first)

	second, secondSlug := app.seedArticle(author, "published")
	app.exec(`INSERT INTO article_series_items (series_id, article_id, position) VALUES ($1, $2, 20)`, series, second)

	slug := seriesSlug(t, NewSeriesStore(app.pool), series)

	t.Run("оглавление считает уроки, а не страницы", func(t *testing.T) {
		rec := app.do("GET", "/course/"+slug, nil)
		if rec.Code != 200 {
			t.Fatalf("карта курса ответила %d", rec.Code)
		}
		body := rec.Body.String()
		if !strings.Contains(body, "/read/"+firstSlug) {
			t.Error("на карте нет ссылки на первый урок")
		}
		// Two lessons and an announcement: the numbers shown are 1 and 2, and
		// the announcement carries no number at all.
		if strings.Count(body, `class="lesson__no" aria-hidden="true">3<`) > 0 {
			t.Error("анонс посчитан уроком: на карте появился номер 3")
		}
		for _, want := range []string{`aria-hidden="true">1<`, `aria-hidden="true">2<`} {
			if !strings.Contains(body, want) {
				t.Errorf("на карте нет номера %q", want)
			}
		}
	})

	t.Run("урок предлагает проверку и ведёт дальше", func(t *testing.T) {
		rec := app.do("GET", "/read/"+firstSlug, nil, withCookie(cookie))
		if rec.Code != 200 {
			t.Fatalf("урок ответил %d", rec.Code)
		}
		body := rec.Body.String()
		if !strings.Contains(body, `data-check="`+firstSlug+`"`) {
			t.Error("на уроке с заданием нет блока проверки")
		}
		if !strings.Contains(body, "/read/"+secondSlug) {
			t.Error("с урока нет пути к следующему")
		}
	})

	t.Run("кнопка «причесать» отвечает на решение", func(t *testing.T) {
		f := url.Values{}
		f.Set("solution", "def average(series):   \n    return sum(series) / len(series)")
		rec := app.do(http.MethodPost, "/read/"+firstSlug+"/format", f, withCookie(cookie))
		if rec.Code != 200 {
			t.Fatalf("проверка формата ответила %d: %s", rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), "def average") {
			t.Errorf("решение не вернулось: %.120s", rec.Body.String())
		}
	})
}
