# Step 9 — three countries and one grouping

The state of the digest after lesson 28, *Grouping and aggregates*.

Until now the digest watched one country. Three of them arrive in a single
request — the bank takes the codes separated by semicolons — and the table is
long: one row per country and year.

That shape is what makes the report a single question instead of a loop:

```
table.groupby("country").agg(...)
```

The peak year comes from `idxmax` inside the grouping: the label of the largest
value is the row it sits in, and that row already knows its year.

```
pip install -r requirements.txt
python3 main.py        # first run: fetches and saves
python3 main.py        # second run: reads from disk
```

`data/` is what the program produces and is not kept in the repository.
