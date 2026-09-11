"""The exercise of lesson 31: a weekly report and honest smoothing.

A weekday series with one jump in it. Counted by weeks, smoothed over five days,
and then shown what the smoothing did to that jump.
"""

import pandas as pd

days, rates = [], []
start = pd.Timestamp("2026-01-05")
for i in range(40):
    day = start + pd.Timedelta(days=i)
    if day.weekday() >= 5:
        continue
    value = 512.0 + i * 0.35 + ((i * 7) % 11 - 5) * 0.6
    if i == 21:
        value += 14.0
    days.append(day)
    rates.append(value)
rate = pd.Series(rates, index=pd.DatetimeIndex(days, name="day"), name="rate").round(2)

print("points in the series:", len(rate), "| calendar days:", len(rate.asfreq("D")),
      "| with no observation:", int(rate.asfreq("D").isna().sum()))

print()
print("by week:")
weekly = rate.resample("W").agg(["count", "mean", "min", "max"]).round(2)
print(weekly.to_string())

print()
smooth = rate.rolling(5).mean().round(2)
print("the five-day rolling mean:")
print("  empty at the start:", int(smooth.isna().sum()))
print("  range of the series:", round(rate.max() - rate.min(), 2),
      "| range of the smoothed one:", round(smooth.max() - smooth.min(), 2))
print("  peak of the series:", rate.idxmax().date(), rate.max(),
      "| peak of the smoothed one:", smooth.idxmax().date(), smooth.max())

print()
jump = rate.diff().idxmax()
print("the largest one-day jump:", jump.date(), "+", round(rate.diff().max(), 2))
print("that day in the smoothed series:", smooth.loc[jump], "— short of the real one by",
      round(rate.loc[jump] - smooth.loc[jump], 2))
