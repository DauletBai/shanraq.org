"""Step 1 -- after lesson 1, "Why count for yourself".

The first state of the digest: it goes to the World Bank for the price index
and the inflation series and prints what it found. No storage, no schedule and
no error handling yet -- each of those arrives in the lesson that needs it.

Run:  python3 main.py
"""

import json
import urllib.request

# Дүниежүзілік банк елдер бойынша көрсеткіштерді кілтсіз әрі тіркеусіз береді.
# FP.CPI.TOTL — тұтыну бағаларының индексі, FP.CPI.TOTL.ZG — пайызбен инфляция.
API = "https://api.worldbank.org/v2/country/{}/indicator/{}?format=json&per_page=80&date={}:{}"


def rows(country, indicator, first, last):
    """{жыл: мән} қайтарады — Дүниежүзілік банк не жауап берсе, сол."""
    url = API.format(country, indicator, first, last)
    with urllib.request.urlopen(url, timeout=30) as answer:
        data = json.load(answer)
    if len(data) < 2 or not data[1]:
        return {}
    return {int(item["date"]): item["value"] for item in data[1] if item["value"] is not None}


def main():
    print("== Қазақстанда бағалар неше есе өсті")
    prices = rows("KZ", "FP.CPI.TOTL", 2010, 2025)
    base, last = min(prices), max(prices)
    times = prices[last] / prices[base]
    print(f"{base} жылғы баға индексі: {prices[base]:.1f}")
    print(f"{last} жылғы баға индексі: {prices[last]:.1f}")
    print(f"бағалар {times:.2f} есе өсті")
    # Бір қатынас екі жаққа да оқылады, әрі оларды шатастыру қымбатқа түседі.
    print(f"{last} жылғы 1000 теңге = {base} жылғы бағамен {1000 / times:.0f} теңге")
    print(f"{base} жылы 1000 теңге тұрған нәрсе {1000 * times:.0f} теңге тұрады")

    print()
    print("== әлем бір, бағалар басқа: жылдар бойынша инфляция, %")
    countries = {"KZ": "Қазақстан", "WLD": "әлем", "GE": "Грузия",
                 "AM": "Армения", "PL": "Польша"}
    years = range(2021, 2026)
    print("ел          " + "".join(f"{year:>7}" for year in years))
    for code, name in countries.items():
        inflation = rows(code, "FP.CPI.TOTL.ZG", 2021, 2025)
        line = "".join(f"{inflation[year]:7.1f}" if year in inflation else f"{'—':>7}"
                       for year in years)
        print(f"{name:<12}{line}")


main()
