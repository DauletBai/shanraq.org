# Release 2.1: tests, instructions, and product handoff

_Lead (summary):_ **A classmate downloads the project and asks “What runs first? What if the DB exists? When can I trust a model suggestion?” A release answers with runnable instructions and checks.**

## Where we are on the map

Lesson 71 of 72. AI and project release block (63–72). The assistant keeps its earlier reminder rules and local database; the new feature only suggests a task category and can abstain.

![Release 2.1: tests, instructions, and product handoff](/static/course/informatics/map-71-release-testing-docs-en.svg)

## Situation and question

A classmate downloads the project and asks “What runs first? What if the DB exists? When can I trust a model suggestion?” A release answers with runnable instructions and checks.

## New words without gaps

**Version 2.1** is a fixed set of files and behaviour. A **regression test** checks that a new feature did not break an older rule. **Documentation** lists inputs, commands, results and limits. **Product handoff** means another person can reproduce the path without the author present.

## The lesson's support signal

Read the three cells left to right, then cover the third. State the input and action in each cell; predict how the path ends and which observable fact checks it. Compare the drawing with a project file or worked example. If a number on the map cannot be derived from the data, correct the map rather than changing data for a pretty picture.

## Work through it step by step

In `step-07`, run the tests: they cover earlier reminders and SQLite, backup and recovery, and the separate classifier. Create the DB only with `init`; run the 45/0/25 report; check `classifier.py "Кітап оқу"` and an unfamiliar phrase. Set the 4/4 held-out score beside a clear warning that the sample is tiny. Open the README for your language and give it to a classmate.

## Predict and check

Write the expected output before running code. Execute it in `step-06` or `step-07`, compare every character and explain each line. Then change one safe fictional input and repeat “predict — run — change — explain — verify.” Do not turn an observation on a tiny dataset into a promise about real people.

```python
from classifier import read_examples, train, evaluate
rows = read_examples("labelled_tasks.csv")
print(len(rows), sum(evaluate(rows, train(rows)).values()))
```

## Expected output

```text
16 4
```

## Catch the error

Passing tests do not prove the absence of all defects. Four correct predictions do not establish performance on real learners. The web server remains local and educational; the backup remains unencrypted.

## Project change

Release 2.1 preserves the DB, report, reminders and backup. `labelled_tasks.csv` contains fictional examples only. `classifier.py` suggests `math`, `reading` or `review` but writes no decision to the database. Check each model output against held-out examples and let a person reconsider it.

## Task and evidence

Write a short handoff protocol: commands, exact output, backup location, security limits and a way to report an error. Ask another learner to follow it without spoken hints.

## Transfer to a new setting

A club wants to roll the app out across a school. Which tests and decisions are needed beyond our local classroom version?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name one boundary of the idea. After 7 days explain it to a friend without reading the lesson and answer 7 of 10 block questions; analyse mistakes with concrete inputs. After 30 days open a fresh project copy, reproduce the check and transfer the idea to a new case. Keep a record of input, expected output and actual output.

## Primary reference to check

[Official reference](https://docs.python.org/3/library/unittest.html)

## Next lesson

[Project defence: demonstrate a solution, an error, and a limitation](/read/informatics-72-capstone-defense?lang=en)
