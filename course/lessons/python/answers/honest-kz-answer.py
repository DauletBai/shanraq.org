"""33-сабақтың тапсырмасының шешімі: өзін өзі тексеретін график.

Сол инфляция қатары. Адал саламыз — нөлден, белгілерімен, дереккөзімен және
нүкте санымен, — содан кейін дайын осьтерден әр ереже орындалды ма деп сұраймыз.
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
# Бағана мәнді ұзындықпен кодтайды, сондықтан оның басы — тек нөл.
ax.set_ylim(0, 16)
ax.set_title(f"Қазақстандағы инфляция, {years[0]}—{years[-1]} (N = {len(kz)})")
ax.set_xlabel("жыл")
ax.set_ylabel("инфляция, жылына %")
ax.grid(True, axis="y", linewidth=0.4, alpha=0.5)
# Дереккөз — суреттің бөлігі: оны тұрған бетінсіз де жіберіп жібереді.
fig.text(0.01, -0.02, SOURCE, fontsize=8, color="#5c5c5c")

out = HERE / "inflyaciya.png"
fig.savefig(out, dpi=150, bbox_inches="tight")

print("графикті тексеру:")
print("  Y осі нөлден:", ax.get_ylim()[0] == 0)
print("  Y осінің белгісі:", repr(ax.get_ylabel()))
print("  белгіде өлшем бірлігі бар:", "%" in ax.get_ylabel())
print("  тақырып:", repr(ax.get_title()))
print("  нүкте саны аталған:", f"N = {len(kz)}" in ax.get_title())
print("  дереккөз суретте:", any(SOURCE in t.get_text() for t in fig.texts))
print("  бағана саны:", len(ax.patches), "| файл:", out.exists())
plt.close(fig)

print()
print("бұл график не дейді:")
print("  бес жылда:", kz[0], "→", kz[-1], "| өзгеріс", round(kz[-1] - kz[0], 1), "п.т.")
print("  ең жоғары:", max(kz), "—", years[kz.index(max(kz))], "| ең төмен:", min(kz), "—", years[kz.index(min(kz))])
