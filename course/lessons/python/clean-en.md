# Dirty data: gaps, duplicates and types

_Лид (summary):_ **The thirtieth lesson of the Python course. Data arrives broken: numbers as text, a city spelled three ways, a twin row and an empty cell. `isna`, `to_numeric` with `errors="coerce"`, `astype`, `drop_duplicates` — and the question the lesson turns on: drop a gap or fill it, and what filling it costs.**

## Why this matters

The tables in these lessons have been clean so far, because we wrote them. Real ones arrive differently: an amount is text with a space inside it, a city is spelled three different ways, one row doubled during the export, and a cell says "n/a".

Cleaning is not a one-off chore but part of the program, and it has a rule: **every step leaves a trace**. How many rows there were, how many there are, what exactly changed. A table nobody can say what was done to is worse than a dirty one.

## The whole thing first

The file is `uborka.py`. Eight rows in which everything that usually breaks is broken.

```python
"""Lesson 30: repairing what arrived broken.

Eight rows in which everything that usually breaks is broken: gaps, a twin row,
numbers as text, a city spelled three ways. We mend them one at a time and count
what changed.
"""

import pandas as pd

rows = [
    ("Kostanay", "food", "4 200", "2026-01-03"),
    ("kostanay ", "fuel", "18 500", "2026-01-05"),
    ("Rudny", "food", "2 600", "2026-01-07"),
    ("Rudny", "food", "2 600", "2026-01-07"),
    ("KOSTANAY", "phone", "n/a", "2026-01-09"),
    ("Rudny", "rent", "62 000", "2026-01-11"),
    ("Kostanay", "rent", "95 000", None),
    ("Astana", "food", "7 300", "2026-01-15"),
]
df = pd.DataFrame(rows, columns=["city", "kind", "amount", "day"])
print(df)
print("types:", dict(df.dtypes.astype(str)))

print()
print("== what is broken: look before mending")
print("gaps per column:", df.isna().sum().to_dict())
print("whole duplicates:", int(df.duplicated().sum()))
print("cities before:", sorted(df["city"].unique()))

print()
print("== numbers: text into a number, junk into a gap")
df["amount"] = pd.to_numeric(df["amount"].str.replace(" ", "", regex=False), errors="coerce")
print("the type now:", df["amount"].dtype, "| gaps:", int(df["amount"].isna().sum()))

print()
print("== dates: the same move, errors='coerce' once more")
df["day"] = pd.to_datetime(df["day"], errors="coerce")
print("the type now:", df["day"].dtype, "| gaps:", int(df["day"].isna().sum()))

print()
print("== cities: three spellings, one city")
df["city"] = df["city"].str.strip().str.capitalize()
print("cities after:", sorted(df["city"].unique()))

print()
print("== duplicates: count them, then drop them")
before = len(df)
df = df.drop_duplicates()
print(f"rows: {before} -> {len(df)}")

print()
print("== gaps: dropping or filling is a decision, not a habit")
print("rows with no amount:", int(df["amount"].isna().sum()))
paid = df.dropna(subset=["amount"])
print("counting over:", len(paid), "rows | sum:", int(paid["amount"].sum()))
print("and filled with zero:", int(df["amount"].fillna(0).sum()), "-- the same sum, but the average now lies")
print("average over what is there:", round(paid["amount"].mean(), 1),
      "| with the zero:", round(df["amount"].fillna(0).mean(), 1))

print()
print("== a category instead of strings: same table, less memory")
print("as strings:", int(df["city"].memory_usage(deep=True)),
      "| as a category:", int(df["city"].astype("category").memory_usage(deep=True)))
```

It prints:

```text
        city   kind  amount         day
0   Kostanay   food   4 200  2026-01-03
1  kostanay    fuel  18 500  2026-01-05
2      Rudny   food   2 600  2026-01-07
3      Rudny   food   2 600  2026-01-07
4   KOSTANAY  phone     n/a  2026-01-09
5      Rudny   rent  62 000  2026-01-11
6   Kostanay   rent  95 000         NaN
7     Astana   food   7 300  2026-01-15
types: {'city': 'str', 'kind': 'str', 'amount': 'str', 'day': 'str'}

== what is broken: look before mending
gaps per column: {'city': 0, 'kind': 0, 'amount': 0, 'day': 1}
whole duplicates: 1
cities before: ['Astana', 'KOSTANAY', 'Kostanay', 'Rudny', 'kostanay ']

== numbers: text into a number, junk into a gap
the type now: float64 | gaps: 1

== dates: the same move, errors='coerce' once more
the type now: datetime64[us] | gaps: 1

== cities: three spellings, one city
cities after: ['Astana', 'Kostanay', 'Rudny']

== duplicates: count them, then drop them
rows: 8 -> 7

== gaps: dropping or filling is a decision, not a habit
rows with no amount: 1
counting over: 6 rows | sum: 189600
and filled with zero: 189600 -- the same sum, but the average now lies
average over what is there: 31600.0 | with the zero: 27085.7

== a category instead of strings: same table, less memory
as strings: 447 | as a category: 229
```

