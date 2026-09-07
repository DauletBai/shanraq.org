# Теги: третья таблица, которой не видно на странице

_Лид (summary):_ **Тридцать третий урок курса по Go. У статьи много тегов, у тега много статей — и хранится это отдельной таблицей, а не списком в колонке. Составной ключ, on conflict ради id, каскад при удалении. И обещанный способ спросить у ошибки её код вместо того, чтобы читать её текст.**

## Зачем это нужно

Теги выглядят безобидно: пара слов под заголовком. Соблазн сделать их колонкой — `tags text` со списком через запятую — велик, и он ошибочен.

Стоит спросить у такой таблицы «дай статьи с тегом go», и начинается:

```sql
select slug, tags from a where tags like '%go%';
```

```
┌────────┬────────────────┐
│  slug  │      tags      │
├────────┼────────────────┤
│ dala   │ go,степь       │
│ golang │ golang,новости │
└────────┴────────────────┘
```

`golang` попал в ответ, потому что внутри него есть `go`. Дальше выяснится, что тег нельзя переименовать одной командой, нельзя посчитать, сколько раз он встречается, и нельзя не дать положить его дважды.

Причина простая: связь здесь **многие-ко-многим**. У статьи несколько тегов, у тега несколько статей. Такую связь колонкой не выразить — для неё заводят третью таблицу.

## Сразу целиком

Новая папка, `go mod init sabaq30`:

```go
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	sqlite "modernc.org/sqlite"
)

// uniqueViolation — код, которым SQLite отвечает на нарушение unique.
// Список кодов — на sqlite.org/rescode.html.
const uniqueViolation = 2067

// ErrTaken — адрес уже занят. Обработчику важно именно это, а не код базы.
var ErrTaken = errors.New("адрес занят")

func main() {
	os.Remove("blog.db")
	db, err := sql.Open("sqlite", "blog.db?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	setup(db)

	fmt.Println("добавили dala: ", add(db, "dala", "О степи", "go", "степь"))
	fmt.Println("добавили salem:", add(db, "salem", "Привет", "go", "блог"))
	fmt.Println("тот же адрес:  ", add(db, "dala", "Другая", "go"))

	fmt.Println("теги dala:      ", tagsOf(db, "dala"))
	fmt.Println("статьи с go:    ", byTag(db, "go"))
	fmt.Println("статьи со степь:", byTag(db, "степь"))

	db.Exec(`delete from articles where slug = 'dala'`)
	fmt.Println("после удаления dala:")
	fmt.Println("  статьи с go:  ", byTag(db, "go"))
	fmt.Println("  связей всего: ", count(db, `select count(*) from article_tags`))
	fmt.Println("  тегов всего:  ", count(db, `select count(*) from tags`))
}

// add кладёт статью и её теги. Две таблицы, одно действие — значит,
// транзакция: статья без тегов или тег без статьи никому не нужны.
func add(db *sql.DB, slug, title string, tags ...string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`insert into articles (slug, title) values (?, ?)`, slug, title)
	if err != nil {
		return wrap(err, slug)
	}
	articleID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	for _, name := range tags {
		var tagID int64
		// on conflict ... do update — единственный способ получить id и для
		// нового тега, и для уже существующего: do nothing вернул бы пусто.
		err := tx.QueryRow(`insert into tags (name) values (?)
			on conflict (name) do update set name = excluded.name
			returning id`, name).Scan(&tagID)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(
			`insert into article_tags (article_id, tag_id) values (?, ?)`,
			articleID, tagID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// wrap спрашивает у ошибки её код, а не читает её текст.
func wrap(err error, slug string) error {
	var serr *sqlite.Error
	if errors.As(err, &serr) && serr.Code() == uniqueViolation {
		return fmt.Errorf("%q: %w", slug, ErrTaken)
	}
	return err
}

func tagsOf(db *sql.DB, slug string) string {
	var out string
	err := db.QueryRow(`select group_concat(t.name, ', ')
		from tags t
		join article_tags at on at.tag_id = t.id
		join articles a on a.id = at.article_id
		where a.slug = ?`, slug).Scan(&out)
	if err != nil {
		return "ошибка: " + err.Error()
	}
	return out
}

func byTag(db *sql.DB, tag string) string {
	rows, err := db.Query(`select a.slug
		from articles a
		join article_tags at on at.article_id = a.id
		join tags t on t.id = at.tag_id
		where t.name = ?
		order by a.id`, tag)
	if err != nil {
		return "ошибка: " + err.Error()
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return "ошибка: " + err.Error()
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return "ошибка: " + err.Error()
	}
	return fmt.Sprint(out)
}

func count(db *sql.DB, q string) int {
	var n int
	if err := db.QueryRow(q).Scan(&n); err != nil {
		log.Fatal(err)
	}
	return n
}

func setup(db *sql.DB) {
	for _, q := range []string{
		`create table articles (
			id    integer primary key,
			slug  text not null unique,
			title text not null
		) strict`,
		`create table tags (
			id   integer primary key,
			name text not null unique
		) strict`,
		`create table article_tags (
			article_id integer not null references articles (id) on delete cascade,
			tag_id     integer not null references tags (id),
			primary key (article_id, tag_id)
		) strict`,
	} {
		if _, err := db.Exec(q); err != nil {
			log.Fatal(err)
		}
	}
}
```

