"""Step 27 -- the complete digest, now guarded by regression tests.

One main() that reads like a table of contents: take, clean, count, draw,
collect. Around it the four things a program run by a scheduler needs and a
program run by a person does not.

  --dry-run   count everything and write nothing
  --offline   never ask the bank; fail if there is nothing on disk
  exit code   0 when the report was rebuilt, 1 when it was not
  the log     one line per step, so the morning after is readable

A source that does not answer leaves yesterday's report where it was. An empty
report is worse than an old one: an old one is at least true about an older day.

The application behaviour is unchanged from step 26. The new work lives in
test_policy.py and test_page.py: project promises can now be checked with one
pytest run before a report is published.

Run:  python3 main.py [--dry-run] [--offline]
"""

import sys
from pathlib import Path

import pandas as pd

from sholu import anyqtama, bet, derekkoz, esep, model, saqtau, suret, tazalau

HERE = Path(__file__).parent
STORE = HERE / "data"
SERIES = STORE / "inflation.csv"
AQSHA = STORE / "aqsha.csv"
NOTES = STORE / "indicators.json"
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


def explain(indicators, offline=False, dry=False):
    """What the source says its indicators are, off the disk or off the network.

    The definitions almost never change, so they are read from disk when they
    are there, and a source that will not answer costs the page its notes and
    nothing else.
    """
    stored = saqtau.load_json(NOTES)
    if stored and sorted(n["код"] for n in stored) == sorted(indicators):
        return stored, f"дискіден оқылды: {NOTES.name}"
    # --offline promises the bank will not be asked, and a note is not worth
    # breaking that promise for.
    if offline:
        raise RuntimeError("дискіде анықтама жоқ, ал желіге тыйым салынған")
    notes = [derekkoz.about(code) for code in indicators]
    if dry:
        return notes, "банктен алынды, бірақ құрғақ жүрісте сақталмады: indicators.json"
    saqtau.save_json(notes, NOTES)
    return notes, f"банктен алынып, сақталды: {NOTES.name}"


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

    # The notes are asked for after the tables are counted: they decorate the
    # report, and nothing about them is worth failing a run over.
    try:
        notes, notes_line = explain([derekkoz.INFLATION, derekkoz.MONEY],
                                    offline=offline, dry=dry)
        log.append("анықтама: " + notes_line)
    except Exception as error:                      # network, disk, format
        notes = None
        log.append(f"анықтама алынбады: {error} — есеп анықтамасыз жиналады")

    aqsha = esep.aqsha(money)
    # A country the money series has nothing for says so in the log too. A line
    # that prints "nan тиын" teaches the reader to ignore the line.
    counted = [f"{row.аты} — " + (esep.UNKNOWN if pd.isna(row.тиын) else f"{row.тиын} тиын")
               for row in aqsha.itertuples(index=False)]
    log.append("ЖІӨ теңгесіне ақша: " + ", ".join(counted))

    link = esep.baylanys(report, aqsha)
    log.append("баға мен ақшаның байланысы: " + str(link.iloc[0]["қорытынды"]))

    trend = esep.trend(table)
    # to_dict here rather than itertuples: one of the columns is two words, and
    # a namedtuple would rename it to a number.
    log.append("тренд: " + ", ".join(
        f"{row['аты']} — {row['жылдық']} есе/жыл, өз жылдарында {row['ең үлкен қате']}, "
        f"көрмеген жылда {row['соңғы жылда']}"
        for row in trend.to_dict("records")))

    limits = esep.shekteu(report, aqsha, link, trend)
    log.append("шектеулер: " + str(len(limits)) + " жол")

    source_record = model.source_from(report)
    verified = model.verify(model.sample_response(source_record), source_record)
    trust_note = model.trust_note(source_record)
    log.append("модель жауабы дереккөзбен салыстырылды: " + source_record["source_id"])

    if dry:
        log.append("құрғақ жүріс: файлдар тиылмады")
    else:
        saqtau.write_text(report.to_csv(index=False), REPORT)
        suret.draw(table, PICTURE, anyqtama.table().set_index("country")["аты"].to_dict())
        bet.write(report, PICTURE, PAGE, TITLE, SOURCE, spread=spread, index=index,
                  money=aqsha, money_source=AQSHA_SOURCE, link=link, trend=trend,
                  limits=limits, about=notes, verified=verified, trust_note=trust_note)
        log.append(f"жазылды: {REPORT.name}, {PICTURE.name}, {PAGE.name}")

    for line in log:
        print(f"  {line}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
