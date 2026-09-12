"""The page: the report as one file anybody can open.

A CSV needs a spreadsheet and a PNG carries no numbers. A page carries both, it
opens in any browser, it prints to PDF, and it can be sent as a single file --
which is why the picture goes inside it rather than beside it.

Everything that came from the data is escaped on the way into the markup, and
every number goes through sholu/pishim.py: the reader sees "11,54" rather than
11.539999999999999, and the columns line up because the separators agree.
"""

import base64
import html
from datetime import date
from string import Template

from sholu import pishim

PAGE = Template("""<!doctype html>
<meta charset="utf-8">
<title>$title</title>
<style>
 body { font: 15px/1.5 system-ui, sans-serif; max-width: 760px; margin: 32px auto; padding: 0 16px; }
 table { border-collapse: collapse; width: 100%; }
 th, td { border: 1px solid #d6d3ce; padding: 6px 10px; text-align: right; }
 th:first-child, td:first-child { text-align: left; }
 td { font-variant-numeric: tabular-nums; }
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
<h2>$spread_title</h2>
<table>
$spread_head
$spread_body
</table>
<p class="source">$source</p>
""")

# Which formatter each column of the report needs. A column not named here is
# text and goes through escaping only.
FORMATS = {
    "орташа": pishim.number,
    "ең_жоғары": pishim.number,
    "соңғы_мән": pishim.number,
    "өткен_жылға": pishim.signed,
    "медиана": pishim.number,
    "шашырау": pishim.number,
    "ширекаралық": pishim.number,
    "мәні": pishim.number,
}


def _cell(name, value):
    """One cell: formatted if it is a number, escaped either way."""
    shape = FORMATS.get(name)
    text = shape(value) if shape else ("" if value is None else str(value))
    if text in ("", "nan", "None"):
        text = pishim.DASH
    return f"<td>{html.escape(text)}</td>"


def _table(frame):
    """A head row and a body, escaped and formatted the same way everywhere."""
    head = "<tr>" + "".join(f"<th>{html.escape(str(name))}</th>" for name in frame.columns) + "</tr>"
    body = "\n".join(
        "<tr>" + "".join(_cell(name, value) for name, value in zip(frame.columns, row)) + "</tr>"
        for row in frame.itertuples(index=False)
    )
    return head, body


def write(report, picture, path, title, source, spread=None, day=None):
    """Writes the report, its picture and the spread table into one HTML file."""
    assert not report.empty, "есепте бірде-бір жол жоқ"
    path.parent.mkdir(parents=True, exist_ok=True)

    head, body = _table(report)
    # The second block is the whole series, not its last year: the mean, the
    # median and the distance between the quarters. Without it the page says
    # what happened once and stays silent about how usual that is.
    spread_head, spread_body = _table(spread) if spread is not None else ("", "")

    page = PAGE.substitute(
        title=html.escape(title),
        prepared="Дайындалды",
        rows_word="Жол саны",
        day=(day or date.today()).isoformat(),
        count=len(report),
        picture=base64.b64encode(picture.read_bytes()).decode("ascii"),
        head=head,
        body=body,
        spread_title="Барлық жылдар бойынша",
        spread_head=spread_head,
        spread_body=spread_body,
        source=html.escape(source),
    )
    path.write_text(page, encoding="utf-8")
    return path