## Going through it

### Look first, mend second

Three lines that start any work with somebody else's table:

```text
df.isna().sum()        how many gaps in each column
df.duplicated().sum()  how many twin rows
df["city"].unique()    how many distinct values there really are
```

The third is usually the eye-opener: five "cities" instead of three, because of a space here and capitals there. `value_counts()` shows the same with numbers — and then it is plain which spellings are rare, which is to say probably wrong.

### Numbers: `to_numeric`, not `astype`

`astype(float)` over a column that holds "n/a" falls over entirely: one bad cell and not a single number. `pd.to_numeric(..., errors="coerce")` converts what it can and turns what it cannot into `NaN`. That is the honest move: a value that makes no sense becomes **known to make no sense** rather than becoming a zero.

The spaces inside a number come out before the conversion: `str.replace(" ", "", regex=False)`. The same line is where a currency sign, a per cent mark and the non-breaking space that comes out of Excel — and looks exactly like an ordinary one — are removed.

Dates are mended the same way: `pd.to_datetime(..., errors="coerce")`.

### Types: `astype`, and integers with gaps

`astype` is for conversion once the data is in order: `astype("category")`, `astype("str")`, `astype("Int64")`.

The trap everybody steps in: `astype(int)` over a column with a gap **falls over**, because an ordinary integer cannot be empty:

```text
IntCastingNaNError: Cannot convert non-finite values (NA or inf) to integer
```

The way out is `Int64` with a capital letter: an integer that can be `<NA>`. The ordinary `int64` cannot.

### Duplicates: whole ones, and ones by key

`df.duplicated()` finds rows identical to one already seen. `drop_duplicates()` removes them.

But what usually breaks things is the other kind: a duplicate **by key**, where the rows differ in one column while their key is the same. Then `subset` and `keep` decide what counts as a duplicate and which row survives:

```text
df.drop_duplicates(subset=["city", "day"], keep="last")
```

`keep="last"` — when a later record counts as a correction of an earlier one. That is a decision about the meaning of the data rather than a technical detail: choosing `first` or `last`, you are answering which of the two records is the truer.

What it costs to miss was shown in [lesson twenty-nine](/read/py-kesteler-merge-join-kilt): a duplicate in a directory's key multiplies rows during a join.

### Gaps: drop or fill

Here the lesson reaches the place where technique ends and honesty begins.

`dropna()` throws away rows with gaps — whole or by `subset`. `fillna(value)` fills them. Both are decisions, not habits.

Look at the numbers in the example. The sum did not change: `fillna(0)` added a zero, and a zero adds nothing to a sum. But the average fell from 31,600 to 27,086 — because an empty receipt became a receipt for zero tenge that never existed. The same operation is harmless for one figure and a lie for another.

When filling is fair:

- the default value is known and real: no quantity given means one;
- a series over time where a value holds until the next change (`ffill`): the rate on a weekend equals Friday's rate, because it really does;
- the gap is an "unknown" that can be put into words, and you set a **marker** rather than a number: `"unknown"`, `"not stated"`. A marker is not an invention; it tells the truth.

When it is not:

- filling a numeric gap with a mean, a zero or the neighbouring value so that the formula will compute. Those are invented data, and afterwards nobody can tell them from the real ones.

A rule worth writing down: **a gap is a fact, not an obstacle**. How many there are, in which rows, and what you did about them is part of the answer rather than rubbish in front of it.

### While we are here: a category instead of strings

`astype("category")` for a column with few values and many rows: five cities are stored once and the column keeps their numbers. In the example 447 bytes become 229, and over a million rows [lesson twenty-four](/read/py-pandas-nege-kerek-olsheu) showed 155 MB against 9.

## The map of this lesson

