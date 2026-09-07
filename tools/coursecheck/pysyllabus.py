#!/usr/bin/env python3
"""Keep the Python course from forgetting a piece of the language.

In the Go course the gaps -- arrays, constants, closures, recursion -- were not
postponed but forgotten, and closing them afterwards cost more than planning
them would have. So here every element is assigned to a lesson in
docs/py-syllabus.md before the lesson is written, and this checks two things:

    plan            every element has a lesson, and no lesson is out of range
    lesson N file   the elements assigned to lesson N appear in that text

A marker is a word or a fragment of code the explanation cannot avoid. It is a
coarse instrument on purpose: it does not judge whether the explanation is good,
only whether it happened at all.
"""

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
SYLLABUS = ROOT / "docs" / "py-syllabus.md"
PLAN = ROOT / "docs" / "py-course.md"

ROW = re.compile(r"^\|\s*(?P<what>[^|]+?)\s*\|\s*(?P<lesson>\d+)\s*\|\s*(?P<markers>[^|]+?)\s*\|\s*$", re.M)


def rows():
    """Every assignment in the table: what, which lesson, which markers."""
    out = []
    for m in ROW.finditer(SYLLABUS.read_text(encoding="utf-8")):
        markers = [x.strip().strip("`") for x in m.group("markers").split(",")]
        out.append((m.group("what"), int(m.group("lesson")), [x for x in markers if x]))
    return out


def lesson_count():
    """How many lessons the plan actually has, counted from its tables."""
    numbers = [int(m) for m in re.findall(r"^\|\s*(\d+)\s*\|", PLAN.read_text(encoding="utf-8"), re.M)]
    return max(numbers) if numbers else 0


def check_plan():
    items = rows()
    top = lesson_count()
    bad = [f"  {what}: урок {n}, а в плане всего {top}" for what, n, _ in items if n > top]
    print(f"элементов в списке: {len(items)}, уроков в плане: {top}")
    if bad:
        print("назначено мимо плана:")
        print("\n".join(bad))
        return 1
    covered = sorted({n for _, n, _ in items})
    empty = [n for n in range(3, top + 1) if n not in covered]
    if empty:
        print("уроки без назначенных элементов (это нормально для обзорных):")
        print("  " + ", ".join(str(n) for n in empty))
    print("все элементы назначены урокам")
    return 0


def check_lesson(number, path):
    text = Path(path).read_text(encoding="utf-8")
    mine = [(what, markers) for what, n, markers in rows() if n == number]
    if not mine:
        print(f"уроку {number} элементы не назначены — проверять нечего")
        return 0
    missing = []
    for what, markers in mine:
        if not any(re.search(m, text) if any(c in m for c in ".*[]()\\") else m in text for m in markers):
            missing.append(f"  {what}: не нашёл ни одного из {markers}")
    print(f"урок {number}: элементов назначено {len(mine)}, не найдено {len(missing)}")
    if missing:
        print("\n".join(missing))
        return 1
    return 0


def main(argv):
    if not argv or argv[0] not in {"plan", "lesson"}:
        print(__doc__)
        return 2
    if argv[0] == "plan":
        return check_plan()
    if len(argv) < 3:
        print("нужно: pysyllabus.py lesson НОМЕР файл.md")
        return 2
    return check_lesson(int(argv[1]), argv[2])


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
