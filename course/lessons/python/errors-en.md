# Errors and exceptions: no file, or rubbish inside it

_Лид (summary):_ **The tenth lesson of the Python course. Three of five rows from somebody else's export made it through: a comma is a full stop, while "no data" and an empty string are not numbers. `try/except/else/finally`, an error name of your own in one line, and the skill that matters — reading a traceback from its last line upwards.**

## Why this is needed

In the previous lesson `assert` checked our own assumptions about the calculation. Today something else begins: data that somebody else sent.

Where your code ends and another party's file begins, an error stops being a sign of a bad programmer and becomes an ordinary event. The file is not there. A word sits where a number should be. The disk is full. The network dropped halfway.

A program that falls over on every such row is useless. A program that silently skips everything is dangerous. What separates them is the ability to **name** the error and decide what to do with it.

## The whole thing at once

The file is `qate.py`. Run it with `python qate.py` from inside the environment.

The required part is the first block: `try`, an `except` with the error's name, `else`, and the accumulator. The second block, about a file and `finally`, is taken apart below; files themselves come in the next lesson.

```python
"""Lesson 10: an error is not the end of the program but a decision to make.

Data from somebody else's file arrives as it is: with a comma instead of a full
stop, with an empty row, with a word where a number should be. The program has
to know what to do with each of them.
"""


class BadRow(ValueError):
    """A row that could not be parsed. Our own name for our own error."""


# This is what an export from someone's table looks like: five rows, three good.
rows = ["8.0", "15,0", "no data", "", "11.4"]


def to_number(text):
    """Text into a number: a comma counts as a full stop, the rest is BadRow."""
    clean = text.strip().replace(",", ".")
    try:
        return float(clean)
    except ValueError:
        raise BadRow(f"not a number: {text!r}")


print("== parsing the rows")
total = 0.0
count = 0
for text in rows:
    try:
        value = to_number(text)
    except BadRow as error:
        print(f"skipped — {error}")
    else:
        total += value
        count += 1
        print(f"taken: {value}")
print(f"parsed {count} of {len(rows)}, sum {total:.1f}")

print()
print("== a missing file and an unreadable one are different errors")
try:
    with open("dannye.csv", encoding="utf-8") as source:
        print(source.readline())
except FileNotFoundError:
    print("no file: we will take the data from the network")
except OSError as error:
    print("the file is there but will not read:", error)
finally:
    print("finally runs whatever happens")
```

It prints:

```
== parsing the rows
taken: 8.0
taken: 15.0
skipped — not a number: 'no data'
skipped — not a number: ''
taken: 11.4
parsed 3 of 5, sum 34.4

== a missing file and an unreadable one are different errors
no file: we will take the data from the network
finally runs whatever happens
```

## Taking it apart

### An error has a name, and the names are a family

`ValueError`, `FileNotFoundError`, `KeyError`, `ZeroDivisionError` are types, like `int` and `str`. And they are not scattered about but arranged in a tree: `FileNotFoundError` is a particular case of `OSError`, and `OSError` a particular case of `Exception`.

A useful consequence follows: catching `OSError` catches "no file", "no permission" and "the disk went away" alike. Catching `Exception` catches everything — including your own typo in a variable's name.

> **Picture it.** A doctor in casualty. "My stomach hurts" and "my arm is broken" are treated by different people, and the first thing done is to name what happened. A diagnosis of "unwell" cannot be treated.

### An `except` with no name is almost always a mistake

```
try:
    ...
except:            ← nobody writes this
    pass
```

A bare `except` catches everything: the error in the data, the `NameError` from a typo, and the Ctrl-C you are pressing to stop the program. Along with the error it swallows the reason — and a `pass` inside says "I do not care what happened".

The rule is simple: **catch what you know how to handle**. Let the rest fall over — falling with a clear message is more honest than quiet nonsense in a report.

### `else` on a `try`: only what can break stays under guard

```python
try:
    value = to_number(text)
except BadRow as error:
    print(f"skipped — {error}")
else:
    total += value
    count += 1
```

There is one line inside `try` — the one that may not work. Everything done **after it succeeds** has moved into `else`.

That is not for tidiness. Put `total += value` inside the `try` and an error in the addition would land there too — and be caught by an `except BadRow` written for something else entirely. `else` holds the line: what is risky goes in `try`, what follows a risk that did not fire goes in `else`.

### The order of the `except` branches matters

```python
except FileNotFoundError:
    ...
except OSError as error:
    ...
```

