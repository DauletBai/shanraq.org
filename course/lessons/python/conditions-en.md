# Conditions: if, elif, else — and the number that is not there

_Лид (summary):_ **The seventh lesson of the Python course. A number on its own says nothing: is 11.4% a lot or a little? A condition compares, and the comparison picks the road. Measured on inflation: Kazakhstan at 11.4% in 2025, the world at 3.0%. Plus the first rule of working with data: a gap in a series is not a zero, and `is None` is how you ask.**

## Why this is needed

So far the program has done the same arithmetic whatever the numbers were. Real work with data starts where the numbers differ and the road depends on them: skip a year that has no value instead of counting it as zero, mark a figure above the norm, set aside a row that arrived as rubbish.

All of that is one construct. It compares and it chooses, and it has three words: `if`, `elif`, `else`.

Today also brings a rule that is worth more than syntax in a course about data: **a missing number and a zero are not the same thing**. A report that counts a gap as zero is wrong silently, and you find out from a reader.

## The whole thing at once

The file is `shart.py`. Run it with `python shart.py` from inside the environment.

```python
"""Lesson 7: conditions, on inflation that has something to be compared with.

A single number means nothing until there is something to compare it against.
A condition is that comparison, after which the program takes one road and not
the other.
"""

# Inflation for the year, %. The figures come from lesson one, the World Bank.
# 2026 has no value: the year is not over. That is a gap, not a zero.
kz = {2022: 15.0, 2023: 14.5, 2024: 8.7, 2025: 11.4, 2026: None}
world = {2022: 8.1, 2023: 5.8, 2024: 3.0, 2025: 3.0}

here = kz[2025]
there = world[2025]

print("== the year 2025")
print(f"Kazakhstan: {here} | the world: {there}")
if here is None:
    print("there is no figure for this year yet")
elif here < there:
    print("below the world's")
elif here < 2 * there:
    print("above the world's, but less than twice")
else:
    print("twice the world's or more")

print()
print("== the year 2026")
future = kz[2026]
if future is None:
    print("no figure yet — nothing to compare with")
else:
    print("there is a figure:", future)

print()
print("== two comparisons in one condition")
if 2.0 <= there < 4.0:
    print(f"the world: {there} — between two and four")
if here > there and kz[2024] > world[2024]:
    print("Kazakhstan: above the world's for the second year running")
```

It prints:

```
== the year 2025
Kazakhstan: 11.4 | the world: 3.0
twice the world's or more

== the year 2026
no figure yet — nothing to compare with

== two comparisons in one condition
the world: 3.0 — between two and four
Kazakhstan: above the world's for the second year running
```

## Taking it apart

### A colon and an indent instead of brackets

```python
if here < there:
    print("below the world's")
```

The colon at the end of the line means "a block follows". The block itself is marked by the **indent** — four spaces. Most languages need curly brackets for this; in Python the empty space on the left does the job.

Hence the one beginner's mistake you can catch by eye: a line drifts to the wrong level and runs at a moment you did not intend. An editor puts the four spaces in for you; mixing them with tabs is not allowed, and Python complains about that separately.

> **Picture it.** A paragraph in a letter. You do not write "paragraph begins" and "paragraph ends" — you indent, and it shows.

### `elif` is not "another if"

The difference people trip over is visible in our own output. The chain is checked from the top down: `11.4` is not `None`, not below the world's, not less than twice the world's. No condition matched, so it reached the `else`, and one line was printed.

An `if — elif — else` chain is checked from the top down and **stops at the first match**. The rest are not checked at all.

Three separate `if` statements are a different thing entirely:

```
if here > there:        → matches
if here > 2 * there:    → matches too
```

Two lines instead of one. That is not a fault of the language; they are different tools. A chain picks **one** of the options; separate `if` statements ask independent questions.

The `else` at the end means "in every other case". It is optional, but without it nothing happens when no condition matched — and that, too, is a decision to make deliberately.

### A gap is not a zero

The most expensive line of the lesson:

```python
if future is None:
```

The series holds `None` for 2026 — not because inflation was zero but because the year has not happened. Those are different things, and in a report they give different answers.

The way to ask is `is None`, not this:

```
if future:          → False for None and for 0.0 alike
```

Remember from the lesson on types: `bool(0.0)` is false. A year of zero inflation would pass such a check as "no data" and vanish from the calculation without a word. `is None` asks exactly what is meant: the value is absent.

The word `is` is not accidental here. It compares not values but **whether this is the same object**, and `None` exists in a single copy in Python — which is what makes `is None` the right check.

### Two comparisons in one condition

```python
if 2.0 <= there < 4.0:
```

It reads as it does in mathematics: from two inclusive up to four. Most languages will not take that and make you write `there >= 2.0 and there < 4.0`; Python understands a chain of comparisons directly and evaluates `there` once.

### `and`, `or`, and the order of the checks

```python
if here > there and kz[2024] > world[2024]:
```

`and` is true when both are; `or` when at least one is; `not` turns it over.

What matters more: **`and` stops at the first false**, `or` at the first true. The right-hand side is not evaluated at all. Hence the habit the rest of the course will need constantly:

