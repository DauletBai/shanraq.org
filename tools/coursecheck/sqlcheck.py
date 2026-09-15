#!/usr/bin/env python3
"""Execute every SQL fence from the SQL course against its sample database."""

from __future__ import annotations

import argparse
import re
import sqlite3
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
SCHEMA = (ROOT / "course/sql-budget/schema.sql").read_text(encoding="utf-8")
SEED = (ROOT / "course/sql-budget/seed.sql").read_text(encoding="utf-8")
FENCE = re.compile(r"```sql\s*\n(.*?)```", re.DOTALL | re.IGNORECASE)


def database() -> sqlite3.Connection:
    db = sqlite3.connect(":memory:")
    db.executescript(SCHEMA)
    db.executescript(SEED)
    return db


def runnable(sql: str) -> str:
    # The reference schema already contains the course index. IF NOT EXISTS
    # keeps its later teaching example executable without weakening that file.
    return re.sub(
        r"\bCREATE\s+INDEX\s+(?!IF\s+NOT\s+EXISTS)",
        "CREATE INDEX IF NOT EXISTS ",
        sql,
        flags=re.IGNORECASE,
    )


def check(path: Path) -> list[str]:
    errors: list[str] = []
    text = path.read_text(encoding="utf-8")
    blocks = FENCE.findall(text)
    if not blocks:
        return [f"{path}: no SQL fences"]
    for number, source in enumerate(blocks, 1):
        if not sqlite3.complete_statement(source.strip()):
            errors.append(f"{path}: SQL block {number} is incomplete")
            continue
        if re.search(r"\bCREATE\s+TABLE\s+transactions\b", source, re.IGNORECASE):
            db = sqlite3.connect(":memory:")
        else:
            db = database()
        try:
            db.executescript(runnable(source))
        except sqlite3.Error as exc:
            errors.append(f"{path}: SQL block {number}: {exc}")
        finally:
            db.close()
    return errors


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("paths", nargs="+", type=Path)
    args = parser.parse_args()
    errors = [error for path in args.paths for error in check(path)]
    if errors:
        print("\n".join(errors), file=sys.stderr)
        return 1
    print(f"SQL blocks execute: {len(args.paths)} files")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
