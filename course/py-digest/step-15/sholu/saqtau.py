"""The store: the long table to a CSV and back, one line each way.

The table is long now -- country, year, value -- so the row labels are just
numbers and have no business in the file. index=False on the way out, nothing to
ask for on the way in.
"""

import pandas as pd


def save(table, path):
    """Writes the table as country,year,value. A missing value is an empty cell."""
    path.parent.mkdir(parents=True, exist_ok=True)
    table.to_csv(path, index=False, float_format="%.2f")


def load(path):
    """Reads the saved table; an absent file means there is nothing yet."""
    if not path.exists():
        return None
    return pd.read_csv(path)
