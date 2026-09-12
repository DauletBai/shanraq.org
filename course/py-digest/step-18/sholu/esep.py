"""The report: the grouping, the join, the year before -- and the spread.

diff() inside a grouping answers "how much more than last year" without a loop
and without the off-by-one that a loop over sorted years invites. It needs the
rows in order, which is what the cleaning step guarantees.

Resampling and a rolling mean, the other half of lesson 31, are not here on
purpose: five yearly points are not a series to smooth. They arrive with a
monthly one.

ozara() is lesson 37: over all the years a country has, it counts the mean, the
median, the spread and the quarters, and picks the one measure that goes in the
heading. The choice is a rule and not a taste -- see headline().

indeks() is lesson 38: yearly rates are multipliers, so the growth over a run of
years is their product and never their sum. It is the first count here that uses
every year at once, which is why it is also the first to refuse a series with a
gap in it: an index over a broken run of years means nothing.
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


def headline(row):
    """Which measure describes a series: a rule, not a taste.

    When the mean has moved away from the median by more than a tenth of the
    median, the series has a tail -- three of these countries have spike years,
    and a mean over those describes the spike rather than the country.
    """
    mean, median = row.mean(), row.median()
    if median and abs(mean - median) / median > 0.1:
        return "медиана", round(median, 2)
    return "орташа", round(mean, 2)


def ozara(table):
    """Per country, over every year it has: the centre and the spread.

    Four numbers instead of one, because one hides what it costs: a mean of 9
    with a spread of 4 and a mean of 9 with a spread of 1 are different
    countries, and the report used to print them the same way.
    """
    assert not table.empty, "кестеде бірде-бір жол жоқ"
    rows = []
    for country, part in table.groupby("country", observed=True):
        values = part["value"].dropna()
        if values.empty:
            continue
        quarters = values.quantile([0.25, 0.75])
        measure, value = headline(values)
        rows.append({
            "country": str(country),
            "жылдар": int(values.count()),
            "орташа": round(values.mean(), 2),
            "медиана": round(values.median(), 2),
            "шашырау": round(values.std(), 2),
            "ширекаралық": round(quarters[0.75] - quarters[0.25], 2),
            "өлшем": measure,
            "мәні": value,
        })
    out = pd.DataFrame(rows)
    named = out.merge(anyqtama.table(), on="country", how="left", validate="one_to_one")
    named["аты"] = named["аты"].fillna(named["country"])
    return named[["country", "аты", "жылдар", "орташа", "медиана", "шашырау",
                  "ширекаралық", "өлшем", "мәні"]]


def indeks(table, base=100.0):
    """Per country: the price index, the multiplier and what a thousand became.

    The years have to run without a gap. A missing year is not a rounding
    problem -- the product would quietly skip a multiplier and the index would
    come out lower than the truth -- so it is reported instead of counted.
    """
    assert not table.empty, "кестеде бірде-бір жол жоқ"
    rows = []
    for country, part in table.sort_values("year").groupby("country", observed=True):
        years = part["year"].astype(int).tolist()
        values = part["value"].tolist()
        gap = [y for y in range(years[0], years[-1] + 1) if y not in years]
        factors = [1 + v / 100 for v in values[1:]]      # the first year is the base
        total = 1.0
        for f in factors:
            total *= f
        rows.append({
            "country": str(country),
            "база": years[0],
            "соңғы": years[-1],
            "үзіліс": ", ".join(str(y) for y in gap) if gap else "жоқ",
            "индекс": round(base * total, 2),
            "есе": round(total, 3),
            "мың теңге": round(1000 / total, 2),
        })
    out = pd.DataFrame(rows)
    named = out.merge(anyqtama.table(), on="country", how="left", validate="one_to_one")
    named["аты"] = named["аты"].fillna(named["country"])
    return named[["country", "аты", "база", "соңғы", "үзіліс", "индекс", "есе", "мың теңге"]]


def write(table, path):
    """Writes the report as a CSV a spreadsheet can open."""
    path.parent.mkdir(parents=True, exist_ok=True)
    report(table).to_csv(path, index=False)
