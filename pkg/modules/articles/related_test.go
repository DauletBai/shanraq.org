package articles

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

// До этого блока страница статьи не ссылалась ни на одну другую статью. Вся
// ссылочная схема сайта шла с главной наружу и обрывалась: дочитавшему читателю
// некуда было идти, а каждая статья была листом, из которого не ведёт ничего.
func TestArticleOffersSomewhereToGoNext(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	authorID := app.createUser("relauthor@example.com", "Parol123!")
	id, slug := app.seedArticle(authorID, "published")
	app.exec(`UPDATE articles SET published_at = NOW() - interval '1 hour' WHERE id = $1`, id)

	// Свежее опубликованного — иначе предложены будут чужие статьи, которых в
	// общей тестовой базе всегда хватает.
	var others []string
	for i := 0; i < 3; i++ {
		oid, s := app.seedArticle(authorID, "published")
		app.exec(`UPDATE articles SET published_at = NOW() WHERE id = $1`, oid)
		others = append(others, s)
	}

	w := app.do(http.MethodGet, "/read/"+slug, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("страница статьи: %d", w.Code)
	}
	body := w.Body.String()

	i := strings.Index(body, `class="related"`)
	if i < 0 {
		t.Fatal("блока «читайте также» на странице нет")
	}
	// Судим только по самому блоку: ниже на странице есть форма комментариев,
	// которая ссылается на текущую статью, и по всей странице проверять нельзя.
	block := body[i:]
	if end := strings.Index(block, "</section>"); end > 0 {
		block = block[:end]
	}

	found := 0
	for _, s := range others {
		if strings.Contains(block, "/read/"+s) {
			found++
		}
	}
	if found == 0 {
		t.Errorf("блок есть, а свежих статей в нём нет: %s", block)
	}
	if strings.Contains(block, "/read/"+slug) {
		t.Error("статья предлагает прочитать саму себя")
	}
}

// Машинные колонки открыты тому, кто их выбрал, и закрыты поисковикам. Подсунуть
// их под человеческой статьёй — значит навязать читателю мнение машины, о котором
// он не просил, и потратить визит робота на страницу, которую всё равно нельзя
// проиндексировать.
func TestRelatedNeverOffersMachineColumns(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	authorID := app.createUser("relmix@example.com", "Parol123!")
	humanID, humanSlug := app.seedArticle(authorID, "published")
	machineID, machineSlug := app.seedArticle(authorID, "published")
	app.exec(`UPDATE articles SET indexable = false WHERE id = $1`, machineID)

	rel, err := NewStore(app.pool).RelatedPublished(
		context.Background(), humanID, "economy", "", 10, nil)
	if err != nil {
		t.Fatalf("RelatedPublished: %v", err)
	}
	for _, a := range rel {
		if a.Slug == machineSlug {
			t.Error("в «читайте также» попала статья, закрытая от индексации")
		}
		if a.ID == humanID {
			t.Error("статья предложена сама себе")
		}
		if !a.Indexable {
			t.Errorf("предложена неиндексируемая статья %s", a.Slug)
		}
	}
	_ = humanSlug
}

// Ближайшее — первым: сначала та же подрубрика, потом та же рубрика, потом
// просто свежее. Иначе «читайте также» превращается в случайную выдачу.
func TestRelatedPutsTheNearestFirst(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	authorID := app.createUser("relnear@example.com", "Parol123!")
	baseID, _ := app.seedArticle(authorID, "published")
	app.exec(`UPDATE articles SET category='economy', subcategory='banks' WHERE id=$1`, baseID)

	farID, farSlug := app.seedArticle(authorID, "published")
	app.exec(`UPDATE articles SET category='society', subcategory='' WHERE id=$1`, farID)

	nearID, nearSlug := app.seedArticle(authorID, "published")
	app.exec(`UPDATE articles SET category='economy', subcategory='banks' WHERE id=$1`, nearID)

	rel, err := NewStore(app.pool).RelatedPublished(
		context.Background(), baseID, "economy", "banks", 4, nil)
	if err != nil {
		t.Fatalf("RelatedPublished: %v", err)
	}
	if len(rel) < 2 {
		t.Fatalf("предложено %d статей, ожидалось хотя бы две", len(rel))
	}
	if rel[0].Slug != nearSlug {
		t.Errorf("первой предложена %q, а должна статья из той же подрубрики (%q)", rel[0].Slug, nearSlug)
	}
	_ = farSlug
}

// Дата в sitemap должна быть настоящей. Ленты рубрик её имеют — они меняются,
// когда в раздел выходит статья; служебные страницы не имеют, и выдуманная
// дата там научила бы Google, что наши даты ничего не значат.
func TestCategoryFreshnessComesFromRealArticles(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	authorID := app.createUser("relfresh@example.com", "Parol123!")
	id, _ := app.seedArticle(authorID, "published")
	app.exec(`UPDATE articles SET category='economy', published_at = NOW() WHERE id=$1`, id)

	fresh, err := NewStore(app.pool).CategoryFreshness(context.Background())
	if err != nil {
		t.Fatalf("CategoryFreshness: %v", err)
	}
	if fresh["economy"].IsZero() {
		t.Error("у рубрики с опубликованной статьёй нет даты")
	}
	// Закрытая от индексации статья дату не двигает: её нет в sitemap.
	hidden, _ := app.seedArticle(authorID, "published")
	app.exec(`UPDATE articles SET category='sport', indexable=false, published_at=NOW() WHERE id=$1`, hidden)
	fresh2, err := NewStore(app.pool).CategoryFreshness(context.Background())
	if err != nil {
		t.Fatalf("CategoryFreshness: %v", err)
	}
	if !fresh2["sport"].IsZero() {
		t.Error("неиндексируемая статья задала дату рубрике, которой нет в sitemap")
	}
}
