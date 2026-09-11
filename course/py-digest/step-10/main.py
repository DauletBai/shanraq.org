"""Step 10 -- after lesson 29, on joining tables.

The report stops speaking in codes. sholu/anyqtama.py holds what the bank does
not send -- the names and the regions -- and the report joins it on the country
code with a left join, so a country the directory has never heard of stays in
the report under its code.

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
