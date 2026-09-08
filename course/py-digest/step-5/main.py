"""Step 5 -- after lesson 12, on CSV.

The store becomes a real table. The series is written with csv.DictWriter and
read back with csv.DictReader, so a column is found by its name rather than by
its position, and a value holding a separator cannot break the file.

Run:  python3 main.py
"""

import csv
import json
import urllib.request
from pathlib import Path

API = "https://api.worldbank.org/v2/country/{}/indicator/{}?format=json&per_page=80&date={}:{}"
INFLATION = "FP.CPI.TOTL.ZG"

HERE = Path(__file__).parent
STORE = HERE / "data"
SERIES = STORE / "inflation-kz.csv"
REPORT = STORE / "report.csv"
FIELDS = ["year", "value", "note"]


def fetch(country, indicator, first, last):
    """{жыл: мән} -- банктен алынған қатар, олқылықтарымен қоса."""
    url = API.format(country, indicator, first, last)
    with urllib.request.urlopen(url, timeout=30) as answer:
        data = json.load(answer)
    if len(data) < 2 or not data[1]:
        return {}
    return {int(item["date"]): item["value"] for item in data[1]}


def save(series, path):
    """Қатарды CSV-ға жазады: жыл, мән, ескертпе. Олқылық — n/a."""
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", encoding="utf-8", newline="") as target:
        writer = csv.DictWriter(target, fieldnames=FIELDS)
        writer.writeheader()
        for year in sorted(series):
            value = series[year]
            if value is None:
                writer.writerow({"year": year, "value": "n/a", "note": "дерек жоқ"})
            else:
                writer.writerow({"year": year, "value": f"{value:.2f}", "note": ""})


def load(path):
    """Сақталған CSV-ды оқиды; файл жоқ болса, бос сөздік қайтарады."""
    if not path.exists():
        return {}
    series = {}
    with path.open(encoding="utf-8-sig", newline="") as source:
        for row in csv.DictReader(source):
            value = row["value"]
            series[int(row["year"])] = None if value == "n/a" else float(value)
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

    with REPORT.open("w", encoding="utf-8", newline="") as target:
        writer = csv.writer(target)
        writer.writerow(["metric", "value"])
        writer.writerow(["жылдар", len(series)])
        writer.writerow(["орташа", f"{average(series):.2f}"])
    print(f"есеп жазылды: {REPORT.name}")
    print(REPORT.read_text(encoding="utf-8"), end="")


main()
