# CRUD: четыре действия над статьёй и жизнь без среза

_Лид (summary):_ **Тридцатый урок курса по Go. Статьи переезжают из среза в базу и переживают перезапуск. Тип Store, который прячет базу от остального кода, четыре действия над статьёй и одна неочевидная вещь, стоящая целого раздела: обновить то, чего нет, — не ошибка, и узнать об этом можно только по RowsAffected.**

## Зачем это нужно

Блог до сих пор держит статьи в срезе. Форма из урока про формы пишет в память, и всё написанное живёт до `Ctrl+C`.

В прошлом уроке программа научилась читать из базы одну строку. Сегодня она научится всему остальному — создавать, менять, удалять, — и срез уйдёт совсем.

Эти четыре действия так и называют, по первым буквам: create, read, update, delete. CRUD — не какая-то технология, а просто список того, что вообще можно сделать с записью.

## Сразу целиком

Новая папка, `go mod init sabaq27`, драйвер тот же:

```
go get modernc.org/sqlite
```

`main.go`:

```go
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

type Article struct {
	ID    int64
	Slug  string
	Title string
	Words int
}

var ErrNotFound = errors.New("статья не найдена")

// Store — всё, что блог умеет делать со статьями. Снаружи видны только
// действия; то, что внутри база, за пределы этого типа не выходит.
type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("открыть %s: %w", path, err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("база не отвечает: %w", err)
	}
	if _, err := db.Exec(`create table if not exists articles (
		id    integer primary key,
		slug  text    not null unique,
		title text    not null,
		words integer not null default 0
	) strict`); err != nil {
		db.Close()
		return nil, fmt.Errorf("схема: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

// Create — «C»: положить статью и узнать её номер.
func (s *Store) Create(a Article) (int64, error) {
	res, err := s.db.Exec(
		`insert into articles (slug, title, words) values (?, ?, ?)`,
		a.Slug, a.Title, a.Words)
	if err != nil {
		return 0, fmt.Errorf("создать %q: %w", a.Slug, err)
	}
	return res.LastInsertId()
}

// Get — «R» для одной статьи.
func (s *Store) Get(slug string) (Article, error) {
	var a Article
	err := s.db.QueryRow(
		`select id, slug, title, words from articles where slug = ?`, slug).
		Scan(&a.ID, &a.Slug, &a.Title, &a.Words)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return Article{}, fmt.Errorf("%q: %w", slug, ErrNotFound)
	case err != nil:
		return Article{}, fmt.Errorf("прочитать %q: %w", slug, err)
	}
	return a, nil
}

// All — «R» для списка.
func (s *Store) All() ([]Article, error) {
	rows, err := s.db.Query(`select id, slug, title, words from articles order by id`)
	if err != nil {
		return nil, fmt.Errorf("список: %w", err)
	}
	defer rows.Close()

	var out []Article
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.Slug, &a.Title, &a.Words); err != nil {
			return nil, fmt.Errorf("разобрать строку: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// Update — «U». Обратите внимание на RowsAffected: без него мы не узнаем,
// что менять было нечего.
func (s *Store) Update(a Article) error {
	res, err := s.db.Exec(
		`update articles set title = ?, words = ? where slug = ?`,
		a.Title, a.Words, a.Slug)
	if err != nil {
		return fmt.Errorf("изменить %q: %w", a.Slug, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("%q: %w", a.Slug, ErrNotFound)
	}
	return nil
}

// Delete — «D», и та же история с RowsAffected.
func (s *Store) Delete(slug string) error {
	res, err := s.db.Exec(`delete from articles where slug = ?`, slug)
	if err != nil {
		return fmt.Errorf("удалить %q: %w", slug, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("%q: %w", slug, ErrNotFound)
	}
	return nil
}

func main() {
	os.Remove("blog.db")
	store, err := Open("blog.db")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer store.Close()

	id, err := store.Create(Article{Slug: "dala", Title: "О степи", Words: 400})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("создали, номер:", id)

	if _, err := store.Create(Article{Slug: "dala", Title: "Ещё раз о степи"}); err != nil {
		fmt.Println("второй раз тот же адрес:", err)
	}

	a, _ := store.Get("dala")
	fmt.Printf("прочитали: %q, %d слов\n", a.Title, a.Words)

	if err := store.Update(Article{Slug: "dala", Title: "О степи весной", Words: 450}); err != nil {
		log.Fatal(err)
	}
	a, _ = store.Get("dala")
	fmt.Printf("после правки: %q, %d слов\n", a.Title, a.Words)

	err = store.Update(Article{Slug: "kokek", Title: "Ничего"})
	fmt.Println("правка несуществующей:", errors.Is(err, ErrNotFound), "—", err)

	if err := store.Delete("dala"); err != nil {
		log.Fatal(err)
	}
	err = store.Delete("dala")
	fmt.Println("удаление дважды:", errors.Is(err, ErrNotFound), "—", err)

	list, _ := store.All()
	fmt.Println("осталось статей:", len(list))
}
```

