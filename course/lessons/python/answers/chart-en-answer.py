"""The exercise of lesson 32: two files out of one chart.

The same series, two lines, the day of the news marked -- saved as a PNG for an
email and an SVG for a page. Not one window: the program has to work on a
schedule.
"""

import matplotlib

matplotlib.use("Agg")

import matplotlib.pyplot as plt
import pandas as pd
from pathlib import Path

HERE = Path(__file__).resolve().parent

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
smooth = rate.rolling(5).mean()

fig, ax = plt.subplots(figsize=(9, 4))
ax.plot(rate.index, rate.values, color="#b03a2e", linewidth=1.6, label="rate")
ax.plot(rate.index, smooth, color="#7f8c8d", linewidth=1.2, linestyle="--",
        label="rolling mean, 5 days")

# The day of the news is marked with a dot and a label: the smoothed line does
# not show it.
peak_day = rate.idxmax()
ax.plot([peak_day], [rate.max()], marker="o", color="#b03a2e")
ax.annotate(f"{peak_day.date()}: {rate.max()}", xy=(peak_day, rate.max()),
            xytext=(6, 6), textcoords="offset points")

ax.set_title("The rate: the real series and the smoothed one")
ax.set_xlabel("day")
ax.set_ylabel("tenge per dollar")
ax.grid(True, linewidth=0.4, alpha=0.5)
ax.legend(loc="upper left")
fig.autofmt_xdate()

files = []
for name in ("kurs.png", "kurs.svg"):
    path = HERE / name
    fig.savefig(path, dpi=150, bbox_inches="tight")
    files.append(path)
plt.close(fig)

print("points:", len(rate), "| lines on the axes:", len(ax.get_lines()))
print("title:", ax.get_title())
print("marked day:", peak_day.date(), "| value:", rate.max())
for path in files:
    print(f"{path.name}: created {path.exists()}, not empty {path.stat().st_size > 0}")
