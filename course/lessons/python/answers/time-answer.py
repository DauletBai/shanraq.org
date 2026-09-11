"""Решение задания урока 31: недельный отчёт и честное сглаживание.

Ряд по будням, в нём один скачок. Считаем неделями, сглаживаем по пяти дням и
показываем, что сглаживание с этим скачком сделало.
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

print("дней в ряду:", len(rate), "| календарных дней:", len(rate.asfreq("D")),
      "| без наблюдения:", int(rate.asfreq("D").isna().sum()))

print()
print("по неделям:")
weekly = rate.resample("W").agg(["count", "mean", "min", "max"]).round(2)
print(weekly.to_string())

print()
smooth = rate.rolling(5).mean().round(2)
print("скользящее среднее по пяти дням:")
print("  пустых в начале:", int(smooth.isna().sum()))
print("  размах ряда:", round(rate.max() - rate.min(), 2),
      "| размах сглаженного:", round(smooth.max() - smooth.min(), 2))
print("  пик ряда:", rate.idxmax().date(), rate.max(),
      "| пик сглаженного:", smooth.idxmax().date(), smooth.max())

print()
jump = rate.diff().idxmax()
print("самый большой скачок за день:", jump.date(), "+", round(rate.diff().max(), 2))
print("в сглаженном ряду этот день:", smooth.loc[jump], "— меньше настоящего на",
      round(rate.loc[jump] - smooth.loc[jump], 2))
