# Release 1.2: the assistant with SQLite

_Lead (summary):_ **Release assistant 1.2: SQLite data, unchanged reminder rules, and the same local browser page.**

## Where we are on the map

Lesson 54 of 72, Data, tables and databases block (47–54). We continue assistant 1.1 with three fictional tasks. By the end, records live in local SQLite while tested reminder rules remain unchanged.

![Release 1.2: the assistant with SQLite](/static/course/informatics/map-54-data-release-en.svg)

## Situation and question

A classmate asks, 'We moved tasks into a database; did the old reminders survive?' Check by comparing the three records, the five possible outcomes and the browser view, rather than promising they did.

## New words without gaps

A **migration** moves data from an older format to a new one under an explicit rule. A **transaction** groups changes: after an error, the set must not be treated as complete. A **checkpoint** contains code, source files, build command and a checkable result. **Read-only** means opening the page does not change the DB file. SQLite stores tasks and sessions in one file on your own computer; it is neither a cloud nor a shared database for site visitors.

## The lesson's support signal

Read the map left to right. Name each arrow's input, transformation and check. Cover the third step and predict it. Then locate the matching row in CSV, SQL or the report: an arrow must not replace real values or relationships.

## Work through it step by step

From `step-05`, run `python3 data_assistant.py init`, then `python3 data_assistant.py report`: expect `t-01 45`, `t-02 0`, `t-03 25`. Run `python3 web_assistant.py --port 8765` and open the date 2026-10-09. Statuses remain `REMIND`, `DONE`, `NO_DATE`; `ru`, `kz`, `en` change interface labels only. Compare the three `tasks.json` rows with those read from SQLite. Run the tests. Running `init` again must stop without replacing the existing database.

## Predict and check

Write the exact expected output, including row order, before running anything. Execute the snippet in `step-05` and compare characters. Change one input, then predict, run and explain again. The example proves only the stated rule; fictional teaching data do not become measurements of real people.

```python
from pathlib import Path
from tempfile import TemporaryDirectory
from datetime import date
from data_store import create_database, load_document_from_db
from assistant_core import reminder_status
with TemporaryDirectory() as folder:
    db = Path(folder) / "test.db"
    create_database(db, "tasks.json", "study_sessions.csv")
    for task in load_document_from_db(db)["tasks"]:
        print(task["id"], reminder_status(task, date(2026, 10, 9)))
```

## Expected output

```text
t-01 REMIND
t-02 DONE
t-03 NO_DATE
```

## Catch the error

Do not put real learner data in this classroom file or expose `http.server` to the public internet. SQLite does not automatically give you accounts, backups, permissions or protection from unwanted requests. An unchanged DB after GET is not proof that a future public system is secure.

## Project change

Keep `step-05`: the unchanged 1.0 core, fictional JSON and CSV, SQLite schema, read-only web view and tests. Each learner creates `assistant.db` locally; it is not committed to Git. The next block studies threats and recovery separately.

## Task and evidence

Show the three statuses and the 45/0/25 report. Explain the zero for `t-02`, the invalid-date refusal and the refusal to import twice. For mastery answer 8 of 10 questions and let another learner open the page from the instructions without your help.

## Ten questions for self-check

1. How does an observation differ from a conclusion recalled later?
2. How do you distinguish unknown duration from zero minutes?
3. Why do the four records total 70 rather than 60?
4. Why keep the raw CSV before cleaning?
5. Which key prevents duplicate sessions, and which prevents links to missing tasks?
6. Why enable foreign-key checks for every SQLite connection?
7. How do COUNT(*) and SUM(minutes) differ on an empty selection?
8. Why does LEFT JOIN keep t-02 while INNER JOIN loses it?
9. What three checks make a chart honest?
10. What survived migration to 1.2, and why is the local server not public?

## Transfer to a new setting

A school wants families in different homes to use this assistant. What must change before moving the local SQLite file and teaching HTTP server to a shared host? Name at least access control, data protection and backup recovery.

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name a case where the simple analogy fails. After 7 days explain a fresh table without the text and answer at least 7 of the final lesson's 10 questions. After 30 days repeat import into a temporary folder and check three tasks, four sessions, total 70 and unchanged reminders. Record any unclear term and return to its first definition.

## Primary reference to check

[Official documentation](https://www.sqlite.org/atomiccommit.html)

## All course lessons

[All lessons](/course/informatics?lang=en)
