"""The answer to lesson 16's exercise: a file read one value at a time.

The generator is the whole point: nothing here holds the file, and the rows that
cannot be parsed are dropped where they are read rather than collected into a
second list. One pass counts everything the report needs, because a generator
gives no second pass.
"""

from pathlib import Path

HERE = Path(__file__).parent
data = HERE / "vygruzka.csv"

rows = ["2024;8.7", "2025;n/a", "2025;11.4", "2026;15.0", "мусор", "2026;12.0"]
data.write_text("\n".join(rows) + "\n", encoding="utf-8")


def values(path):
    """Отдаёт (год, значение) по одному; негодные строки пропускает."""
    with path.open(encoding="utf-8") as source:
        for line in source:
            year, _, text = line.strip().partition(";")
            try:
                yield int(year), float(text)
            except ValueError:
                continue


stream = values(data)
print("первое:", next(stream))
print("второе:", next(stream))

total = 0.0
count = 0
years = []
for year, value in values(data):
    total += value
    count += 1
    years.append(year)
print(f"строк взято: {count}, среднее: {total / count:.2f}")
print("годы после 2024:", [year for year in years if year > 2024])

data.unlink()
