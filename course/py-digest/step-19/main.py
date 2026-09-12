"""Step 19 -- after lesson 40, the digest that reads a second indicator.

One main() that reads like a table of contents: take, clean, count, draw,
collect. Around it the four things a program run by a scheduler needs and a
program run by a person does not.

  --dry-run   count everything and write nothing
  --offline   never ask the bank; fail if there is nothing on disk
  exit code   0 when the report was rebuilt, 1 when it was not
  the log     one line per step, so the morning after is readable

A source that does not answer leaves yesterday's report where it was. An empty
report is worse than an old one: an old one is at least true about an older day.

New in this step: a fourth table, and the first indicator that is not inflation
-- broad money against GDP, that is, how much money a country has for a tenge of
its yearly output. Two series mean the fetching stops being about prices:
collect() takes the indicator and the file it belongs to as arguments, so a
third series would cost a line rather than a fork of this program.

Run:  python3 main.py [--dry-run] [--offline]
"""

import sys
from pathlib import Path

import pandas as pd

from sholu import anyqtama, bet, derekkoz, esep, saqtau, suret, tazalau

HERE = Path(__file__).parent
STORE = HERE / "data"
SERIES = STORE / "inflation.csv"
AQSHA = STORE / "aqsha.csv"
REPORT = STORE / "report.csv"
PICTURE = STORE / "inflation.png"
PAGE = STORE / "report.html"
COUNTRIES = ["KZ", "UZ", "RU"]
TITLE = "Инфляция: Қазақстан, Өзбекстан, Ресей"
SOURCE = "Дереккөз: Дүниежүзілік банк, FP.CPI.TOTL.ZG"
AQSHA_SOURCE = "Ақша дереккөзі: Дүниежүзілік банк, FM.LBL.BMNY.GD.ZS"
FIRST, LAST = 2021, 2025


def collect(indicator, path, offline=False, dry=False):
    """Step one: a table, off the disk or off the network.

    The indicator and the file are arguments now: the digest reads two series,
    and the only thing that differs between them is the code in the address and
    the name on disk.

    A dry run writes nothing at all -- not even the series it just fetched.
    Half a dry run is not a dry run.
    """
    table = saqtau.load(path)
    if table is not None:
        return table, f"дискіден оқылды: {path.name}"
    if offline:
        raise RuntimeError(f"дискіде дерек жоқ ({path.name}), ал желіге тыйым салынған")
    table = derekkoz.fetch(COUNTRIES, indicator, FIRST, LAST)
    if dry:
        return table, f"банктен алынды, бірақ құрғақ жүрісте сақталмады: {path.name}"
    saqtau.save(table, path)
    return table, f"банктен алынып, сақталды: {path.name}"


def main(argv=()):
    """The whole pipeline, in the order it happens."""
    dry = "--dry-run" in argv
    offline = "--offline" in argv
    log = []

    try:
        table, line = collect(derekkoz.INFLATION, SERIES, offline=offline, dry=dry)
        money, money_line = collect(derekkoz.MONEY, AQSHA, offline=offline, dry=dry)
    except Exception as error:                      # network, disk, format
        print(f"  дереккөз: {error} — есеп бұрынғы күйінде қалды", file=sys.stderr)
        return 1
    log.append(line)
    log.append("ақша: " + money_line)

    table, cleaning = tazalau.clean(table)
    log.extend(cleaning)

    # The same cleaning, the same log -- with a word in front, so a line about
    # the second series is not read as a line about the first.
    money, money_cleaning = tazalau.clean(money)
    log.extend("ақша: " + line for line in money_cleaning)

    report = esep.report(table)
    log.append(f"есеп жолдары: {len(report)}")

    # The second block of the report: over all the years, not over the last one.
    spread = esep.ozara(table)
    log.append("барлық жылдар: " + ", ".join(
        f"{row.аты} — {row.өлшем} {row.мәні}" for row in spread.itertuples(index=False)))

    index = esep.indeks(table)
    log.append("баға индексі: " + ", ".join(
        f"{row.аты} — {row.есе} есе" for row in index.itertuples(index=False)))

    aqsha = esep.aqsha(money)
    # A country the money series has nothing for says so in the log too. A line
    # that prints "nan тиын" teaches the reader to ignore the line.
    counted = [f"{row.аты} — " + (esep.UNKNOWN if pd.isna(row.тиын) else f"{row.тиын} тиын")
               for row in aqsha.itertuples(index=False)]
    log.append("ЖІӨ теңгесіне ақша: " + ", ".join(counted))

    if dry:
        log.append("құрғақ жүріс: файлдар тиылмады")
    else:
        saqtau.write_text(report.to_csv(index=False), REPORT)
        suret.draw(table, PICTURE, anyqtama.table().set_index("country")["аты"].to_dict())
        bet.write(report, PICTURE, PAGE, TITLE, SOURCE, spread=spread, index=index,
                  money=aqsha, money_source=AQSHA_SOURCE)
        log.append(f"жазылды: {REPORT.name}, {PICTURE.name}, {PAGE.name}")

    for line in log:
        print(f"  {line}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
