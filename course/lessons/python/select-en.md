# Selecting and filtering: keeping only what is needed

_Лид (summary):_ **The twenty-seventh lesson of the Python course. A mask is a column of yes and no, and a whole table is filtered with one. `df[mask]`, `.loc[mask, columns]`, `&`, `|`, `~` and the brackets that are not decoration, `isin`, `between`, `str` — and the labels a filter leaves exactly as they were, which is where `loc` and `iloc` finally part company.**

## Why this matters

A table is rarely wanted whole. What is wanted is the receipts of one city, the spending of the second half of the month, the rows where the amount is suspiciously large. Over lists that is a loop with an `if` and a new list; over a table it is one expression, and it reads as a condition rather than as a walk.

This is also where the promise of [lesson twenty-five](/read/py-dataframe-series-indeks-qoltanba) comes due: `loc` and `iloc` part company exactly where the rows have been filtered.

## The whole thing first

The file is `vyborka.py`. Ten receipts from January and eight ways to take what is needed out of them.

```python
"""Lesson 27: taking the rows and the columns you need.

Ten receipts from January. We pick from them -- the columns first, then the
rows, then both at once.
"""

import pandas as pd

rows = [
    ("2026-01-03", "Kostanay", "food", 4200),
    ("2026-01-05", "Kostanay", "fuel", 18500),
    ("2026-01-07", "Rudny", "food", 2600),
    ("2026-01-09", "Kostanay", "phone", 4990),
    ("2026-01-12", "Astana", "food", 7300),
    ("2026-01-15", "Kostanay", "rent", 95000),
    ("2026-01-18", "Rudny", "fuel", 16200),
    ("2026-01-21", "Astana", "phone", 4990),
    ("2026-01-24", "Kostanay", "food", 3100),
    ("2026-01-27", "Rudny", "rent", 62000),
]
df = pd.DataFrame(rows, columns=["day", "city", "kind", "amount"])
df["day"] = pd.to_datetime(df["day"])
print(df)

print()
print("== columns: one and several")
print("df['amount'] →", type(df["amount"]).__name__)
print("df[['city', 'amount']] →", type(df[["city", "amount"]]).__name__)

print()
print("== a mask is a column of yes and no")
big = df["amount"] > 10000
print(big.tolist())
print("rows that matched:", int(big.sum()), "out of", len(df))

print()
print("== filtering the rows")
print(df[big])

print()
print("== two conditions: & | ~ and the brackets you must write")
kostanay_food = (df["city"] == "Kostanay") & (df["kind"] == "food")
print("Kostanay and food:", int(kostanay_food.sum()))
print("not food:", int((~(df["kind"] == "food")).sum()))
print("food or phone:", int(df["kind"].isin(["food", "phone"]).sum()))

print()
print("== rows and columns in one go")
print(df.loc[kostanay_food, ["day", "amount"]])

print()
print("== three more conditions off the shelf")
print("between 4000 and 20000:", int(df["amount"].between(4000, 20000).sum()))
print("the city starts with K:", int(df["city"].str.startswith("K").sum()))
print("the second half of the month:", int((df["day"] >= "2026-01-15").sum()))

print()
print("== a filter does not renumber the labels")
rudny = df[df["city"] == "Rudny"]
print("labels:", rudny.index.tolist())
print("loc[2]:", rudny.loc[2, "kind"], "| iloc[2]:", rudny.iloc[2]["kind"])

print()
print("== writing through loc")
df.loc[df["kind"] == "phone", "kind"] = "mobile phone"
print(df["kind"].value_counts().to_dict())
```

It prints:

```text
         day      city   kind  amount
0 2026-01-03  Kostanay   food    4200
1 2026-01-05  Kostanay   fuel   18500
2 2026-01-07     Rudny   food    2600
3 2026-01-09  Kostanay  phone    4990
4 2026-01-12    Astana   food    7300
5 2026-01-15  Kostanay   rent   95000
6 2026-01-18     Rudny   fuel   16200
7 2026-01-21    Astana  phone    4990
8 2026-01-24  Kostanay   food    3100
9 2026-01-27     Rudny   rent   62000

== columns: one and several
df['amount'] → Series
df[['city', 'amount']] → DataFrame

== a mask is a column of yes and no
[False, True, False, False, False, True, True, False, False, True]
rows that matched: 4 out of 10

== filtering the rows
         day      city  kind  amount
1 2026-01-05  Kostanay  fuel   18500
5 2026-01-15  Kostanay  rent   95000
6 2026-01-18     Rudny  fuel   16200
9 2026-01-27     Rudny  rent   62000

== two conditions: & | ~ and the brackets you must write
Kostanay and food: 2
not food: 6
food or phone: 6

== rows and columns in one go
         day  amount
0 2026-01-03    4200
8 2026-01-24    3100

== three more conditions off the shelf
between 4000 and 20000: 6
the city starts with K: 5
the second half of the month: 5

== a filter does not renumber the labels
labels: [2, 6, 9]
loc[2]: food | iloc[2]: rent

== writing through loc
{'food': 4, 'fuel': 2, 'mobile phone': 2, 'rent': 2}
```

