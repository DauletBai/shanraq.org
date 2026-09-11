# Step 14 — the report becomes a page

The state of the digest after lesson 34, *The report as a page*.

Until now the digest produced three files and whoever received them had to put
them together: a CSV that needs a spreadsheet, a PNG with no numbers on it, and
a log on somebody's screen. `sholu/bet.py` collects the table, the picture and
the source into one HTML file.

Three decisions in it are worth reading:

- **the picture goes inside the page** as a base64 data URI. The page is then one
  file: it can be attached to an email and it will not arrive with a broken
  image;
- **everything from the data is escaped** on the way into the markup. Our own
  country names hold nothing dangerous today; the rule is not about today;
- **the source and the date are on the page**, which settles what lesson 33
  called a debt: a picture sent by itself said nothing about where its numbers
  came from.

```
pip install -r requirements.txt
python3 main.py        # first run: fetches, saves, draws, collects
python3 main.py        # second run: reads from disk
open data/report.html  # or double-click it
```

`data/` is what the program produces and is not kept in the repository.
