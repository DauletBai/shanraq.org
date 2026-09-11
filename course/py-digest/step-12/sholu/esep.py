"""The report: the grouping, the join, and now the year before.

diff() inside a grouping answers "how much more than last year" without a loop
and without the off-by-one that a loop over sorted years invites. It needs the
rows in order, which is what the cleaning step guarantees.

Resampling and a rolling mean, the other half of lesson 31, are not here on
purpose: five yearly points are not a series to smooth. They arrive with a
monthly one.
"""

import pandas as pd

from sholu import anyqtama

UNKNOWN = "белгісіз"


def report(table):
    """Per country: the name, the years, the gaps, the average, the peak and the
    change over the year before the last one."""
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

    # The last year of each country, and how it differs from the year before.
    ordered = table.sort_values(["country", "date"]).copy()
    ordered["айырма"] = ordered.groupby("country", observed=True)["value"].diff()
    last = ordered.groupby("country", observed=True).tail(1).set_index("country")
    out["соңғы_жыл"] = last["year"]
    out["соңғы_мән"] = last["value"]
    out["өткен_жылға"] = last["айырма"]

    out = out.reset_index()
    out["country"] = out["country"].astype("str")

    named = out.merge(anyqtama.table(), on="country", how="left", validate="one_to_one")
    named["аты"] = named["аты"].fillna(named["country"])
    # A marker, not a guess: "unknown" says what is true, an empty cell says
    # nothing and a made-up region would say something false.
    named["өңір"] = named["өңір"].fillna(UNKNOWN)
    columns = ["country", "аты", "өңір", "жылдар", "дерегі_жоқ", "орташа", "ең_жоғары",
               "ең_жоғары_жыл", "соңғы_жыл", "соңғы_мән", "өткен_жылға"]
    return named[columns].round({"орташа": 2, "ең_жоғары": 2, "соңғы_мән": 2, "өткен_жылға": 2})


def write(table, path):
    """Writes the report as a CSV a spreadsheet can open."""
    path.parent.mkdir(parents=True, exist_ok=True)
    report(table).to_csv(path, index=False)
