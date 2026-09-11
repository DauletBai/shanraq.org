"""Решение задания урока 32: два файла из одного графика.

Тот же ряд, две линии, отмеченный день новостей — и сохранение в PNG для письма
и в SVG для страницы. Ни одного окна: программа должна работать по расписанию.
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
ax.plot(rate.index, rate.values, color="#b03a2e", linewidth=1.6, label="курс")
ax.plot(rate.index, smooth, color="#7f8c8d", linewidth=1.2, linestyle="--",
        label="скользящее среднее, 5 дней")

# День новостей отмечен точкой и подписью: на сглаженной линии его не видно.
peak_day = rate.idxmax()
ax.plot([peak_day], [rate.max()], marker="o", color="#b03a2e")
ax.annotate(f"{peak_day.date()}: {rate.max()}", xy=(peak_day, rate.max()),
            xytext=(6, 6), textcoords="offset points")

ax.set_title("Курс доллара: настоящий ряд и сглаженный")
ax.set_xlabel("день")
ax.set_ylabel("тенге за доллар")
ax.grid(True, linewidth=0.4, alpha=0.5)
ax.legend(loc="upper left")
fig.autofmt_xdate()

files = []
for name in ("kurs.png", "kurs.svg"):
    path = HERE / name
    fig.savefig(path, dpi=150, bbox_inches="tight")
    files.append(path)
plt.close(fig)

print("точек:", len(rate), "| линий на осях:", len(ax.get_lines()))
print("заголовок:", ax.get_title())
print("отмечен день:", peak_day.date(), "| значение:", rate.max())
for path in files:
    print(f"{path.name}: создан {path.exists()}, не пустой {path.stat().st_size > 0}")
