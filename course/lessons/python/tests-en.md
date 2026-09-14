# `pytest` tests: verify the pipeline's promises

_Lead (summary):_ **Lesson fifty-two of the Python course. Turn project rules into automated tests: check the normal case, boundaries, and expected failures, then run the whole suite with one command.**

## Why this matters

A manual run shows that the program worked **once**. A test preserves a project promise and repeats the check after every change. An unknown source must be rejected, a causal claim must wait for review, and an empty report must not overwrite the previous one.

Tests do not prove that no bugs exist. Each test confirms only the example it states. A useful suite therefore covers important behavioural boundaries, not a target percentage of source lines.

## The whole thing first

Create two files:

```python
# policy.py
def route(claim_type, evidence_ids, known_ids, formula_id=None):
    if claim_type not in {"observation", "calculation", "cause"}:
        return "reject"
    if not set(evidence_ids).issubset(known_ids):
        return "reject"
    if claim_type == "cause":
        return "review"
    if claim_type == "calculation" and not formula_id:
        return "review"
    return "accept"
```

```python
# test_policy.py
import pytest

from policy import route


@pytest.mark.parametrize(
    ("claim_type", "evidence", "formula_id", "expected"),
    [
        ("observation", ["cpi-2025"], None, "accept"),
        ("cause", ["cpi-2025"], None, "review"),
        ("calculation", ["cpi-2025"], None, "review"),
        ("calculation", ["cpi-2025"], "change-v1", "accept"),
        ("guess", ["cpi-2025"], None, "reject"),
        ("observation", ["missing"], None, "reject"),
    ],
)
def test_route(claim_type, evidence, formula_id, expected):
    assert route(claim_type, evidence, {"cpi-2025"}, formula_id) == expected
```

Install the development tool and run it:

```console
python -m pip install pytest
python -m pytest -q
```

`6 passed` means six cases passed. `python -m pytest` deliberately uses the same interpreter selected by `python`.

## The shape of a test

The names `test_policy.py` and `test_route` let `pytest` discover the test. Its body follows Arrange–Act–Assert:

1. **Arrange** the input;
2. **Act** by calling one operation;
3. **Assert** that the result matches the promise.

Do not combine unrelated reasons for failure in one test. When it fails, the broken rule should be obvious.

## A table of cases instead of copies

`@pytest.mark.parametrize` runs one test function for every row. When adding a rule, first add an example that fails, change the production code, and run the tests again. This is the short “red → green → refactor” feedback loop.

Avoid a giant table of every possible combination. Choose behaviour classes: a normal case, an empty value, a boundary, an invalid type, and a bug found previously.

## An expected exception

An exception can be the correct result:

```python
import pytest

def percent(value):
    if not 0 <= value <= 100:
        raise ValueError("percentage outside range")
    return value

def test_percent_rejects_101():
    with pytest.raises(ValueError, match="outside range"):
        percent(101)
```

The test fails if no exception is raised or its type differs. Avoid broad `pytest.raises(Exception)`: it can mistake an unrelated crash for correct behaviour.

## Files without litter

The built-in `tmp_path` fixture gives every test an isolated temporary directory:

```python
def test_write_report(tmp_path):
    report = tmp_path / "report.txt"
    report.write_text("ready\n", encoding="utf-8")
    assert report.read_text(encoding="utf-8") == "ready\n"
```

The test does not depend on the working directory or damage a real report. Unit tests should replace the network, current time, and an external API with supplied data or a small test double. Separate integration tests can exercise real boundaries, but they are slower and less stable.

## Test behaviour, not implementation

A brittle test knows which internal helpers were called and in what order. A useful test observes a decision, file, report line, or exception. Internal code may then improve without rewriting tests while its promise remains unchanged.

Every fixed bug deserves a regression test: reproduce the defect, confirm that the test fails against the old code, fix the cause, and keep the test.

## Lesson map

![Lesson map: promise, examples, and result](/static/course/py/map-tests-en.svg)

Recall cue: **promise → normal case + boundary + failure → one run → readable result**.

## Say it in your own words

1. Why does one successful manual run not replace a test?
2. What does `@pytest.mark.parametrize` provide?
3. Why use `tmp_path`?
4. Why is `pytest.raises(Exception)` too broad?

## Warm-up

**1. Predict.** How many test cases run?

<!-- drill 1 -->
```python
import pytest

values = [0, 50, 100]

@pytest.mark.parametrize("value", values)
def test_range(value):
    assert 0 <= value <= 100

print(len(values))
```

**2. Fill the blank.**

```python
import pytest

with pytest.___(ValueError):
    percent(101)
```

**3. Fix it.** This test writes to the real report:

```python
def test_report():
    path = Path("data/report.txt")
    path.write_text("test", encoding="utf-8")
```

## Assignment

**Required.** Create `policy.py` and `test_policy.py`. Cover the six cases above. In a separate test, confirm that percentage validation rejects `-1` and `101` with `ValueError`. Run `python -m pytest -q`.

<!-- task out -->
```text
8 passed
```

**With your own data.** Choose three promises made by your pipeline. Give each a normal case, a boundary, and invalid input.

**Optional.** Test HTML output through `tmp_path` and confirm that user-controlled text is escaped.

## Where this fits in the project

Step 27 gains regression tests for model policy and the report page. A single command can now check a change. The next lesson adds type annotations and `mypy`, catching some mismatches before tests run.

## Answers

1. A manual run preserves neither input, expected output, nor a repeatable check.
2. It runs one rule over several explicit cases.
3. It isolates test files and protects project data.
4. Any unrelated exception could be accepted as the expected outcome.

<!-- drill 1 out -->
```text
3
```

2. The method is `raises`.

<!-- drill 2 -->
```python
import pytest

def percent(value):
    if not 0 <= value <= 100:
        raise ValueError("percentage outside range")

with pytest.raises(ValueError):
    percent(101)
print("ValueError")
```

<!-- drill 2 out -->
```text
ValueError
```

3. Accept `tmp_path` and create the file beneath it.

<!-- drill 3 -->
```python
def test_report(tmp_path):
    path = tmp_path / "report.txt"
    path.write_text("test", encoding="utf-8")
    assert path.read_text(encoding="utf-8") == "test"
```

<!-- drill 3 out -->
```text
```

## Sources

- [pytest: Get Started](https://docs.pytest.org/en/stable/getting-started.html) — test discovery and execution.
- [pytest: parametrization](https://docs.pytest.org/en/stable/how-to/parametrize.html) — tables of cases.
- [pytest: temporary directories](https://docs.pytest.org/en/stable/how-to/tmp_path.html) — isolated temporary files.
