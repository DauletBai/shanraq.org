# Step 17 — how usual was that year

The state of the digest after lesson 37, *The mean, the median and the spread*,
and the start of the fifth module.

The report used to end at the last year: 11.39 % in Kazakhstan, and nothing about
whether that is a lot for this country. The page now carries a second table —
over every year the digest holds, per country:

- **the observations**, because 11.54 over five years and over fifty are
  different statements;
- **the mean and the median**, side by side: when they disagree, the series has
  a tail and the mean is describing it rather than the country;
- **the spread and the interquartile range**: one number for how far the values
  sit from the mean, and one for the band the middle half falls into;
- **the measure for the heading**, chosen by a rule in `sholu/esep.py`
  (`headline`) rather than by taste.

On today's data the rule picks the mean for all three countries — their
five-year series are symmetrical enough — and it will say the opposite the day a
spike arrives in one of them.

```
pip install -r requirements.txt
python3 main.py --dry-run     # count, write nothing
python3 main.py               # rebuild the report
```

The files it writes are the same three as before — `data/report.csv`,
`data/inflation.png`, `data/report.html` — and the page now has two tables in it.
