#!/usr/bin/env python3
"""check_language.py — fail the build when Russian prose reaches the code.

The project is read on GitHub by people who do not read Russian, so its
comments, its test failure messages and its commit subjects are English. That
rule was agreed repeatedly and broken repeatedly, because it lived only in
somebody's memory. It lives here now.

What counts as a violation is narrow on purpose. Cyrillic in this codebase is
usually legitimate and must survive:

  * strings shown to a reader — the i18n dictionary, the Telegram bot's
    replies, the text of e-mails;
  * data under test — the bank's own "ДОЛЛАР США" in the parser fixtures, the
    Kazakh names in the name tests;
  * examples quoted inside an English sentence — a comment explaining that the
    parser reads Russian "стоит" as "сто" has to print both.

So the test is not "does this line contain Cyrillic". It is "does the prose
itself, once every quoted fragment is removed, contain Cyrillic". Only two
kinds of prose are checked: comments, and the message a failing test prints.

A line that genuinely needs Russian prose can say so with a `lang:ok` marker on
the line or the line above it. The marker is meant to be rare and to be argued
for in review; it exists so the check can be satisfied honestly rather than
switched off.

Usage:
    scripts/check_language.py                 # check the working tree
    scripts/check_language.py --commits A..B  # also check commit subjects
    scripts/check_language.py --commits A..B --subjects-only
"""

import argparse
import re
import subprocess
import sys

CYRILLIC = re.compile(r"[Ѐ-ӿ]")

# Quoted fragments are data, not prose: a Go string literal, a raw literal, and
# the guillemets used when an English sentence cites a Russian word.
QUOTED = re.compile(r'"(?:[^"\\]|\\.)*"' r"|`[^`]*`" r"|«[^»]*»")

TEST_CALL = re.compile(r"\bt\.(?:Errorf|Fatalf|Error|Fatal|Logf|Log|Skipf)\s*\(")

ALLOW = "lang:ok"


def prose_only(line: str) -> str:
    """Everything outside quotes: the sentence the author wrote themselves."""
    return QUOTED.sub(" ", line)


def format_string(line: str) -> str:
    """The message a failing test prints — the first string of the t.* call.

    Later arguments carry the values under test and may be Cyrillic by nature;
    only the sentence around them is ours to write.
    """
    m = TEST_CALL.search(line)
    if not m:
        return ""
    rest = line[m.end():]
    lit = re.match(r'\s*"((?:[^"\\]|\\.)*)"', rest)
    return lit.group(1) if lit else ""


def go_files() -> list[str]:
    out = subprocess.run(
        ["git", "ls-files", "*.go"], capture_output=True, text=True, check=True
    ).stdout.split()
    return [f for f in out if not f.endswith("_gen.go")]


def scan_file(path: str) -> list[tuple[int, str, str]]:
    found = []
    try:
        lines = open(path, encoding="utf-8").read().splitlines()
    except OSError:
        return found
    for i, line in enumerate(lines, 1):
        if not CYRILLIC.search(line):
            continue
        if ALLOW in line or (i > 1 and ALLOW in lines[i - 2]):
            continue
        stripped = line.strip()
        if stripped.startswith("//") and CYRILLIC.search(prose_only(line)):
            found.append((i, "comment", stripped))
            continue
        msg = format_string(line)
        # A cited fragment inside the message is data; the rest is the sentence.
        if msg and CYRILLIC.search(QUOTED.sub(" ", msg.replace('\\"', '"'))):
            found.append((i, "test message", stripped))
    return found


def check_commits(rng: str) -> list[str]:
    try:
        out = subprocess.run(
            ["git", "log", "--format=%h %s", rng],
            capture_output=True, text=True, check=True,
        ).stdout.splitlines()
    except subprocess.CalledProcessError:
        # A shallow clone or a first push has no range to walk. Nothing to say.
        return []
    return [c for c in out if CYRILLIC.search(c.split(" ", 1)[-1])]


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--commits", help="commit range to check subjects of, e.g. A..B")
    ap.add_argument(
        "--subjects-only",
        action="store_true",
        help="check only commit subjects, leaving the code backlog for later",
    )
    args = ap.parse_args()

    bad = {} if args.subjects_only else {f: hits for f in go_files() if (hits := scan_file(f))}
    commits = check_commits(args.commits) if args.commits else []

    if not bad and not commits:
        if args.subjects_only:
            print("language check: commit subjects are English")
        else:
            print("language check: comments, test messages and commit subjects are English")
        return 0

    if bad:
        total = sum(len(v) for v in bad.values())
        print(f"Russian prose in {total} line(s) across {len(bad)} file(s).")
        print("Comments and test failure messages are English; strings shown to a")
        print(f"reader and data under test are not. Mark a real exception `{ALLOW}`.\n")
        for f, hits in sorted(bad.items()):
            for line_no, kind, text in hits:
                print(f"  {f}:{line_no}  [{kind}] {text[:100]}")
        print()

    if commits:
        print(f"Russian in {len(commits)} commit subject(s):")
        for c in commits:
            print(f"  {c[:100]}")
        print()

    return 1


if __name__ == "__main__":
    sys.exit(main())
