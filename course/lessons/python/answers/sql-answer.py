"""The answer to lesson 22's exercise: the report the database writes itself.

Nothing here is counted in Python. The days, the extremes and the averages come
back from one query each; the names come from a second table through a join; the
month is a group made out of the date's first seven characters. And the currency
nobody quoted answers None rather than zero, because no rows is not the same as
a rate of nothing.
"""

import sqlite3
from pathlib import Path

HERE = Path(__file__).parent
FILE = HERE / "otchet.db"
FILE.unlink(missing_ok=True)

conn = sqlite3.connect(FILE)
conn.row_factory = sqlite3.Row
conn.executescript("""
    CREATE TABLE rates (
        day   TEXT NOT NULL,
        code  TEXT NOT NULL,
        value REAL NOT NULL,
        PRIMARY KEY (day, code)
    );
    CREATE TABLE currencies (
        code  TEXT PRIMARY KEY,
        name  TEXT NOT NULL,
        quant INTEGER NOT NULL
    );
""")
conn.executemany("INSERT INTO rates VALUES (?, ?, ?)", [
    ("2026-01-15", "USD", 510.43), ("2026-01-15", "UZS", 4.24),
    ("2026-01-16", "USD", 511.02), ("2026-01-16", "UZS", 4.21),
    ("2026-02-02", "USD", 515.70), ("2026-02-02", "UZS", 4.19),
])
conn.executemany("INSERT INTO currencies VALUES (?, ?, ?)", [
    ("USD", "доллар США", 1), ("UZS", "узбекский сум", 100),
])
conn.commit()

print("== по валютам")
for row in conn.execute("""
    SELECT code, count(*) AS days, min(value) AS low, max(value) AS high,
           round(avg(value), 2) AS mid
    FROM rates
    GROUP BY code
    ORDER BY code
"""):
    print(f"{row['code']}: {row['days']} дн., от {row['low']} до {row['high']}, среднее {row['mid']}")

print("== за один день, с названиями")
for row in conn.execute("""
    SELECT c.name, round(r.value / c.quant, 4) AS unit
    FROM rates AS r
    JOIN currencies AS c ON c.code = r.code
    WHERE r.day = ?
    ORDER BY unit DESC
""", ("2026-01-16",)):
    print(f"{row['name']}: {row['unit']} за единицу")

print("== доллар по месяцам")
for row in conn.execute("""
    SELECT substr(day, 1, 7) AS month, round(avg(value), 2) AS mid
    FROM rates
    WHERE code = ?
    GROUP BY month
    ORDER BY month
""", ("USD",)):
    print(f"{row['month']}: {row['mid']}")

mid = conn.execute("SELECT avg(value) AS mid FROM rates WHERE code = ?", ("GBP",)).fetchone()["mid"]
print("фунт:", "данных нет" if mid is None else round(mid, 2))

conn.close()
FILE.unlink()
