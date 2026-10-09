# Lists and dictionaries for many tasks

_Lead (summary):_ **Collect tasks in a list, read dictionary fields, and reject a duplicate `id` before it corrupts an answer.**

## Where we are on the map

This is lesson 32 of 72 in the Python block (27–38). The support map shows verifiable transitions. We continue the 0.3 paper contract: five reminder outcomes and the original fictional tasks.

![Lists and dictionaries for many tasks](/static/course/informatics/map-32-python-collections-en.svg)

## Situation and question

On paper each task was a separate card. In a file it becomes fields: `id`, `title`, `done`, later `due_date`. How do we find one when titles repeat? Why is the same `id` on two different cards dangerous? Choose data structures by the job they must do.

## New words without gaps

A **list** `[]` keeps tasks in order and lets us visit them all. A **dictionary** `{}` maps a key to a value: `task["id"]` reads one task's field. Here a **key** is a field name; do not confuse it with a task's unique identifier. An **index** is a list position starting at zero. A **set** `set()` stores unique values and catches a repeated `id`. A dictionary indexed by task ID, `by_id`, makes lookup convenient, but building it before checking duplicates silently overwrites an earlier record.

## The lesson's support signal

Input → check the rule → change state → observable result. Point to where the program reads data, compares it with the contract, and merely reports the result. If a step is absent from the map, find it in the code and add it to your own trace.

## Work through it step by step

Lay out two cards: `id=t-01, done=False` and `id=t-02, done=True`. The list preserves their order; `by_id` maps each ID to its card. The code below prints `True` for `t-02`. Add a third card also named `t-02`: a plain dictionary construction hides one of them. Walk the list with a `seen` set first and stop on the duplicate.

## Predict and check

Cover the output block. Trace names line by line and write the exact predicted text, including case and line order. Only then run the code with `python3` from `step-03` and compare character by character. Change one input and predict again before running. If results differ, find the first divergent line instead of adjusting the prediction afterwards.

```python
tasks = [{"id": "t-01", "done": False}, {"id": "t-02", "done": True}]
for task in tasks:
    print(task["id"], task["done"])
```

## Expected output

```text
t-01 False
t-02 True
```

## Catch the error

Two tasks may both be titled “Кітап оқу” and have different IDs; that is not a duplicate. Two records with one ID and different titles are duplicates. Checking titles fails in both directions. Another trap: `tasks[1]` means the second item, not the task whose `id` is 1.

## Project change

Move the three original tasks into a list of dictionaries without changing their values. Validate uniqueness before counting or reminding: bad input must not produce a convincing-looking answer.

## Task and evidence

Make a list of the three original tasks and a duplicate-ID check using `seen`. Test the empty list, the original three, `t-01,t-02,t-01`, and two equal titles with distinct IDs. Draw the set contents after each step. Build `by_id` only after validation succeeds.

## Transfer to a new setting

A library can hold two books with one title but different inventory IDs. Which structure preserves checkout order, and which finds a book by ID quickly? Where will you reject a repeated ID?

## Return after 1, 7, and 30 days

After 1 day, recall the rule and one boundary case without this page. After 7 days, explain a new error example to a classmate. After 30 days, rerun the project test, check earlier records and output still match, and transfer the rule to a different task again.

## Next lesson

[Continue: Text, dates, and Kazakh Unicode](/read/informatics-33-python-text-dates?lang=en)
