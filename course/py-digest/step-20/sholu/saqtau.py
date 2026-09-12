"""The store: the long table to a CSV and back, and every write atomic.

A write that is not atomic is a file somebody can open halfway through. It costs
nothing to avoid: write beside the target, then rename. os.replace swaps the
name in one motion, so a reader sees either yesterday's file or today's and
never half of either.
"""

import json
import os

import pandas as pd


def _replace(text, path):
    """Writes text through a temporary file next to the target."""
    path.parent.mkdir(parents=True, exist_ok=True)
    temporary = path.with_suffix(path.suffix + ".tmp")
    temporary.write_text(text, encoding="utf-8")
    os.replace(temporary, path)
    return path


def save(table, path):
    """Writes the table as country,year,value. A missing value is an empty cell."""
    return _replace(table.to_csv(index=False, float_format="%.2f"), path)


def load(path):
    """Reads the saved table; an absent file means there is nothing yet."""
    if not path.exists():
        return None
    return pd.read_csv(path)


def write_text(text, path):
    """The same atomic write for anything else the digest produces."""
    return _replace(text, path)


def save_json(data, path):
    """The same atomic write for what is not a table: the indicator notes."""
    return _replace(json.dumps(data, ensure_ascii=False, indent=2) + "\n", path)


def load_json(path):
    """Reads them back; an absent or broken file means there is nothing yet.

    A definition is a nicety, not the data, so a file somebody edited by hand
    into invalid json must not take the whole report down with it.
    """
    if not path.exists():
        return None
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except (ValueError, OSError):
        return None
