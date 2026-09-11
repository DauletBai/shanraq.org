# Time in a table: resampling, rolling means and honesty

_Лид (summary):_ **The thirty-first lesson of the Python course. An index of dates can do what a list cannot: `asfreq` shows the days that are not there, `resample` folds days into weeks, `rolling` smooths. And the measured price of smoothing: the range halves, and the day of the real peak disappears from the series.**

## Why this matters

Series over time turn up in almost any job: a rate by day, sales by hour, a meter read by minute. Three things get done to them — they are coarsened, smoothed, and compared with the period before — and all three come ready-made in pandas.

There is a fourth thing that comes in no library: not lying. A smoothed series looks more convincing than the real one, which is exactly why it is so easy to show instead of the real one. By the end of the lesson that is measured in tenge.

## The whole thing first

The file is `vremya.py`. Forty calendar days of a rate, thirty of them working days, and one day with news on it.

```python
"""Lesson 31: time in a table.

Forty days of an exchange rate: weekdays present, weekends absent. What
resampling does with that, what a rolling mean does, and where both begin to
lie.
"""

import pandas as pd

# Built from a counter, the same for everybody. No weekends in it, exactly as
# in a real statement off an exchange.
days, rates = [], []
start = pd.Timestamp("2026-01-05")
for i in range(40):
    day = start + pd.Timedelta(days=i)
    if day.weekday() >= 5:
        continue
    value = 512.0 + i * 0.35 + ((i * 7) % 11 - 5) * 0.6
    if i == 21:
        value += 14.0          # a day of news: one jump in the middle of the series
    days.append(day)
    rates.append(value)

rate = pd.Series(rates, index=pd.DatetimeIndex(days, name="day"), name="rate").round(2)
print(rate.head(3))
print("points:", len(rate), "| from", rate.index.min().date(), "to", rate.index.max().date())

print()
print("== an index of dates can do what a list cannot")
print("January:", len(rate.loc["2026-01"]), "points | one date:", rate.loc["2026-01-19"])
print("weekdays:", sorted(set(rate.index.day_name())))

print()
print("== asfreq('D') shows the holes")
daily = rate.asfreq("D")
print("points now:", len(daily), "| of them empty:", int(daily.isna().sum()))
print(daily.loc["2026-01-09":"2026-01-12"])

print()
print("== resample: a week instead of a day")
weekly = rate.resample("W").agg(["mean", "min", "max", "count"]).round(2)
print(weekly.head(3))

print()
print("== rolling: a mean over five days")
smooth = rate.rolling(5).mean().round(2)
print("the first four are empty:", int(smooth.head(4).isna().sum()))
print(smooth.head(6))

print()
print("== what the smoothing costs")
print("the real range:", round(rate.max() - rate.min(), 2))
print("the smoothed range:", round(smooth.max() - smooth.min(), 2))
print("the day of the peak: real", rate.idxmax().date(), "| smoothed", smooth.idxmax().date())
print("the jump on 26 January: in the series", rate.loc["2026-01-26"], "| in the smoothed one", smooth.loc["2026-01-26"])

print()
print("== the change against the day before")
print((rate.diff().round(2)).head(3))
print("days up:", int((rate.diff() > 0).sum()), "| days down:", int((rate.diff() < 0).sum()))
```

It prints:

```text
day
2026-01-05    509.00
2026-01-06    513.55
2026-01-07    511.50
Name: rate, dtype: float64
points: 30 | from 2026-01-05 to 2026-02-13

== an index of dates can do what a list cannot
January: 20 points | one date: 519.9
weekdays: ['Friday', 'Monday', 'Thursday', 'Tuesday', 'Wednesday']

== asfreq('D') shows the holes
points now: 40 | of them empty: 10
day
2026-01-09    514.00
2026-01-10       NaN
2026-01-11       NaN
2026-01-12    514.45
Freq: D, Name: rate, dtype: float64

== resample: a week instead of a day
              mean    min     max  count
day                                     
2026-01-11  512.82  509.0  516.05      5
2026-01-18  514.31  512.4  516.95      5
2026-01-25  518.44  515.8  520.35      5

== rolling: a mean over five days
the first four are empty: 4
day
2026-01-05       NaN
2026-01-06       NaN
2026-01-07       NaN
2026-01-08       NaN
2026-01-09    512.82
2026-01-12    513.91
Name: rate, dtype: float64

== what the smoothing costs
the real range: 23.75
the smoothed range: 12.73
the day of the peak: real 2026-01-26 | smoothed 2026-02-13
the jump on 26 January: in the series 532.75 | in the smoothed one 521.01

== the change against the day before
day
2026-01-05     NaN
2026-01-06    4.55
2026-01-07   -2.05
Name: rate, dtype: float64
days up: 14 | days down: 15
```

