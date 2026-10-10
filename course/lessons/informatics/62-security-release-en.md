# Release 1.3: a protected and recoverable assistant

_Lead (summary):_ **You hand the assistant to a classmate.**

## Where we are on the map

Lesson 62 of 72. Digital security block (55–62). We continue a local assistant with three fictional tasks. Each control has a defined boundary: the classroom version does not become a public service.

![Release 1.3: a protected and recoverable assistant](/static/course/informatics/map-62-security-release-en.svg)

## Situation and question

You hand the assistant to a classmate. Without spoken hints, they should open the page, obtain the same report and restore a damaged teaching copy. That is the release 1.3 test.

## New words without gaps

A **checkpoint** is concrete code, data, commands and expected output. **Recovery** rebuilds a working copy from a verified one. A **threat model** explains what is protected and what lies outside the boundary. **Residual risk** is trouble still possible after controls are applied.

## The lesson's support signal

Read the three cells left to right, then cover the third. State the input and action in each cell; predict how the path ends and which observable fact checks it. Compare the drawing with a project file or worked example. If a number on the map cannot be derived from the data, correct the map rather than changing data for a pretty picture.

## Work through it step by step

In `step-06`, run `python3 data_assistant.py init`, then `python3 security_assistant.py verify assistant.db`: expect `tasks=3 sessions=4`. Make a backup, restore into a new file and compare the 45/0/25 report. Check the page on 127.0.0.1: it still shows three tasks and `REMIND`, `DONE`, `NO_DATE`. Stop the server. Never commit the database, backup or checksum.

## Predict and check

Write the expected output before running code. Execute it in `step-06` or `step-07`, compare every character and explain each line. Then change one safe fictional input and repeat “predict — run — change — explain — verify.” Do not turn an observation on a tiny dataset into a promise about real people.

```python
from pathlib import Path
from tempfile import TemporaryDirectory
from data_store import create_database, minutes_report
from security_assistant import backup, restore, verify
with TemporaryDirectory() as folder:
    a, b, c = (Path(folder) / name for name in ("a.db", "b.db", "c.db"))
    create_database(a, "tasks.json", "study_sessions.csv")
    backup(a, b); restore(b, c)
    print(verify(c), minutes_report(c))
```

## Expected output

```text
(3, 4) [('t-01', 45), ('t-02', 0), ('t-03', 25)]
```

## Catch the error

Release 1.3 is protected only within its classroom boundary: local viewing, no overwrite, integrity checks and recovery. It has no accounts, encrypted backup or public-server protection. Do not put it on the internet.

## Project change

Release 1.3 preserves the three tasks, four sessions and reminder rules. New operations live in `security_assistant.py`: validation, a separate backup and restore only into a new file. State which threat each change reduces and which it does not; all teaching data stay fictional.

## Task and evidence

Show five pieces of evidence: three tasks, four sessions, the 45/0/25 report, a verified restore and refusal of a tampered backup. Answer 8 of 10 review questions now and 7 without hints in a week.

## Ten questions for self-check

1. What does mode=ro protect?
2. What is the asset?
3. How do threat and vulnerability differ?
4. Why use distinct salts?
5. Why are two passwords not two factors?
6. What should you do with a phishing link?
7. How does a hash differ from encryption?
8. What does 3/4 mean?
9. Why restore to a new file?
10. What remains outside release 1.3?

## Transfer to a new setting

A friend suggests exposing the server port to the whole class. What new requirements must be met first?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name one boundary of the idea. After 7 days explain it to a friend without reading the lesson and answer 7 of 10 block questions; analyse mistakes with concrete inputs. After 30 days open a fresh project copy, reproduce the check and transfer the idea to a new case. Keep a record of input, expected output and actual output.

## Primary reference to check

[Official reference](https://www.sqlite.org/pragma.html#pragma_integrity_check)

## Next lesson

[Rules, algorithms, and trainable models](/read/informatics-63-rules-algorithms-models?lang=en)
