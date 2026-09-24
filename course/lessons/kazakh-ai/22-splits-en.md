# Lesson 22. Separate training, tuning, and testing

## Why this matters

If a pupil has already seen the test answers, a high mark proves little. Put the questions into three envelopes: one for learning, one for adjusting the rule, and one to open only for the final check.

## Before the program

The new `split` field names an envelope: `train` for learning, `tune` for adjusting the rule, and `test` for final testing. This file contains 8 training and 5 tuning questions. Six labeled test questions are kept in a separate step 25 file: do not open it before the final check. We assigned them by hand, so the result is reproducible. The same text must not appear in two envelopes. `seen` stores questions already visited; `valid` records whether we found an error. `True` means valid and `False` means invalid.

```json
[
  {
    "text": "Шахмат қашан?",
    "label": "уақыт",
    "split": "train"
  },
  {
    "text": "Сурет қашан?",
    "label": "уақыт",
    "split": "train"
  },
  {
    "text": "Шахмат уақыты?",
    "label": "уақыт",
    "split": "train"
  },
  {
    "text": "Сурет уақыты?",
    "label": "уақыт",
    "split": "train"
  },
  {
    "text": "Шахмат қайда?",
    "label": "орын",
    "split": "train"
  },
  {
    "text": "Сурет қайда?",
    "label": "орын",
    "split": "train"
  },
  {
    "text": "Шахмат орны?",
    "label": "орын",
    "split": "train"
  },
  {
    "text": "Сурет орны?",
    "label": "орын",
    "split": "train"
  },
  {
    "text": "Шахмат қашан өтеді?",
    "label": "уақыт",
    "split": "tune"
  },
  {
    "text": "Сурет қайда өтеді?",
    "label": "орын",
    "split": "tune"
  },
  {
    "text": "Шахмат нешеде?",
    "label": "уақыт",
    "split": "tune"
  },
  {
    "text": "Шахмат неге?",
    "label": "белгісіз",
    "split": "tune"
  },
  {
    "text": "Шахмат қашан және қайда өтеді? Уақыты қандай?",
    "label": "белгісіз",
    "split": "tune"
  }
]
```

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-22
python3 splits.py
```

On Windows, use `py splits.py`.

```python
import json

with open("examples.json", encoding="utf-8") as file:
    rows = json.load(file)
counts = {"train": 0, "tune": 0, "test": 0}
seen = []
valid = True
for row in rows:
    if row["text"] in seen:
        print("Қайталанған сұрақ:", row["text"])
        valid = False
    seen.append(row["text"])
    split = row["split"]
    if split in counts:
        counts[split] += 1
    else:
        print("Қате бөлік:", split)
        valid = False
print("Оқыту:", counts["train"])
print("Баптау:", counts["tune"])
print("Бақылау:", counts["test"])
print("Жинақ жарамды:", valid)
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-22).

## How the code works

The program checks whether the exact text appeared before, then `seen.append` stores it. It increments the counter for its envelope. An unknown envelope name or a duplicate sets `valid = False`. This check compares text exactly, so a change of case or punctuation could hide a duplicate; that is a limitation of this exercise. We will not use test labels to choose rules before lesson 25.

## Support map

13 open cards → train 8 / tune 5 / test 0; six test cards remain sealed → duplicate check → fair later test.

## Recall and check

Hide the code and explain each envelope. Copy one question record to the end of the file. Hint: `seen` already contains that text. Expected: `Қайталанған сұрақ:` and `Жинақ жарамды: False`; the counts sum to 14. Restore the file afterwards. Common mistake: changing the rule after reading test answers while still calling the test independent.
