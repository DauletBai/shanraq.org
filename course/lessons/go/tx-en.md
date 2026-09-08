# Transactions and SQL injection: a query is not a string you glue together

_Лид (summary):_ **The thirty-second lesson of the Go course. Two rules that keep a database whole. Data is never glued into the text of a query: a tampered address returns three rows instead of one, and through a parameter it returns none. And what has to happen together is wrapped in a transaction: without one, a counter went up for a move that never took place.**

## Why this matters

The queries in these lessons have been short and parameterised so far, and that looked like a matter of style. It is not style.

A string sent by a reader is not text but somebody else's will. Glue it into a query and the reader adds their own query to yours. That is how more than one database has been carried off, and the way it is done fits in a single form field.

The second rule is about something else. One action in a blog is often two queries: move an article and record that in a counter, take money off one account and add it to another. If the second query fails, the first is already done — and the database is left in a state that does not occur in life.

## The whole thing at once

A new folder, `go mod init sabaq29`, the same driver:

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	os.Remove("blog.db")
	db, err := sql.Open("sqlite", "blog.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	setup(db)

	fmt.Println("== glued strings")
	fmt.Println("ordinary address:", bad(db, "dala"))
	fmt.Println("tampered with:   ", bad(db, "dala' or '1'='1"))

	fmt.Println("== a parameter")
	fmt.Println("ordinary address:", good(db, "dala"))
	fmt.Println("tampered with:   ", good(db, "dala' or '1'='1"))

	fmt.Println("== ordering")
	fmt.Println("order by ?        ->", byParam(db, "title"))
	fmt.Println("from the list     ->", byAllowed(db, "title"))
	fmt.Println("a field of theirs ->", byAllowed(db, "views; drop table articles"))
}

// bad glues the value into the text of the query. Never write this.
func bad(db *sql.DB, slug string) string {
	q := fmt.Sprintf(`select count(*) from articles where slug = '%s'`, slug)
	var n int
	if err := db.QueryRow(q).Scan(&n); err != nil {
		return "error: " + err.Error()
	}
	return fmt.Sprint("rows: ", n)
}

// good hands the value over separately, as a value.
func good(db *sql.DB, slug string) string {
	var n int
	if err := db.QueryRow(`select count(*) from articles where slug = ?`, slug).Scan(&n); err != nil {
		return "error: " + err.Error()
	}
	return fmt.Sprint("rows: ", n)
}

func byParam(db *sql.DB, col string) string {
	return list(db, `select slug from articles order by ?`, col)
}

// byAllowed checks the column name against a list we wrote ourselves. A name
// that is not on it is not explained, it is refused.
func byAllowed(db *sql.DB, col string) string {
	allowed := map[string]bool{"slug": true, "title": true, "views": true}
	if !allowed[col] {
		return "refused: no such field"
	}
	return list(db, `select slug from articles order by `+col)
}

