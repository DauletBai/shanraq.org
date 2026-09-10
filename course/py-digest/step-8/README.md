# Step 8 — the disk speaks tables

The state of the digest after lesson 26, *Reading and writing: CSV, JSON, SQL*.

Step 7 held the series as a `DataFrame` in memory and still wrote it to disk row
by row with the `csv` module. Here the store speaks tables both ways:

- `saqtau.save(table, path)` is `to_csv` with `float_format` and an empty cell
  for a missing year;
- `saqtau.load(path)` is `read_csv` with `index_col="year"`, so the year comes
  back as a label rather than as a column;
- `esep.write` writes the report with `to_csv(index=False)` — the row numbers
  are ours and have no business in the file.

The `note` column is gone. It existed to spell "no data" in words next to an
`n/a` marker; an empty cell says the same thing, and `read_csv` turns it into
`NaN` without being asked.

The network is untouched: the bank answers with one small JSON, and `json.load`
is still the right tool for that.

```
pip install -r requirements.txt
python3 main.py        # first run: fetches and saves
python3 main.py        # second run: reads from disk
```

`data/` is what the program produces and is not kept in the repository.