`go run .`:

```
создали, номер: 1
второй раз тот же адрес: создать "dala": constraint failed: UNIQUE constraint failed: articles.slug (2067)
прочитали: "О степи", 400 слов
после правки: "О степи весной", 450 слов
правка несуществующей: true — "kokek": статья не найдена
удаление дважды: true — "dala": статья не найдена
осталось статей: 0
```

## Разбор

### `Store` — тип, который прячет базу

```go
type Store struct{ db *sql.DB }
```

Одно поле, и оно с маленькой буквы — значит, снаружи пакета его не видно. Наружу выходят только действия: `Create`, `Get`, `All`, `Update`, `Delete`.

Это тот же приём, что и в уроке про пакеты, но теперь у него есть цена и выгода, которые видно. Обработчик, который зовёт `store.Get("dala")`, не знает ни про SQL, ни про SQLite, ни про то, что база вообще есть. Захотите завтра Postgres — поменяется один файл. Захотите в тестах хранилище в памяти — оно подставится, если описать эти пять действий интерфейсом.

И обратная сторона: пока такой тип не заведён, куски SQL расползаются по обработчикам, и собрать их обратно потом стоит дня.

### `Create`: `insert` и номер, который дала база

```go
res, err := s.db.Exec(
	`insert into articles (slug, title, words) values (?, ?, ?)`,
	a.Slug, a.Title, a.Words)
...
return res.LastInsertId()
```

`Exec` — для запросов, которые ничего не возвращают: строк в ответе нет, есть только результат. Из него берут `LastInsertId` — номер, который база присвоила сама.

В SQLite это работает, а вот в Postgres `LastInsertId` не поддерживается вовсе: там берут `QueryRow` с `returning id`. Знать об этом стоит заранее, чтобы переезд не стал сюрпризом.

### `Get`: одна строка и своя ошибка

```go
case errors.Is(err, sql.ErrNoRows):
	return Article{}, fmt.Errorf("%q: %w", slug, ErrNotFound)
```

`sql.ErrNoRows` мы переводим в `ErrNotFound` — это разбирали в прошлом уроке. Здесь важно другое: `ErrNotFound` объявлен **в этом же пакете**, рядом со `Store`, и именно его будут проверять обработчики. Дальше по коду слово `sql` встречаться не должно.

### `All`: список и ошибка после цикла

```go
return out, rows.Err()
```

`rows.Err()` вернули последней строкой вместе со списком. Если чтение сломалось на середине, вызывающий получит и обрубок, и ошибку — и по ошибке поймёт, что списку верить нельзя.

### `Update`: почему `RowsAffected` обязателен

Вот главное в уроке, и это не про Go, а про базы.

Попросите базу изменить статью, которой нет. Ошибки **не будет**:

```
update dala    ошибка=<nil>  RowsAffected=1
update kokek   ошибка=<nil>  RowsAffected=0
delete dala    ошибка=<nil>  RowsAffected=1
delete kokek   ошибка=<nil>  RowsAffected=0
```

`ошибка=<nil>` в обоих случаях. База сделала ровно то, что просили: нашла подходящие строки — их оказалось ноль — и изменила все ноль. С её точки зрения всё прошло успешно.

Единственный способ узнать правду — спросить, сколько строк задело:

```go
n, err := res.RowsAffected()
if n == 0 {
	return fmt.Errorf("%q: %w", a.Slug, ErrNotFound)
}
```