![The map of this lesson: look, mend, write it down](/static/course/py/map-clean-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why is `to_numeric(errors="coerce")` better than `astype(float)` for a column out of somebody else's file?
2. Why did `fillna(0)` leave the sum alone and change the average?
3. When is filling a gap honest, and when is it inventing data?

## The warm-up

Three short steps before the exercise: predict, fill in, fix. The answers are at the end, but answer for yourself first.

**1. Predict.** What does this program print?

<!-- drill 1 -->
```python
import pandas as pd

amounts = pd.Series(["1200", "n/a", "4990"])
numbers = pd.to_numeric(amounts, errors="coerce")
print(numbers.tolist())
print("sum:", numbers.sum(), "| average:", numbers.mean(), "| count:", int(numbers.count()))
```

**2. Fill in the blank.** In place of `...` make "n/a" become a gap without the program falling over.

```python
# one bad value must not break the whole column
import pandas as pd

amounts = pd.Series(["1200", "n/a", "4990"])
numbers = pd.to_numeric(amounts, ...)
print(numbers.isna().sum(), "| type:", numbers.dtype)
```

**3. Fix it.** The program falls over with `IntCastingNaNError`. The year has to stay an integer and the gap has to stay a gap.

```python
# an ordinary integer cannot be empty
import pandas as pd

years = pd.Series([2024.0, 2025.0, None])
print(years.astype(int).tolist())
```

## The exercise

**Required.** Here is an export of eight rows: a city spelled three ways, amounts with spaces in them, "n/a" and "—" where numbers should be, and one twin row. Put it in order and print a **cleaning log**: how many rows came in, how many cities there were and how many there are, how many values did not become numbers, how many duplicates were dropped and how many rows reached the report. Then compute, per city, the number of receipts, the sum and the average — over the rows whose amount is known.

The expected output:

<!-- task out -->
```text
the cleaning log:
  rows in: 8
  cities in: 5
  cities out: 3 — ['Astana', 'Kostanay', 'Rudny']
  not numbers: 2
  duplicates dropped: 1
  rows into the report: 5 out of 7

per city:
          receipts     total  average
city                                 
Kostanay         3  117700.0  39233.0
Rudny            2   64600.0  32300.0

the sums agree: True
```

Done when: the output matches line for line; the numbers come from `to_numeric(errors="coerce")` rather than from replacements one at a time; the cities are brought to one spelling before the distinct ones are counted; the rows with no amount are excluded rather than filled with zero; and the log makes it visible that Astana never reached the report — its only receipt had no amount.

**On your own data.** Take any file somebody sent you and write a cleaning log for it: the three "look" lines before and the same three after. If the row count fell by more than a couple of per cent, stop and look at what exactly you threw away.

**If you feel like it.**

- Compare `keep="first"` and `keep="last"` on your own duplicates and decide which record is the truer.
- Measure `memory_usage(deep=True)` before and after `astype("category")` on a text column of yours.
- Try `ffill()` on a daily series and explain to yourself when that is honest.

## Where this fits the project

Step eleven: a cleaning step appears between the source and the report. `sholu/tazalau.py` returns two things — the table fit to count, and a log of what it did.

Its rules are the lesson's: a row without a key is dropped (that is not a gap in the data but the absence of the row itself), so is a repeated country-and-year pair, what is not a number becomes a gap rather than a zero, and the country code becomes a `category`.

The report no longer decides any of that; it only counts. One decision stayed with it, and it is about words rather than numbers: a country the directory does not know gets the marker `белгісіз` instead of an empty cell. That settles the debt of [lesson twenty-nine](/read/py-kesteler-merge-join-kilt): an empty cell in a report is a question nobody answers, and a marker answers it honestly.

Debts. The cleaning log is printed to the screen and disappears with it. Its place is beside the report — in a file that can be compared with the last run; we will get there where we get to the server.

## The answers

### To the questions

1. Because `astype(float)` falls over at the first value it cannot read and converts nothing, while `to_numeric(errors="coerce")` converts everything it can and makes the rest `NaN` — that is, known to be unreadable. That can be worked with: counted, shown, decided about.
2. Because a zero adds nothing to a sum but adds another term to an average. An empty receipt became a receipt for zero tenge that never existed, and the average fell from 31,600 to 27,086.
3. Honest when a real default is being filled in, when a series over time holds its previous value until the next change, and when a marker of "unknown" stands in place of a number. Invention when a gap is filled with a mean or a zero so that a formula will compute: such data cannot afterwards be told apart from the real kind.

### To the warm-up

1. `errors="coerce"` turns "n/a" into `NaN`, and the whole column becomes floating-point. `sum()` and `mean()` pass over the gap, and `count()` shows how many values they were computed over.

<!-- drill 1 out -->
```text
[1200.0, nan, 4990.0]
sum: 6190.0 | average: 3095.0 | count: 2
```

2. `errors="coerce"`. Without it, `to_numeric` falls over at the first value it cannot convert.

<!-- drill 2 -->
```python
import pandas as pd

amounts = pd.Series(["1200", "n/a", "4990"])
numbers = pd.to_numeric(amounts, errors="coerce")
print(numbers.isna().sum(), "| type:", numbers.dtype)
```

<!-- drill 2 out -->
```text
1 | type: float64
```

3. `astype("Int64")` — the capital-letter integer that can be empty. The ordinary `int` cannot, which is why it falls over.

<!-- drill 3 -->
```python
import pandas as pd

years = pd.Series([2024.0, 2025.0, None])
print(years.astype("Int64").tolist())
```

<!-- drill 3 out -->
```text
[2024, 2025, <NA>]
```

### To the exercise

The cities are brought to one spelling **before** the distinct ones are counted: otherwise the report gets three Kostanays, each with a part of the sum.

The rows with no amount are excluded rather than filled with zero, and that is exactly why Astana did not reach the report: its only receipt had no number. That is the right result, but it has to be visible — which is what the log is printed for. A report a city vanished from silently is no different from a report it was never in.

## Sources

- [Working with missing data](https://pandas.pydata.org/docs/user_guide/missing_data.html) — `isna`, `dropna`, `fillna` and how `NaN` behaves in arithmetic.
- [Types and casting](https://pandas.pydata.org/docs/user_guide/basics.html#dtypes) — `astype`, `to_numeric` and the integers that can be empty.
- [Categorical data](https://pandas.pydata.org/docs/user_guide/categorical.html) — when a column of words is better stored as numbers.
