# SELECT: find the needed rows without changing them

_Lead (summary):_ **Learn projection, filtering, stable sorting, and limiting a result.**

## Why this matters

A database is useful because it returns exactly what is needed: recent spending, purchases above a threshold, or one month's events.

> **Image.** FROM puts the receipt box on the table; WHERE keeps matching receipts; SELECT covers unused fields; ORDER BY arranges them; LIMIT takes the first few.

## See the whole thing first

```sql
SELECT happened_on, amount_tiyn / 100.0 AS amount_kzt, note
FROM transactions
WHERE amount_tiyn >= 1000000
ORDER BY happened_on DESC, id DESC
LIMIT 5;
```

## Explanation

The logical order begins with FROM and WHERE even though SELECT is written first. SELECT * is fine for exploration, but application contracts should name columns. Without ORDER BY row order is undefined; add id as a tie-breaker. LIKE rules, especially case handling, differ across systems.

## Lesson map

![SELECT](/static/course/sql/map-select-en.svg)

## Say it in your own words

1. Why is SELECT written first but not evaluated first?
2. What resolves equal dates?
3. When is SELECT * risky?

## Exercise

Return September's three largest transactions with date, tenge amount, and note; break ties by newest row.

## Where this fits in the project

This query pattern powers transaction history and previews rows before a bulk change.

## Answers

```sql
SELECT happened_on, amount_tiyn / 100.0 AS amount_kzt, note
FROM transactions
WHERE happened_on >= '2026-09-01' AND happened_on < '2026-10-01'
ORDER BY amount_tiyn DESC, happened_on DESC, id DESC
LIMIT 3;
```

## Sources

- [SQLite: SELECT](https://sqlite.org/lang_select.html)
