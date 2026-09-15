# CASE, COALESCE, and dates: turn stored facts into a useful view

_Lead (summary):_ **Compute display columns, handle NULL honestly, and define month boundaries.**

## Why this matters

The database stores stable facts—tiyn, ISO dates, and keys—while a reader needs tenge, labels, and a limit status.

> **Image.** A stored fact is an ingredient; a SELECT expression plates it. Do not store an ingredient pre-cut for every possible dish.

## See the whole thing first

```sql
SELECT happened_on,
       amount_tiyn / 100.0 AS amount_kzt,
       COALESCE(note, 'No note') AS note,
       CASE
           WHEN amount_tiyn >= 5000000 THEN 'large'
           WHEN amount_tiyn >= 1000000 THEN 'medium'
           ELSE 'small'
       END AS size
FROM transactions;
```

## Explanation

AS names a result. Dividing by 100.0 produces a display value while the stored money stays integral. COALESCE returns the first non-NULL value without rewriting the fact. CASE checks conditions top-down. Pass period boundaries as parameters; SQLite date modifiers are useful but not portable syntax.

## Lesson map

![CASE, COALESCE, and dates](/static/course/sql/map-expressions-en.svg)

## Say it in your own words

1. Why keep tiyn after displaying tenge?
2. What does COALESCE preserve?
3. Why does CASE order matter?

## Exercise

Add a label showing whether each expense is within its category limit.

## Where this fits in the project

These expressions make the report readable without duplicating source facts.

## Answers

```sql
SELECT id,
       CASE WHEN NULLIF(trim(note), '') IS NULL THEN 'no' ELSE 'yes' END AS has_note
FROM transactions;
```

## Sources

- [SQLite: выражения CASE и COALESCE](https://sqlite.org/lang_expr.html)
- [SQLite: функции даты и времени](https://sqlite.org/lang_datefunc.html)
