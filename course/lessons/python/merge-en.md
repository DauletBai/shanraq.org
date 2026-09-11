# Joining tables: gluing two sources on a key

_Лид (summary):_ **The twenty-ninth lesson of the Python course. `merge` joins two tables on a key — and by default throws away, in silence, every row that found no match. `inner` and `left`, `indicator`, keys with different names, `join` by label, and the duplicated key that multiplies rows and the sums along with them.**

## Why this matters

Data almost never sits in one table. The bank hands over country codes; the names live in a directory. Sales in one file, prices in another. Orders apart, customers apart.

Lesson 28 ended on exactly that debt: the report came out in the codes `KAZ`, `UZB`, `RUS`, because we had no names for them. Today we put the names beside them — and see what a join can break along the way.

If you remember `JOIN` from [lesson twenty-two](/read/py-sql-group-by-join-suranys) — this is it, over tables in memory. The traps are the same ones.

## The whole thing first

The file is `sklejka.py`. Inflation by code on the left; on the right a directory, deliberately incomplete and with one code to spare.

```python
"""Lesson 29: joining two tables on a key.

On the left, inflation by country code; on the right, a directory of names. The
key is the same one, and what becomes of the rows that find no match depends on
how the join is made.
"""

import pandas as pd

# Inflation, % a year. World Bank figures, rounded to one decimal.
data = pd.DataFrame(
    [("KAZ", 2023, 14.5), ("KAZ", 2024, 8.7),
     ("UZB", 2023, 10.0), ("UZB", 2024, 9.6),
     ("RUS", 2023, 5.9), ("RUS", 2024, 8.4)],
    columns=["code", "year", "value"],
)
# The directory: not every code has a name here, and one code is spare.
names = pd.DataFrame(
    [("KAZ", "Kazakhstan", "Central Asia"),
     ("UZB", "Uzbekistan", "Central Asia"),
     ("KGZ", "Kyrgyzstan", "Central Asia")],
    columns=["code", "name", "region"],
)
print("rows on the left:", len(data), "| on the right:", len(names))

print()
print("== inner: only what was found on both sides stays")
inner = data.merge(names, on="code")
print(inner)
print("rows now:", len(inner), "-- lost:", len(data) - len(inner))

print()
print("== left: everything on the left stays, NaN appears on the right")
left = data.merge(names, on="code", how="left")
print(left.tail(3))
print("rows without a name:", int(left["name"].isna().sum()))

print()
print("== indicator shows where each row came from")
both = data.merge(names, on="code", how="outer", indicator=True)
print(both["_merge"].value_counts().to_dict())

print()
print("== keys with different names")
other = names.rename(columns={"code": "iso"})
print(data.merge(other, left_on="code", right_on="iso").columns.tolist())

print()
print("== by label: join")
a = data.set_index("code")
b = names.set_index("code")
print(a.join(b, how="left").head(2))

print()
print("== a duplicate in the key multiplies the rows")
twice = pd.concat([names, names.iloc[[0]]], ignore_index=True)
print("directory rows:", len(twice), "| after the join:", len(data.merge(twice, on="code")))
try:
    data.merge(twice, on="code", validate="many_to_one")
except pd.errors.MergeError as err:
    print("validate caught it:", str(err).split("\n")[0])
```

It prints:

```text
rows on the left: 6 | on the right: 3

== inner: only what was found on both sides stays
  code  year  value        name        region
0  KAZ  2023   14.5  Kazakhstan  Central Asia
1  KAZ  2024    8.7  Kazakhstan  Central Asia
2  UZB  2023   10.0  Uzbekistan  Central Asia
3  UZB  2024    9.6  Uzbekistan  Central Asia
rows now: 4 -- lost: 2

== left: everything on the left stays, NaN appears on the right
  code  year  value        name        region
3  UZB  2024    9.6  Uzbekistan  Central Asia
4  RUS  2023    5.9         NaN           NaN
5  RUS  2024    8.4         NaN           NaN
rows without a name: 2

== indicator shows where each row came from
{'both': 4, 'left_only': 2, 'right_only': 1}

== keys with different names
['code', 'year', 'value', 'iso', 'name', 'region']

== by label: join
      year  value        name        region
code                                       
KAZ   2023   14.5  Kazakhstan  Central Asia
KAZ   2024    8.7  Kazakhstan  Central Asia

== a duplicate in the key multiplies the rows
directory rows: 4 | after the join: 6
validate caught it: Merge keys are not unique in right dataset; not a many-to-one merge
```

## Going through it

### `inner` by default, and it is the most expensive line in the lesson

`data.merge(names, on="code")` keeps only the rows whose key was found **on both sides**. Six rows became four in the example: Russia is not in the directory, and both of its rows vanished. No error, no warning.

