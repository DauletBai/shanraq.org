# The report as a page: a table, a picture and a source in one file

_Лид (summary):_ **The thirty-fourth lesson of the Python course. HTML built out of data: a template instead of glued strings, `html.escape` for everything that came from the table, the picture inside the file — and a page that checks itself. One file that opens in any browser, prints to PDF and travels as a single attachment.**

## Why this matters

The digest already produces three files: a CSV with the numbers, a PNG with the picture and a log on somebody's screen. Whoever receives them has to put them together: open the table in something that reads CSV, look at the picture separately and work out that the two are about the same thing.

A page settles that in one file. It opens in any browser with nothing to install, it prints to PDF from that same browser, it travels as a single attachment and it looks the same for everybody.

No knowledge of HTML is needed for it: five tags will do. What is needed is something else — understanding that data **cannot go into markup as it is**.

## The whole thing first

The file is `otchet.py`. The data deliberately holds a name that breaks markup.

```python
"""Lesson 34: the report as a page.

The table, the picture and the source are collected into one HTML file. Data
reaches the markup only through escaping -- otherwise a single "&" breaks the
page.
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

# The name deliberately holds what breaks markup: angle brackets and an
# ampersand.
rows = [
    ("Kazakhstan", 2024, 8.7),
    ("Uzbekistan", 2024, 9.6),
    ("<b>Russia</b> & Co", 2024, 8.4),
]
table = pd.DataFrame(rows, columns=["country", "year", "inflation"])
SOURCE = "Source: World Bank, FP.CPI.TOTL.ZG"

print("== why data cannot go in as it is")
raw = f"<td>{rows[2][0]}</td>"
safe = f"<td>{html.escape(rows[2][0])}</td>"
print("as it is:", raw)
print("escaped:  ", safe)

# The picture goes straight into the page: one file that can be emailed.
fig, ax = plt.subplots(figsize=(6, 3))
ax.bar(table["country"].str.replace("<b>", "", regex=False).str.replace("</b>", "", regex=False),
       table["inflation"], color="#b03a2e")
ax.set_ylim(0, 12)
ax.set_ylabel("inflation, %")
ax.set_title("Inflation, 2024")
png = HERE / "kartinka.png"
fig.savefig(png, dpi=110, bbox_inches="tight")
plt.close(fig)
picture = base64.b64encode(png.read_bytes()).decode("ascii")

PAGE = Template("""<!doctype html>
<meta charset="utf-8">
<title>$title</title>
<h1>$title</h1>
<p>Prepared: $day. Observations: $count.</p>
<img alt="$title" src="data:image/png;base64,$picture">
<table border="1" cellspacing="0" cellpadding="6">
$head
$body
</table>
<p>$source</p>
""")

head = "<tr>" + "".join(f"<th>{html.escape(name)}</th>" for name in table.columns) + "</tr>"
body = "\n".join(
    "<tr>" + "".join(f"<td>{html.escape(str(cell))}</td>" for cell in row) + "</tr>"
    for row in table.itertuples(index=False)
)

page = PAGE.substitute(
    title=html.escape("Inflation in three countries"),
    day=date(2026, 9, 11).isoformat(),
    count=len(table),
    picture=picture,
    head=head,
    body=body,
    source=html.escape(SOURCE),
)

out = HERE / "otchet.html"
out.write_text(page, encoding="utf-8")

print()
print("== what came out")
print("rows in the table:", page.count("<tr>") - 1, "| columns:", page.count("<th>"))
print("the picture is inside the file:", "data:image/png;base64," in page)
print("the dangerous name became a tag:", "<b>Russia</b>" in page)
print("the reader sees it as text:", "&lt;b&gt;Russia&lt;/b&gt; &amp; Co" in page)
print("file:", out.name, "| created:", out.exists())
```

It prints:

```text
== why data cannot go in as it is
as it is: <td><b>Russia</b> & Co</td>
escaped:   <td>&lt;b&gt;Russia&lt;/b&gt; &amp; Co</td>

== what came out
rows in the table: 3 | columns: 3
the picture is inside the file: True
the dangerous name became a tag: False
the reader sees it as text: True
file: otchet.html | created: True
```

And `otchet.html` appears beside it — a page with the picture inside it.

## Going through it

### Escaping is not over-caution

The first block of the output is the whole point of the lesson:

```text
as it is:  <td><b>Russia</b> & Co</td>
escaped:   <td>&lt;b&gt;Russia&lt;/b&gt; &amp; Co</td>
```

