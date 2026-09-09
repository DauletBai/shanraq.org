"""Дереккөз: Дүниежүзілік банктен қатарды алу.

Желі осы файлда ғана. Қалған кодтың интернет туралы білуінің қажеті жоқ.
"""

import json
import urllib.request

API = "https://api.worldbank.org/v2/country/{}/indicator/{}?format=json&per_page=80&date={}:{}"
INFLATION = "FP.CPI.TOTL.ZG"


def fetch(country, indicator, first, last):
    """{жыл: мән} -- банктен алынған қатар, олқылықтарымен қоса."""
    url = API.format(country, indicator, first, last)
    with urllib.request.urlopen(url, timeout=30) as answer:
        data = json.load(answer)
    if len(data) < 2 or not data[1]:
        return {}
    return {int(item["date"]): item["value"] for item in data[1]}
