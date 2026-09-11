"""Step 13 -- after lesson 32, on the first chart.

The report gains a picture. sholu/suret.py draws one line per country and saves
it beside the CSV; the names come from the same directory the report reads, so
the legend says "Қазақстан" rather than "KAZ".

Run:  python3 main.py
"""

from pathlib import Path

from sholu import anyqtama, derekkoz, esep, saqtau, suret, tazalau

HERE = Path(__file__).parent
STORE = HERE / "data"
SERIES = STORE / "inflation.csv"
REPORT = STORE / "report.csv"
PICTURE = STORE / "inflation.png"
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

    names = anyqtama.table().set_index("country")["аты"].to_dict()
    suret.draw(table, PICTURE, names)
    print(f"сурет салынды: {PICTURE.name}")

    print(REPORT.read_text(encoding="utf-8"), end="")


if __name__ == "__main__":
    main()
