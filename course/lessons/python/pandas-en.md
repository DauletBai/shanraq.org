# Why pandas, when lists already work: what the measurement showed

_Лид (summary):_ **The twenty-fourth lesson of the Python course. The pandas module opens with the question it ought to open with: is the library worth it. A measurement answers, not belief. Where pandas is forty times faster, where it is merely twice, where it loses to an ordinary loop, what the memory costs — and the rule for when a table is needed.**

## Why this matters

The digest can already take data off the network, put it in a database and answer with a query. What comes next are questions that never come one at a time: how much per category, how that moved month by month, what happens converted into another currency, and now the same thing without weekends. Over lists, each of those is a fresh loop and a fresh dictionary, and ten questions later the program is made of them.

pandas is usually described as fast. That is half the truth, and the other half is the more useful one. So the module begins with a stopwatch rather than a promise: one question, asked two ways.

## The whole thing first

Installation comes first, once for the whole module. The environment was set up in [lesson two](/read/py-jumys-orny-venv), and the library goes into it:

```text
python -m venv .venv
source .venv/bin/activate      # Windows: .venv\Scripts\activate
pip install pandas==3.0.5
```

The version is pinned deliberately. pandas prints a table its own way, and under another version the output here may differ from yours — through no mistake of yours.

The file is `zamer.py`. It asks one question — the average receipt per category — first by hand, then through pandas, and checks that the answers agree.

```python
"""Lesson 24: why pandas, when lists already work.

One question -- the average receipt per category -- asked twice: by hand over
lists, and through pandas. The answers have to agree first; only then is there
anything worth arguing about.
"""

import pandas as pd

KINDS = ["food", "fuel", "rent", "phone", "other"]
BASE = {"food": 2500, "fuel": 9000, "rent": 60000, "phone": 4000, "other": 1500}
N = 300_000

# Built from a counter, with no random numbers: the table is the same for all.
rows = [(KINDS[i % 5], BASE[KINDS[i % 5]] + (i * 37) % 900 - 450) for i in range(N)]
print("rows:", len(rows))

print()
print("== by hand: two dictionaries and one pass")
total, count = {}, {}
for kind, amount in rows:
    total[kind] = total.get(kind, 0) + amount
    count[kind] = count.get(kind, 0) + 1
by_hand = {kind: total[kind] / count[kind] for kind in sorted(total)}
for kind, average in by_hand.items():
    print(f"{kind:6} {average:9.2f}")

print()
print("== pandas: one line")
df = pd.DataFrame(rows, columns=["kind", "amount"])
by_pandas = df.groupby("kind")["amount"].mean()
print(by_pandas)

print()
print("== the answers agree:", all(round(by_hand[k], 6) == round(by_pandas[k], 6) for k in by_hand))
print("table:", type(df).__name__, "| columns:", list(df.columns))
print("column types:", dict(df.dtypes.astype(str)))
print("answer:", type(by_pandas).__name__, "-- not a list, a column with labels")

print()
print("== a whole column, no loop")
with_vat = df["amount"] * 1.12
print("pandas:", [round(x, 2) for x in with_vat.head(3)])
print("by hand:", [round(amount * 1.12, 2) for _, amount in rows[:3]])
```

It prints:

```text
rows: 300000

== by hand: two dictionaries and one pass
food     2497.49
fuel     8999.49
other    1500.48
phone    3998.49
rent    60001.50

== pandas: one line
kind
food      2497.485
fuel      8999.490
other     1500.480
phone     3998.490
rent     60001.495
Name: amount, dtype: float64

== the answers agree: True
table: DataFrame | columns: ['kind', 'amount']
column types: {'kind': 'str', 'amount': 'int64'}
answer: Series -- not a list, a column with labels

== a whole column, no loop
pandas: [2296.0, 9617.44, 66778.88]
by hand: [2296.0, 9617.44, 66778.88]
```

## Going through it

### The same answer first, the stopwatch second

Comparing the speed of two programs that compute different things is the commonest mistake in measurement. So the first thing `zamer.py` does is prove that both ways produced the same numbers. Anyone can compute the wrong thing quickly.

