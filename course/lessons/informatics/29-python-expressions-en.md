# Expressions without hidden conversions

_Lead (summary):_ **Write the reminder condition as a readable expression and check operation order without hidden text-to-number conversion.**

## Where we are on the map

This is lesson 29 of 72 in the Python block (27–38). The support map shows verifiable transitions. We continue the 0.3 paper contract: five reminder outcomes and the original fictional tasks.

![Expressions without hidden conversions](/static/course/informatics/map-29-python-expressions-en.svg)

## Situation and question

The assistant should remind us when a task is unfinished and zero to two days remain. A learner writes one long expression and gets an unexpected `True`. Which comparison runs first? Instead of guessing, split the rule into small questions, as with the paper cards.

## New words without gaps

An **expression** computes a value: `days + 1` yields a number; `days <= 2` yields a Boolean. An **operator** says what to do: `+`, `<=`, `and`, `not`. **Operands** are the values it acts on. **Precedence** sets the order: multiplication precedes addition; parentheses make intent explicit. `and` needs both parts true; `not done` reverses a Boolean. The chained comparison `0 <= days <= 2` includes both ends. Python will not add string `"2"` to number 1 as arithmetic; convert its type first.

## The lesson's support signal

Input → check the rule → change state → observable result. Point to where the program reads data, compares it with the contract, and merely reports the result. If a step is absent from the map, find it in the code and add it to your own trace.

## Work through it step by step

Try `done=False, days=2`: `not done` is `True` and `0 <= 2 <= 2` is also `True`, so the whole condition is `True`. Change only `days` to −1: the lower boundary becomes false. Then set `done=True`: even with two days left, no reminder is needed. In the example's second line `3 * 4` happens before adding 2, giving 14 rather than 20. Draw two small tables before running it.

## Predict and check

Cover the output block. Trace names line by line and write the exact predicted text, including case and line order. Only then run the code with `python3` from `step-03` and compare character by character. Change one input and predict again before running. If results differ, find the first divergent line instead of adjusting the prediction afterwards.

```python
days = 2
limit = 2
print(days + 1)
print(days <= limit)
```

## Expected output

```text
3
True
```

## Catch the error

The expression `days <= 2` accidentally includes −1: that task is overdue, not upcoming. Repair it by adding the lower bound `0 <= days`. This expression does not yet handle a missing date (`None`); do not compare `None` with a number. Check whether a number exists first.

## Project change

Record `not done and 0 <= days <= 2` as a temporary check of numbers on test cards. Until the file gains a date, do not present that number as a stored deadline.

## Task and evidence

Make a table for `days = −1, 0, 1, 2, 3` with each value of `done`. Fill ten answers in pencil before running code in a loop or by hand. For each mistaken prediction, identify the exact part of the expression that changed the result. “Python decided so” is not an explanation.

## Transfer to a new setting

A school club accepts applications through 18:00 inclusive. Which comparisons distinguish 17:59, 18:00, and 18:01? Minutes and time zone matter here; our calendar-day expression transfers as a boundary idea, not ready-made code.

## Return after 1, 7, and 30 days

After 1 day, recall the rule and one boundary case without this page. After 7 days, explain a new error example to a classmate. After 30 days, rerun the project test, check earlier records and output still match, and transfer the rule to a different task again.

## Next lesson

[Continue: Python conditions and boundary checks](/read/informatics-30-python-conditions?lang=en)
