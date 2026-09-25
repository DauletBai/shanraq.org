# Lesson 18. Source and check date

## Why this matters

A class time on a notice board is useful when we know which notice it came from and when someone checked it. Add those details to each card.

## Before running

Use this step’s files. New fields: `source` is an exercise ID for the club sheet; `checked_on` is the check date in `year-month-day` form. `2026-09-01` means 1 September 2026; leading zeros keep dates consistent. `club-sheet-01` is an exercise label, not a link to a real document.

```json
[
  {"club": "шахмат", "kind": "уақыт", "value": "бейсенбі, 15:00", "source": "club-sheet-01", "checked_on": "2026-09-01"},
  {"club": "шахмат", "kind": "орын", "value": "203-бөлме", "source": "club-sheet-01", "checked_on": "2026-09-01"}
]
```

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-18
python3 source.py
```

On Windows, use `py source.py` instead of the last command.

```python
import json

with open("facts.json", encoding="utf-8") as file:
    facts = json.load(file)
club = input("Үйірме: ").lower()
kind = input("Мәлімет түрі: ").lower()
matches = []
for fact in facts:
    if fact["club"] == club and fact["kind"] == kind:
        matches.append(fact)
if len(matches) == 1:
    fact = matches[0]
    print("Дерек:", fact["value"])
    print("Дереккөз:", fact["source"])
    print("Тексерілген күні:", fact["checked_on"])
elif len(matches) == 0:
    print("білмеймін: дерек жоқ")
else:
    print("тоқта: бірнеше дерек табылды")
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-18).

## How the program works

Lookup is unchanged. With one match, `fact = matches[0]` gives the selected card a short name. We print its value, source, and date. A check date does not mean the record is still valid today; lesson 20 adds an expiry date. The source itself can also be wrong. The program only shows where its answer came from.

## Support map

Keys → one card → value + exercise source ID + check date.

![Lesson 18 support map](/static/course/kazakh-ai/map-18-provenance-en.svg)

## Recall and check

Enter `шахмат` and `орын`. Hint: find the card matching both fields. Expected lines: `Дерек: 203-бөлме`, `Дереккөз: club-sheet-01`, `Тексерілген күні: 2026-09-01`. Change only that card’s `source` and inspect the new line. Common mistake: treating the check date as an expiry date.
