# Step 10 — the report stops speaking in codes

The state of the digest after lesson 29, *Joining tables*.

The bank sends codes: `KAZ`, `UZB`, `RUS`. The names are not its business and
they are not in its answer, so they are ours: `sholu/anyqtama.py` is a table of
three columns — code, name, region — and the report joins it on the code.

Two decisions worth reading the code for:

- the join is a **left** one. The data decides which rows exist; the directory
  only adds words to them. A country missing from the directory stays in the
  report under its code instead of vanishing from it;
- the join is **checked**: `validate="one_to_one"`. A code repeated in the
  directory would silently multiply the rows it matched, and the totals would
  grow without anything being added.

```
pip install -r requirements.txt
python3 main.py        # first run: fetches and saves
python3 main.py        # second run: reads from disk
```

`data/` is what the program produces and is not kept in the repository.
