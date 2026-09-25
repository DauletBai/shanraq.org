# Lesson 20. Stop on conflicts and expired records

## Why this matters

If a notice board has two different notices for one activity, choosing the first is unsafe. Once a notice has expired, its old schedule cannot be presented as current. Combine question parsing, lookup, and these two checks.

## Before running

`valid_until` is the last calendar day on which the card is valid, inclusive. It uses the same `year-month-day` form. `from datetime import date` imports Python’s calendar date type. `date.fromisoformat(...)` converts text into a date we can compare. This exercise fixes the date at `2026-09-24`, so its result stays the same in every year; the final program will accept a date from the user. `!=` means “not equal”; `or` means either condition is enough.

```json
[
  {"club": "шахмат", "kind": "уақыт", "value": "бейсенбі, 15:00", "source": "club-sheet-01", "checked_on": "2026-09-01", "valid_until": "2026-12-31"},
  {"club": "шахмат", "kind": "орын", "value": "203-бөлме", "source": "club-sheet-01", "checked_on": "2026-09-01", "valid_until": "2026-12-31"}
]
```

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-20
python3 model.py
```

On Windows, use `py model.py` instead of the last command.

```python
import json
from datetime import date

with open("facts.json", encoding="utf-8") as file:
    facts = json.load(file)
question = input("Сұрақ: ").lower().replace("?", "").replace(",", "")
words = question.split()
entities = []
intents = []
for word in words:
    if word in ["шахмат", "сурет"]:
        entities.append(word)
    if word == "қашан":
        intents.append("уақыт")
    if word == "қайда":
        intents.append("орын")
if len(entities) != 1 or len(intents) != 1:
    print("нақтылаңыз: бір үйірме және бір сұрақ түрі керек")
else:
    club = entities[0]
    kind = intents[0]
    matches = []
    for fact in facts:
        if fact["club"] == club and fact["kind"] == kind:
            matches.append(fact)
    if len(matches) == 0:
        print("білмеймін: дерек жоқ")
    elif len(matches) > 1:
        print("тоқта: бірнеше дерек табылды")
    else:
        fact = matches[0]
        today = date.fromisoformat("2026-09-24")
        until = date.fromisoformat(fact["valid_until"])
        if today > until:
            print("білмеймін: дерек ескірген")
        elif kind == "уақыт":
            print(fact["club"], "үйірмесінің уақыты:", fact["value"])
            print("Дереккөз:", fact["source"], "| тексерілген күні:", fact["checked_on"])
        else:
            print(fact["club"], "үйірмесінің орны:", fact["value"])
            print("Дереккөз:", fact["source"], "| тексерілген күні:", fact["checked_on"])
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-20).

## How the program works

First, the lesson 15 logic extracts exactly one activity and one information type. Then we scan the whole catalog. Zero records means missing information; two or more mean a possible conflict or duplicate that a person should review. Only for one card do we compare dates: `today > until` becomes true the day after `valid_until`. Only then do we print the answer and source. `checked_on` records a past check; `valid_until` sets the expiry. They serve different purposes. A bad source sheet and limited vocabulary can still cause errors; the model does not guarantee truth.

## Support map

Question → exactly two keys → 0/1/many records → expiry of one record → sourced answer or clear stop.

![Lesson 20 support map](/static/course/kazakh-ai/map-20-expiry-en.svg)

## Recall and check

Try `Шахмат қашан?`, then `Сурет қайда?`. The first answer is `шахмат үйірмесінің уақыты: бейсенбі, 15:00` plus a source line. The second is `білмеймін: дерек жоқ`. To test a conflict, add a second `шахмат / уақыт` card with another value: expect `тоқта: бірнеше дерек табылды`. To test expiry, change the exercise date in the code to `2027-01-01`: expect `білмеймін: дерек ескірген`. Common mistake: deleting an old card before finding the cause of the conflict.

## Gate 4: answer only from a card

Show four outcomes: one valid card, no card, two cards with the same key, and one expired card. For each, name the first check that permits or stops an answer. Pass when a factual answer appears only in the first case and includes its source. After an error, rebuild `card count → expiry → answer` and retry.
