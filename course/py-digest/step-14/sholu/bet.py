"""The page: the report as one file anybody can open.

A CSV needs a spreadsheet and a PNG carries no numbers. A page carries both, it
opens in any browser, it prints to PDF, and it can be sent as a single file --
which is why the picture goes inside it rather than beside it.

Everything that came from the data is escaped on the way into the markup. Our
own country names are harmless today; the rule is not about today.
"""

import base64
import html
from datetime import date
from string import Template

PAGE = Template("""<!doctype html>
<meta charset="utf-8">
<title>$title</title>
<style>
 body { font: 15px/1.5 system-ui, sans-serif; max-width: 760px; margin: 32px auto; padding: 0 16px; }
 table { border-collapse: collapse; width: 100%; }
 th, td { border: 1px solid #d6d3ce; padding: 6px 10px; text-align: right; }
 th:first-child, td:first-child { text-align: left; }
 img { max-width: 100%; }
 .source { color: #5c5c5c; font-size: 13px; }
</style>
<h1>$title</h1>
<p>$prepared: $day. $rows_word: $count.</p>
<img alt="$title" src="data:image/png;base64,$picture">
<table>
$head
$body
</table>
<p class="source">$source</p>
""")


def _cell(value):
    """One cell: a gap becomes a dash, everything else is escaped."""
    text = "" if value is None else str(value)
    if text in ("", "nan", "None"):
        return "<td>—</td>"
    return f"<td>{html.escape(text)}</td>"


def write(report, picture, path, title, source, day=None):
    """Writes the report table and its picture into a single HTML file."""
    assert not report.empty, "есепте бірде-бір жол жоқ"
    path.parent.mkdir(parents=True, exist_ok=True)

    head = "<tr>" + "".join(f"<th>{html.escape(str(name))}</th>" for name in report.columns) + "</tr>"
    body = "\n".join("<tr>" + "".join(_cell(value) for value in row) + "</tr>"
                     for row in report.itertuples(index=False))

    page = PAGE.substitute(
        title=html.escape(title),
        prepared="Дайындалды",
        rows_word="Жол саны",
        day=(day or date.today()).isoformat(),
        count=len(report),
        picture=base64.b64encode(picture.read_bytes()).decode("ascii"),
        head=head,
        body=body,
        source=html.escape(source),
    )
    path.write_text(page, encoding="utf-8")
    return path
