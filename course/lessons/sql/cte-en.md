# Subqueries and WITH: answer a complex question in steps

_Lead (summary):_ **Break the savings-rate calculation into named, testable CTE stages.**

## Why this matters

Copying one long formula makes review difficult. Named intermediate steps reveal where a wrong total entered the calculation.

> **Image.** CTEs are clear mixing bowls: prepare income and spending separately, inspect them, then combine them.

## See the whole thing first

```sql
WITH monthly AS (
    SELECT c.kind, SUM(t.amount_tiyn) AS amount_tiyn
    FROM transactions AS t
    JOIN categories AS c ON c.id = t.category_id
    WHERE t.happened_on >= '2026-09-01' AND t.happened_on < '2026-10-01'
    GROUP BY c.kind
), totals AS (
    SELECT COALESCE(SUM(amount_tiyn) FILTER (WHERE kind='income'), 0) AS income,
           COALESCE(SUM(amount_tiyn) FILTER (WHERE kind='expense'), 0) AS expense
    FROM monthly
)
SELECT income / 100.0 AS income_kzt,
       expense / 100.0 AS expense_kzt,
       ROUND(100.0 * (income - expense) / NULLIF(income, 0), 1) AS saving_rate
FROM totals;
```

## Explanation

A CTE exists for one statement; it is not a stored table. Read several CTEs top-down as calculation stages. NULLIF(income,0) prevents a false percentage when income is zero. Correlated subqueries are not automatically slow; choose between them, joins, and pre-aggregation using clarity and the query plan.

## Lesson map

![Subqueries and WITH](/static/course/sql/map-cte-en.svg)

## Say it in your own words

1. How does a CTE differ from a table?
2. Why use NULLIF here?
3. Why not declare every subquery slow?

## Exercise

Add a CTE that returns spending as a percentage of income and NULL when income is absent.

## Where this fits in the project

Named stages become the final report and let us test each calculation separately.

## Answers

```sql
WITH totals AS (
    SELECT SUM(amount_tiyn) FILTER (WHERE kind='income') AS income,
           SUM(amount_tiyn) FILTER (WHERE kind='expense') AS expense
    FROM transactions JOIN categories ON categories.id=category_id
), ratio AS (
    SELECT 100.0 * expense / NULLIF(income, 0) AS expense_rate FROM totals
)
SELECT expense_rate FROM ratio;
```

## Sources

- [SQLite: WITH](https://sqlite.org/lang_with.html)
