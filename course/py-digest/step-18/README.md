# Step 18 — how many times prices grew

The state of the digest after lesson 38, *Inflation as a multiplier*.

The report knew each year on its own and, since the last step, the whole run of
them as a centre and a spread. What it could not say is the one thing a reader
asks first: **how much have prices grown since the series began**. That is not a
sum of the rates — percentages are multipliers — so the page now carries a third
table:

- **the base and the last year**, because an index without its base means
  nothing;
- **the index and the multiplier**: 1.595 for Kazakhstan over 2021–2025, which
  no sum of the yearly rates would give;
- **what a thousand of the base year is worth** at the end of the run;
- **the gaps**, if any: an index over a broken run of years is meaningless, so
  `sholu/esep.py` reports a missing year rather than skipping its multiplier in
  silence.

```
pip install -r requirements.txt
python3 main.py --dry-run     # count, write nothing
python3 main.py               # rebuild the report
```

The page now has three tables: the last year per country, the whole series as a
centre and a spread, and the run of years as one multiplier.
