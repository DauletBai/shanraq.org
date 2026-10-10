# Bias, fairness, and privacy in AI

_Lead (summary):_ **The training rows contain Kazakh and English words, but not all subjects or languages.**

## Where we are on the map

Lesson 67 of 72. AI and project release block (63–72). The assistant keeps its earlier reminder rules and local database; the new feature only suggests a task category and can abstain.

![Bias, fairness, and privacy in AI](/static/course/informatics/map-67-bias-fairness-privacy-en.svg)

## Situation and question

The training rows contain Kazakh and English words, but not all subjects or languages. Claiming equal quality for every learner without checking could leave some people with worse suggestions.

## New words without gaps

**Bias** is a systematic mismatch between data or results and the intended population or task. **Fairness** asks us to examine effects across groups rather than one overall percentage. **Privacy** limits use of personal data; this toy uses fictional rows. A **quality slice** evaluates one chosen group separately. Very small slices give unstable results.

## The lesson's support signal

Read the three cells left to right, then cover the third. State the input and action in each cell; predict how the path ends and which observable fact checks it. Compare the drawing with a project file or worked example. If a number on the map cannot be derived from the data, correct the map rather than changing data for a pretty picture.

## Work through it step by step

Split the four test rows by language: two Kazakh and two English. Count mistakes separately but always state the denominator of two for each group. Invent an example in a language absent from training: unfamiliar words should return `review`. Adding real learners’ messages will not automatically fix the gap and introduces privacy risk. Prefer permitted fictional variants and test them separately.

## Predict and check

Write the expected output before running code. Execute it in `step-06` or `step-07`, compare every character and explain each line. Then change one safe fictional input and repeat “predict — run — change — explain — verify.” Do not turn an observation on a tiny dataset into a promise about real people.

```python
from classifier import read_examples, train, suggest
rows = read_examples("labelled_tasks.csv")
model = train(rows)
print([suggest(r["title"], model) for r in rows if r["split"] == "test"])
```

## Expected output

```text
['reading', 'reading', 'math', 'math']
```

## Catch the error

Scores of 2/2 and 2/2 in tiny samples do not establish equal reliability. Do not infer a learner’s language from their name. Do not publish private messages to improve the model.

## Project change

Release 2.1 preserves the DB, report, reminders and backup. `labelled_tasks.csv` contains fictional examples only. `classifier.py` suggests `math`, `reading` or `review` but writes no decision to the database. Check each model output against held-out examples and let a person reconsider it.

## Task and evidence

Make two slices and state their denominators. Design a safe new example set without personal data. Say who should resolve a disputed label.

## Transfer to a new setting

An automatic translator works well on one dialect. What should be measured before using it for all parents at a school?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name one boundary of the idea. After 7 days explain it to a friend without reading the lesson and answer 7 of 10 block questions; analyse mistakes with concrete inputs. After 30 days open a fresh project copy, reproduce the check and transfer the idea to a new case. Keep a record of input, expected output and actual output.

## Primary reference to check

[Official reference](https://www.nist.gov/itl/ai-risk-management-framework)

## Next lesson

[How a language model continues text](/read/informatics-68-generative-ai-llm?lang=en)
