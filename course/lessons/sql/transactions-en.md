# Transactions and UPDATE: change everything or nothing

_Lead (summary):_ **Make safe corrections and treat a transfer as one atomic database operation.**

## Why this matters

A transfer needs at least two related records. If debit succeeds and credit fails, the database invents or loses money.

> **Image.** A transaction is a lock with two gates: either the boat passes completely or boat and water return to their initial state.

## See the whole thing first

```sql
BEGIN IMMEDIATE;

UPDATE transactions
SET note = 'Utilities and rent'
WHERE id = 2;

SELECT changes() AS changed_rows;
COMMIT;
```

## Explanation

ACID summarises atomicity, consistency, isolation, and durability; exact guarantees vary by DBMS. Before UPDATE or DELETE, run SELECT with the same WHERE and then check the affected-row count. COMMIT makes changes durable; ROLLBACK cancels unfinished work; SAVEPOINT rolls back part. A transfer between your accounts is not household spending.

## Lesson map

![Transactions and UPDATE](/static/course/sql/map-transactions-en.svg)

## Say it in your own words

1. What is atomicity?
2. How do you test WHERE before changing data?
3. Why is an internal transfer not spending?

## Exercise

Inside a transaction set Transport's limit to 40,000 tenge, read it, ROLLBACK, and confirm the old value returned.

## Where this fits in the project

Every related budget change is atomic and can be cancelled before commit.

## Answers

```sql
BEGIN;
UPDATE categories SET monthly_limit_tiyn=4000000 WHERE name='Transport';
SELECT monthly_limit_tiyn FROM categories WHERE name='Transport';
ROLLBACK;
SELECT monthly_limit_tiyn FROM categories WHERE name='Transport';
```

## Sources

- [SQLite: транзакции](https://sqlite.org/lang_transaction.html)
- [SQLite: изоляция](https://sqlite.org/isolation.html)
