"""The answer to lesson 13's exercise: an API answer, and a round trip through JSON.

The year becomes a number while the answer is being parsed rather than after,
null lands among the gaps rather than in the average, and the keys are brought
back to numbers on the way in -- because JSON has no other kind of key.
"""

import json

ANSWER = """[
  {"page": 1, "pages": 1, "per_page": 4, "total": 4},
  [
    {"country": {"id": "KZ", "value": "Kazakhstan"}, "date": "2025", "value": 11.39},
    {"country": {"id": "KZ", "value": "Kazakhstan"}, "date": "2024", "value": 8.69},
    {"country": {"id": "KZ", "value": "Kazakhstan"}, "date": "2023", "value": null},
    {"country": {"id": "KZ", "value": "Kazakhstan"}, "date": "2022", "value": 15.0}
  ]
]"""

data = json.loads(ANSWER)
print("records:", len(data[1]))

series = {}
gaps = []
for row in data[1]:
    year = int(row["date"])
    if row["value"] is None:
        gaps.append(year)
        continue
    series[year] = row["value"]

total = 0.0
for value in series.values():
    total += value
print(f"years with a number: {len(series)}, average: {total / len(series):.2f}")
print("gaps:", gaps)

text = json.dumps(series, ensure_ascii=False)
back = json.loads(text)
print("key before:", type(list(series)[0]).__name__, "| after:", type(list(back)[0]).__name__)

fixed = {}
for key, value in back.items():
    fixed[int(key)] = value
print("after the mend:", fixed[2025])
