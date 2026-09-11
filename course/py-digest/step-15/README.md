# Step 15 — the numbers are formatted where the reader sees them

The state of the digest after lesson 35, *Numbers in a report*.

Until now the page printed whatever `str()` gave it: `11.54` in one cell and
`8.690822` in another, with a full stop where this country writes a comma.

`sholu/pishim.py` is the one place where a number turns into text. Three
functions and one rule: everything above this file keeps full precision, and
everything the reader sees passes through here.

- `number` — a space in the thousands and a comma in the decimals;
- `signed` — a change carries its sign, and a rounded zero is not "+0";
- `percent` — a share with the same comma, because a table whose separators
  disagree reads as two tables.

A gap becomes a dash in all three, so an empty cell never has to be guessed at.

The page also asks for `font-variant-numeric: tabular-nums`, which makes the
digits the same width — otherwise a column of numbers does not line up however
carefully it was formatted.

```
pip install -r requirements.txt
python3 main.py
open data/report.html
```

`data/` is what the program produces and is not kept in the repository.
