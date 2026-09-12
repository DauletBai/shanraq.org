# The mean, the median and the spread: when the mean lies

_Лид (summary):_ **Lesson thirty-seven of the Python course, and the start of the module about money. Ten countries, one year: the mean is 20.48 %, the median 13.88 % — and the mean is higher than what eight of the ten actually had. `mean`, `median`, `std`, the quarters, and the rule that picks a measure before anyone has seen the result.**

## Why this is needed

A series of eleven numbers does not fit in a headline, so it is folded into one. The question is not whether to fold it, but **what exactly is lost** when you do — and whether the truth goes with it.

"Average inflation across the region: 20 %" sounds like a description of the region. On the real numbers of 2022 it describes exactly one country out of ten: eight of them were below that mean, and what dragged it up there was Turkiye at seventy-two per cent.

This lesson is about three measures that answer three different questions, and about the rule for which one goes in the headline.

## The whole thing at once

The file `average.py`. Two real series — inflation in Kazakhstan by year, and inflation in ten countries in 2022 — both from [the World Bank](https://data.worldbank.org/indicator/FP.CPI.TOTL.ZG).

```python
"""Lesson 37: the mean, the median and the spread -- and when the mean lies.

Two real series: yearly inflation in Kazakhstan, and inflation in ten countries
in 2022. Both come from the World Bank, indicator FP.CPI.TOTL.ZG.
"""

import numpy as np
import pandas as pd

# Kazakhstan, yearly inflation in per cent, eleven years.
kazakhstan = pd.Series(
    [6.85, 6.68, 14.36, 7.44, 6.16, 5.33, 6.72, 8.04, 15.03, 14.53, 8.69],
    index=[2014, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024],
    name="inflation",
)

# The same indicator, a different slice: one year and ten countries.
year_2022 = pd.Series(
    {
        "Kazakhstan": 15.03, "Russia": 13.74, "Uzbekistan": 11.45,
        "Kyrgyzstan": 13.92, "Turkiye": 72.31, "Georgia": 11.90,
        "Armenia": 8.64, "Azerbaijan": 13.85, "Belarus": 15.21,
        "Moldova": 28.74,
    },
    name="inflation 2022",
)

print("== Eleven years of one country")
print("  mean:", round(kazakhstan.mean(), 2))
print("  median:", round(kazakhstan.median(), 2))
print("  spread:", round(kazakhstan.std(), 2))
print(f"  from {kazakhstan.min()} to {kazakhstan.max()}, range {round(kazakhstan.max() - kazakhstan.min(), 2)}")

print()
print("== One year, ten countries")
mean_2022 = year_2022.mean()
median_2022 = year_2022.median()
print("  mean:", round(mean_2022, 2))
print("  median:", round(median_2022, 2))
print("  below the mean:", int((year_2022 < mean_2022).sum()), "of", len(year_2022))
print("  below the median:", int((year_2022 < median_2022).sum()), "of", len(year_2022))

print()
print("== What one country does to them")
without = year_2022.drop("Turkiye")
print("  with Turkiye:    mean", round(mean_2022, 2), "| median", round(median_2022, 2))
print("  without Turkiye: mean", round(without.mean(), 2), "| median", round(without.median(), 2))
# The shift is computed from the numbers the reader sees above, not from hidden
# ones: otherwise "moved by 0.04" would not match 13.88 - 13.85.
print("  the mean moved by", round(round(mean_2022, 2) - round(without.mean(), 2), 2))
print("  the median moved by", round(round(median_2022, 2) - round(without.median(), 2), 2))

print()
print("== Spread: two answers to one question")
quarters = year_2022.quantile([0.25, 0.5, 0.75])
print("  std:", round(year_2022.std(), 2))
print("  quarters:", round(quarters[0.25], 2), "|", round(quarters[0.5], 2), "|", round(quarters[0.75], 2))
print("  interquartile range:", round(quarters[0.75] - quarters[0.25], 2))

print()
print("== One word, two numbers")
print("  pandas .std():", round(kazakhstan.std(), 4), "- divides by n - 1")
print("  numpy  .std():", round(float(np.std(kazakhstan)), 4), "- divides by n")
print("  pandas std(ddof=0):", round(kazakhstan.std(ddof=0), 4))

print()
print("== One line instead of six")
print(year_2022.describe().round(2).to_string())
```

It prints:

```text
== Eleven years of one country
  mean: 9.08
  median: 7.44
  spread: 3.69
  from 5.33 to 15.03, range 9.7

== One year, ten countries
  mean: 20.48
  median: 13.88
  below the mean: 8 of 10
  below the median: 5 of 10

== What one country does to them
  with Turkiye:    mean 20.48 | median 13.88
  without Turkiye: mean 14.72 | median 13.85
  the mean moved by 5.76
  the median moved by 0.03

== Spread: two answers to one question
  std: 18.97
  quarters: 12.36 | 13.88 | 15.16
  interquartile range: 2.81

== One word, two numbers
  pandas .std(): 3.6851 - divides by n - 1
  numpy  .std(): 3.5136 - divides by n
  pandas std(ddof=0): 3.5136

== One line instead of six
count    10.00
mean     20.48
std      18.97
min       8.64
25%      12.36
50%      13.88
75%      15.16
max      72.31
```

## How it works

### The mean answers "how much each would get if it were shared out"

The mean is the sum divided equally. Of the three measures it is the only one with a direct meaning when the **sum** is what matters: the total bill, the total spend, the total tax collected. If the question is "how much altogether and across how many", the mean answers it and nothing else will.

The price of that is sensitivity. One large value pulls the mean towards itself, the harder the larger it is and the fewer observations there are.

### The median answers "and what about the middle"

The median is the value with half the series above it and half below. No single number shifts it, however large: all that matters is which side of the middle it is on.

The program measures this. Drop Turkiye from the ten:

```text
with Turkiye:    mean 20.48 | median 13.88
without Turkiye: mean 14.72 | median 13.85
```

The mean moved by 5.76 percentage points, the median by 0.03. One country out of ten rewrote the "average" answer by almost a third, and the "middle" answer never noticed.

### When the mean lies

Not always, and not by itself: the mean lies when it is offered as the **typical** value and the series is lopsided. The check takes one line:

```text
below the mean: 8 of 10
below the median: 5 of 10
```

The median splits a series in half by definition — half is always below it. If eight of ten came out below the mean, the series has a long tail upwards, and "20 % on average" describes not the region but its single extreme.

The same skew shows in Kazakhstan's eleven years: a mean of 9.08 against a median of 7.44 — three spikes out of eleven years lifted the mean nearly a point and a half above the middle.

### Spread: one number against four

`std` is the average distance of the values from the mean, and it has exactly the same weakness: an outlier inflates it along with the mean. For the ten countries `std` is 18.97 while the median is 13.88 — the spread is larger than the middle itself.

The quarters work differently. The first quarter is the value a quarter of the series falls below; the third is the one three quarters fall below; the distance between them is called the interquartile range:

```text
quarters: 12.36 | 13.88 | 15.16
interquartile range: 2.81
```

Eight countries of ten fit into a band less than three points wide. That is the real spread of this series — and 18.97 describes not the spread but the distance to Turkiye.

The range (`max - min`) is useful for something else: it names the **limits**, not the typical distance. In a headline it is honest only next to the number of observations.

### One word, two numbers

The pair of lines behind figures that disagree between two people holding the same series:

```text
pandas .std(): 3.6851 - divides by n - 1
numpy  .std(): 3.5136 - divides by n
```

Neither library is wrong. When the series is the **whole** population, you divide by `n`; when it is a sample used to judge something larger, you divide by `n − 1`. pandas assumes your data is a sample by default, numpy assumes it is the population.

A difference in the third decimal looks like a detail right up to the day someone who used the other library sees it. So a report says which spread it printed, and the code passes `ddof` explicitly.

### `describe()` — one line instead of six

`Series.describe()` prints the count, the mean, the spread, the minimum, the three quarters and the maximum. It is the first thing to do with an unfamiliar series: eight numbers show the centre, the spread and the skew at once.

A report takes two or three of them. Looking at the data yourself, you read all eight.

### What goes in the headline

A rule worth writing into the code rather than keeping in your head:

1. A question about the **sum** — the mean. "How much was collected and how much per person" is answered by nothing else.
2. A question about the **typical** value — the median, if the mean has visibly moved away from it; otherwise it makes no difference, and then the mean is the more familiar of the two.
3. Next to any measure — the **number of observations**: 20 % across ten countries and 20 % across a hundred are different statements.
4. If the series is lopsided, one measure is not enough: the median with the quarters says what the mean with `std` hides.

## The lesson map

![The lesson map: the mean, the median and the spread](/static/course/py/map-average-en.svg)

## Say it in your own words

Without looking, answer aloud or on paper. The answers are at the end of the lesson.

1. Why did the median barely move when Turkiye was dropped from the series, while the mean moved by 5.76?
2. What does "eight of ten below the mean" mean, and what does it say about the series?
3. Why do `.std()` in pandas and in numpy give different numbers for the same series?

## Warm-up

Three short steps before the exercise: predict, fill in, fix. The answers are at the end of the lesson, but answer for yourself first.

**1. Predict.** What does this program print?

<!-- drill 1 -->
```python
import pandas as pd

pay = pd.Series([120000, 130000, 140000, 150000, 900000])
print("mean:", int(pay.mean()))
print("median:", int(pay.median()))
print("below the mean:", int((pay < pay.mean()).sum()), "of", len(pay))
```

**2. Fill in the blank.** In place of `...` put the measure one expensive receipt cannot shift.

```python
# five receipts, one of them a fridge
import pandas as pd

prices = pd.Series([320, 340, 350, 360, 4800])
measure = ...
print("the measure one expensive receipt did not shift:", measure)
```

**3. Fix it.** The program prints two different spreads for one series and declares one of the numbers a bug.

```python
# "one of the libraries counts it wrong" - no, they count different things
import numpy as np
import pandas as pd

row = pd.Series([6.85, 6.68, 14.36, 7.44, 6.16])
print("pandas:", round(row.std(), 4))
print("numpy: ", round(float(np.std(row)), 4))
print("equal:", round(row.std(), 10) == round(float(np.std(row)), 10))
```

## Exercise

**Required.** Take three series: Kazakhstan's eleven years, the calm five of them (2017–2021), and the ten countries in 2022. For each print the number of observations, the mean, the median, the spread, the limits and the interquartile range — and the measure that goes in the headline. The measure is chosen by a **rule in the code**: if the mean has moved away from the median by more than a tenth of the median, the headline takes the median. At the end, `describe()` for all three side by side.

The expected output:

<!-- task out -->
```text
== Kazakhstan, 2014-2024
  observations: 11
  mean: 9.08 | median: 7.44
  spread: 3.69 | from 5.33 to 15.03
  interquartile range: 4.82
  for the headline: median 7.44 %
  below the mean: 8 of 11

== Kazakhstan, 2017-2021
  observations: 5
  mean: 6.74 | median: 6.72
  spread: 1.06 | from 5.33 to 8.04
  interquartile range: 1.28
  for the headline: mean 6.74 %
  below the mean: 3 of 5

== Ten countries, 2022
  observations: 10
  mean: 20.48 | median: 13.88
  spread: 18.97 | from 8.64 to 72.31
  interquartile range: 2.81
  for the headline: median 13.88 %
  below the mean: 8 of 10

== Why not one measure for everything
       11 years  5 calm  10 countries
count     11.00    5.00         10.00
mean       9.08    6.74         20.48
std        3.69    1.06         18.97
min        5.33    5.33          8.64
25%        6.70    6.16         12.36
50%        7.44    6.72         13.88
75%       11.52    7.44         15.16
max       15.03    8.04         72.31
```

Done when: the output matches line for line; a function picks the measure, not your eye; the rule can say "mean" too — on the calm five years it does; the number of observations is printed next to the measure rather than apart from it.

**On your own data.** Take any series of your own — daily spending, travel time, weight, electricity bills. Compute the mean and the median. If they differ by more than a tenth, find the value that separates them and decide: is it a typing error, a rare event, or ordinary life? The answer decides whether it goes or stays — and that decision is worth writing down beside the number.

**If you feel like it.**

- Compute `std` with `ddof=0` and `ddof=1` on your own series and see at how many observations the difference stops being visible.
- Build `quantile([0.1, 0.5, 0.9])` on your series and compare the "eight of ten" band with the `min`–`max` limits.
- Take the series with missing values from [lesson thirty](/read/py-las-derek-isna-fillna-astype) and check that `mean()` counts the non-empty values rather than the length of the series.

## Where this goes in the project

The seventeenth step, and the fifth module starts with it: the digest stops showing only the latest year.

`sholu/esep.py` gains `ozara` — per country it counts the observations, the mean, the median, the spread and the interquartile range, and puts beside them the measure chosen by the same rule as in the exercise. The report page gains a second block, "over all the years", under the table of the last one.

What the rule said on the digest's real data:

```text
KAZ  5 years  mean 11.54  median 11.39  spread 3.22  -> the mean
RUS  5 years  mean 8.69   median 8.43   spread 3.06  -> the mean
UZB  5 years  mean 10.14  median 9.96   spread 1.04  -> the mean
```

It chose the **mean** for all three: their five-year series are symmetrical enough that the median sits a per cent or two from the mean. That is the use of a rule written into the code — it does not slip in a "careful" median where the data does not ask for one, and it will say the opposite the day a spike like Turkiye's appears in a series.

Debts. Five points is a short series: a median over five observations is coarse, and a spread over five years is an estimate with a wide error the digest names nowhere. That is the same conversation about the number of observations beside a measure, and it returns in the lesson on reading official statistics.

## Answers

### To the questions

1. Because the median is set by a value's **position**, not its size. Turkiye is above the middle either way, and what exactly stands there — 28 or 72 — the middle does not care. The mean adds every value up, so one large number pulls it upwards the harder the larger it is.
2. That the series is lopsided: it has a long tail upwards. Half is always below the median — that is its definition; if eight of ten came out below the mean, the mean has moved off towards the tail and is no longer the typical value.
3. Because they divide by different things: pandas treats the series as a sample by default and divides by `n − 1`, numpy treats it as the whole population and divides by `n`. Neither is a mistake; the mistake is not saying in the report which was computed, and leaving `ddof` at its default.

### To the warm-up

1. One salary of nine hundred thousand dragged the mean above what four of the five actually earn.

<!-- drill 1 out -->
```text
mean: 288000
median: 140000
below the mean: 4 of 5
```

2. `prices.median()`. The median looks at position, not size: the fridge stays to the right of the middle whatever its price.

<!-- drill 2 -->
```python
import pandas as pd

prices = pd.Series([320, 340, 350, 360, 4800])
measure = prices.median()
print("the measure one expensive receipt did not shift:", measure)
```

<!-- drill 2 out -->
```text
the measure one expensive receipt did not shift: 350.0
```

3. It is not a bug but `ddof`: give both libraries the same divisor.

<!-- drill 3 -->
```python
import numpy as np
import pandas as pd

row = pd.Series([6.85, 6.68, 14.36, 7.44, 6.16])
print("pandas:", round(row.std(ddof=0), 4))
print("numpy: ", round(float(np.std(row)), 4))
print("equal:", round(row.std(ddof=0), 10) == round(float(np.std(row)), 10))
```

<!-- drill 3 out -->
```text
pandas: 3.0584
numpy:  3.0584
equal: True
```

### To the exercise

The rule is in a `headline` function not for tidiness: while it lives in your head, every series gets whichever measure looks better. Written into the code, it answers all three the same way — and on the calm five years it honestly says "mean", though the median would have been the more cautious choice.

The quarters in the `describe()` table show what neither the mean nor the spread does: half of the ten countries fit between 12.36 and 15.16, while an `std` of 18.97 describes the distance to one value rather than to the other nine.

## Sources

- [Descriptive statistics in pandas](https://pandas.pydata.org/docs/user_guide/basics.html#descriptive-statistics) — `mean`, `median`, `std`, `quantile`, `describe` and what they do with missing values.
- [numpy.std](https://numpy.org/doc/stable/reference/generated/numpy.std.html) — the `ddof` parameter and why it defaults to zero.
- [Consumer price inflation, World Bank](https://data.worldbank.org/indicator/FP.CPI.TOTL.ZG) — the source of every number in this lesson.
