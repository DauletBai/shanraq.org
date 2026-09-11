#!/usr/bin/env python3
"""Check the checkers.

The lesson checks are the only thing standing between a reader and a program
that does not do what the page says. They have been wrong twice: langcheck's
fence pattern read prose as code and never opened a ```python block at all, and
pysyllabus counted an equals sign as proof that a lesson explains functions.
Both were found by a person reading the code, which is exactly the kind of luck
a course should not depend on.

    python3 tools/coursecheck/selftest.py

Every case here is a small file written to a temporary directory and handed to
the real checker, so what is tested is the tool itself rather than a copy of its
logic.
"""

import importlib.util
import io
import os
import subprocess
import sys
import tempfile

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.dirname(os.path.dirname(HERE))
failures = []


def load(name):
    """Import a checker by file, since the tools are scripts rather than a package."""
    spec = importlib.util.spec_from_file_location(name, os.path.join(HERE, name + ".py"))
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def check(name, got, want):
    if got != want:
        failures.append(f"{name}: получено {got!r}, ожидалось {want!r}")
        print(f"  ! {name}: {got!r} вместо {want!r}")
    else:
        print(f"  · {name}")


def write(directory, name, text):
    path = os.path.join(directory, name)
    os.makedirs(os.path.dirname(path), exist_ok=True)
    io.open(path, "w", encoding="utf-8").write(text)
    return path


def run(script, *args):
    """Run a checker as CI runs it and return (code, output)."""
    r = subprocess.run([sys.executable, os.path.join(HERE, script), *args],
                       capture_output=True, text=True)
    return r.returncode, r.stdout + r.stderr


LESSON = '''# Заголовок

_Лид (summary):_ **Лид.**

## Сразу целиком

```python
{program}
```

Выводит:

```
{output}
```

## Задание

**Обязательное.** Сделайте это.
'''


def test_langcheck():
    print("langcheck:")
    lang = load("langcheck")
    with tempfile.TemporaryDirectory() as tmp:
        # A Python fence must be read: it was invisible to the old pattern.
        p = write(tmp, "urok-kz.md", LESSON.format(program='print("нет данных")', output="нет данных"))
        words, _ = lang.check(p)
        check("русское слово в казахском коде", bool(words), True)

        # Prose between two fences is prose, not code.
        p = write(tmp, "urok2-kz.md", LESSON.format(program='print("дерек жоқ")', output="дерек жоқ")
                  + "\n\nОбычный текст между заборами не код.\n\n```\nещё блок\n```\n")
        words, _ = lang.check(p)
        check("проза не читается как код", words, [])

        # The Kazakh prose checks: a russism and an informal address.
        p = write(tmp, "urok3-kz.md", LESSON.format(program='print("дерек жоқ")', output="дерек жоқ")
                  + "\n\nБұл выгрузка туралы.\n")
        check("русизм в казахской прозе", bool(lang.check_kazakh_prose(p)), True)
        # A word that is Kazakh and still wrong here: unknown to the reader, or
        # breaking vowel harmony.
        p = write(tmp, "urok3-kz.md", LESSON.format(program='print("дерек жоқ")', output="дерек жоқ")
                  + "\n\nБөгде файл келді.\n")
        check("неудачное слово поймано", bool(lang.check_kazakh_prose(p)), True)
        p = write(tmp, "urok5-kz.md", LESSON.format(program='print("дерек жоқ")', output="дерек жоқ")
                  + "\n\nДеректер экспортының көлемі өсті.\n")
        check("похожее правильное слово не поймано", lang.check_kazakh_prose(p), [])
        p = write(tmp, "urok4-kz.md", LESSON.format(program='print("дерек жоқ")', output="дерек жоқ")
                  + "\n\nОны өзің істейсің.\n")
        check("обращение на «сен»", bool(lang.check_kazakh_prose(p)), True)
        # A Russian lesson is not held to the Kazakh rules.
        p = write(tmp, "urok.md", LESSON.format(program='print("нет данных")', output="нет данных"))
        check("русский урок не проверяется на казахские правила", lang.check_kazakh_prose(p), [])


