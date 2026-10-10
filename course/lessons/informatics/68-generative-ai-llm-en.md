# How a language model continues text

_Lead (summary):_ **A phone suggests the next word.**

## Where we are on the map

Lesson 68 of 72. AI and project release block (63–72). The assistant keeps its earlier reminder rules and local database; the new feature only suggests a task category and can abstain.

![How a language model continues text](/static/course/informatics/map-68-generative-ai-llm-en.svg)

## Situation and question

A phone suggests the next word. A large language model does a far more complex version over a longer context, yet a fluent sentence can still be wrong.

## New words without gaps

A **token** is a piece of text processed by a model; it is not always a whole word. **Context** is the part of the conversation and instructions available to the model. **Training** on large corpora fits parameters for predicting continuation. **Generation** selects further tokens in steps. A **prompt** is user input. The model does not automatically gain a verified source of truth.

## The lesson's support signal

Read the three cells left to right, then cover the third. State the input and action in each cell; predict how the path ends and which observable fact checks it. Compare the drawing with a project file or worked example. If a number on the map cannot be derived from the data, correct the map rather than changing data for a pretty picture.

## Work through it step by step

Write two completions of “Tomorrow we have…”: “class” and “a holiday”. Both are grammatical; the timetable decides which is true. Compare with our `step-07` classifier: it returns only `math`, `reading` or `review`, while a generative model creates free text. Do not send someone else’s tasks or details to an outside service without permission.

## Predict and check

Write the expected output before running code. Execute it in `step-06` or `step-07`, compare every character and explain each line. Then change one safe fictional input and repeat “predict — run — change — explain — verify.” Do not turn an observation on a tiny dataset into a promise about real people.

```python
from collections import Counter
next_words = Counter(("class", "holiday"))
print(next_words["class"], next_words["holiday"])
```

## Expected output

```text
1 1
```

## Catch the error

A polite, coherent answer is no guarantee of fact. “Large” does not grant a model authority to decide for a person. Generating text and locating a current official document are different operations.

## Project change

Release 2.1 preserves the DB, report, reminders and backup. `labelled_tasks.csv` contains fictional examples only. `classifier.py` suggests `math`, `reading` or `review` but writes no decision to the database. Check each model output against held-out examples and let a person reconsider it.

## Task and evidence

For three statements separate “plausible” from “verified.” Name the document or observation needed for each check.

## Transfer to a new setting

A model gives a convincing date for a school competition. What should a family verify before buying tickets?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name one boundary of the idea. After 7 days explain it to a friend without reading the lesson and answer 7 of 10 block questions; analyse mistakes with concrete inputs. After 30 days open a fresh project copy, reproduce the check and transfer the idea to a new case. Keep a record of input, expected output and actual output.

## Primary reference to check

[Official reference](https://arxiv.org/abs/1706.03762)

## Next lesson

[A confident answer is not a reliable answer](/read/informatics-69-ai-verification-sources?lang=en)
