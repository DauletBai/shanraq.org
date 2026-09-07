#!/usr/bin/env python3
"""Compile every Python program printed in a lesson.

The Go course has gocheck.py for the same reason: a lesson's code is not what
the author meant to write but what the file says, and a regular expression that
edited a text once corrupted a program that still looked right on the page. This
one hands every fenced ```python block that starts a program to the compiler.

Only blocks that look like a whole program are checked: a fragment showing three
lines of a function has no imports and would fail for the wrong reason.
"""

import py_compile
import re
import sys
import tempfile
from pathlib import Path

BLOCK = re.compile(r"```python\n(.*?)```", re.S)
# A whole program starts with a docstring, an import or a shebang -- a fragment
# usually starts in the middle of something.
WHOLE = re.compile(r'\A\s*(?:"""|#!|import |from )')


def check(path):
    text = Path(path).read_text(encoding="utf-8")
    checked = broken = 0
    for block in BLOCK.findall(text):
        if not WHOLE.match(block):
            continue
        checked += 1
        with tempfile.NamedTemporaryFile("w", suffix=".py", delete=False, encoding="utf-8") as f:
            f.write(block)
            name = f.name
        try:
            py_compile.compile(name, doraise=True, cfile=tempfile.mktemp())
        except py_compile.PyCompileError as err:
            broken += 1
            print(f"{path}: программа не компилируется\n  {err.msg.strip()}")
    return checked, broken


def main(paths):
    total = failed = 0
    for path in paths:
        checked, broken = check(path)
        total += checked
        failed += broken
    print(f"проверено программ: {total}, сломанных: {failed}")
    return 1 if failed else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
