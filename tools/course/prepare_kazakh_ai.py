#!/usr/bin/env python3
"""Prepare one atomic SQL publication for the Kazakh AI preface and lessons 1–20."""

import argparse
import hashlib
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/kazakh-ai"
STEMS = ("preface", "01-goal", "02-first-program", "03-language-types",
         "04-word-parts", "05-manual-model", "06-cyrillic",
         "07-words", "08-roots", "09-plurals", "10-order",
         "11-cases", "12-ambiguity", "13-entities", "14-intent", "15-abstain", "16-catalog", "17-lookup", "18-provenance",
         "19-template", "20-expiry")
SLUGS = ("kazakh-ai-before-start",) + tuple("kazakh-ai-" + s for s in STEMS[1:])
LANGS = {"ru": "", "kz": "-kz", "en": "-en"}
META = {
    "ru": ("ИИ без LLM: казахская модель с нуля",
           "Создаём локальную модель для казахского текста: разбираем слова, проверяем источники и учимся отвечать «не знаю». Открыты первые двадцать уроков из запланированных 30."),
    "kz": ("Үлкен тілдік үлгісіз ЖИ: қазақша мәтін үлгісі",
           "Қазақша сөздерді талдап, дереккөзді тексеретін жергілікті үлгі құрастырамыз. Жоспарланған 30 сабақтың алғашқы жиырмасы ашық."),
    "en": ("AI without an LLM: a Kazakh text model",
           "Build a local Kazakh text model that analyzes words, checks sources, and can say “I don't know”. The first twenty of 30 planned lessons are open."),
}


def literal(value: str) -> str:
    delimiter = "$kazakh_ai_course$"
    if delimiter in value:
        raise ValueError("SQL delimiter occurs in course content")
    return delimiter + value + delimiter


def lesson(path: Path) -> tuple[str, str, str]:
    text = path.read_text(encoding="utf-8").strip()
    first, _, body = text.partition("\n")
    if not first.startswith("# ") or not body.strip():
        raise ValueError(f"missing title/body: {path}")
    if "Статус:" in body or "Status:" in body:
        raise ValueError(f"editorial note in public body: {path}")
    title = first[2:].strip()
    body = body.strip() + "\n"
    # Summaries are short and independent of any editorial heading.
    paragraphs = [p.strip() for p in body.split("\n\n") if p.strip()]
    summary = next(p for p in paragraphs if not p.startswith("#"))
    summary = summary.replace("**", "").replace("`", "")
    if len(summary) > 450:
        summary = summary[:447].rsplit(" ", 1)[0] + "…"
    return title, summary, body


def prepare():
    sql = ["BEGIN;", "SELECT pg_advisory_xact_lock(hashtext('shanraq-kazakh-ai-course'));"]
    slugs = ",".join(literal(s) for s in SLUGS)
    sql.append(f"""DO $guard$
BEGIN
  IF (SELECT count(*) FROM auth_users WHERE email='baimurza.daulet@gmail.com') <> 1 THEN
    RAISE EXCEPTION 'Expected course author not found';
  END IF;
  IF EXISTS (SELECT 1 FROM articles a JOIN auth_users u ON u.id=a.author_id
             WHERE a.slug IN ({slugs}) AND u.email <> 'baimurza.daulet@gmail.com') THEN
    RAISE EXCEPTION 'Course slug belongs to another author';
  END IF;
  IF EXISTS (SELECT 1 FROM article_series_items i
             JOIN article_series s ON s.id=i.series_id
             JOIN articles a ON a.id=i.article_id
             WHERE a.slug IN ({slugs}) AND s.slug <> 'kazakh-ai') THEN
    RAISE EXCEPTION 'Lesson belongs to another course';
  END IF;
END $guard$;""")
    cover = literal("/static/brand/kazakh-ai-mark.svg")
    sql.append(f"""INSERT INTO article_series(slug,cover_url,status,code_lang)
VALUES('kazakh-ai',{cover},'published','python')
ON CONFLICT(slug) DO UPDATE SET cover_url=EXCLUDED.cover_url,status='published',
code_lang='python',updated_at=now();""")
    for lang, (title, summary) in META.items():
        sql.append(f"""INSERT INTO article_series_i18n(series_id,lang,title,summary)
SELECT id,'{lang}',{literal(title)},{literal(summary)} FROM article_series
WHERE slug='kazakh-ai' ON CONFLICT(series_id,lang) DO UPDATE
SET title=EXCLUDED.title,summary=EXCLUDED.summary;""")
    expected = []
    for position, (stem, slug_name) in enumerate(zip(STEMS, SLUGS)):
        slug = literal(slug_name)
        sql.append(f"""INSERT INTO articles(author_id,slug,original_lang,category,subcategory,cover_url,status,published_at)
SELECT id,{slug},'ru','it','programming','', 'published',now()
FROM auth_users WHERE email='baimurza.daulet@gmail.com'
ON CONFLICT(slug) DO UPDATE SET cover_url='',status='published',
updated_at=now(),published_at=COALESCE(articles.published_at,now());""")
        for lang, suffix in LANGS.items():
            title, summary, body = lesson(LESSONS / f"{stem}{suffix}.md")
            sql.append(f"""INSERT INTO article_translations(article_id,lang,title,summary,body_md,source,status)
SELECT id,'{lang}',{literal(title)},{literal(summary)},{literal(body)},'ai','ready'
FROM articles WHERE slug={slug}
ON CONFLICT(article_id,lang) DO UPDATE SET title=EXCLUDED.title,summary=EXCLUDED.summary,
body_md=EXCLUDED.body_md,source='ai',status='ready',updated_at=now();""")
            expected.append({"slug": slug_name, "lang": lang, "title": title,
                             "body_sha256": hashlib.sha256(body.encode()).hexdigest(),
                             "summary_sha256": hashlib.sha256(summary.encode()).hexdigest()})
        pos = 1 if position == 0 else position * 10
        sql.append(f"""INSERT INTO article_series_items(series_id,article_id,position)
SELECT s.id,a.id,{pos} FROM article_series s,articles a
WHERE s.slug='kazakh-ai' AND a.slug={slug}
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
    print(f"Prepared {len(expected)} localized pages in one transaction; no remote changes")


if __name__ == "__main__":
    main()
