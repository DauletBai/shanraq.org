"""The exercise of lesson 33: a chart that checks itself.

The same inflation series. Drawn honestly -- from zero, with labels, a source
and the number of points -- and then asked, of the finished axes, whether every
rule was kept.
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
# A bar encodes its value as a length, so it can only begin at zero.
ax.set_ylim(0, 16)
ax.set_title(f"Inflation in Kazakhstan, {years[0]}-{years[-1]} (N = {len(kz)})")
ax.set_xlabel("year")
ax.set_ylabel("inflation, % a year")
ax.grid(True, axis="y", linewidth=0.4, alpha=0.5)
# The source belongs to the picture: it will be forwarded without the page it
# was sitting on.
fig.text(0.01, -0.02, SOURCE, fontsize=8, color="#5c5c5c")

out = HERE / "inflyaciya.png"
fig.savefig(out, dpi=150, bbox_inches="tight")

print("checking the chart:")
print("  the Y axis starts at zero:", ax.get_ylim()[0] == 0)
print("  the Y label:", repr(ax.get_ylabel()))
print("  the label carries units:", "%" in ax.get_ylabel())
print("  the title:", repr(ax.get_title()))
print("  the number of points is named:", f"N = {len(kz)}" in ax.get_title())
print("  the source is on the picture:", any(SOURCE in t.get_text() for t in fig.texts))
print("  bars:", len(ax.patches), "| файл:", out.exists())
plt.close(fig)

print()
print("what this chart says:")
print("  over five years:", kz[0], "→", kz[-1], "| change", round(kz[-1] - kz[0], 1), "pp")
print("  the maximum:", max(kz), "in", years[kz.index(max(kz))], "| the minimum:", min(kz), "in", years[kz.index(min(kz))])
