"""The store: the long table to a CSV and back, and every write atomic.

A write that is not atomic is a file somebody can open halfway through. It costs
nothing to avoid: write beside the target, then rename. os.replace swaps the
name in one motion, so a reader sees either yesterday's file or today's and
never half of either.
"""

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
