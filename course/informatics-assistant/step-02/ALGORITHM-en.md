# My Digital Assistant algorithms, version 0.3

This is a paper contract for the future program. The three original fictional
tasks in `../step-01/data/tasks.json` remain unchanged. Version 0.2 has no due
date field. Day counts arrive on separate `cases.json` test cards; writing a
date into the file requires an explicit format change in the next block.

## 1. Validate input

Open the file as UTF-8 and parse JSON. Check `version="0.2"`, the `tasks`
array, nonempty unique `id` and `title`, and Boolean `done`. On failure,
return a clear message and do not show partial data as trustworthy. To find
duplicate IDs, scan left to right with an empty set `seen`: if the next ID
is already present, stop with an error; otherwise add it. Invariant: `seen`
contains exactly the IDs already visited.

## 2. Count completed tasks

Start with `position=0`, `completed=0`. While `position < len(tasks)`, add 1 to
`completed` if `tasks[position].done=true`, then advance `position` by 1.
Return `completed`. An empty list returns 0. Invariant: the counter equals
the number of done tasks in the visited prefix. Position advances on every
pass, so the algorithm terminates.

## 3. Decide whether to remind

`should_remind(done, days_to_due)` accepts a Boolean and an integer count of
calendar days or a missing date. It **does not change the file**.

| Check in order | Return value |
|---|---|
| `done=true` | `DONE` |
| date missing | `NO_DATE` |
| `days_to_due < 0` | `OVERDUE` |
| `0 ≤ days_to_due ≤ 2` | `REMIND` |
| remaining dates | `NOT_YET` |

Use fixed date `2026-10-09` and time zone `Asia/Qostanay` when reproducing the
`cases.json` cards. Do not confuse calendar days with exact 24-hour periods.
A missing date is not zero.

## 4. Trace and acceptance

For `false, true, false`, the count follows `0 → 0 → 1 → 1` and returns 1.
For `done=false` and days `−1, 0, 1, 2, 3`, expect `OVERDUE, REMIND, REMIND,
REMIND, NOT_YET`. For `done=true`, always expect `DONE`. Reject a duplicate
`id` before counting or reminding.

Another learner must reproduce these results using only this document and
`cases.json`. Score lesson 26's ten items: 8/10 now, including duplicate-ID
checking and transfer to 48 hours; 7/10 on a new case after seven days.
Repair any mismatch and try an equivalent test. This version sends no
notifications.
