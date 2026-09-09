"""The answer to lesson 17's exercise: one program laid out in files.

The module and the package are written by the program itself, so that the answer
can be run with one command; in real work a person creates those files once and
the program only imports them. The point is what the imports show: a name for
every piece, and __name__ telling a module apart from the file that was started.
"""

import importlib
import shutil
from pathlib import Path

HERE = Path(__file__).parent

(HERE / "sana.py").write_text('''"""Расчёты: среднее и годы выше предела."""


def average(values):
    """Среднее по ряду; пропуски не в счёт."""
    numbers = [value for value in values if value is not None]
    return sum(numbers) / len(numbers)


def above(series, limit):
    """Годы, где значение больше предела."""
    return [year for year, value in series.items()
            if value is not None and value > limit]
''', encoding="utf-8")

package = HERE / "svodka"
package.mkdir(exist_ok=True)
(package / "__init__.py").write_text('"""Пакет сводки."""\n', encoding="utf-8")
(package / "otchet.py").write_text('''"""Печать отчёта."""


def title(text):
    """Заголовок отчёта."""
    return text.upper()


def line(year, value):
    """Одна строка отчёта."""
    return f"{year}: {value:.1f}%"
''', encoding="utf-8")

importlib.invalidate_caches()

import sana
from svodka import otchet

series = {2023: 14.5, 2024: 8.7, 2025: 11.4, 2026: None}

print(otchet.title("инфляция"))
for year, value in series.items():
    if value is None:
        print(f"{year}: данных нет")
        continue
    print(otchet.line(year, value))
print(f"среднее: {sana.average(series.values()):.2f}")
print("выше десяти:", sana.above(series, 10))
print("модуль:", sana.__name__, "| эта программа:", __name__)

for name in ("sana.py",):
    (HERE / name).unlink()
shutil.rmtree(package)
shutil.rmtree(HERE / "__pycache__", ignore_errors=True)