This is the commonest silent loss of data in table work. The report comes out plausible, the sums are smaller than the truth, and there is exactly one way to notice: count the rows before and after.

The habit: **compare the length after every join**. `len(data)` and `len(joined)` differ — you either lost rows or multiplied them; there is no third possibility.

### `left`: the data decides, the directory adds

`how="left"` keeps every row of the left table; where no match was found, the right-hand columns hold `NaN`. That is the sane choice when the data is on the left and a directory on the right: the data decides which rows exist, the directory only adds words to them.

There is also `how="right"` (the mirror image) and `how="outer"` — keep everything from both sides. `outer` earns its place less in a report than in a check: with `indicator=True` it shows how many rows were found on both sides, how many were left unmatched on the left, and how many on the right.

```text
{'both': 4, 'left_only': 2, 'right_only': 1}
```

Three numbers worth looking at **before** joining in earnest. `left_only` is what the directory does not have; `right_only` is what the data does not have.

### When the key is called something else

`left_on="code", right_on="iso"` joins on columns with different names. Both columns stay in the answer, and the spare one is yours to drop.

When the key is the row labels rather than a column, `join` does the work: `a.join(b)` joins on the index. It is the same as `merge(left_index=True, right_index=True)`, only shorter; [lesson twenty-five](/read/py-dataframe-series-indeks-qoltanba) showed that arithmetic between series is matched by label too — `join` is that, for tables.

Columns with the same name on both sides do not clash: they get `_x` and `_y` attached. Better to set those yourself — `suffixes=("_data", "_directory")` — or in a month nobody will remember which is which.

### A duplicate in the key multiplies the rows

The second trap, the mirror image of the first. If a code appears twice in the directory, every row of data with that code is **doubled**. There are more rows than there were, the sums grow, and again nothing is said.

The drill shows it in numbers: a sum of 300 turns into 400 — not because anything was added, but because one row was counted twice.

The cure is `validate`:

- `validate="one_to_one"` — the key is unique on both sides;
- `validate="many_to_one"` — it may repeat on the left but not on the right (the commonest case: data and a directory);
- `validate="one_to_many"` — the other way round.

If it does not hold, a `MergeError` instead of a silent multiplication. This is one of those checks worth writing every time: it costs nothing and catches the mistake where it happened.

### `merge` and `concat` are different things

`merge` joins **on a key** and adds columns. `concat` stacks tables **one under another** (or side by side) and asks about no key at all: that is how twelve monthly files become one series. If you are looking for "combine tables", decide first which you want: more columns, or more rows.

## The map of this lesson

