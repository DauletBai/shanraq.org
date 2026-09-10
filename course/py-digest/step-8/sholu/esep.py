"""The report: the table answers, and pandas writes the answer out.

The five numbers are the same as in step 7. What changed is the last line: the
report is a table too, so it is written by to_csv rather than by a writer and a
loop over rows.
"""

import pandas as pd


def frame(series):
    """A table out of {year: value}: one column, the years as its labels."""
    table = pd.DataFrame({"value": pd.Series(series, dtype="float64")})
    table.index.name = "year"
    return table.sort_index()


def report(table):
    """The five numbers the report is made of, as a table of its own."""
    value = table["value"]
    assert value.notna().any(), "қатарда бірде-бір сан жоқ"
    # idxmax answers with the label, which here is the year itself.
    peak = value.idxmax()
    rows = [
        ("жылдар", f"{len(value)}"),
        ("дерегі жоқ", f"{int(value.isna().sum())}"),
        # mean() leaves the gaps out on its own.
        ("орташа", f"{value.mean():.2f}"),
        ("ең жоғары", f"{value.max():.2f}"),
        ("ең жоғары жыл", f"{peak}"),
    ]
    return pd.DataFrame(rows, columns=["metric", "value"])


def write(table, path):
    """Writes the report as a CSV a spreadsheet can open."""
    path.parent.mkdir(parents=True, exist_ok=True)
    # index=False: the row numbers are ours, and a file that carries them grows
    # a nameless first column the next reader has to ask about.
    report(table).to_csv(path, index=False)
