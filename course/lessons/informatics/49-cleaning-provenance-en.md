# Cleaning data without losing its provenance

_Lead (summary):_ **Correct a bad row while preserving the original and the reason for every decision.**

## Where we are on the map

Lesson 49 of 72, Data, tables and databases block (47–54). We continue assistant 1.1 with three fictional tasks. By the end, records live in local SQLite while tested reminder rules remain unchanged.

![Cleaning data without losing its provenance](/static/course/informatics/map-49-cleaning-provenance-en.svg)

## Situation and question

Another copy of the journal contains `s-04` twice and a date `2026-02-29`. Can we just delete the ugly rows? That would hide evidence that an export was duplicated or damaged.

## New words without gaps

**Cleaning** checks and repairs data under stated rules. **Provenance** is the path from source to result. A **raw file** preserves what arrived; a **cleaned set** contains accepted rows. A **quarantine** lists rejected rows with reasons. Here a **duplicate** means the same `session_id`, even if other fields differ. **Validation** checks real calendar dates, integer minutes from 0 to 180, known `task_id`, and `paper`/`app` sources.

## The lesson's support signal

Read the map left to right. Name each arrow's input, transformation and check. Cover the third step and predict it. Then locate the matching row in CSV, SQL or the report: an arrow must not replace real values or relationships.

## Work through it step by step

Make a copy of the CSV and leave the original untouched. Inspect rows in order: accept `s-01`, then quarantine the second `s-01` as a repeated ID. Reject 29 February 2026: that year is not a leap year. Reject `t-99` as a link to a task that does not exist. In code, `read_sessions` validates all rows before returning them to the database builder. One bad observation stops import; at this stage we require a deliberate source correction.

## Predict and check

Write the exact expected output, including row order, before running anything. Execute the snippet in `step-05` and compare characters. Change one input, then predict, run and explain again. The example proves only the stated rule; fictional teaching data do not become measurements of real people.

```python
from data_store import read_sessions
rows = read_sessions("study_sessions.csv", {"t-01", "t-02", "t-03"})
print(len(rows), sum(row[3] for row in rows))
```

## Expected output

```text
4 70
```

## Catch the error

Deleting a row without a decision log makes the result impossible to reproduce. Automatically replacing an unknown date with today invents a fact. Turning a blank duration into zero silently claims the session took zero minutes.

## Project change

Keep the four original fictional rows. Make a paper quarantine for three damaged examples and repair only a working copy. The 1.2 importer refuses to create a database until the whole CSV passes validation.

## Task and evidence

Write a separate rejection reason for a repeated ID, an impossible date, and an unknown `task_id`. What evidence would justify repairing each row, and which file stays unchanged?

## Transfer to a new setting

Two sites publish different figures for one event. Why should you keep the URLs, retrieval dates and selection rules before drawing one combined chart?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name a case where the simple analogy fails. After 7 days explain a fresh table without the text and answer at least 7 of the final lesson's 10 questions. After 30 days repeat import into a temporary folder and check three tasks, four sessions, total 70 and unchanged reminders. Record any unclear term and return to its first definition.

## Primary reference to check

[Official documentation](https://docs.python.org/3/library/csv.html)

## Next lesson

[How a chart can lie with true numbers](/read/informatics-50-charts-honesty?lang=en)
