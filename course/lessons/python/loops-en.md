# Loops: for, range, and what zip does silently

_Лид (summary):_ **The eighth lesson of the Python course. One comparison instead of thirty: a loop walks the series itself, `continue` steps over the year with no data, an accumulator gives the average — 11.52% across five years. Plus two traps: `zip` silently cuts to the shorter series, and a `while` without a way out walks off the edge.**

## Why this is needed

In the previous lesson the same comparison was written out three times. Five years would be fifteen lines; thirty years would be a program nobody can read or correct.

A loop removes the repetition entirely: the rule is written once and the series may hold as many rows as it likes. Real work with data starts here — real data has thousands of numbers, not three.

Today also brings two places where a loop is wrong without saying so. Both are measured on our own series, and both are worth remembering at once.

## The whole thing at once

The file is `cikl.py`. Run it with `python cikl.py` from inside the environment.

The program is longer than the previous ones — forty-four lines. The required part is the first block: walking the series, `continue`, and the accumulator. The other three show `zip`, `enumerate`, `break` and `while`; those are enough to recognise by sight and come back to when you need them.

```python
"""Lesson 8: a loop — the program walks the series, not you.

In the previous lesson one comparison was written three times. Here it is
written once, and there may be thirty years instead of three.
"""

# Inflation for the year, %: Kazakhstan and the world. Figures from lesson one,
# the World Bank.
kz = {2021: 8.0, 2022: 15.0, 2023: 14.5, 2024: 8.7, 2025: 11.4, 2026: None}
world = {2021: 3.5, 2022: 8.1, 2023: 5.8, 2024: 3.0, 2025: 3.0}

print("== the whole series, the gap left out")
total = 0.0
count = 0
for year, value in kz.items():
    if value is None:
        print(f"{year}: no data")
        continue
    total += value
    count += 1
    print(f"{year}: {value:5.1f}%")
print(f"years with a figure: {count}, average: {total / count:.2f}%")

print()
print("== two series side by side")
for number, (here, there) in enumerate(zip(kz.values(), world.values()), start=1):
    print(f"{number}. Kazakhstan {here:5.1f} | the world {there:4.1f}")

print()
print("== the first year above ten per cent")
for year in range(2021, 2026):
    if kz[year] > 10:
        print(f"{year}: {kz[year]}%")
        break

print()
print("== how many years running, from the end, above the world's")
year = 2025
streak = 0
# The condition also checks that the year is in the series at all: without it
# the loop steps into 2020 and falls over with KeyError. A for ends by itself;
# a while ends only when you say so.
while year in world and kz[year] > world[year]:
    streak += 1
    year -= 1
print(f"years running: {streak}, the last such year: {year + 1}")
```

It prints:

```
== the whole series, the gap left out
2021:   8.0%
2022:  15.0%
2023:  14.5%
2024:   8.7%
2025:  11.4%
2026: no data
years with a figure: 5, average: 11.52%

== two series side by side
1. Kazakhstan   8.0 | the world  3.5
2. Kazakhstan  15.0 | the world  8.1
3. Kazakhstan  14.5 | the world  5.8
4. Kazakhstan   8.7 | the world  3.0
5. Kazakhstan  11.4 | the world  3.0

== the first year above ten per cent
2022: 15.0%

== how many years running, from the end, above the world's
years running: 5, the last such year: 2021
```

## Taking it apart

### `for` takes the items, it does not count the numbers

```python
for year, value in kz.items():
```

In most languages a loop sets up a counter and reaches in by index. Python's `for` takes **the items themselves**: values from a list, keys from a dictionary, and from `.items()` the pairs, which are unpacked into two variables on the spot. That is the same unpacking as in the lesson on tuples.

Hence a rule that saves half the mistakes: if you want an item, take the item, not its number.

> **Picture it.** A stack of receipts in your hand. You do not ask for "receipt number four" — you take the next one until the stack runs out.

### `continue` steps over, `break` stops

```python
if value is None:
    print(f"{year}: no data")
    continue
```

`continue` abandons the current turn and moves to the next item. That is exactly what the gap from the previous lesson needs: the year with no figure gets a line of its own and stays out of the average.

`break` is different: it leaves the loop altogether. In the third block it stops the walk at the very first year above ten per cent, so 2023 to 2025 are never examined.

### An accumulator is set up before the loop

```python
total = 0.0
count = 0
for ...:
    total += value
    count += 1
```

Two variables outside, `+=` inside — that is what counting across a series looks like. Setting them up **before** the loop is not optional: inside it, they would be reset on every turn.

We count two numbers rather than one: the sum and the count. Because the division is not by the length of the series but by the number of years that have a figure — the gap is left out, and that is why the average came out at `11.52` rather than `9.6`.

### `zip` silently cuts to the shorter one

Look at the second block of output: five rows, although `kz` holds six years.

```python
for number, (here, there) in enumerate(zip(kz.values(), world.values()), start=1):
```

`zip` joins the series in pairs and **stops at whichever ran out first**. We have six values in `kz` and five in `world`, so 2026 disappeared from the output. No error, no warning.

Here that suits us: 2026 has no figure anyway. But the habit is dangerous: if two series come from different sources and one is shorter, `zip` quietly cuts the tail off and the report comes out on incomplete data. When the length matters, it is checked outright — `len(a) == len(b)` — or you take `zip(a, b, strict=True)`, which raises on a mismatch instead of saying nothing.

And the thing that matters more than the cut: `zip` joins **by position, not by year**. Here the keys of both series run in order, so the pairs came out right. Let the series come from different sources and the order of the keys may differ, and `zip` will quietly pair 2021 with 2022. Neither equal length nor `strict=True` saves you from that: they count elements, they do not match keys.

