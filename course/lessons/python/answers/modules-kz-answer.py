"""The answer to lesson 17's exercise, Kazakh: one program laid out in files.

The module and the package are written by the program itself, so that the answer
can be run with one command; in real work a person creates those files once and
the program only imports them. The point is what the imports show: a name for
every piece, and __name__ telling a module apart from the file that was started.
"""

import importlib
import shutil
from pathlib import Path

HERE = Path(__file__).parent

(HERE / "sana.py").write_text('''"""Есептеулер: орташа және шектен жоғары жылдар."""


def average(values):
    """Қатар бойынша орташа; олқылықтар есепке кірмейді."""
    numbers = [value for value in values if value is not None]
    return sum(numbers) / len(numbers)


def above(series, limit):
    """Мәні шектен үлкен жылдар."""
    return [year for year, value in series.items()
            if value is not None and value > limit]
''', encoding="utf-8")

package = HERE / "sholu"
package.mkdir(exist_ok=True)
(package / "__init__.py").write_text('"""Шолу пакеті."""\n', encoding="utf-8")
(package / "esep.py").write_text('''"""Есепті басып шығару."""


def title(text):
    """Есептің тақырыбы."""
    return text.upper()


def line(year, value):
    """Есептің бір жолы."""
    return f"{year}: {value:.1f}%"
''', encoding="utf-8")

importlib.invalidate_caches()

import sana
from sholu import esep

series = {2023: 14.5, 2024: 8.7, 2025: 11.4, 2026: None}

print(esep.title("инфляция"))
for year, value in series.items():
    if value is None:
        print(f"{year}: дерек жоқ")
        continue
    print(esep.line(year, value))
print(f"орташа: {sana.average(series.values()):.2f}")
print("оннан жоғары:", sana.above(series, 10))
print("модуль:", sana.__name__, "| бұл бағдарлама:", __name__)

for name in ("sana.py",):
    (HERE / name).unlink()
shutil.rmtree(package)
shutil.rmtree(HERE / "__pycache__", ignore_errors=True)
