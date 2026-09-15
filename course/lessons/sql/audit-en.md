# Data reconciliation: find a gap before making a decision

_Lead (summary):_ **Use control totals, duplicate candidates, and relationship checks to test completeness.**

## Why this matters

A query without an error does not prove complete data. One receipt imported twice creates a plausible but wrong total.

> **Image.** Reconciliation is a scale at the warehouse exit: individual records look valid, but their total weight must agree.

## See the whole thing first

```sql
SELECT happened_on, account_id, category_id, amount_tiyn, COUNT(*) AS copies
FROM transactions
GROUP BY happened_on, account_id, category_id, amount_tiyn, COALESCE(note, '')
HAVING COUNT(*) > 1;
```

## Explanation

Duplicate identity is a business rule, not simply equal ids; equal dates and amounts may still be two purchases. foreign_key_check finds broken links and integrity_check checks file structure, but neither knows whether spending is sensible. A row CHECK cannot enforce a rule across a whole month. Opening balances plus income minus spending must reconcile with calculated balances.

## Lesson map

![Data reconciliation](/static/course/sql/map-audit-en.svg)

## Say it in your own words

1. Why does a match only identify a duplicate candidate?
2. How does integrity_check differ from a business check?
3. What equation reconciles the budget?

## Exercise

Record transaction count and expense total, insert a copied transaction with a new id, and show which checks detect it.

## Where this fits in the project

Run reconciliation after imports and before publishing a report.

## Answers

```sql
SELECT COUNT(*) AS rows_total FROM transactions;
SELECT SUM(amount_tiyn) AS expense_tiyn
FROM transactions JOIN categories ON categories.id=category_id
WHERE kind='expense';
```

## Sources

- [SQLite: PRAGMA integrity_check и foreign_key_check](https://sqlite.org/pragma.html)