Hence the rule: **two series are compared by key**, not by position. That is what the years are for: `for year in kz: ... world[year]`. `zip` is for the places where the position is itself the meaning — as in our output, where two numbers for the same year are printed side by side because both series were taken from one place, in one order.

### `enumerate` counts the turns for you

```python
enumerate(..., start=1)
```

If you need a row number, do not keep a counter by hand. `enumerate` hands you the pair "number, item", and `start=1` says that people read a numbering that begins at one rather than at zero.

### `range` is a recipe, not a list

```python
for year in range(2021, 2026):
```

`range(2021, 2026)` gives 2021, 2022, 2023, 2024 and 2025 — **the right edge is not included**, as with the slices in the lesson on lists. It is the same rule, and there is nothing new to memorise.

What does matter: `range` does not build a list in memory. `range(1_000_000)` takes as much room as `range(5)` — it hands out numbers one at a time. This is the first meeting with laziness; the lesson on generators takes it further.

### A `while` ends only when you say so

```python
while year in world and kz[year] > world[year]:
```

A `for` ends by itself: the items run out and the loop is over. A `while` turns for as long as its condition holds, and watching for the way out is left to you.

The first version of this program was written without `year in world` — and fell over:

```
KeyError: 2020
```

All five years turned out to be above world inflation, so the loop reached 2021, stepped into 2020, which is not in the series, and asked for a key that does not exist. The `year in world` check is the edge the loop has to stop at.

## The map of the lesson

![The map of the lesson: the walk, the skip and the edge](/static/course/py/map-loops-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. How does `continue` differ from `break`?
2. Why is the average worked out by dividing by `count` rather than by the length of the series?
3. What does `zip` do when one series is shorter than the other, and why is that dangerous?

## Warm-up

Three short steps before the exercise: predict, fill in, fix. The answers are at the end of the lesson, but answer them yourself first.

**1. Predict.** What does this program print?

<!-- drill 1 -->
```python
total = 0
for value in [8.0, None, 15.0]:
    if value is None:
        continue
    total += value
print(total)
```

**2. Fill in the gap.** In place of `...` put what prints the years from 2021 to 2025 inclusive.

```python
for year in ...:
    print(year, end=" ")
```

**3. Fix it.** The program prints two years out of three. Find where the first one went.

```python
years = [2021, 2022, 2023]
for i in range(1, len(years)):
    print(years[i])
```

## Exercise

**Required.** Given:

```python
prices = [520.0, 546.0, None, 498.0, 515.0]
```

Walk the list and print the number of the month (counting from one) and its price; print a month without a number on its own line and leave it out of the count. At the end print how many months have a number and the average over those. Then find the **first** month above 530 and print its number — and stop looking.

The expected output:

<!-- task out -->
```
1: 520.0
2: 546.0
3: no data
4: 498.0
5: 515.0
months with a number: 4, average: 519.75
first above 530: month 2
```

Done when: the output matches line by line; the gap got into neither the sum nor the counter; the search for the first month stops itself instead of walking the list to the end.

**On your own data.** Take a series of five to seven numbers of your own — the price of one item month by month, your own spending, anything you can measure — and put a `None` in one place. Walk it with a loop: print the gap on a line of its own and leave it out of the count; for the rest print the value and work out the sum, the count and the average.

**Optional.**

- Use `break` to find the first month whose figure went above the average.
- Put two of your own series through `zip` and make them different lengths on purpose — see how many rows get printed.
- Replace `zip(a, b)` with `zip(a, b, strict=True)` and read what Python says.
- Swap two pairs around in the second series, keeping its length, and look at the output: the same numbers, different pairs, no error.

## Where this goes in the project

The digest stops depending on how many years the series holds. One and the same walk works out the average, finds the maximum and sets the gaps aside — over five years or over thirty.

Debts. The walk is still written inside the program as one lump: to count the same thing over a second series it would have to be copied. The next lesson gives that lump a name — a function.

## The answers

### To the questions

1. `continue` abandons the current turn and moves to the next item; the loop goes on. `break` leaves the loop altogether and the remaining items are never examined.
2. Because the series holds a year with no figure. Dividing by the length of the series would mean counting the gap as a zero — which is what the previous lesson was about.
3. `zip` stops at the shorter series and silently drops the tail of the longer one — no error, no warning. If the series came from different sources, the report comes out on incomplete data; so either the lengths are checked or `strict=True` is used.

### To the warm-up

1. `23.0`. `continue` drops the current turn and goes to the next value, so the `None` never reaches the sum while `8.0 + 15.0` do.

<!-- drill 1 out -->
```
23.0
```

2. `range(2021, 2026)`. The right-hand bound is not included, so the last year is written one higher than the one you want — this is the first place a beginner loses a year.

<!-- drill 2 -->
```python
for year in range(2021, 2026):
    print(year, end=" ")
```

<!-- drill 2 out -->
```
2021 2022 2023 2024 2025
```

3. `range(1, len(years))` starts at the second element: a list counts from zero. When it is the element you want rather than its number, walk the list directly:

<!-- drill 3 -->
```python
years = [2021, 2022, 2023]
for year in years:
    print(year)
```

<!-- drill 3 out -->
```
2021
2022
2023
```

## Sources

- [Python: the for statement and range](https://docs.python.org/3/tutorial/controlflow.html#for-statements)
- [Python: break, continue and else on loops](https://docs.python.org/3/tutorial/controlflow.html#break-and-continue-statements)
- [Python: zip](https://docs.python.org/3/library/functions.html#zip)
- [Python: enumerate](https://docs.python.org/3/library/functions.html#enumerate)
- [World Bank: inflation, annual %](https://data.worldbank.org/indicator/FP.CPI.TOTL.ZG)
