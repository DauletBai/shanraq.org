# From an observation to a data row

_Lead (summary):_ **Turn one study observation into a readable row: what was measured, when, and for which task.**

## Where we are on the map

Lesson 47 of 72, Data, tables and databases block (47–54). We continue assistant 1.1 with three fictional tasks. By the end, records live in local SQLite while tested reminder rules remain unchanged.

![From an observation to a data row](/static/course/informatics/map-47-observations-data-schema-en.svg)

## Situation and question

You say, 'I studied maths yesterday.' A friend writes 20 minutes, but you remember 25. Was the duration timed, estimated, or recalled later? Without that distinction, a neat table can look more precise than the evidence.

## Where to get the project files

Open the [1.2 project folder](https://github.com/DauletBai/shanraq.org/tree/main/course/informatics-assistant/step-05), download the repository with Code → Download ZIP and find `course/informatics-assistant/step-05`. Open a terminal there; `python3 --version` should show an installed Python. Begin with CSV and JSON; do not run the database file as a program.

## New words without gaps

An **observation** is a fact recorded under a chosen rule. A **row** holds one observation; a **column** holds one property across rows. A **schema** defines column names and types. The **unit** here is minutes. **Source** records whether a value came from a paper journal or the teaching app. `session_id` distinguishes records; `task_id` links a session to a task. Neither ID identifies a person.

## The lesson's support signal

Read the map left to right. Name each arrow's input, transformation and check. Cover the third step and predict it. Then locate the matching row in CSV, SQL or the report: an arrow must not replace real values or relationships.

## Work through it step by step

Open the fictional `study_sessions.csv`. The row `s-01,t-01,2026-10-07,20,paper` describes one study session: task `t-01`, 7 October, 20 minutes, paper source. Commas separate fields; the first line names them. Compare the next row: it is the same task but a different date, duration and source. Four rows mean four observations, not four learners. Check `tasks.json` really contains `t-01` and `t-03`; otherwise the link would be imaginary.

## Predict and check

Write the exact expected output, including row order, before running anything. Execute the snippet in `step-05` and compare characters. Change one input, then predict, run and explain again. The example proves only the stated rule; fictional teaching data do not become measurements of real people.

```python
import csv
with open("study_sessions.csv", encoding="utf-8", newline="") as stream:
    first = next(csv.DictReader(stream))
print(first["session_id"], first["task_id"], first["minutes"])
```

## Expected output

```text
s-01 t-01 20
```

## Catch the error

A duration recalled from memory is not a precisely timed measurement. Do not mix hours and minutes in one column. Two rows for the same task are not automatically duplicates; `session_id` distinguishes them.

## Project change

Add a separate `study_sessions.csv` with four fictional observations to version 1.1. Keep the tasks and reminder rules unchanged. Write a data passport: who recorded each row, what a minute means, accepted values, and what we deliberately do not collect (names, marks, exact location).

## Task and evidence

Draw the five columns. For `s-03`, locate its date, duration and `task_id`. Invent a new fictional row with a new `session_id`, and explain how you would validate its date and unit.

## Transfer to a new setting

A school club counts attendance. Why is it wrong to enter zero minutes for an attendee whose duration is unknown? How can the table show 'unknown' without pretending it means zero?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name a case where the simple analogy fails. After 7 days explain a fresh table without the text and answer at least 7 of the final lesson's 10 questions. After 30 days repeat import into a temporary folder and check three tasks, four sessions, total 70 and unchanged reminders. Record any unclear term and return to its first definition.

## Primary reference to check

[Official documentation](https://docs.python.org/3/library/csv.html)

## Next lesson

[A table, a formula, and a verifiable calculation](/read/informatics-48-spreadsheets-formulas?lang=en)
