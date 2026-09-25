# Lesson 19. Build an answer from a template

## Why this matters

A form has fixed wording and one blank for a checked value. Such a template does not invent a timetable: it only displays the selected card clearly.

## Before running

Use the cards from step 18. `kind` selects one of two prepared phrases, one for time and one for place. `value` fills the blank. The `|` inside quotes is a printed separator between source and date, not a Python operation.

```json
[
  {"club": "шахмат", "kind": "уақыт", "value": "бейсенбі, 15:00", "source": "club-sheet-01", "checked_on": "2026-09-01"},
  {"club": "шахмат", "kind": "орын", "value": "203-бөлме", "source": "club-sheet-01", "checked_on": "2026-09-01"}
]
```

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-19
python3 answer.py
```

On Windows, use `py answer.py` instead of the last command.

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
    if kind == "уақыт":
        print(fact["club"], "үйірмесінің уақыты:", fact["value"])
    elif kind == "орын":
        print(fact["club"], "үйірмесінің орны:", fact["value"])
    print("Дереккөз:", fact["source"], "| тексерілген күні:", fact["checked_on"])
elif len(matches) == 0:
    print("білмеймін: дерек жоқ")
else:
    print("тоқта: бірнеше дерек табылды")
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-19).

## How the program works

The name comes from `fact["club"]`, so the wording follows the selected card. `print` inserts spaces between parts. No answer is built for zero or multiple records. This still uses two separate input keys; the next lesson will connect it to a full question.

## Support map

Two keys → one card → time/place sentence + value → source + date.

![Lesson 19 support map](/static/course/kazakh-ai/map-19-template-en.svg)

## Recall and check

Enter `шахмат` and `уақыт`. Hint: the first template branch applies. Expected: `шахмат үйірмесінің уақыты: бейсенбі, 15:00` and `Дереккөз: club-sheet-01 | тексерілген күні: 2026-09-01`. Change the place value and check that the time answer stays the same. Common mistake: hard-coding the time in the template, which can then disagree with the catalog.
