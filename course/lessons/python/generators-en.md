# Iterators and generators: four megabytes against four hundred bytes

_Лид (summary):_ **The sixteenth lesson of the Python course. A hundred thousand squares take almost four megabytes as a list and 376 bytes as a generator — measured, with every number counted in. The difference is not in how you count them but in when the values appear. Plus the mistake that matters: a generator is good for one pass, and the second one quietly gives you nothing.**

## Why this is needed

Lesson five ended on a promise: list comprehensions come in the sixteenth. Lesson eleven added the debt about a gigabyte file, lesson twelve the one about an export of hundreds of thousands of rows. All of it is one question, and today it is settled.

The question is **when the values appear**. A list makes all of them at once and holds them in memory. A generator makes them one at a time, when asked, and holds nothing.

While the data is small the difference is invisible. When it is large, it decides whether the program works or falls over for want of memory.

## The whole thing at once

The file is `agyn.py`. Run it with `python agyn.py` from inside the environment.

The required part is the second and third blocks: a generator with `yield`, and the difference between a list and a generator. The first shows what `for` has been doing all along, the fourth shows the mistake the lesson is worth finishing for.

```python
"""Lesson 16: a generator is data that is not in memory yet.

A file of a hundred thousand rows will fit in memory; one of a gigabyte will
not. The difference is not in how you read it but in when the values appear:
all at once, or one at a time.
"""

import sys
import tracemalloc
from pathlib import Path

HERE = Path(__file__).parent
data = HERE / "vygruzka.csv"

# A practice export: a hundred thousand rows of "year;value".
with data.open("w", encoding="utf-8") as target:
    for year in range(2000, 2100):
        for number in range(1000):
            target.write(f"{year};{number % 20 + 1}.5\n")

print("== what a for loop actually does")
numbers = [1, 2, 3]
step = iter(numbers)
print("one at a time:", next(step), next(step), next(step))
try:
    next(step)
except StopIteration:
    print("StopIteration — this is where a for loop simply ends")

print()
print("== a generator hands the values over one at a time")


def values(path):
    """Yields (year, value) one at a time, without reading the whole file."""
    with path.open(encoding="utf-8") as source:
        for line in source:
            year, _, text = line.strip().partition(";")
            yield int(year), float(text)


stream = values(data)
print("type:", type(stream).__name__)
print("the first value:", next(stream))
print("the second value:", next(stream))

print()
print("== the brackets decide: a list or a generator")
# getsizeof measures the object itself, and a list is only references to the
# numbers. What is actually taken up is counted by tracemalloc. We print
# megabytes: tracemalloc's exact byte differs from machine to machine, while
# "four megabytes" is the same everywhere.
tracemalloc.start()
base = tracemalloc.get_traced_memory()[0]
squares_list = [number * number for number in range(100_000)]
list_all = tracemalloc.get_traced_memory()[0] - base
squares_gen = (number * number for number in range(100_000))
gen_all = tracemalloc.get_traced_memory()[0] - base - list_all
tracemalloc.stop()
print(f"list:      the object itself {sys.getsizeof(squares_list)} bytes, with every number {list_all / 1_000_000:.1f} MB")
print(f"generator: the object itself {sys.getsizeof(squares_gen)} bytes, with all it holds {gen_all} bytes")
print("the sum is the same:", sum(squares_list) == sum(squares_gen))

print()
print("== a generator is good for one pass")
stream = values(data)
total = 0.0
count = 0
for _, value in stream:
    total += value
    count += 1
print(f"rows: {count}, sum: {total:.1f}")
try:
    print(max(value for _, value in stream))
except ValueError as error:
    print("a second pass over the same generator:", error)

data.unlink()
```

It prints:

```
== what a for loop actually does
one at a time: 1 2 3
StopIteration — this is where a for loop simply ends

== a generator hands the values over one at a time
type: generator
the first value: (2000, 1.5)
the second value: (2000, 2.5)

== the brackets decide: a list or a generator
list:      the object itself 800984 bytes, with every number 4.0 MB
generator: the object itself 208 bytes, with all it holds 376 bytes
the sum is the same: True

== a generator is good for one pass
rows: 100000, sum: 1100000.0
a second pass over the same generator: max() iterable argument is empty
```

