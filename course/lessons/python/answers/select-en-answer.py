"""The exercise of lesson 27: three questions to one table of receipts.

The expensive receipts of one city, a "large" mark set through loc, and the sum
per city over the large ones only -- all of them selections, not loops.
"""

import pandas as pd

rows = [
    ("2026-01-03", "Kostanay", "food", 4200),
    ("2026-01-05", "Kostanay", "fuel", 18500),
    ("2026-01-07", "Rudny", "food", 2600),
    ("2026-01-09", "Kostanay", "phone", 4990),
    ("2026-01-12", "Astana", "food", 7300),
    ("2026-01-15", "Kostanay", "rent", 95000),
    ("2026-01-18", "Rudny", "fuel", 16200),
    ("2026-01-21", "Astana", "phone", 4990),
    ("2026-01-24", "Kostanay", "food", 3100),
    ("2026-01-27", "Rudny", "rent", 62000),
]
df = pd.DataFrame(rows, columns=["day", "city", "kind", "amount"])
df["day"] = pd.to_datetime(df["day"])

print("over 5000 in Kostanay:")
mask = (df["city"] == "Kostanay") & (df["amount"] > 5000)
print(df.loc[mask, ["day", "amount"]].to_string(index=False))

print()
print("the size of a receipt:")
df["size"] = "ordinary"
df.loc[df["amount"] >= 20000, "size"] = "large"
print(df["size"].value_counts().to_string())

print()
print("the large ones per city:")
big = df.loc[df["size"] == "large", ["city", "amount"]]
print(big.groupby("city")["amount"].sum().to_string())
