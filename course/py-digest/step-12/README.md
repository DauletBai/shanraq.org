# Step 12 — the year becomes a date

The state of the digest after lesson 31, *Time in a table*.

Two things changed, and the smaller one matters more later.

`sholu/tazalau.py` turns the year into a real date — the first day of that year.
A year as a number sorts and subtracts by luck; a date says what it is, and the
day the digest reads a monthly series nothing above this line will have to
change.

`sholu/esep.py` gains the question a yearly series can answer: how much the last
year differs from the one before it. `diff()` inside the grouping does it in one
line and without the off-by-one a loop over sorted years invites — it needs the
rows in order, which the cleaning step guarantees.

What is deliberately **not** here is the other half of lesson 31: `resample` and
`rolling`. Five yearly points are not a series to smooth, and a five-point
rolling mean would be a picture with nothing behind it. They arrive when the
digest reads something denser than one figure a year.

```
pip install -r requirements.txt
python3 main.py        # first run: fetches and saves
python3 main.py        # second run: reads from disk
```

`data/` is what the program produces and is not kept in the repository.
