# Step 13 — the report gets a picture

The state of the digest after lesson 32, *The first chart*.

`sholu/suret.py` draws one line per country and saves it next to the CSV. The
names for the legend come from the same lookup table the report reads, so the
picture says "Қазақстан" where the data says "KAZ".

Two things in that file are there because the digest runs on a schedule rather
than in front of a person:

- `matplotlib.use("Agg")` before `pyplot` is imported. There is no screen at
  four in the morning, and a backend that wants one turns a working program into
  a failing one;
- `plt.close(fig)` after the save. A figure nobody closes stays in memory, and a
  program drawing one a day finds that out on the twenty-first day.

```
pip install -r requirements.txt
python3 main.py        # first run: fetches, saves, draws
python3 main.py        # second run: reads from disk
```

`data/` is what the program produces — the CSV, the report and now the picture —
and is not kept in the repository.
