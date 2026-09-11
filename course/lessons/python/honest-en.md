# How not to lie with a chart: the scale, the zero and the labels

_Лид (summary):_ **The thirty-third lesson of the Python course. The same two numbers: from zero the bar is 1.31 times taller, with the axis starting at eight it is 4.86. The same segment of a line: 14 degrees on a wide sheet and 73 on a narrow one. Per cent against percentage points, the choice of window — and a chart that checks itself.**

## Why this matters

The last lesson taught drawing. This one is about the fact that a drawing is almost always more convincing than the numbers deserve.

Lying with a chart is easy, and it is hardly ever done on purpose. The axis did not start at zero because "the difference shows better that way". The last three months were taken because "the rest is out of date". "Up 31 %" was written because that is how it came out. Every step looks reasonable, and together they produce a picture the reader takes something out of that the data does not hold.

So this lesson is not about taste but about counting: every distortion here is measured.

## The whole thing first

The file is `chestno.py`. The same numbers, the distortions computed, and two pictures side by side.

```python
"""Lesson 33: how not to lie with a chart.

The same series drawn honestly and dishonestly. Every distortion is counted
rather than described: how many times taller a bar looks than it is, and how
many degrees steeper a line becomes when the proportions change.
"""

import math

import matplotlib

matplotlib.use("Agg")

import matplotlib.pyplot as plt
from pathlib import Path

HERE = Path(__file__).resolve().parent

# Inflation, % a year. World Bank figures.
years = [2021, 2022, 2023, 2024, 2025]
kz = [8.0, 15.0, 14.5, 8.7, 11.4]

print("== bars: where the axis starts")
low, high = 8.7, 11.4
print("in truth:", high, "against", low, "— by", round(high / low, 2), "times")
for bottom in (0.0, 8.0):
    look = (high - bottom) / (low - bottom)
    print(f"  axis from {bottom:>4}: the bar looks taller by {look:.2f} times")

fig, axes = plt.subplots(1, 2, figsize=(9, 3.5))
for ax, bottom, title in ((axes[0], 0.0, "axis from zero"), (axes[1], 8.0, "axis from eight")):
    ax.bar(["2024", "2025"], [low, high], color=["#7f8c8d", "#b03a2e"])
    ax.set_ylim(bottom, 12)
    ax.set_title(title)
    ax.set_ylabel("inflation, %")
fig.savefig(HERE / "stolbiki.png", dpi=120, bbox_inches="tight")
plt.close(fig)

print()
print("== a line: the same numbers at another angle")
# The angle is computed honestly: from the data into the figure, through the
# size of the axes.
def slope_degrees(width, height, ylow, yhigh):
    x_per_inch = (len(years) - 1) / width
    y_per_inch = (yhigh - ylow) / height
    return math.degrees(math.atan(((kz[-1] - kz[-2]) / y_per_inch) / (1 / x_per_inch)))

for size, ylim in (((8, 3), (0, 16)), ((4, 5), (8, 12))):
    angle = slope_degrees(size[0], size[1], ylim[0], ylim[1])
    print(f"  sheet {size[0]}x{size[1]} inches, axis {ylim}: the last segment {angle:.1f} deg")

print()
print("== per cent and percentage points")
print("it was", low, "%, it is", high, "% — that is a rise of", round(high - low, 1), "percentage points")
print("and \"by", round((high - low) / low * 100, 1), "%\" is another statement about other numbers")

print()
print("== which stretch of the series to show")
print("the whole series:", kz[0], "→", kz[-1], "| change", round(kz[-1] - kz[0], 1), "pp")
print("the last two years:", kz[-2], "→", kz[-1], "| change", round(kz[-1] - kz[-2], 1), "pp")
print("both are true, and the choice is made before anything is drawn")
```

It prints:

```text
== bars: where the axis starts
in truth: 11.4 against 8.7 — by 1.31 times
  axis from  0.0: the bar looks taller by 1.31 times
  axis from  8.0: the bar looks taller by 4.86 times

== a line: the same numbers at another angle
  sheet 8x3 inches, axis (0, 16): the last segment 14.2 deg
  sheet 4x5 inches, axis (8, 12): the last segment 73.5 deg

== per cent and percentage points
it was 8.7 %, it is 11.4 % — that is a rise of 2.7 percentage points
and "by 31.0 %" is another statement about other numbers

== which stretch of the series to show
the whole series: 8.0 → 11.4 | change 3.4 pp
the last two years: 8.7 → 11.4 | change 2.7 pp
both are true, and the choice is made before anything is drawn
```

`stolbiki.png` appears beside it — those two charts: the honest one on the left, the other on the right.

## Going through it

### A bar is measured by its length, so it starts at zero

A bar encodes its value as a **length**. The reader compares not the tops but the whole heights: "this one is twice that one". Start the axis anywhere but zero and the length stops matching the number, and the eye does not know it.

