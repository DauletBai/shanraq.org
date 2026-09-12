# Step 19 — how much money there is for a tenge of output

The state of the digest after lesson 40, *Where money comes from*.

Until now the digest read one series. From this step it reads two, and the
second one is not prices at all: broad money as a share of GDP, that is, how
much money a country has for every tenge of a year of its output.

What changed, in the order it matters:

- **`collect()` takes the indicator and the file**, instead of knowing about
  inflation and nothing else. A third series would now cost a line rather than a
  copy of this program; `data/` holds `inflation.csv` and `aqsha.csv` side by
  side.
- **`sholu/esep.py`, `aqsha()`** turns a share of GDP into tiyn per tenge —
  33 for Kazakhstan, 20 for Uzbekistan — and puts the first year of the run next
  to the last one, so the row says which way the country moved.
- **A country the source knows nothing about keeps its row.** The World Bank has
  no broad money for Russia in these years, and the report says `белгісіз`
  rather than dropping the country and looking complete.
- **The page carries its second source line.** Two indicators and one credit
  would be a page claiming something it cannot back.

```
pip install -r requirements.txt
python3 main.py --dry-run     # count, write nothing
python3 main.py               # rebuild the report
```

The page now has four tables: the last year per country, the whole series as a
centre and a spread, the run of years as one multiplier, and the money a country
holds against a year of its output.
