# Step 7 — the series becomes a table

The state of the digest after lesson 25, *DataFrame and Series*.

The dictionary `{year: value}` is still what comes back from the bank and still
what goes to disk. What changed is the middle: `sholu/esep.py` turns it into a
`DataFrame` with the years as row labels, and the report is built from that
table rather than from three loops.

Three loops the table removed:

- the average that had to be told to skip the years without a figure — `mean()`
  leaves them out on its own;
- the count of the gaps — `isna().sum()`;
- the search for the largest value that had to carry its year along — `idxmax()`
  answers with the year, because the year is the label.

```
pip install -r requirements.txt
python3 main.py        # first run: fetches and saves
python3 main.py        # second run: reads from disk
```

The disk layer is untouched: the CSV is still written by `csv.DictWriter` and
read by `csv.DictReader`. `read_csv` and `to_csv` come in the next lesson.

`data/` is what the program produces and is not kept in the repository.