The comparison goes through `round(..., 6)` rather than a plain `==`: these are floating-point numbers, and [lesson four](/read/py-aynymalylar-men-tipter) showed why comparing them head-on does not work.

### What a DataFrame actually is

A list of tuples stores **rows**: three hundred thousand small objects, each with its own header and its own references. A `DataFrame` stores **columns**: three hundred thousand amounts lie in one block of same-typed numbers, one after another, with no wrappers. Everything else follows from that.

Which is why `df["amount"] * 1.12` does not walk three hundred thousand Python objects one at a time. It hands the whole block to code that is not written in Python, and the multiplication happens there. The loop did not disappear — it merely stopped being a Python loop.

The `Series` follows from the same fact. The answer to a grouping is not a list of numbers but a column with labels: categories on the left, values on the right, and the two stay together. A list out of `groupby` would have to be labelled by hand — with exactly the dictionary we wrote ourselves.

### The measurement

The same table, but a million rows, four questions and `timeit`. The program prints nothing about correctness — correctness was settled above:

```python
# the measurement: one question, asked two ways. The numbers will be yours.
import timeit

import pandas as pd

KINDS = ["food", "fuel", "rent", "phone", "other"]
BASE = {"food": 2500, "fuel": 9000, "rent": 60000, "phone": 4000, "other": 1500}
rows = [(KINDS[i % 5], BASE[KINDS[i % 5]] + (i * 37) % 900 - 450) for i in range(1_000_000)]
df = pd.DataFrame(rows, columns=["kind", "amount"])


def hand_group():
    total, count = {}, {}
    for kind, amount in rows:
        total[kind] = total.get(kind, 0) + amount
        count[kind] = count.get(kind, 0) + 1
    return {k: total[k] / count[k] for k in total}


def best(work):
    # Five runs, keep the quickest: the slow ones are the machine looking away.
    return min(timeit.repeat(work, number=1, repeat=5))


print("average receipt: by hand %.3f s, pandas %.3f s"
      % (best(hand_group), best(lambda: df.groupby("kind")["amount"].mean())))
print("a new column:    by hand %.3f s, pandas %.3f s"
      % (best(lambda: [a * 1.12 for _, a in rows]), best(lambda: df["amount"] * 1.12)))
print("sum of one category: by hand %.3f s, pandas %.3f s"
      % (best(lambda: sum(a for k, a in rows if k == "food")),
         best(lambda: df.loc[df["kind"] == "food", "amount"].sum())))
print("sorting:         by hand %.3f s, pandas %.3f s"
      % (best(lambda: sorted(rows, key=lambda r: r[1])), best(lambda: df.sort_values("amount"))))
```

On my machine — Python 3.14.5, pandas 3.0.5, a million rows:

| What is computed | By hand | pandas |
|---|---|---|
| average receipt per category | 0.124 s | 0.062 s |
| a column with 12 % added | 0.045 s | 0.001 s |
| the sum of one category | 0.021 s | 0.045 s |
| sorting by amount | 0.173 s | 0.089 s |
| memory for the same table | 100 MB | 155 MB |

Your numbers will be different: another machine, another version of pandas, other column types, another load on the box. The distances between the columns move too — more in one place, less in another. What holds is something else: **which operation puts the difference into orders of magnitude and which leaves it in multiples**. That is the part to remember; the numbers are to be measured again on your own machine.

### What those numbers show

**Forty times, where a column does the work.** `df["amount"] * 1.12` against a list comprehension: 0.001 against 0.045. A million multiplications done once as a block instead of a million times one by one. That is the main argument for pandas, and the only place where the difference is one of order rather than degree.

**Twice, where everything has to be walked anyway.** Grouping and sorting win by two, not by forty. The work is the same; it is simply carried out by faster code than ours. Twice is good, and twice is not the difference libraries are installed for.

**Slower, where we asked for more than we needed.** `df.loc[df["kind"] == "food", "amount"].sum()` first builds a mask: a million yes-or-no answers, a whole new column, and only then adds. Meanwhile `sum(a for k, a in rows if k == "food")` walks the data once and builds nothing. On one question the list wins. The win appears once the mask is **kept and reused**: `food = df["kind"] == "food"`, and then `df.loc[food, "amount"].sum()`, `df.loc[food, "amount"].mean()`, `df.loc[food]` — a million comparisons made once for ten questions. Write the condition out afresh every time and pandas will build the mask afresh too.

