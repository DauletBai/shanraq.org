"""The answer to lesson 23's exercise: a run nobody is watching.

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

parser = argparse.ArgumentParser(description="Сводка за один день")
parser.add_argument("--day", default="2026-01-15", help="день в виде ГГГГ-ММ-ДД")
parser.add_argument("--dry-run", action="store_true", help="посчитать, но не записывать")
args = parser.parse_args(["--day", "2026-01-16"])

logging.basicConfig(filename=LOG, filemode="w", encoding="utf-8", level=logging.INFO,
                    format="%(asctime)s | %(levelname)s | %(message)s")
log = logging.getLogger("svodka")


def locked():
    """Замок занят кем-то другим?"""
    try:
        LOCK.touch(exist_ok=False)
        return False
    except FileExistsError:
        return True


def run(day, dry):
    """Одна ночная работа: считает и, если не dry-run, записывает."""
    if locked():
        log.warning("день %s: другой запуск уже идёт, выходим", day)
        return 1
    try:
        log.info("день %s: начали", day)
        if dry:
            log.info("день %s: считаем без записи", day)
        else:
            log.info("день %s: записано 2 строки", day)
        return 0
    finally:
        LOCK.unlink()


print("первый запуск, код:", run(args.day, args.dry_run))
LOCK.touch()
print("запуск при занятом замке, код:", run(args.day, args.dry_run))
LOCK.unlink()
print("запуск без записи, код:", run(args.day, True))

print("журнал:")
for line in LOG.read_text(encoding="utf-8").splitlines():
    print("  " + line.split(" | ", 1)[1])

LOG.unlink()
