#!/usr/bin/env python3
"""Check that every Go program printed in a lesson still parses.

A lesson's code is assembled by scripts, and a script that rewrites a block can
corrupt it silently -- re.sub, for one, reads backslash sequences in its
replacement, which turned `"\\""` into something that no longer compiled. The
lesson still rendered and still looked right; only the reader who typed it in
would have found out.

So every fenced go block that starts with `package` is handed to gofmt -e,
which parses it and says where it breaks. Fragments -- a function or two out of
context -- are skipped: they are not meant to parse on their own.

    python3 tools/coursecheck/gocheck.py <файлы>...
"""

import os
import re
import subprocess
import sys
import tempfile

BLOCK = re.compile(r"```go\n(.*?)\n```", re.S)


def main(argv):
    paths = argv[1:]
    if not paths:
        print(__doc__.strip())
        return 2

    checked = bad = 0
    for path in paths:
        with open(path, encoding="utf-8") as f:
            text = f.read()
        for i, block in enumerate(BLOCK.findall(text)):
            if not block.lstrip().startswith("package "):
                continue  # фрагмент, а не программа целиком
            checked += 1
            tmp = tempfile.NamedTemporaryFile(
                "w", suffix=".go", delete=False, encoding="utf-8")
            tmp.write(block + "\n")
            tmp.close()
            r = subprocess.run(["gofmt", "-e", tmp.name],
                               capture_output=True, text=True)
            os.unlink(tmp.name)
            if r.returncode != 0:
                bad += 1
                first = (r.stderr.strip().splitlines() or ["gofmt молчит"])[0]
                print(f"{os.path.basename(path)}: программа №{i + 1} не разбирается")
                print(f"    {first.split(':', 1)[-1].strip()}")

    print(f"проверено программ: {checked}, сломанных: {bad}")
    return 1 if bad else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
