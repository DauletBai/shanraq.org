"""The source: a series from the World Bank.

The network lives in this file only. The rest of the code has no business
knowing there is an internet.

One request now brings several countries: the bank takes them separated by
semicolons, and asking three times for what can be asked once is three times the
waiting and three times the chance of being turned away.
"""

import json
import urllib.request

import pandas as pd

API = "https://api.worldbank.org/v2/country/{}/indicator/{}?format=json&per_page=400&date={}:{}"
INFLATION = "FP.CPI.TOTL.ZG"


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
