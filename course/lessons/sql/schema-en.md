# CREATE TABLE: a schema that rejects impossible data

_Lead (summary):_ **Create the family-budget tables and move essential validation into the database.**

## Why this matters

An interface checks only one input path. Imports and scripts can bypass it; database constraints protect every path.

> **Image.** A schema is an ice tray: water may arrive from any jug, but the mould will not produce an impossible cube.

## See the whole thing first

```sql
PRAGMA foreign_keys = ON;

CREATE TABLE transactions (
    id INTEGER PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(id),
    category_id INTEGER NOT NULL REFERENCES categories(id),
    happened_on TEXT NOT NULL CHECK (
        date(happened_on) IS NOT NULL AND happened_on = date(happened_on)
    ),
    amount_tiyn INTEGER NOT NULL CHECK (amount_tiyn > 0),
    note TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Explanation

NOT NULL requires a value; UNIQUE prevents duplicate account names; CHECK validates a row; REFERENCES rejects a missing account or category. SQLite foreign keys must be enabled for each connection with PRAGMA foreign_keys=ON. INTEGER PRIMARY KEY already generates a row identifier; AUTOINCREMENT is rarely needed. SQLite stores our ISO date as TEXT because it has no dedicated DATE type.

## Lesson map

![CREATE TABLE](/static/course/sql/map-schema-en.svg)

## Say it in your own words

1. Why are NOT NULL and CHECK different?
2. What happens when account_id refers to no account?
3. Why does the date constraint make two checks?

## Exercise

Create goals with a key, unique non-empty name, positive target_tiyn, and an optional valid ISO due date.

## Where this fits in the project

The complete contract lives in course/sql-budget/schema.sql; later queries rely on it.

## Answers

```sql
CREATE TABLE goals (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE CHECK (trim(name) <> ''),
    target_tiyn INTEGER NOT NULL CHECK (target_tiyn > 0),
    due_on TEXT CHECK (due_on IS NULL OR (
        date(due_on) IS NOT NULL AND due_on = date(due_on)
    ))
);
```

## Sources

- [SQLite: CREATE TABLE](https://sqlite.org/lang_createtable.html)
- [SQLite: AUTOINCREMENT](https://sqlite.org/autoinc.html)
