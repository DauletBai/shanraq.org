# On a schedule: arguments, a journal, a lock and an exit code

_Лид (summary):_ **The twenty-third lesson of the Python course. A program that works through the night without you: arguments instead of editing the source, a journal with levels instead of `print`, a lock against two runs at once, and the exit code — the one thing the scheduler actually reads. Plus cron, launchd and systemd.**

## Why this is needed

The digest has learned to fetch data, put it into a database and answer with queries. One thing is left: all of that happening without you.

A program started by a scheduler lives in another world. Nobody sees its output. Nobody edits the day in its source. It may start a second time while the first is still running. And the only thing the system will learn about it is the number it returns.

So this lesson is not about cron. It is about the four things without which cron is useless.

## The whole thing at once

The file is `kesteme.py`. Run it with `python kesteme.py` from inside the environment.

The required part is the first three blocks: the arguments, the journal and the lock. The fourth shows the exit code, which is the reason a scheduler looks at a program at all.

```python
"""Lesson 23: a program that works while you are away.

It is time the digest ran on a schedule. Three things separate that from hope:
arguments instead of editing the source, a journal instead of print, and one
copy running instead of two.
"""

import argparse
import logging
import subprocess
import sys
from pathlib import Path

HERE = Path(__file__).resolve().parent
LOG = HERE / "svodka.log"
LOCK = HERE / "svodka.lock"

print("== arguments instead of editing the source")
parser = argparse.ArgumentParser(description="The digest for one day")
parser.add_argument("--day", default="today", help="the day as YYYY-MM-DD")
parser.add_argument("--dry-run", action="store_true", help="count without writing")
parser.add_argument("--verbose", action="store_true", help="a fuller journal")
# The arguments are handed in outright, so the lesson prints the same for all.
args = parser.parse_args(["--day", "2026-01-15", "--verbose"])
print("day:", args.day, "| dry run:", args.dry_run, "| verbose:", args.verbose)
print("defaults:", parser.parse_args([]).day, parser.parse_args([]).dry_run)

print()
print("== a journal instead of print")
logging.basicConfig(
    filename=LOG, filemode="w", encoding="utf-8",
    level=logging.DEBUG if args.verbose else logging.INFO,
    format="%(asctime)s | %(levelname)s | %(message)s",
)
log = logging.getLogger("svodka")
log.debug("detail: parsing %s", args.day)
log.info("day %s: 3 values taken", args.day)
log.warning("no rate for 2026-01-14, skipping")
try:
    1 / 0
except ZeroDivisionError:
    log.exception("division by zero while averaging")

# The time in the journal belongs to the run, so the lesson prints what comes
# after it, and does not print the stack lines log.exception added at all.
for line in LOG.read_text(encoding="utf-8").splitlines():
    if " | " in line:
        print(line.split(" | ", 1)[1])
print("and after it the file holds the stack — everybody's is their own")

print()
print("== one copy running instead of two")


def start():
    """Takes the lock, or refuses to work."""
    try:
        LOCK.touch(exist_ok=False)
        return True
    except FileExistsError:
        return False


print("the first run took the lock:", start())
print("a second run while it is held:", start())
LOCK.unlink()
print("after it is released:", start())
LOCK.unlink()

print()
print("== the scheduler reads the exit code")
worker = HERE / "rabota.py"
worker.write_text(
    "import sys\n"
    "sys.stderr.write('the source did not answer\\n')\n"
    "raise SystemExit(2)\n", encoding="utf-8")
done = subprocess.run([sys.executable, worker], capture_output=True, text=True)
print("exit code:", done.returncode, "| on stderr:", done.stderr.strip())
print("zero means all is well:", subprocess.run([sys.executable, "-c", "pass"]).returncode)

worker.unlink()
LOG.unlink()
```

It prints:

```
== arguments instead of editing the source
day: 2026-01-15 | dry run: False | verbose: True
defaults: today False

== a journal instead of print
DEBUG | detail: parsing 2026-01-15
INFO | day 2026-01-15: 3 values taken
WARNING | no rate for 2026-01-14, skipping
ERROR | division by zero while averaging
and after it the file holds the stack — everybody's is their own

== one copy running instead of two
the first run took the lock: True
a second run while it is held: False
after it is released: True

== the scheduler reads the exit code
exit code: 2 | on stderr: the source did not answer
zero means all is well: 0
```

## Taking it apart

### Arguments instead of editing the source

```python
parser.add_argument("--day", default="today", help="the day as YYYY-MM-DD")
args = parser.parse_args(["--day", "2026-01-15", "--verbose"])
```

