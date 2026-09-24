# Lesson 30. Release the command-line model and check our claims

## Why this matters

A first-aid kit is useful when each item has a name, expiry, and instructions. Our release is also a complete kit: program, rules, exercise cards, a run command, and checkable boundaries. It answers only two kinds of short question about a fictional club.

## Before the code

Copy **all five files** from step 30: `main.py`, `engine.py`, `checks.py`, `facts.json`, and `examples.json`. The command takes a date and question; quotes keep a spaced question as one argument. Put optional `--trace` **before the date** to print the decision trail. `sys.argv` stores arguments; `len` selects one of two accepted command forms. `sys.exit(2)` stops bad usage with an error code. `try` attempts to parse the date, while `except ValueError` catches a bad format and prints a clear message. The date is the **day being asked about**, not a claim that an exercise record is true now.

Files in this step: `examples.json`, `checks.py`, `engine.py`, `facts.json`, `main.py`.

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-30
python3 main.py 2026-09-24 "Шахмат қашан?"
```

On Windows, replace `python3` with `py`.

```python
import sys
from datetime import date

from engine import answer

if len(sys.argv) == 3:
    show_trace = False
    raw_date = sys.argv[1]
    question = sys.argv[2]
elif len(sys.argv) == 4 and sys.argv[1] == "--trace":
    show_trace = True
    raw_date = sys.argv[2]
    question = sys.argv[3]
else:
    print('Қолдану: python3 main.py [--trace] YYYY-MM-DD "Сұрақ"')
    sys.exit(2)
try:
    on_date = date.fromisoformat(raw_date)
except ValueError:
    print("Қате күн: YYYY-MM-DD түрінде жазыңыз")
    sys.exit(2)
result = answer(question, on_date)
print(result["text"])
if show_trace:
    for line in result["trace"]:
        print("Із:", line)
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-30).

## How the program works

From the step folder, run `python3 main.py 2026-09-24 "Шахмат қашан?"`. It prints `шахмат үйірмесінің уақыты: бейсенбі, 15:00 | 2026-09-24 күнгі дерек | дереккөз: club-sheet-01`. `Сурет қайда?` yields `білмеймін: дерек жоқ`; `Шахмат қашан қайда?` asks for one type; date `2027-01-01` refuses an expired record. Add `--trace` to see five source-check steps. The small held-out set in lesson 25 tested question types only; it does not establish real service quality. We built a narrow local model, not universal AI. No paid API call in exercise code does not mean zero total cost; avoiding free-form generation does not eliminate source errors.

## Support map

Command + date + question → word guard → type → one card → valid date → sourced answer; otherwise an explained refusal.

## Recall and check

Hide the code and reconstruct the support map. Try four commands: a known question, unknown club, expired date, and `--trace`. Hint: `2026-09-24` is between `checked_on` and `valid_until`, while `2027-01-01` is not. Expected for the expired date: `білмеймін: бұл күнге жарамды дерек жоқ`; for an unknown club: `білмеймін: таныс емес сөз`. **Final defense:** show how the answer depends on a card and date, name two justified refusals, explain why lesson 25’s 1.0 is not perfection, and why “a thousand times” needs a comparable experiment. Common mistake: presenting an exercise date as a real club’s current schedule.
