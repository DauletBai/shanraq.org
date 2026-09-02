package articles

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// A course is only worth building if the chain through it never breaks. These
// tests cover the two ways it can: lessons coming out in the wrong order, and
// the chain walking a reader into a draft.

// seedSeries makes a published course and returns its id.
func (a *testApp) seedSeries(slug string) uuid.UUID {
	a.t.Helper()
	id, err := NewSeriesStore(a.pool).Save(context.Background(), nil, slug, "", SeriesPublished,
		map[string]string{LangRU: "Курс Go", LangKZ: "Go курсы"},
		map[string]string{LangRU: "С нуля до своего блога"})
	if err != nil {
		a.t.Fatalf("seedSeries: %v", err)
	}
	return id
}

// Lessons must come out in the order the editor set, not the order they were
// added -- a course whose third lesson was written first would otherwise open
// with it.
func TestSeriesOrdersByPosition(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	author := app.createUser("series-order@example.com", "Passw0rd!x")
	st := NewSeriesStore(app.pool)
	ctx := context.Background()

	sid := app.seedSeries("go-order-" + uuid.NewString()[:6])
	third, _ := app.seedArticle(author, "published")
	first, _ := app.seedArticle(author, "published")
	second, _ := app.seedArticle(author, "published")

	// Added out of order, on purpose.
	mustAttach(t, st, sid, third, 30)
	mustAttach(t, st, sid, first, 10)
	mustAttach(t, st, sid, second, 20)

	s, err := st.BySlug(ctx, seriesSlug(t, st, sid), LangRU)
	if err != nil {
		t.Fatalf("BySlug: %v", err)
	}
	want := []uuid.UUID{first, second, third}
	if len(s.Items) != len(want) {
		t.Fatalf("got %d lessons, want %d", len(s.Items), len(want))
	}
	for i, id := range want {
		if s.Items[i].ArticleID != id {
			t.Errorf("lesson %d: got %s, want %s", i+1, s.Items[i].ArticleID, id)
		}
	}
	if s.Lessons() != 3 {
		t.Errorf("Lessons() = %d, want 3", s.Lessons())
	}
}

// The chain must step over an unpublished lesson rather than stopping at it or
// linking to it: a course under construction has gaps, and a reader following
// "next" through one must still land on something readable.
func TestSeriesNavSkipsUnpublished(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	author := app.createUser("series-skip@example.com", "Passw0rd!x")
	st := NewSeriesStore(app.pool)
	ctx := context.Background()

	sid := app.seedSeries("go-skip-" + uuid.NewString()[:6])
	one, _ := app.seedArticle(author, "published")
	hole, _ := app.seedArticle(author, "draft")
	three, _ := app.seedArticle(author, "published")

	mustAttach(t, st, sid, one, 10)
	mustAttach(t, st, sid, hole, 20)
	mustAttach(t, st, sid, three, 30)

	places, err := st.ForArticle(ctx, one, LangRU)
	if err != nil {
		t.Fatalf("ForArticle: %v", err)
	}
	if len(places) != 1 {
		t.Fatalf("got %d courses, want 1", len(places))
	}
	p := places[0]
	if p.Total != 2 {
		t.Errorf("Total = %d, want 2 (the draft must not be counted)", p.Total)
	}
	if p.Number != 1 {
		t.Errorf("Number = %d, want 1", p.Number)
	}
	if p.Prev != nil {
		t.Errorf("first lesson has a previous: %s", p.Prev.Slug)
	}
	if p.Next == nil {
		t.Fatal("first lesson has no next; the draft broke the chain")
	}
	if p.Next.ArticleID != three {
		t.Errorf("next = %s, want the published lesson %s", p.Next.ArticleID, three)
	}

	// And the draft itself reports no place, so its page shows no strip.
	holePlaces, err := st.ForArticle(ctx, hole, LangRU)
	if err != nil {
		t.Fatalf("ForArticle(draft): %v", err)
	}
	if len(holePlaces) != 1 || holePlaces[0].Number != 0 {
		t.Errorf("draft reports a lesson number; it must not")
	}
}

// Deleting a course must not take its lessons down: they are ordinary articles
// that readers may already have linked to.
func TestSeriesDeleteKeepsArticles(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	author := app.createUser("series-del@example.com", "Passw0rd!x")
	st := NewSeriesStore(app.pool)
	ctx := context.Background()

	sid := app.seedSeries("go-del-" + uuid.NewString()[:6])
	lesson, slug := app.seedArticle(author, "published")
	mustAttach(t, st, sid, lesson, 10)

	if err := st.Delete(ctx, sid); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := NewStore(app.pool).GetPublishedBySlug(ctx, slug); err != nil {
		t.Fatalf("lesson vanished with its course: %v", err)
	}
}

func mustAttach(t *testing.T, st *SeriesStore, sid, aid uuid.UUID, pos int) {
	t.Helper()
	if err := st.Attach(context.Background(), sid, aid, pos); err != nil {
		t.Fatalf("Attach: %v", err)
	}
}

func seriesSlug(t *testing.T, st *SeriesStore, id uuid.UUID) string {
	t.Helper()
	s, err := st.ByID(context.Background(), id, LangRU)
	if err != nil {
		t.Fatalf("ByID: %v", err)
	}
	return s.Slug
}

// A draft course must not announce itself through a lesson that went live.
// The strip and the link on an article page come from ForArticle, which found
// every course a lesson belonged to and asked nothing about its status — so a
// course kept back deliberately was published by its own first lesson.
func TestDraftSeriesStaysOffPublishedLesson(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	author := app.createUser("series-draft@example.com", "Passw0rd!x")
	st := NewSeriesStore(app.pool)
	ctx := context.Background()

	sid, err := st.Save(ctx, nil, "hidden-"+uuid.NewString()[:6], "", SeriesDraft,
		map[string]string{LangRU: "Черновой курс"}, map[string]string{})
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	lesson, _ := app.seedArticle(author, "published")
	mustAttach(t, st, sid, lesson, 10)

	places, err := st.ForArticle(ctx, lesson, LangRU)
	if err != nil {
		t.Fatalf("ForArticle: %v", err)
	}
	if len(places) != 0 {
		t.Fatalf("черновой курс виден на опубликованном уроке: %d", len(places))
	}

	// And it appears the moment the course itself is published.
	if _, err := st.Save(ctx, &sid, seriesSlug(t, st, sid), "", SeriesPublished,
		map[string]string{LangRU: "Черновой курс"}, map[string]string{}); err != nil {
		t.Fatalf("publish: %v", err)
	}
	if places, err = st.ForArticle(ctx, lesson, LangRU); err != nil || len(places) != 1 {
		t.Fatalf("опубликованный курс не виден: %d, %v", len(places), err)
	}
}