![The map of this lesson: two sources and one key](/static/course/py/map-merge-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. What happens to the rows that found no match, under `inner` and under `left`?
2. How do you find out, in one line, how many rows found no match on each side?
3. Why can a join leave you with more rows than you started with?

## The warm-up

Three short steps before the exercise: predict, fill in, fix. The answers are at the end, but answer for yourself first.

**1. Predict.** What does this program print?

<!-- drill 1 -->
```python
import pandas as pd

left = pd.DataFrame({"code": ["KAZ", "UZB", "RUS"], "value": [8.7, 9.6, 8.4]})
right = pd.DataFrame({"code": ["KAZ", "UZB"], "name": ["Kazakhstan", "Uzbekistan"]})
joined = left.merge(right, on="code")
print(len(left), "→", len(joined))
print(joined.to_string(index=False))
```

**2. Fill in the blank.** In place of `...` make it so that three rows remain and Russia has a gap where its name would be.

```python
# the data is on the left and must not be lost
import pandas as pd

left = pd.DataFrame({"code": ["KAZ", "UZB", "RUS"], "value": [8.7, 9.6, 8.4]})
right = pd.DataFrame({"code": ["KAZ", "UZB"], "name": ["Kazakhstan", "Uzbekistan"]})
joined = left.merge(right, on="code", ...)
print(len(joined), "| without a name:", int(joined["name"].isna().sum()))
```

**3. Fix it.** There were two rows and now there are three, and the sum grew from three hundred to four hundred. Nothing was added.

```python
# the code appears twice in the directory
import pandas as pd

data = pd.DataFrame({"code": ["KAZ", "UZB"], "value": [100, 200]})
names = pd.DataFrame({"code": ["KAZ", "KAZ", "UZB"], "name": ["Kazakhstan", "Kazakhstan", "Uzbekistan"]})
joined = data.merge(names, on="code")
print("rows:", len(joined), "| sum:", int(joined["value"].sum()))
```

## The exercise

**Required.** You are given eight rows of inflation across four codes and a directory of four countries — not quite the same four. Join them so that no row of data is lost, and print:

1. how many rows there were and how many there are, along with the codes that found no name;
2. the average inflation per country, with the name beside the code;
3. the average per region, counting only the rows whose region is known.

The expected output:

<!-- task out -->
```text
rows before: 8 | after: 8
without a name: ['RUS']

average inflation per country:
code       name  value
 KAZ Kazakhstan  11.60
 KGZ Kyrgyzstan   8.55
 RUS        NaN   7.15
 UZB Uzbekistan   9.80

per region:
              countries  average
region                          
Central Asia          3     9.98
```

Done when: the output matches line for line; the number of rows is the same before and after; the join is checked with `validate`; the rows without a region are left out of the per-region grouping rather than filed under somebody else's region.

**On your own data.** Take two of your tables with a key in common — anything with a code or an identifier. Look at `merge(..., how="outer", indicator=True)` first and read off the three numbers. Then join the way your task needs and compare the length before and after.

**If you feel like it.**

- Run a join with `validate="one_to_one"` over data whose key repeats and read the whole message.
- Join the same tables with `join` on the index and compare what came out.
- Stack two tables with `concat` and explain to yourself how that differs from `merge`.

## Where this fits the project

Step ten: the report stops speaking in codes. `sholu/anyqtama.py` is a table of three columns — code, name, region — and the report joins it on the code.

Two decisions there are worth reading in the code. The join is a **left** one: the data decides which rows exist, the directory only adds words, and a country the directory has never heard of stays in the report under its code. And the join is **checked**: `validate="one_to_one"`, because a repeated code in the directory would multiply the rows and the totals would grow on their own.

Debts. Our directory is incomplete on purpose: Russia is not in it, and in the report it stayed a code with an empty region. That is the right behaviour, but an empty cell in a report is a question somebody has to answer; what to fill such gaps with is the subject of the next lesson, on dirty data.

## The answers

### To the questions

1. Under `inner` they disappear — from both sides, silently. Under `left` every row of the left table stays and the right-hand columns are filled with `NaN` where no match was found.
2. `merge(..., how="outer", indicator=True)` adds a `_merge` column holding `both`, `left_only` and `right_only`; `value_counts()` over it gives the three numbers.
3. Because the key is not unique in the right-hand table: every row on the left is multiplied by the number of matches on the right. The `validate` argument catches it and raises a `MergeError` instead.

### To the warm-up

1. The default join is `inner`, and the row with the code `RUS` disappears because the directory does not have it. Two rows remain out of three.

<!-- drill 1 out -->
```text
3 → 2
code  value       name
 KAZ    8.7 Kazakhstan
 UZB    9.6 Uzbekistan
```

2. `how="left"`. The left table is kept whole, and Russia's `name` becomes a gap.

<!-- drill 2 -->
```python
import pandas as pd

left = pd.DataFrame({"code": ["KAZ", "UZB", "RUS"], "value": [8.7, 9.6, 8.4]})
right = pd.DataFrame({"code": ["KAZ", "UZB"], "name": ["Kazakhstan", "Uzbekistan"]})
joined = left.merge(right, on="code", how="left")
print(len(joined), "| without a name:", int(joined["name"].isna().sum()))
```

<!-- drill 2 out -->
```text
3 | without a name: 1
```

3. Drop the duplicate from the directory — `names.drop_duplicates("code")` — and add the check, so that next time it is the program that notices rather than the report.

<!-- drill 3 -->
```python
import pandas as pd

data = pd.DataFrame({"code": ["KAZ", "UZB"], "value": [100, 200]})
names = pd.DataFrame({"code": ["KAZ", "KAZ", "UZB"], "name": ["Kazakhstan", "Kazakhstan", "Uzbekistan"]})
joined = data.merge(names.drop_duplicates("code"), on="code", validate="one_to_one")
print("rows:", len(joined), "| sum:", int(joined["value"].sum()))
```

<!-- drill 3 out -->
```text
rows: 2 | sum: 300
```

### To the exercise

The join is a left one and it is checked: `how="left", validate="many_to_one"`. Left, because the data matters more than the directory; checked, because the directory came from outside, which means nobody is answerable for the uniqueness of its key.

The rows without a region are dropped from the per-region grouping with `dropna(subset=["region"])`. Not because they are unwanted, but because their region is not empty but unknown: filing them under "Central Asia" would be inventing data.

## Sources

- [Merge, join and concatenate](https://pandas.pydata.org/docs/user_guide/merging.html) — `merge`, `join`, `concat` and all their arguments.
- [The validate argument](https://pandas.pydata.org/docs/reference/api/pandas.DataFrame.merge.html) — the check that catches multiplied rows.
- [Consumer price inflation, World Bank](https://data.worldbank.org/indicator/FP.CPI.TOTL.ZG) — the source of this lesson's numbers.