func list(db *sql.DB, q string, args ...any) string {
	rows, err := db.Query(q, args...)
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

func setup(db *sql.DB) {
	if _, err := db.Exec(`create table articles (
		id    integer primary key,
		slug  text not null unique,
		title text not null,
		views integer not null default 0
	) strict`); err != nil {
		log.Fatal(err)
	}
	for _, a := range [][2]string{{"dala", "The yurt"}, {"salem", "The alphabet"}, {"shanyraq", "The world"}} {
		if _, err := db.Exec(`insert into articles (slug, title) values (?, ?)`, a[0], a[1]); err != nil {
			log.Fatal(err)
		}
	}
}
```

The output:

```
== glued strings
ordinary address: rows: 1
tampered with:    rows: 3
== a parameter
ordinary address: rows: 1
tampered with:    rows: 0
== ordering
order by ?        -> [dala salem shanyraq]
from the list     -> [salem shanyraq dala]
a field of theirs -> refused: no such field
```

## Taking it apart

### Three rows instead of one

Look at the first two lines of the output. The address `dala` found one article, which is right. The address `dala' or '1'='1` found three — that is, all of them.

What happened is plain once the query is put together by hand:

```sql
select count(*) from articles where slug = 'dala' or '1'='1'
```

The quote closed the address, and after it came a condition that is always true. The check on the address simply stopped existing. In exactly the same way `or 1=1` is added to a login form to get in without a password — and a `;` lets someone add not a condition but a whole second statement.

> **Picture it.** A filled-in form with an extra line written into it — and that line reads as part of the form itself. Not a note in the margin, but an item like the printed ones.

### `?` is not substitution, it is a separate delivery

The third and fourth lines: the same tampered address through `?` gives **none**.

None rather than an error, and that is the thing to understand correctly. The database received the query and the value **apart**. The text of the query was parsed before the value ever reached it, so quotes inside the value close nothing: they are letters. The database honestly looked for an article addressed `dala' or '1'='1` and found none.

Hence a rule with no exceptions: **everything that comes from outside goes through `?`**. Not "if it looks dangerous", not "if it is a string" — everything. Escaping quotes by hand is unnecessary and unwanted: that is the driver's work, and it does it better.

Postgres writes the places for values differently — `$1`, `$2` — but means the same thing.

### What cannot go through a `?`

The fifth line of the output is the most treacherous thing in the lesson:

```
order by ?        -> [dala salem shanyraq]
```

There is no error. There is no sorting either: the addresses came out in the order they were lying in. Through the allow-list the same `title` gives a different order — `[salem shanyraq dala]` — and that one is right: "The alphabet", "The world", "The yurt".

The reason is the reason it is safe. A value is data; a column name is part of the text of the query. The database got `order by 'title'` — sorting by a constant string, the same for every row. Sorting by that is meaningless, so the order stayed as it was.

While on sorting: SQLite compares strings by character codes rather than by the alphabet of a language. For Latin script the two are nearly the same; for Kazakh they are not, because `Қ` and `Ә` sit far from where the Kazakh alphabet puts them. An order that is right for a human being is separate work, and in a database it is called collation.

A `?` cannot carry a table name, a column name, `asc`/`desc`, or in some databases a `limit`. This is exactly where injection turns up in the code of people who "use parameters everywhere".

### An allow-list instead of ingenuity

If a column name does come from outside — from a link like `?sort=title` — it cannot be escaped and cannot be cleaned. It can only be **checked against a list you wrote**:

```go
allowed := map[string]bool{"slug": true, "title": true, "views": true}
	if !allowed[col] {
		return "refused: no such field"
	}
	return list(db, `select slug from articles order by `+col)
}
```

The last line of the output shows the refusal of `views; drop table articles`. Notice what this check does not contain: any attempt to guess what is dangerous. What is listed is allowed and the rest is not. A list of three words is more reliable than any cleaning, because it does not depend on your imagination about what a stranger will think of.

### A transaction: all or nothing

The second rule of the lesson. Take an action made of two queries: move an article to a new address and mark that in a counter.

```go
func move(db *sql.DB, from, to string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`update articles set views = views + 1 where slug = 'dala'`); err != nil {
		return err
	}
	if _, err := tx.Exec(`update articles set slug = ? where slug = ?`, to, from); err != nil {
		return err
	}
	return tx.Commit()
}
```

First, what happens **without** `Begin` — the same two `db.Exec` calls in a row. We move to an address that is taken, then to a free one:

```
before: dala=0 shanyraq=0
move: constraint failed: UNIQUE constraint failed: articles.slug (2067)
after: dala=1 shanyraq=0
move: <nil>
after: dala=2 koktem=0
```

The counter went up to `1` for a move that did not happen: the first query passed, the second fell over. By the end it says `2` for one real move. The database is not broken; it is worse than broken — it lies.

Now with a transaction, the same sequence:

```
before: dala=0 shanyraq=0
move: constraint failed: UNIQUE constraint failed: articles.slug (2067)
after: dala=0 shanyraq=0
move: <nil>
after: dala=1 koktem=0
```

After the failure `dala=0`: the counter is untouched and the article is where it was. After the success `dala=1`. A transaction is a promise from the database: either both queries or neither.

### `defer tx.Rollback()` — one line for every exit

There are four ways out of the function, and three of them have to roll the transaction back. Writing `tx.Rollback()` before every `return` is a reliable way to forget it once.

So one line goes in right after `Begin`. It works because **a rollback after a successful `Commit` does nothing**: the transaction is already over, `Rollback` returns `sql.ErrTxDone`, and that is not an error but a message saying "too late, it is all finished". We leave the value unchecked deliberately.

And the thing that is easy to forget: a transaction holds **one** connection from the pool for as long as it lives. Inside `move` you must not reach for `db`, only for `tx`; a query through `db` goes out on another connection, cannot see the unfinished changes, and may queue behind the transaction itself. A long transaction holds its connection and gives it to nobody.

### Two writers: `database is locked`

SQLite allows only one writer. Let us check: two transactions from two separate connections to one file.

```
first transaction writes: ok
second transaction writes: database is locked (5) (SQLITE_BUSY)
result: 2
```

The second did not wait; it was refused at once. The first calmly wrote its own, and the `+100` from the second was lost — rightly, since nobody committed it.

Readers, meanwhile, may be as many as you like. For a blog with one author writing, this is no limit at all; when there are many writers it becomes the reason to move to Postgres, where transactions do not get in each other's way so bluntly.

## The lesson map

![Lesson map: the value travels apart from the query, and the two actions together](/static/course/go/map-tx-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why does `?` protect against injection when nothing in your code escapes anything?
2. What does `order by ?` return, and why is that more dangerous than an error?
3. Why is `defer tx.Rollback()` written even where a `Commit` is coming?

## Exercise

**Required.** Find every place in your blog where a value reaches a query through `fmt.Sprintf` or `+`, and move them onto `?`. Then make the list sortable: `/?sort=title` and `/?sort=updated`, the column name through an allow-list, an unknown value quietly giving the default order.

The sorting is already written in [step-17](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-17) — compare against it once you have written your own.

**Optional.**

- Add a `deletions` table and delete an article together with a row in it, in one transaction. Break the second half on purpose and make sure the article is still there.
- Write a test that hands `Get` the address `x' or '1'='1` and demands `ErrNotFound`. Such a test catches gluing coming back in a later edit.
- Measure what happens if code inside a transaction reaches for `db` instead of `tx`.

## Where this goes in your blog

There is no gluing in the blog — parameters have been there since the very first query. What appears instead is a sortable list, and it is the first place where a column name arrives from outside.

Debts. There is nowhere serious to put a transaction yet: they are genuinely needed in the lesson on tags, where one article is written into two tables at once. And `Update` still writes every field at once.

## Answers

1. Because the value never enters the text of the query at all: the database gets a parsed query and the value separately. Quotes inside the value stay letters instead of becoming part of the SQL.
2. The order the rows were lying in — that is, no sorting at all. It is more dangerous than an error because nothing falls over: a wrong order is easy to miss, and trying to "fix" it by gluing opens the way to injection.
3. Because there are several ways out of the function, and forgetting the rollback once is enough. After a successful `Commit` that `Rollback` does nothing and returns `sql.ErrTxDone`.

## Sources

- [SQLite: SQL injection and parameters](https://sqlite.org/lang_expr.html#varparam)
- [Go: the database/sql package](https://pkg.go.dev/database/sql)