Without escaping the angle brackets out of the data became **markup**: part of the text vanished, part went bold, and an `&` can break the page outright. `html.escape` turns `<`, `>` and `&` into safe sequences, and the reader sees exactly what the data held.

It is the same thought as `params=` in [lesson twenty-one](/read/py-sqlite-keste-kilt-tarih) and `query` in [lesson twenty-seven](/read/py-tandau-suzgi-maska-loc): **data must not become code**. In SQL that is injection; in HTML it is broken layout at best and somebody else's script on your page at worst.

The rule is simple: **everything** that did not come from you is escaped — headings, cells, names, the source. Your own is safe today, and tomorrow the lookup table gains "Zarya & Co LLP".

### A template instead of glue

`string.Template` from the standard library: the text holds `$title`, `$body`, and `substitute` puts the values in. What that gains over `+` and page-long f-strings:

- the markup lies there whole, visible and editable without touching the code;
- a forgotten substitution is an error rather than a quietly empty space;
- a `$` inside the data is harmless: only the named places are filled.

When a template outgrows `Template` — loops, conditions, one piece included in another — people take Jinja2. It works the same way, does more, and has escaping on by default. In this course the standard library is enough.

### A table out of a `DataFrame`

In two lines: the heading from `table.columns`, the body from `itertuples`. pandas has `to_html` ready-made, and it escapes by default — but it brings its own classes and index along, and bending it into the shape you want usually takes longer than writing the two lines yourself.

A gap in a cell is better shown as a dash than as emptiness: an empty cell looks like a value somebody forgot, while a dash says "there is no data" — the distinction from [lesson thirty](/read/py-las-derek-isna-fillna-astype).

### The picture inside the file

`base64.b64encode(png.read_bytes())` puts the picture into the page itself: `src="data:image/png;base64,…"`. The page then stays **one file** — it gets forwarded by email and the picture does not fall off.

The price: the page grows by about a third of the picture's weight (base64 adds a third), and it cannot be cached separately. For a report with one or two pictures that is the right trade; for a page with twenty it is not, and there the pictures stay beside it as files.

### What a page must carry

Everything the chart in [the last lesson](/read/py-grafikpen-otirik-aitpau-shkala) had, plus one:

1. a title — what this is;
2. the date it was prepared — when;
3. the number of observations — how many;
4. the source — where the numbers come from;
5. **the table itself** — so that exact values stand beside the picture.

The picture answers "what does this look like", the table answers "how much exactly", and together they cover both.

## The map of this lesson

