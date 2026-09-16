#!/usr/bin/env python3
"""Create or update one trilingual article draft on the live site.

The Markdown sources are the reviewable copy in the repository; PostgreSQL is
the copy shown in the author's studio.  This command never publishes and
refuses to modify a slug that is already live.

    python3 tools/article/draft.py \
      --slug python-course-money-mechanics \
      --author baimurza.daulet@gmail.com \
      --cover /static/covers/it/python-course-money-mechanics.webp \
      content/drafts/python-course-money-mechanics-{ru,kz,en}.md
"""

import argparse
import io
import os
import re
import subprocess
import sys
import uuid
from pathlib import Path


HOST = os.environ.get("SHANRAQ_HOST", "root@85.202.192.61")
PSQL = ("cd /opt/shanraq && docker compose -f docker-compose.prod.yml exec -T db "
        "psql -U shanraq -d shanraq -v ON_ERROR_STOP=1 -q")
LEAD = re.compile(r"_[^_]+:_\s*\*\*(.+)\*\*\s*$")


def parse(path: Path):
    lines = io.open(path, encoding="utf-8").read().split("\n")
    if len(lines) < 5 or not lines[0].startswith("# "):
        raise ValueError(f"{path}: first line must be an H1")
    match = LEAD.match(lines[2].strip())
    if not match:
        raise ValueError(f"{path}: third line must be a bold lead")
    suffix = path.stem.rsplit("-", 1)[-1]
    if suffix not in {"ru", "kz", "en"}:
        raise ValueError(f"{path}: filename must end in -ru.md, -kz.md or -en.md")
    return suffix, lines[0][2:].strip(), match.group(1).strip(), "\n".join(lines[4:]).strip() + "\n"


def quote(text: str) -> str:
    marker = "$shanraq_draft$"
    if marker in text:
        raise ValueError("article contains the SQL quoting marker")
    return marker + text + marker


def remote(sql: str) -> str:
    result = subprocess.run(["ssh", HOST, PSQL], input=sql, text=True, capture_output=True)
    if result.returncode:
        raise RuntimeError(result.stderr.strip() or "remote psql failed")
    return result.stdout


def main(argv):
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--slug", required=True)
    parser.add_argument("--author", required=True)
    parser.add_argument("--cover", required=True)
    parser.add_argument("--category", default="it")
    parser.add_argument("--subcategory", default="")
    parser.add_argument("files", nargs=3)
    args = parser.parse_args(argv[1:])

    translations = dict((row[0], row[1:]) for row in map(lambda p: parse(Path(p)), args.files))
    if set(translations) != {"ru", "kz", "en"}:
        parser.error("provide exactly one RU, KZ and EN source")

    article_id = str(uuid.uuid4())
    statements = ["BEGIN;"]
    statements.append(
        "DO $$ BEGIN "
        f"IF NOT EXISTS (SELECT 1 FROM auth_users WHERE lower(email)=lower({quote(args.author)})) THEN "
        "RAISE EXCEPTION 'author not found'; END IF; END $$;"
    )
    statements.append(
        "DO $$ BEGIN "
        f"IF EXISTS (SELECT 1 FROM articles WHERE slug={quote(args.slug)} AND status <> 'draft') THEN "
        "RAISE EXCEPTION 'refusing to replace a non-draft article'; END IF; END $$;"
    )
    statements.append(
        "INSERT INTO articles (id, author_id, slug, original_lang, category, subcategory, cover_url, status) "
        f"SELECT '{article_id}'::uuid, id, {quote(args.slug)}, 'ru', {quote(args.category)}, {quote(args.subcategory)}, {quote(args.cover)}, 'draft' "
        f"FROM auth_users WHERE lower(email)=lower({quote(args.author)}) "
        "ON CONFLICT (slug) DO UPDATE SET original_lang='ru', category=EXCLUDED.category, "
        "subcategory=EXCLUDED.subcategory, cover_url=EXCLUDED.cover_url, updated_at=NOW();"
    )
    for lang in ("ru", "kz", "en"):
        title, summary, body = translations[lang]
        statements.append(
            "INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status) "
            f"SELECT id, '{lang}', {quote(title)}, {quote(summary)}, {quote(body)}, 'human', 'ready' "
            f"FROM articles WHERE slug={quote(args.slug)} "
            "ON CONFLICT (article_id, lang) DO UPDATE SET title=EXCLUDED.title, summary=EXCLUDED.summary, "
            "body_md=EXCLUDED.body_md, source='human', status='ready', updated_at=NOW();"
        )
    statements.append("COMMIT;")
    statements.append(
        "SELECT a.id || '|' || a.status || '|' || count(t.lang) "
        "FROM articles a JOIN article_translations t ON t.article_id=a.id "
        f"WHERE a.slug={quote(args.slug)} GROUP BY a.id, a.status;"
    )
    print(remote("\n".join(statements)))
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