Без этих трёх строк форма правки скажет читателю «сохранено», не сохранив ничего. Ошибку не покажет ни база, ни Go, ни журнал — она проявится через неделю, когда человек не найдёт свою правку.

> **Образ.** Письмо, отправленное на несуществующий адрес, которое почта молча уносит. Квитанция есть, доставки нет, и узнать об этом можно только спросив, сколько писем доставлено.

### `Delete`: та же история

```
delete kokek   ошибка=<nil>  RowsAffected=0
```

Удаление несуществующей статьи — тоже не ошибка. Проверка та же, и по той же причине: кнопка «удалить» не должна отвечать «удалено», когда удалять было нечего.

### Уникальность проверяет база, а не вы

```
второй раз тот же адрес: создать "dala": constraint failed: UNIQUE constraint failed: articles.slug (2067)
```

Соблазн проверить в коде — сначала `Get`, и если не нашли, тогда `Create` — понятен и неправилен. Между двумя запросами кто-то другой успеет вставить ту же строку, и проверка не спасёт. База проверяет это атомарно, потому что `unique` записан в схеме.

Правильно — попробовать и разобрать ошибку. Разбирать её по тексту сообщения плохо; в следующих уроках появится способ спросить у ошибки её код. Пока достаточно понимать, кто здесь главный.

## Карта урока

![Карта урока: четыре действия, и одно отвечает молча](/static/course/go/map-crud-ru.svg)

## Скажите своими словами

Не подглядывая, ответьте вслух или на бумаге. Ответы — в конце урока.

1. Зачем `Store` прячет `*sql.DB` в неэкспортируемом поле?
2. Что вернёт `Exec`, если `update` не нашёл ни одной строки?
3. Почему уникальность адреса проверяет база, а не код перед вставкой?

## Задание

**Обязательное.** Добавьте в блог страницу правки. `GET /edit/{slug}` показывает форму с уже заполненным заголовком, `POST /edit/{slug}` меняет статью через `Update` и отвечает редиректом на `/read/{slug}`. Если статьи нет — `404`, и не важно, узнали вы об этом из `Get` или из `RowsAffected`.

Сам переезд блога на базу разбирать не нужно: он уже сделан в [step-15](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-15) — сверьтесь с ним после того, как напишете страницу правки.

**По желанию.**

- Уберите проверку `RowsAffected` из `Update` и посмотрите, что покажет форма при правке несуществующей статьи.
- Сделайте `Delete` мягким: вместо удаления ставьте `deleted_at` и не показывайте такие статьи в списке.
- Добавьте кнопку «удалить» на страницу статьи. Подумайте, каким методом она должна отправлять запрос и почему не `GET`.

## Куда это встанет в блоге

Срез уходит, база остаётся: это и есть шаг, ради которого затевался весь модуль. Блог теперь переживает перезапуск.

Долги. Схему по-прежнему создаёт сама программа при старте — так нельзя, когда таблиц станет десять и они начнут меняться; это следующий урок, про миграции. Тексты запросов лежат прямо в методах, и когда их станет двадцать, захочется отдельного места. И `Update` меняет все поля разом: если два человека правят одну статью, второй затрёт первого, не заметив.

## Ответы

1. Чтобы остальной код не знал, что внутри база. Снаружи видны только действия, поэтому замена SQLite на Postgres или на хранилище в памяти для тестов затрагивает один файл, а не все обработчики.
2. Ничего плохого: `err` будет `nil`, потому что запрос выполнился успешно — просто подошло ноль строк. Отличить «изменили» от «нечего было менять» можно только через `res.RowsAffected()`.
3. Потому что проверка в коде состоит из двух запросов, а между ними есть промежуток: другой запрос успеет вставить ту же строку, и обе проверки пройдут. `unique` в схеме проверяется атомарно, поэтому правильный порядок — попробовать вставить и разобрать ошибку.

## Источники

- [database/sql: Result](https://pkg.go.dev/database/sql#Result)
- [database/sql: DB.Exec](https://pkg.go.dev/database/sql#DB.Exec)
- [SQLite: изменение данных](https://sqlite.org/lang_update.html)
- [SQLite: коды ошибок](https://sqlite.org/rescode.html)
