#!/usr/bin/env python3
"""Prepare one atomic SQL publication for the released mathematics routes."""
import argparse
import hashlib
import json
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/mathematics"
# The source filenames retain their historical lesson numbers, while position
# follows the dependency graph.  The foundations block therefore sits between
# diagnosis and fractions without changing any published URL.
ROUTE = (
    ("preface", "mathematics-before-start", 1),
    ("01-diagnostic", "math-01-diagnostic", 10),
    ("10-quantity-counting", "math-10-quantity-counting", 11),
    ("11-place-value", "math-11-place-value", 12),
    ("12-addition-subtraction", "math-12-addition-subtraction", 13),
    ("13-multiplication-division", "math-13-multiplication-division", 14),
    ("14-order-estimation", "math-14-order-estimation", 15),
    ("15-negative-numbers", "math-15-negative-numbers", 16),
    ("16-divisibility-primes", "math-16-divisibility-primes", 17),
    ("17-decimal-fractions", "math-17-decimal-fractions", 18),
    ("18-foundations-mastery", "math-18-foundations-mastery", 19),
    ("02-fraction-meaning", "math-02-fraction-meaning", 20),
    ("03-equivalent-fractions", "math-03-equivalent-fractions", 30),
    ("04-compare-fractions", "math-04-compare-fractions", 40),
    ("05-fraction-operations", "math-05-fraction-operations", 50),
    ("06-ratio", "math-06-ratio", 60),
    ("07-percent", "math-07-percent", 70),
    ("08-proportion", "math-08-proportion", 80),
    ("09-mastery", "math-09-mastery", 90),
    ("19-variables", "math-19-variables", 100),
    ("20-expressions", "math-20-expressions", 110),
    ("21-equations-basic", "math-21-equations-basic", 120),
    ("22-coordinates", "math-22-coordinates", 130),
    ("23-inequalities", "math-23-inequalities", 140),
    ("24-prealgebra-mastery", "math-24-prealgebra-mastery", 150),
    ("25-linear-functions", "math-25-linear-functions", 160),
    ("26-systems", "math-26-systems", 170),
    ("27-powers-roots", "math-27-powers-roots", 180),
    ("28-polynomials", "math-28-polynomials", 190),
    ("29-quadratics", "math-29-quadratics", 200),
    ("30-exponential-log", "math-30-exponential-log", 210),
    ("31-sequences", "math-31-sequences", 220),
    ("32-algebra-mastery", "math-32-algebra-mastery", 230),
)
STEMS = tuple(item[0] for item in ROUTE)
SLUGS = tuple(item[1] for item in ROUTE)
MASTERY_STEMS = ("18-foundations-mastery", "09-mastery", "24-prealgebra-mastery", "32-algebra-mastery")
TEACHING_STEMS = tuple(
    stem for stem in STEMS if stem != "preface" and stem not in MASTERY_STEMS
)
COVER = "/static/covers/school/mathematics/mathematics-foundations.webp"
META = {
    "ru": (
        "Математика: от фундамента к высшей математике",
        "Бесплатный курс по карте зависимостей, а не по классам. 32 занятия ведут от чисел и дробей через предалгебру к функциям, системам, многочленам, параболам и показательному росту; каждый блок завершается проверкой переноса.",
    ),
    "kz": (
        "Математика: іргетастан жоғары математикаға дейін",
        "Сыныптарға емес, ұғымдардың тәуелділік картасына құрылған тегін курс. 32 сабақ сандар мен бөлшектерден бастап, алгебра алдындағы ұғымдар, функциялар, жүйелер, көпмүшелер, параболалар және көрсеткіштік өсуге дейін жетелейді.",
    ),
    "en": (
        "Mathematics: from foundations to higher mathematics",
        "A free course organized by idea dependencies rather than grade levels. Its 32 lessons lead from numbers and fractions through prealgebra to functions, systems, polynomials, parabolas, and exponential growth.",
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
    if len(ROUTE) != 33 or len(STEMS) != len(SLUGS):
        raise ValueError("expected one preface and thirty-two lessons")
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
    for stem, slug_name, position in ROUTE:
        slug = literal(slug_name)
        sql.append(f"""INSERT INTO articles(author_id,slug,original_lang,category,subcategory,cover_url,status,published_at)
SELECT id,{slug},'ru','society','education',{literal(COVER)},'published',now()
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
    print(f"Prepared {len(expected)} localized pages in one transaction; no remote changes")


if __name__ == "__main__":
    main()
