"""The exercise of lesson 30: a cleaning with a log.

Every step leaves a trace: how many there were, how many there are, what exactly
changed. A cleaning without a log is data nobody can say what was done to.
"""

import pandas as pd

rows = [
    ("Kostanay", "food", "4 200"),
    ("kostanay ", "fuel", "18 500"),
    ("Rudny", "food", "2 600"),
    ("Rudny", "food", "2 600"),
    ("KOSTANAY", "phone", "n/a"),
    ("Rudny", "rent", "62 000"),
    ("Kostanay", "rent", "95 000"),
    ("Astana", "food", "—"),
]
df = pd.DataFrame(rows, columns=["city", "kind", "amount"])

print("the cleaning log:")
print("  rows in:", len(df))
print("  cities in:", len(df["city"].unique()))

df["city"] = df["city"].str.strip().str.capitalize()
print("  cities out:", len(df["city"].unique()), "—", sorted(df["city"].unique()))

df["amount"] = pd.to_numeric(df["amount"].str.replace(" ", "", regex=False), errors="coerce")
print("  not numbers:", int(df["amount"].isna().sum()))

before = len(df)
df = df.drop_duplicates()
print(f"  duplicates dropped: {before - len(df)}")

known = df.dropna(subset=["amount"])
print("  rows into the report:", len(known), "out of", len(df))

print()
print("per city:")
report = known.groupby("city").agg(
    receipts=("amount", "count"),
    total=("amount", "sum"),
    average=("amount", "mean"),
).round(0).sort_values("total", ascending=False)
print(report.to_string())

print()
print("the sums agree:", int(report["total"].sum()) == int(known["amount"].sum()))
