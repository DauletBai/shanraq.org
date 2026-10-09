# Mastery checkpoint: the assistant's algorithms on paper

_Lead (summary):_ **Assemble the assistant's paper version 0.3 and test it on ordinary, boundary, and surprising inputs.**

## Where we are on the map

This is the third block of one project. We use version 0.2 data but describe actions on paper first, where mistakes are easier to find before writing Python.

![Mastery checkpoint: the assistant's algorithms on paper](/static/course/informatics/map-26-algorithms-mastery-en.svg)

## Situation and question

A friend receives only `ALGORITHM-en.md` and the three tasks from version 0.2. They want to know when a reminder appears and how many tasks are done. If they must ask you what a missing date means, the algorithm is not ready to become Python.

## New words without gaps

A **specification** is a written promise about inputs, outputs, and failures. A **test case** is a concrete input with expected output. A **regression** is an old fault returning after a change. An **oracle** is an expected answer written before execution, not invented afterwards. **Acceptance** decides readiness using observable criteria. **Version 0.3** here means tested paper algorithms; executable code arrives in the next block. Version 0.2 `tasks.json` remains the source of three fictional records; dates are supplied on separate test cards, not silently added to the old format.

## The lesson's support signal

`contract → trace → boundary → invariant → independent check → version 0.3`

## Work through it step by step

Take `t-01` through the rules. Its test card says `done=false`, `days_to_due=2`; the function returns `REMIND`. Original `t-02` has `done=true`, so even a test card saying 0 days produces `DONE`. `t-03` has no due-date card: return `NO_DATE`. Counting the original list gives exactly 1 completed task. Then add a fourth fictional record with duplicate `id=t-02` only to the test list: the uniqueness check must reject it before reminding. Do not alter the original three records to force a convenient answer.

## Check it by hand

Run a short defence without the author. A partner opens `cases.json`, hides the `expected` column, and fills it from your rules. They compare it with the oracle and mark the first mismatch, not only the number correct. Ask for one new case outside the file, such as a completed task with an overdue date. If the contract cannot determine the answer, return to lesson 24 and clarify the checking order. Record the mismatch and repair in the project log.

## Predict and check

Before viewing answers, predict outputs for −1, 0, 1, 2, 3 days with `done=false`, then the same days with `done=true`. They should be `OVERDUE, REMIND, REMIND, REMIND, NOT_YET` and five `DONE` answers. A missing date gives `NO_DATE`; a done task with missing date gives `DONE` under our chosen contract.

## Catch the error

“The three original tasks opened, so 0.3 is ready.” Opening checks the 0.2 format but not branches, loop termination, duplicate identifiers, or another person's understanding. Show counterexamples and repairs, not one successful screen.

## Transfer to a new setting

A school club switches from days to hours and reminds at 48 hours inclusive. Which parts of your plan can be reused, and which requirement must change? Explain why two calendar days are not always exactly 48 hours if times of day differ.

## Project change

Submit `ALGORITHM-en.md`, trace tables, `cases.json` with fictional inputs, and a log of one repair. The version 0.3 README must say the paper algorithm does not yet run or send notifications. Another learner should obtain the same expected answers without your help.

## Task and evidence

Check ten points: (1) goal and input; (2) output and failures; (3) state before/after; (4) step order; (5) −1/0/1/2/3 cases; (6) missing date; (7) done task; (8) empty list and stopping; (9) duplicate `id` and invariant; (10) independent transfer to 48 hours. Passing requires 8/10 now, including 9 and 10, and 7/10 on a new case after seven days. Repair each failed item, then retry an equivalent problem.

## Return after 1, 7, and 30 days

Tomorrow recall the function's five answers. In seven days perform a new acceptance check without this page. In a month compare 0.3 promises with Python 1.0 and repair differences before release.

[Next lesson: First program: state becomes code](/read/informatics-27-python-first-state?lang=en)