def test_pyrun():
    print("pyrun:")
    with tempfile.TemporaryDirectory() as tmp:
        good = write(tmp, "ok.md", LESSON.format(program="import sys\n\nprint(2 + 2)", output="4"))
        code, out = run("pyrun.py", good)
        check("совпавший вывод принят", code, 0)

        bad = write(tmp, "bad.md", LESSON.format(program="import sys\n\nprint(2 + 2)", output="5"))
        code, out = run("pyrun.py", bad)
        check("разошедшийся вывод пойман", code, 1)
        check("названо расхождение", "разошёлся" in out, True)

        # A drill lives apart from its answer and is tied by markers.
        drill = write(tmp, "drill.md", LESSON.format(program="import sys\n\nprint(1)", output="1")
                      + "\n<!-- drill 1 -->\n```python\nprint(7 * 6)\n```\n\n<!-- drill 1 out -->\n```\n41\n```\n")
        code, out = run("pyrun.py", drill)
        check("разминка проверяется", code, 1)
        check("названа разминка", "разминка" in out, True)

        # The exercise's promised output is checked against the solution beside it.
        task = write(tmp, "task.md", LESSON.format(program="import sys\n\nprint(1)", output="1")
                     + "\nОжидаемый вывод:\n\n<!-- task out -->\n```\nдва\n```\n")
        code, out = run("pyrun.py", task)
        check("обещанный вывод без решения — ошибка", code, 1)
        write(tmp, "answers/task-answer.py", 'print("два")\n')
        code, out = run("pyrun.py", task)
        check("обещанный вывод сверен с решением", code, 0)

        # From the pandas module on, a lesson needs a library that is not part
        # of Python. Absent on this machine, it is a machine that is not set up;
        # absent in the standard library, it is a mistake in the lesson.
        absent = write(tmp, "lib.md", LESSON.format(
            program="import nesushchestvuyushchaya_biblioteka\n\nprint(1)", output="1"))
        code, out = run("pyrun.py", absent)
        check("неустановленная библиотека — пропуск, не ошибка", code, 0)
        check("сказано, чего не хватает", "не установлена библиотека" in out, True)

        code, out = run("pyrun.py", "--libraries-required", absent)
        check("с --libraries-required пропуск библиотеки — ошибка", code, 1)

        # A module Python itself ships is a different matter: missing, it means
        # the lesson names it wrongly, and that must stay an error.
        pyrun = load("pyrun")
        check("отсутствующая библиотека названа",
              pyrun.missing_library("ModuleNotFoundError: No module named 'pandas'"), "pandas")
        check("стандартный модуль не считается неустановленной библиотекой",
              pyrun.missing_library("ModuleNotFoundError: No module named 'sqlite3'"), "")
        check("другая ошибка библиотекой не объясняется",
              pyrun.missing_library("NameError: name 'x' is not defined"), "")


def test_pycheck():
    print("pycheck:")
    with tempfile.TemporaryDirectory() as tmp:
        broken = write(tmp, "broken.md", LESSON.format(program="import sys\n\nprint(2 +", output=""))
        code, out = run("pycheck.py", broken)
        check("непарсящаяся программа поймана", code, 1)


def test_linkcheck():
    print("linkcheck:")
    link = load("linkcheck")
    slugs = {"py-jumys-orny-venv"}
    with tempfile.TemporaryDirectory() as tmp:
        p = write(tmp, "urok.md", "Идите в [следующий урок](/read/py-jumys-orny-venv?lang=ru).\n")
        check("настоящий урок принят", link.internal_problems(p, slugs, ROOT), [])
        p = write(tmp, "urok2.md", "Идите в [никуда](/read/py-net-takogo-uroka).\n")
        check("выдуманный урок пойман", len(link.internal_problems(p, slugs, ROOT)), 1)
        p = write(tmp, "urok3.md", "![карта](/static/course/py/net-takoi-kartinki.svg)\n")
        check("пропавшая картинка поймана", len(link.internal_problems(p, slugs, ROOT)), 1)
        p = write(tmp, "urok4.md", "![карта](/static/course/py/map-prices-ru.svg)\n")
        check("существующая картинка принята", link.internal_problems(p, slugs, ROOT), [])


def test_pysyllabus():
    print("pysyllabus:")
    ps = load("pysyllabus")
    check("слабый маркер отброшен", ps.telling(["def ", "return", "="]), ["def ", "return"])
    check("строка из одних слабых маркеров сохраняется", ps.telling(["=", "+"]), ["=", "+"])
    text = "Функция объявляется словом def и возвращает значение через return."
    check("сильный маркер найден", ps.found("def ", text), True)
    check("слабый маркер не спасает", any(ps.found(m, "x = 1") for m in ps.telling(["def ", "return", "="])), False)


def main():
    for test in (test_langcheck, test_pyrun, test_pycheck, test_linkcheck, test_pysyllabus):
        test()
    if failures:
        print(f"\nсломанных проверок: {len(failures)}")
        return 1
    print("\nпроверки проверок пройдены")
    return 0


if __name__ == "__main__":
    sys.exit(main())