## Going through it

### A column and a list of columns

`df["amount"]` is a `Series`, one column. `df[["city", "amount"]]` is a `DataFrame` of two: inside the square brackets there is a list, so the answer is a table even when the list holds a single name. The difference matters: `Series` and `DataFrame` have different methods, and half of all "why does this not work" is a `Series` where a table was expected.

### A mask is a column of yes and no

`df["amount"] > 10000` does not answer "yes" or "no". It answers with a **column**: one answer per row, carrying the same labels. Everything else follows:

- `big.sum()` counts the rows that matched, because `True` is a one ([lesson four](/read/py-aynymalylar-men-tipter));
- `df[big]` keeps the rows where `True` stands;
- the mask is matched to the table by label rather than by position, like everything in pandas.

### `&`, `|`, `~` and brackets that are not decoration

Two conditions are joined by signs, not by words: `&` instead of `and`, `|` instead of `or`, `~` instead of `not`. The words do not work here, and that is not a whim: `and` wants one answer, yes or no, and a column has ten of them. Python says so plainly:

```text
ValueError: The truth value of a Series is ambiguous.
```

The brackets around each condition are compulsory, because `&` runs before the comparison does: without them `df["city"] == "Kostanay" & df["kind"] == "food"` is read as `df["city"] == ("Kostanay" & df["kind"]) == "food"` — and falls over without explaining that the brackets are to blame.

### Conditions off the shelf instead of long chains

Three things that get written every day:

- `df["kind"].isin(["food", "phone"])` — instead of `(... == "food") | (... == "phone")`, and the list can be any length;
- `df["amount"].between(4000, 20000)` — instead of two comparisons, with both ends included;
- `df["city"].str.startswith("K")` — and beside it `str.contains`, `str.lower`, `str.len`: everything strings can do is on a column through `.str`.

Gaps are checked separately: `df["amount"].isna()` and `.notna()`. A comparison with `NaN` is always false, `NaN == NaN` included, so looking for gaps with a comparison finds nothing.

### Rows and columns in one go

`df.loc[mask, ["day", "amount"]]` — the left side picks the rows, the right side the columns. That is both shorter and cheaper than `df[mask][["day", "amount"]]`, which builds an intermediate table out of every column only to throw most of it away.

There is a worded form too: `df.query("amount > 5000 and city == 'Kostanay'")` — inside that string `and` and `or` do work. It suits long conditions that you wrote yourself.

What it must not be used for: **building the query string out of what a user typed**. `query` parses the expression as code, and text put into it is executed — the same hole as the SQL injection of [lesson twenty-one](/read/py-sqlite-keste-kilt-tarih). Values go in through `@variable` (`df.query("amount > @limit")`), and a condition that arrives from outside is safer built as a mask, which executes nothing.

### A filter does not renumber the labels

Three rows were filtered out and the labels stayed `[2, 6, 9]`, because these are the same rows. Which is why `loc[2]` and `iloc[2]` differ in the example: one asks for the row labelled 2, the other for the third one along.

This is the very place [lesson twenty-five](/read/py-dataframe-series-indeks-qoltanba) meant by "the difference shows up as soon as the rows are filtered". If the labels are in the way, `reset_index(drop=True)` numbers the rows afresh; `drop=True` because otherwise the old labels become a new column.

### Writing through `loc`, and why not through a filter

Marking the rows that matched goes through `loc`:

```text
df.loc[df["kind"] == "phone", "kind"] = "mobile phone"
```

The variant `df[df["kind"] == "phone"]["kind"] = "mobile phone"` looks reasonable and **does nothing**. `df[mask]` returns a new table, and the assignment changes that one rather than the original. pandas notices and prints `ChainedAssignmentError` — a warning that is easy to miss in a long run, because it does not go where you are watching for the program's output.

The rule is simple: **if there is more than one pair of square brackets to the left of the `=`, it should be `loc`**.

## The map of this lesson

