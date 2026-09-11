"""The exercise of lesson 34: the report on one page.

A title, a date, the number of observations, the picture inside the file, the
table and the source. At the end the page checks itself, the way the chart did
in the last lesson.
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
TITLE = "Inflation in Kazakhstan, 2021-2025"
SOURCE = "Source: World Bank, FP.CPI.TOTL.ZG, taken 2026-09-11"
DAY = date(2026, 9, 11)

table = pd.DataFrame(
    {"year": [2021, 2022, 2023, 2024, 2025],
     "inflation, %": [8.0, 15.0, 14.5, 8.7, 11.4],
     "against last year, pp": [None, 7.0, -0.5, -5.8, 2.7]},
)

fig, ax = plt.subplots(figsize=(7, 3.2))
ax.bar(table["year"].astype(str), table["inflation, %"], color="#b03a2e")
ax.set_ylim(0, 16)
ax.set_title(TITLE)
ax.set_xlabel("year")
ax.set_ylabel("inflation, % a year")
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
<p>Prepared: $day. Observations: $count.</p>
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

print("checking the page:")
print("  data rows:", page.count("<tr>") - 1, "| columns:", page.count("<th>"))
print("  the date is there:", DAY.isoformat() in page)
print("  the size of the data is named:", f"Observations: {len(table)}" in page)
print("  the source is on the page:", html.escape(SOURCE) in page)
print("  the picture is inside the file:", "data:image/png;base64," in page)
print("  a gap is shown as a dash:", "<td>—</td>" in page)
print("  file:", out.name, "| created:", out.exists())