`argparse` turns a program into a tool: the day, the mode and the verbosity are chosen **at the start**, not by editing a file. For a schedule that is not convenience but necessity — you cannot edit code from a line of cron.

Three things it gives for nothing: a `--help` with every description in it, a clear error on an unknown argument, and defaults that are visible in one place. `action="store_true"` is a flag: it is either there or it is not.

Usually one writes `parser.parse_args()` with nothing in the brackets, and the real command line is parsed. The lesson hands the list in outright so that the output is the same for everybody; in your program the brackets will be empty.

> **Picture it.** A car key. A program with arguments is a car you can start; a program with the number written inside is a car whose wheel is welded in one position.

### A journal instead of `print`

```
DEBUG | detail: parsing 2026-01-15
INFO | day 2026-01-15: 3 values taken
WARNING | no rate for 2026-01-14, skipping
ERROR | division by zero while averaging
```

`print` writes into nowhere: at night nobody is looking at a terminal, and cron will at best post the output to you. `logging` writes into a file, and every line has a **level** and a **time**.

The levels are not decoration. `DEBUG` is turned on while something is being investigated; `INFO` is the ordinary course of things; `WARNING` is "we carry on, but notice this"; `ERROR` is "it did not work". The level is set once, and the `--verbose` from the first block does exactly that: `level=DEBUG` instead of `INFO`.

`log.exception(...)` inside an `except` writes the **stack of the error** into the file — the one you read in lesson ten. That is the main reason to keep a journal: in the morning you have not "the program fell over" but the line it fell over on.

And a small thing that matters: `log.info("day %s: %d taken", day, n)` — with percent signs rather than an f-string. That way the text is assembled only if that level is on.

### A lock: one copy running instead of two

```
the first run took the lock: True
a second run while it is held: False
```

The night's run ran long, the time for the next one came — and now there are two. Both read the same source, both write into the same table; at best you pay in double traffic, at worst you get half the data from one and half from the other.

`LOCK.touch(exist_ok=False)` creates the file and **fails** if it is already there. That is a lock: the check for "is it taken" and the taking, in one indivisible act. An `if LOCK.exists()` on its own will not do — a second run fits between the check and the creation.

The lock is always released in a `finally`: a program that left its lock behind after a crash will never start again, and mending that is a job for hands.

### The exit code: the one thing the scheduler reads

```
exit code: 2 | on stderr: the source did not answer
zero means all is well: 0
```

The scheduler does not read your journal. It looks at the number the program returned: `0` means it worked, anything else means it did not. Notifications, retries and monitoring are all built on that.

So a night-time program has a contract: `raise SystemExit(1)` when the work was not done, and zero when it was. The error itself is better printed to `stderr` — cron keeps the two streams apart.

### The schedule itself

One line. On Linux and macOS it is `crontab -e`:

```
# minute hour day month weekday  command
30 7 * * *  /home/user/digest/.venv/bin/python /home/user/digest/main.py --day today >> /home/user/digest/cron.log 2>&1
```

Three rules everybody trips over:

- **Full paths only.** cron has neither your `PATH` nor your current folder — hence `Path(__file__).resolve().parent` in the program and the full path to the environment's Python here.
- **The timezone is the system's**, not yours. `30 7` on a server in UTC is half past twelve in Almaty.
- **Redirect the output.** Without `>> ... 2>&1` cron will try to post it, and there is usually no mail on a server.

On macOS the same is done by `launchd` (a `.plist` file in `~/Library/LaunchAgents`), on a modern Linux by a `systemd` timer (a `.service` and a `.timer`), on Windows by Task Scheduler. They differ in syntax; everything this lesson is about is the same for all of them.

## The map of the lesson

