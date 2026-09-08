# Search: FTS5, and what a parameter does not save you from

_Лид (summary):_ **The thirty-fourth lesson of the Go course. A search that works: an FTS5 virtual table, three triggers, ranking and highlighting. Plus a measured surprise: a parameter protects you from injection, but not from a person typing the word OR into the search box.**

## Why this matters

A search over a blog usually starts as `like '%word%'` — and the tags lesson showed where that ends: the query `%go%` found `golang`.

There are other troubles. `like` cannot search title and body at once with different weights. It cannot say which article fits better. It does not highlight what it found. And it reads the whole table: across three articles nobody notices, across three thousand they do.

SQLite has a separate tool for this — **FTS5**, full-text search. It is built in, nothing has to be installed, and it works through ordinary `database/sql`.

## The whole thing at once

A new folder, `go mod init sabaq31`, the same driver:

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

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

	for _, q := range []string{"steppe", "STEPPE", "step", "steppe green", "steppe OR", "\"", "nosuchword"} {
		fmt.Printf("%-12q -> %s\n", q, search(db, q))
	}
}

// ftsQuery turns what a person typed into a query FTS5 accepts. Every word
// goes in quotes, so an AND, an OR or a star inside stays a letter, and a star
// outside makes it match the beginning of a word.
func ftsQuery(s string) string {
	var parts []string
	for _, w := range strings.Fields(s) {
		w = strings.ReplaceAll(w, `"`, `""`)
		parts = append(parts, `"`+w+`"*`)
	}
	return strings.Join(parts, " AND ")
}

func search(db *sql.DB, q string) string {
	fts := ftsQuery(q)
	if fts == "" {
		return "empty query"
	}
	rows, err := db.Query(`select a.slug, snippet(articles_fts, 1, '[', ']', '…', 5)
		from articles_fts
		join articles a on a.id = articles_fts.rowid
		where articles_fts match ?
		order by bm25(articles_fts, 10.0, 1.0)`, fts)
	if err != nil {
		return "error: " + err.Error()
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var slug, frag string
		if err := rows.Scan(&slug, &frag); err != nil {
			return "error: " + err.Error()
		}
		out = append(out, slug+": "+frag)
	}
	if err := rows.Err(); err != nil {
		return "error: " + err.Error()
	}
	if len(out) == 0 {
		return "nothing found"
	}
	return strings.Join(out, " | ")
}

