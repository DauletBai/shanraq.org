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
<h2>$index_title</h2>
<table>
$index_head
$index_body
</table>
<h2>$money_title</h2>
<table>
$money_head
$money_body
</table>
<p class="source">$source</p>
<p class="source">$money_source</p>
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
    "индекс": pishim.number,
    "есе": pishim.number,
    "мың теңге": pishim.number,
    "% ЖІӨ": pishim.number,
    "өзгерді": pishim.signed,
}


def _cell(name, value):
    """One cell: formatted if it is a number, escaped either way.

    A missing value has several spellings by the time it reaches here -- nan
    from a float column, None from a plain object, <NA> from a nullable
    integer -- and all of them mean the same thing to a reader: a dash.
    """
    shape = FORMATS.get(name)
    text = shape(value) if shape else ("" if value is None else str(value))
    if text in ("", "nan", "None", "<NA>"):
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


def write(report, picture, path, title, source, spread=None, index=None,
          money=None, money_source="", day=None):
    """Writes the report, its picture and the three summary tables into one file."""
    assert not report.empty, "есепте бірде-бір жол жоқ"
    path.parent.mkdir(parents=True, exist_ok=True)

    head, body = _table(report)
    # The second block is the whole series, not its last year: the mean, the
    # median and the distance between the quarters. Without it the page says
    # what happened once and stays silent about how usual that is.
    spread_head, spread_body = _table(spread) if spread is not None else ("", "")
    # The third block is the whole run of years as one number: percentages are
    # multipliers, so what a reader wants is "how many times", not a sum.
    index_head, index_body = _table(index) if index is not None else ("", "")
    # The fourth block comes from a different indicator altogether: how much
    # money a country has for a tenge of its yearly output. It carries its own
    # source line, because a page with two sources and one credit is a page that
    # says something it cannot back.
    money_head, money_body = _table(money) if money is not None else ("", "")

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
        index_title="Баға индексі: базадан бері",
        index_head=index_head,
        index_body=index_body,
        money_title="Ақша: ЖІӨ теңгесіне қанша",
        money_head=money_head,
        money_body=money_body,
        source=html.escape(source),
        money_source=html.escape(money_source),
    )
    path.write_text(page, encoding="utf-8")
    return path
