#!/usr/bin/env python3
"""Prepare one atomic SQL publication for the first mathematics route."""
import argparse
import hashlib
import json
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/mathematics"
STEMS = (
    "preface", "01-diagnostic", "02-fraction-meaning",
    "03-equivalent-fractions", "04-compare-fractions",
    "05-fraction-operations", "06-ratio", "07-percent",
    "08-proportion", "09-mastery",
)
SLUGS = ("mathematics-before-start",) + tuple("math-" + stem for stem in STEMS[1:])
COVER = "/static/covers/school/mathematics/mathematics-foundations.webp"
META = {
    "ru": (
        "Математика: от фундамента к высшей математике",
        "Бесплатный курс по карте зависимостей, а не по классам. Первый открытый маршрут из 9 занятий связывает дроби, отношения, проценты и пропорции и завершает их проверкой переноса.",
    ),
    "kz": (
        "Математика: іргетастан жоғары математикаға дейін",
        "Сыныптармен емес, ұғымдар тәуелділігінің картасымен құрылған тегін курс. 9 сабақтан тұратын алғашқы бағыт қазір орыс тілінде ашық; қазақша нұсқа редакциялық тексеруден кейін қосылады.",
    ),
    "en": (
        "Mathematics: from foundations to higher mathematics",
        "A free course organized by idea dependencies rather than grade levels. Its first 9-lesson route is currently open in Russian; reviewed English localization will follow.",
    ),
}
LEAD = re.compile(r"_[^_]+:_\s*\*\*(.+)\*\*\s*$")


def literal(value: str) -> str:
    delimiter = "$mathematics_course$"
    if delimiter in value:
        raise ValueError("SQL delimiter occurs in mathematics content")
    return delimiter + value + delimiter


def lesson(path: Path):
    lines = path.read_text(encoding="utf-8").splitlines()
    if len(lines) < 5 or not lines[0].startswith("# "):
        raise ValueError(f"missing title/body: {path}")
    lead = LEAD.fullmatch(lines[2])
    if not lead:
        raise ValueError(f"missing summary lead: {path}")
    body = "\n".join(lines[4:]).strip() + "\n"
    return lines[0][2:].strip(), lead.group(1), body


def prepare():
    if len(STEMS) != 10 or len(SLUGS) != 10:
        raise ValueError("expected one preface and nine lessons")
    sql = ["BEGIN;", "SELECT pg_advisory_xact_lock(hashtext('shanraq-mathematics-course'));" ]
    slugs = ",".join(literal(s) for s in SLUGS)
    sql.append(f"""DO $guard$
BEGIN
  IF (SELECT count(*) FROM auth_users WHERE email='baimurza.daulet@gmail.com') <> 1 THEN
    RAISE EXCEPTION 'Expected course author not found';
  END IF;
  IF EXISTS (SELECT 1 FROM articles a JOIN auth_users u ON u.id=a.author_id
             WHERE a.slug IN ({slugs}) AND u.email <> 'baimurza.daulet@gmail.com') THEN
    RAISE EXCEPTION 'Mathematics slug belongs to another author';
  END IF;
  IF EXISTS (SELECT 1 FROM article_series_items i
             JOIN article_series s ON s.id=i.series_id JOIN articles a ON a.id=i.article_id
             WHERE a.slug IN ({slugs}) AND s.slug <> 'mathematics') THEN
    RAISE EXCEPTION 'Mathematics lesson belongs to another course';
  END IF;
END $guard$;""")
    sql.append(f"""INSERT INTO article_series(slug,cover_url,status,code_lang)
VALUES('mathematics',{literal(COVER)},'published','math')
ON CONFLICT(slug) DO UPDATE SET cover_url=EXCLUDED.cover_url,status='published',
code_lang='math',updated_at=now();""")
    for lang, (title, summary) in META.items():
        sql.append(f"""INSERT INTO article_series_i18n(series_id,lang,title,summary)
SELECT id,'{lang}',{literal(title)},{literal(summary)} FROM article_series
WHERE slug='mathematics' ON CONFLICT(series_id,lang) DO UPDATE
SET title=EXCLUDED.title,summary=EXCLUDED.summary;""")
    expected = []
    for index, (stem, slug_name) in enumerate(zip(STEMS, SLUGS)):
        title, summary, body = lesson(LESSONS / f"{stem}.md")
        slug = literal(slug_name)
        sql.append(f"""INSERT INTO articles(author_id,slug,original_lang,category,subcategory,cover_url,status,published_at)
SELECT id,{slug},'ru','society','education',{literal(COVER)},'published',now()
FROM auth_users WHERE email='baimurza.daulet@gmail.com'
ON CONFLICT(slug) DO UPDATE SET category='society',subcategory='education',
cover_url=EXCLUDED.cover_url,status='published',updated_at=now(),
published_at=COALESCE(articles.published_at,now());""")
        sql.append(f"""INSERT INTO article_translations(article_id,lang,title,summary,body_md,source,status)
SELECT id,'ru',{literal(title)},{literal(summary)},{literal(body)},'ai','ready'
FROM articles WHERE slug={slug}
ON CONFLICT(article_id,lang) DO UPDATE SET title=EXCLUDED.title,summary=EXCLUDED.summary,
body_md=EXCLUDED.body_md,source='ai',status='ready',updated_at=now();""")
        expected.append({
            "slug": slug_name,
            "lang": "ru",
            "title": title,
            "body_sha256": hashlib.sha256(body.encode()).hexdigest(),
            "summary_sha256": hashlib.sha256(summary.encode()).hexdigest(),
        })
        position = 1 if index == 0 else index * 10
        sql.append(f"""INSERT INTO article_series_items(series_id,article_id,position)
SELECT s.id,a.id,{position} FROM article_series s,articles a
WHERE s.slug='mathematics' AND a.slug={slug}
ON CONFLICT(series_id,article_id) DO UPDATE SET position=EXCLUDED.position;""")
    sql.append("COMMIT;")
    return "\n".join(sql) + "\n", expected


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--sql", type=Path, required=True)
    parser.add_argument("--expected", type=Path, required=True)
    args = parser.parse_args()
    sql, expected = prepare()
    args.sql.write_text(sql, encoding="utf-8")
    args.expected.write_text(json.dumps(expected, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(f"Prepared {len(expected)} Russian pages in one transaction; no remote changes")


if __name__ == "__main__":
    main()