The numbers from the output: 8.7 and 11.4 differ by **1.31 times**. On an axis from zero that is how it looks. On an axis from eight the right bar looks taller by **4.86 times**. The gap between the honest picture and the other one is not "a bit" here, it is nearly fourfold.

The rule is simple and has no exceptions: **the axis of a bar chart starts at zero**. If the difference then does not show, that is the answer: the difference is small. To show it large, show the numbers themselves or the difference — not bars.

### A line is different, but not free

A line encodes its value by the **position** of a point rather than by a length, so zero is not compulsory for it: nobody draws temperature from zero. But a truncated axis still changes the impression, and there is no need to hide it: if the axis does not start at zero, say so.

### The proportions of the sheet are a statement too

The least noticeable distortion. In the output, the same last segment of the series:

- a sheet of 8×3 inches, axis from zero to 16 — a slope of **14.2°**, and the series looks calm;
- a sheet of 4×5 inches, axis from 8 to 12 — a slope of **73.5°**, and the same series looks like a collapse.

Not one number changed. So the proportions are chosen for the data rather than for the impression: a wide sheet for a long series, and the same proportions for charts a reader will compare with each other.

### Per cent and percentage points

Inflation was 8.7 %, it is 11.4 %. There are two correct phrases, and they differ:

- it rose by **2.7 percentage points** — that is the difference;
- the figure itself grew by **31 %** — that is the relative change.

Writing "rose by 31 %" while meaning the first is an ordinary mistake, and it makes the statement eleven times louder. For everything measured in per cent already — inflation, rates, shares — the label must carry "pp" or "%" deliberately.

### Which stretch of the series to show

The whole series: from 8.0 to 11.4, a change of 3.4 pp. The last two years: from 8.7 to 11.4, a change of 2.7 pp. Both are true and they tell different stories, while a third choice — 2022 to 2024 — shows a fall.

The honest rule: **the window is chosen before the result is looked at**, and for an outside reason — "that is how much data there is", "that is how long the contract runs", "year against year". A window picked once it was clear which picture looked better is no longer analysis.

### What belongs on every chart

Five things, without which a picture should not be sent anywhere:

1. **A title** — what is shown.
2. **A Y label with units** — "%", "₸", "items". A number with no units means nothing.
3. **An X label** — what runs across, and over what period.
4. **The source and the date it was taken** — on the picture itself: it will be forwarded without the page it sat on.
5. **How many observations** — `N = 5`. A chart of three points and one of three hundred look alike and are worth different things.

A sixth as the case requires: if the axis does not start at zero, if the series is smoothed, if part of the data was dropped, that is written beside it. A smoothed series, as [the lesson on time](/read/py-uaqyt-resample-rolling-asfreq) showed, belongs next to the real one anyway.

## The map of this lesson

