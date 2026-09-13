# Step 22 — the digest fits its first model

The state of the digest after lesson 43, *The first model: linear regression*.

Until now the digest counted: an average, an index, a share, a link. This step
fits — a straight line through the logarithm of each country's price index,
which is the shape prices actually have, because they multiply rather than add.

What the page gains is two numbers per country:

```
аты        жыл  жылдық  ең үлкен қате  қорытынды
Қазақстан    5    1,12           3,56  саналды
RUS          5    1,09           2,95  саналды
Өзбекстан    5    1,10           1,02  саналды
```

- **`жылдық`** is the multiplier the line implies for one year — 1.12 for
  Kazakhstan, that is about 12 % a year over the run. The log line carries it
  unrounded (1.1222), the page rounds like every other number on it.
- **`ең үлкен қате`** is the worst the line misses by, in points of the index.
  It is printed next to the multiplier on purpose: a trend without its error is
  a claim without a cost, and one of these countries is described three times
  better than another.
- **`қорытынды`** says `аз: N жыл` when a country has fewer years than the
  floor. Five is the floor and five is what the digest has, so the line is thin
  on purpose and the report says how thin.

**No forecast is printed.** The line could produce one in a single call, and
that is exactly why it is worth saying no: a model is worth a forecast only
after it has been checked on years it did not see. That check is the next
lesson, and until it exists the digest reports what the line says about the
years it was given and nothing about the ones it was not.

```
pip install -r requirements.txt
python3 main.py --dry-run     # count, write nothing
python3 main.py               # rebuild the report
```

`requirements.txt` gains scikit-learn, pinned like the rest: a model whose
library version is not written down is a model nobody can reproduce.