```
if value is not None and value > 5:
```

First make sure there is a number, only then compare. Swap the two and the program falls over on a gap, because `None > 5` cannot be compared.

### Why the same thing is written three times

Look at the program again: three blocks, and each of them has another `print` and another comparison. If there were thirty years instead of three, writing it this way would be impossible.

That is what the next two lessons are for: a **loop** will walk all the years by itself, and a **function** will give the comparison a name so it can be called as often as needed. For now it is done by hand, so that what they grow out of is visible.

## The map of the lesson

![The map of the lesson: a fork, a chain and a gap](/static/course/py/map-conditions-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why does an `if — elif — else` chain print one line while three separate `if` statements print two?
2. How does `if value is None` differ from `if not value`, and why does it matter in data?
3. What does `and` do when the left-hand side turns out to be false?

## Warm-up

Three short steps before the exercise: predict, fill in, fix. The answers are at the end of the lesson, but answer them yourself first.

**1. Predict.** What does this program print? Say it out loud, then run it.

<!-- drill 1 -->
```python
value = None
if value == 0:
    print("zero")
elif value is None:
    print("no number")
else:
    print("number:", value)
```

**2. Fill in the gap.** In place of `...` put the condition "differs from the average by no more than 5% either way".

```python
value = 8.3
average = 8.0
if ...:
    print("about the average")
else:
    print("far from the average")
```

**3. Fix it.** This program falls over on the second pair. Find the line that causes it and make both lines print.

```python
kz = {2024: 8.7, 2025: None}
for year, value in kz.items():
    if value > 8:
        print(year, "above eight")
    elif value is None:
        print(year, "no number")
```

## Exercise

**Required.** Given:

```python
prices = {"January": 520.0, "February": 580.0, "March": None, "April": 498.0}
```

Print the average over the months that have a number, and then each month with a label: "about the average" if the price differs from it by no more than 5% (exactly 5% counts as "about"), otherwise "above the average" or "below the average"; a month without a number is "no data". The checks go in this order: the gap first, then "about", then above and below.

The expected output:

<!-- task out -->
```
average: 532.67
January: 520.0 — about the average
February: 580.0 — above the average
March: no data
April: 498.0 — below the average
```

Done when: the output matches line by line; the average is worked out over three months rather than four; no comparison ever meets a `None`.

**On your own data.** Take three numbers of your own — the price of one item in different shops, or your own inflation from your receipts, or anything you can measure. Compare each with the average and print an `if — elif — else` chain: "below average", "about average" (within 5%), "above average". Make one of the three `None` and check for it first, before any comparison.

**Optional.**

- Replace the chain with three separate `if` statements and see how many lines get printed.
- Write a condition with `and` whose right-hand side would fall over if it came first.
- Check it on your own series: how many of the five years were above world inflation?

## Where this goes in the project

The digest stops counting everything indiscriminately. From this lesson on it has judgement: a year with a gap stays out of the average, a figure above the norm is marked, and a row of rubbish is set aside instead of spoiling the total.

Still open. The conditions are still written out one per year — with a loop that becomes three lines. And "above the norm" is decided by eye: where the norm runs is written down nowhere.

## The answers

### To the questions

1. Because the chain is checked from the top down and stops at the first condition that matched: the rest are not checked at all. Separate `if` statements are independent questions, and each of them can match.
2. `is None` asks whether the value is absent. `not value` is true for `None`, for zero and for an empty string alike — a year of zero inflation would disappear from the calculation as "no data".
3. Nothing: the right-hand side is not evaluated. That is why `value is not None and value > 5` is safe while the other order falls over on a gap.

### To the warm-up

1. `no number`. `None == 0` is `False`: a gap is not equal to zero, so the first branch does not match, and the second one asks `is None` and does.

<!-- drill 1 out -->
```
no number
```

2. `abs(value - average) <= average * 0.05`. The absolute difference is what "either way" means, and `<=` includes exactly 5%: the boundary has to be decided once and written down, or two people will get different answers from the same data.

<!-- drill 2 -->
```python
value = 8.3
average = 8.0
if abs(value - average) <= average * 0.05:
    print("about the average")
else:
    print("far from the average")
```

<!-- drill 2 out -->
```
about the average
```

3. The line with `None` is not the culprit, the order is: `value > 8` is checked first, and on 2025 that comparison meets a `None` — `TypeError`. The gap is checked before any comparison:

<!-- drill 3 -->
```python
kz = {2024: 8.7, 2025: None}
for year, value in kz.items():
    if value is None:
        print(year, "no number")
    elif value > 8:
        print(year, "above eight")
```

<!-- drill 3 out -->
```
2024 above eight
2025 no number
```

## Sources

- [Python: the if statement](https://docs.python.org/3/tutorial/controlflow.html#if-statements)
- [Python: comparisons and their chains](https://docs.python.org/3/reference/expressions.html#comparisons)
- [Python: truth value testing](https://docs.python.org/3/library/stdtypes.html#truth-value-testing)
- [World Bank: inflation, annual %](https://data.worldbank.org/indicator/FP.CPI.TOTL.ZG)
