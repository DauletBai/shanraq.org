"""The store: the table to a CSV and back, one line each way.

The csv module wrote the file row by row and read it back row by row, turning
every value from text by hand. read_csv and to_csv do the same work with the
column names, the missing values and the decimals as arguments -- and what comes
back is the table itself rather than a dictionary to rebuild one from.
"""

import pandas as pd

MISSING = ""


def save(table, path):
    """Writes the table as year,value. A missing year is an empty cell."""
    path.parent.mkdir(parents=True, exist_ok=True)
    table.to_csv(path, float_format="%.2f", na_rep=MISSING)


def load(path):
    """Reads the saved table; an absent file means there is nothing yet."""
    if not path.exists():
        return None
    # index_col is what keeps the year a label rather than becoming a column of
    # its own -- the round trip through a file forgets what the index was.
    return pd.read_csv(path, index_col="year")
