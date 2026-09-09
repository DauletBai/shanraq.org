"""The answer to lesson 23's exercise, Kazakh: a run nobody is watching.

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

parser = argparse.ArgumentParser(description="Бір күндік шолу")
parser.add_argument("--day", default="2026-01-15", help="күн ЖЖЖЖ-АА-КК түрінде")
parser.add_argument("--dry-run", action="store_true", help="санау, бірақ жазбау")
args = parser.parse_args(["--day", "2026-01-16"])

logging.basicConfig(filename=LOG, filemode="w", encoding="utf-8", level=logging.INFO,
                    format="%(asctime)s | %(levelname)s | %(message)s")
log = logging.getLogger("svodka")


def locked():
    """Құлыпты басқа біреу алып қойған ба?"""
    try:
        LOCK.touch(exist_ok=False)
        return False
    except FileExistsError:
        return True


def run(day, dry):
    """Бір түнгі жұмыс: санайды, ал dry-run болмаса, жазады."""
    if locked():
        log.warning("%s күні: басқа іске қосу жүріп жатыр, шығамыз", day)
        return 1
    try:
        log.info("%s күні: бастадық", day)
        if dry:
            log.info("%s күні: жазбай санаймыз", day)
        else:
            log.info("%s күні: 2 жол жазылды", day)
        return 0
    finally:
        LOCK.unlink()


print("бірінші іске қосу, коды:", run(args.day, args.dry_run))
LOCK.touch()
print("құлып бос емес кездегі іске қосу, коды:", run(args.day, args.dry_run))
LOCK.unlink()
print("жазбай іске қосу, коды:", run(args.day, True))

print("журнал:")
for line in LOG.read_text(encoding="utf-8").splitlines():
    print("  " + line.split(" | ", 1)[1])

LOG.unlink()
