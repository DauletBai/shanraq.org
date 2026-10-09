# Python conditions and boundary checks

_Lead (summary):_ **Translate reminder boundaries into `if`, `elif`, and `else`, ensuring a completed task never takes the wrong branch.**

## Where we are on the map

This is lesson 30 of 72 in the Python block (27–38). The support map shows verifiable transitions. We continue the 0.3 paper contract: five reminder outcomes and the original fictional tasks.

![Python conditions and boundary checks](/static/course/informatics/map-30-python-conditions-en.svg)

## Situation and question

The previous expression answered only yes or no. The assistant needs five outcomes: done, no date, overdue, remind, or not yet. How do we choose exactly one without losing an exceptional case? The paper contract made checking order important; in Python it becomes line order.

## New words without gaps

A **condition** is an expression yielding `True` or `False`. **Branching** chooses a path from that answer. `if` checks first, `elif` checks another condition only if earlier ones failed, and `else` takes what remains. Lines at the same indentation form a **block**; in Python indentation is part of syntax, the language's writing rules. A **boundary case** sits exactly on a threshold, such as two days. **Mutually exclusive outcomes** mean one task receives exactly one status.

## The lesson's support signal

Input → check the rule → change state → observable result. Point to where the program reads data, compares it with the contract, and merely reports the result. If a step is absent from the map, find it in the code and add it to your own trace.

## Work through it step by step

Lay out cards −1, 0, 2, 3. Walk each through the sample branches: −1 stops at `days < 0`; 0 and 2 reach `days <= 2`; 3 falls to `else`. Now imagine `done=True`: the paper contract demands `DONE` before any date comparison. The short sample below omits `done` to isolate boundaries; put it first in the project function.

## Predict and check

Cover the output block. Trace names line by line and write the exact predicted text, including case and line order. Only then run the code with `python3` from `step-03` and compare character by character. Change one input and predict again before running. If results differ, find the first divergent line instead of adjusting the prediction afterwards.

```python
for days in (-1, 0, 2, 3):
    if days < 0:
        status = "OVERDUE"
    elif days <= 2:
        status = "REMIND"
    else:
        status = "NOT_YET"
    print(days, status)
```

## Expected output

```text
-1 OVERDUE
0 REMIND
2 REMIND
3 NOT_YET
```

## Catch the error

If `days <= 2` comes first, overdue −1 becomes `REMIND`. Four independent `if` statements may print several answers. Repair the order with an `if`/`elif`/`else` chain, and handle a missing date before comparing it with a number.

## Project change

Turn the version 0.3 table into a function using `return`; do not print inside each branch. Leave display to the caller so the same function can later serve the CLI and tests.

## Task and evidence

Write a complete function that takes `done` and `days` and returns one of the paper contract's five outcomes. Before running it, draw a table for −1, 0, 1, 2, 3, `None`, and `done=True`. Give a partner only the code and test cards; they should obtain your answers without spoken help.

## Transfer to a new setting

A ticket machine might return “paid,” “price missing,” “not enough money,” or “ready.” Which condition must come first if the purchase is already paid? Explain how that resembles `done=True`.

## Return after 1, 7, and 30 days

After 1 day, recall the rule and one boundary case without this page. After 7 days, explain a new error example to a classmate. After 30 days, rerun the project test, check earlier records and output still match, and transfer the rule to a different task again.

## Next lesson

[Continue: Python loops that know when to stop](/read/informatics-31-python-loops?lang=en)
