"""The answer to lesson 12's exercise: an export with a comma inside a field.

It is written by csv.writer rather than by gluing strings, so the field with a
comma in it is quoted for us, and read by DictReader, so a column is found by
its name. The row without a number is skipped rather than counted as zero.
"""

import csv
from pathlib import Path

HERE = Path(__file__).parent
data = HERE / "vygruzka.csv"

rows = [
    ["month", "price", "note"],
    ["январь", "520", "обычный месяц"],
    ["февраль", "546,5", "шок, спрос, логистика"],
    ["март", "n/a", "цены нет"],
]

with data.open("w", encoding="utf-8", newline="") as target:
    csv.writer(target).writerows(rows)

line = data.read_text(encoding="utf-8").splitlines()[2]
print("split(','):", len(line.split(",")), "кусков")
with data.open(encoding="utf-8", newline="") as source:
    print("csv.reader:", len(list(csv.reader(source))[2]), "поля")

total = 0.0
count = 0
with data.open(encoding="utf-8", newline="") as source:
    for row in csv.DictReader(source):
        if row["price"] == "n/a":
            print(f"{row['month']}: пропуск — {row['note']}")
            continue
        total += float(row["price"].replace(",", "."))
        count += 1
        print(f"{row['month']}: {row['price']} — {row['note']}")
print(f"месяцев взято: {count}, среднее {total / count:.2f}")

data.unlink()
