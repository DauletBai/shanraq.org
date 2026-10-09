# Text, dates, and Kazakh Unicode

_Lead (summary):_ **Preserve Kazakh text and turn a date string into a calendar-day distance from a deadline.**

## Where we are on the map

This is lesson 33 of 72 in the Python block (27–38). The support map shows verifiable transitions. We continue the 0.3 paper contract: five reminder outcomes and the original fictional tasks.

![Text, dates, and Kazakh Unicode](/static/course/informatics/map-33-python-text-dates-en.svg)

## Situation and question

A task says “Әліппе оқу.” After saving, Ә becomes unreadable, while someone calls `2026-10-11` “two days away” without naming today's date. These are separate faults: lost characters and an unstated reference date. Investigate both with reproducible data.

## New words without gaps

**Unicode** assigns code points to characters, letting Ә and A remain different letters. **UTF-8** encodes those characters as bytes; reading and writing must agree on encoding. The **string** `"2026-10-11"` is still just text. A `date` is a calendar day that Python can subtract. The **ISO format** `YYYY-MM-DD` fixes year–month–day order. A **time zone** determines what calendar day it is where the user lives. This experiment explicitly sets `today` to 9 October 2026; it does not pretend that the old 0.2 file stores a due date.

## The lesson's support signal

Input → check the rule → change state → observable result. Point to where the program reads data, compares it with the contract, and merely reports the result. If a step is absent from the map, find it in the code and add it to your own trace.

## Work through it step by step

Write `2026-10-09` and `2026-10-11` on two cards and mark two calendar transitions. In code `date.fromisoformat` turns strings into dates, subtraction yields a `timedelta`, and `.days` gives integer 2. Ә prints correctly when the source file and terminal use UTF-8. Change the deadline to `2026-02-29`: Python rejects a nonexistent day; this is useful evidence, not a reason to round the date.

## Predict and check

Cover the output block. Trace names line by line and write the exact predicted text, including case and line order. Only then run the code with `python3` from `step-03` and compare character by character. Change one input and predict again before running. If results differ, find the first divergent line instead of adjusting the prediction afterwards.

```python
from datetime import date
today = date.fromisoformat("2026-10-09")
due = date.fromisoformat("2026-10-11")
print("Әліппе оқу", (due - today).days)
```

## Expected output

```text
Әліппе оқу 2
```

## Catch the error

Two calendar days are not always exactly 48 hours: a deadline may be in the morning while it is now evening. Version 0.3 specifies calendar dates. An alarm or stopwatch needs time of day and a time zone; do not reuse this calculation blindly.

## Project change

Version 1.0 will store a due date in a separate `due_date` field as an ISO date string or `null`. Do not confuse it with the paper cards' `days_to_due`: that number is computed from the date and chosen `today`.

## Task and evidence

Use three fictional titles, one containing Ә, one Я, and one Latin A. Save and read them in UTF-8 without loss. For fixed `today=2026-10-09`, compute deadlines on 8, 9, 11, and 12 October: expect −1, 0, 2, 3. Show that 29 February 2026 is rejected before any reminder is made.

## Transfer to a new setting

An international class is scheduled for “10 October, 09:00” without a time zone. Can a learner in Qostanay calculate exact hours remaining? Name the missing input that must be requested.

## Return after 1, 7, and 30 days

After 1 day, recall the rule and one boundary case without this page. After 7 days, explain a new error example to a classmate. After 30 days, rerun the project test, check earlier records and output still match, and transfer the rule to a different task again.

## Next lesson

[Continue: Python functions with a clear contract](/read/informatics-34-python-functions?lang=en)
