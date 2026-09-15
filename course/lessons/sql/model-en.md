# Tables, rows, columns, and NULL: model the budget first

_Lead (summary):_ **We identify the project's accounts, categories, transactions, and limits, and separate a missing value from zero. Before learning syntax, understand the single fact represented by each row.**

## Why this matters

A poor model forces every query to guess. A person may decipher `Food / card / 49800` in one cell; a database cannot know which part is the category, account, or amount.

> **Image.** One row is one complete sentence: “Transaction 7 happened on 10 September, on the Card account, in Transport, for 6,800 tenge.” A row that tells several unrelated stories is difficult to verify.

## See the whole thing first

```sql
SELECT id, happened_on, amount_tiyn, note
FROM transactions
ORDER BY happened_on, id;
```

A table contains rows of one kind; a column gives a value its role. Primary key `id` names a row reliably. Foreign keys `account_id` and `category_id` connect it to other tables.

## Explanation

`accounts` are places where money is kept. `categories` explain income or spending and may define a limit. `transactions` are events that happened. A monthly report is a computed result, not another copy of those facts.

`NULL` means unknown or not applicable. A salary limit is `NULL` because an expense limit does not apply; zero would mean spending is forbidden. Use `IS NULL`, never `= NULL`:

```sql
SELECT name
FROM categories
WHERE monthly_limit_tiyn IS NULL;
```

SQL uses three-valued logic: `TRUE`, `FALSE`, and `UNKNOWN`. That is why our date constraint also checks `date(...) IS NOT NULL`.

## Lesson map

![One row is one fact; NULL is not zero](/static/course/sql/map-model-en.svg)

## Say it in your own words

1. What fact does one `transactions` row store?
2. How does `NULL` differ from zero and an empty string?
3. Why does a row need a stable key?

## Exercise

Return only `name` and `kind` for categories without a monthly limit.

## Where this fits in the project

The model becomes the database schema. Reports are calculated from source transactions, so two copies of one total never need synchronising.

## Answers

```sql
SELECT name, kind
FROM categories
WHERE monthly_limit_tiyn IS NULL;
```

## Sources

- [SQLite expressions and NULL](https://sqlite.org/lang_expr.html)
- [SQLite foreign keys](https://sqlite.org/foreignkeys.html)