![The map of this lesson: zero, proportions and labels](/static/course/py/map-honest-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why must a bar chart's axis start at zero while a line's need not?
2. How does "rose by 2.7 percentage points" differ from "rose by 31 %"?
3. What changes in a chart if no number is touched and only the proportions of the sheet are?

## The warm-up

Three short steps before the exercise: predict, fill in, fix. The answers are at the end, but answer for yourself first.

**1. Predict.** What does this program print?

<!-- drill 1 -->
```python
low, high = 8.7, 11.4
for bottom in (0.0, 7.0, 8.5):
    look = (high - bottom) / (low - bottom)
    print(f"axis from {bottom}: taller by {look:.2f} (in truth by {high / low:.2f})")
```

**2. Fill in the blank.** In place of `...` set the limits that make the bars honest.

```python
# a bar is measured by its length
import matplotlib

matplotlib.use("Agg")

import matplotlib.pyplot as plt

fig, ax = plt.subplots(figsize=(4, 3))
ax.bar(["2024", "2025"], [8.7, 11.4], color="#b03a2e")
ax.set_ylim(...)
ax.set_ylabel("inflation, %")
print("the bottom of the axis:", ax.get_ylim()[0], "| label:", ax.get_ylabel())
plt.close(fig)
```

**3. Fix it.** The program computes the relative change and calls it percentage points.

```python
# 2.7 and 31 are different numbers and different statements
low, high = 8.7, 11.4
change = (high - low) / low * 100
print(f"inflation rose by {change:.1f} percentage points")
```

## The exercise

**Required.** Draw an honest chart of five years of inflation — bars, from zero, with a title, labels on both axes, the number of observations in the title and the source on the picture itself. Save it to a file. Then **ask the finished axes** whether every rule was kept, and print the check: does the axis start at zero, does the label carry units, is the size of the data named, is the source there. At the end — what this chart says: the change over five years, the maximum and the minimum with their years.

The expected output:

<!-- task out -->
```text
checking the chart:
  the Y axis starts at zero: True
  the Y label: 'inflation, % a year'
  the label carries units: True
  the title: 'Inflation in Kazakhstan, 2021-2025 (N = 5)'
  the number of points is named: True
  the source is on the picture: True
  bars: 5 | файл: True

what this chart says:
  over five years: 8.0 → 11.4 | change 3.4 pp
  the maximum: 15.0 in 2022 | the minimum: 8.0 in 2021
```

Done when: the output matches line for line; the check asks the object itself (`ax.get_ylim()`, `ax.get_ylabel()`, `fig.texts`) rather than repeating what you wrote above; the source is added with `fig.text` rather than in the title; the file is saved after all the drawing.

**On your own data.** Take any chart of yours — last lesson's will do — and walk the five points of the list. Add whatever is missing. Then draw the same series a second time, dishonestly on purpose: truncate the axis, stretch the sheet. Put the two pictures side by side and see how one truth becomes two different messages.

**If you feel like it.**

- Compute the distortion factor for your own truncated axis, with the formula from the first drill.
- Draw one series against another on two Y axes and see that almost any two things can be made to "show a connection".
- Sort the bars by size instead of alphabetically and decide which order is the honest one for your task.

## Where this fits the project

There is no step: the digest draws one picture, and the rules of this lesson are already in it — it is a line, so a zero axis is not required, and the labels and the legend have been there since the last step.

But the lesson settles one debt. The digest's picture carried no **source or date**, so it could not be sent apart from the report: whoever received it did not know where the numbers came from. In the next step, where the report becomes a page, the source appears both on the picture and under the table — one line, taken from the same place as the data.

## The answers

### To the questions

1. Because a bar encodes its value as a length and the eye compares whole lengths. Truncate the axis and the length stops matching the number: 1.31 times becomes 4.86. A line encodes its value by position rather than length, so zero is not compulsory for it — but a truncated axis is labelled all the same.
2. The first is the difference between two quantities that are themselves in per cent: 11.4 − 8.7. The second is the relative change: by how many per cent the quantity itself grew. For inflation those are 2.7 and 31.0 — an elevenfold difference, and confusing them makes the statement louder than the data allows.
3. The slope of the lines, and with it the impression. The same segment is 14.2° on a wide sheet with a zero axis and 73.5° on a narrow one with a truncated axis. Not one number changed, and a calm series became a collapse.

### To the warm-up

1. The higher the axis starts, the stronger the exaggeration: from zero the true 1.31 times, from 8.5 fourteen and a half.

<!-- drill 1 out -->
```text
axis from 0.0: taller by 1.31 (in truth by 1.31)
axis from 7.0: taller by 2.59 (in truth by 1.31)
axis from 8.5: taller by 14.50 (in truth by 1.31)
```

2. `ax.set_ylim(0, 16)`. Zero at the bottom is compulsory; the top is chosen so that the tallest bar does not run into the edge.

<!-- drill 2 -->
```python
import matplotlib

matplotlib.use("Agg")

import matplotlib.pyplot as plt

fig, ax = plt.subplots(figsize=(4, 3))
ax.bar(["2024", "2025"], [8.7, 11.4], color="#b03a2e")
ax.set_ylim(0, 16)
ax.set_ylabel("inflation, %")
print("the bottom of the axis:", ax.get_ylim()[0], "| label:", ax.get_ylabel())
plt.close(fig)
```

<!-- drill 2 out -->
```text
the bottom of the axis: 0.0 | label: inflation, %
```

3. Percentage points are the difference: `high - low`. The relative change is computed separately and called by its own name.

<!-- drill 3 -->
```python
low, high = 8.7, 11.4
points = high - low
percent = (high - low) / low * 100
print(f"inflation rose by {points:.1f} percentage points")
print(f"that is, the figure itself grew by {percent:.1f} %")
```

<!-- drill 3 out -->
```text
inflation rose by 2.7 percentage points
that is, the figure itself grew by 31.0 %
```

### To the exercise

The check asks the chart rather than the programmer. That matters: a label can be set in one place and forgotten in another, while `ax.get_ylabel()` answers for what is actually drawn. Such a check is worth leaving in a report's code for good — it costs four lines and catches a forgotten label before the reader does.

The source goes in through `fig.text` rather than in the title: the title answers "what is shown" and the source answers "where the numbers come from", and those are different questions. On the sheet it lives at the bottom left, small and grey.

## Sources

- [Caveats of data visualisation](https://www.data-to-viz.com/caveats.html) — a collection of the usual distortions with examples.
- [Axes and limits in matplotlib](https://matplotlib.org/stable/api/_as_gen/matplotlib.axes.Axes.set_ylim.html) — `set_ylim` and what it does to a picture.
- [Consumer price inflation, World Bank](https://data.worldbank.org/indicator/FP.CPI.TOTL.ZG) — the source of this lesson's numbers.
