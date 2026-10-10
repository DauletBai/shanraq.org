# Backups, updates, and event logs

_Lead (summary):_ **Imagine a spare house key: unless it has been tried, you do not know whether it opens the lock.**

## Where we are on the map

Lesson 60 of 72. Digital security block (55–62). We continue a local assistant with three fictional tasks. Each control has a defined boundary: the classroom version does not become a public service.

![Backups, updates, and event logs](/static/course/informatics/map-60-backup-updates-logs-en.svg)

## Situation and question

Imagine a spare house key: unless it has been tried, you do not know whether it opens the lock. A backup is useful only after a real restore test.

## New words without gaps

A **backup** is a separate copy made for recovery. A **SHA-256 checksum** detects a changed file if the expected digest was kept trustworthy. An **event log** describes actions, but should not contain passwords or private text. An **update** fixes known faults; get it from a trusted source and retain a rollback plan.

## The lesson's support signal

Read the three cells left to right, then cover the third. State the input and action in each cell; predict how the path ends and which observable fact checks it. Compare the drawing with a project file or worked example. If a number on the map cannot be derived from the data, correct the map rather than changing data for a pretty picture.

## Work through it step by step

In `step-06`, make a database with `data_assistant.py init`. Run `python3 security_assistant.py verify assistant.db`; expect `tasks=3 sessions=4`. Then make `backups/first.db` with `backup` and restore into a new `restored.db`; compare the 45/0/25 report. The command refuses to overwrite an existing destination. Add one byte to a teaching copy and check that SHA-256 verification rejects restore. Never damage your only working database.

## Predict and check

Write the expected output before running code. Execute it in `step-06` or `step-07`, compare every character and explain each line. Then change one safe fictional input and repeat “predict — run — change — explain — verify.” Do not turn an observation on a tiny dataset into a promise about real people.

```python
from pathlib import Path
from tempfile import TemporaryDirectory
from data_store import create_database
from security_assistant import backup, restore, verify
with TemporaryDirectory() as folder:
    db, saved, restored = (Path(folder) / name for name in ("live.db", "backup.db", "restored.db"))
    create_database(db, "tasks.json", "study_sessions.csv")
    backup(db, saved)
    print(verify(db), restore(saved, restored))
```

## Expected output

```text
(3, 4) (3, 4)
```

## Catch the error

A matching SHA-256 stored beside the backup does not prove authorship: an attacker may replace both. Creating a copy does not prove it can become a working database. Keep secrets out of logs.

## Project change

Release 1.3 preserves the three tasks, four sessions and reminder rules. New operations live in `security_assistant.py`: validation, a separate backup and restore only into a new file. State which threat each change reduces and which it does not; all teaching data stay fictional.

## Task and evidence

Show commands, expected counts 3/4 and 45/0/25, refusal after tampering and refusal to overwrite. Explain where to keep a separate copy and how you would test an update first.

## Transfer to a new setting

A family keeps photographs on one disk and a “backup” in another folder on that same disk. What happens when the disk breaks?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name one boundary of the idea. After 7 days explain it to a friend without reading the lesson and answer 7 of 10 block questions; analyse mistakes with concrete inputs. After 30 days open a fresh project copy, reproduce the check and transfer the idea to a new case. Keep a record of input, expected output and actual output.

## Primary reference to check

[Official reference](https://docs.python.org/3/library/sqlite3.html#sqlite3.Connection.backup)

## Next lesson

[Privacy, authorship, and digital wellbeing](/read/informatics-61-privacy-rights-wellbeing?lang=en)
