"""The source: a series from the World Bank.

The network lives in this file only. The rest of the code has no business
knowing there is an internet.

One request now brings several countries: the bank takes them separated by
semicolons, and asking three times for what can be asked once is three times the
waiting and three times the chance of being turned away.

The indicator is an argument, not a constant: since lesson 40 the digest reads
two of them -- prices and broad money -- through the same function, and the only
difference between the two series is the code in the address.

Since lesson 41 the same source is also asked what its numbers mean. A rate of
inflation is an average over a year in one publication and December against
December in another, and a report that prints the figure without saying which
one it is has told the reader less than it thinks.
"""

import json
import urllib.request

import pandas as pd

API = "https://api.worldbank.org/v2/country/{}/indicator/{}?format=json&per_page=400&date={}:{}"
ABOUT = "https://api.worldbank.org/v2/indicator/{}?format=json"
INFLATION = "FP.CPI.TOTL.ZG"
MONEY = "FM.LBL.BMNY.GD.ZS"      # broad money as a share of GDP, %


def fetch(countries, indicator, first, last):
    """A long table: one row per country and year, gaps included."""
    url = API.format(";".join(countries), indicator, first, last)
    with urllib.request.urlopen(url, timeout=30) as answer:
        data = json.load(answer)
    if len(data) < 2 or not data[1]:
        return pd.DataFrame(columns=["country", "year", "value"])
    rows = [
        (item["countryiso3code"], int(item["date"]), item["value"])
        for item in data[1]
    ]
    table = pd.DataFrame(rows, columns=["country", "year", "value"])
    table["value"] = table["value"].astype("float64")
    return table.sort_values(["country", "year"]).reset_index(drop=True)


def about(indicator):
    """What the source says its own indicator is: the name, the note, the base.

    The answer is a paragraph the bank wrote, not a sentence we invented, and
    that is the point -- the definition travels with the number instead of
    living in somebody's head.
    """
    with urllib.request.urlopen(ABOUT.format(indicator), timeout=30) as answer:
        data = json.load(answer)
    if len(data) < 2 or not data[1]:
        raise RuntimeError(f"{indicator}: анықтама жоқ")
    item = data[1][0]
    return {
        "код": item["id"],
        "аты": item["name"],
        "дереккөз": (item.get("source") or {}).get("value", ""),
        "анықтама": " ".join((item.get("sourceNote") or "").split()),
    }
