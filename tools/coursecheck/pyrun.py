#!/usr/bin/env python3
"""Run every program a Python lesson prints and compare it with the shown output.

pycheck.py proves that a program compiles, and a lesson promises more than that:
this program, typed in exactly as printed, prints exactly this. The two are not
the same check, and the difference has already cost a real error -- a ratio
computed correctly and described backwards, which every existing check passed.

    python3 tools/coursecheck/pyrun.py course/lessons/python/*.md
    python3 tools/coursecheck/pyrun.py --steps

A program that needs the network is run too. If the network is unreachable the
program is reported as skipped rather than failed: a course must not go red
because somebody else's server is down.

The warm-up drills are checked the same way. A drill shows a program in one
place and its output in another -- the reader is asked to predict it, or to fix
the program first -- so the two are tied together by a marker instead of by
being next to each other:

    <!-- drill 1 -->        before the program that is meant to run
    <!-- drill 1 out -->    before the output it is meant to print

The exercise is held to the same standard. A lesson sets its required task on
fixed data and prints the output the reader is working towards; that output is
a promise, and a promise nobody can keep is the worst kind of error a course can
make. So the task's own solution lives beside the lessons, out of the reader's
way, and the printed result is compared with what it actually prints:

    course/lessons/python/answers/<lesson>-answer.py   the reference solution
    <!-- task out -->                             before the promised output

Both markers are HTML comments, so a reader never sees them, and the answer at
the end of the lesson is held to what the machine actually prints.
"""

import os
import re
import subprocess
import sys
import tempfile

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
FENCE = re.compile(r"```(\w*)\n(.*?)```", re.S)
# A whole program starts with a docstring, an import or a shebang; a fragment of
# three lines out of a function has no imports and would fail for its own reason.
WHOLE = ("\"\"\"", "'''", "import ", "from ", "#!/")
NETWORK = ("URLError", "urlopen error", "Temporary failure in name resolution",
           "Network is unreachable", "timed out", "getaddrinfo failed")

# From the pandas module on, a lesson needs a library that is not part of Python.
# A machine without it is a machine that is not set up, and that is not the same
# thing as a broken lesson -- the same rule the network already gets. A missing
# standard-library module stays an error, because that is a typo in the lesson.
NO_MODULE = re.compile(r"ModuleNotFoundError: No module named '([\w.]+)'")


def missing_library(err):
    """The uninstalled library a program needed, or "" when that was not why."""
    m = NO_MODULE.search(err)
    if not m:
        return ""
    name = m.group(1).split(".")[0]
    return "" if name in sys.stdlib_module_names else name


DRILL = re.compile(r"<!--\s*drill (\d+)(?:\s+(out))?\s*-->")


def drills(text):
    """Yield (number, program, expected output) for the lesson's warm-up drills.

    A drill's program and its output live apart -- the whole point is that the
    reader answers before they see the answer -- so each is named by a marker
    and the pair is put back together here.
    """
    fences = list(FENCE.finditer(text))

    def after(pos):
        """The first fence that starts after this marker."""
        for f in fences:
            if f.start() >= pos:
                return f
        return None

    program, printed = {}, {}
    for m in DRILL.finditer(text):
        f = after(m.end())
        if f is None:
            continue
        (printed if m.group(2) else program)[m.group(1)] = f
    for number in sorted(program, key=int):
        if number in printed:
            yield number, program[number].group(2), printed[number].group(2)


def drill_bodies(text):
    """The fences that belong to a drill, so the plain scan leaves them alone."""
    return {f.group(2) for _, f, _ in _drill_fences(text)}


def _drill_fences(text):
    fences = list(FENCE.finditer(text))
    for m in DRILL.finditer(text):
        for f in fences:
            if f.start() >= m.end():
                yield m.group(1), f, bool(m.group(2))
                break


def blocks(text):
    """Yield (program, printed output) for every whole program in a lesson."""
    found = FENCE.findall(text)
    taken = drill_bodies(text)
    for i, (lang, body) in enumerate(found):
        if lang != "python" or not body.lstrip().startswith(WHOLE):
            continue
        if body in taken:
            continue  # a drill: checked by its marker, not by what follows it
        # The output is the first plain fence after the program: the lesson's
        # own layout, "here is the program, here is what it prints".
        for lang2, out in found[i + 1:]:
            if lang2 == "python":
                break
            yield body, out
            break


TASK_OUT = re.compile(r"<!--\s*task out\s*-->")


def task(path, text):
    """Return (reference solution, promised output) for the lesson's exercise.

    Both halves are optional: a lesson without a fixed-data task has neither,
    and a lesson that has one without a solution beside it is reported, because
    that is how a promise goes unchecked.
    """
    m = TASK_OUT.search(text)
    if not m:
        return None, None
    printed = None
    for f in FENCE.finditer(text):
        if f.start() >= m.end():
            printed = f.group(2)
            break
    # One solution per lesson file, not per lesson: the printed result is in
    # the lesson's own language, and that is exactly where a translation drifts.
    # The "-answer" is not decoration: a file called csv.py beside a program is
    # what "import csv" finds, and the solution would break on its own name.
    name = os.path.basename(path).removesuffix(".md")
    solution = os.path.join(os.path.dirname(path), "answers", name + "-answer.py")
    if not os.path.isfile(solution):
        return "", printed
    with open(solution, encoding="utf-8") as f:
        return f.read(), printed


def normalise(s):
    return "\n".join(line.rstrip() for line in s.strip().split("\n"))


