# Step 20 — the report says what its numbers mean

The state of the digest after lesson 41, *How to read official statistics*.

The digest has printed a rate of inflation since its first step, and never said
which rate. An annual figure is an average over the year against the average
over the year before in one publication, and December against December in
another; for Kazakhstan in 2022 those two differ by five points. A reader who
does not know which one is on the page cannot use the page.

So from this step the definition travels with the number:

- **`sholu/derekkoz.py`, `about()`** asks the World Bank what its own indicator
  is — `https://api.worldbank.org/v2/indicator/<code>` — and keeps the name, the
  source and the note in the source's own words.
- **The notes are cached** in `data/indicators.json` and re-read from disk; they
  change about as often as the definition of inflation does.
- **A note is a nicety, not the data.** A source that will not answer costs the
  page its notes and nothing else: the log says what is missing and the report
  is rebuilt anyway. `--offline` is honoured rather than quietly broken for a
  paragraph of text.
- **`sholu/bet.py`** puts them at the foot of the page, under everything they
  explain.

```
pip install -r requirements.txt
python3 main.py --dry-run     # count, write nothing
python3 main.py               # rebuild the report
```

The page now ends with what the bank says its indicators are — including the
line that matters most for this course: a basket "that may be fixed or changed
at specified intervals".
