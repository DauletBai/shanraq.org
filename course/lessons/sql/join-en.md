# Keys and JOIN: connect transactions, accounts, and categories

_Lead (summary):_ **Normalise the budget model and join tables by stable keys without multiplying rows.**

## Why this matters

Repeating a category name in every receipt creates spelling variants and conflicting rules. Store the meaning once and reference its key.

> **Image.** A transaction holds a cloakroom token; JOIN retrieves the exact coat by token number, not by a similar colour.

## See the whole thing first

```sql
SELECT t.happened_on, a.name AS account, c.name AS category,
       t.amount_tiyn / 100.0 AS amount_kzt
FROM transactions AS t
JOIN accounts AS a ON a.id = t.account_id
JOIN categories AS c ON c.id = t.category_id
ORDER BY t.happened_on, t.id;
```

## Explanation

INNER JOIN keeps matches; LEFT JOIN preserves every left row and uses NULL for missing right data. Join on keys, not similar text. Before trusting totals, verify key uniqueness and compare row counts: duplicate lookup keys can multiply facts.

## Lesson map

![Keys and JOIN](/static/course/sql/map-join-en.svg)

## Say it in your own words

1. How do INNER and LEFT JOIN differ?
2. Why not join names?
3. How can a join multiply money?

## Exercise

List every account and its transaction count, including accounts with none.

## Where this fits in the project

Joins add readable names to compact facts while foreign keys protect relationships.

## Answers

```sql
SELECT c.name, COALESCE(SUM(t.amount_tiyn), 0) / 100.0 AS spent_kzt
FROM categories AS c
LEFT JOIN transactions AS t
  ON t.category_id = c.id
 AND t.happened_on >= '2026-09-01'
 AND t.happened_on <  '2026-10-01'
WHERE c.kind = 'expense'
GROUP BY c.id, c.name
ORDER BY c.name;
```

## Sources

- [SQLite: JOIN в SELECT](https://sqlite.org/lang_select.html)
- [PostgreSQL: соединения таблиц](https://www.postgresql.org/docs/current/queries-table-expressions.html)
