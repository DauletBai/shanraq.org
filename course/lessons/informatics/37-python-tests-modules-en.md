# Tests and modules protect changes

_Lead (summary):_ **Turn paper examples into executable tests and split the program into clear parts.**

## Where we are on the map

This is lesson 37 of 72 in the Python block (27–38). The support map shows verifiable transitions. We continue the 0.3 paper contract: five reminder outcomes and the original fictional tasks.

![Tests and modules protect changes](/static/course/informatics/map-37-python-tests-modules-en.svg)

## Situation and question

Yesterday reminders worked. After adding `done`, the assistant started reminding users about finished tasks. How can we catch this before sharing it? Have the computer repeat earlier checks after every change.

## New words without gaps

A **module** is an importable Python file: `assistant_core.py` holds rules and `assistant.py` commands. An **import** makes a module's name available elsewhere. A **test** runs an input and checks the expected result; `assert` reports disagreement. A **unit test** checks a small function alone; an **end-to-end test** runs a command through a file and inspects the result. A **regression** is a returning bug after a change. A **fixture** is prepared test data. A wrong expected value makes a test useless, so compare it with the paper contract first.

## The lesson's support signal

Input → check the rule → change state → observable result. Point to where the program reads data, compares it with the contract, and merely reports the result. If a step is absent from the map, find it in the code and add it to your own trace.

## Work through it step by step

Take the 0.3 table: completed → `DONE`, no date → `NO_DATE`, past date → `OVERDUE`, today or within two days → `REMIND`, later → `NOT_YET`. Check each case separately; the example shows a minimal `assert`. Change `<= 2` to `< 2`: the two-day boundary test must fail. Restore the rule and rerun every case. The duplicate-ID test is separate from date tests because it protects a different contract.

## Predict and check

Cover the output block. Trace names line by line and write the exact predicted text, including case and line order. Only then run the code with `python3` from `step-03` and compare character by character. Change one input and predict again before running. If results differ, find the first divergent line instead of adjusting the prediction afterwards.

```python
from datetime import date
from assistant_core import reminder_status
task = {"done": False, "due_date": "2026-10-11"}
actual = reminder_status(task, date(2026, 10, 9))
assert actual == "REMIND"
print("Boundary test passed")
```

## Expected output

```text
Boundary test passed
```

## Catch the error

`assert result == result` always passes, even for a wrong result. A test using the real current date can pass today and fail tomorrow. Write the exact expected output first and pass `today` explicitly.

## Project change

Run `tests/test_assistant.py` with `python3 -m unittest discover -s tests -v`. Rules, files, and CLI are checked separately. The three fictional 0.2 records remain source data; new tests use temporary copies rather than replacing them.

## Task and evidence

Write tests for five reminder outcomes, counts of 0 and 3 finished tasks, duplicate ID, invalid February, and preserving Ә. Add one end-to-end `add` → `done` → reload test. Deliberately break a rule, observe the failing test, restore it, and record the command that runs the whole suite.

## Transfer to a new setting

A school quiz changes its scoring. Which old test would catch a correct answer once again receiving zero? Give an input, expected output, and boundary that needs its own test.

## Return after 1, 7, and 30 days

After 1 day, recall the rule and one boundary case without this page. After 7 days, explain a new error example to a classmate. After 30 days, rerun the project test, check earlier records and output still match, and transfer the rule to a different task again.

## Next lesson

[Continue: Release 1.0: a command-line digital assistant](/read/informatics-38-python-cli-release?lang=en)
