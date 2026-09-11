"""Решение задания урока 36: конвейер с журналом, сухим прогоном и кодом возврата.

Один main() читается как оглавление: забрать, привести, посчитать, показать.
Каждый шаг пишет строку в журнал, --dry-run считает, но не пишет, а упавший
источник оставляет вчерашний отчёт на месте.
"""

import os
import sys
from pathlib import Path

HERE = Path(__file__).resolve().parent
REPORT = HERE / "svodka.txt"


def fetch(broken=False):
    if broken:
        raise RuntimeError("источник не ответил")
    return [("2023", 14.5), ("2024", 8.7), ("2025", None), ("2025", 11.4)]


def clean(rows):
    """Без значения — не строка; повтор года — берём последний."""
    by_year = {}
    for year, value in rows:
        if value is not None:
            by_year[year] = value
    return sorted(by_year.items())


def count(rows):
    values = [value for _, value in rows]
    return {"лет": len(values), "среднее": round(sum(values) / len(values), 2)}


def render(rows, totals):
    lines = [f"{year}: {value}" for year, value in rows]
    lines.append(f"среднее за {totals['лет']} года: {totals['среднее']}")
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
        log.append(f"источник: {error}")
        print("\n".join(f"  {line}" for line in log))
        return 1
    log.append(f"взято строк: {len(rows)}")

    rows = clean(rows)
    log.append(f"после уборки: {len(rows)}")

    totals = count(rows)
    log.append(f"посчитано: {totals}")

    if dry:
        log.append("сухой прогон: файл не тронут")
    else:
        save(render(rows, totals), REPORT)
        log.append(f"записано: {REPORT.name}")

    print("\n".join(f"  {line}" for line in log))
    return 0


print("== обычный запуск")
code = main([])
saved = REPORT.read_text(encoding="utf-8")
print("  код:", code)

print()
print("== сухой прогон")
code = main(["--dry-run"])
print("  код:", code, "| файл не изменился:", REPORT.read_text(encoding="utf-8") == saved)

print()
print("== источник упал")
code = main(["--broken"])
print("  код:", code, "| вчерашний отчёт на месте:", REPORT.read_text(encoding="utf-8") == saved)

print()
print("отчёт:")
print(saved, end="")
sys.exit(0)
