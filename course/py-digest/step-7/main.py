"""Step 7 -- after lesson 25, on DataFrame and Series.

The series stops being a dictionary. sholu/esep.py holds it as a table whose
rows are labelled by year, and the report is built from that table: the average
skips the missing years by itself, and the peak comes back as a year rather than
as a position in a list.

The disk has not changed: the CSV is still written and read by the csv module.
That is the next lesson.

Run:  python3 main.py
"""

from pathlib import Path

from sholu import derekkoz, esep, saqtau

HERE = Path(__file__).parent
STORE = HERE / "data"
SERIES = STORE / "inflation-kz.csv"
REPORT = STORE / "report.csv"


def main():
    series = saqtau.load(SERIES)
    if series:
        print(f"дискіден оқылды: {SERIES.name}")
    else:
        series = derekkoz.fetch("KZ", derekkoz.INFLATION, 2021, 2025)
        saqtau.save(series, SERIES)
        print(f"банктен алынып, сақталды: {SERIES.name}")

    table = esep.frame(series)
    print(table)

    esep.write(table, REPORT)
    print(f"есеп жазылды: {REPORT.name}")
    print(REPORT.read_text(encoding="utf-8"), end="")


if __name__ == "__main__":
    main()
