"""The report: the grouping and the join, over a table that has been cleaned.

Nothing here decides what to do with a gap any more -- tazalau.py did that
before the report was built. What is left is a marker instead of an empty cell:
a country the directory does not know gets its code for a name and the word
"unknown" for a region, because an empty cell in a report is a question nobody
answers.
"""

import pandas as pd

from sholu import anyqtama

UNKNOWN = "белгісіз"


def report(table):
    """Per country: the name, the years, the gaps, the average and the peak."""
    assert not table.empty, "кестеде бірде-бір жол жоқ"
    out = table.groupby("country", observed=True).agg(
        жылдар=("value", "size"),
        дерегі_жоқ=("value", lambda column: int(column.isna().sum())),
        орташа=("value", "mean"),
        ең_жоғары=("value", "max"),
    )
    # The label of the largest value in each group is the row it sits in, and
    # that row knows its year.
    peaks = table.loc[table.groupby("country", observed=True)["value"].idxmax()]
    out["ең_жоғары_жыл"] = peaks.set_index("country")["year"]
    out = out.reset_index()
    out["country"] = out["country"].astype("str")

    named = out.merge(anyqtama.table(), on="country", how="left", validate="one_to_one")
    named["аты"] = named["аты"].fillna(named["country"])
    # A marker, not a guess: "unknown" says what is true, an empty cell says
    # nothing and a made-up region would say something false.
    named["өңір"] = named["өңір"].fillna(UNKNOWN)
    columns = ["country", "аты", "өңір", "жылдар", "дерегі_жоқ", "орташа", "ең_жоғары", "ең_жоғары_жыл"]
    return named[columns].round({"орташа": 2, "ең_жоғары": 2})


def write(table, path):
    """Writes the report as a CSV a spreadsheet can open."""
    path.parent.mkdir(parents=True, exist_ok=True)
    report(table).to_csv(path, index=False)
