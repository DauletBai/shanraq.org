# Lesson 25. Measure correct answers and refusals

## Why this matters

A teacher distinguishes a wrong answer from an honest “I do not know yet.” Too many refusals are still a problem. One attractive number is not enough: count correct answers and known questions left unanswered.

## Before the program

Now open the six `test` cards. Do not change the training code or lesson 24 rule after seeing the result. `answered` counts predictions, `correct` counts predictions matching the human label, `known` counts test questions truly about time or place, `refused` counts refusals, and `unknown_refused` counts appropriate refusals on other topics. **Answer precision** is `correct / answered`; **recall of known questions** is `correct / known`; **refusal rate** is `refused / total`. `/` divides numbers. `round(x, 2)` rounds to two decimal places. A ratio with a zero denominator is undefined, so the code checks before division.

Copy [`examples.json`](https://github.com/DauletBai/shanraq.org/blob/main/course/kazakh-ai/step-25/examples.json) beside the program. Its 19 cards use the format explained in lesson 22.

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-25
python3 evaluate.py
```

On Windows, use `py evaluate.py`.

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
answered = 0
correct = 0
known = 0
refused = 0
unknown_refused = 0
for row in rows:
    if row["split"] == "test":
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
        if row["label"] != "белгісіз":
            known += 1
        if prediction == "білмеймін":
            refused += 1
            if row["label"] == "белгісіз":
                unknown_refused += 1
        else:
            answered += 1
            if prediction == row["label"]:
                correct += 1
        print(row["text"], "→", prediction, "|", row["label"])
print("Жауап:", answered, "Дұрыс:", correct, "Белгілі:", known)
print("Бас тарту:", refused, "Белгісізге дұрыс бас тарту:", unknown_refused)
if answered == 0:
    print("Дәлдік: анықталмайды")
else:
    print("Дәлдік:", round(correct / answered, 2))
if known == 0:
    print("Толықтық: анықталмайды")
else:
    print("Толықтық:", round(correct / known, 2))
total = answered + refused
if total == 0:
    print("Бас тарту үлесі: анықталмайды")
else:
    print("Бас тарту үлесі:", round(refused / total, 2))
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-25).

## How the code works

The model predicts a type for two of six questions and gets both right: type precision 1.0. It covers only two of four known-type questions: recall 0.5. Four refusals to choose a type out of six round to 0.67; two are appropriate for other topics, while two expose missing forms (`қай күні`, `қай жерде`). This tiny set cannot estimate real-world quality reliably. Precision here concerns the *question type*, not the truth of a fact. After changing the rule, use a fresh unseen test set.

## Support map

Held-out test → answer/refusal per question → 2/2 correct answers, 2/4 known covered, 4/6 refusals → missing forms.

## Recall and check

Hide the code and reconstruct all three fractions with their denominators. Find the two known-label refusals: `Шахмат қай күні?` and `Сурет қай жерде?`. Hint: they lower recall but leave answered-question precision unchanged. Expected: `Дәлдік: 1.0`, `Толықтық: 0.5`, `Бас тарту үлесі: 0.67`. Common mistake: reading 1.0 as “the model is always right”; it answered only one third of the questions, and this test did not check any timetable facts.
