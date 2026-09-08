# Step 12 — the store becomes a table

The state of the digest after lesson 12, *CSV*.

The saved series is a real CSV now: written with `csv.DictWriter`, read back
with `csv.DictReader` through the names of its columns, and opened as
`utf-8-sig` so that a file touched by Excel still parses. The report is a CSV
too — a spreadsheet can open it without anyone converting anything.

```
python3 main.py        # first run: fetches and saves
python3 main.py        # second run: reads from disk
```

`data/` is what the program produces and is not kept in the repository.
