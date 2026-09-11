"""34-сабақтың тапсырмасының шешімі: бір беттегі есеп.

Тақырып, күн, бақылау саны, файл ішіндегі сурет, кесте және дереккөз. Соңында
бет өзін өзі тексереді — өткен сабақтағы график сияқты.
"""

import base64
import html
from datetime import date
from pathlib import Path
from string import Template

import matplotlib

matplotlib.use("Agg")

import matplotlib.pyplot as plt
import pandas as pd

HERE = Path(__file__).resolve().parent
TITLE = "Қазақстандағы инфляция, 2021—2025"
SOURCE = "Дереккөз: Дүниежүзілік банк, FP.CPI.TOTL.ZG, 2026-09-11"
DAY = date(2026, 9, 11)

table = pd.DataFrame(
    {"жыл": [2021, 2022, 2023, 2024, 2025],
     "инфляция, %": [8.0, 15.0, 14.5, 8.7, 11.4],
     "өткен жылға, п.т.": [None, 7.0, -0.5, -5.8, 2.7]},
)

fig, ax = plt.subplots(figsize=(7, 3.2))
ax.bar(table["жыл"].astype(str), table["инфляция, %"], color="#b03a2e")
ax.set_ylim(0, 16)
ax.set_title(TITLE)
ax.set_xlabel("жыл")
ax.set_ylabel("инфляция, жылына %")
ax.grid(True, axis="y", linewidth=0.4, alpha=0.5)
png = HERE / "otchet.png"
fig.savefig(png, dpi=120, bbox_inches="tight")
plt.close(fig)

def cell(value):
    """Ячейка таблицы: пропуск — прочерком, всё остальное — экранированным."""
    if value is None or (isinstance(value, float) and pd.isna(value)):
        return "<td>—</td>"
    return f"<td>{html.escape(str(value))}</td>"

head = "<tr>" + "".join(f"<th>{html.escape(name)}</th>" for name in table.columns) + "</tr>"
body = "\n".join("<tr>" + "".join(cell(v) for v in row) + "</tr>"
                 for row in table.itertuples(index=False))

PAGE = Template("""<!doctype html>
<meta charset="utf-8">
<title>$title</title>
<style>
 body { font: 15px/1.5 system-ui, sans-serif; max-width: 720px; margin: 32px auto; padding: 0 16px; }
 table { border-collapse: collapse; width: 100%; }
 th, td { border: 1px solid #d6d3ce; padding: 6px 10px; text-align: right; }
 th:first-child, td:first-child { text-align: left; }
 img { max-width: 100%; }
 .source { color: #5c5c5c; font-size: 13px; }
</style>
<h1>$title</h1>
<p>Дайындалды: $day. Бақылау саны: $count.</p>
<img alt="$title" src="data:image/png;base64,$picture">
<table>
$head
$body
</table>
<p class="source">$source</p>
""")

page = PAGE.substitute(
    title=html.escape(TITLE),
    day=DAY.isoformat(),
    count=len(table),
    picture=base64.b64encode(png.read_bytes()).decode("ascii"),
    head=head,
    body=body,
    source=html.escape(SOURCE),
)
out = HERE / "otchet.html"
out.write_text(page, encoding="utf-8")

print("бетті тексеру:")
print("  дерек жолы:", page.count("<tr>") - 1, "| баған:", page.count("<th>"))
print("  күні көрсетілген:", DAY.isoformat() in page)
print("  дерек көлемі аталған:", f"Бақылау саны: {len(table)}" in page)
print("  дереккөз бетте:", html.escape(SOURCE) in page)
print("  сурет файл ішінде:", "data:image/png;base64," in page)
print("  жетіспейтін мән сызықшамен:", "<td>—</td>" in page)
print("  файл:", out.name, "| жасалды:", out.exists())
