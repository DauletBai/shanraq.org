"""The answer to lesson 11's exercise: a file beside the program, and a report.

Every path is built from the program's own folder rather than the folder it was
started from, both files name their encoding, and the practice files are removed
at the end -- a program that makes a file should know how to take it away.
"""

from pathlib import Path

HERE = Path(__file__).parent
data = HERE / "ceny.csv"
report = HERE / "otchet.txt"

rows = ["January;520", "February;546,5", "March;no price", "April;498"]
data.write_text("\n".join(rows) + "\n", encoding="utf-8")
print(f"written: {data.name}, {data.stat().st_size} bytes")

good = []
skipped = []
with data.open(encoding="utf-8") as source:
    for line in source:
        month, _, text = line.strip().partition(";")
        try:
            good.append((month, float(text.replace(",", "."))))
        except ValueError:
            skipped.append(month)

lines = []
for month, value in good:
    lines.append(f"{month}: {value:.2f}")
report.write_text("\n".join(lines) + "\n", encoding="utf-8")
print(report.read_text(encoding="utf-8"), end="")
print("skipped:", ", ".join(skipped))

data.unlink()
report.unlink()
print("cleared, files beside it:", len(list(HERE.glob("ceny.csv"))))
