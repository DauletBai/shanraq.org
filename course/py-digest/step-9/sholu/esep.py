"""The report: one grouping instead of a loop per country.

Until now the digest looked at one country and counted it with three calls. With
three countries a loop would have appeared -- the same three calls inside it, and
a dictionary to collect the results. groupby is that loop, written as a question:
split by country, count each piece, put the pieces back together as a table.
"""

import pandas as pd


def report(table):
    """Per country: the years, the gaps, the average, the peak and its year."""
    assert not table.empty, "кестеде бірде-бір жол жоқ"
    out = table.groupby("country").agg(
        жылдар=("value", "size"),
        дерегі_жоқ=("value", lambda column: int(column.isna().sum())),
        орташа=("value", "mean"),
        ең_жоғары=("value", "max"),
    )
    # The label of the largest value in each group is the row it sits in, and
    # that row knows its year -- so the peak year costs one lookup, not a loop.
    peaks = table.loc[table.groupby("country")["value"].idxmax()]
    out["ең_жоғары_жыл"] = peaks.set_index("country")["year"]
    return out.round({"орташа": 2, "ең_жоғары": 2}).reset_index()


def write(table, path):
    """Writes the report as a CSV a spreadsheet can open."""
    path.parent.mkdir(parents=True, exist_ok=True)
    report(table).to_csv(path, index=False)
