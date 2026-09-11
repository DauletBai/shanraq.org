"""Решение задания урока 33: график, который проверяет сам себя.

Тот же ряд инфляции. Рисуем честно — от нуля, с подписями, источником и числом
точек, — а потом спрашиваем у готовых осей, выполнено ли каждое правило.
"""

import matplotlib

matplotlib.use("Agg")

import matplotlib.pyplot as plt
from pathlib import Path

HERE = Path(__file__).resolve().parent

years = [2021, 2022, 2023, 2024, 2025]
kz = [8.0, 15.0, 14.5, 8.7, 11.4]
SOURCE = "Дереккөз: Дүниежүзілік банк, FP.CPI.TOTL.ZG, 2026-09"

fig, ax = plt.subplots(figsize=(8, 4))
ax.bar([str(year) for year in years], kz, color="#b03a2e")
# Столбик кодирует значение длиной, поэтому его начало — только ноль.
ax.set_ylim(0, 16)
ax.set_title(f"Инфляция в Казахстане, {years[0]}—{years[-1]} (N = {len(kz)})")
ax.set_xlabel("год")
ax.set_ylabel("инфляция, % за год")
ax.grid(True, axis="y", linewidth=0.4, alpha=0.5)
# Источник — часть картинки: её перешлют без страницы, на которой она лежала.
fig.text(0.01, -0.02, SOURCE, fontsize=8, color="#5c5c5c")

out = HERE / "inflyaciya.png"
fig.savefig(out, dpi=150, bbox_inches="tight")

print("проверка графика:")
print("  ось Y от нуля:", ax.get_ylim()[0] == 0)
print("  подпись оси Y:", repr(ax.get_ylabel()))
print("  в подписи есть единицы:", "%" in ax.get_ylabel())
print("  заголовок:", repr(ax.get_title()))
print("  число точек названо:", f"N = {len(kz)}" in ax.get_title())
print("  источник на картинке:", any(SOURCE in t.get_text() for t in fig.texts))
print("  столбиков:", len(ax.patches), "| файл:", out.exists())
plt.close(fig)

print()
print("что говорит этот график:")
print("  за пять лет:", kz[0], "→", kz[-1], "| изменение", round(kz[-1] - kz[0], 1), "п.п.")
print("  максимум:", max(kz), "в", years[kz.index(max(kz))], "| минимум:", min(kz), "в", years[kz.index(min(kz))])
