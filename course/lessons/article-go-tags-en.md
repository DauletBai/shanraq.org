# Tags: the third table nobody sees on the page

_Лид (summary):_ **The thirty-third lesson of the Go course. An article has many tags and a tag has many articles — and that is kept in a table of its own rather than a list in a column. A composite key, on conflict for the sake of an id, a cascade on delete. And the promised way to ask an error for its code instead of reading its text.**

## Why this matters

Tags look harmless: a couple of words under the title. The temptation to make them a column — `tags text` holding a comma-separated list — is strong, and it is a mistake.

Ask such a table for "the articles tagged go" and it begins:

```sql
select slug, tags from a where tags like '%go%';
```

```
┌────────┬─────────────┐
│  slug  │    tags     │
├────────┼─────────────┤
│ dala   │ go,steppe   │
│ golang │ golang,news │
└────────┴─────────────┘
```

`golang` came back because `go` is inside it. Then it turns out a tag cannot be renamed in one statement, cannot be counted, and cannot be stopped from being added twice.

The reason is simple: this link is **many-to-many**. An article has several tags and a tag has several articles. A column cannot express that — such a link gets a third table.

## The whole thing at once

A new folder, `go mod init sabaq30`:

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

// uniqueViolation is the code SQLite answers a broken unique with.
// The codes are listed at sqlite.org/rescode.html.
const uniqueViolation = 2067

// ErrTaken says the address is in use. That is what the handler needs,
// not a number out of the database.
var ErrTaken = errors.New("address is taken")

