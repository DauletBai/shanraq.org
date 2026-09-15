PRAGMA foreign_keys = ON;

CREATE TABLE accounts (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    opening_balance_tiyn INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE categories (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    kind TEXT NOT NULL CHECK (kind IN ('income', 'expense')),
    monthly_limit_tiyn INTEGER CHECK (monthly_limit_tiyn IS NULL OR monthly_limit_tiyn >= 0)
);

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

CREATE INDEX transactions_month_category
    ON transactions(happened_on, category_id);
