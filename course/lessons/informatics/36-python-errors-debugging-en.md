# An error as an observable fact

_Lead (summary):_ **Read errors as evidence and keep tasks safe when input is wrong.**

## Where we are on the map

This is lesson 36 of 72 in the Python block (27–38). The support map shows verifiable transitions. We continue the 0.3 paper contract: five reminder outcomes and the original fictional tasks.

![An error as an observable fact](/static/course/informatics/map-36-python-errors-debugging-en.svg)

## Situation and question

A friend enters `2026-02-29`. The assistant must neither silently turn it into 1 March nor erase its file. Did input, date validation, or saving fail? Find the first point of divergence.

## New words without gaps

An **exception** is Python's signal that an operation cannot continue normally. `ValueError` means an invalid value; `FileNotFoundError` means a missing file. A **traceback** shows calls and the failing line; start at its last line and compare with the source. `try` encloses a risky operation; `except ValueError` catches only an expected value error. **Debugging** means finding a cause from observations and checking the fix. Our **data invariant** is that a failure before saving leaves stored tasks unchanged. A blanket `except:` hides programming defects.

## The lesson's support signal

Input → check the rule → change state → observable result. Point to where the program reads data, compares it with the contract, and merely reports the result. If a step is absent from the map, find it in the code and add it to your own trace.

## Work through it step by step

First write the prediction: an invalid date prints `INVALID DATE` and leaves the file byte for byte unchanged. Run the example and trace its `except ValueError` branch. Then damage one character in a copy of the JSON: that is a `JSONDecodeError`, not an invalid date. `assistant_core` validates the whole document and date before making the changed copy that can be saved. Compare the file hash before and after a failed `add`.

## Predict and check

Cover the output block. Trace names line by line and write the exact predicted text, including case and line order. Only then run the code with `python3` from `step-03` and compare character by character. Change one input and predict again before running. If results differ, find the first divergent line instead of adjusting the prediction afterwards.

```python
from assistant_core import parse_date
try:
    parse_date("2026-02-29")
except ValueError:
    print("INVALID DATE")
```

## Expected output

```text
INVALID DATE
```

## Catch the error

`except Exception: pass` pretends everything worked. If a file cannot be read, showing an empty list as if no tasks exist is false. Distinguish “file missing,” “invalid JSON,” and “no tasks,” and report the actual cause without inventing data.

## Project change

`assistant.py` gives clear messages and a nonzero exit code for expected errors; `assistant_core.py` never saves an unchecked document. Record the failure cases as tests so a later improvement cannot turn an error into lost data.

## Task and evidence

On a copy of `tasks.json`, try an invalid date, duplicate ID, broken JSON, and missing file. Record the expected error, actual message, and whether the file changed for each. Fix one fault and rerun an earlier success case: the new path must not break the old one.

## Transfer to a new setting

A school gradebook will not open: perhaps the password is wrong, perhaps the network is down. Why would “there are zero marks” be dangerous? Suggest two distinct messages and name the data each failure must leave untouched.

## Return after 1, 7, and 30 days

After 1 day, recall the rule and one boundary case without this page. After 7 days, explain a new error example to a classmate. After 30 days, rerun the project test, check earlier records and output still match, and transfer the rule to a different task again.

## Next lesson

[Continue: Tests and modules protect changes](/read/informatics-37-python-tests-modules?lang=en)
