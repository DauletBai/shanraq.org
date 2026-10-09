# Conditions and boundary cases

_Lead (summary):_ **Test a condition at its boundary so the assistant does not miss a task exactly two days before its due date.**

## Where we are on the map

This is the third block of one project. We use version 0.2 data but describe actions on paper first, where mistakes are easier to find before writing Python.

![Conditions and boundary cases](/static/course/informatics/map-22-conditions-boundaries-en.svg)

## Situation and question

You agreed to remind about unfinished tasks when no more than two days remain. The assistant stays silent at exactly two days, then wakes someone after the deadline. One comparison sign is wrong. Find it with a table rather than a guess.

## New words without gaps

A **condition** is a yes/no question. A **branch** chooses actions from the answer. A **boundary** is a value where the answer changes: here 2 days. A **boundary case** tests 2 itself and its neighbours 1 and 3. A **Boolean** is `true` or `false`; **AND** requires both parts and **OR** at least one. **Overdue** means the date has passed; it is a separate state, not “another zero days.” Compare dates under one calendar rule and one time zone, or the same instant can give different results.

## The lesson's support signal

`unfinished AND 0 ≤ days_to_due ≤ 2 → remind; past due → overdue`

## Work through it step by step

Let `today = 2026-10-09` and count whole calendar days in one time zone. Due 12 October means 3 days: no reminder; 11 October means 2: remind; 10 October means 1: remind; 9 October means 0: remind today; 8 October means −1: show “overdue” under a separate rule. If `done=true`, none of those dates triggers a reminder. Make a five-row table and add a `done` column: the same date with another state gives another answer.

## Check it by hand

Mark −1, 0, 1, 2, and 3 on a number line. Shade only the interval from 0 through 2, including both ends. Translate the picture into `0 ≤ days_to_due ≤ 2`: each part protects one boundary. Remove the left part and test −1; remove the right part and test 3. Restore both. Then place a `done=true` card over point 1: there is still no reminder. These are two independent reasons for a different answer.

## Predict and check

What does `days < 2` return at exactly 2 days? `false`: it wrongly excludes the boundary. What does `days <= 2` return at −1? `true`: it mixes an overdue task with future ones. Add a lower bound of 0 and a separate past-due branch. If the due date is missing, do not compare at all; report “no date set.”

## Catch the error

“Changing `≤` to `<` affects only one day.” That exact day is included by “no more than two.” Check a fix with three neighbouring values, not one convenient example.

## Transfer to a new setting

A school club accepts applications through 18:00 inclusive. Test 17:59, 18:00, and 18:01. Clarify the time zone and whether 18:00:30 is included: otherwise the requirement is ambiguous.

## Project change

In `ALGORITHM-en.md`, write checks in order: record found, date present, task unfinished, not overdue, 0–2 days remain. Add tests for −1, 0, 1, 2, 3, and `done=true`. State the calendar date and time zone so the experiment can be repeated.

## Task and evidence

Repair `days < 2` and prepare at least seven tests including missing date and done task. Evidence: inputs, expected messages, and reasons for each boundary.

## Return after 1, 7, and 30 days

Tomorrow explain why 2 is included but −1 is not. In seven days test another rule's boundary. In a month compare the paper cases with real date-handling code.

[Next lesson: Repetition without copying and a stopping condition](/read/informatics-23-loops-invariants?lang=en)
