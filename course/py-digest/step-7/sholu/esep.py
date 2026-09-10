"""The report: the series becomes a table, and the table answers.

Until now the series was a dictionary, and every question about it was a loop of
its own: one to skip the years without a figure, one to add up, one to find the
largest. Here it is a DataFrame. The years are the labels of the rows, so the
year of the peak comes back as a year rather than as a position, and the gaps
count themselves.
"""

import csv

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
    # idxmax answers with the label, which here is the year itself. Over a
    # dictionary this took a loop that carried the year along with the maximum.
    peak = value.idxmax()
    rows = [
        ("жылдар", f"{len(value)}"),
        ("дерегі жоқ", f"{int(value.isna().sum())}"),
        # mean() leaves the gaps out on its own -- the same rule the hand-written
        # average() had to be told.
        ("орташа", f"{value.mean():.2f}"),
        ("ең жоғары", f"{value.max():.2f}"),
        ("ең жоғары жыл", f"{peak}"),
    ]
    return pd.DataFrame(rows, columns=["metric", "value"])


def write(table, path):
    """Writes the report as a CSV a spreadsheet can open."""
    path.parent.mkdir(parents=True, exist_ok=True)
    out = report(table)
    with path.open("w", encoding="utf-8", newline="") as target:
        writer = csv.writer(target)
        writer.writerow(out.columns.tolist())
        for _, row in out.iterrows():
            writer.writerow([row["metric"], row["value"]])
