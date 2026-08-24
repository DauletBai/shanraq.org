package articles

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

// Before this block an article page linked to no other article at all. The site's
// whole link structure ran outward from the front page and stopped dead: a reader
// who finished had nowhere to go, and every article was a leaf leading nowhere.
func TestArticleOffersSomewhereToGoNext(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	authorID := app.createUser("relauthor@example.com", "Parol123!")
	id, slug := app.seedArticle(authorID, "published")
	app.exec(`UPDATE articles SET published_at = NOW() - interval '1 hour' WHERE id = $1`, id)

	// Newer than the published one — otherwise other articles get suggested, and the
	// shared test database always holds plenty of those.
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
	// Judge by the block alone: further down the page there is a comment form that
	// links to the current article, so the whole page cannot be checked.
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

// Machine columns are open to whoever chose them and closed to search engines.
// Slipping them under a human article means forcing a machine's opinion on a reader
// who did not ask for it, and spending a robot's visit on a page that cannot be
// indexed anyway.
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

// Nearest first: the same subsection, then the same section, then simply the
// newest. Otherwise "read also" turns into a random draw.
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

// The date in a sitemap has to be a real one. Section feeds have one — they change
// when an article lands in the section; service pages do not, and an invented date
// there would teach Google that our dates mean nothing.
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
	// An article closed to indexing does not move the date: it is not in the sitemap.
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
