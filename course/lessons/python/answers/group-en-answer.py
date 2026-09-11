"""The exercise of lesson 28: a report by month and category.

One groupby over two keys, column names of its own, the share inside a month
through transform -- and the sorting that makes a report readable.
"""

import pandas as pd

rows = [
    ("2026-01-03", "Kostanay", "food", 4200),
    ("2026-01-05", "Kostanay", "fuel", 18500),
    ("2026-01-07", "Rudny", "food", 2600),
    ("2026-01-09", "Kostanay", "phone", 4990),
    ("2026-01-15", "Kostanay", "rent", 95000),
    ("2026-01-18", "Rudny", "fuel", 16200),
    ("2026-01-24", "Kostanay", "food", 3100),
    ("2026-02-02", "Rudny", "food", 2800),
    ("2026-02-06", "Kostanay", "fuel", 19100),
    ("2026-02-11", "Kostanay", "phone", 4990),
    ("2026-02-15", "Kostanay", "rent", 95000),
    ("2026-02-20", "Rudny", "rent", 62000),
]
df = pd.DataFrame(rows, columns=["day", "city", "kind", "amount"])
df["day"] = pd.to_datetime(df["day"])
df["month"] = df["day"].dt.to_period("M")

report = df.groupby(["month", "kind"]).agg(
    receipts=("amount", "count"),
    total=("amount", "sum"),
).reset_index()

# The share inside its own month, not of the whole period.
report["share"] = (report["total"] / report.groupby("month")["total"].transform("sum") * 100).round(1)
report = report.sort_values(["month", "total"], ascending=[True, False])

print(report.to_string(index=False))

print()
totals = df.groupby("month")["amount"].sum()
print("the totals per month:", {str(k): int(v) for k, v in totals.items()})
print("the month that cost most:", totals.idxmax())