## Taking it apart

### `for` has been calling `next` all along

```
one at a time: 1 2 3
StopIteration — this is where a for loop simply ends
```

`for` is not magic. It takes an iterator from the object (`iter`), then pulls `next` until a `StopIteration` arrives, and that is where it ends.

Hence something worth knowing: **`for` does not work with lists, it works with anything that can hand values over one at a time**. A file, a dictionary, a `range`, a string, a generator — all of them are walked the same way, and `range(1_000_000)` does not build a million numbers, it counts them along the way.

> **Picture it.** A queue at a till. The cashier does not ask how many people there are altogether; they serve the next one while there is a next one.

### `yield` — a function that hands over one value and remembers its place

```python
def values(path):
    with path.open(encoding="utf-8") as source:
        for line in source:
            ...
            yield int(year), float(text)
```

A function with a `yield` in it is a generator function. Calling it computes nothing: `values(data)` returns a generator object, and at that moment the file is not even open. The values appear at the first `next` — and one at a time after that.

`return` ends a function; `yield` **suspends** it: at the next `next` it carries on from the same place, with the same variables and the same open file. That is how a generator reads a gigabyte file in one row of memory: there is only ever the current row in it.

Such a generator is also the place for the sieve. A row that could not be parsed need never be handed outwards at all, and whoever reads the generator need not know about it.

### The brackets decide: a list or a generator

```python
squares_list = [number * number for number in range(100_000)]
squares_gen = (number * number for number in range(100_000))
```

Square brackets make a **list comprehension**: the short form of a loop that collects a list. `[x * 2 for x in prices]` is a `for` with an `append`, written on one line. There is a form with a condition too: `[x for x in prices if x > 500]`.

Round brackets make a **generator expression**: the same thing, except the values are not collected but handed over one at a time. The difference shows in the measurement:

```
list:      the object itself 800984 bytes, with every number 4.0 MB
generator: the object itself 208 bytes, with all it holds 376 bytes
```

There are two measurements here and both are needed. `sys.getsizeof` shows **the object itself**: for the list that is 800 kilobytes — a hundred thousand references and nothing more, because the numbers themselves lie elsewhere. `tracemalloc` counts **everything that was allocated**, and there the truth shows: some four megabytes against three hundred and seventy-six bytes.

The difference is nearly ten thousandfold, and only one side of it grows with the data. A generator takes the same room over any number of values: it holds not the values but the place it stopped at — its own local variables and the current step. "It holds nothing" would be untrue; the truth is that it holds no **finished result**.

The rule for choosing is simple: **if the result is needed whole and more than once, a list; if it is needed once and in passing, a generator**. `sum(x * x for x in range(100_000))` builds no list at all.

### A generator is good for one pass

```
rows: 100000, sum: 1100000.0
a second pass over the same generator: max() iterable argument is empty
```

Here is the mistake the lesson is worth finishing for. A generator runs out. Walk it once and you have spent it, and a second pass gives you **nothing** — not an error, zero elements.

Here `max` complained out loud, because nothing has no maximum. But `sum` over nothing gives `0`, `list` gives `[]`, and a loop simply does not run. The report comes out full of zeroes and nobody says why.

The cure is one of two decisions: either **everything needed is counted in one pass** (as in the example: the sum and the count together), or the generator is made afresh — `values(data)` — and you pay for a second reading of the file. Which is cheaper depends on the size of the file; which is honest is always visible in the code.

### What we are not taking today

`itertools` is the library for streams: `islice` takes the first N, `chain` joins, `groupby` groups what runs together. All of it works with generators and all of it one value at a time. Knowing the module exists is enough for today.

## The map of the lesson

