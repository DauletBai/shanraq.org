# Lesson 29. Audit locality and explain each answer

## Why this matters

A food label tells us where an item came from, when it was checked, and when it expires. A fact card needs similar information. For each answer, record a decision trail: what was recognized, how many cards matched, and why an answer was allowed.

## Before the code

Copy `engine.py`, `checks.py`, `examples.json`, and `facts.json` together. `from checks import classify` loads the lesson 27 function; its demo does not run on import because of the `__main__` guard. `answer(question, on_date)` takes question text and a calendar date, then returns a dictionary with `text` and `trace`. `trace` is a list of short decisions, `str(...)` converts a number to text, and `+` joins strings. `on_date.isoformat()` writes a `year-month-day` date. We check question type, then one fact card, check date, and expiry. The chess fact is fictional.

Files in this step: `examples.json`, `checks.py`, `engine.py`, `facts.json`.

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-29
python3 engine.py
```

On Windows, replace `python3` with `py`.

```python
import json
from datetime import date

from checks import classify


def answer(question, on_date):
    trace = []
    request = classify(question)
    if request["status"] != "ready":
        trace.append("Сұрақ түрі анықталмады")
        return {"text": request["reason"], "trace": trace}
    club = request["club"]
    kind = request["kind"]
    trace.append("Кілт: " + club + " / " + kind)
    with open("facts.json", encoding="utf-8") as file:
        facts = json.load(file)
    matches = []
    for fact in facts:
        if fact["club"] == club and fact["kind"] == kind:
            matches.append(fact)
    trace.append("Жазба саны: " + str(len(matches)))
    if len(matches) == 0:
        return {"text": "білмеймін: дерек жоқ", "trace": trace}
    if len(matches) > 1:
        return {"text": "тоқта: бірнеше дерек табылды", "trace": trace}
    fact = matches[0]
    checked = date.fromisoformat(fact["checked_on"])
    until = date.fromisoformat(fact["valid_until"])
    trace.append("Дереккөз: " + fact["source"])
    trace.append("Тексерілген: " + fact["checked_on"])
    trace.append("Жарамды: " + fact["valid_until"] + " дейін")
    if on_date < checked or on_date > until:
        return {"text": "білмеймін: бұл күнге жарамды дерек жоқ", "trace": trace}
    if kind == "уақыт":
        label = "уақыты"
    else:
        label = "орны"
    message = (fact["club"] + " үйірмесінің " + label + ": " + fact["value"]
               + " | " + on_date.isoformat() + " күнгі дерек"
               + " | дереккөз: " + fact["source"])
    return {"text": message, "trace": trace}


if __name__ == "__main__":
    day = date.fromisoformat("2026-09-24")
    for question in ["Шахмат қашан?", "Сурет қайда?", "Шахмаат қашан?"]:
        result = answer(question, day)
        print(question, "→", result["text"])
        print("Із:", result["trace"])
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-29).

## How the program works

For `Шахмат қашан?` on 24 September 2026, the trail shows the key, one card, source, and two dates. `Сурет қайда?` has no card; the typo `Шахмаат` stops before lookup. A day before `checked_on` or after `valid_until` yields no answer. `<` means earlier/less, and `>` means later/greater. The program reads local JSON files and contains no network API call. That audits the visible code path, not isolation of the whole operating system. The trail can be inspected, but an incorrect source card can still produce an incorrect answer.

## Support map

Question → type → key → 0/1/many cards → two dates → sourced answer and trail or refusal.

![Lesson 29 support map](/static/course/kazakh-ai/map-29-audit-en.svg)

## Recall and check

Hide the code and recall the order of checks. Expected: chess yields an answer and five trail steps; art has no card and yields `білмеймін: дерек жоқ`; the typo yields `білмеймін: таныс емес сөз`. Hint: two early `return` statements stop before date checks. Change `valid_until` to `2026-09-01` and expect `білмеймін: бұл күнге жарамды дерек жоқ`. Common mistake: treating a trace as proof that the original sheet was true.
