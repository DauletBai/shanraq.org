#!/usr/bin/env python3
"""Prepare the completed Kazakh-language blocks as one atomic publication."""
import argparse
import hashlib
import json
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/kazakh-language"
ROUTE = (
    ("01-map-first-conversation", "kazakh-language-01-map-first-conversation", 10),
    ("02-nine-letters-sounds", "kazakh-language-02-nine-letters-sounds", 20),
    ("03-vowel-harmony", "kazakh-language-03-vowel-harmony", 30),
    ("04-sentence-skeleton", "kazakh-language-04-sentence-skeleton", 40),
    ("05-men-introduction", "kazakh-language-05-men-introduction", 50),
    ("06-questions-answers", "kazakh-language-06-questions-answers", 60),
    ("07-repair-politeness", "kazakh-language-07-repair-politeness", 70),
    ("08-first-contact-mastery", "kazakh-language-08-first-contact-mastery", 80),
    ("09-short-vowels-y-i", "kazakh-language-09-short-vowels-y-i", 90),
    ("10-family-who", "kazakh-language-10-family-who", 100),
    ("11-belonging-possessives", "kazakh-language-11-belonging-possessives", 110),
    ("12-plural-families", "kazakh-language-12-plural-families", 120),
    ("13-numbers-age", "kazakh-language-13-numbers-age", 130),
    ("14-home-locative", "kazakh-language-14-home-locative", 140),
    ("15-there-is-have", "kazakh-language-15-there-is-have", 150),
    ("16-people-home-mastery", "kazakh-language-16-people-home-mastery", 160),
)
SERIES_COVER = "/static/covers/school/kazakh-language/foundations/01-first-conversation.webp"
BLOCK_COVERS = {
    1: SERIES_COVER,
    2: "/static/covers/school/kazakh-language/people-home/02-people-home.webp",
}
META = {
    "ru": (
        "Казахский язык: начинаем говорить с первой встречи",
        "Бесплатный практический курс казахского языка от первых звуков до самостоятельного общения. Первые 16 уроков учат знакомиться, слышать естественные ы/і, рассказывать о людях и доме, а также развивают проект «Моя среда».",
    ),
    "kz": (
        "Қазақ тілі: алғашқы кездесуден бастап сөйлейміз",
        "Алғашқы дыбыстардан дербес қарым-қатынасқа дейінгі тегін тәжірибелік қазақ тілі курсы. Алғашқы 16 сабақ танысуды, табиғи ы/і айтылымын, адамдар мен үй туралы сөйлеуді және «Менің ортам» жобасын дамытады.",
    ),
    "en": (
        "Kazakh: start speaking from the first meeting",
        "A free practical Kazakh course from first sounds to independent communication. The first 16 lessons cover introductions, natural ы/і, people and home, and the continuing My World project.",
    ),
}
LEAD = re.compile(r"_[^_]+:_\s*\*\*(.+)\*\*\s*$")


def literal(value: str) -> str:
    delimiter = "$kazakh_language_course$"
    if delimiter in value:
        raise ValueError("SQL delimiter occurs in Kazakh-language content")
    return delimiter + value + delimiter


def lesson(path: Path):
    lines = path.read_text(encoding="utf-8").splitlines()
    if len(lines) < 5 or not lines[0].startswith("# "):
        raise ValueError(f"missing title/body: {path}")
    lead = LEAD.fullmatch(lines[2])
    if not lead:
        raise ValueError(f"missing summary lead: {path}")
    return lines[0][2:].strip(), lead.group(1), "\n".join(lines[4:]).strip() + "\n"


def prepare():
    if len(ROUTE) != 16:
        raise ValueError("the current release must contain exactly sixteen lessons")
    sql = ["BEGIN;", "SELECT pg_advisory_xact_lock(hashtext('shanraq-kazakh-language-course'));"]
    slugs = ",".join(literal(slug) for _, slug, _ in ROUTE)
    sql.append(f"""DO $guard$
BEGIN
  IF (SELECT count(*) FROM auth_users WHERE email='baimurza.daulet@gmail.com') <> 1 THEN
    RAISE EXCEPTION 'Expected course author not found';
  END IF;
  IF EXISTS (SELECT 1 FROM articles a JOIN auth_users u ON u.id=a.author_id
             WHERE a.slug IN ({slugs}) AND u.email <> 'baimurza.daulet@gmail.com') THEN
    RAISE EXCEPTION 'Kazakh-language slug belongs to another author';
  END IF;
  IF EXISTS (SELECT 1 FROM article_series_items i
             JOIN article_series s ON s.id=i.series_id JOIN articles a ON a.id=i.article_id
             WHERE a.slug IN ({slugs}) AND s.slug <> 'kazakh-language') THEN
    RAISE EXCEPTION 'Kazakh-language lesson belongs to another course';
  END IF;
END $guard$;""")
    sql.append(f"""INSERT INTO article_series(slug,cover_url,status,code_lang)
VALUES('kazakh-language',{literal(SERIES_COVER)},'published','kazakh')
ON CONFLICT(slug) DO UPDATE SET cover_url=EXCLUDED.cover_url,status='published',
code_lang='kazakh',updated_at=now();""")
    for lang, (title, summary) in META.items():
        sql.append(f"""INSERT INTO article_series_i18n(series_id,lang,title,summary)
SELECT id,'{lang}',{literal(title)},{literal(summary)} FROM article_series
WHERE slug='kazakh-language' ON CONFLICT(series_id,lang) DO UPDATE
SET title=EXCLUDED.title,summary=EXCLUDED.summary;""")

    expected = []
    for stem, slug_name, position in ROUTE:
        slug = literal(slug_name)
        cover = BLOCK_COVERS[1 if position <= 80 else 2]
        sql.append(f"""INSERT INTO articles(author_id,slug,original_lang,category,subcategory,cover_url,status,published_at)
SELECT id,{slug},'ru','society','education',{literal(cover)},'published',now()
FROM auth_users WHERE email='baimurza.daulet@gmail.com'
ON CONFLICT(slug) DO UPDATE SET category='society',subcategory='education',
cover_url=EXCLUDED.cover_url,status='published',updated_at=now(),
published_at=COALESCE(articles.published_at,now());""")
        for lang in ("ru", "kz", "en"):
            suffix = "" if lang == "ru" else f"-{lang}"
            title, summary, body = lesson(LESSONS / f"{stem}{suffix}.md")
            sql.append(f"""INSERT INTO article_translations(article_id,lang,title,summary,body_md,source,status)
SELECT id,'{lang}',{literal(title)},{literal(summary)},{literal(body)},'ai','ready'
FROM articles WHERE slug={slug}
ON CONFLICT(article_id,lang) DO UPDATE SET title=EXCLUDED.title,summary=EXCLUDED.summary,
body_md=EXCLUDED.body_md,source='ai',status='ready',updated_at=now();""")
            expected.append({
                "slug": slug_name,
                "lang": lang,
                "title": title,
                "body_sha256": hashlib.sha256(body.encode()).hexdigest(),
                "summary_sha256": hashlib.sha256(summary.encode()).hexdigest(),
            })
        sql.append(f"""INSERT INTO article_series_items(series_id,article_id,position)
SELECT s.id,a.id,{position} FROM article_series s,articles a
WHERE s.slug='kazakh-language' AND a.slug={slug}
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
    print(f"Prepared {len(ROUTE)} lessons and {len(expected)} localized translations")


if __name__ == "__main__":
    main()
