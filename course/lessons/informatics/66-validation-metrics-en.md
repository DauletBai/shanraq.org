# Model validation and a confusion matrix

_Lead (summary):_ **The classifier gets four out of four cards right.**

## Where we are on the map

Lesson 66 of 72. AI and project release block (63–72). The assistant keeps its earlier reminder rules and local database; the new feature only suggests a task category and can abstain.

![Model validation and a confusion matrix](/static/course/informatics/map-66-validation-metrics-en.svg)

## Situation and question

The classifier gets four out of four cards right. It sounds perfect, but there are only four cards. A pleasing percentage can hide an error, so display the individual outcomes.

## New words without gaps

A **test set** consists of examples excluded from training. A **confusion matrix** counts “true label → suggestion” pairs. **Accuracy** is 4/4 on this tiny sample, not a guarantee for tomorrow. A **false classification** assigns the wrong category and deserves separate examination. **Uncertainty** is substantial with so few examples.

## The lesson's support signal

Read the three cells left to right, then cover the third. State the input and action in each cell; predict how the path ends and which observable fact checks it. Compare the drawing with a project file or worked example. If a number on the map cannot be derived from the data, correct the map rather than changing data for a pretty picture.

## Work through it step by step

Run `classifier.py`: `held-out` contains two correct reading and two correct math predictions, four total. Draw a 2×2 table with true labels as rows and suggestions as columns. In this sample the diagonal is 2 and 2 and other cells are zero. Invent a new title with words from both categories; it may yield `review` or a mistake. Do not put it in train until a fair new evaluation is finished.

## Predict and check

Write the expected output before running code. Execute it in `step-06` or `step-07`, compare every character and explain each line. Then change one safe fictional input and repeat “predict — run — change — explain — verify.” Do not turn an observation on a tiny dataset into a promise about real people.

```python
from classifier import read_examples, train, evaluate
rows = read_examples("labelled_tasks.csv")
matrix = evaluate(rows, train(rows))
print(matrix[("reading", "reading")], matrix[("math", "math")], sum(matrix.values()))
```

## Expected output

```text
2 2 4
```

## Catch the error

4/4 is not proven 100% reliability. If you tune the vocabulary after seeing test rows, those rows cease to be independent evaluation. Count `review` separately rather than hiding it as success.

## Project change

Release 2.1 preserves the DB, report, reminders and backup. `labelled_tasks.csv` contains fictional examples only. `classifier.py` suggests `math`, `reading` or `review` but writes no decision to the database. Check each model output against held-out examples and let a person reconsider it.

## Task and evidence

Show the matrix and the separate `review` count. Add one new held-out example and recalculate the fraction correct by hand. Name the error you consider most consequential.

## Transfer to a new setting

A medical test works perfectly for four people. Can that result justify a promise of error-free testing for a city?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name one boundary of the idea. After 7 days explain it to a friend without reading the lesson and answer 7 of 10 block questions; analyse mistakes with concrete inputs. After 30 days open a fresh project copy, reproduce the check and transfer the idea to a new case. Keep a record of input, expected output and actual output.

## Primary reference to check

[Official reference](https://scikit-learn.org/stable/modules/model_evaluation.html#confusion-matrix)

## Next lesson

[Bias, fairness, and privacy in AI](/read/informatics-67-bias-fairness-privacy?lang=en)
