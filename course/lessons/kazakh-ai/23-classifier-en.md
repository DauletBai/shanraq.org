# Lesson 23. Learn the question type from examples

## Why this matters

A pupil notices that time cards often contain “қашан” while place cards contain “қайда.” Let the program count such patterns. This small classifier chooses a known question type. It has no club facts and does not generate free text.

## Before the program

Use the step 23 `examples.json`: it contains the same 19 cards and three envelopes. A **feature** is an observable part of a question. For each word we use both the full word and its first four letters. `word[:4]` takes positions from the start up to, but not including, position 4. A short word can contribute the same feature twice; this exercise keeps both counts. `weights` stores how often each feature occurred for each training label. `feature not in weights` means a record must be created.

Copy [`examples.json`](https://github.com/DauletBai/shanraq.org/blob/main/course/kazakh-ai/step-23/examples.json) beside the program. Its 19 cards use the format explained in lesson 22.

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-23
python3 classifier.py
```

On Windows, use `py classifier.py`.

```python
import json

with open("examples.json", encoding="utf-8") as file:
    rows = json.load(file)
weights = {}
for row in rows:
    if row["split"] == "train":
        words = row["text"].lower().replace("?", "").replace(",", "").split()
        for word in words:
            for feature in [word, word[:4]]:
                if feature not in weights:
                    weights[feature] = {"уақыт": 0, "орын": 0}
                weights[feature][row["label"]] += 1
for row in rows:
    if row["split"] == "tune":
        scores = {"уақыт": 0, "орын": 0}
        words = row["text"].lower().replace("?", "").replace(",", "").split()
        for word in words:
            for feature in [word, word[:4]]:
                if feature in weights:
                    scores["уақыт"] += weights[feature]["уақыт"]
                    scores["орын"] += weights[feature]["орын"]
        if scores["уақыт"] > scores["орын"]:
            prediction = "уақыт"
        elif scores["орын"] > scores["уақыт"]:
            prediction = "орын"
        else:
            prediction = "білмеймін"
        print(row["text"], "→", prediction)
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-23).

## How the code works

The code learns only from `train`. `weights[feature][row["label"]] += 1` increments the feature count for the human label. For each `tune` question, familiar counts are added into `scores`. The higher score wins; a tie yields `білмеймін`. These are counts, **not probabilities or calibrated confidence**. For `Шахмат қашан және қайда өтеді? Уақыты қандай?`, the classifier incorrectly picks time because two time cues outweigh one place cue. The next lesson adds a stop rule.

## Support map

train → words and first 4 letters → label counts → tune → compare two scores → type or refusal.

## Recall and check

Hide the code and explain why `Шахмат нешеде?` gets a refusal: `нешеде` was absent from training and the club name appeared equally often with both labels. Expected predictions for the five `tune` questions are `уақыт`, `орын`, `білмеймін`, `білмеймін`, `уақыт`. Hint: count the `қашан`, `қайда`, and `уақыты` features. Common mistake: treating a high count as proof that the question has only one meaning.
