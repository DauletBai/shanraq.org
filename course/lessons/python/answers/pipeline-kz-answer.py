"""36-сабақтың тапсырмасының шешімі: журналы, құрғақ жүрісі және қайтару коды бар конвейер.

Бір main() мазмұндама болып оқылады: алу, реттеу, санау, көрсету. Әр қадам
журналға жол жазады, --dry-run санайды да жазбайды, ал құлаған дереккөз кешегі
есепті орнында қалдырады.
"""

import os
import sys
from pathlib import Path

HERE = Path(__file__).resolve().parent
REPORT = HERE / "svodka.txt"


def fetch(broken=False):
    if broken:
        raise RuntimeError("дереккөз жауап бермеді")
    return [("2023", 14.5), ("2024", 8.7), ("2025", None), ("2025", 11.4)]


def clean(rows):
    """Мәні жоқ — жол емес; жыл қайталанса — соңғысын аламыз."""
    by_year = {}
    for year, value in rows:
        if value is not None:
            by_year[year] = value
    return sorted(by_year.items())


def count(rows):
    values = [value for _, value in rows]
    return {"жыл": len(values), "орташа": round(sum(values) / len(values), 2)}


def render(rows, totals):
    lines = [f"{year}: {value}" for year, value in rows]
    lines.append(f"{totals['жыл']} жылдың орташасы: {totals['орташа']}")
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
        log.append(f"дереккөз: {error}")
        print("\n".join(f"  {line}" for line in log))
        return 1
    log.append(f"алынған жол: {len(rows)}")

    rows = clean(rows)
    log.append(f"тазалаудан кейін: {len(rows)}")

    totals = count(rows)
    log.append(f"есептелді: {totals}")

    if dry:
        log.append("құрғақ жүріс: файл тиылмады")
    else:
        save(render(rows, totals), REPORT)
        log.append(f"жазылды: {REPORT.name}")

    print("\n".join(f"  {line}" for line in log))
    return 0


print("== кәдімгі жүріс")
code = main([])
saved = REPORT.read_text(encoding="utf-8")
print("  код:", code)

print()
print("== құрғақ жүріс")
code = main(["--dry-run"])
print("  код:", code, "| файл өзгерген жоқ:", REPORT.read_text(encoding="utf-8") == saved)

print()
print("== дереккөз құлады")
code = main(["--broken"])
print("  код:", code, "| кешегі есеп орнында:", REPORT.read_text(encoding="utf-8") == saved)

print()
print("есеп:")
print(saved, end="")
sys.exit(0)