**Memory is not free.** The same table: a hundred megabytes as a list, a hundred and fifty-five as a table. The numbers are not to blame, the words are: the `kind` column holds five distinct values, but they lie there one per row. `df["kind"].astype("category")` turns those 155 MB into 9 MB — the five words are stored once and the column keeps their numbers. It is the same move as in [the lesson on sets](/read/py-jiyndar-set-qiylysu-ayyrma): do not store the same thing twice.

**The import costs a third of a second.** For a program that runs for a minute that is nothing. For one that a scheduler starts every minute it is the main expense — and a reason not to drag pandas into a small script.

### The real reason people reach for it

Eight lines by hand against one is what actually settles it. Not the seconds, but the fact that the second, third and tenth question to the same table cost one line each, and that the answer stays a table you can ask the next question of.

You have seen the same argument in [the lesson on SQL](/read/py-sql-group-by-join-suranys): `GROUP BY` is better than a loop not because it is faster but because the question is written as a question. pandas is the same move, with the table sitting in your program's memory rather than in a database.

### When not to reach for it

- Five hundred rows and one question: a loop is simpler, and the import takes longer than the work.
- The data is not a table: nested JSON or a tree has to be talked into being flat first.
- You are computing this once in your life: a dependency you install and pin will outlive the task.

The rule that holds: **a table with more than one question to it — take pandas; one answer over little data — do not**.

## The map of this lesson

