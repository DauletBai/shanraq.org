# Step 16 — the digest updates itself

The state of the digest after lesson 36, *A digest that updates itself*, and the
end of the fourth module.

`main()` now reads like a table of contents — take, clean, count, draw, collect
— and around it are the four things a program run by a scheduler needs while a
program run by a person does not:

- **`--dry-run`**: everything is computed and nothing is written. The first
  thing to run after a change, and the only safe way to try it at four in the
  afternoon on a Friday;
- **`--offline`**: never ask the bank. Useful on a plane, in a test, and every
  time the question is "is it my code or their server";
- **an exit code**: `0` when the report was rebuilt, `1` when it was not. That
  number is all a scheduler reads ([lesson 23](../../lessons/python/schedule.md));
- **a log line per step**, so the morning after is readable.

And one rule that is not a flag: a source that does not answer leaves yesterday's
report where it was. An empty report is worse than an old one — an old one is at
least true about an older day.

`sholu/saqtau.py` writes everything through a temporary file and `os.replace`,
which renames in one motion: a reader opening the report while it is being
rebuilt sees either the old one or the new one, never half of either.

```
pip install -r requirements.txt
python3 main.py --dry-run     # count, write nothing
python3 main.py               # rebuild the report
python3 main.py --offline     # never touch the network
echo $?                       # 0 done, 1 not done
```

A line for cron, at five past six every morning:

```
5 6 * * * cd /home/you/digest && .venv/bin/python main.py >> data/run.log 2>&1
```

`data/` is what the program produces and is not kept in the repository.
