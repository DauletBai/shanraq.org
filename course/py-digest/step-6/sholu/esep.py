"""Есеп: санау және адам көретін нәрсені жазу."""

import csv


def average(series):
    """Қатар бойынша орташа; олқылықтар есепке кірмейді."""
    total = 0.0
    count = 0
    for value in series.values():
        if value is None:
            continue
        total += value
        count += 1
    assert count > 0, "қатарда бірде-бір сан жоқ"
    return total / count


def write(series, path):
    """Есепті CSV-ға жазады: кестеде ашуға келетін екі баған."""
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", encoding="utf-8", newline="") as target:
        writer = csv.writer(target)
        writer.writerow(["metric", "value"])
        writer.writerow(["жылдар", len(series)])
        writer.writerow(["орташа", f"{average(series):.2f}"])