Вывод:

```
добавили dala:  <nil>
добавили salem: <nil>
тот же адрес:   "dala": адрес занят
теги dala:       go, степь
статьи с go:     [dala salem]
статьи со степь: [dala]
после удаления dala:
  статьи с go:   [salem]
  связей всего:  2
  тегов всего:   3
```

## Разбор

### Третья таблица и составной ключ

`article_tags` не хранит ничего своего. В ней два столбца, и оба — ссылки: на статью и на тег. Одна строка означает «у этой статьи есть этот тег».

Такую таблицу называют связующей, и первичный ключ у неё **составной** — из двух колонок сразу:

```sql
primary key (article_id, tag_id)
```

Это не мелочь оформления. Составной ключ означает, что пара «статья + тег» может встретиться только один раз, и база не даст положить тег дважды:

```
Error: stepping, UNIQUE constraint failed: article_tags.article_id, article_tags.tag_id (19)
```

Проверять это в коде не нужно — и, как мы уже знаем из урока про инъекции, нельзя: между проверкой и вставкой успеет вклиниться другой запрос.

> **Образ.** Список приглашённых, где каждая строка — «этот человек на этой свадьбе». Ни человек, ни свадьба в этой строке не описаны: строка описывает только их встречу.

### `on conflict ... do update` ради `id`

Тег может быть новым, а может уже существовать — а `id` нужен в обоих случаях. Отсюда странная на вид строка:

```sql
insert into tags (name) values (?)
on conflict (name) do update set name = excluded.name
returning id
```

Читается так: попробуй вставить; если такое имя уже есть, «обнови» его тем же самым значением; в любом случае верни `id`.

Обновление здесь бессмысленное нарочно — оно нужно только затем, чтобы строка считалась затронутой и `returning` вернул её `id`. Написать `do nothing` не выйдет: конфликтующая строка тогда не считается затронутой, `returning` не вернёт ничего, и `Scan` получит `sql.ErrNoRows`. Ловушка известная, и обходят её именно так.

### Транзакция здесь не украшение

В прошлом уроке транзакции было некуда приложить. Теперь есть: `add` пишет в две таблицы, и половина работы хуже, чем ничего. Статья без тегов — терпимо, а вот тег, созданный для статьи, которая не вставилась, останется мусором навсегда.

Поэтому `Begin`, `defer tx.Rollback()` и один `Commit` в конце — ровно то, что разбирали, и впервые по-настоящему нужное.

### `errors.As`: спросить у ошибки код, а не читать текст

В уроке про CRUD я обещал: способ спросить у ошибки её код появится позже. Вот он.

```go
var serr *sqlite.Error
if errors.As(err, &serr) && serr.Code() == uniqueViolation {
	return fmt.Errorf("%q: %w", slug, ErrTaken)
}
```

`errors.Is` отвечает на вопрос «это та самая ошибка?». `errors.As` — на другой: «а не лежит ли внутри ошибка вот такого типа? если да, дай её мне». Она разворачивает цепочку `%w` так же, как `Is`, но кладёт найденное в переменную, и у этой переменной уже есть методы — здесь `Code()`.

Сравнивать текст ошибки нельзя: он меняется от версии к версии и от языка к языку. Код `2067` не меняется никогда — это `SQLITE_CONSTRAINT_UNIQUE`, и он записан в документации базы.

