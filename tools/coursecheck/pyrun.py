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


def blocks(text):
    """Yield (program, printed output) for every whole program in a lesson."""
    found = FENCE.findall(text)
    for i, (lang, body) in enumerate(found):
        if lang != "python" or not body.lstrip().startswith(WHOLE):
            continue
        # The output is the first plain fence after the program: the lesson's
        # own layout, "here is the program, here is what it prints".
        for lang2, out in found[i + 1:]:
            if lang2 == "python":
                break
            yield body, out
            break


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
