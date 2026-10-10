# Rules, algorithms, and trainable models

_Lead (summary):_ **The assistant already applies an exact rule to remind us of due dates.**

## Where we are on the map

Lesson 63 of 72. AI and project release block (63–72). The assistant keeps its earlier reminder rules and local database; the new feature only suggests a task category and can abstain.

![Rules, algorithms, and trainable models](/static/course/informatics/map-63-rules-algorithms-models-en.svg)

## Situation and question

The assistant already applies an exact rule to remind us of due dates. Now it may suggest whether a new task concerns reading or maths. The due-date rule stays unchanged: a category suggestion never controls reminders.

## Where to get the project files

Open the [checkpoint folder](https://github.com/DauletBai/shanraq.org/tree/main/course/informatics-assistant/step-07), download the repository with Code → Download ZIP and find `course/informatics-assistant/step-07`. Open a terminal there. `python3 --version` shows Python; then run the lesson check. All tasks are fictional; enter no real personal data.

## New words without gaps

An **algorithm** is a finite sequence of steps. A **rule** is specified by a person in advance, such as “remind when at most two days remain.” A **trainable model** derives parameters from examples; here they are word counts. A **prediction** is a suggestion for a new record, not proof. `review` means the evidence cannot distinguish categories.

## The lesson's support signal

Read the three cells left to right, then cover the third. State the input and action in each cell; predict how the path ends and which observable fact checks it. Compare the drawing with a project file or worked example. If a number on the map cannot be derived from the data, correct the map rather than changing data for a pretty picture.

## Work through it step by step

Compare two paths on one card. For `t-01`, a precise due date yields `REMIND` under the 1.0 rule. The title `Кітап оқу` may yield `reading` through the model’s word counts. Change the title to unknown words: the model should return `review` while the due-date rule still works. Model quality and reminder correctness need separate checks.

## Predict and check

Write the expected output before running code. Execute it in `step-06` or `step-07`, compare every character and explain each line. Then change one safe fictional input and repeat “predict — run — change — explain — verify.” Do not turn an observation on a tiny dataset into a promise about real people.

```python
from datetime import date
from assistant_core import reminder_status
from classifier import read_examples, train, suggest
rows = read_examples("labelled_tasks.csv")
print(reminder_status({"done": False, "due_date": "2026-10-10"}, date(2026,10,9)))
print(suggest("Кітап оқу", train(rows)))
```

## Expected output

```text
REMIND
reading
```

## Catch the error

Code does not become AI merely by using `if`. A trained model does not understand a word as a person does. Do not replace an exact due-date rule with a statistical guess.

## Project change

Release 2.1 preserves the DB, report, reminders and backup. `labelled_tasks.csv` contains fictional examples only. `classifier.py` suggests `math`, `reading` or `review` but writes no decision to the database. Check each model output against held-out examples and let a person reconsider it.

## Task and evidence

Draw two tracks for one task: due date → status and title → suggested category. Name the input, output and person responsible for the final choice.

## Transfer to a new setting

A thermostat starts heating below a chosen temperature, while a system forecasts tomorrow’s weather. Which path uses a fixed threshold?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name one boundary of the idea. After 7 days explain it to a friend without reading the lesson and answer 7 of 10 block questions; analyse mistakes with concrete inputs. After 30 days open a fresh project copy, reproduce the check and transfer the idea to a new case. Keep a record of input, expected output and actual output.

## Primary reference to check

[Official reference](https://www.nist.gov/itl/ai-risk-management-framework)

## Next lesson

[Features, labels, and an honest dataset](/read/informatics-64-features-labels-datasets?lang=en)