![The map of the lesson: the arguments, the journal and the lock](/static/course/py/map-schedule-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why is the day chosen with an argument rather than a line in the source?
2. How does `logging` differ from `print` for a program that runs at night?
3. What is a lock for, and why does `if LOCK.exists()` not replace it?

## Warm-up

Three short steps before the exercise: predict, fill in, fix. The answers are at the end of the lesson, but answer them yourself first.

**1. Predict.** What does this program print?

<!-- drill 1 -->
```python
import logging
import sys

logging.basicConfig(stream=sys.stdout, level=logging.INFO,
                    format="%(levelname)s: %(message)s")
log = logging.getLogger("demo")
log.debug("a detail")
log.info("done")
log.warning("something is off")
```

**2. Fill in the gap.** In place of `...` set the default that gets printed when the command line is empty.

```python
import argparse

parser = argparse.ArgumentParser()
parser.add_argument("--day", ...)
print(parser.parse_args([]).day)
```

**3. Fix it.** The scheduler starts the program from a folder of its own, and the file beside it stops being found. Make the path such that it does not matter.

```python
# the program puts a file beside itself and then walks off into another folder —
# which is exactly how a scheduler starts it
import os
from pathlib import Path

(Path(__file__).resolve().parent / "dannye.csv").write_text("2026;8.0\n", encoding="utf-8")
os.chdir("/")                      # a scheduler starts it from its own folder
print("found the file:", Path("dannye.csv").exists())
```

## Exercise

**Required.** Write the digest's night run. It needs: the arguments `--day` (with a default) and `--dry-run`; a journal in a file, with levels and times; a function `run(day, dry)` that takes the lock, writes into the journal what it is doing, releases the lock in a `finally`, and returns `0` on success and `1` when the lock is held.

Call it three times: an ordinary run, a run while the lock is held (make the lock by hand), and a run with `--dry-run`. Print the exit code of each, and at the end the journal without its times.

The expected output:

<!-- task out -->
```
the first run, code: 0
a run while the lock is held, code: 1
a dry run, code: 0
the journal:
  INFO | day 2026-01-16: started
  INFO | day 2026-01-16: 2 rows written
  WARNING | day 2026-01-16: another run is under way, leaving
  INFO | day 2026-01-16: started
  INFO | day 2026-01-16: counting without writing
```

Done when: the output matches line by line; the lock is released in a `finally` rather than after the successful branch; a held lock gives a `WARNING` and a code of `1` rather than a crash; every path is built from `Path(__file__).resolve().parent`.

**On your own data.** Put your own program on a schedule on your own machine — cron, launchd or Task Scheduler — for five minutes from now. Wait for it and read the journal. That is the only way to find out whether you wrote the path to Python correctly.

**Optional.**

- Add a `--log-level` and learn to turn `DEBUG` on from the command line.
- Make the program fall over on purpose and see what `log.exception` put in the journal.
- Start two copies at once (`python main.py & python main.py`) and check that the second left with a code of `1`.

## Where this goes in the project

The digest becomes a service. Once a day it fetches the rate itself, puts it into the database, writes into the journal what it did, and returns zero — and you hear about it only when the zero stops coming.

Still open. Our lock is a file: it will not survive the power going out mid-run, and in the morning it will have to be removed by hand. And the journal grows without limit — rotation waits where the server does.

## The answers

### To the questions

1. Because a line in the source cannot be changed from a schedule: cron can start a command, not edit a file. An argument turns the program into a tool with modes, and `--help` into the place where those modes are described.
2. `print` writes where nobody is at night. `logging` writes into a file, marks importance with a level and the moment with a time, and `log.exception` puts the stack of the error beside it. In the morning that is the difference between "it fell over" and "it fell over here".
3. A lock keeps two runs from doing one job twice. `if LOCK.exists()` does not replace it because a second run fits between the check and the creation of the file: what is needed is one indivisible act, and `touch(exist_ok=False)` is that act.

### To the warm-up

1. Two lines: `INFO: done` and `WARNING: something is off`. `DEBUG` is below the level that was set, so it is not seen — the level is the switch for detail.

<!-- drill 1 out -->
```
INFO: done
WARNING: something is off
```

2. `default="today"`. A default is the answer to "and if it was not given", and it belongs where the argument is declared.

<!-- drill 2 -->
```python
import argparse

parser = argparse.ArgumentParser()
parser.add_argument("--day", default="today")
print(parser.parse_args([]).day)
```

<!-- drill 2 out -->
```
today
```

3. The path `Path("dannye.csv")` is relative to the current folder, and the scheduler sets that. Build from the program's own folder, and where it was called from stops mattering:

<!-- drill 3 -->
```python
import os
from pathlib import Path

HERE = Path(__file__).resolve().parent
(HERE / "dannye.csv").write_text("2026;8.0\n", encoding="utf-8")
os.chdir("/")                      # a scheduler starts it from its own folder
print("found the file:", (HERE / "dannye.csv").exists())
(HERE / "dannye.csv").unlink()
```

<!-- drill 3 out -->
```
found the file: True
```

## Sources

- [Python: argparse](https://docs.python.org/3/library/argparse.html)
- [Python: logging](https://docs.python.org/3/library/logging.html)
- [Python: the logging how-to](https://docs.python.org/3/howto/logging.html)
- [crontab: the format of a schedule line](https://man7.org/linux/man-pages/man5/crontab.5.html)