Заметьте, что делает `wrap`: он **переводит** ошибку базы в ошибку блога. Наружу уходит `ErrTaken`, и обработчику не нужно знать ни про SQLite, ни про числа. Захотите Postgres — поменяется одна функция.

И ещё: драйвер теперь подключён по имени, а не через `_`. Подчёркивание в уроке про SQLite означало «нужен только ради регистрации»; сегодня нам понадобился его тип, и имя вернулось.

### `on delete cascade`: связи уходят, теги остаются

Что происходит с тегами, когда статью удаляют, видно в выводе. Связей было четыре, осталось две. Их убрала база — сама, потому что в схеме написано:

```sql
article_id integer not null references articles (id) on delete cascade
```

А тегов осталось три, хотя «степь» больше никому не нужен. И это правильно: `cascade` стоит только на ссылке к статье. Удалять тег вместе с последней статьёй — решение, которое принимает блог, а не база.

Обратная сторона: тег, на который ссылаются, база удалить не даст:

```
Error: stepping, FOREIGN KEY constraint failed (19)
```

Сначала уберите связи, потом тег. Работает это, разумеется, только при включённой прагме `foreign_keys` — той самой, что стоит в строке подключения с урока про SQLite.

### Собрать теги обратно: `group_concat`

Показать теги статьи строкой — работа базы, а не Go:

```sql
select group_concat(t.name, ', ') from tags t join article_tags at ...
```

`group_concat` склеивает значения всех подошедших строк через разделитель. Получается `go, степь` — то, что и нужно вывести под заголовком.

Осторожно с двумя вещами. Порядок внутри `group_concat` не гарантирован — если он важен, задавайте его явно. И если тегов нет ни одного, функция вернёт `NULL`, а `Scan` в обычную строку такого не примет; в блоге это лечится `coalesce` или чтением в `sql.NullString`.

## Карта урока

![Карта урока: две таблицы и связь между ними третьей](/static/course/go/map-tags-ru.svg)

## Скажите своими словами

Не подглядывая, ответьте вслух или на бумаге. Ответы — в конце урока.

1. Почему теги не хранят списком в колонке статьи?
2. Зачем в `on conflict` написано бессмысленное обновление вместо `do nothing`?
3. Чем `errors.As` отличается от `errors.Is`?

## Задание

**Обязательное.** Добавьте теги в блог. Схему меняет миграция: две новые таблицы и связующая. Форма принимает теги строкой через запятую, `Store.Add` кладёт статью и теги одной транзакцией, страница статьи показывает теги, а `/tag/{name}` — список статей с этим тегом.

Всё это уже сделано в [step-18](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-18) — сверьтесь после того, как напишете сами.

**По желанию.**

- Приведите теги к одному виду перед вставкой: обрежьте пробелы и опустите регистр, чтобы `Go` и `go` не стали двумя тегами.
- Покажите на главной облако тегов: имя и число статей. Понадобится `group by` и `count`.
- Уберите `on delete cascade` и повторите удаление статьи. Прочитайте ошибку и объясните её вслух.

## Куда это встанет в блоге

Третья таблица — первая в блоге, у которой нет своей страницы. Она нужна не читателю, а связям, и это нормально: в настоящих базах таких таблиц много.

Долги. Теги пока не ищутся по кусочку слова — этим займётся следующий урок про поиск. Облако тегов считается каждый раз заново. И теги, введённые с разным регистром, остаются разными, пока вы не сделаете задание по желанию.

## Ответы

1. Потому что связь многие-ко-многим: у статьи несколько тегов, у тега несколько статей. В колонке её не выразить — поиск превращается в `like '%go%'`, который находит `golang`, а переименование и подсчёт становятся невозможными.
2. Потому что `do nothing` не считает конфликтующую строку затронутой, и `returning` не вернёт `id`. Бессмысленное обновление делает строку затронутой, и `id` приходит и для нового тега, и для существующего.
3. `errors.Is` отвечает «это та самая ошибка?», сравнивая со значением. `errors.As` ищет в цепочке ошибку нужного **типа** и кладёт её в переменную — после чего у неё можно спросить, например, код.

## Источники

- [SQLite: ON CONFLICT и upsert](https://sqlite.org/lang_upsert.html)
- [SQLite: коды ошибок](https://sqlite.org/rescode.html)
- [Go: пакет errors](https://pkg.go.dev/errors)