def run(program):
    """Return (stdout, error) with error set when the program did not finish."""
    with tempfile.TemporaryDirectory() as tmp:
        path = os.path.join(tmp, "lesson.py")
        with open(path, "w", encoding="utf-8") as f:
            f.write(program)
        env = dict(os.environ, PYTHONIOENCODING="utf-8", PYTHONHASHSEED="0")
        try:
            r = subprocess.run([sys.executable, path], cwd=tmp, env=env,
                               capture_output=True, text=True, timeout=180)
        except subprocess.TimeoutExpired:
            return "", "не уложилась в 180 секунд"
        if r.returncode != 0:
            return r.stdout, r.stderr.strip().split("\n")[-1]
        return r.stdout, ""


def check_lessons(paths):
    ran = skipped = bad = 0
    for path in paths:
        with open(path, encoding="utf-8") as f:
            text = f.read()
        for n, (program, printed) in enumerate(blocks(text), 1):
            out, err = run(program)
            if err and any(mark in err for mark in NETWORK):
                skipped += 1
                print(f"  ~ {os.path.basename(path)} #{n}: пропущено, сеть недоступна")
                continue
            if err and missing_library(err):
                skipped += 1
                print(f"  ~ {os.path.basename(path)} #{n}: пропущено, "
                      f"не установлена библиотека {missing_library(err)}")
                continue
            if err:
                bad += 1
                print(f"  ! {os.path.basename(path)} #{n}: {err}")
                continue
            ran += 1
            if normalise(out) != normalise(printed):
                bad += 1
                print(f"  ! {os.path.basename(path)} #{n}: вывод разошёлся с уроком")
                for line in diff(normalise(printed), normalise(out)):
                    print("      " + line)
        solution, promised = task(path, text)
        if promised is not None:
            if not solution:
                bad += 1
                print(f"  ! {os.path.basename(path)} задание: нет решения в answers/, "
                      f"обещанный вывод никем не проверен")
            else:
                out, err = run(solution)
                if err and any(mark in err for mark in NETWORK):
                    skipped += 1
                    print(f"  ~ {os.path.basename(path)} задание: пропущено, сеть недоступна")
                elif err and missing_library(err):
                    skipped += 1
                    print(f"  ~ {os.path.basename(path)} задание: пропущено, "
                          f"не установлена библиотека {missing_library(err)}")
                elif err:
                    bad += 1
                    print(f"  ! {os.path.basename(path)} задание: {err}")
                else:
                    ran += 1
                    if normalise(out) != normalise(promised):
                        bad += 1
                        print(f"  ! {os.path.basename(path)} задание: обещанный вывод "
                              f"разошёлся с решением")
                        for line in diff(normalise(promised), normalise(out)):
                            print("      " + line)
        for number, program, printed in drills(text):
            out, err = run(program)
            if err and any(mark in err for mark in NETWORK):
                skipped += 1
                print(f"  ~ {os.path.basename(path)} разминка {number}: пропущена, сеть недоступна")
                continue
            if err and missing_library(err):
                skipped += 1
                print(f"  ~ {os.path.basename(path)} разминка {number}: пропущена, "
                      f"не установлена библиотека {missing_library(err)}")
                continue
            if err:
                bad += 1
                print(f"  ! {os.path.basename(path)} разминка {number}: {err}")
                continue
            ran += 1
            if normalise(out) != normalise(printed):
                bad += 1
                print(f"  ! {os.path.basename(path)} разминка {number}: ответ разошёлся с выводом")
                for line in diff(normalise(printed), normalise(out)):
                    print("      " + line)
    print(f"выполнено программ: {ran}, пропущено: {skipped}, разошлось: {bad}")
    return 1 if bad else 0


def diff(want, got):
    """The first few lines where the two differ, as the lesson would show them."""
    a, b = want.split("\n"), got.split("\n")
    out = []
    for i in range(max(len(a), len(b))):
        x = a[i] if i < len(a) else "(нет строки)"
        y = b[i] if i < len(b) else "(нет строки)"
        if x != y:
            out.append(f"в уроке: {x}")
            out.append(f"на деле: {y}")
            if len(out) >= 6:
                break
    return out


def check_steps():
    """Run every step of the course project the way its README says to."""
    root = os.path.join(ROOT, "course", "py-digest")
    if not os.path.isdir(root):
        print("шагов проекта пока нет")
        return 0
    ran = skipped = bad = 0
    for name in sorted(os.listdir(root)):
        step = os.path.join(root, name)
        main = os.path.join(step, "main.py")
        if not os.path.isfile(main):
            continue
        env = dict(os.environ, PYTHONIOENCODING="utf-8")
        try:
            r = subprocess.run([sys.executable, "main.py"], cwd=step, env=env,
                               capture_output=True, text=True, timeout=180)
        except subprocess.TimeoutExpired:
            bad += 1
            print(f"  ! {name}: не уложился в 180 секунд")
            continue
        if r.returncode != 0:
            last = r.stderr.strip().split("\n")[-1]
            if any(mark in r.stderr for mark in NETWORK):
                skipped += 1
                print(f"  ~ {name}: пропущен, сеть недоступна")
            elif missing_library(r.stderr):
                skipped += 1
                print(f"  ~ {name}: пропущен, не установлена библиотека "
                      f"{missing_library(r.stderr)}")
            else:
                bad += 1
                print(f"  ! {name}: {last}")
            continue
        ran += 1
    print(f"шагов проекта выполнено: {ran}, пропущено: {skipped}, сломанных: {bad}")
    return 1 if bad else 0


def main(argv):
    args = argv[1:]
    if not args:
        print(__doc__.strip())
        return 2
    if args[0] == "--steps":
        return check_steps()
    return check_lessons(args)


if __name__ == "__main__":
    sys.exit(main(sys.argv))
