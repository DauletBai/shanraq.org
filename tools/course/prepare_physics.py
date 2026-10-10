#!/usr/bin/env python3
"""Prepare the first physics block as one guarded SQL transaction; no remote write."""
import argparse
import hashlib
import json
from pathlib import Path
import re

try:
    from .build_physics_block import DATA, LESSONS, ROOT
except ImportError:  # direct invocation: python3 tools/course/prepare_physics.py
    from build_physics_block import DATA, LESSONS, ROOT

COVER = "/static/covers/school/physics/investigation/01-observe-measure.webp"
META = {
    "ru": ("Физика вокруг нас: от первого опыта до своей лаборатории",
           "Бесплатный курс физики без пропущенных шагов: наблюдаем, измеряем и проверяем объяснения. Первые восемь уроков посвящены исследованию маятника."),
    "kz": ("Айналамыздағы физика: алғашқы тәжірибеден өз зертханамызға дейін",
           "Тегін физика курсы: бақылаймыз, өлшейміз, түсіндірмені тексереміз. Алғашқы сегіз сабақта маятникті зерттейміз."),
    "en": ("Physics around us: from a first experiment to your own lab",
           "A free physics course with no skipped steps: observe, measure, and test explanations. The first eight lessons investigate a pendulum."),
}
LEAD = re.compile(r"_Лид \(summary\):_ \*\*(.+)\*\*$")


def literal(value):
    delimiter = "$physics_course$"
    if delimiter in value:
        raise ValueError("SQL delimiter occurs in lesson text")
    return delimiter + value + delimiter


def read_lesson(path):
    lines = path.read_text(encoding="utf-8").splitlines()
    if len(lines) < 6 or not lines[0].startswith("# "):
        raise ValueError(f"missing lesson title: {path}")
    lead = LEAD.fullmatch(lines[2])
    if not lead:
        raise ValueError(f"missing summary: {path}")
    return lines[0][2:].strip(), lead.group(1), "\n".join(lines[4:]).strip() + "\n"


def prepare():
    if len(DATA) != 8:
        raise ValueError("first block must contain eight lessons")
    route = [(item["stem"], f"physics-{i:02d}-{item['stem'][3:]}", i*10)
             for i, item in enumerate(DATA, 1)]
    slug_list = ",".join(literal(slug) for _, slug, _ in route)
    sql = ["BEGIN;", "SELECT pg_advisory_xact_lock(hashtext('shanraq-physics-course'));",
           f"""DO $guard$
BEGIN
  IF (SELECT count(*) FROM auth_users WHERE email='baimurza.daulet@gmail.com') <> 1 THEN
    RAISE EXCEPTION 'Expected course author not found';
  END IF;
  IF EXISTS (SELECT 1 FROM articles a JOIN auth_users u ON u.id=a.author_id
             WHERE a.slug IN ({slug_list}) AND u.email <> 'baimurza.daulet@gmail.com') THEN
    RAISE EXCEPTION 'Physics slug belongs to another author';
  END IF;
  IF EXISTS (SELECT 1 FROM article_series_items i
             JOIN article_series s ON s.id=i.series_id JOIN articles a ON a.id=i.article_id
             WHERE a.slug IN ({slug_list}) AND s.slug <> 'physics') THEN
    RAISE EXCEPTION 'Physics lesson belongs to another course';
  END IF;
END $guard$;""",
           f"""INSERT INTO article_series(slug,cover_url,status,code_lang)
VALUES('physics',{literal(COVER)},'published','math')
ON CONFLICT(slug) DO UPDATE SET cover_url=EXCLUDED.cover_url,status='published',code_lang='math',updated_at=now();"""]
    for lang, (title, summary) in META.items():
        sql.append(f"""INSERT INTO article_series_i18n(series_id,lang,title,summary)
SELECT id,'{lang}',{literal(title)},{literal(summary)} FROM article_series WHERE slug='physics'
ON CONFLICT(series_id,lang) DO UPDATE SET title=EXCLUDED.title,summary=EXCLUDED.summary;""")
    expected = []
    for stem, slug_name, position in route:
        slug = literal(slug_name)
        sql.append(f"""INSERT INTO articles(author_id,slug,original_lang,category,subcategory,cover_url,status,published_at)
SELECT id,{slug},'ru','society','education',{literal(COVER)},'published',now()
FROM auth_users WHERE email='baimurza.daulet@gmail.com'
ON CONFLICT(slug) DO UPDATE SET category='society',subcategory='education',cover_url=EXCLUDED.cover_url,
status='published',updated_at=now(),published_at=COALESCE(articles.published_at,now());""")
        for lang in ("ru", "kz", "en"):
            suffix = "" if lang == "ru" else f"-{lang}"
            title, summary, body = read_lesson(LESSONS / f"{stem}{suffix}.md")
            sql.append(f"""INSERT INTO article_translations(article_id,lang,title,summary,body_md,source,status)
SELECT id,'{lang}',{literal(title)},{literal(summary)},{literal(body)},'ai','ready'
FROM articles WHERE slug={slug}
ON CONFLICT(article_id,lang) DO UPDATE SET title=EXCLUDED.title,summary=EXCLUDED.summary,
body_md=EXCLUDED.body_md,source='ai',status='ready',updated_at=now();""")
            expected.append({"slug": slug_name, "lang": lang, "title": title,
                             "summary_sha256": hashlib.sha256(summary.encode()).hexdigest(),
                             "body_sha256": hashlib.sha256(body.encode()).hexdigest()})
        sql.append(f"""INSERT INTO article_series_items(series_id,article_id,position)
SELECT s.id,a.id,{position} FROM article_series s,articles a
WHERE s.slug='physics' AND a.slug={slug}
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
    print(f"Prepared {len(expected)} localized pages; no remote changes")


if __name__ == "__main__":
    main()
