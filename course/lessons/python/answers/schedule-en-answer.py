"""The answer to lesson 23's exercise, English: a run nobody is watching.

Three things make the difference between a program that is scheduled and a
program that merely runs: arguments, so the day is chosen without editing the
source; a journal, so the morning can be told what the night did; and a lock,
so two runs started a minute apart do not do the same work twice. The exit code
is the fourth: it is the only thing the scheduler itself reads.
"""

import argparse
import logging
from pathlib import Path

HERE = Path(__file__).resolve().parent
LOG = HERE / "svodka.log"
LOCK = HERE / "svodka.lock"

parser = argparse.ArgumentParser(description="The digest for one day")
parser.add_argument("--day", default="2026-01-15", help="the day as YYYY-MM-DD")
parser.add_argument("--dry-run", action="store_true", help="count without writing")
args = parser.parse_args(["--day", "2026-01-16"])

logging.basicConfig(filename=LOG, filemode="w", encoding="utf-8", level=logging.INFO,
                    format="%(asctime)s | %(levelname)s | %(message)s")
log = logging.getLogger("svodka")


def locked():
    """Is the lock held by somebody else?"""
    try:
        LOCK.touch(exist_ok=False)
        return False
    except FileExistsError:
        return True


def run(day, dry):
    """One night's work: counts, and writes unless this is a dry run."""
    if locked():
        log.warning("day %s: another run is under way, leaving", day)
        return 1
    try:
        log.info("day %s: started", day)
        if dry:
            log.info("day %s: counting without writing", day)
        else:
            log.info("day %s: 2 rows written", day)
        return 0
    finally:
        LOCK.unlink()


print("the first run, code:", run(args.day, args.dry_run))
LOCK.touch()
print("a run while the lock is held, code:", run(args.day, args.dry_run))
LOCK.unlink()
print("a dry run, code:", run(args.day, True))

print("the journal:")
for line in LOG.read_text(encoding="utf-8").splitlines():
    print("  " + line.split(" | ", 1)[1])

LOG.unlink()
