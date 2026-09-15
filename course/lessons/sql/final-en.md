# Indexes, EXPLAIN, and the final report

_Lead (summary):_ **Assemble the report, inspect its plan, add a justified index, and prove a backup by restoring it.**

## Why this matters

A correct query may become slow; a fast query may be wrong; an untested backup may be unusable.

> **Image.** An index is a book index: faster lookup costs space and update work. A backup is a spare key tested in the real door.

## See the whole thing first

```sql
EXPLAIN QUERY PLAN
SELECT category_id, SUM(amount_tiyn)
FROM transactions
WHERE happened_on >= '2026-09-01' AND happened_on < '2026-10-01'
GROUP BY category_id;

CREATE INDEX transactions_month_category
    ON transactions(happened_on, category_id);
```

## Explanation

Match an index to a measured query path; do not index every column or judge using eight rows. EXPLAIN QUERY PLAN output is diagnostic and version-dependent. SCAN is not always bad for a small table. A VIEW stores a query, not its result. Back up an open SQLite database with .backup or the backup API, then open the copy, run integrity_check, and compare control totals. Never commit real financial data.

## Lesson map

![Indexes, EXPLAIN, and the final report](/static/course/sql/map-final-en.svg)

## Say it in your own words

1. What does an index cost?
2. Why is SCAN not always bad?
3. When is a backup verified?

## Exercise

Create a backup, open it separately, run integrity_check and the monthly report, and compare all three totals.

## Where this fits in the project

The finished project has protected facts, explainable calculations, a measured index, and a recoverable copy.

## Answers

```text
sqlite3 budget.db ".read course/sql-budget/report.sql"
```

```sql
CREATE VIEW expense_by_category AS
SELECT category_id, SUM(amount_tiyn) AS amount_tiyn
FROM transactions
JOIN categories ON categories.id=category_id
WHERE categories.kind='expense'
GROUP BY category_id;
```

```text
sqlite3 budget.db ".backup budget-backup.db"
sqlite3 budget-backup.db "PRAGMA integrity_check;"
sqlite3 budget-backup.db ".read course/sql-budget/report.sql"
```

## Sources

- [SQLite: план запроса](https://sqlite.org/eqp.html)
- [SQLite: планировщик запросов](https://sqlite.org/queryplanner.html)
- [SQLite: резервное копирование](https://sqlite.org/backup.html)