## Going through it

### An index of dates is not labels that look like dates

`pd.DatetimeIndex` is the index of [lesson twenty-five](/read/py-dataframe-series-indeks-qoltanba), except that it understands the labels to be time. Everything else follows:

- `rate.loc["2026-01"]` is the whole of January, though no such label exists in the index: pandas reads it as a stretch;
- `rate.index.day_name()`, `.month`, `.quarter` — every label has parts of its own;
- the series can be asked about its frequency, shifted by a week, laid out by quarter.

With the strings `"2026-01-05"` none of that works: to `resample` they are text, and it answers with a `TypeError`. The first thing done to somebody else's file is `pd.to_datetime` on the date column ([lesson twenty-six](/read/py-oqu-jazu-csv-json-sql), the `parse_dates` argument), and only then `set_index`.

### `asfreq`: the days that are not there

The series has thirty points and the calendar has forty days. Ten of them are weekends with no observation, and in a list that is invisible: a list does not know there is a hole between Friday and Monday.

`rate.asfreq("D")` rebuilds the series onto the calendar's grid: the days with no observation appear, holding `NaN`. That is the first thing worth doing to any series over time — not to fill the holes but **to learn that they are there**. What follows is a decision: weekends are as expected, while a Tuesday missing mid-week is a reason to ask the source what happened.

### `resample`: a week instead of a day

`resample("W")` cuts the series into weeks and computes whatever you asked for in each. It is the `groupby` of [lesson twenty-eight](/read/py-toptau-groupby-agg-transform), with a calendar stretch for a key and nothing to build by hand.

Frequencies are written short: `D` a day, `W` a week, `ME` a month end, `QE` a quarter, `YE` a year, `h` an hour. A week's label is by default its **last** day (`2026-01-11` in the output is the week of 5 to 11 January). That is worth checking every time: a weekly report labelled by its last day reads easily as "as of 11 January" when it means "for the week to 11 January".

The `count` in the aggregate is not decoration. It shows how many points fell into the week: with five in one and two in another, the two averages are not comparable.

### `rolling`: the moving average

`rate.rolling(5).mean()` is the mean of the last five points, stepping along the series. The first four values are empty, and rightly: five points have not yet accumulated, and computing over three while calling it a five-point mean is already untrue. If the empty start is in the way there is `min_periods=3`, but then the first values are computed over a different number of points, and that has to be remembered.

The window looks **backwards** by default: the value on 9 January is the mean of 5 to 9 January. So a smoothed series lags behind the real one. `center=True` puts the window symmetrically, and then it does not lag — but the empty values at the start and the end are twice as many.

### What the smoothing costs

This is what the lesson exists for. The series has one day with news on it — 26 January, when the rate jumped by 14 tenge. What became of it:

| | the real series | smoothed |
|---|---|---|
| range | 23.75 | 12.73 |
| day of the peak | 26 January | 13 February |
| the value on 26 January | 532.75 | 521.01 |

The smoothing did exactly what it was asked to: it removed the outliers. But the outlier here is not noise, it is an **event**. On the smoothed chart 26 January stands out in no way at all, the range is half the real one, and the peak has moved two and a half weeks forward — to where the series merely runs higher on average.

Three rules follow, and they are worth taking as they are:

1. A smoothed series is a **conclusion**, not data. The real one belongs beside it — at the very least as a pale line on the same chart.
2. The range, the maximum and the minimum are computed on the real series, not the smoothed one. Otherwise a report ends up saying the rate never passed 526 when it stood at 533.
3. Smoothing has a parameter — the width of the window — and it is chosen before the result is looked at. Otherwise the window is picked to give the curve the shape that was wanted, and that is no longer analysis.

