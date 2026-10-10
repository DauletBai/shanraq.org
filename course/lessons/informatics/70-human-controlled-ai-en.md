# Adding AI while keeping the human in control

_Lead (summary):_ **The model suggests `math`, but the learner titles a task “Read the history of mathematics.**

## Where we are on the map

Lesson 70 of 72. AI and project release block (63–72). The assistant keeps its earlier reminder rules and local database; the new feature only suggests a task category and can abstain.

![Adding AI while keeping the human in control](/static/course/informatics/map-70-human-controlled-ai-en.svg)

## Situation and question

The model suggests `math`, but the learner titles a task “Read the history of mathematics.” A person sees both topics and may choose reading. The assistant must not silently edit the card.

## New words without gaps

A **human in the decision loop** sees a suggestion, may correct it and has clear responsibility. The **abstention threshold** here is simple: tied or zero scores return `review`. An **explanation** shows why a suggestion arose, such as matching words and counts. An **automatic action** would change data without confirmation; this teaching version has none.

## The lesson's support signal

Read the three cells left to right, then cover the third. State the input and action in each cell; predict how the path ends and which observable fact checks it. Compare the drawing with a project file or worked example. If a number on the map cannot be derived from the data, correct the map rather than changing data for a pretty picture.

## Work through it step by step

Run `classifier.py "Кітап оқу"`: get `reading`; an unknown title gives `review`. Inspect the original task card: calling the classifier changes neither `title`, `due_date`, `done` nor the reminder. A learner may accept, correct or skip advice on a paper decision form. Record the suggested label and your choice without a learner’s name.

## Predict and check

Write the expected output before running code. Execute it in `step-06` or `step-07`, compare every character and explain each line. Then change one safe fictional input and repeat “predict — run — change — explain — verify.” Do not turn an observation on a tiny dataset into a promise about real people.

```python
from pathlib import Path
from tempfile import TemporaryDirectory
from data_store import create_database
from classifier import read_examples, train, suggest
with TemporaryDirectory() as folder:
    db = Path(folder) / "example.db"
    create_database(db, "tasks.json", "study_sessions.csv")
    before = db.read_bytes()
    print(suggest("Кітап оқу", train(read_examples("labelled_tasks.csv"))))
    print(before == db.read_bytes())
```

## Expected output

```text
reading
True
```

## Catch the error

A hidden automatic choice after `review` would defeat human control. Free text from a generative model cannot be called a verified recommendation without sources. Agreeing to one category does not permit collecting other private data.

## Project change

Release 2.1 preserves the DB, report, reminders and backup. `labelled_tasks.csv` contains fictional examples only. `classifier.py` suggests `math`, `reading` or `review` but writes no decision to the database. Check each model output against held-out examples and let a person reconsider it.

## Task and evidence

For three titles show the suggestion, its explanation and a person’s choice. Intentionally reject the model once. Verify that tasks and statuses in the DB remain unchanged.

## Transfer to a new setting

A system proposes a school bus route. Which decisions may be suggested to a driver and which should not be executed without confirmation?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name one boundary of the idea. After 7 days explain it to a friend without reading the lesson and answer 7 of 10 block questions; analyse mistakes with concrete inputs. After 30 days open a fresh project copy, reproduce the check and transfer the idea to a new case. Keep a record of input, expected output and actual output.

## Primary reference to check

[Official reference](https://www.nist.gov/itl/ai-risk-management-framework)

## Next lesson

[Release 2.1: tests, instructions, and product handoff](/read/informatics-71-release-testing-docs?lang=en)
