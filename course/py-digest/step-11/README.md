# Step 11 — a cleaning step, and a log of it

The state of the digest after lesson 30, *Dirty data*.

The bank is a tidy source and still sends a year with no figure in it, because
the year is not over yet. A file read back off the disk can hold a row twice if
a run was interrupted. Until now every report decided what to do about that on
its own; now one place decides, before the report is built.

`sholu/tazalau.py` returns two things: the table fit to count, and a log:

```
жол келді: 15
саны жоқ жыл: 0
есепке кететін жол: 15
```

The rules it follows are the lesson's:

- a row without a key is dropped — it is not a gap in the data but the absence
  of the row itself;
- a repeated `(country, year)` is dropped, the later one kept;
- what is not a number becomes a gap rather than a zero: a year not yet counted
  is not a year of zero inflation;
- the country code becomes a `category` — five words in a column of many rows
  are stored once.

The report keeps one decision of its own, and it is about words rather than
numbers: a country the directory does not know gets the marker `белгісіз` for a
region instead of an empty cell.

```
pip install -r requirements.txt
python3 main.py        # first run: fetches and saves
python3 main.py        # second run: reads from disk
```

`data/` is what the program produces and is not kept in the repository.