![The map of this lesson: lists against a table](/static/course/py/map-pandas-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why does the program prove the answers agree before it measures any time?
2. Why is multiplying a column forty times faster than a loop, while grouping is only twice?
3. In which case did pandas turn out slower than a plain list, and why?

## The warm-up

Three short steps before the exercise: predict, fill in, fix. The answers are at the end, but answer for yourself first.

**1. Predict.** What does this program print?

<!-- drill 1 -->
```python
import pandas as pd

cheques = pd.Series([1200, 800, 4000], index=["food", "food", "phone"])
print(cheques * 1.12)
print("sum:", cheques.sum(), "| average:", cheques.mean())
```

**2. Fill in the blank.** In place of `...`, ask for as many runs as it takes for `len(times)` to print `5`.

```python
# how many times to measure: one run is not a measurement
import timeit

times = timeit.repeat(lambda: sum(range(100_000)), number=1, repeat=...)
print("measurements:", len(times), "| the best is no worse than the worst:", min(times) <= max(times))
```

**3. Fix it.** Two ways count the rows of the same table and disagree by one. Find the extra row and remove it.

```python
# by hand it comes out one row longer than in pandas
import csv
import io

import pandas as pd

TEXT = "kind,amount\nfood,1200\nfood,800\nphone,4000\n"
reader = csv.reader(io.StringIO(TEXT))
print("rows by hand:", len(list(reader)))
print("rows in pandas:", len(pd.read_csv(io.StringIO(TEXT))))
```

## The exercise

**Required.** Take the same 300 000-row table from the lesson. For each category compute three numbers: the sum, the average receipt and the share of the total. By hand first — dictionaries in one pass — then through pandas: `df.groupby("kind")["amount"].agg(["sum", "mean"])` and a share column.

Print the table sorted by sum from the top down, a "total" line and a line that reports the check. The check has to agree to the sixth decimal.

The expected output:

<!-- task out -->
```text
category           sum    average   share
rent        3600089700   60001.50   77.9%
fuel         539969400    8999.49   11.7%
phone        239909400    3998.49    5.2%
food         149849100    2497.49    3.2%
other         90028800    1500.48    1.9%
total       4619846400
they agree: yes
```

Done when: the output matches line for line; the shares are computed from the total sum rather than from the sum of the averages; the check compares rounded numbers instead of floating-point ones head-on; the table is sorted by sum rather than alphabetically.

**On your own data.** Take a file of yours — a bank export, a table of expenses, anything a few thousand rows long — and ask it one question two ways. Measure with `timeit`. Write the two numbers down and answer for yourself: on your data, was the library worth that difference.

**If you feel like it.**

- Measure what `import pandas` alone costs on your machine: `python -X importtime -c "import pandas"`.
- Make the table ten times bigger and see which of the four differences grew and which stayed where it was.
- Convert the `kind` column with `astype("category")` and compare `df.memory_usage(deep=True).sum()` before and after.

## Where this fits the project

Nowhere yet, and that is on purpose. The lesson answers "is it worth it", not "how"; the digest changes from the next lesson on, where its table becomes a `DataFrame`. The only thing added now is the line `pandas==3.0.5` in the project's `requirements.txt`: the version is pinned because the output in these lessons was printed with that one.

Still open. A library is a dependency: it has to be installed, updated and one day repaired after an update. While the digest counts exchange rates, where rows number in the thousands, all pandas wins there is shorter code rather than speed. Saying so is the honest thing, rather than pretending we got faster.

## The answers

### To the questions

1. Because the speed of two programs is only comparable when they give the same answer. Otherwise the measurement measures something else: a program that computes wrongly is usually faster too — it is doing less work.
2. Multiplying a whole column leaves Python and runs along a block of memory in order, so it wins by an order. Grouping still has to look at every row and drop it into a bucket: the work is the same and merely carried out faster — hence "twice" rather than "forty times".
3. On the sum of one category. `df["kind"] == "food"` builds a whole new column of a million yes-or-no values before adding anything, while the generator walks the data once and creates nothing. As soon as the table gets many questions, that payment is repaid.

### To the warm-up

1. The multiplication reaches every value, the labels stay where they are, and the column's type becomes floating-point. `sum()` and `mean()` count over the whole column at once.

<!-- drill 1 out -->
```text
food     1344.0
food      896.0
phone    4480.0
dtype: float64
sum: 6000 | average: 2000.0
```

2. `repeat=5`. `timeit.repeat` returns a list with as many runs as you asked for; you take the smallest of them rather than the average, because the extra time is always somebody else's work on the same machine and not your program.

<!-- drill 2 -->
```python
import timeit

times = timeit.repeat(lambda: sum(range(100_000)), number=1, repeat=5)
print("measurements:", len(times), "| the best is no worse than the worst:", min(times) <= max(times))
```

<!-- drill 2 out -->
```text
measurements: 5 | the best is no worse than the worst: True
```

3. The extra row is the header. `csv.reader` hands it over as an ordinary row of data, while `read_csv` understands that it holds the column names. Skip it with `next(reader)`.

<!-- drill 3 -->
```python
import csv
import io

import pandas as pd

TEXT = "kind,amount\nfood,1200\nfood,800\nphone,4000\n"
reader = csv.reader(io.StringIO(TEXT))
next(reader)
print("rows by hand:", len(list(reader)))
print("rows in pandas:", len(pd.read_csv(io.StringIO(TEXT))))
```

<!-- drill 3 out -->
```text
rows by hand: 3
rows in pandas: 3
```

That is also the answer to why `read_csv` exists when the `csv` module from [lesson twelve](/read/py-csv-utir-tyrnaqsha-tanba) is right there: it knows about the header, about column types and about missing values — that is, about everything a loop leaves you to remember yourself.

### To the exercise

The share has to be computed from the sum of the sums, not from the sum of the averages: the categories are of different sizes, and an average of averages answers a different question. In pandas that is `table["sum"] / table["sum"].sum() * 100`; by hand it is a division by `sum(total.values())`, computed once before the printing loop.

Sorting by sum is `table.sort_values("sum", ascending=False)`. Without it the table comes out alphabetically, because `groupby` sorts by the key.

## Sources

- [10 minutes to pandas](https://pandas.pydata.org/docs/user_guide/10min.html) — the official introduction; the first two sections are the ones to read.
- [What is new in pandas 3.0](https://pandas.pydata.org/docs/whatsnew/v3.0.0.html) — the version this lesson's output was printed with.
- [The timeit module](https://docs.python.org/3/library/timeit.html) — how to measure time so that you measure the program and not the machine.
- [pandas: scaling to large datasets](https://pandas.pydata.org/docs/user_guide/scale.html) — what to do when the table no longer fits in memory.
