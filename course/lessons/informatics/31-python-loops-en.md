# Python loops that know when to stop

_Lead (summary):_ **Replace three copied commands with a loop, trace the count, and prove stopping on empty and nonempty lists.**

## Where we are on the map

This is lesson 31 of 72 in the Python block (27–38). The support map shows verifiable transitions. We continue the 0.3 paper contract: five reminder outcomes and the original fictional tasks.

![Python loops that know when to stop](/static/course/informatics/map-31-python-loops-en.svg)

## Situation and question

Three tasks can be checked with three separate lines. What about 300? A loop repeats one check for each element. Yet a learner who cannot say when it stops may write a program that hangs or misses the last task. First trace three cards and an empty list.

## New words without gaps

A **list** stores values in a definite order. `for done in values` takes the next value; one run of the loop body is an **iteration**. The **counter** `total` holds how many `True` values were found. An **invariant** is a promise true before and after every iteration: `total` equals the completed cards already visited. **Termination** follows from finite list length: nothing remains after the last item. For an empty list the body never runs and the answer is zero.

## The lesson's support signal

Input → check the rule → change state → observable result. Point to where the program reads data, compares it with the contract, and merely reports the result. If a step is absent from the map, find it in the code and add it to your own trace.

## Work through it step by step

Lay out `False, True, False` from left to right. After the first card `total=0`; after the second it is 1; after the third it remains 1. The code prints that trace. Replace `values` with `[]`: there are no intermediate lines, only `Completed: 0`. Compare the paper table with the screen. In a `for` loop Python moves to the next item itself; you do not increment a separate position.

## Predict and check

Cover the output block. Trace names line by line and write the exact predicted text, including case and line order. Only then run the code with `python3` from `step-03` and compare character by character. Change one input and predict again before running. If results differ, find the first divergent line instead of adjusting the prediction afterwards.

```python
tasks = [{"done": False}, {"done": True}, {"done": True}]
count = 0
for task in tasks:
    if task["done"]:
        count += 1
print(count)
```

## Expected output

```text
2
```

## Catch the error

Putting `total = 0` inside the loop restarts the count for every card. Incrementing it for every item without checking `done` counts list length, not completed tasks. Use `False, False, True` as a counterexample to both mistakes.

## Project change

Replace the 0.3 paper counter with `count_done`. Compare it with four acceptance cards in `cases.json`; the original three tasks should yield 1.

## Task and evidence

Write `count_done` for the project's list of task dictionaries. Check four cases: empty, all `False`, a three-item mixture, and all `True`. Show every intermediate count for the mixture. Then remove printing from the function and return the number so a test can compare it with the expected result.

## Transfer to a new setting

Count books marked `borrowed=True` instead of completed tasks. What changes in the condition, and what stays the same in the invariant? Explain why an empty shelf still gives zero.

## Return after 1, 7, and 30 days

After 1 day, recall the rule and one boundary case without this page. After 7 days, explain a new error example to a classmate. After 30 days, rerun the project test, check earlier records and output still match, and transfer the rule to a different task again.

## Next lesson

[Continue: Lists and dictionaries for many tasks](/read/informatics-32-python-collections?lang=en)
