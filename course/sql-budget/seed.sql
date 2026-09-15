INSERT INTO accounts (id, name, opening_balance_tiyn) VALUES
    (1, 'Card', 18000000),
    (2, 'Cash', 2500000);

INSERT INTO categories (id, name, kind, monthly_limit_tiyn) VALUES
    (1, 'Salary', 'income', NULL),
    (2, 'Food', 'expense', 9000000),
    (3, 'Home', 'expense', 6500000),
    (4, 'Transport', 'expense', 3500000),
    (5, 'Learning', 'expense', 2500000);

INSERT INTO transactions
    (id, account_id, category_id, happened_on, amount_tiyn, note)
VALUES
    (1, 1, 1, '2026-09-01', 42000000, 'September salary'),
    (2, 1, 3, '2026-09-02', 5500000, 'Utilities and rent'),
    (3, 1, 2, '2026-09-03', 1865000, 'Groceries'),
    (4, 2, 4, '2026-09-04', 420000, 'Bus card'),
    (5, 1, 5, '2026-09-06', 1200000, 'Books'),
    (6, 2, 2, '2026-09-08', 975000, 'Market'),
    (7, 1, 4, '2026-09-10', 680000, 'Fuel'),
    (8, 1, 2, '2026-09-12', 2140000, 'Groceries');

