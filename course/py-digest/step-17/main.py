"""Step 17 -- after lesson 37, the digest that says how usual a year was.

One main() that reads like a table of contents: take, clean, count, draw,
collect. Around it the four things a program run by a scheduler needs and a
program run by a person does not.

  --dry-run   count everything and write nothing
  --offline   never ask the bank; fail if there is nothing on disk
  exit code   0 when the report was rebuilt, 1 when it was not
  the log     one line per step, so the morning after is readable

A source that does not answer leaves yesterday's report where it was. An empty
report is worse than an old one: an old one is at least true about an older day.

New in this step: the page carries a second table -- over every year a country
has, the mean, the median, the spread and the quarters, with the measure for
the heading chosen by a rule (sholu/esep.py, headline).

Run:  python3 main.py [--dry-run] [--offline]
"""

import sys
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


def collect(offline=False, dry=False):
    """Step one: the table, off the disk or off the network.

    A dry run writes nothing at all -- not even the series it just fetched.
    Half a dry run is not a dry run.
    """
    table = saqtau.load(SERIES)
    if table is not None:
        return table, f"дискіден оқылды: {SERIES.name}"
    if offline:
        raise RuntimeError("дискіде дерек жоқ, ал желіге тыйым салынған")
    table = derekkoz.fetch(COUNTRIES, derekkoz.INFLATION, 2021, 2025)
    if dry:
        return table, "банктен алынды, бірақ құрғақ жүрісте сақталмады"
    saqtau.save(table, SERIES)
    return table, f"банктен алынып, сақталды: {SERIES.name}"


def main(argv=()):
    """The whole pipeline, in the order it happens."""
    dry = "--dry-run" in argv
    offline = "--offline" in argv
    log = []

    try:
        table, line = collect(offline=offline, dry=dry)
    except Exception as error:                      # network, disk, format
        print(f"  дереккөз: {error} — есеп бұрынғы күйінде қалды", file=sys.stderr)
        return 1
    log.append(line)

    table, cleaning = tazalau.clean(table)
    log.extend(cleaning)

    report = esep.report(table)
    log.append(f"есеп жолдары: {len(report)}")

    # The second block of the report: over all the years, not over the last one.
    spread = esep.ozara(table)
    log.append("барлық жылдар: " + ", ".join(
        f"{row.аты} — {row.өлшем} {row.мәні}" for row in spread.itertuples(index=False)))

    if dry:
        log.append("құрғақ жүріс: файлдар тиылмады")
    else:
        saqtau.write_text(report.to_csv(index=False), REPORT)
        suret.draw(table, PICTURE, anyqtama.table().set_index("country")["аты"].to_dict())
        bet.write(report, PICTURE, PAGE, TITLE, SOURCE, spread=spread)
        log.append(f"жазылды: {REPORT.name}, {PICTURE.name}, {PAGE.name}")

    for line in log:
        print(f"  {line}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
