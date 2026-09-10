"""The exercise of lesson 26: an accountant's export becomes a table.

The file came out of Excel: semicolons, a comma inside the decimals, "n/a" where
a value is missing and a BOM at the front. None of that is a loop; all of it is
an argument to read_csv.
"""

import sqlite3
from pathlib import Path

import pandas as pd

HERE = Path(__file__).resolve().parent
DIRTY = HERE / "export.csv"
CLEAN = HERE / "clean.csv"
DB = HERE / "digest.db"

# This is how the file looks on disk. The program writes it itself, so that the
# exercise is the same for everybody.
DIRTY.write_text(
    "year;category;amount;note\n"
    "2024;food;1 200,50;\n"
    "2024;phone;4 000,00;the tariff changed\n"
    "2025;food;1 380,75;\n"
    "2025;phone;n/a;the bill never came\n",
    encoding="utf-8-sig",
)

df = pd.read_csv(
    DIRTY,
    sep=";",
    decimal=",",
    thousands=" ",
    encoding="utf-8-sig",
    na_values=["n/a"],
    usecols=["year", "category", "amount"],
)
print("read:")
print(df)
print("types:", dict(df.dtypes.astype(str)))
print("missing:", int(df["amount"].isna().sum()))

print()
print("clean:")
df.to_csv(CLEAN, index=False, encoding="utf-8")
print(CLEAN.read_text(encoding="utf-8"), end="")

print()
db = sqlite3.connect(DB)
df.to_sql("spending", db, if_exists="replace", index=False)
answer = pd.read_sql(
    "SELECT year, SUM(amount) AS amount FROM spending GROUP BY year ORDER BY year", db
)
db.close()
print("by year:")
print(answer.to_string(index=False))
