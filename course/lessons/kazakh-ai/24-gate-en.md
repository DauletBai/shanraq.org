# Lesson 24. Require enough evidence before choosing

## Why this matters

A receptionist sees both “when” and “where” in a note. Even if “when” appears twice, the receptionist should ask which question the writer wants answered. The previous counter made exactly that mistake.

## Before the program

Keep the learned counts but check the question before choosing. `time_cue` means an explicit time word is present; `place_cue` means an explicit place word is present. We allow only the already studied forms `қашан`, `уақыты`, `қайда`, and `орны`. `or` accepts either condition, `and` requires both, and `not` reverses a yes/no value. Choose a type only when its cue is present, the other type’s cue is absent, and its score is higher. Two kinds of cue or no clear cue lead to refusal.

Copy [`examples.json`](https://github.com/DauletBai/shanraq.org/blob/main/course/kazakh-ai/step-24/examples.json) beside the program. Its 13 cards use the format explained in lesson 22.

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-24
python3 gate.py
```

On Windows, use `py gate.py`.

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
        time_cue = "қашан" in words or "уақыты" in words
        place_cue = "қайда" in words or "орны" in words
        if time_cue and not place_cue and scores["уақыт"] > scores["орын"]:
            prediction = "уақыт"
        elif place_cue and not time_cue and scores["орын"] > scores["уақыт"]:
            prediction = "орын"
        else:
            prediction = "білмеймін"
        print(row["text"], "→", prediction)
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-24).

## How the code works

The last tuning question now yields `білмеймін` even though the time score is larger: `place_cue` is true too. `Шахмат нешеде?` is still unanswered because that form is not yet covered by examples and rules. We can add examples later and check the change on a new test set. **Enough evidence to choose a question type** is not enough to answer a timetable question; lessons 16–20 still require a fact card, source, and valid date.

## Support map

Learned score + explicit cue for one type + no cue for the other → type; otherwise → refuse. Fact lookup remains separate.

![Lesson 24 support map](/static/course/kazakh-ai/map-24-gate-en.svg)

## Recall and check

Hide the code and recall the three conditions for choosing time. Run the five tuning questions. Expected predictions: `уақыт`, `орын`, then `білмеймін` three times. Hint: on the last question both cues are true; `not place_cue` blocks the choice. Common mistake: dropping the other-cue check because the time score is higher. That recreates the lesson 23 error.
