# A confident answer is not a reliable answer

_Lead (summary):_ **The assistant confidently says “the library is closed tomorrow.**

## Where we are on the map

Lesson 69 of 72. AI and project release block (63–72). The assistant keeps its earlier reminder rules and local database; the new feature only suggests a task category and can abstain.

![A confident answer is not a reliable answer](/static/course/informatics/map-69-ai-verification-sources-en.svg)

## Situation and question

The assistant confidently says “the library is closed tomorrow.” The library may have changed its hours an hour ago. Without a source and check time, the sentence is a poor basis for a decision.

## New words without gaps

A **source** lets you independently check a claim. A **primary source** publishes information from the responsible organisation or participant. An **update date** shows freshness. A **testable claim** can be compared with a document, measurement or record. An AI **hallucination** is plausible but unsupported output; the term does not mean human perception.

## The lesson's support signal

Read the three cells left to right, then cover the third. State the input and action in each cell; predict how the path ends and which observable fact checks it. Compare the drawing with a project file or worked example. If a number on the map cannot be derived from the data, correct the map rather than changing data for a pretty picture.

## Work through it step by step

Take a fictional model reply: “all work on `t-02` was completed in zero minutes.” Check the SQL report: `t-02` has no study-session row, so `0` means no recorded sessions, not a measured zero-minute task. Write the precise statement “no sessions recorded.” For an outside fact, open the organisation’s official calendar or page and note the check date. If no source is accessible, say “unverified.”

## Predict and check

Write the expected output before running code. Execute it in `step-06` or `step-07`, compare every character and explain each line. Then change one safe fictional input and repeat “predict — run — change — explain — verify.” Do not turn an observation on a tiny dataset into a promise about real people.

```python
from pathlib import Path
from tempfile import TemporaryDirectory
from data_store import create_database, minutes_report
with TemporaryDirectory() as folder:
    db = Path(folder) / "example.db"
    create_database(db, "tasks.json", "study_sessions.csv")
    print(dict(minutes_report(db))["t-02"])
```

## Expected output

```text
0
```

## Catch the error

Two pages repeating one mistake are not independent confirmation. A URL alone does not prove that a document supports a specific statement. Do not disguise a guess with model “confidence.”

## Project change

Release 2.1 preserves the DB, report, reminders and backup. `labelled_tasks.csv` contains fictional examples only. `classifier.py` suggests `math`, `reading` or `review` but writes no decision to the database. Check each model output against held-out examples and let a person reconsider it.

## Task and evidence

Check three claims: one against CSV, one against an official timetable and one with no accessible source. Record source, check time and status for each.

## Transfer to a new setting

An AI says a course is free “forever.” How can you check the terms without turning today’s information into a permanent guarantee?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name one boundary of the idea. After 7 days explain it to a friend without reading the lesson and answer 7 of 10 block questions; analyse mistakes with concrete inputs. After 30 days open a fresh project copy, reproduce the check and transfer the idea to a new case. Keep a record of input, expected output and actual output.

## Primary reference to check

[Official reference](https://www.nist.gov/itl/ai-risk-management-framework)

## Next lesson

[Adding AI while keeping the human in control](/read/informatics-70-human-controlled-ai?lang=en)
