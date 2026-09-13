# Step 23 — the model is asked about a year it did not see

The state of the digest after lesson 44, *Training and checking*.

The trend of the previous step was measured on the years it was fitted to, which
is the one measurement a model can always pass. From this step it is taught a
second time without the last year and then asked about it, and both errors stand
on the page:

```
аты        жылдық  ең үлкен қате  соңғы жылда
Қазақстан    1,12           3,55         4,85
RUS          1,09           2,95         1,32
Өзбекстан    1,10           1,02         2,56
```

- **`ең үлкен қате`** is the worst miss on the years the line learned from.
- **`соңғы жылда`** is the miss on the year it did not get. For Kazakhstan it is
  larger, as it usually is. For Russia it is *smaller* — and that is the honest
  reason the column is not called "the real error": **one held-out year is a
  start, not a verdict.** Five points leave four for the fitting, and a check of
  a single number goes the way it goes.
- **Still no forecast.** A line that misses a known year by five points does not
  become trustworthy about an unknown one by being asked nicely.

The store stopped rounding in this step, which is a fix rather than a feature.
`saqtau.save` wrote two decimals, so the same report said 3.55 on a run that
fetched and 3.56 on a run that read the file back: enough precision for a rate
of inflation, not enough for a line fitted through sixteen of them. A store that
changes the answer is not a store.

```
pip install -r requirements.txt
python3 main.py --dry-run     # count, write nothing
python3 main.py               # rebuild the report
```
