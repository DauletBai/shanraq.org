package articles

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// kacharID возвращает узел справочника для посёлка Качар — тот самый пример,
// ради которого всё затевалось: председатель ЖКХ публикует для своего посёлка.
func kacharID(t *testing.T, app *testApp) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := app.pool.QueryRow(context.Background(),
		`SELECT id FROM geo_nodes WHERE name_ru = 'Качар' LIMIT 1`).Scan(&id); err != nil {
		t.Skipf("в тестовой базе нет справочника мест: %v", err)
	}
	return id
}

// Место — одно поле, а не четыре. Республику и область из него выводит дерево,
// и они не могут разойтись с городом, как разошлись бы четыре отдельные колонки.
func TestUserPlaceIsOneNodeAndItsAncestryGivesTheRest(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	userID := app.createUser("place1@example.com", "Parol123!")
	node := kacharID(t, app)
	geo := NewGeoStore(app.pool)

	if err := geo.SetUserPlace(context.Background(), userID, &node); err != nil {
		t.Fatalf("SetUserPlace: %v", err)
	}
	got, err := geo.UserPlace(context.Background(), userID)
	if err != nil || got == nil || *got != node {
		t.Fatalf("UserPlace вернул %v, %v", got, err)
	}

	chain, err := geo.Ancestry(context.Background(), node, LangRU)
	if err != nil {
		t.Fatalf("Ancestry: %v", err)
	}
	if len(chain) < 3 {
		t.Fatalf("цепочка предков короткая: %+v", chain)
	}
	if chain[0].Kind != "country" {
		t.Errorf("цепочка начинается не со страны: %+v", chain[0])
	}
	if chain[len(chain)-1].Name != "Качар" {
		t.Errorf("цепочка кончается не выбранным местом: %+v", chain[len(chain)-1])
	}

	// Подпись читателю: место и область, чтобы тёзку нельзя было спутать.
	label, err := geo.PlaceLabel(context.Background(), node, LangRU)
	if err != nil {
		t.Fatalf("PlaceLabel: %v", err)
	}
	if !strings.Contains(label, "Качар") || !strings.Contains(label, "область") {
		t.Errorf("подпись места не называет область: %q", label)
	}
}

// Место выбрано в спешке или человек переехал — он должен поправить сам, без
// письма в поддержку. И должен уметь стереть его совсем.
func TestReaderCanChangeAndClearTheirPlace(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	app.createUser("place2@example.com", "Parol123!")
	cookie := app.login("place2@example.com", "Parol123!")
	node := kacharID(t, app)

	w := app.do(http.MethodPost, "/studio/place",
		url.Values{"geo_node_id": {node.String()}}, withCookie(cookie))
	if w.Code != http.StatusSeeOther {
		t.Fatalf("сохранение места: %d (%s)", w.Code, w.Body.String())
	}

	page := app.do(http.MethodGet, "/studio/profile", nil, withCookie(cookie))
	if !strings.Contains(page.Body.String(), "Качар") {
		t.Error("профиль не показывает выбранное место")
	}

	// Пустое значение стирает выбор, а не ломает форму.
	if w := app.do(http.MethodPost, "/studio/place",
		url.Values{"geo_node_id": {""}}, withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("очистка места: %d", w.Code)
	}
	var left int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM user_places WHERE user_id = (SELECT id FROM auth_users WHERE email = $1)`,
		"place2@example.com").Scan(&left); err != nil {
		t.Fatalf("проверка: %v", err)
	}
	if left != 0 {
		t.Errorf("после очистки осталось %d записей", left)
	}
}

// Мусор в поле не должен ни падать, ни молча сохраняться.
func TestBadPlaceIsRefusedWithoutBreakingTheProfile(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	app.createUser("place3@example.com", "Parol123!")
	cookie := app.login("place3@example.com", "Parol123!")

	w := app.do(http.MethodPost, "/studio/place",
		url.Values{"geo_node_id": {"не-идентификатор"}}, withCookie(cookie))
	if w.Code != http.StatusSeeOther {
		t.Fatalf("получили %d, ожидался переход обратно в профиль", w.Code)
	}
	if loc := w.Header().Get("Location"); !strings.Contains(loc, "place_bad") {
		t.Errorf("читателю не сказали, что не вышло: %s", loc)
	}
}

// Поле необязательное, и это принципиально: регистрация — узкое место, и
// требование выбрать посёлок до создания аккаунта стоило бы дороже, чем даёт.
// Кто не выбрал — видит всё, ровно как до появления мест.
func TestRegistrationWorksWithoutAPlace(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	app.emails = append(app.emails, "placefree@example.com")
	w := app.do(http.MethodPost, "/studio/register", url.Values{
		"first_name": {"Тест"}, "last_name": {"Пользователь"},
		"email": {"placefree@example.com"}, "password": {"Parol123!"},
		"consent": {"on"},
	})
	if w.Code != http.StatusSeeOther {
		t.Fatalf("регистрация без места: %d (%s)", w.Code, w.Body.String())
	}
	var n int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM auth_users WHERE email = $1`, "placefree@example.com").Scan(&n); err != nil {
		t.Fatalf("проверка: %v", err)
	}
	if n != 1 {
		t.Errorf("аккаунт не создан: %d", n)
	}
}

// Каждое поле формы обязано объяснять, зачем оно. Место — тем более: человек
// сообщает, где живёт, и вправе знать, что с этим будет.
func TestPlaceFieldExplainsItself(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	w := app.do(http.MethodGet, "/studio/register", nil)
	body := w.Body.String()
	if !strings.Contains(body, `data-geo`) {
		t.Fatal("в форме регистрации нет выбора места")
	}
	for _, must := range []string{"Где вы живёте", "Необязательно"} {
		if !strings.Contains(body, must) {
			t.Errorf("нет пояснения к полю места: %q", must)
		}
	}
	if !strings.Contains(body, "fhelp") {
		t.Error("у поля места нет подсказки «?»")
	}
}
