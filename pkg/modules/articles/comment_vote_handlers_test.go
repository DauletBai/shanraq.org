package articles

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// Голосование за комментарии целиком, через настоящие HTTP-запросы: раньше на
// этом месте стояла модель, которая читала каждый комментарий и решала, увидит
// его кто-нибудь или нет. Теперь решают читатели, и проверять надо именно то,
// что они могут сделать руками.
func TestCommentVotingEndToEnd(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	authorID := app.createUser("comauthor@example.com", "Parol123!")
	app.createUser("comreader@example.com", "Parol123!")

	articleID, slug := app.seedArticle(authorID, "published")

	// Комментарий пишет автор статьи — голосовать за него он не сможет.
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

	// Новый комментарий начинается с нуля.
	if got := score(); got != 0 {
		t.Errorf("новый комментарий имеет счёт %d, а должен ноль", got)
	}

	// Свой комментарий не оценивают. Запрос не падает — он просто ничего не меняет.
	if got := vote(authorCookie, "1"); got != 0 {
		t.Errorf("автор проголосовал за свой комментарий: счёт %d", got)
	}

	readerCookie := app.login("comreader@example.com", "Parol123!")

	// Плюс от читателя без кармы весит единицу.
	if got := vote(readerCookie, "1"); got != 1 {
		t.Errorf("после плюса счёт %d, ожидался 1", got)
	}
	// Повторный клик по той же стрелке снимает голос.
	if got := vote(readerCookie, "1"); got != 0 {
		t.Errorf("повторный плюс не снял голос: счёт %d", got)
	}
	// Минус.
	if got := vote(readerCookie, "-1"); got != -1 {
		t.Errorf("после минуса счёт %d, ожидался -1", got)
	}
	// Смена мнения заменяет голос, а не добавляет второй.
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

	// Гость голосовать не может — его отправляют на вход.
	w := app.do(http.MethodPost, "/read/"+slug+"/comment/"+commentID.String()+"/vote",
		url.Values{"value": {"1"}})
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/studio/login" {
		t.Errorf("гостя не отправили на вход: %d -> %s", w.Code, w.Header().Get("Location"))
	}
}

// Заминусованный комментарий сворачивается, но не исчезает: строка с оценкой
// остаётся на странице, и раскрыть её может любой.
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
	// Свёрнут — значит спрятан за строкой, а не удалён со страницы.
	if !strings.Contains(body, "Непопулярное мнение") {
		t.Error("текст свёрнутого комментария пропал со страницы")
	}
}