func setup(db *sql.DB) {
	for _, q := range []string{
		`create table articles (
			id    integer primary key,
			slug  text not null unique,
			title text not null,
			body  text not null
		) strict`,
		// content='articles' means the search table stores no text of its own
		// but looks into articles; content_rowid says by which column.
		`create virtual table articles_fts using fts5(
			title, body, content='articles', content_rowid='id'
		)`,
		`create trigger articles_ai after insert on articles begin
			insert into articles_fts (rowid, title, body)
			values (new.id, new.title, new.body);
		end`,
		`create trigger articles_ad after delete on articles begin
			insert into articles_fts (articles_fts, rowid, title, body)
			values ('delete', old.id, old.title, old.body);
		end`,
		`create trigger articles_au after update on articles begin
			insert into articles_fts (articles_fts, rowid, title, body)
			values ('delete', old.id, old.title, old.body);
			insert into articles_fts (rowid, title, body)
			values (new.id, new.title, new.body);
		end`,
	} {
		if _, err := db.Exec(q); err != nil {
			log.Fatal(err)
		}
	}
	for _, a := range [][3]string{
		{"steppe", "The steppe", "The steppe is yellow in summer and green in spring"},
		{"hello", "Hello, world", "The first entry. The blog starts here"},
		{"go", "The Go language", "Go is a compiled language"},
		{"shanyraq", "What a shanyraq is", "A shanyraq is the ring at the top of a yurt"},
		{"spring", "Spring", "In spring the steppe is green"},
		{"test", "What a test is", "A test says what broke when a change is made"},
	} {
		if _, err := db.Exec(
			`insert into articles (slug, title, body) values (?, ?, ?)`, a[0], a[1], a[2]); err != nil {
			log.Fatal(err)
		}
	}
}
```

The output:

```
"steppe"     -> steppe: The [steppe] is yellow in… | spring: In spring the [steppe] is…
"STEPPE"     -> steppe: The [steppe] is yellow in… | spring: In spring the [steppe] is…
"step"       -> steppe: The [steppe] is yellow in… | spring: In spring the [steppe] is…
"steppe green" -> steppe: The [steppe] is yellow in… | spring: …spring the [steppe] is [green]
"steppe OR"  -> nothing found
"\""         -> nothing found
"nosuchword" -> nothing found
```

## Taking it apart

### A virtual table: an index, not a copy

`create virtual table ... using fts5` makes a table that does not exist on disk in the ordinary sense. Inside is an inverted index: for every word, the list of articles it occurs in. That is how a book finds a page by a word in a second rather than by leafing through.

The two important words are in the brackets:

```sql
create virtual table articles_fts using fts5(
	title, body, content='articles', content_rowid='id'
)
```

`content='articles'` says: do **not** store the text itself, go to `articles` for it. Without that line the blog would keep every article twice. `content_rowid='id'` says which column the rows are matched by.

> **Picture it.** The index at the back of a book. It does not retell the book, it only says which page to look at. Tear out the index and the book remains; take away the book and the index is useless.

### Three triggers: an index that cannot fall behind

Since the text is not copied, keeping the two in step is our job. That is what triggers are for — small rules the database runs itself whenever a table changes.

There are exactly three: after insert, after delete, after update. The delete is written in a way that looks odd:

```sql
insert into articles_fts (articles_fts, rowid, title, body)
values ('delete', old.id, old.title, old.body);
```

It is not a mistake. For tables with `content=`, removal from the index is written as an insert of the `'delete'` command — together with the **old** values, or the index will not find what to take out. An update is a delete plus an insert, which is why the third trigger does both in a row.

Forget the triggers and the search will find deleted articles and miss new ones, in silence.

### `match`, not `like`

A query to the search is written with `match`:

```sql
where articles_fts match ?
```

On the left, the name of the search table; on the right, the query. Matching is by whole words rather than pieces of a string, so `steppe` will not find `steppes` but will find `Steppe` with a capital letter: FTS5 does not care about case. That is the second line of the output — `STEPPE` found the same thing.

### What a person types is not what FTS5 reads

Here is the heart of the lesson.

FTS5 has a query language of its own: `AND`, `OR`, `NOT`, quotes for a phrase, a star for the beginning of a word, `title:word` to search one column. A person in a search box knows nothing about that and types whatever they like.

Take `ftsQuery` out and hand over what was typed as it is:

```
"step"       -> nothing found
"steppe OR"  -> error: SQL logic error: fts5: syntax error near "" (1)
"\""         -> error: SQL logic error: unterminated string (1)
```

Notice that this is **not** injection: the value travelled through `?`, it never entered the text of the query, and no reader can add SQL of their own. But they broke the query all the same — because what parses it now is not SQL but FTS5, which has rules of its own.

Hence the rule: what was typed is always prepared before it reaches `match`. Our `ftsQuery` does the minimum that suffices:

```go
w = strings.ReplaceAll(w, `"`, `""`)
parts = append(parts, `"`+w+`"*`)
```

Every word goes in quotes — inside quotes an `OR` stops being a command and becomes a word. A quote inside a word is doubled, or it would close ours. The star outside means "a word that starts with this", so `step` finds `Steppe`. The words are joined with `AND`: what is found has all of them.

### `order by rank`: the smaller, the better

The order is set by `bm25`, a measure decades old that weighs how often a word occurs, how rare it is, and how long the text is.

```
┌───────┬────────┐
│ rowid │  rank  │
├───────┼────────┤
│ 1     │ -0.784 │
│ 5     │ -0.687 │
└───────┴────────┘
```

The numbers are negative, and **smaller means better** — which is why `order by rank` carries no `desc`. The first article fits more strongly than the second.

`bm25` takes weights, one per column:

```sql
order by bm25(articles_fts, 10.0, 1.0)
```

The title counts ten times what the body counts, and the gap widens: `-1.157` against `-0.687`. That is what a search is expected to do: an article with the word in its title should come first.

One oddity that makes it easy to decide ranking is broken: if there are very few articles and the word is in nearly all of them, `bm25` returns zero for every one. I saw that on three articles. It is not the code — a rare word is valuable because it is rare, and a tiny database has no rare words.

### `snippet`: show where it was found

A link with no explanation of why it was found is half an answer. `snippet` cuts a piece out around the word and frames it:

```sql
snippet(articles_fts, 1, '[', ']', '…', 5)
```

The arguments: the column (`1` is `body`, counting from zero), what to open with, what to close with, what marks the break, and how many words to show. A real blog puts `<mark>` and `</mark>` there instead of brackets — but then the template has to print that as markup rather than text, and the care from the templates lesson applies: only what the database returned may be highlighted, never what a reader sent.

### Kazakh letters: what works and what does not

Let us check how FTS5 breaks up words that are not Latin:

```sql
create virtual table t using fts5(w);
insert into t values ('café'), ('ёлка'), ('әже');
```

- searching for `cafe` finds `café`;
- searching for `елка` does not find `ёлка`;
- searching for `аже` does not find `әже`.

The reason is the `remove_diacritics` setting, which is on by default. It strips the mark where a letter is built out of a base and a mark — like `é` out of `e`. But `ё`, `ә`, `ң`, `қ` are not decorated letters, they are **letters of their own alphabet**, and there is nothing there to strip.

For a Kazakh blog that means something simple: a reader who types `алем` will not find `әлем`. The cure is not a setting but work: bringing the query and the text to one shape before indexing. We will not do it here — but knowing why a search "cannot find the obvious" is worth knowing in advance.

## The lesson map

![Lesson map: a word, an index and the order of the answer](/static/course/go/map-search-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. What does `content='articles'` do, and what would happen without it?
2. Why does `?` not save you from a person typing `OR` into the search box?
3. Why is `order by rank` written without `desc`?

## Exercise

**Required.** Give the blog a search. A migration creates the search table and the three triggers — and fills the index from the articles already written. A `GET /search?q=` form appears on the front page, the results come out highlighted, and an empty result says so.

All of it is written in [step-19](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-19) — compare against it once you have written your own.

**Optional.**

- Add tags to the search: a third column in the search table and the triggers that fill it.
- Print how many were found, and page through them with `limit` and `offset`.
- Measure `like '%steppe%'` against `match` on ten thousand rows. The difference explains what all this was for.

## Where this goes in your blog

Search is the first part of the blog that works with every article at once. And the first whose answer depends on how well somebody else's input was prepared.

Debts. The index knows nothing about tags. Kazakh letters are searched literally. And the search page has neither paging nor a memory of the query — the lesson on layout sees to that, when every page gets a frame of its own.

## Answers

1. It tells the search table not to store the text again but to take it from `articles` by `rowid`. Without it every article would sit in the database twice.
2. Because `?` protects **the text of the SQL query**: the value never enters it. But the value itself is then parsed by FTS5 under its own rules, and `OR` is a command to it. What breaks is not the SQL but the search query.
3. Because `bm25` returns negative numbers, and the better the match the smaller the number. Ascending means the best ones first.

## Sources

- [SQLite: FTS5](https://sqlite.org/fts5.html)
- [SQLite: the snippet and bm25 functions](https://sqlite.org/fts5.html#the_snippet_function)
