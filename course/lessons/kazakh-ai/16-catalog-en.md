# Lesson 16. A catalog of checked records

## Why this matters

A library gives each book a card with a title and shelf. We will give each exercise fact a similar record. Here “checked” means the fictional club author entered it into our learning catalog; it is not a real timetable.

## Before running

JSON is a text format for storing data. Picture a stack of cards: square brackets `[ ]` enclose the stack; braces `{ }` enclose one card. Commas separate cards and fields, a colon separates a field name from its value, and text uses double quotes. `club` names an activity, `kind` gives the information type, and `value` holds the fact. Do not repeat a field name on one card. Save `facts.json` as UTF-8 beside the program.

```json
[
  {"club": "шахмат", "kind": "уақыт", "value": "бейсенбі, 15:00"},
  {"club": "шахмат", "kind": "орын", "value": "203-бөлме"}
]
```

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-16
python3 catalog.py
```

On Windows, use `py catalog.py` instead of the last command.

```python
import json

with open("facts.json", encoding="utf-8") as file:
    facts = json.load(file)
print("Жазба саны:", len(facts))
for fact in facts:
    print(fact["club"], fact["kind"], fact["value"])
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-16).

## How the program works

`import json` brings in Python’s built-in JSON tools. `with open(...) as file` keeps the file open for the indented block, then closes it. `encoding="utf-8"` preserves Kazakh letters. `json.load(file)` turns the cards into a Python list. `fact["club"]` selects a field on one card. `len` counts the cards.

## Support map

JSON file → list of cards → three fields per card → display; no lookup or answer yet.

![Lesson 16 support map](/static/course/kazakh-ai/map-16-catalog-en.svg)

## Recall and check

Hide the code and point to one card’s boundaries. Add a `сурет / орын / 105-бөлме` card and predict the count. Hint: separate cards with a comma. Expected: `Жазба саны: 3` and `сурет орын 105-бөлме`. Common mistake: omitting JSON string quotes makes the file unreadable.
