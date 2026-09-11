"""Step 12 -- after lesson 31, on time in a table.

The year stops being a number and becomes a date, and the report gains the
question a yearly series can actually answer: how much the last year differs
from the one before it. diff() inside the grouping does that in one line,
because the cleaning step leaves the rows in order.

Run:  python3 main.py
"""

from pathlib import Path

from sholu import derekkoz, esep, saqtau, tazalau

HERE = Path(__file__).parent
STORE = HERE / "data"
SERIES = STORE / "inflation.csv"
REPORT = STORE / "report.csv"
COUNTRIES = ["KZ", "UZ", "RU"]


def main():
    table = saqtau.load(SERIES)
    if table is not None:
        print(f"дискіден оқылды: {SERIES.name}")
    else:
        table = derekkoz.fetch(COUNTRIES, derekkoz.INFLATION, 2021, 2025)
        saqtau.save(table, SERIES)
        print(f"банктен алынып, сақталды: {SERIES.name}")

    table, log = tazalau.clean(table)
    for line in log:
        print(f"  {line}")

    esep.write(table, REPORT)
    print(f"есеп жазылды: {REPORT.name}")
    print(REPORT.read_text(encoding="utf-8"), end="")


if __name__ == "__main__":
    main()
