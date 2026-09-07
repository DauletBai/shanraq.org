# CRUD: four actions on an article, and life without a slice

_Лид (summary):_ **The thirtieth lesson of the Go course. Articles move out of a slice and into a database, and survive a restart. A Store type that hides the database from the rest of the code, four actions on an article, and one non-obvious thing worth a section of its own: updating what is not there is not an error, and RowsAffected is the only way to find out.**

## Why this matters

The blog still keeps its articles in a slice. The form from the forms lesson writes into memory, and everything written lives until `Ctrl+C`.

Last lesson the program learned to read one row from a database. Today it learns the rest — create, change, delete — and the slice goes for good.

Those four actions are named by their initials: create, read, update, delete. CRUD is not a technology, just the list of what can be done to a record at all.

## The whole thing at once

A new folder, `go mod init sabaq27`, the same driver:

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

var ErrNotFound = errors.New("article not found")

// Store is everything the blog can do with articles. Only the actions are
// visible from outside; that there is a database inside stays inside.
type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("database does not answer: %w", err)
	}
	if _, err := db.Exec(`create table if not exists articles (
		id    integer primary key,
		slug  text    not null unique,
		title text    not null,
		words integer not null default 0
	) strict`); err != nil {
		db.Close()
		return nil, fmt.Errorf("schema: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

// Create is the C: put an article in and learn its number.
func (s *Store) Create(a Article) (int64, error) {
	res, err := s.db.Exec(
		`insert into articles (slug, title, words) values (?, ?, ?)`,
		a.Slug, a.Title, a.Words)
	if err != nil {
		return 0, fmt.Errorf("create %q: %w", a.Slug, err)
	}
	return res.LastInsertId()
}

// Get is the R for a single article.
func (s *Store) Get(slug string) (Article, error) {
	var a Article
	err := s.db.QueryRow(
		`select id, slug, title, words from articles where slug = ?`, slug).
		Scan(&a.ID, &a.Slug, &a.Title, &a.Words)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return Article{}, fmt.Errorf("%q: %w", slug, ErrNotFound)
	case err != nil:
		return Article{}, fmt.Errorf("read %q: %w", slug, err)
	}
	return a, nil
}

// All is the R for a list.
func (s *Store) All() ([]Article, error) {
	rows, err := s.db.Query(`select id, slug, title, words from articles order by id`)
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	defer rows.Close()

	var out []Article
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.Slug, &a.Title, &a.Words); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// Update is the U. Note RowsAffected: without it we never learn that there
// was nothing to change.
func (s *Store) Update(a Article) error {
	res, err := s.db.Exec(
		`update articles set title = ?, words = ? where slug = ?`,
		a.Title, a.Words, a.Slug)
	if err != nil {
		return fmt.Errorf("update %q: %w", a.Slug, err)
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

// Delete is the D, and RowsAffected tells the same story.
func (s *Store) Delete(slug string) error {
	res, err := s.db.Exec(`delete from articles where slug = ?`, slug)
	if err != nil {
		return fmt.Errorf("delete %q: %w", slug, err)
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

	id, err := store.Create(Article{Slug: "dala", Title: "About the steppe", Words: 400})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("created, number:", id)

	if _, err := store.Create(Article{Slug: "dala", Title: "The steppe again"}); err != nil {
		fmt.Println("the same address twice:", err)
	}

	a, _ := store.Get("dala")
	fmt.Printf("read back: %q, %d words\n", a.Title, a.Words)

	if err := store.Update(Article{Slug: "dala", Title: "The steppe in spring", Words: 450}); err != nil {
		log.Fatal(err)
	}
	a, _ = store.Get("dala")
	fmt.Printf("after the update: %q, %d words\n", a.Title, a.Words)

	err = store.Update(Article{Slug: "kokek", Title: "Nothing"})
	fmt.Println("updating a missing one:", errors.Is(err, ErrNotFound), "—", err)

	if err := store.Delete("dala"); err != nil {
		log.Fatal(err)
	}
	err = store.Delete("dala")
	fmt.Println("deleting twice:", errors.Is(err, ErrNotFound), "—", err)

	list, _ := store.All()
	fmt.Println("articles left:", len(list))
}
```

`go run .`:

```
created, number: 1
the same address twice: create "dala": constraint failed: UNIQUE constraint failed: articles.slug (2067)
read back: "About the steppe", 400 words
after the update: "The steppe in spring", 450 words
updating a missing one: true — "kokek": article not found
deleting twice: true — "dala": article not found
articles left: 0
```

## Taking it apart

### `Store`, the type that hides the database

```go
type Store struct{ db *sql.DB }
```

One field, and lower case, so it is invisible outside the package. What leaves are the actions only: `Create`, `Get`, `All`, `Update`, `Delete`.

It is the trick from the packages lesson, but now its cost and its payoff are visible. A handler calling `store.Get("dala")` knows nothing about SQL, about SQLite, or that a database exists at all. Want Postgres tomorrow? One file changes. Want an in-memory store for tests? It drops in, once those five actions are described by an interface.

And the other side: until such a type exists, fragments of SQL spread through the handlers, and gathering them back costs a day.

### `Create`: `insert` and the number the database gave

```go
res, err := s.db.Exec(
	`insert into articles (slug, title, words) values (?, ?, ?)`,
	a.Slug, a.Title, a.Words)
...
return res.LastInsertId()
```

`Exec` is for statements that return nothing: there are no rows in the answer, only a result. From it you take `LastInsertId`, the number the database assigned itself.

This works in SQLite; in Postgres `LastInsertId` is not supported at all, and you use `QueryRow` with `returning id` instead. Worth knowing in advance, so a move is not a surprise.

### `Get`: one row and an error of our own

```go
case errors.Is(err, sql.ErrNoRows):
	return Article{}, fmt.Errorf("%q: %w", slug, ErrNotFound)
```

We translate `sql.ErrNoRows` into `ErrNotFound`, which the previous lesson covered. What matters here is different: `ErrNotFound` is declared **in this package**, next to `Store`, and it is what the handlers will check. Past this point the word `sql` should not appear.

### `All`: a list and the error after the loop

```go
return out, rows.Err()
```

`rows.Err()` is returned on the last line together with the list. If the read broke halfway, the caller gets both a stump and an error — and the error says the list is not to be trusted.

### `Update`: why `RowsAffected` is not optional

Here is the heart of the lesson, and it is about databases, not about Go.

Ask the database to change an article that is not there. There will be **no error**:

```
update dala    err=<nil>  RowsAffected=1
update kokek   err=<nil>  RowsAffected=0
delete dala    err=<nil>  RowsAffected=1
delete kokek   err=<nil>  RowsAffected=0
```

`err=<nil>` in both cases. The database did exactly what was asked: it found the matching rows — there were zero — and changed all zero of them. From where it stands, that succeeded.

The only way to learn the truth is to ask how many rows were touched:

```go
n, err := res.RowsAffected()
if n == 0 {
	return fmt.Errorf("%q: %w", a.Slug, ErrNotFound)
}
```

Without those three lines an edit form tells the reader "saved" having saved nothing. No error comes from the database, from Go, or from the log — it surfaces a week later, when someone cannot find their edit.

> **Picture it.** A letter sent to an address that does not exist, which the post office quietly carries away. There is a receipt and no delivery, and the only way to find out is to ask how many letters arrived.

### `Delete`: the same story

```
delete kokek   err=<nil>  RowsAffected=0
```

Deleting an article that is not there is not an error either. The same check, for the same reason: a delete button must not answer "deleted" when there was nothing to delete.

### The database checks uniqueness, not you

```
the same address twice: create "dala": constraint failed: UNIQUE constraint failed: articles.slug (2067)
```

The temptation to check in code — `Get` first, and if nothing was found then `Create` — is understandable and wrong. Between the two queries somebody else inserts the same row, and the check saves nothing. The database checks it atomically, because `unique` is written in the schema.

The right way is to try and then read the error. Reading it by its message text is poor; later lessons bring a way to ask an error for its code. For now it is enough to know who is in charge here.

## The lesson map

![Lesson map: four actions, and one of them answers in silence](/static/course/go/map-crud-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why does `Store` hide `*sql.DB` in an unexported field?
2. What does `Exec` return if the `update` matched no rows?
3. Why does the database check the address for uniqueness rather than code before the insert?

## Exercise

**Required.** Add an edit page to the blog. `GET /edit/{slug}` shows a form with the title already filled in; `POST /edit/{slug}` changes the article through `Update` and answers with a redirect to `/read/{slug}`. If the article is not there, `404` — and it does not matter whether you learned that from `Get` or from `RowsAffected`.

You do not have to work out the move to a database yourself: it is done in [step-15](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-15) — compare against it once you have written the edit page.

**Optional.**

- Remove the `RowsAffected` check from `Update` and see what the form shows when editing an article that does not exist.
- Make `Delete` soft: set a `deleted_at` instead of removing the row, and keep such articles out of the list.
- Add a delete button to the article page. Work out which method it should send with, and why not `GET`.

## Where this goes in your blog

The slice goes and the database stays: this is the step the whole module was for. The blog now survives a restart.

Debts. The schema is still created by the program at start-up, which stops working once there are ten tables and they begin to change — that is the next lesson, on migrations. The query texts sit inside the methods, and at twenty of them you will want a place of their own. And `Update` writes every field at once: if two people edit one article, the second overwrites the first without noticing.

## Answers

1. So that the rest of the code does not know there is a database inside. Only the actions are visible from outside, so swapping SQLite for Postgres, or for an in-memory store in tests, touches one file rather than every handler.
2. Nothing alarming: `err` is `nil`, because the statement succeeded — zero rows simply matched. Telling "changed" from "there was nothing to change" is only possible through `res.RowsAffected()`.
3. Because a check in code is two queries with a gap between them: another query inserts the same row in that gap and both checks pass. `unique` in the schema is enforced atomically, so the right order is to attempt the insert and read the error.

## Sources

- [database/sql: Result](https://pkg.go.dev/database/sql#Result)
- [database/sql: DB.Exec](https://pkg.go.dev/database/sql#DB.Exec)
- [SQLite: UPDATE](https://sqlite.org/lang_update.html)
- [SQLite: result codes](https://sqlite.org/rescode.html)
