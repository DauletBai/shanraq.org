#!/usr/bin/env python3
"""Send the lesson files in course/lessons to the live site.

The site's database is where readers see a lesson; course/lessons is where it is
written and reviewed. This keeps the two the same, in that direction only: it
updates the body and the summary of lessons that already exist, and it refuses
to create anything. Titles are left alone -- the site's headline carries a
subtitle the source file does not.

    python3 tools/course/publish.py --check     # list what differs
    python3 tools/course/publish.py             # send the differences
    python3 tools/course/publish.py go/http.md ...          # only these

The host comes from SHANRAQ_HOST (default: the production host in docs/BACKUPS.md)
and every statement runs through the same psql the deploy runbook uses. A file
is compared by the sha256 of its body, so a run that changes nothing sends
nothing, and after sending, every hash is read back and compared again.
"""

import argparse
import hashlib
import io
import json
import os
import re
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course" / "lessons"
SLUGS = Path(__file__).resolve().parent / "lesson-slugs.json"
HOST = os.environ.get("SHANRAQ_HOST", "root@85.202.192.61")

PSQL = ("cd /opt/shanraq && docker compose -f docker-compose.prod.yml exec -T db "
        "psql -U shanraq -d shanraq -v ON_ERROR_STOP=1 -q")

LEAD = re.compile(r"_[^_]+:_\s*\*\*(.+)\*\*\s*$")


def parse(path):
    """Split a lesson file into its headline, its summary and its body."""
    lines = io.open(path, encoding="utf-8").read().split("\n")
    if not lines[0].startswith("# "):
        raise ValueError(f"{path}: the first line is not a headline")
    m = LEAD.match(lines[2].strip())
    if not m:
        raise ValueError(f"{path}: the third line is not a summary")
    return lines[0][2:].strip(), m.group(1).strip(), "\n".join(lines[4:]).strip() + "\n"


def lesson_key(arg, slugs):
    """Turn whatever was typed on the command line into a key of the map."""
    path = Path(arg)
    stem = path.name[:-3] if path.name.endswith(".md") else path.name
    if path.parent.name and f"{path.parent.name}/{stem}" in slugs:
        return f"{path.parent.name}/{stem}"
    matches = [k for k in slugs if k.split("/")[-1] == stem]
    return matches[0] if len(matches) == 1 else stem


def sha(text):
    return hashlib.sha256(text.encode()).hexdigest()


def remote(sql_or_query, stdin=None):
    """Run psql on the server. Query mode returns the rows, script mode the log."""
    if stdin is None:
        cmd = ["ssh", HOST, PSQL + ' -tAc "' + sql_or_query + '"']
        return subprocess.run(cmd, capture_output=True, text=True, check=True).stdout
    cmd = ["ssh", HOST, PSQL]
    r = subprocess.run(cmd, input=stdin, capture_output=True, text=True)
    if r.returncode != 0:
        sys.exit("psql: " + (r.stderr.strip() or "failed"))
    return r.stdout


def live():
    """The hash of every published body and summary, by slug and language."""
    rows = remote(
        "select a.slug || '|' || t.lang || '|' "
        "|| encode(sha256(convert_to(t.body_md,'UTF8')),'hex') || '|' "
        "|| encode(sha256(convert_to(t.summary,'UTF8')),'hex') "
        "from article_translations t join articles a on a.id = t.article_id;")
    out = {}
    for line in rows.splitlines():
        if line.strip():
            slug, lang, body, summary = line.split("|")
            out[(slug, lang)] = (body, summary)
    return out


def main(argv):
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--check", action="store_true", help="list the differences and stop")
    ap.add_argument("files", nargs="*", help="lesson files (default: all of them)")
    args = ap.parse_args(argv[1:])

    slugs = json.load(io.open(SLUGS, encoding="utf-8"))
    # A file may be named as it is passed on the command line -- with or
    # without the course folder, with or without the leading path -- and the
    # map is keyed by "<course>/<name>".
    names = [lesson_key(f, slugs) for f in args.files] if args.files else sorted(slugs)
    published = live()

    todo, sql = [], ["BEGIN;"]
    for name in names:
        if name not in slugs:
            print(f"  ? {name}: not in lesson-slugs.json — skipped")
            continue
        slug, lang = slugs[name]
        path = LESSONS / (name + ".md")
        if not path.exists():
            print(f"  ? {name}: no such file — skipped")
            continue
        try:
            _, summary, body = parse(path)
        except ValueError as e:
            print(f"  ? {e}")
            continue
        if (slug, lang) not in published:
            print(f"  ? {name} -> {slug}/{lang}: not published — skipped")
            continue
        body_live, summary_live = published[(slug, lang)]
        fields = []
        if sha(body) != body_live:
            fields.append("body")
        if sha(summary) != summary_live:
            fields.append("summary")
        if not fields:
            continue
        # $lesson$ quoting keeps the text literal; a lesson that contained the
        # marker itself would break the script, so it is checked rather than
        # escaped -- no lesson has ever needed it.
        for text in (body, summary):
            if "$lesson$" in text:
                sys.exit(f"{name}: the text contains the quoting marker")
        todo.append((name, slug, lang, fields))
        sets = []
        if "body" in fields:
            sets.append("body_md = $lesson$" + body + "$lesson$")
        if "summary" in fields:
            sets.append("summary = $lesson$" + summary + "$lesson$")
        sql.append("UPDATE article_translations t SET " + ", ".join(sets) + ", updated_at = now()\n"
                   f"FROM articles a WHERE a.id = t.article_id AND a.slug = '{slug}' "
                   f"AND t.lang = '{lang}';")
    sql.append("COMMIT;")

    for name, slug, lang, fields in todo:
        print(f"  -> {name} -> {slug}/{lang} ({', '.join(fields)})")
    if not todo:
        print("the site matches course/lessons")
        return 0
    if args.check:
        print(f"{len(todo)} lessons differ")
        return 1

    remote(None, stdin="\n".join(sql) + "\n")
    published = live()
    bad = []
    for name, slug, lang, _ in todo:
        _, summary, body = parse(LESSONS / (name + ".md"))
        if published[(slug, lang)] != (sha(body), sha(summary)):
            bad.append(name)
    print(f"sent: {len(todo)}, still different: {bad or 'none'}")
    return 1 if bad else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
