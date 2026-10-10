# Joining tables and making an honest report

_Lead (summary):_ **Join two tables and keep the task with no study sessions in the report.**

## Where we are on the map

Lesson 53 of 72, Data, tables and databases block (47–54). We continue assistant 1.1 with three fictional tasks. By the end, records live in local SQLite while tested reminder rules remain unchanged.

![Joining tables and making an honest report](/static/course/informatics/map-53-joins-reports-en.svg)

## Situation and question

An ordinary matching join shows only `t-01` and `t-03`. If `t-02` vanishes, a reader may think that task never existed. The kind of join changes the report's question.

## New words without gaps

A **JOIN** combines rows under a matching condition. **INNER JOIN** keeps matches in both tables. **LEFT JOIN** retains every left-hand row and supplies `NULL` where the right-hand side has no match. `COALESCE(SUM(...), 0)` displays zero instead of an absent total; the caption must say 'no recorded sessions'. **GROUP BY** aggregates by task ID. The report's **grain** is one row per task, not one per session.

## The lesson's support signal

Read the map left to right. Name each arrow's input, transformation and check. Cover the third step and predict it. Then locate the matching row in CSV, SQL or the report: an arrow must not replace real values or relationships.

## Work through it step by step

Put `tasks` on the left and `study_sessions` on the right. INNER JOIN yields four session rows and loses `t-02`. LEFT JOIN yields five: four sessions plus `t-02` with `NULL` session columns. `GROUP BY t.task_id` leaves three rows. Totals are `t-01` 45, `t-02` 0, `t-03` 25. Check 45+0+25=70, but do not claim zero means a measured zero-minute session.

## Predict and check

Write the exact expected output, including row order, before running anything. Execute the snippet in `step-05` and compare characters. Change one input, then predict, run and explain again. The example proves only the stated rule; fictional teaching data do not become measurements of real people.

```python
from pathlib import Path
from tempfile import TemporaryDirectory
from data_store import create_database, minutes_report
with TemporaryDirectory() as folder:
    db = Path(folder) / "test.db"
    create_database(db, "tasks.json", "study_sessions.csv")
    for task_id, minutes in minutes_report(db):
        print(task_id, minutes)
```

## Expected output

```text
t-01 45
t-02 0
t-03 25
```

## Catch the error

If you put `WHERE s.source = 'app'` after a LEFT JOIN, the task without sessions disappears: a condition on `NULL` is not true. For 'all tasks, app sessions only', put the right-table filter in `ON` and verify the result.

## Project change

`minutes_report` returns all three tasks. Record source-session count (4), report-row count (3) and check total (70) beside it. These reveal a category lost silently by a query.

## Task and evidence

Draw INNER and LEFT JOIN results before grouping. Why does one have four rows and the other five? Check the report for all three IDs and name the absence for `t-02`.

## Transfer to a new setting

You have a list of classes and a table of competition results. How would you show a class with no entrants without turning a missing result into a poor mark?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name a case where the simple analogy fails. After 7 days explain a fresh table without the text and answer at least 7 of the final lesson's 10 questions. After 30 days repeat import into a temporary folder and check three tasks, four sessions, total 70 and unchanged reminders. Record any unclear term and return to its first definition.

## Primary reference to check

[Official documentation](https://www.sqlite.org/lang_select.html)

## Next lesson

[Release 1.2: the assistant with SQLite](/read/informatics-54-data-release?lang=en)
