# Project defence: demonstrate a solution, an error, and a limitation

_Lead (summary):_ **At the board you show a working solution, not merely attractive slides.**

## Where we are on the map

Lesson 72 of 72. AI and project release block (63–72). The assistant keeps its earlier reminder rules and local database; the new feature only suggests a task category and can abstain.

![Project defence: demonstrate a solution, an error, and a limitation](/static/course/informatics/map-72-capstone-defense-en.svg)

## Situation and question

At the board you show a working solution, not merely attractive slides. A viewer can give an unexpected task title, damage a copy and ask what the assistant cannot do.

## New words without gaps

A **project defence** demonstrates a question, solution, evidence and limit of use. **Reproducibility** means another person obtains the same result from instructions. A **counterexample** is an input on which a claim fails. **Success criteria** are set before the demo: the task opens, the report is correct, recovery works and an unfamiliar title goes to a person.

## The lesson's support signal

Read the three cells left to right, then cover the third. State the input and action in each cell; predict how the path ends and which observable fact checks it. Compare the drawing with a project file or worked example. If a number on the map cannot be derived from the data, correct the map rather than changing data for a pretty picture.

## Work through it step by step

Show four episodes: 1) three fictional tasks and their status on 2026-10-09; 2) the 45/0/25 report, explaining the missing session for `t-02`; 3) backup and restore; 4) `reading`, `math` and `review` for three titles. Let a classmate choose a new unknown title and decide its category personally. Then name the missing accounts, encrypted backup, public server and proof of model accuracy.

## Predict and check

Write the expected output before running code. Execute it in `step-06` or `step-07`, compare every character and explain each line. Then change one safe fictional input and repeat “predict — run — change — explain — verify.” Do not turn an observation on a tiny dataset into a promise about real people.

```python
from classifier import read_examples, train, suggest
model = train(read_examples("labelled_tasks.csv"))
print(*(suggest(title, model) for title in ("Кітап оқу", "Геометрия есебі", "unknown")))
```

## Expected output

```text
reading math review
```

## Catch the error

Do not call the teaching project a ready system for real children. Do not conceal abstentions or empty results. Do not promise that 4/4 on a tiny sample will hold on new data.

## Project change

Release 2.1 preserves the DB, report, reminders and backup. `labelled_tasks.csv` contains fictional examples only. `classifier.py` suggests `math`, `reading` or `review` but writes no decision to the database. Check each model output against held-out examples and let a person reconsider it.

## Task and evidence

Prepare a five-minute demonstration and evidence card. Answer ten block questions, correct each mistake, repeat at least 7/10 after a week and restore the project in a fresh folder after 30 days. Ask an independent learner to reproduce the outcome.

## Ten questions for self-check

1. How does a rule differ from a model?
2. What is a feature?
3. What is a label?
4. Why is test excluded from training?
5. What matrix came from four rows?
6. Why is 4/4 no guarantee?
7. What does review mean?
8. Who makes the final decision?
9. How do you verify an AI fact?
10. What is the boundary of 2.1?

## Transfer to a new setting

How would your defence change if a client requested a real multiuser app containing personal data?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name one boundary of the idea. After 7 days explain it to a friend without reading the lesson and answer 7 of 10 block questions; analyse mistakes with concrete inputs. After 30 days open a fresh project copy, reproduce the check and transfer the idea to a new case. Keep a record of input, expected output and actual output.

## Primary reference to check

[Official reference](https://www.nist.gov/itl/ai-risk-management-framework)

## All course lessons

[All lessons](/course/informatics?lang=en)
