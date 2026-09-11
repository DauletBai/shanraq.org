# The module's vocabulary: eight words a table speaks in

_Лид (summary):_ **Not a lesson but a reference page, standing in front of the heavy part of the pandas module. Eight words — a table, a column, a label and a position, a mask, a missing value, a key, a type and a shape — on one table of three rows. Read it once now and come back when a word turns up in a lesson.**

## How to use this page

Seven lessons follow one after another, and they are denser than anything so far. Each brings in ideas of its own, but eight words run through all of them, and confusing those costs more than not knowing any single function.

Here they are together, on one small table. This is not a lesson: there is no warm-up, no exercise and no walk-through, only the vocabulary. Skim it now without trying to memorise anything, and come back when a word turns up in earnest. Each of the eight is explained properly in a lesson of its own, and which one is said beside it.

pandas was installed in [the previous lesson](/read/py-pandas-nege-kerek-olsheu), so the program on this page runs straight away.

## Eight words on one table

```python
"""The module's vocabulary: eight words on one small table.

Three receipts, two columns and one missing value. That is enough to show every
word the pandas lessons speak in.
"""

import pandas as pd

df = pd.DataFrame(
    {"city": ["Kostanay", "Rudny", "Astana"], "amount": [4200.0, 16200.0, None]},
    index=[101, 102, 103],
)
df.index.name = "receipt"
print(df)

print()
print("1. the table:", type(df).__name__, "| 2. a column:", type(df["city"]).__name__)
print("3. row labels:", df.index.tolist(), "| positions: 0, 1, 2")
print("   loc[102] — by label:   ", df.loc[102, "city"])
print("   iloc[1]  — by position:", df.iloc[1]["city"])
print("4. a mask:", (df["amount"] > 5000).tolist())
print("5. missing values:", int(df["amount"].isna().sum()), "| the sum passes over them:", df["amount"].sum())
print("6. a key is a column without repeats; city repeats:", int(df["city"].duplicated().sum()))
print("7. column types:", dict(df.dtypes.astype(str)))
print("8. the shape:", df.shape, "— rows and columns")
```

It prints:

```text
             city   amount
receipt                   
101      Kostanay   4200.0
102         Rudny  16200.0
103        Astana      NaN

1. the table: DataFrame | 2. a column: Series
3. row labels: [101, 102, 103] | positions: 0, 1, 2
   loc[102] — by label:    Rudny
   iloc[1]  — by position: Rudny
4. a mask: [False, True, False]
5. missing values: 1 | the sum passes over them: 20400.0
6. a key is a column without repeats; city repeats: 0
7. column types: {'city': 'str', 'amount': 'float64'}
8. the shape: (3, 2) — rows and columns
```

![The map of this page: eight words](/static/course/py/map-glossary-en.svg)

## What each of them means

**1. A table — `DataFrame`.** Not a list of lists: the columns have names, the rows have labels, and every column has a type of its own. In full in [lesson twenty-five](/read/py-dataframe-series-indeks-qoltanba).

**2. A column — `Series`.** Values plus those same labels. `df["city"]` is a `Series`, while `df[["city", "amount"]]` is a table again: there is a list inside the brackets. Half of all confusion in somebody else's code is a `Series` where a table was expected.

**3. A label and a position.** A label is a row's name (here it is the receipt number: 101, 102, 103). A position is its place in the order (0, 1, 2). `loc` asks by label, `iloc` by position. While the labels run from zero in order the difference is invisible; it shows up after a sort or a filter — [lesson twenty-seven](/read/py-tandau-suzgi-maska-loc).

**4. A mask.** A comparison against a column gives not a yes or a no but a **column** of yes and no, one answer per row. That is what a table is filtered with: `df[df["amount"] > 5000]`. Masks are joined by `&`, `|`, `~`, each part in brackets of its own.

**5. A missing value — `NaN`.** Not a zero and not an empty string: "there is no number here". It is contagious in arithmetic (`NaN + 5` gives `NaN`), and `sum` and `mean` pass over it — which is why the example prints a sum of 20,400 rather than an error. What to do about it is [lesson thirty](/read/py-las-derek-isna-fillna-astype).

**6. A key.** The column two tables are joined on: a country code, a receipt number, a date. One that does not repeat will do — otherwise the join multiplies rows. The check is a single line: `df["code"].duplicated().sum()`. Joining is [lesson twenty-nine](/read/py-kesteler-merge-join-kilt).

**7. A column's type — `dtype`.** `int64`, `float64`, `str`, `datetime64`, `category`. The type decides what can be done with the column: a string has no month, and comparing a string with a number raises. It is the first thing to look at in somebody else's file.

**8. The shape — `shape`.** The pair "rows, columns". It is looked at before and after every step: if the number of rows changed otherwise than you expected after a join or a cleaning, counting anything further is pointless.

## Three questions for any table somebody sends you

Before computing anything at all:

```text
df.shape          how many rows and columns
df.dtypes         what pandas takes each column to be
df.head(3)        what it looks like
```

The fourth is `df.index`: what the rows are labelled by. Half of the oddities in somebody else's code are explained by an index that was forgotten about.

## What comes next

The next lesson is `DataFrame` and `Series` in full: the index, slices, matching by label. After that: reading and writing, selecting, grouping, joining, dirty data and time in a table. If a word stops making sense in some lesson, this page is not going anywhere.
