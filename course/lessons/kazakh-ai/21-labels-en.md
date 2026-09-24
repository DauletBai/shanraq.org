# Lesson 21. Label questions by hand

## Why this matters

A teacher prepares answer cards before marking a pupil’s work. A model likewise needs a question paired with a correct exercise label. Without labeled examples, it has nothing to learn from and we have no way to inspect its mistakes.

## Before the program

`examples.json` holds 19 fictional questions. Each card has `text`, the question, and `label`, the human exercise label: `уақыт` for time, `орын` for place, and `белгісіз` for another topic or two requests at once. A label says what the question asks; it is not a timetable fact. Before labeling, ask, “What information does this person want?” If neither type fits, do not force a label. Save the UTF-8 JSON beside `labels.py`.

```json
[
  {
    "text": "Шахмат қашан?",
    "label": "уақыт"
  },
  {
    "text": "Сурет қашан?",
    "label": "уақыт"
  },
  {
    "text": "Шахмат уақыты?",
    "label": "уақыт"
  },
  {
    "text": "Сурет уақыты?",
    "label": "уақыт"
  },
  {
    "text": "Шахмат қайда?",
    "label": "орын"
  },
  {
    "text": "Сурет қайда?",
    "label": "орын"
  },
  {
    "text": "Шахмат орны?",
    "label": "орын"
  },
  {
    "text": "Сурет орны?",
    "label": "орын"
  },
  {
    "text": "Шахмат қашан өтеді?",
    "label": "уақыт"
  },
  {
    "text": "Сурет қайда өтеді?",
    "label": "орын"
  },
  {
    "text": "Шахмат нешеде?",
    "label": "уақыт"
  },
  {
    "text": "Шахмат неге?",
    "label": "белгісіз"
  },
  {
    "text": "Шахмат қашан және қайда өтеді? Уақыты қандай?",
    "label": "белгісіз"
  },
  {
    "text": "Сурет қашан басталады?",
    "label": "уақыт"
  },
  {
    "text": "Шахмат қайда болады?",
    "label": "орын"
  },
  {
    "text": "Шахмат қай күні?",
    "label": "уақыт"
  },
  {
    "text": "Сурет қай жерде?",
    "label": "орын"
  },
  {
    "text": "Сурет неге?",
    "label": "белгісіз"
  },
  {
    "text": "Шахмат кімге?",
    "label": "белгісіз"
  }
]
```

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-21
python3 labels.py
```

On Windows, use `py labels.py`.

```python
import json

with open("examples.json", encoding="utf-8") as file:
    rows = json.load(file)
counts = {"уақыт": 0, "орын": 0, "белгісіз": 0}
for row in rows:
    label = row["label"]
    if label in counts:
        counts[label] += 1
    else:
        print("Қате белгі:", label)
print("Уақыт:", counts["уақыт"])
print("Орын:", counts["орын"])
print("Белгісіз:", counts["белгісіз"])
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-21).

## How the code works

`counts` holds three counters. `for row in rows` visits every card, and `row["label"]` reads its human label. `label in counts` checks that the label is allowed; `+= 1` adds one to the previous count. An unfamiliar label is reported and not counted. We are checking the annotations, not training a model yet.

## Support map

Question → human label → JSON list → allowed-label check → counts 8 / 7 / 4.

## Recall and check

Hide the code and recall the three labels. Add `Сурет нешеде?` with label `уақыт`. Hint: separate the new card with a comma. Expected counts: `Уақыт: 9`, `Орын: 7`, `Белгісіз: 4`. Then try the invalid label `уақ` and find `Қате белгі: уақ`. Common mistake: treating `белгісіз` as a missing catalog fact; here it labels the *question*, while fact lookup happens separately.