Python goes down the branches and takes the first that fits. `FileNotFoundError` is a particular case of `OSError`, so the particular one stands **above** the general one. Swap them and the first branch takes everything while the second never runs.

### `finally` runs whatever happens

`finally` fires after success, after an error, and even when the `try` was left through a `return`. It is where you close what has to be closed in any case.

In our example the file needs no closing — `with` does that: it closes the file whatever happens inside. That is why files are almost always opened through `with`, and `finally` is left for what `with` cannot do.

### `raise` and an error name of your own

```python
class BadRow(ValueError):
    """A row that could not be parsed."""
```

One line, and we have an error name of our own. Inheriting from `ValueError` says "this is a particular case of a wrong value": code that catches `ValueError` catches ours too.

Why bother when `ValueError` already exists? Because in parsing an export a `ValueError` can arrive from anywhere — from `float`, from `int`, from somebody's library. `BadRow` says **this row of ours failed to parse**, and it can be caught precisely.

`raise` raises the error. Note where it stands: inside an `except`. Python remembers that and shows both when it falls:

```
ValueError: could not convert string to float: 'no data'

During handling of the above exception, another exception occurred:
...
BadRow: not a number: 'no data'
```

That is convenient: you see both what happened and what we called it. When the second half is in the way, people write `raise BadRow(...) from None`.

### How to read a traceback

A real error looks like this (the paths are shortened):

```
Traceback (most recent call last):
  File "tb.py", line 12, in <module>
    print(average(["8.0", "no data"]))
          ~~~~~~~^^^^^^^^^^^^^^^^^^^^
  File "tb.py", line 8, in average
    total += to_number(text)
             ~~~~~~~~~^^^^^^
  File "tb.py", line 2, in to_number
    return float(text)
ValueError: could not convert string to float: 'no data'
```

It is read **from the bottom up**, and that is the skill this lesson is for.

**The last line** is what happened: the type of the error and its message. Here `float` was handed a string it could not turn into a number, and it named the string outright.

**The lines above** are the road the program travelled to get there: the lowest frame is where it broke, the topmost is where it all started. The `~~~^^^` marks under a line point at the expression at fault when a line holds several.

Hence a habit that saves hours: do not be scared by the length. Look at the last line, then at the nearest frame **with your own file in it** — the mistake is almost always there rather than deep inside somebody's library.

## The map of the lesson

![The map of the lesson: the error's name, the branch and the traceback](/static/course/py/map-errors-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why does `total += value` sit in `else` rather than inside `try`?
2. What happens if `except OSError` is put above `except FileNotFoundError`?
3. Which line of a traceback do you read first, and what does it tell you?

## Exercise

**Required.** Take a list of strings as an export delivers them: a number, a number with a comma, an empty string, a word. Write `to_number(text)` that raises an exception of your own on a row it cannot parse, and a loop that adds the good ones and prints the bad ones with the reason. At the end: how many of how many were parsed.

**Optional.**

- Add an `except ZeroDivisionError` branch to the average and try it on an empty list.
- Write `raise BadRow(...) from None` and compare the traceback with the previous one.
- Open a missing file without a `try` and read the traceback aloud: what happened, where, and where it came from.

## Where this goes in the project

The digest stops falling over on one bad row. Parsing an export runs to the end, bad rows are named one by one, and the report gains a count: so many parsed, so many skipped — and that number is itself a measure of the source's quality.

Debts. Skipped rows are printed to the screen for now. Their place is in a log beside the report, so that tomorrow it is visible what failed to parse yesterday.

## The answers

1. Because only what can break is held inside a `try`. Otherwise an error in the addition lands in an `except` written for parsing a row, and gets explained by the wrong cause.
2. The `FileNotFoundError` branch never runs: `OSError` is the general case, and Python takes the first branch that fits from the top. The particular always stands above the general.
3. The last one: it holds the type of the error and its message — what actually happened. Above it is the chain of calls, and there you look for the nearest frame with your own file.

## Sources

- [Python: errors and exceptions](https://docs.python.org/3/tutorial/errors.html)
- [Python: the exception hierarchy](https://docs.python.org/3/library/exceptions.html#exception-hierarchy)
- [Python: the raise statement](https://docs.python.org/3/reference/simple_stmts.html#the-raise-statement)
- [World Bank: inflation, annual %](https://data.worldbank.org/indicator/FP.CPI.TOTL.ZG)
