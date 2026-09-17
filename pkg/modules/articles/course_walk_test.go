package articles

import (
	"context"
	"fmt"
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

func TestRustSelfCheckDoesNotUseGoFormatter(t *testing.T) {
	if _, err := formatSolution("fn main() {}", CodeRust); err == nil {
		t.Fatal("Rust must not claim a Go formatter checked the answer")
	}
	if !strings.Contains(checkSystem(LangKZ, CodeRust), "Rust") {
		t.Fatal("Rust review must not identify the course as Go")
	}
	if got := string(highlightCode("fn main() {}", CodeRust)); !strings.Contains(got, "chroma") {
		t.Fatal("Rust syntax was not highlighted")
	}
}

func TestRustFirstBatchStaysOutOfFeed(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	email := "rust-" + uuid.NewString() + "@t.test"
	author := app.createUser(email, "Sup3r-Secret-Pass!")
	courseSlug := "rust-test-" + uuid.NewString()
	titles := map[string]string{"ru": "Rust RU", "kz": "Rust KZ", "en": "Rust EN"}
	series, err := NewSeriesStore(app.pool).Save(context.Background(), nil, courseSlug, "", SeriesPublished, CodeRust, titles, titles)
	if err != nil {
		t.Fatal(err)
	}
	defer app.exec(`DELETE FROM article_series WHERE id=$1`, series)
	var slugs []string
	for n := 0; n <= 5; n++ {
		id, slug := app.seedArticle(author, "published")
		slugs = append(slugs, slug)
		pos := n * 10
		if n == 0 {
			pos = 1
		}
		app.exec(`INSERT INTO article_series_items(series_id,article_id,position) VALUES($1,$2,$3)`, series, id, pos)
		for _, lang := range []string{"kz", "ru", "en"} {
			heading := map[string]string{"kz": "Тапсырма", "ru": "Задание", "en": "Exercise"}[lang]
			body := "## " + heading + "\n\nExplain your result and compare it with the reference answer.\n"
			app.exec(`INSERT INTO article_translations(article_id,lang,title,summary,body_md,source,status)
                VALUES($1,$2,$3,'Rust lesson',$4,'ai','ready') ON CONFLICT(article_id,lang)
                DO UPDATE SET title=EXCLUDED.title,body_md=EXCLUDED.body_md,status='ready'`, id, lang, fmt.Sprintf("Rust %d %s", n, lang), body)
		}
		app.exec(`UPDATE articles SET published_at='2099-02-01',score=100000 WHERE id=$1`, id)
	}
	for _, lang := range []string{"kz", "ru", "en"} {
		course := app.do(http.MethodGet, "/course/"+courseSlug+"?lang="+lang, nil)
		if course.Code != http.StatusOK {
			t.Fatalf("course: %d", course.Code)
		}
		body := course.Body.String()
		for _, slug := range slugs {
			if !strings.Contains(body, "/read/"+slug) {
				t.Errorf("%s: missing %s", lang, slug)
			}
		}
		if !strings.Contains(body, `aria-hidden="true">5<`) || strings.Contains(body, `aria-hidden="true">6<`) {
			t.Errorf("%s: preface incorrectly counted as lesson", lang)
		}
		lesson := app.do(http.MethodGet, "/read/"+slugs[1]+"?lang="+lang, nil)
		if lesson.Code != http.StatusOK || !strings.Contains(lesson.Body.String(), "/read/"+slugs[2]) {
			t.Errorf("%s: missing next-lesson navigation", lang)
		}
		if strings.Contains(lesson.Body.String(), `data-check="`+slugs[1]+`"`) {
			t.Errorf("%s: unsupported answer checker shown", lang)
		}
		for _, sort := range []string{"", "top"} {
			feed := app.do(http.MethodGet, "/?lang="+lang+"&sort="+sort, nil)
			if feed.Code != http.StatusOK {
				t.Fatalf("feed: %d", feed.Code)
			}
			for _, slug := range slugs {
				if feedSlugs(feed.Body.String())[slug] {
					t.Errorf("Rust lesson leaked to feed: %s", slug)
				}
			}
		}
	}
	cookie := app.login(email, "Sup3r-Secret-Pass!")
	form := url.Values{"solution": {"fn main() {}"}}
	response := app.do(http.MethodPost, "/read/"+slugs[1]+"/format", form, withCookie(cookie))
	if response.Code != http.StatusNotImplemented {
		t.Errorf("Rust formatter status: %d", response.Code)
	}
}
