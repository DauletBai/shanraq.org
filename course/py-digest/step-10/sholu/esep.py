"""The report: the grouping of step 9, now with names beside the codes.

The join is a left one and it is checked. Left, because the data decides which
rows exist and the directory only adds words to them: a country missing from the
directory has to stay in the report under its code rather than disappear from
it. Checked, because a duplicated code in the directory would quietly multiply
the rows it matched -- validate says so instead.
"""

import pandas as pd

from sholu import anyqtama


def report(table):
    """Per country: the name, the years, the gaps, the average and the peak."""
    assert not table.empty, "кестеде бірде-бір жол жоқ"
    out = table.groupby("country").agg(
        жылдар=("value", "size"),
        дерегі_жоқ=("value", lambda column: int(column.isna().sum())),
        орташа=("value", "mean"),
        ең_жоғары=("value", "max"),
    )
    # The label of the largest value in each group is the row it sits in, and
    # that row knows its year.
    peaks = table.loc[table.groupby("country")["value"].idxmax()]
    out["ең_жоғары_жыл"] = peaks.set_index("country")["year"]
    out = out.reset_index()

    named = out.merge(anyqtama.table(), on="country", how="left", validate="one_to_one")
    # A country the directory does not know keeps its code as its name: an empty
    # cell in a report is a question, and the code at least answers it.
    named["аты"] = named["аты"].fillna(named["country"])
    columns = ["country", "аты", "өңір", "жылдар", "дерегі_жоқ", "орташа", "ең_жоғары", "ең_жоғары_жыл"]
    return named[columns].round({"орташа": 2, "ең_жоғары": 2})


def write(table, path):
    """Writes the report as a CSV a spreadsheet can open."""
    path.parent.mkdir(parents=True, exist_ok=True)
    report(table).to_csv(path, index=False)
