"""Step 9 -- after lesson 28, on grouping and aggregates.

The digest stops looking at one country. Three come back in one request as a
long table -- country, year, value -- and the report is a single groupby over
it: the years, the gaps, the average, the peak and the year of the peak, one row
per country.

Run:  python3 main.py
"""

from pathlib import Path

from sholu import derekkoz, esep, saqtau

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

    print(f"жол саны: {len(table)}, ел саны: {table['country'].nunique()}")

    esep.write(table, REPORT)
    print(f"есеп жазылды: {REPORT.name}")
    print(REPORT.read_text(encoding="utf-8"), end="")


if __name__ == "__main__":
    main()
