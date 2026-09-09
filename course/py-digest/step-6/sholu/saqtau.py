"""Сақтау: қатарды CSV-ға жазу және кері оқу.

Дискі туралы білетін жалғыз орын. Бағанның аттары да осында тұрады.
"""

import csv

FIELDS = ["year", "value", "note"]


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
