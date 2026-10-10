# Training a simple classifier

_Lead (summary):_ **You teach a younger friend to sort cards: “read” goes with books; “geometry” with maths.**

## Where we are on the map

Lesson 65 of 72. AI and project release block (63–72). The assistant keeps its earlier reminder rules and local database; the new feature only suggests a task category and can abstain.

![Training a simple classifier](/static/course/informatics/map-65-train-classifier-en.svg)

## Situation and question

You teach a younger friend to sort cards: “read” goes with books; “geometry” with maths. Our code performs a very simple version by counting words in labelled teaching cards.

## New words without gaps

**Training** here means counting word occurrences by train labels, not understanding a subject. **Tokenisation** extracts words from a title. A **word weight** counts how many training examples of one category contain it. A **category score** sums these weights for a new title. A **tie**, including no known words, returns `review`.

## The lesson's support signal

Read the three cells left to right, then cover the third. State the input and action in each cell; predict how the path ends and which observable fact checks it. Compare the drawing with a project file or worked example. If a number on the map cannot be derived from the data, correct the map rather than changing data for a pretty picture.

## Work through it step by step

In `step-07`, run `python3 classifier.py "Кітап оқу"`: expect `suggestion: reading`. For `Геометрия есебі`, expect `math`. Open `classifier.py`: `train` reads only `split=train`; `suggest` sums counters. For unfamiliar “Task without examples”, both scores are zero and the answer is `review`. Change one fictional training row and predict which score will move.

## Predict and check

Write the expected output before running code. Execute it in `step-06` or `step-07`, compare every character and explain each line. Then change one safe fictional input and repeat “predict — run — change — explain — verify.” Do not turn an observation on a tiny dataset into a promise about real people.

```python
from classifier import read_examples, train, suggest
model = train(read_examples("labelled_tasks.csv"))
print(suggest("Кітап оқу", model), suggest("Геометрия есебі", model), suggest("unknown", model))
```

## Expected output

```text
reading math review
```

## Catch the error

This is a toy classifier, not a large language model. It does not understand inflections, negation or context. Moving a convenient test row into training before evaluation breaks the evaluation.

## Project change

Release 2.1 preserves the DB, report, reminders and backup. `labelled_tasks.csv` contains fictional examples only. `classifier.py` suggests `math`, `reading` or `review` but writes no decision to the database. Check each model output against held-out examples and let a person reconsider it.

## Task and evidence

Compute scores by hand for titles with one familiar word. Show a tie. Then run tests and explain any discrepancy.

## Transfer to a new setting

A mail sorter knows “bill” but encounters “billing”. Why might the word-splitting rule change its answer?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name one boundary of the idea. After 7 days explain it to a friend without reading the lesson and answer 7 of 10 block questions; analyse mistakes with concrete inputs. After 30 days open a fresh project copy, reproduce the check and transfer the idea to a new case. Keep a record of input, expected output and actual output.

## Primary reference to check

[Official reference](https://docs.python.org/3/library/collections.html#collections.Counter)

## Next lesson

[Model validation and a confusion matrix](/read/informatics-66-validation-metrics?lang=en)
