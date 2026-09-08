"""Step 4 -- after lesson 11, on files and pathlib.

The digest gets a disk. The series it fetched is saved beside the program, and
the next run reads the file instead of asking the bank again. The report is
written as a second file.

Paths start at Path(__file__).parent, never at the current folder: the program
has to work when it is started from somewhere else.

Run:  python3 main.py
"""

import json
import urllib.request
from pathlib import Path

API = "https://api.worldbank.org/v2/country/{}/indicator/{}?format=json&per_page=80&date={}:{}"
INFLATION = "FP.CPI.TOTL.ZG"

HERE = Path(__file__).parent
STORE = HERE / "data"
SERIES = STORE / "inflation-kz.txt"
REPORT = STORE / "report.txt"


def fetch(country, indicator, first, last):
    """{жыл: мән} -- банктен алынған қатар, олқылықтарымен қоса."""
    url = API.format(country, indicator, first, last)
    with urllib.request.urlopen(url, timeout=30) as answer:
        data = json.load(answer)
    if len(data) < 2 or not data[1]:
        return {}
    return {int(item["date"]): item["value"] for item in data[1]}


def save(series, path):
    """Қатарды 'жыл;мән' түрінде сақтайды; олқылық n/a болып жазылады."""
    lines = []
    for year in sorted(series):
        value = series[year]
        lines.append(f"{year};{'n/a' if value is None else value}")
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text("\n".join(lines) + "\n", encoding="utf-8")


def load(path):
    """Сақталған қатарды оқиды; жоқ болса, бос сөздік қайтарады."""
    if not path.exists():
        return {}
    series = {}
    with path.open(encoding="utf-8") as source:
        for line in source:
            year, _, value = line.strip().partition(";")
            series[int(year)] = None if value == "n/a" else float(value)
    return series


def average(series):
    total = 0.0
    count = 0
    for value in series.values():
        if value is None:
            continue
        total += value
        count += 1
    assert count > 0, "қатарда бірде-бір сан жоқ"
    return total / count


def main():
    series = load(SERIES)
    if series:
        print(f"дискіден оқылды: {SERIES.name}")
    else:
        series = fetch("KZ", INFLATION, 2021, 2025)
        save(series, SERIES)
        print(f"банктен алынып, сақталды: {SERIES.name}")

    # Тізім өрнегі әлі өтілген жоқ: 16-сабаққа дейін бәрі кәдімгі циклмен.
    lines = []
    for year in sorted(series):
        value = series[year]
        if value is None:
            lines.append(f"{year}: —")
        else:
            lines.append(f"{year}: {value:.2f}")
    lines.append(f"орташа: {average(series):.2f}")
    REPORT.write_text("\n".join(lines) + "\n", encoding="utf-8")

    print(f"есеп жазылды: {REPORT.name}, жол саны: {len(lines)}")
main()
