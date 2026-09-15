# Type annotations and `mypy`: fail before running

_Lead (summary):_ **Lesson fifty-three of the Python course. Describe data shapes in signatures, narrow `None` with a check, and run `mypy` so incompatible values are found before the program executes.**

## Why this matters

Python remains dynamic: an annotation does not stop a caller passing text instead of a number at runtime. An editor and `mypy` can trace that value beforehand. Types are most valuable at module boundaries: parameters, return values, and record structures.

## The whole thing first

```python
from dataclasses import dataclass

@dataclass(frozen=True)
class Observation:
    source_id: str
    year: int
    value: float

def latest(rows: list[Observation]) -> Observation | None:
    if not rows:
        return None
    return max(rows, key=lambda row: row.year)

def label(row: Observation | None) -> str:
    if row is None:
        return "no data"
    return f"{row.year}: {row.value:.1f}%"

print(label(latest([Observation("cpi-kz", 2025, 11.4)])))
```

Output:

```text
2025: 11.4%
```

The `-> Observation | None` signature admits an empty list honestly. After `if row is None`, the checker narrows the type: below it, `row` is an `Observation`.

## Run `mypy`

Install the course version and check the file:

```console
python -m pip install mypy==2.3.1
python -m mypy --strict digest.py
```

With `Observation("cpi-kz", "2025", 11.4)`, the checker reports that argument two must be `int`. Python may still run that code; static analysis and runtime validation are different layers.

Start with your modules rather than every dependency. A useful minimal `pyproject.toml` is:

```toml
[tool.mypy]
python_version = "3.12"
strict = true
files = ["sholu", "main.py"]
```

This is the oldest Python version promised by the project, not whichever version happens to be on the author's machine.

## Collections and readable aliases

```python
SourceID = str
Evidence = dict[SourceID, float]

def missing(required: set[SourceID], known: Evidence) -> set[SourceID]:
    return required - known.keys()
```

An alias communicates intent, but `SourceID = str` does not create a distinct type. If confusing two strings is dangerous, introduce `NewType` later.

## Do not hide uncertainty in `Any`

`Any` disables checking for a value and spreads downstream. A JSON response is genuinely unknown at the boundary: accept it as `object`, validate its structure, and then construct a typed record. Do not label the whole application `dict[str, Any]` merely to silence errors.

An annotation does not validate a range, unit, or existing `source_id`. Tests and runtime validation still own those promises. Types complement lesson 52; they do not replace it.

## Lesson map

![Lesson map: value, annotation, and static check](/static/course/py/map-typing-en.svg)

Recall cue: **boundary → precise type → narrowing → `mypy` → runtime test**.

## Say it in your own words

1. Why is an annotation not runtime validation?
2. What does `Observation | None` mean?
3. Why check `None` before accessing fields?
4. Why confine `Any` to a system boundary?

## Warm-up

**1. Predict the output.**

<!-- drill 1 -->
```python
def size(values: list[int]) -> int:
    return len(values)

print(size([3, 5]))
```

**2. Fill the return annotation.**

```python
def title(year: int) ... str:
    return f"Report {year}"
```

**3. Fix it.** `row` may be `None`:

```python
def show(row: Observation | None) -> str:
    return str(row.year)
```

## Assignment

**Required.** Type `latest` and `label`, then check a populated and an empty list. Run `python -m mypy --strict solution.py`, followed by the file itself.

<!-- task out -->
```text
2025: 11.4%
no data
```

**With your data.** Annotate three public digest functions and fix errors without `# type: ignore`.

**Optional.** Introduce `NewType("SourceID", str)` and find a swapped identifier.

## Where this fits in the project

Step 28 adds types at the model boundary and strict `mypy` configuration. The next lesson pins Python and dependencies so the same analysis and tests can run on another machine.

## Answers

1. Python normally does not enforce annotations; a checker reads them.
2. The function returns a record or no result.
3. The check narrows the type and prevents field access on no object.
4. `Any` disables useful checks and propagates uncertainty.

<!-- drill 1 out -->
```text
2
```

2. Use `->`.

<!-- drill 2 -->
```python
def title(year: int) -> str:
    return f"Report {year}"

print(title(2025))
```

<!-- drill 2 out -->
```text
Report 2025
```

3. Handle `None` first.

<!-- drill 3 -->
```python
def show(row: Observation | None) -> str:
    if row is None:
        return "no data"
    return str(row.year)
```

<!-- drill 3 out -->
```text
```

## Sources

- [Python: typing](https://docs.python.org/3/library/typing.html) — standard annotations.
- [mypy: Getting started](https://mypy.readthedocs.io/en/stable/getting_started.html) — static checking.
- [mypy: The `Any` type](https://mypy.readthedocs.io/en/stable/dynamic_typing.html) — dynamic typing boundaries.

