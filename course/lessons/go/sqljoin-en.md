# SQL: two tables, JOIN and foreign keys

_Лид (summary):_ **The twenty-eighth lesson of the Go course. An article acquires comments, and they do not fit inside it. A second table and the key between them, join against left join, cascading deletes, and the SQLite trap that leaves foreign keys silent and rows orphaned.**

## Why this matters

Last lesson there was one table and everything was simple: a row was an article.

Now comments arrive. One article may have none of them and may have fifty, and you cannot know in advance. There is nowhere to put them in the same row: columns named `comment1`, `comment2`, `comment3` are not a schema but a dead end.

The right answer is a second table that **points at** the first. Today we work out how to tie it on, how to collect both back into one answer, and what happens to the comments when the article is deleted.

## The whole thing at once

Open the same `blog.db` as last lesson, with the same three lines of setup:

```
.mode box
.headers on
pragma foreign_keys = on;
```

The second table:

```sql
create table comments (
    id         integer primary key,
    article_id integer not null references articles(id) on delete cascade,
    author     text    not null,
    body       text    not null
) strict;
```

There is one key line here, and it is the first thing we take apart.

## Taking it apart

### The second table points at the first

Comments do not live inside an article. They live in a table of their own and **point at** the article by number:

```sql
article_id integer not null references articles(id) on delete cascade
```

That is the foreign key: `comments.article_id` holds an `articles.id`. `references` tells the database that no stranger's number belongs there, and `on delete cascade` that an article takes its comments with it.

There are no comments yet — let us add three, to the first article and the second:

```sql
insert into comments (article_id, author, body) values
  (1, 'Asel', 'Good piece'),
  (1, 'Dana', 'Thanks'),
  (2, 'Asel', 'Will there be one on the crown?');
```

### `join`: only those with a pair

Collecting the comments together with the articles is the job of `join`:

```sql
select a.title, c.author, c.body
from articles a
join comments c on c.article_id = a.id
order by a.id, c.id;
```

```
┌────────────────────┬────────┬─────────────────────────────────┐
│       title        │ author │              body               │
├────────────────────┼────────┼─────────────────────────────────┤
│ About the steppe   │ Asel   │ Good piece                      │
│ About the steppe   │ Dana   │ Thanks                          │
│ What a shanyraq is │ Asel   │ Will there be one on the crown? │
└────────────────────┴────────┴─────────────────────────────────┘
```

`join` keeps only the pairs that matched. Articles with no comments will not appear at all — and a front page needs them, even with a zero.

### `left join`: and those without one

It takes every row of the left table and puts nothing on the right where there is no pair.

```sql
select a.title, count(c.id) as comments
from articles a
left join comments c on c.article_id = a.id
group by a.id, a.title
order by comments desc, a.id;
```

```
┌────────────────────┬──────────┐
│       title        │ comments │
├────────────────────┼──────────┤
│ About the steppe   │ 2        │
│ What a shanyraq is │ 1        │
│ Hello, world       │ 0        │
│ Draft              │ 0        │
└────────────────────┴──────────┘
```

`count(c.id)` counts the comments and `group by` gathers rows into groups per article. Note the zeros on the last two: a plain `join` would simply not have shown them.

### The SQLite trap: foreign keys are off

Delete an article and check whether its comments went with it:

```sql
delete from articles where slug = 'dala';
select count(*) as comments_left from comments;
```

```
┌───────────────┐
│ comments_left │
├───────────────┤
│ 1             │
└───────────────┘
```

There were three comments, two belonged to the deleted article, and one remains. `on delete cascade` did its job.

It did so because at the very beginning we wrote:

```sql
pragma foreign_keys = on;
```

Without that line, nothing written above about keys works at all. Check it on a fresh connection:

```sql
pragma foreign_keys;
```

```
┌──────────────┐
│ foreign_keys │
├──────────────┤
│ 0            │
└──────────────┘
```

Zero. For historical reasons SQLite keeps foreign-key checking **off by default**, and it has to be switched on again in every connection. Silently, with no error: `references` is recorded in the schema and simply not enforced, and `on delete cascade` leaves the comments orphaned.

It is the first thing you will do next lesson, when Go connects to the database.

## The lesson map

![Lesson map: two tables and the key that ties them together](/static/course/go/map-sql-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why do the comments live in a table of their own rather than as columns in `articles`?
2. How does `join` differ from `left join`, and when does the difference show?
3. What happens to `on delete cascade` if `pragma foreign_keys = on` is forgotten?

## Exercise

**Required.** Add a third table, `tags`, and connect it to the articles. One article can carry several tags and one tag can sit on several articles, so the connection needs a table of its own with two columns. Write a query that lists articles with their tags, and a second one that counts how many articles each tag has.

**Optional.**

- Insert a comment for an article number that does not exist and read what the database says. Then switch `pragma foreign_keys` off and try again.
- Replace the `left join` with a `join` in the counting query and explain where the two rows went.
- Delete an article that has no comments and confirm that the comment count does not change.

## Where this goes in your blog

The schema from these two lessons is the blog's schema. Next lesson Go connects to it: `database/sql`, a driver, and the first row read from a file rather than from a slice.

Still open. We typed the queries by hand in a console: from Go they have to be sent so that a stranger's text cannot become part of the command, which is the lesson on injections. And `pragma foreign_keys` has to be switched on in every connection, of which a program has several — the next lesson says so separately.

## Answers

1. Because the number of comments is not known in advance while a table has a fixed number of columns. Columns `comment1`, `comment2` run out at the third comment, and you cannot search or count by them. A separate table removes the limit: it holds as many rows as are needed.
2. `join` keeps only the pairs that matched on the key, while `left join` takes every row of the left table and puts nothing where there is no pair. The difference shows exactly when there are unpaired rows: an article with no comments does not appear in a `join` at all, and appears with a zero in a `left join`.
3. Nothing happens — and that is the worst outcome. There is no error, `references` stays recorded in the schema, and it simply stops being enforced: deleting an article leaves its comments pointing at a number that no longer exists. That is why `pragma foreign_keys = on` is switched on in every connection.

## Sources

- [SQLite: foreign keys](https://sqlite.org/foreignkeys.html)
- [SQLite: the query language](https://sqlite.org/lang.html)
- [SQLite: STRICT tables](https://sqlite.org/stricttables.html)
- [The PostgreSQL tutorial: SQL from scratch](https://www.postgresql.org/docs/current/tutorial-sql.html)
