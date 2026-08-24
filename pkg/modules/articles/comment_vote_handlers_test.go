package articles

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// Comment voting end to end, over real HTTP requests: what used to stand here was
// a model that read every comment and decided whether anyone would see it. Readers
// decide now, and what has to be tested is exactly what they can do by hand.
func TestCommentVotingEndToEnd(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	authorID := app.createUser("comauthor@example.com", "Parol123!")
	app.createUser("comreader@example.com", "Parol123!")

	articleID, slug := app.seedArticle(authorID, "published")

	// The article's author writes the comment — so they will not be able to vote on it.
	authorCookie := app.login("comauthor@example.com", "Parol123!")
	if w := app.do(http.MethodPost, "/read/"+slug+"/comment",
		url.Values{"body": {"Мой комментарий"}}, withCookie(authorCookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("создание комментария: получили %d (%s)", w.Code, w.Body.String())
	}

	var commentID uuid.UUID
	if err := app.pool.QueryRow(context.Background(),
		`SELECT id FROM comments WHERE article_id = $1`, articleID).Scan(&commentID); err != nil {
		t.Fatalf("комментарий не сохранился: %v", err)
	}

	score := func() int {
		var s int
		if err := app.pool.QueryRow(context.Background(),
			`SELECT score FROM comments WHERE id = $1`, commentID).Scan(&s); err != nil {
			t.Fatalf("счёт: %v", err)
		}
		return s
	}
	vote := func(c *http.Cookie, value string) int {
		t.Helper()
		w := app.do(http.MethodPost, "/read/"+slug+"/comment/"+commentID.String()+"/vote",
			url.Values{"value": {value}}, withCookie(c))
		if w.Code != http.StatusSeeOther {
			t.Fatalf("голос %q: получили %d (%s)", value, w.Code, w.Body.String())
		}
		return score()
	}

	// A new comment starts at zero.
	if got := score(); got != 0 {
		t.Errorf("новый комментарий имеет счёт %d, а должен ноль", got)
	}

	// You do not rate your own comment. The request does not fail — it simply changes nothing.
	if got := vote(authorCookie, "1"); got != 0 {
		t.Errorf("автор проголосовал за свой комментарий: счёт %d", got)
	}

	readerCookie := app.login("comreader@example.com", "Parol123!")

	// An upvote from a reader with no karma weighs one.
	if got := vote(readerCookie, "1"); got != 1 {
		t.Errorf("после плюса счёт %d, ожидался 1", got)
	}
	// Clicking the same arrow again withdraws the vote.
	if got := vote(readerCookie, "1"); got != 0 {
		t.Errorf("повторный плюс не снял голос: счёт %d", got)
	}
	// A downvote.
	if got := vote(readerCookie, "-1"); got != -1 {
		t.Errorf("после минуса счёт %d, ожидался -1", got)
	}
	// Changing your mind replaces the vote rather than adding a second one.
	if got := vote(readerCookie, "1"); got != 1 {
		t.Errorf("смена голоса дала %d, ожидался 1", got)
	}
	var votes int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM comment_votes WHERE comment_id = $1`, commentID).Scan(&votes); err != nil {
		t.Fatalf("подсчёт голосов: %v", err)
	}
	if votes != 1 {
		t.Errorf("у одного читателя %d голосов вместо одного", votes)
	}

	// A guest cannot vote — they are sent to the login page.
	w := app.do(http.MethodPost, "/read/"+slug+"/comment/"+commentID.String()+"/vote",
		url.Values{"value": {"1"}})
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/studio/login" {
		t.Errorf("гостя не отправили на вход: %d -> %s", w.Code, w.Header().Get("Location"))
	}
}

// A downvoted comment folds away but does not disappear: the row with its score
// stays on the page, and anyone can unfold it.
func TestBuriedCommentIsFoldedNotRemoved(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	authorID := app.createUser("foldauthor@example.com", "Parol123!")
	articleID, slug := app.seedArticle(authorID, "published")

	cookie := app.login("foldauthor@example.com", "Parol123!")
	if w := app.do(http.MethodPost, "/read/"+slug+"/comment",
		url.Values{"body": {"Непопулярное мнение"}}, withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("создание комментария: %d", w.Code)
	}
	app.exec(`UPDATE comments SET score = $2 WHERE article_id = $1`, articleID, commentCollapseScore-1)

	w := app.do(http.MethodGet, "/read/"+slug, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("страница статьи: %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "comment--folded") {
		t.Error("заминусованный комментарий не свёрнут")
	}
	// Folded means hidden behind a row, not removed from the page.
	if !strings.Contains(body, "Непопулярное мнение") {
		t.Error("текст свёрнутого комментария пропал со страницы")
	}
}
