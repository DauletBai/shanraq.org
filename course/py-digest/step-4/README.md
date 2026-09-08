# Step 4 — the digest gets a disk

The state of the digest after lesson 11, *Files and pathlib*.

The series fetched from the bank is saved into `data/inflation-kz.txt` beside
the program, and the next run reads that file instead of asking again. The
report is written as a second file. Every path starts at
`Path(__file__).parent`, so the program works when it is started from another
folder.

```
python3 main.py        # first run: fetches and saves
python3 main.py        # second run: reads from disk
```

`data/` is not kept in the repository: it is what the program produces, and it
is rebuilt by deleting the folder and running again.
