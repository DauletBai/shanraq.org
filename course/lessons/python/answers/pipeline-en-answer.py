"""The exercise of lesson 36: a pipeline with a log, a dry run and an exit code.

One main() reads like a table of contents: take, clean, count, show. Every step
writes a line into the log, --dry-run counts and writes nothing, and a source
that fell over leaves yesterday report where it was.
"""

import os
import sys
from pathlib import Path

HERE = Path(__file__).resolve().parent
REPORT = HERE / "svodka.txt"


def fetch(broken=False):
    if broken:
        raise RuntimeError("the source did not answer")
    return [("2023", 14.5), ("2024", 8.7), ("2025", None), ("2025", 11.4)]


def clean(rows):
    """No value, no row; a repeated year keeps the later one."""
    by_year = {}
    for year, value in rows:
        if value is not None:
            by_year[year] = value
    return sorted(by_year.items())


def count(rows):
    values = [value for _, value in rows]
    return {"years": len(values), "mean": round(sum(values) / len(values), 2)}


def render(rows, totals):
    lines = [f"{year}: {value}" for year, value in rows]
    lines.append(f"the mean over {totals['years']} years: {totals['mean']}")
    return "\n".join(lines) + "\n"


def save(text, path):
    temporary = path.with_suffix(path.suffix + ".tmp")
    temporary.write_text(text, encoding="utf-8")
    os.replace(temporary, path)


def main(argv=()):
    dry = "--dry-run" in argv
    broken = "--broken" in argv
    log = []
    try:
        rows = fetch(broken=broken)
    except RuntimeError as error:
        log.append(f"the source: {error}")
        print("\n".join(f"  {line}" for line in log))
        return 1
    log.append(f"rows taken: {len(rows)}")

    rows = clean(rows)
    log.append(f"after cleaning: {len(rows)}")

    totals = count(rows)
    log.append(f"counted: {totals}")

    if dry:
        log.append("a dry run: the file was not touched")
    else:
        save(render(rows, totals), REPORT)
        log.append(f"written: {REPORT.name}")

    print("\n".join(f"  {line}" for line in log))
    return 0


print("== an ordinary run")
code = main([])
saved = REPORT.read_text(encoding="utf-8")
print("  code:", code)

print()
print("== a dry run")
code = main(["--dry-run"])
print("  code:", code, "| the file did not change:", REPORT.read_text(encoding="utf-8") == saved)

print()
print("== the source fell over")
code = main(["--broken"])
print("  code:", code, "| yesterday report is in place:", REPORT.read_text(encoding="utf-8") == saved)

print()
print("the report:")
print(saved, end="")
sys.exit(0)
