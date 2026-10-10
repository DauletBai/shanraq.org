# Related tables, keys, and no duplicates

_Lead (summary):_ **Link tasks and study sessions by keys so no record points to nowhere.**

## Where we are on the map

Lesson 51 of 72, Data, tables and databases block (47–54). We continue assistant 1.1 with three fictional tasks. By the end, records live in local SQLite while tested reminder rules remain unchanged.

![Related tables, keys, and no duplicates](/static/course/informatics/map-51-relational-keys-en.svg)

## Situation and question

If every session repeats the title 'Кітап оқу', one typo can make it look like a new task. Keeping the title once and referring to its ID is more reliable.

## New words without gaps

A **database** is organised storage with rules for reading and writing. A **relational table** describes one kind of entity: ours are `tasks` and `study_sessions`. A **primary key** is unique within its table: `task_id` or `session_id`. The **foreign key** `study_sessions.task_id` requires an existing task. **One-to-many** means one task may have many sessions, but each session belongs to one task. `NULL` for a due date means 'not assigned', not an empty string or a zero date.

## The lesson's support signal

Read the map left to right. Name each arrow's input, transformation and check. Cover the third step and predict it. Then locate the matching row in CSV, SQL or the report: an arrow must not replace real values or relationships.

## Work through it step by step

Draw three tasks on the left and four sessions on the right. Connect `s-01` and `s-02` to `t-01`; `s-03` and `s-04` to `t-03`. `t-02` has no arrow. Create two SQLite tables with primary keys and `REFERENCES tasks(task_id)`. Enable `PRAGMA foreign_keys = ON` for each Python connection; never assume this check is already enabled. Inserting `t-99` must fail.

## Predict and check

Write the exact expected output, including row order, before running anything. Execute the snippet in `step-05` and compare characters. Change one input, then predict, run and explain again. The example proves only the stated rule; fictional teaching data do not become measurements of real people.

```python
import sqlite3
db = sqlite3.connect(":memory:")
db.execute("PRAGMA foreign_keys = ON")
db.execute("CREATE TABLE tasks (task_id TEXT PRIMARY KEY)")
db.execute("CREATE TABLE sessions (task_id TEXT REFERENCES tasks(task_id))")
print(db.execute("PRAGMA foreign_keys").fetchone()[0])
db.close()
```

## Expected output

```text
1
```

## Catch the error

The same `task_id` in two sessions is an allowed relationship, not a duplicate. A duplicate would repeat a `session_id` or a task row's `task_id`. Foreign-key enforcement only works on connections where it is enabled.

## Project change

`step-05/data_store.py` contains the two tables, constraints and an index for sessions by task. Import builds a temporary database file; if any check fails, no completed database appears.

## Task and evidence

Label the keys on your drawing. Imagine adding `s-05,t-99`: which rule rejects it? Why can `t-02` have no sessions without violating the relationship?

## Transfer to a new setting

A library has `books` and `loans` tables. If one book may be loaned many times, which table should hold the foreign key, and what row must be rejected?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name a case where the simple analogy fails. After 7 days explain a fresh table without the text and answer at least 7 of the final lesson's 10 questions. After 30 days repeat import into a temporary folder and check three tasks, four sessions, total 70 and unchanged reminders. Record any unclear term and return to its first definition.

## Primary reference to check

[Official documentation](https://www.sqlite.org/foreignkeys.html)

## Next lesson

[SQL asks data precise questions](/read/informatics-52-sql-queries?lang=en)
