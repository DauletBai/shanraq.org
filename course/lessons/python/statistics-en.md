# How to read official statistics

_Lead (summary):_ **The forty-first lesson of the Python course. One row of an official table carries five different numbers about the same prices, and all five are right. Which question each of them answers, how to recover the weights of the official basket that the release does not print, and why the World Bank ends up with a different figure than the news.**

## Why this matters

The last two lessons left a debt. In the thirty-ninth we counted our own inflation and compared it with the official one, promising to work out what the official basket is made of. In the fortieth the digest said "prices grew 1.6 times" — and that turned out not to be the number people quote in the news.

This lesson pays both debts. The work here is less programming than reading: open the release, understand what is written in it, and pull out what is there but not printed.

## The whole thing first

The file is `statistika.py`. Every number comes from one release — the [Bureau of National Statistics](https://stat.gov.kz/en/industries/economy/prices/), "Consumer price index and derived indicators", August 2026. The series used for changing the base is the World Bank's `FP.CPI.TOTL`.

```python
"""Lesson 41: how to read official statistics.

Every number comes from one release: the Bureau of National Statistics,
"Consumer price index and derived indicators", August 2026. The series used to
change the base is the World Bank's FP.CPI.TOTL, based at 2010 = 100.
"""

import pandas as pd

# Table 1, the row "Goods and services": the same month, different bases.
INDICES = [
    ("against July 2026", 100.6, "what happened in a month"),
    ("against December 2025", 106.4, "how much has piled up this year"),
    ("against August 2025", 109.8, "year on year — the news headline"),
    ("against December 2020", 185.8, "over five years and eight months"),
    ("January-August against January-August", 110.8, "average to average — the World Bank's way"),
]

# Table 6: the rise since the start of the year and each division's share of it.
# The names are shortened; the full ones are in the release.
DIVISIONS = [
    ("Food and drink", 4.9, 1.90),
    ("Alcohol and tobacco", 10.4, 0.15),
    ("Clothing and footwear", 6.6, 0.61),
    ("Housing, water, energy", 8.6, 0.83),
    ("Household goods", 6.0, 0.33),
    ("Health", 11.5, 0.67),
    ("Transport", 3.6, 0.32),
    ("Communications", 6.5, 0.31),
    ("Recreation and culture", 10.4, 0.35),
    ("Education", 2.7, 0.07),
    ("Restaurants and hotels", 6.1, 0.11),
    ("Insurance and finance", 3.0, 0.02),
    ("Personal care and other", 9.6, 0.68),
]
HEADLINE = 6.4           # the rise of "goods and services" since January, table 6

# Table 3: core inflation, August against August.
CORE = [
    ("without fruit, vegetables, petrol and coal", 110.7),
    ("without fruit and vegetables", 110.6),
    ("without those, utilities, transport and telecoms", 111.2),
]

# Table 2: the same month by region, August against August.
REGIONS = [("Republic of Kazakhstan", 109.8),
           ("North Kazakhstan", 112.0),
           ("Karaganda", 108.1)]

# The World Bank, FP.CPI.TOTL: the price index, 2010 = 100.
SERIES = pd.Series({2019: 189.30, 2020: 202.02, 2021: 218.27, 2022: 251.07,
                    2023: 287.54, 2024: 312.53, 2025: 348.12})

print("== One row, five numbers")
for base, value, question in INDICES:
    print(f"  {value:6.1f}  {base:38} — {question}")

print()
print("== What the basket is made of")
# The weights are not printed in the release, but they show: a contribution is
# the weight times the rise, so the weight is the contribution divided by it.
basket = pd.DataFrame(DIVISIONS, columns=["division", "rise", "contribution"])
basket["weight"] = (basket["contribution"] / basket["rise"] * 100).round(1)
print(basket.sort_values("weight", ascending=False)[["division", "weight", "rise", "contribution"]]
      .to_string(index=False))
print("  the weights add up to:", round(basket["weight"].sum(), 1), "% — the whole basket")

print()
print("== Who raised the index this year")
top = basket.sort_values("contribution", ascending=False).head(3)
print("  three divisions out of thirteen:", ", ".join(top["division"]))
print("  their contribution:", round(top["contribution"].sum(), 2),
      "of", round(basket["contribution"].sum(), 2))
print("  the contributions add up to:", round(basket["contribution"].sum(), 2),
      "against the headline", HEADLINE, "— the gap is rounding")

print()
print("== Core inflation: the same index without vegetables and petrol")
print(f"  {REGIONS[0][1]:6.1f}  the whole index")
for name, value in CORE:
    print(f"  {value:6.1f}  {name}")

print()
print("== One country, different places")
places = pd.DataFrame(REGIONS, columns=["place", "index"])
print(places.to_string(index=False))
print("  the spread:", round(places["index"].max() - places["index"].min(), 1), "points")

print()
print("== Changing the base: the levels move, the growth does not")
rebased = (SERIES / SERIES[2019] * 100).round(2)
print(pd.DataFrame({"2010 = 100": SERIES, "2019 = 100": rebased}).to_string())
print("  growth 2019 -> 2025 on the first base:", round(SERIES[2025] / SERIES[2019], 3))
print("  on the second base:", round(rebased[2025] / rebased[2019], 3))
```

The output:

```text
== One row, five numbers
   100.6  against July 2026                      — what happened in a month
   106.4  against December 2025                  — how much has piled up this year
   109.8  against August 2025                    — year on year — the news headline
   185.8  against December 2020                  — over five years and eight months
   110.8  January-August against January-August  — average to average — the World Bank's way

== What the basket is made of
               division  weight  rise  contribution
         Food and drink    38.8   4.9          1.90
 Housing, water, energy     9.7   8.6          0.83
  Clothing and footwear     9.2   6.6          0.61
              Transport     8.9   3.6          0.32
Personal care and other     7.1   9.6          0.68
                 Health     5.8  11.5          0.67
        Household goods     5.5   6.0          0.33
         Communications     4.8   6.5          0.31
 Recreation and culture     3.4  10.4          0.35
              Education     2.6   2.7          0.07
 Restaurants and hotels     1.8   6.1          0.11
    Alcohol and tobacco     1.4  10.4          0.15
  Insurance and finance     0.7   3.0          0.02
  the weights add up to: 99.7 % — the whole basket

== Who raised the index this year
  three divisions out of thirteen: Food and drink, Housing, water, energy, Personal care and other
  their contribution: 3.41 of 6.35
  the contributions add up to: 6.35 against the headline 6.4 — the gap is rounding

== Core inflation: the same index without vegetables and petrol
   109.8  the whole index
   110.7  without fruit, vegetables, petrol and coal
   110.6  without fruit and vegetables
   111.2  without those, utilities, transport and telecoms

== One country, different places
                 place  index
Republic of Kazakhstan  109.8
      North Kazakhstan  112.0
             Karaganda  108.1
  the spread: 3.9 points

== Changing the base: the levels move, the growth does not
      2010 = 100  2019 = 100
2019      189.30      100.00
2020      202.02      106.72
2021      218.27      115.30
2022      251.07      132.63
2023      287.54      151.90
2024      312.53      165.10
2025      348.12      183.90
  growth 2019 -> 2025 on the first base: 1.839
  on the second base: 1.839
```

## The walk-through

### Five numbers are not five opinions

```text
100.6  against July 2026
106.4  against December 2025
109.8  against August 2025
185.8  against December 2020
110.8  January-August against January-August
```

One row of a table, one month, one office. Nobody is arguing and nobody made a mistake: an index always has a **base**, a moment it is compared against, and the question "what is inflation" has no answer without one.

After that it is simply a thing to keep in mind. "Inflation is 9.8 %" is August against August. "Prices are up 6.4 % since the start of the year" is August against December. "Prices nearly doubled in five years" is the 185.8 against December 2020. All three sentences are about one table, and none of them refutes the others.

### The weights are not printed, but they show

The release carries a table of contributions: for every division, the rise in prices and that division's share of the total rise. The weights are not in it. But a contribution is counted exactly as in lesson thirty-nine: `contribution = weight × rise`. So the weight comes out of a division:

```text
               division  weight  rise  contribution
         Food and drink    38.8   4.9          1.90
 Housing, water, energy     9.7   8.6          0.83
```

And the check that says we understood it right: the weights add up to 99.7 %. That is the whole basket, to the accuracy of one decimal place.

There is the answer to the question of lesson thirty-nine. In the official basket food takes 38.8 % and housing with utilities takes 9.7 %. For somebody who rents, the share of housing is four to eight times larger, which is why their inflation is under no obligation to match this one.

### Who raised the index

```text
  three divisions out of thirteen: Food and drink, Housing, water, energy, Personal care and other
  their contribution: 3.41 of 6.35
```

The same decomposition by contribution as on a personal basket, only at the level of a country — and the result has the same shape: half of the rise is explained by three divisions out of thirteen. Notice that health rose the most of all (11.5 %) and stands fourth in the contributions: a weight of 5.8 % does not let it climb higher.

The contributions add up to 6.35 against a headline of 6.4. That is neither an error nor a sleight of hand: both are printed rounded to one decimal, and the sum of thirteen roundings is under no obligation to match the rounding of a sum. The 0.05 of a point is what printing costs.

### Core inflation: what is taken out, and why

```text
109.8  the whole index
110.7  without fruit, vegetables, petrol and coal
111.2  without those, utilities, transport and telecoms
```

Core inflation is the same index without the items that jump for reasons that have nothing to do with money: a harvest, the world price of oil, a decision about a tariff. It is counted to see the persistent part of the rise.

Notice the direction: without vegetables and petrol the figure is **larger**, not smaller. The habit of thinking that "core" means "smaller and prettier" does not work here: in this month it was the volatile items that were pulling the index down.

### "Inflation in the country" is an average over places too

```text
                 place  index
Republic of Kazakhstan  109.8
      North Kazakhstan  112.0
             Karaganda  108.1
```

The spread between regions is 3.9 points, and this is the very same month. The office counts and publishes regional indices separately; the national one is their aggregate, weighted by region. So "inflation of 9.8 %" averages not only over goods but over geography.

### What a revision does

A revision is when what changes is not the price but the thing that measured it. Two things change regularly.

**The base.** The current index is built on "December 2020 = 100"; at the World Bank the same Kazakhstan sits on "2010 = 100". The arithmetic is simple and the program shows all of it: on another base every level changes, while the growth between two dates does not change at all — 1.839 either way. When two sources give different "indices", the first thing to compare is the bases and the second is the period.

**The basket.** The weights are revised because people change what they buy. The source says so itself — here is the definition the World Bank gives its own inflation indicator: "…the cost to the average consumer of acquiring a basket of goods and services **that may be fixed or changed at specified intervals**, such as yearly".

What that does to the numbers follows from the same formula. If a contribution is `weight × rise`, then with the very same prices a different set of weights gives a different headline. There is nothing in this lesson to compare two vintages of weights against: we have one release in hand. But knowing that the headline depends on more than prices is not optional.

### Why the World Bank's number is not the headline

The last of the five columns is "January-August against January-August", 110.8. That is the average over a period against the average over the same period a year earlier, and that is exactly the annual inflation the World Bank publishes: `FP.CPI.TOTL.ZG`.

Hence a divergence that otherwise looks like somebody's mistake. The project's digest takes the World Bank series and says one number; the news takes December against December and says another. Both are about one year and one country. Different questions, different answers.

### What this count cannot do

The recovered weights are the weights **of today's release**, and they are good for understanding the structure rather than for recounting history. Their accuracy is exactly what rounding to one decimal allows: the weight of a division with a small rise comes out of dividing two short numbers, and the error there is larger than for the big divisions. The check by the sum is the only thing standing between you and a gross mistake.

## The lesson map

![Lesson map: one table, five questions, weights out of contributions](/static/course/py/map-statistics-en.svg)

## Say it in your own words

Answer out loud or on paper without looking. The answers are at the end of the lesson.

1. Why do five numbers in one row of an official table not contradict each other?
2. How do you get a division's weight if it is not printed?
3. What changes and what does not when the base of an index changes?

## Warm-up

Three short steps before the task: predict, fill in, fix. The answers are at the end of the lesson, but answer for yourself first.

**1. Predict.** What does the program print, and will the two "growths" be the same?

<!-- drill 1 -->
```python
index = {2019: 189.30, 2025: 348.12}
print("growth:", round(index[2025] / index[2019], 3))
rebased = {year: value / index[2019] * 100 for year, value in index.items()}
print("on the 2019 base:", {year: round(value, 2) for year, value in rebased.items()})
print("growth:", round(rebased[2025] / rebased[2019], 3))
```

**2. Fill in the blank.** In place of `...` recover the division's weight from its contribution and its rise.

```python
# housing services: the rise since January and its share of the total rise
rate, contribution = 8.6, 0.83
weight = ...
print("the division's weight:", round(weight, 1), "%")
```

**3. Fix it.** The program takes the index "August 2026 against December 2020" and calls it annual inflation.

```python
index = 185.8
print("annual inflation:", round(index - 100, 1), "%")
```

## The task

**Required.** From the table of contributions, recover the divisions' weights and print them by contribution, largest first. Count how many divisions from the top it takes to cover half of the rise, and what share of the rise they give. Compare the sum of the contributions with the headline rise and name the gap. And separately: put the World Bank series on the 2019 base and make sure the growth did not change because of it.

The expected output:

<!-- task out -->
```text
== The weights, recovered from the contributions
               division  weight  rise  contribution
         Food and drink    38.8   4.9          1.90
 Housing, water, energy     9.7   8.6          0.83
Personal care and other     7.1   9.6          0.68
                 Health     5.8  11.5          0.67
  Clothing and footwear     9.2   6.6          0.61
 Recreation and culture     3.4  10.4          0.35
        Household goods     5.5   6.0          0.33
              Transport     8.9   3.6          0.32
         Communications     4.8   6.5          0.31
    Alcohol and tobacco     1.4  10.4          0.15
 Restaurants and hotels     1.8   6.1          0.11
              Education     2.6   2.7          0.07
  Insurance and finance     0.7   3.0          0.02
  the weights add up to: 99.7 %

== Half of the rise
  divisions needed: 3
  they are: Food and drink, Housing, water, energy, Personal care and other
  their share of the rise: 53.7 %

== The headline and the sum of the contributions
  the contributions add up to: 6.35
  the headline: 6.4
  the gap: 0.05 points — rounding

== Changing the base
  base 2010: 1.839
  base 2019: 1.839
  the growth did not change: True
```

Done means: the output matches line for line; the weights are counted from the contributions rather than taken from somewhere else; "half of the rise" is found by a running sum rather than picked out by eye; the gap with the headline is printed as a number.

**On your own data.** Download the latest CPI release and repeat the count on its numbers: the weights move from month to month, because the contributions are counted from the start of the year. Compare the weights you get with the ones in the lesson and see which divisions move the most. Then take your own basket from lesson thirty-nine and stand its weights next to the official ones — the difference between your inflation and the headline lives exactly there.

**If you want more.**

- Count what the headline would be if food had a weight of 20 % and housing 30 %, with the same rises by division.
- Pull the regional table out of the release and find the region furthest from the country — in both directions.
- Do to core inflation what the lesson does to the overall index: count how many points each of the three versions differs from the whole index by.

## Where this goes in the project

Step twenty, and it is about the honesty of the report rather than about new numbers. The digest has printed an annual rate of inflation since its first step and has never once said which one.

Now it says. `derekkoz.about()` asks the World Bank what its own indicator is and puts the answer next to the series, in `data/indicators.json`; `bet.py` prints those definitions at the foot of the page, under everything they explain. The definition travels with the number instead of living in somebody's head.

The step also shows how to treat what is optional: a note is decoration, not data. A source that does not answer costs the page its footer and nothing more — the log carries a line about what is missing and the report is built anyway. And `--offline` stays a promise: with no network and no cache the digest does not go to the network for a paragraph of text.

## The answers

### To the questions

1. Because each has its own base: month against month, against December, year on year, against December 2020, and the average of a period against an average. An index without its base is not a number but half of one.
2. Weight = contribution ÷ rise. A division's contribution to the total rise is its weight times its own rise, so the inverse operation gives the weight back. The check is that the weights add up to about 100 %.
3. The levels change: on the 2019 base the series starts at a hundred rather than at 189.3. The growth between any two dates does not change: 1.839 stays 1.839. That is why comparing an "index" from two sources is useless until the bases are lined up, while comparing rates of growth is fine.

### To the warm-up

1. Both growths are the same. Changing the base is dividing the whole series by one number, and the ratio of two points does not change because of that.

<!-- drill 1 out -->
```text
growth: 1.839
on the 2019 base: {2019: 100.0, 2025: 183.9}
growth: 1.839
```

2. `contribution / rate * 100`. The contribution and the rise are in per cent, and the weight is wanted as a per cent of the basket — hence the hundred.

<!-- drill 2 -->
```python
# housing services: the rise since January and its share of the total rise
rate, contribution = 8.6, 0.83
weight = contribution / rate * 100
print("the division's weight:", round(weight, 1), "%")
```

<!-- drill 2 out -->
```text
the division's weight: 9.7 %
```

3. It is not a year but 68 months. Over the whole stretch prices rose by 85.8 %; to get "a year on average" takes a root rather than a division — the multiplier is raised to the power of 12/68.

<!-- drill 3 -->
```python
# the index "August 2026 against December 2020" covers 68 months, not a year
index = 185.8
months = 68
print("over the whole stretch:", round(index - 100, 1), "%")
print("a year on average:", round(((index / 100) ** (12 / months) - 1) * 100, 1), "%")
```

<!-- drill 3 out -->
```text
over the whole stretch: 85.8 %
a year on average: 11.6 %
```

### To the task

Half of the rise is covered by three divisions out of thirteen — 53.7 %. It is found with a running sum: sort by contribution, walk from the top, and count how many steps it takes to pass the half. By eye the answer would be the same, but only on a table of thirteen rows; on a hundred items the eyes run out.

The gap between the sum of the contributions and the headline is 0.05 of a point. Telling that kind of gap apart from an error is a skill: it is smaller than one unit of the last digit and it does not grow. Had the sum differed from the headline by half a point, the matter would not have been rounding but the fact that we added up the wrong things.

## Sources

- [Consumer price index and derived indicators, Bureau of National Statistics](https://stat.gov.kz/en/industries/economy/prices/) — the release every number in the lesson comes from.
- [Inflation, consumer prices (annual %), World Bank](https://data.worldbank.org/indicator/FP.CPI.TOTL.ZG) — the definition of the indicator the project's digest prints.
- [Consumer price index (2010 = 100), World Bank](https://data.worldbank.org/indicator/FP.CPI.TOTL) — the series the change of base is shown on.
- [ILO consumer price index manual](https://www.ilo.org/publications/consumer-price-index-manual-theory-and-practice) — how weights, bases and revisions are built, at length.