### Comparing with the period before

`diff()` is the difference from the previous point, `pct_change()` the same as a fraction, `shift(1)` the series itself moved along by one. All three work on neighbours **in the order the series happens to lie in**, so it is sorted by its index first; an unsorted series gives numbers that look plausible and mean nothing.

For a table with several series in it — as in the project, where there are three countries — `diff` is taken inside the group: `table.groupby("country")["value"].diff()`. Without the grouping, the first row of each next country subtracts the last year of the previous one.

## The map of this lesson

![The map of this lesson: days, weeks and smoothing](/static/course/py/map-time-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why run `asfreq("D")` if you have no intention of filling the gaps?
2. Why are the first four values of a five-point rolling mean empty?
3. What does a report lose when it shows only the smoothed series?

## The warm-up

Three short steps before the exercise: predict, fill in, fix. The answers are at the end, but answer for yourself first.

**1. Predict.** What does this program print?

<!-- drill 1 -->
```python
import pandas as pd

days = pd.to_datetime(["2026-01-05", "2026-01-06", "2026-01-07", "2026-01-12", "2026-01-13"])
rate = pd.Series([512.0, 514.0, 513.0, 520.0, 522.0], index=days)
print(rate.resample("W").mean().round(2))
print("weeks:", len(rate.resample("W").mean()))
```

**2. Fill in the blank.** In place of `...` put a window width that leaves exactly two empty values at the start.

```python
# as many points in the window, so one fewer empty values at the start
import pandas as pd

days = pd.date_range("2026-01-05", periods=6, freq="D")
rate = pd.Series([512.0, 514.0, 513.0, 520.0, 522.0, 519.0], index=days)
smooth = rate.rolling(...).mean().round(2)
print(smooth.tolist())
print("empty at the start:", int(smooth.isna().sum()))
```

**3. Fix it.** The program raises `Only valid with DatetimeIndex, TimedeltaIndex or PeriodIndex`. There are dates in the index, but they are strings.

```python
# the index looks like dates and pandas sees text
import pandas as pd

rate = pd.Series([512.0, 514.0, 520.0], index=["2026-01-05", "2026-01-06", "2026-01-12"])
print(rate.resample("W").mean().round(2).tolist())
```

## The exercise

**Required.** Take the same series from the lesson — forty calendar days, thirty working ones, the jump on 26 January. Build a weekly report (`count`, `mean`, `min`, `max`) and a five-day rolling mean. Print: how many calendar days were left without an observation; the weekly table; the range of the real series and of the smoothed one; the day and the size of the peak in both; the largest one-day jump and what it became in the smoothed series.

The expected output:

<!-- task out -->
```text
points in the series: 30 | calendar days: 40 | with no observation: 10

by week:
            count    mean     min     max
day                                      
2026-01-11      5  512.82  509.00  516.05
2026-01-18      5  514.31  512.40  516.95
2026-01-25      5  518.44  515.80  520.35
2026-02-01      5  522.73  516.70  532.75
2026-02-08      5  522.74  520.10  524.65
2026-02-15      5  525.55  523.05  528.05

the five-day rolling mean:
  empty at the start: 4
  range of the series: 23.75 | range of the smoothed one: 12.73
  peak of the series: 2026-01-26 532.75 | peak of the smoothed one: 2026-02-13 525.55

the largest one-day jump: 2026-01-26 + 14.45
that day in the smoothed series: 521.01 — short of the real one by 11.74
```

Done when: the output matches line for line; the missing days are found with `asfreq` rather than by subtracting lengths; the weeks come from `resample` rather than from grouping on a week number pulled out of a string; the range and the peak of the real series are computed on the real series; the difference on the day of the jump is computed rather than eyeballed.

**On your own data.** Take a series of your own over time — daily spending, steps, weight, anything. Run `asfreq` and see how many days are really missing from it. Then build rolling means over two different windows and tell yourself honestly which one you would put in a report, and why.

**If you feel like it.**

- Compare `rolling(5).mean()` with `rolling(5, center=True).mean()`: where the series lags, and where it is trimmed at both ends.
- Run `resample("ME")` and see what label pandas gives a month.
- Compute `pct_change()` and find the day that grew most in per cent — see whether it is the same day that grew most in tenge.

## Where this fits the project

Step twelve, and there is less in it than one would like — for an honest reason.

The digest's year stops being a number and becomes a date: `pd.to_datetime(year, format="%Y")`. A number sorts and subtracts by luck; a date does so by meaning, and on the day the digest reads a monthly series nothing above that line will have to change.

The report gains the question a yearly series can actually answer: how much the last year differs from the one before. That is `diff()` inside the grouping — one line, and none of the off-by-one a loop over sorted years invites.

What did **not** go into the project is `resample` and `rolling`, and the step says so: five yearly points are not a series to smooth. A five-point rolling mean would be a picture with nothing behind it. They arrive when the digest starts reading something denser than one figure a year.

Still open. The digest's index is a date now, but the series has no declared frequency: whether it is yearly or monthly the program does not know — it finds out from what arrived. That is fine while there is one source.

## The answers

### To the questions

1. To see what is not there. A missing day shows up in no way in a list — the series is simply shorter than the calendar. `asfreq("D")` puts those days back with `NaN`, and then it is a decision: weekends are expected, a missing Tuesday is a reason to ask the source what happened.
2. Because five points have not yet accumulated. Computing over three and calling it a five-point mean is untrue; the `NaN` here means "no answer yet", which is an honest answer.
3. It loses events. Smoothing removes outliers, and an outlier is often the most important thing that happened: a day of news, a failure, a rush. In the example the range falls from 23.75 to 12.73 and the day of the peak moves two and a half weeks. Which is why a smoothed series is shown beside the real one rather than instead of it.

### To the warm-up

1. `resample("W")` counts by calendar weeks and labels each with its last day, a Sunday. Two weeks: 5 to 11 January and 12 to 18 January.

<!-- drill 1 out -->
```text
2026-01-11    513.0
2026-01-18    521.0
Freq: W-SUN, dtype: float64
weeks: 2
```

2. `rolling(3)`. The empty values at the start are always one fewer than the points in the window: a window of three leaves two.

<!-- drill 2 -->
```python
import pandas as pd

days = pd.date_range("2026-01-05", periods=6, freq="D")
rate = pd.Series([512.0, 514.0, 513.0, 520.0, 522.0, 519.0], index=days)
smooth = rate.rolling(3).mean().round(2)
print(smooth.tolist())
print("empty at the start:", int(smooth.isna().sum()))
```

<!-- drill 2 out -->
```text
[nan, nan, 513.0, 515.67, 518.33, 520.33]
empty at the start: 2
```

3. `rate.index = pd.to_datetime(rate.index)`. Strings that look like dates stay strings until they are turned into dates; `resample` works only with a time index.

<!-- drill 3 -->
```python
import pandas as pd

rate = pd.Series([512.0, 514.0, 520.0], index=["2026-01-05", "2026-01-06", "2026-01-12"])
rate.index = pd.to_datetime(rate.index)
print(rate.resample("W").mean().round(2).tolist())
```

<!-- drill 3 out -->
```text
[513.0, 520.0]
```

### To the exercise

The missing days are counted with `asfreq("D")` rather than by subtracting `40 − 30`: the length of the calendar is known here only because we built it ourselves, and with a real series it would have to come from somewhere.

The range and the peak of the real series are computed **before** the smoothing and on the real series. That is the main thing the exercise checks: the moment a report carries the `max` and `min` of a smoothed series, it begins asserting something that never happened.

## Sources

- [Time series and date functionality](https://pandas.pydata.org/docs/user_guide/timeseries.html) — the date index, frequencies, `resample` and shifts.
- [Windowing operations](https://pandas.pydata.org/docs/user_guide/window.html) — `rolling`, `expanding`, `min_periods` and centring the window.
- [Offset aliases](https://pandas.pydata.org/docs/user_guide/timeseries.html#offset-aliases) — `D`, `W`, `ME`, `QE`, `YE` and the rest.
