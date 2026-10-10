# SQL asks data precise questions

_Lead (summary):_ **Ask the database a precise question: which sessions belong to a task, and how long did they take?**

## Where we are on the map

Lesson 52 of 72, Data, tables and databases block (47–54). We continue assistant 1.1 with three fictional tasks. By the end, records live in local SQLite while tested reminder rules remain unchanged.

![SQL asks data precise questions](/static/course/informatics/map-52-sql-queries-en.svg)

## Situation and question

A friend says, 'Find the reading sessions.' The database does not understand a hint; it needs a table, condition and result order. One wrong character in an ID can return no rows, which differs from zero measured minutes.

## New words without gaps

**SQL** is a language for querying relational data. `SELECT` chooses columns, `FROM` names a table, `WHERE` filters rows, and `ORDER BY` sets their order. `SUM` adds; `COUNT` counts records. A query **parameter** `?` receives a value separately from SQL text, keeping user data separate from commands. **Aggregation** turns several rows into a summary. `task_id = 't-01'` returns two records totalling 45 minutes.

## The lesson's support signal

Read the map left to right. Name each arrow's input, transformation and check. Cover the third step and predict it. Then locate the matching row in CSV, SQL or the report: an arrow must not replace real values or relationships.

## Work through it step by step

Read this query in order: `SELECT minutes FROM study_sessions WHERE task_id = ? ORDER BY observed_on`, parameter `t-01`. The result is 20, then 25. Change the parameter to `t-03`: 15, then 10. In Python pass `(task_id,)`, not a SQL string assembled from user input. `SELECT COUNT(*), SUM(minutes) ...` returns count 2 and total 45 for `t-01`. For `t-02`, `COUNT(*)` is 0 but `SUM` over no rows is `NULL`; they do not mean the same thing.

## Predict and check

Write the exact expected output, including row order, before running anything. Execute the snippet in `step-05` and compare characters. Change one input, then predict, run and explain again. The example proves only the stated rule; fictional teaching data do not become measurements of real people.

```python
from pathlib import Path
from tempfile import TemporaryDirectory
from data_store import create_database, connect_readonly
with TemporaryDirectory() as folder:
    db = Path(folder) / "test.db"
    create_database(db, "tasks.json", "study_sessions.csv")
    with connect_readonly(db) as connection:
        rows = connection.execute("SELECT minutes FROM study_sessions WHERE task_id = ? ORDER BY observed_on", ("t-01",)).fetchall()
    print(*(row[0] for row in rows))
```

## Expected output

```text
20 25
```

## Catch the error

`ORDER BY` changes row order, not durations. Without it, do not rely on a particular order. Inserting user text directly into SQL may let special characters alter a command; parameters separate data from command text.

## Project change

Record three checkable queries in the project journal: list for `t-01`, total for `t-03`, count for `t-02`. Compare every number with the source CSV.

## Task and evidence

Predict the rows for `t-03`. Explain why `(task_id,)` needs the comma to be a one-value Python tuple. Compare `COUNT(*)` and `SUM(minutes)` for a task without sessions.

## Transfer to a new setting

A library query by book ID returns no loans. Does that prove the book exists but was never borrowed? How would you check the book itself?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name a case where the simple analogy fails. After 7 days explain a fresh table without the text and answer at least 7 of the final lesson's 10 questions. After 30 days repeat import into a temporary folder and check three tasks, four sessions, total 70 and unchanged reminders. Record any unclear term and return to its first definition.

## Primary reference to check

[Official documentation](https://docs.python.org/3/library/sqlite3.html)

## Next lesson

[Joining tables and making an honest report](/read/informatics-53-joins-reports?lang=en)
