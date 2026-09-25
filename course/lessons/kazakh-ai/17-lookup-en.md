# Lesson 17. Look up a record by two keys

## Why this matters

An author name alone may identify several library books. Our catalog also needs two keys: activity and information type. Lesson 15 already extracted this pair from a short question.

## Before running

Copy `facts.json` from step 16 into the new folder. This program asks for the two keys separately so we can test lookup first. `and` requires both comparisons to be true. `matches` stores matching cards.

```json
[
  {"club": "шахмат", "kind": "уақыт", "value": "бейсенбі, 15:00"},
  {"club": "шахмат", "kind": "орын", "value": "203-бөлме"}
]
```

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-17
python3 search.py
```

On Windows, use `py search.py` instead of the last command.

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
    print("Табылды:", matches[0]["value"])
elif len(matches) == 0:
    print("білмеймін: дерек жоқ")
else:
    print("тоқта: бірнеше дерек табылды")
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-17).

## How the program works

The loop visits every card. A card enters `matches` only when both fields match. With exactly one match, `[0]` is safe; with none, we say no record exists; with several, we stop because we cannot choose yet. `.lower()` handles capitals but does not fix spelling or extra spaces.

## Support map

Activity + type → scan catalog → 0: missing; 1: value; 2 or more: stop.

![Lesson 17 support map](/static/course/kazakh-ai/map-17-lookup-en.svg)

## Recall and check

Predict the output for `шахмат` + `уақыт`, then `сурет` + `орын`. Hint: both fields must match. Expected: `Табылды: бейсенбі, 15:00`; then `білмеймін: дерек жоқ`. Add a second chess time card and verify the stop. Common mistake: accepting the first match before scanning the whole catalog.