![The map of the lesson: the list, the generator and the single pass](/static/course/py/map-generators-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. What does `for` do underneath, and why can a file be walked the way a list is?
2. How does `[x for x in ...]` differ from `(x for x in ...)`, and when is each taken?
3. Why does a second pass over a generator give nothing, and what is done about it?

## Warm-up

Three short steps before the exercise: predict, fill in, fix. The answers are at the end of the lesson, but answer them yourself first.

**1. Predict.** What does this program print? One and the same expression works out both sums.

<!-- drill 1 -->
```python
squares = (n * n for n in range(3))
print(sum(squares), sum(squares))
```

**2. Fill in the gap.** In place of `...` put the expression that doubles each price.

```python
prices = [260, 620, 1890]
doubled = [... for price in prices]
print(doubled)
```

**3. Fix it.** The program falls over: a generator has no length. Count them without collecting a list.

```python
values = (n for n in range(5))
print(len(values))
```

## Exercise

**Required.** Given:

```python
rows = ["2024;8.7", "2025;n/a", "2025;11.4", "2026;15.0", "rubbish", "2026;12.0"]
```

Write these rows into a file beside the program. Write a generator `values(path)` that hands `(year, value)` over one at a time and silently skips the rows that do not parse. Take the first two values out of it with `next` and print them. Then, **in a single pass**, work out how many rows were taken and the average to two decimal places. At the end print the years after 2024 with a list comprehension. Take the file away after you.

The expected output:

<!-- task out -->
```
the first: (2024, 8.7)
the second: (2025, 11.4)
rows taken: 4, average: 11.78
years after 2024: [2025, 2026, 2026]
```

Done when: the output matches line by line; there is a `yield` in the program and the file is nowhere read whole; the sum and the count are worked out in one pass; the unusable rows are sieved out inside the generator rather than after it.

**On your own data.** Take a file of your own — an export, a log, anything line by line. Write a generator that hands out only what you need from it, and work out one figure over that. Then say out loud how many rows of the file were in memory at once.

**Optional.**

- Compare `sys.getsizeof` for a list and for a generator over a million elements.
- Build `[x for x in range(20) if x % 3 == 0]` and the same thing as a generator, and compare the types.
- Take `itertools.islice(values(path), 3)` and see how many rows of the file were read.

## Where this goes in the project

The digest stops depending on the size of the export. It used to read the file into a list and count over that; now the rows flow through a generator and there is always one of them in memory. The same code will work over a hundred rows and over a gigabyte.

Debts. One pass is a discipline: everything the data is needed for has to be decided in advance. And we named `itertools` without taking it apart; we will come back when the digest has several sources to join.

## The answers

### To the questions

1. It takes an iterator with `iter` and pulls `next` until a `StopIteration` arrives. That is why `for` works the same over anything that can hand values over one at a time — a list, a file, a dictionary, a generator.
2. Square brackets collect a list at once and hold it in memory; round ones give a generator that hands values over one at a time and stores nothing. A list is taken when the result is needed whole and more than once; a generator when it is needed once and in passing.
3. Because a generator remembers the place it stopped at, and it stopped at the end. A second pass is not an error but emptiness — which is more dangerous than an error. Either everything is counted in one pass, or the generator is made afresh.

### To the warm-up

1. `5 0`. The first sum walked the generator to the end — `0 + 1 + 4`. By the second there was nothing left in it, and `sum` over nothing gives zero without a word.

<!-- drill 1 out -->
```
5 0
```

2. `price * 2`. In a comprehension what stands to the left of the `for` is what goes into the list, and what stands to the right is where it comes from.

<!-- drill 2 -->
```python
prices = [260, 620, 1890]
doubled = [price * 2 for price in prices]
print(doubled)
```

<!-- drill 2 out -->
```
[520, 1240, 3780]
```

3. `TypeError: object of type 'generator' has no len()`. A generator has no length: it does not know itself how many values it will hand over. They are counted along the way:

<!-- drill 3 -->
```python
values = (n for n in range(5))
print(sum(1 for _ in values))
```

<!-- drill 3 out -->
```
5
```

## Sources

- [Python: generators in the tutorial](https://docs.python.org/3/tutorial/classes.html#generators)
- [Python: the yield expression](https://docs.python.org/3/reference/expressions.html#yieldexpr)
- [Python: list and generator comprehensions](https://docs.python.org/3/tutorial/datastructures.html#list-comprehensions)
- [Python: the itertools module](https://docs.python.org/3/library/itertools.html)
