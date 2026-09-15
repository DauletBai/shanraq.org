import sqlite3
from pathlib import Path

ROOT = Path(__file__).parent


def database():
    db = sqlite3.connect(":memory:")
    db.executescript((ROOT / "schema.sql").read_text())
    db.executescript((ROOT / "seed.sql").read_text())
    return db


def test_seed_and_totals():
    db = database()
    assert db.execute("SELECT count(*) FROM transactions").fetchone()[0] == 8
    income, expense = db.execute("""
        SELECT sum(amount_tiyn) FILTER (WHERE kind='income'),
               sum(amount_tiyn) FILTER (WHERE kind='expense')
        FROM transactions JOIN categories ON categories.id=category_id
    """).fetchone()
    assert (income, expense, income - expense) == (42000000, 12780000, 29220000)


def test_constraints_reject_bad_money_and_date():
    db = database()
    for sql in (
        "INSERT INTO transactions(account_id,category_id,happened_on,amount_tiyn) VALUES(1,2,'bad',100)",
        "INSERT INTO transactions(account_id,category_id,happened_on,amount_tiyn) VALUES(1,2,'2026-09-15',-1)",
        "INSERT INTO transactions(account_id,category_id,happened_on,amount_tiyn) VALUES(99,2,'2026-09-15',100)",
    ):
        try:
            db.execute(sql)
        except sqlite3.IntegrityError:
            pass
        else:
            raise AssertionError(f"constraint accepted: {sql}")


if __name__ == "__main__":
    test_seed_and_totals()
    test_constraints_reject_bad_money_and_date()
    print("budget project: OK")