![The map of this lesson: a mask, rows and columns](/static/course/py/map-select-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. What is a mask, and why is `df["amount"] > 10000` neither a yes nor a no?
2. Why are two conditions joined by the sign `&` rather than by the word `and`?
3. Why can `loc[2]` and `iloc[2]` give different rows after a filter?

## The warm-up

Three short steps before the exercise: predict, fill in, fix. The answers are at the end, but answer for yourself first.

**1. Predict.** What does this program print?

<!-- drill 1 -->
```python
import pandas as pd

amounts = pd.Series([4200, 18500, 2600, 95000], index=["food", "fuel", "food", "rent"])
mask = amounts > 5000
print(mask.tolist(), "| matched:", int(mask.sum()))
print(amounts[mask].to_dict())
```

**2. Fill in the blank.** In place of `...` put the condition: the city is Kostanay **and** the amount is over 5000.

```python
# two conditions, each in brackets of its own
import pandas as pd

df = pd.DataFrame({"city": ["Kostanay", "Rudny", "Kostanay"], "amount": [4200, 16200, 95000]})
print(df[...].to_string(index=False))
```

**3. Fix it.** The program is meant to rename the phone line to "mobile phone", but the table does not change and a `ChainedAssignmentError` turns up off to the side.

```python
# the assignment lands in a copy rather than in the table
import pandas as pd

df = pd.DataFrame({"kind": ["food", "phone"], "amount": [4200, 4990]})
df[df["amount"] > 4500]["kind"] = "mobile phone"
print(df.to_string(index=False))
```

## The exercise

**Required.** Take the same ten receipts from the lesson and answer three questions with selections rather than loops:

1. The receipts over 5000 in Kostanay — the day and the amount only.
2. Mark the receipts: `large` when the amount is 20000 or more, `ordinary` otherwise. Add the column by assignment and set the mark through `loc`. Print how many of each.
3. The sum of the large receipts per city.

The expected output:

<!-- task out -->
```text
over 5000 in Kostanay:
       day  amount
2026-01-05   18500
2026-01-15   95000

the size of a receipt:
size
ordinary    8
large       2

the large ones per city:
city
Kostanay    95000
Rudny       62000
```

Done when: the output matches line for line; the mark is set through `df.loc[mask, "size"]` rather than through a filter to the left of the `=`; the first answer picks two columns in one go; the sums come from grouping a filtered table rather than from a loop.

**On your own data.** Take a table of yours and write three conditions for it: a simple one, one made of two parts, and one with `isin` or `between`. Print how many rows matched each. If one of them matched nothing, look at `dtypes`. When a column holds strings instead of numbers, `==` quietly answers no for every row while `>` and `<` raise a `TypeError`. The cure is not a different condition but `pd.to_numeric` before the filter.

**If you feel like it.**

- Write the same condition with `query` and compare which reads better.
- Do `df[mask].reset_index(drop=True)` and see what became of the labels.
- Check how many rows `~mask` gives, and that together with the mask it comes to the length of the table.

## Where this fits the project

There is no step this time: the digest has five rows and nothing in it to filter yet. But selection is what the next step will be made of: once there are several series, the report will start with "take only the rows that matter" and only then reach the grouping of the next lesson.

Still open. Our conditions are written inline. Once there are more than three of them they will want a name — and turn into functions like `large(df)`, which read better than three brackets in a row.

## The answers

### To the questions

1. A mask is a column of yes and no, one answer per row, carrying the same labels. Comparing a column with a number compares every value rather than the column as a whole, so the answer comes out as a column too.
2. Because `and` wants one answer and a column has as many as it has rows. `&` is the sign that works element by element and gives back a column of the same shape. The brackets are needed for the same reason: `&` runs before the comparison.
3. Because a filter does not renumber the labels: the rows that stayed keep the ones they had. `loc[2]` asks for the row labelled 2, `iloc[2]` for the third one along, and after a filter those are different rows.

### To the warm-up

1. The comparison gives a column of yes and no, `sum()` counts the matches, and `amounts[mask]` keeps only those — with their labels.

<!-- drill 1 out -->
```text
[False, True, False, True] | matched: 2
{'fuel': 18500, 'rent': 95000}
```

2. `(df["city"] == "Kostanay") & (df["amount"] > 5000)`. The brackets are compulsory and the joiner is a sign, not a word.

<!-- drill 2 -->
```python
import pandas as pd

df = pd.DataFrame({"city": ["Kostanay", "Rudny", "Kostanay"], "amount": [4200, 16200, 95000]})
print(df[(df["city"] == "Kostanay") & (df["amount"] > 5000)].to_string(index=False))
```

<!-- drill 2 out -->
```text
    city  amount
Kostanay   95000
```

3. `df.loc[df["amount"] > 4500, "kind"] = "mobile phone"`. The assignment has to land in the table itself, and `df[mask]` is already another table.

<!-- drill 3 -->
```python
import pandas as pd

df = pd.DataFrame({"kind": ["food", "phone"], "amount": [4200, 4990]})
df.loc[df["amount"] > 4500, "kind"] = "mobile phone"
print(df.to_string(index=False))
```

<!-- drill 3 out -->
```text
        kind  amount
        food    4200
mobile phone    4990
```

### To the exercise

The first answer is `df.loc[mask, ["day", "amount"]]`: rows and columns in one go. The second fills the whole column with the default first and then lets `loc` overwrite it where the condition holds, so that no row is left unmarked. The third is a filter and a grouping: pick out the large ones, then add them up per city.

## Sources

- [Indexing and selecting data](https://pandas.pydata.org/docs/user_guide/indexing.html) — masks, `loc`, `iloc` and the rules they follow.
- [String methods on a column](https://pandas.pydata.org/docs/user_guide/text.html) — everything reachable through `.str`.
- [Copy on write](https://pandas.pydata.org/docs/user_guide/copy_on_write.html) — why an assignment through a filter changes nothing.
