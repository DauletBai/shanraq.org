"""Step 9 -- after the lesson on functions.

The digest has names now. average() walks a series and leaves the years without
a figure out of the count; above() puts two series side by side. An assert stops
an empty series where the mistake is, rather than on a division one floor down.

Run:  python3 main.py
"""

import json
import urllib.request

# The World Bank gives its indicators without a key and without registration.
API = "https://api.worldbank.org/v2/country/{}/indicator/{}?format=json&per_page=80&date={}:{}"
INFLATION = "FP.CPI.TOTL.ZG"


def rows(country, indicator, first, last):
    """{жыл: мән} -- Дүниежүзілік банк қайтарғаны, олқылықтарымен қоса."""
    url = API.format(country, indicator, first, last)
    with urllib.request.urlopen(url, timeout=30) as answer:
        data = json.load(answer)
    if len(data) < 2 or not data[1]:
        return {}
    return {int(item["date"]): item["value"] for item in data[1]}


def average(series):
    """Қатардың орташасы; саны жоқ жылдар есепке кірмейді."""
    total = 0.0
    count = 0
    for value in series.values():
        if value is None:
            continue
        total += value
        count += 1
    assert count > 0, "қатарда бірде-бір сан жоқ"
    return total / count


def above(series, other, since=2021):
    """Since жылынан бастап бірінші қатар екіншісінен жоғары болған жылдар."""
    years = []
    for year, value in series.items():
        if year < since or value is None or year not in other:
            continue
        if value > other[year]:
            years.append(year)
    return sorted(years)


def main():
    kz = rows("KZ", INFLATION, 2021, 2025)
    world = rows("WLD", INFLATION, 2021, 2025)

    print("== инфляция, жылдық %")
    print(f"Қазақстан орташасы: {average(kz):.2f}")
    print(f"әлем орташасы:      {average(world):.2f}")
    print(f"әлемдікінен жоғары жылдар: {above(kz, world)}")


main()
