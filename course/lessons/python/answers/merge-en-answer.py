"""The exercise of lesson 29: data and a directory, joined on a key.

Joined so that not a single row of data is lost, with a count of the ones that
found no name, and only then grouped by region.
"""

import pandas as pd

data = pd.DataFrame(
    [("KAZ", 2023, 14.5), ("KAZ", 2024, 8.7),
     ("UZB", 2023, 10.0), ("UZB", 2024, 9.6),
     ("RUS", 2023, 5.9), ("RUS", 2024, 8.4),
     ("KGZ", 2023, 10.8), ("KGZ", 2024, 6.3)],
    columns=["code", "year", "value"],
)
names = pd.DataFrame(
    [("KAZ", "Kazakhstan", "Central Asia"),
     ("UZB", "Uzbekistan", "Central Asia"),
     ("KGZ", "Kyrgyzstan", "Central Asia"),
     ("TJK", "Tajikistan", "Central Asia")],
    columns=["code", "name", "region"],
)

# The data is on the left and must not be lost: the directory is incomplete.
joined = data.merge(names, on="code", how="left", validate="many_to_one")
print("rows before:", len(data), "| after:", len(joined))
print("without a name:", sorted(joined.loc[joined["name"].isna(), "code"].unique()))

print()
print("average inflation per country:")
by_country = joined.groupby(["code", "name"], dropna=False)["value"].mean().round(2)
print(by_country.reset_index().to_string(index=False))

print()
print("per region:")
# Rows with no region are left out of the grouping: their region is not empty
# but unknown, and filing them under "Central Asia" would be inventing data.
by_region = joined.dropna(subset=["region"]).groupby("region").agg(
    countries=("code", "nunique"),
    average=("value", "mean"),
).round(2)
print(by_region.to_string())
