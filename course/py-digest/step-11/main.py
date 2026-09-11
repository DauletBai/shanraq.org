"""Step 11 -- after lesson 30, on dirty data.

A cleaning step stands between the source and the report. sholu/tazalau.py drops
the rows without a key and the repeated ones, turns what is not a number into a
gap rather than a zero, and says in a log what it did. The report no longer
decides any of that; it only counts.

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
