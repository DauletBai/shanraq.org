# The first model: linear regression

_Lead (summary):_ **The forty-third lesson of the Python course. A line drawn through sixteen years of prices gives an R² of 0.927 and loses to the rule "same as last year" — by a hundredth of a point. Its error is not random but has a shape: plus, minus, plus. A line in logarithms is wrong by half as much, because prices multiply.**

## Why this matters

The link from the last lesson says "these two columns moved together". A model says more: **by how much**, and **what happens at this value**. The first of them, and the most honest, is a straight line.

What it costs is an error that can now be named as a number. This lesson is less about training a model — two lines of code — than about looking at its misses and not being fooled by a pretty R².

## The whole thing first

The file is `model.py`. The series is real: the consumer price index of Kazakhstan since 2010, from the [World Bank](https://data.worldbank.org/indicator/FP.CPI.TOTL). The library is new — scikit-learn, `pip install scikit-learn==1.9.1`.

```python
"""Lesson 43: the first model -- a straight line through prices.

The series is real: the consumer price index of Kazakhstan, 2010 = 100, the
World Bank indicator FP.CPI.TOTL. The library is scikit-learn.
"""

import numpy as np
import pandas as pd
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_absolute_error, r2_score

CPI = pd.Series({
    2010: 100.00, 2011: 108.45, 2012: 114.09, 2013: 120.87, 2014: 129.15,
    2015: 137.77, 2016: 157.56, 2017: 169.28, 2018: 179.72, 2019: 189.30,
    2020: 202.02, 2021: 218.27, 2022: 251.07, 2023: 287.54, 2024: 312.53,
    2025: 348.12,
})

# sklearn expects a table of features: a row per observation, a column per
# feature. There is one feature here, the year, so there is one column -- but a
# column, not a series.
X = CPI.index.to_numpy().reshape(-1, 1)
y = CPI.to_numpy()

print("== A line through sixteen years")
line = LinearRegression().fit(X, y)
print("  the index gains a year:", round(float(line.coef_[0]), 2), "points")
print("  R²:", round(r2_score(y, line.predict(X)), 3))
print("  the average error:", round(mean_absolute_error(y, line.predict(X)), 2), "points")

print()
print("== The error, not only the fit")
table = pd.DataFrame({"index": CPI, "line": line.predict(X).round(1)})
table["error"] = (table["index"] - table["line"]).round(1)
print(table.to_string())
worst = table["error"].abs().idxmax()
print("  the worst year:", worst, "— off by", round(float(table.loc[worst, "error"]), 1), "points")

print()
print("== The error has a shape")
signs = "".join("+" if e > 0 else "-" for e in table["error"])
print("  the signs of the errors by year:", signs)
print("  runs of the same sign:", 1 + sum(1 for a, b in zip(signs, signs[1:]) if a != b))
print("  a random error would have about", round(len(signs) / 2 + 1))

print()
print("== The line in logarithms")
log_line = LinearRegression().fit(X, np.log(y))
growth = float(np.exp(log_line.coef_[0]) - 1)
predicted = np.exp(log_line.predict(X))
print("  a year multiplies prices by:", round(float(np.exp(log_line.coef_[0])), 4),
      "— that is", round(growth * 100, 2), "% a year")
print("  the average error:", round(mean_absolute_error(y, predicted), 2), "points")
print("  the worst miss:", round(float(np.max(np.abs(y - predicted))), 2), "points")

print()
print("== A prediction and what it costs")
next_year = np.array([[2026]])
print("  the line promises for 2026:", round(float(line.predict(next_year)[0]), 1))
print("  the line in logarithms:", round(float(np.exp(log_line.predict(next_year)[0])), 1))
print("  the gap between the promises:",
      round(float(np.exp(log_line.predict(next_year)[0]) - line.predict(next_year)[0]), 1), "points")
```

The output:

```text
== A line through sixteen years
  the index gains a year: 15.45 points
  R²: 0.927
  the average error: 17.17 points

== The error, not only the fit
       index   line  error
2010  100.00   73.2   26.8
2011  108.45   88.7   19.8
2012  114.09  104.1   10.0
2013  120.87  119.6    1.3
2014  129.15  135.0   -5.8
2015  137.77  150.5  -12.7
2016  157.56  165.9   -8.3
2017  169.28  181.4  -12.1
2018  179.72  196.8  -17.1
2019  189.30  212.3  -23.0
2020  202.02  227.7  -25.7
2021  218.27  243.2  -24.9
2022  251.07  258.6   -7.5
2023  287.54  274.1   13.4
2024  312.53  289.5   23.0
2025  348.12  305.0   43.1
  the worst year: 2025 — off by 43.1 points

== The error has a shape
  the signs of the errors by year: ++++---------+++
  runs of the same sign: 3
  a random error would have about 9

== The line in logarithms
  a year multiplies prices by: 1.0849 — that is 8.49 % a year
  the average error: 7.39 points
  the worst miss: 24.09 points

== A prediction and what it costs
  the line promises for 2026: 320.4
  the line in logarithms: 351.5
  the gap between the promises: 31.1 points
```

## The walk-through

### A model is two lines

```python
model = LinearRegression().fit(X, y)
model.predict(X)
```

`fit` picks the line, `predict` asks it. Everything else in the lesson is about working out whether it will do.

One subtlety about the shape of the data: `X` has to be a **table** of features — a row per observation, a column per feature — hence `reshape(-1, 1)`. There is one feature here, the year, but the column is needed all the same: sklearn does not guess whether you have one feature or a hundred.

### R² is pretty and means almost nothing

```text
  R²: 0.927
  the average error: 17.17 points
```

R² is the share of the spread the model explained: 0.927 sounds like "the model is 93 % right". What actually stands beside it is an average error of 17 points on an index that started at a hundred.

The reason is simple: R² compares the model with the weakest thing there is — a horizontal line at the mean. Beating that is easy for any series that grows. So R² is almost never the answer to "is this model any good".

### The error has a shape

```text
  the signs of the errors by year: ++++---------+++
  runs of the same sign: 3
  a random error would have about 9
```

This is the real diagnosis. If a model has the right shape, its misses scatter at random: the sign changes often, and the runs are about half the number of points. Here there are three runs: the model overshoots systematically, then undershoots for as many years, then overshoots again.

That is what **not noise but a missing part of the model** looks like. A straight line cannot bend, and the series bends.

### Prices multiply — so the line belongs in logarithms

Lesson thirty-eight: percentages do not add, they multiply. A series that is multiplied by roughly the same number every year is not a straight line, it is an exponential. A logarithm turns multiplication into addition, and in logarithms that very same line becomes the right shape:

```text
  a year multiplies prices by: 1.0849 — that is 8.49 % a year
  the average error: 7.39 points
```

The error more than halved, and the coefficient became readable: 8.49 % a year is the same average annual rate as in lesson thirty-eight, only arrived at by a model.

Notice that it is wrong too: 7.39 points on average and 24 in the worst year. A series with 2022 in it is not described by one constant rate — and that is an honest limit of both models rather than a reason to go looking for a third.

### A prediction is two numbers, not one

```text
  the line promises for 2026: 320.4
  the line in logarithms: 351.5
```

Thirty-one points of difference is a choice of the model's shape, not of the data: the series is one and the same. Hence a rule worth learning before any formula: **a prediction without its model and its error is an opinion**, and a digit after the decimal point lends it no weight.

### What this count cannot do

Both models were measured on the very years they learned from. That is only good for seeing the shape of the error; for the question "how will it do on a new year" it is no good at all — the model has already seen the answers. Splitting the series honestly is the next lesson.

## The lesson map

![Lesson map: the points, the line and the error](/static/course/py/map-model-en.svg)

## Say it in your own words

Answer out loud or on paper without looking. The answers are at the end of the lesson.

1. Why does `X` have to become a table with one column?
2. Why does an R² of 0.927 not mean the model is good?
3. What do three runs of signs instead of nine mean?

## Warm-up

Three short steps before the task: predict, fill in, fix. The answers are at the end of the lesson, but answer for yourself first.

**1. Predict.** What does the program print when the points lie on a perfect line?

<!-- drill 1 -->
```python
import numpy as np
from sklearn.linear_model import LinearRegression

x = np.array([[1], [2], [3], [4]])
y = np.array([3.0, 5.0, 7.0, 9.0])
model = LinearRegression().fit(x, y)
print("slope:", round(float(model.coef_[0]), 3))
print("intercept:", round(float(model.intercept_), 3))
print("for 10:", round(float(model.predict(np.array([[10]]))[0]), 3))
```

**2. Fill in the blank.** In place of `...` bring the years into the shape sklearn expects.

```python
# sklearn expects a table of features, not a series
import numpy as np
from sklearn.linear_model import LinearRegression

years = np.array([2020, 2021, 2022, 2023])
index = np.array([100.0, 108.0, 117.0, 126.0])
X = ...
model = LinearRegression().fit(X, index)
print("gained a year:", round(float(model.coef_[0]), 2))
```

**3. Fix it.** The program measures the error as an average difference and declares the model flawless.

```python
# errors add up with their signs, and on a line they cancel each other out
import numpy as np

index = np.array([100.0, 108.0, 117.0, 126.0])
guess = np.array([99.0, 109.0, 118.0, 125.0])
print("the average error:", round(float((index - guess).mean()), 2))
```

## The task

**Required.** Fit a line over the years to the price index and print, year by year: the index, the prediction and the error. Report the slope, the R² and the worst year with the size of the miss. Separately, count the signs of the errors and the number of runs of the same sign — next to how many runs a random error would have. At the end compare three ways of predicting a year on the same years: the line, the line in logarithms, and the rule "same as last year" — and name the best.

The expected output:

<!-- task out -->
```text
== The line over the years
       index   line  error
2010  100.00   73.2   26.8
2011  108.45   88.7   19.8
2012  114.09  104.1   10.0
2013  120.87  119.6    1.3
2014  129.15  135.0   -5.8
2015  137.77  150.5  -12.7
2016  157.56  165.9   -8.3
2017  169.28  181.4  -12.1
2018  179.72  196.8  -17.1
2019  189.30  212.3  -23.0
2020  202.02  227.7  -25.7
2021  218.27  243.2  -24.9
2022  251.07  258.6   -7.5
2023  287.54  274.1   13.4
2024  312.53  289.5   23.0
2025  348.12  305.0   43.1
  it gains a year: 15.45 points
  R²: 0.927
  the worst year: 2025 — off by 43.1 points

== The shape of the error
  the signs: ++++---------+++
  runs: 3 against 9 for a random error

== Compared on the same years
the line                  16.53
the line in logarithms     7.58
same as last year         16.54
  years compared: 15
  the best: the line in logarithms
  the line beats doing nothing by: 0.01 points
```

Done means: the output matches line for line; the comparison runs over the same years rather than over whatever each model happens to have; the error is measured in absolute value; the best way is chosen by the code rather than by eye.

**On your own data.** Take any series of your own over time — monthly spending, weight, kilometres, bills. Fit a line, look at the signs of the errors and answer one question for yourself: do they alternate, or do they come in bands? If they come in bands, a straight line is the wrong shape and a logarithm is worth trying.

**If you want more.**

- Train the line on 2010–2019 only and see what it promises for 2022. That is a rehearsal for the next lesson.
- Add a second feature — the square of the year, say — and see how the error falls and what happens to the shape of the signs.
- Compare `LinearRegression` with `numpy.polyfit(x, y, 1)`: the coefficients have to match, and that is a good way to be sure the model does exactly what you think it does.

## Where this goes in the project

Step twenty-two: the digest gets its first model.

`esep.trend` fits a line through the logarithm of each country's index and prints two numbers: the yearly multiplier that line implies, and the **worst miss** in points of the index. The second stands beside the first for a reason — a trend without its error is a claim without a price, and the countries are described differently: Kazakhstan 1.12 a year with a miss of 3.56, Uzbekistan 1.10 with a miss of 1.02.

The digest prints no forecast. It would take one line, and that is exactly why it is worth not writing: a model deserves a forecast only after it has been checked on years it did not see. Until such a check exists, the digest speaks about the years it was given and stays quiet about the ones it was not.

The project's `requirements.txt` gains scikit-learn with a pinned version: a model whose library version is not written down is a model nobody can reproduce.

## The answers

### To the questions

1. Because sklearn takes a table of "observations × features" and does not guess what you meant by a series of sixteen numbers: one feature with sixteen observations, or sixteen features with one. `reshape(-1, 1)` says it outright.
2. Because R² compares the model with a horizontal line at the mean rather than with anything sensible. On a growing series, anything beats that. Beside the R² of 0.927 here stands an average error of 17 points and a loss to the rule "same as last year".
3. That the error is not random: the model overshoots systematically on one stretch and undershoots on another. A random error would change sign often — the runs would be about half the number of points. Three runs mean the shape of the model is wrong.

### To the warm-up

1. The points lie on the line `y = 2x + 1`, so the slope is exactly 2, the intercept exactly 1, and for ten the model gives 21.

<!-- drill 1 out -->
```text
slope: 2.0
intercept: 1.0
for 10: 21.0
```

2. `years.reshape(-1, 1)`. The minus one means "this dimension works itself out", and the one means "a single column".

<!-- drill 2 -->
```python
# sklearn expects a table of features, not a series: a row per observation,
# a column per feature
import numpy as np
from sklearn.linear_model import LinearRegression

years = np.array([2020, 2021, 2022, 2023])
index = np.array([100.0, 108.0, 117.0, 126.0])
X = years.reshape(-1, 1)
model = LinearRegression().fit(X, index)
print("gained a year:", round(float(model.coef_[0]), 2))
```

<!-- drill 2 out -->
```text
gained a year: 8.7
```

3. The errors have to be taken in absolute value. The average difference of a line is zero by construction — it is not a "good model", it is not a measure of quality at all.

<!-- drill 3 -->
```python
# the error is measured in absolute value: on a line the pluses and the minuses
# cancel each other out
import numpy as np

index = np.array([100.0, 108.0, 117.0, 126.0])
guess = np.array([99.0, 109.0, 118.0, 125.0])
errors = index - guess
print("the mean absolute error:", round(float(np.abs(errors).mean()), 2))
print("and the plain mean:", round(float(errors.mean()), 2))
```

<!-- drill 3 out -->
```text
the mean absolute error: 1.0
and the plain mean: 0.0
```

### To the task

The important number in the task is the last one. A line with an R² of 0.927 beats the rule "repeat last year" by one hundredth of a point, which is to say it does not beat it at all. That is the price worth paying for every model: a comparison not with zero but with the simplest thing that can be done without a model.

The line in logarithms wins for real — 7.58 against 16.5 — and it wins not because it is more complicated but because its shape is right. A model fitted to the shape of the data almost always beats a model fitted to the wish for a number.

## Sources

- [LinearRegression, scikit-learn](https://scikit-learn.org/stable/modules/generated/sklearn.linear_model.LinearRegression.html) — the model itself, its `coef_`, `intercept_` and `predict`.
- [Regression metrics, scikit-learn](https://scikit-learn.org/stable/modules/model_evaluation.html#regression-metrics) — what `r2_score` counts and what `mean_absolute_error` counts.
- [Consumer price index, World Bank](https://data.worldbank.org/indicator/FP.CPI.TOTL) — the series the lesson is built on.
