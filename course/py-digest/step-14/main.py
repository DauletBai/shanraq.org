"""Step 14 -- after lesson 34, on the report as a page.

The digest stops handing over three files. sholu/bet.py collects the table, the
picture and the source line into one HTML file: it opens in any browser, prints
to PDF and travels as a single attachment.

The source and the date live on the page, which is the debt lesson 33 left: a
picture sent on its own said nothing about where its numbers came from.

Run:  python3 main.py
"""

from pathlib import Path

from sholu import anyqtama, bet, derekkoz, esep, saqtau, suret, tazalau

HERE = Path(__file__).parent
STORE = HERE / "data"
SERIES = STORE / "inflation.csv"
REPORT = STORE / "report.csv"
PICTURE = STORE / "inflation.png"
PAGE = STORE / "report.html"
COUNTRIES = ["KZ", "UZ", "RU"]
TITLE = "Инфляция: Қазақстан, Өзбекстан, Ресей"
SOURCE = "Дереккөз: Дүниежүзілік банк, FP.CPI.TOTL.ZG"


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

    bet.write(esep.report(table), PICTURE, PAGE, TITLE, SOURCE)
    print(f"бет жиналды: {PAGE.name}")


if __name__ == "__main__":
    main()
