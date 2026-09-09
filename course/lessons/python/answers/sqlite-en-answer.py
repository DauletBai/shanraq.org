"""The answer to lesson 21's exercise, English: a history that survives a second run.

The point of the table is its key: day and currency together. Running the same
day twice updates one row instead of adding a second, so the file grows by what
is new and by nothing else. Every value goes in through a placeholder, and the
rows come back out as a record with names rather than as a tuple to count on
one's fingers.
"""

import sqlite3
from dataclasses import dataclass
from pathlib import Path

HERE = Path(__file__).parent
FILE = HERE / "istoriya.db"
FILE.unlink(missing_ok=True)


@dataclass
class Rate:
    """One row of the history: the day, the currency, the rate."""

    day: str
    code: str
    value: float


rows = [
    Rate("2026-01-15", "USD", 510.43),
    Rate("2026-01-15", "EUR", 594.86),
    Rate("2026-01-16", "USD", 511.02),
    Rate("2026-01-16", "EUR", 596.10),
]

conn = sqlite3.connect(FILE)
conn.row_factory = sqlite3.Row
conn.execute("""
    CREATE TABLE IF NOT EXISTS rates (
        day   TEXT NOT NULL,
        code  TEXT NOT NULL,
        value REAL NOT NULL,
        PRIMARY KEY (day, code)
    )
""")


def save(records):
    """Puts records into the database: new ones added, known ones updated."""
    conn.executemany("""
        INSERT INTO rates (day, code, value) VALUES (?, ?, ?)
        ON CONFLICT (day, code) DO UPDATE SET value = excluded.value
    """, [(r.day, r.code, r.value) for r in records])
    conn.commit()


save(rows)
print("after the first run:", conn.execute("SELECT count(*) FROM rates").fetchone()[0])
save(rows)
print("after the second run:", conn.execute("SELECT count(*) FROM rates").fetchone()[0])

print("the dollar by day:")
for row in conn.execute("SELECT day, value FROM rates WHERE code = ? ORDER BY day", ("USD",)):
    print(f"  {row['day']}: {row['value']:.2f}")

average = conn.execute("SELECT avg(value) FROM rates WHERE code = ?", ("USD",)).fetchone()[0]
print(f"the dollar's average rate: {average:.2f}")

last = conn.execute(
    "SELECT day, code, value FROM rates ORDER BY day DESC, code LIMIT 1").fetchone()
newest = Rate(last["day"], last["code"], last["value"])
print("the last record:", newest)

conn.close()
FILE.unlink()