![The map of this lesson: a template, escaping and one file](/static/course/py/map-report-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. What happens to a page if the name `<b>Russia</b> & Co` goes into it unescaped?
2. Why put the picture inside the file when it could sit beside it?
3. Why is a gap in a table shown as a dash rather than as an empty cell?

## The warm-up

Three short steps before the exercise: predict, fill in, fix. The answers are at the end, but answer for yourself first.

**1. Predict.** What does this program print?

<!-- drill 1 -->
```python
import html

name = '<b>Russia</b> & Co "south"'
print(html.escape(name))
print(html.escape(name, quote=False))
print("length before:", len(name), "| after:", len(html.escape(name)))
```

**2. Fill in the blank.** In place of `...` make sure the data does not become markup.

```python
# what lands in a cell is whatever came out of the table
import html

cells = ["food", "12 & 13", "<b>phone</b>"]
row = "<tr>" + "".join(f"<td>{...}</td>" for cell in cells) + "</tr>"
print(row)
print("tags in the row:", row.count("<td>"), "| markup from the data:", "<b>" in row)
```

**3. Fix it.** The table row is built without escaping, and the data turned into tags.

```python
# "phone" will arrive in bold and "12 & 13" may break the page
cells = ["food", "12 & 13", "<b>phone</b>"]
row = "<tr>" + "".join(f"<td>{cell}</td>" for cell in cells) + "</tr>"
print(row)
print("tags in the row:", row.count("<td>"), "| markup from the data:", "<b>" in row)
```

## The exercise

**Required.** Collect five years of the report into one page: a title, the date it was prepared, the number of observations, the picture inside the file, a table of three columns (year, inflation, change against last year) and the source at the bottom. Show the gap in the change column as a dash. Save `otchet.html` and print the page's check — going down the same list as in the last lesson.

The expected output:

<!-- task out -->
```text
checking the page:
  data rows: 5 | columns: 3
  the date is there: True
  the size of the data is named: True
  the source is on the page: True
  the picture is inside the file: True
  a gap is shown as a dash: True
  file: otchet.html | created: True
```

Done when: the output matches line for line; every value passes through `html.escape`; the page is built from a template rather than glued together; the picture sits inside the file; the check asks the finished text of the page rather than repeating what you wrote. And above all — open `otchet.html` in a browser and look at it.

**On your own data.** Make a page out of a table of yours and send it to yourself in a messenger. Open it on your phone. Everything unreadable there — a table too wide, a picture too small — is what has to be fixed before anybody else sees the report.

**If you feel like it.**

- Print the page to PDF from the browser and see where it breaks across pages.
- Compare the file size with the picture inside and with the picture beside it.
- Build the same page with `table.to_html()` and decide which you prefer.

## Where this fits the project

Step fourteen: the digest stops handing over three files. `sholu/bet.py` collects the table, the picture and the source into one `report.html` — opened in a browser, printed to PDF, sent as a single attachment.

Three decisions in that file are worth reading in the code: the picture inside the page, the escaping of everything that came from the data, and the source line with its date — the very debt [the last lesson](/read/py-grafikpen-otirik-aitpau-shkala) named: a picture sent on its own did not say where its numbers came from.

Still open. The page is built with `string.Template`, and its appearance lives in the same file as the code. While there is one page that is fine; when there are two, the markup will ask for a file of its own — and that is when Jinja2 starts to make sense.

## The answers

### To the questions

1. The angle brackets become markup: "Russia" arrives in bold and the tags themselves vanish from the text. The browser will try to read `&` as the start of a special character. At best the page looks odd; at worst somebody else's code lands on it. `html.escape` turns those characters into safe sequences.
2. So that the page stays one file. It gets forwarded by email and by messenger, and a picture sitting beside it is lost on the way. The price is size: base64 adds about a third to the picture's weight.
3. An empty cell reads as "somebody forgot to fill it in"; a dash reads as "there is no data". Those are different statements, and the second is the true one.

### To the warm-up

1. `html.escape` replaces `<`, `>` and `&`, and with `quote=True` (the default) the quotation marks as well. The string grows by exactly those replacements.

<!-- drill 1 out -->
```text
&lt;b&gt;Russia&lt;/b&gt; &amp; Co &quot;south&quot;
&lt;b&gt;Russia&lt;/b&gt; &amp; Co "south"
length before: 26 | after: 52
```

2. `html.escape(cell)`. Every value is escaped on its own rather than the finished string as a whole — otherwise your own tags get escaped too.

<!-- drill 2 -->
```python
import html

cells = ["food", "12 & 13", "<b>phone</b>"]
row = "<tr>" + "".join(f"<td>{html.escape(cell)}</td>" for cell in cells) + "</tr>"
print(row)
print("tags in the row:", row.count("<td>"), "| markup from the data:", "<b>" in row)
```

<!-- drill 2 out -->
```text
<tr><td>food</td><td>12 &amp; 13</td><td>&lt;b&gt;phone&lt;/b&gt;</td></tr>
tags in the row: 3 | markup from the data: False
```

3. The same: `html.escape(cell)` inside the f-string. The "markup from the data" check answers `False` — the tags from the cells never reached the page.

<!-- drill 3 -->
```python
import html

cells = ["food", "12 & 13", "<b>phone</b>"]
row = "<tr>" + "".join(f"<td>{html.escape(cell)}</td>" for cell in cells) + "</tr>"
print(row)
print("tags in the row:", row.count("<td>"), "| markup from the data:", "<b>" in row)
```

<!-- drill 3 out -->
```text
<tr><td>food</td><td>12 &amp; 13</td><td>&lt;b&gt;phone&lt;/b&gt;</td></tr>
tags in the row: 3 | markup from the data: False
```

### To the exercise

The check again asks the result rather than the intention: `page.count("<tr>")`, and a search for the date and the source in the page's text. Such a check is worth leaving in a report's code — it catches both a forgotten source and a row of the table lost by accident.

The dash goes into a cell **before** escaping and not through it: it is a character of yours, not data. Which is why the `cell` function tests for a gap first and escapes everything else afterwards.

## Sources

- [The html module](https://docs.python.org/3/library/html.html) — `escape` and exactly what it replaces.
- [string.Template](https://docs.python.org/3/library/string.html#template-strings) — simple templates from the standard library.
- [Jinja2](https://jinja.palletsprojects.com/) — templates for when `Template` gets tight.
