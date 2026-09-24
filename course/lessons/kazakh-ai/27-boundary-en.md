# Lesson 27. Check misspellings and unrelated topics

## Why this matters

A receptionist checks a visitor’s name against a list. If the name is unclear or two passes are requested at once, the receptionist asks rather than guessing. Our classifier should likewise stop on an unknown word, a typo, or a question asking for both time and place.

## Before the code

`examples.json` contains the 13 open training and tuning cards from lesson 24; test labels stay out. The code learns the same counts and adds an input guard. Braces without colons make a **set** named `allowed`: it holds permitted words, while `in` and `not in` test membership. `def classify(question):` defines a reusable function with `question` as its input. `return` ends that function and returns a result dictionary. `if __name__ == "__main__":` runs five examples only when this file is executed directly; importing `classify` from another file does not print them. These are Python’s special names; leave them unchanged.

Files in this step: `examples.json`, `checks.py`.

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-27
python3 checks.py
```

On Windows, replace `python3` with `py`.

```python
import json

with open("examples.json", encoding="utf-8") as file:
    rows = json.load(file)
weights = {}
for row in rows:
    if row["split"] == "train":
        words = row["text"].lower().replace("?", "").split()
        for word in words:
            for feature in [word, word[:4]]:
                if feature not in weights:
                    weights[feature] = {"уақыт": 0, "орын": 0}
                weights[feature][row["label"]] += 1


def classify(question):
    text = question.lower().replace("?", "").replace(",", "").replace(".", "")
    words = text.split()
    allowed = {"шахмат", "сурет", "қашан", "қайда", "уақыты", "орны", "өтеді", "болады"}
    for word in words:
        if word not in allowed:
            return {"status": "refuse", "reason": "білмеймін: таныс емес сөз"}
    clubs = []
    for word in words:
        if word in ["шахмат", "сурет"]:
            clubs.append(word)
    if len(clubs) != 1:
        return {"status": "refuse", "reason": "нақтылаңыз: бір үйірме керек"}
    scores = {"уақыт": 0, "орын": 0}
    for word in words:
        for feature in [word, word[:4]]:
            if feature in weights:
                scores["уақыт"] += weights[feature]["уақыт"]
                scores["орын"] += weights[feature]["орын"]
    time_cue = "қашан" in words or "уақыты" in words
    place_cue = "қайда" in words or "орны" in words
    if time_cue and not place_cue and scores["уақыт"] > scores["орын"]:
        return {"status": "ready", "club": clubs[0], "kind": "уақыт"}
    if place_cue and not time_cue and scores["орын"] > scores["уақыт"]:
        return {"status": "ready", "club": clubs[0], "kind": "орын"}
    return {"status": "refuse", "reason": "нақтылаңыз: бір сұрақ түрі керек"}


if __name__ == "__main__":
    for question in ["Шахмат қашан?", "Шахмаат қашан?", "Шахмат неге?",
                     "Шахмат қашан интернет?", "Шахмат қашан қайда?"]:
        print(question, "→", classify(question))
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-27).

## How the program works

Input is lowercased and `?`, `,`, `.` are removed. Every word is then checked against a closed list. `Шахмаат` is not silently corrected to `шахмат`; it yields `білмеймін: таныс емес сөз`. An unrelated topic also stops. Next come the one-club check, learned counts, and explicit cues: two question types request clarification. This strict list will also refuse some phrases a person can understand; that is an honest limit of a tiny model. `ready` means only “question type identified,” not “fact found.”

## Support map

Words → all allowed? → exactly one club? → one clear scored type? → ready; otherwise a refusal reason.

## Recall and check

Hide the code and recall the checks in order. Expected: `Шахмат қашан?` is `ready`; `Шахмаат қашан?` and `Шахмат қашан интернет?` refuse an unfamiliar word; `Шахмат қашан қайда?` asks for one type. Hint: the first failed check returns and ends the function. Common mistake: treating `ready` as a schedule answer; catalog lookup comes later.
