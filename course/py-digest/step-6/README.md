# Step 6 — the digest is laid out in files

The state of the digest after lesson 17, *A module and a package of your own*.

Step 5 was one file that fetched, saved, counted and printed. It still does all
four, and nothing about the result has changed — what changed is where each of
them lives:

```
main.py            what happens, in order
sholu/derekkoz.py  the network, and nothing else knows about it
sholu/saqtau.py    the disk: writing the CSV and reading it back
sholu/esep.py      the arithmetic and the report
```

`main.py` also gains the line the lesson is about: `if __name__ == "__main__":`,
so that importing it does not set the whole run off.

```
python3 main.py        # first run: fetches and saves
python3 main.py        # second run: reads from disk
```

`data/` is what the program produces and is not kept in the repository.
