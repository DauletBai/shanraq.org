"""31-сабақтың тапсырмасының шешімі: апталық есеп және адал тегістеу.

Жұмыс күндері бойынша қатар, ішінде бір секіріс бар. Апталап санаймыз, бес күнмен
тегістейміз және тегістеу сол секіріспен не істегенін көрсетеміз.
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

print("қатардағы күн:", len(rate), "| күнтізбелік күн:", len(rate.asfreq("D")),
      "| бақылаусыз:", int(rate.asfreq("D").isna().sum()))

print()
print("апталар бойынша:")
weekly = rate.resample("W").agg(["count", "mean", "min", "max"]).round(2)
print(weekly.to_string())

print()
smooth = rate.rolling(5).mean().round(2)
print("бес күндік жылжымалы орташа:")
print("  басындағы бос:", int(smooth.isna().sum()))
print("  қатардың ауқымы:", round(rate.max() - rate.min(), 2),
      "| тегістелгеннің ауқымы:", round(smooth.max() - smooth.min(), 2))
print("  қатардың шыңы:", rate.idxmax().date(), rate.max(),
      "| тегістелгеннің шыңы:", smooth.idxmax().date(), smooth.max())

print()
jump = rate.diff().idxmax()
print("бір күндегі ең үлкен секіріс:", jump.date(), "+", round(rate.diff().max(), 2))
print("тегістелген қатарда бұл күн:", smooth.loc[jump], "— нағызынан кем:",
      round(rate.loc[jump] - smooth.loc[jump], 2))