func main() {
	os.Remove("blog.db")
	db, err := sql.Open("sqlite", "blog.db?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	setup(db)

	fmt.Println("added dala  :", add(db, "dala", "The steppe", "go", "steppe"))
	fmt.Println("added salem :", add(db, "salem", "Hello", "go", "blog"))
	fmt.Println("same address:", add(db, "dala", "Another", "go"))

	fmt.Println("tags of dala        :", tagsOf(db, "dala"))
	fmt.Println("articles with go    :", byTag(db, "go"))
	fmt.Println("articles with steppe:", byTag(db, "steppe"))

	db.Exec(`delete from articles where slug = 'dala'`)
	fmt.Println("after deleting dala:")
	fmt.Println("  articles with go:", byTag(db, "go"))
	fmt.Println("  pairs in all    :", count(db, `select count(*) from article_tags`))
	fmt.Println("  tags in all     :", count(db, `select count(*) from tags`))
}

// add files an article and its tags. Two tables, one action, so one
// transaction: neither half is of any use on its own.
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
		// on conflict ... do update is the only way to get an id for a new tag
		// and for one that already exists: do nothing would return nothing.
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

// wrap asks the error for its code instead of reading its text.
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
		return "error: " + err.Error()
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
		return "error: " + err.Error()
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return "error: " + err.Error()
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return "error: " + err.Error()
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

The output:

```
added dala  : <nil>
added salem : <nil>
same address: "dala": address is taken
tags of dala        : go, steppe
articles with go    : [dala salem]
articles with steppe: [dala]
after deleting dala:
  articles with go: [salem]
  pairs in all    : 2
  tags in all     : 3
```

## Taking it apart

### The third table and a composite key

`article_tags` stores nothing of its own. It has two columns and both are references: to an article and to a tag. One row means "this article has this tag".

Such a table is called a join table, and its primary key is **composite** — two columns at once:

```sql
primary key (article_id, tag_id)
```

That is not a decorative detail. A composite key means the pair "article + tag" can occur only once, and the database will not let a tag be added twice:

```
Error: stepping, UNIQUE constraint failed: article_tags.article_id, article_tags.tag_id (19)
```

There is no need to check that in code — and, as the lesson on injection showed, no way either: another query slips in between the check and the insert.

> **Picture it.** A guest list where each line says "this person at this wedding". Neither the person nor the wedding is described in that line: the line describes only their meeting.

### `on conflict ... do update` for the sake of an `id`

A tag may be new or may already exist — and the `id` is needed either way. Hence a line that looks odd:

```sql
insert into tags (name) values (?)
on conflict (name) do update set name = excluded.name
returning id
```

It reads: try to insert; if that name is already there, "update" it with the very same value; either way give back the `id`.

The update is pointless on purpose — it is there only so the row counts as touched and `returning` hands its `id` over. `do nothing` will not do: a conflicting row is then not counted as touched, `returning` returns nothing, and `Scan` gets `sql.ErrNoRows`. It is a well-known trap, and this is how it is stepped around.

### The transaction is not decoration here

In the previous lesson there was nowhere to put a transaction. Now there is: `add` writes to two tables, and half the work is worse than none. An article without tags is survivable; a tag created for an article that never went in is rubbish for ever.

So `Begin`, `defer tx.Rollback()` and one `Commit` at the end — exactly what we took apart, and needed for real for the first time.

### `errors.As`: ask the error for its code, do not read its text

In the CRUD lesson I promised a way to ask an error for its code would come later. Here it is.

```go
var serr *sqlite.Error
if errors.As(err, &serr) && serr.Code() == uniqueViolation {
	return fmt.Errorf("%q: %w", slug, ErrTaken)
}
```

`errors.Is` answers the question "is this that error?". `errors.As` answers a different one: "is there an error of this type inside? if so, hand it to me". It unwraps the `%w` chain just as `Is` does, but puts what it finds into a variable — and that variable has methods, here `Code()`.

Comparing the text of an error is not allowed: it changes between versions and between languages. The code `2067` never changes — it is `SQLITE_CONSTRAINT_UNIQUE`, written down in the database's documentation.

Notice what `wrap` does: it **translates** an error of the database into an error of the blog. What leaves is `ErrTaken`, and the handler needs to know neither SQLite nor numbers. Want Postgres? One function changes.

One more thing: the driver is now imported by name rather than under `_`. The underscore in the SQLite lesson meant "needed only so it registers"; today we needed its type, and the name came back.

### `on delete cascade`: the links go, the tags stay

What happens to tags when an article is deleted is there in the output. There were four pairs; two are left. The database removed them itself, because the schema says:

```sql
article_id integer not null references articles (id) on delete cascade
```

And three tags are left, though "steppe" is of no use to anyone now. That is right: the `cascade` sits only on the reference to the article. Deleting a tag along with its last article is a decision for the blog, not for the database.

The other side of it: a tag that is still referenced cannot be deleted:

```
Error: stepping, FOREIGN KEY constraint failed (19)
```

Remove the links first, the tag second. This works, of course, only while the `foreign_keys` pragma is on — the one that has been in the connection string since the SQLite lesson.

### Putting the tags back together: `group_concat`

Showing an article's tags as one string is the database's work, not Go's:

```sql
select group_concat(t.name, ', ') from tags t join article_tags at ...
```

`group_concat` glues the values of every matching row together with a separator. Out comes `go, steppe` — exactly what goes under the title.

Two things to watch. The order inside `group_concat` is not guaranteed — if it matters, state it. And if there are no tags at all the function returns `NULL`, which `Scan` will not accept into a plain string; in the blog that is cured with `coalesce` or by reading into a `sql.NullString`.

## The lesson map

![Lesson map: two tables, and the link between them in a third](/static/course/go/map-tags-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why are tags not kept as a list in a column of the article?
2. Why does `on conflict` carry a pointless update instead of `do nothing`?
3. How does `errors.As` differ from `errors.Is`?

## Exercise

**Required.** Add tags to the blog. A migration changes the schema: two new tables and the join table. The form takes tags as a comma-separated string, `Store.Add` files the article and its tags in one transaction, the article page shows the tags, and `/tag/{name}` lists the articles carrying one.

All of it is already written in [step-18](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-18) — compare against it once you have written your own.

**Optional.**

- Bring tags to one shape before inserting: trim the spaces and lower the case, so `Go` and `go` do not become two tags.
- Show a tag cloud on the front page: the name and the number of articles. You will need `group by` and `count`.
- Take `on delete cascade` out and delete an article again. Read the error and explain it out loud.

## Where this goes in your blog

The third table is the first in the blog with no page of its own. It is there for the links rather than for the reader, and that is normal: real databases hold many such tables.

Debts. Tags cannot yet be found by part of a word — the next lesson, on search, sees to that. The tag cloud is counted afresh every time. And tags typed in different cases stay different until you do the optional exercise.

## Answers

1. Because the link is many-to-many: an article has several tags and a tag has several articles. A column cannot express it — the search turns into `like '%go%'`, which finds `golang`, and renaming and counting become impossible.
2. Because `do nothing` does not count a conflicting row as touched, so `returning` gives no `id`. The pointless update makes the row touched, and the `id` comes back both for a new tag and for one that exists.
3. `errors.Is` answers "is this that error?" by comparing against a value. `errors.As` looks through the chain for an error of the right **type** and puts it in a variable — after which you can ask it for, say, a code.

## Sources

- [SQLite: ON CONFLICT and upsert](https://sqlite.org/lang_upsert.html)
- [SQLite: result codes](https://sqlite.org/rescode.html)
- [Go: the errors package](https://pkg.go.dev/errors)
