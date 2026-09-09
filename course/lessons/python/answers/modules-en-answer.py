"""The answer to lesson 17's exercise, English: one program laid out in files.

The module and the package are written by the program itself, so that the answer
can be run with one command; in real work a person creates those files once and
the program only imports them. The point is what the imports show: a name for
every piece, and __name__ telling a module apart from the file that was started.
"""

import importlib
import shutil
from pathlib import Path

HERE = Path(__file__).parent

(HERE / "sana.py").write_text('''"""The arithmetic: the average and the years above a limit."""


def average(values):
    """The average over a series; gaps do not count."""
    numbers = [value for value in values if value is not None]
    return sum(numbers) / len(numbers)


def above(series, limit):
    """The years whose value is above the limit."""
    return [year for year, value in series.items()
            if value is not None and value > limit]
''', encoding="utf-8")

package = HERE / "digest"
package.mkdir(exist_ok=True)
(package / "__init__.py").write_text('"""The digest package."""\n', encoding="utf-8")
(package / "report.py").write_text('''"""Printing the report."""


def title(text):
    """The report's heading."""
    return text.upper()


def line(year, value):
    """One line of the report."""
    return f"{year}: {value:.1f}%"
''', encoding="utf-8")

importlib.invalidate_caches()

import sana
from digest import report

series = {2023: 14.5, 2024: 8.7, 2025: 11.4, 2026: None}

print(report.title("inflation"))
for year, value in series.items():
    if value is None:
        print(f"{year}: no data")
        continue
    print(report.line(year, value))
print(f"average: {sana.average(series.values()):.2f}")
print("above ten:", sana.above(series, 10))
print("the module:", sana.__name__, "| this program:", __name__)

for name in ("sana.py",):
    (HERE / name).unlink()
shutil.rmtree(package)
shutil.rmtree(HERE / "__pycache__", ignore_errors=True)
