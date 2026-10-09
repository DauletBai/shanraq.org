#!/usr/bin/env python3
"""Prepare the first three Informatics blocks as one atomic SQL publication."""
import argparse
import hashlib
import json
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/informatics"
ROUTE = (
    ("01-whole-map", "informatics-01-whole-map", 10),
    ("02-diagnostic-product", "informatics-02-diagnostic-product", 20),
    ("03-device-system", "informatics-03-device-system", 30),
    ("04-input-output-sensors", "informatics-04-input-output-sensors", 40),
    ("05-cpu-memory-storage", "informatics-05-cpu-memory-storage", 50),
    ("06-os-process-app", "informatics-06-os-process-app", 60),
    ("07-files-folders-paths", "informatics-07-files-folders-paths", 70),
    ("08-formats-software-licenses", "informatics-08-formats-software-licenses", 80),
    ("09-versions-collaboration-accessibility", "informatics-09-versions-collaboration-accessibility", 90),
    ("10-systems-mastery", "informatics-10-systems-mastery", 100),
    ("11-bits-states", "informatics-11-bits-states", 110),
    ("12-binary-numbers", "informatics-12-binary-numbers", 120),
    ("13-text-unicode", "informatics-13-text-unicode", 130),
    ("14-pixels-color", "informatics-14-pixels-color", 140),
    ("15-sound-video-sampling", "informatics-15-sound-video-sampling", 150),
    ("16-compression", "informatics-16-compression", 160),
    ("17-integrity-errors", "informatics-17-integrity-errors", 170),
    ("18-representation-mastery", "informatics-18-representation-mastery", 180),
    ("19-problem-decomposition", "informatics-19-problem-decomposition", 190),
    ("20-state-variables", "informatics-20-state-variables", 200),
    ("21-sequence-tracing", "informatics-21-sequence-tracing", 210),
    ("22-conditions-boundaries", "informatics-22-conditions-boundaries", 220),
    ("23-loops-invariants", "informatics-23-loops-invariants", 230),
    ("24-functions-contracts", "informatics-24-functions-contracts", 240),
    ("25-correctness-efficiency", "informatics-25-correctness-efficiency", 250),
    ("26-algorithms-mastery", "informatics-26-algorithms-mastery", 260),
)
COVER = "/static/covers/school/informatics/foundations/01-digital-world-computer-project.webp"
REPRESENTATION_COVER = "/static/covers/school/informatics/representation/02-information-data.webp"
ALGORITHMS_COVER = "/static/covers/school/informatics/algorithms/03-algorithmic-thinking.webp"
META = {
    "ru": (
        "Информатика: создаём своего цифрового помощника",
        "26 бесплатных занятий: устройство компьютера, представление информации и алгоритмы. Шаг за шагом создаём и проверяем цифрового помощника.",
    ),
    "kz": (
        "Информатика: өз цифрлық көмекшімізді жасаймыз",
        "26 тегін сабақ: компьютер құрылысы, ақпаратты көрсету және алгоритмдер. Цифрлық көмекшіні қадамдап құрып, тексереміз.",
    ),
    "en": (
        "Informatics: build your own digital assistant",
        "26 free lessons cover computer systems, data representation, and algorithms through one continuing, tested digital-assistant project.",
    ),
}
LEAD = re.compile(r"_[^_]+:_\s*\*\*(.+)\*\*\s*$")


def literal(value: str) -> str:
    delimiter = "$informatics_course$"
    if delimiter in value:
        raise ValueError("SQL delimiter occurs in Informatics content")
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
    if len(ROUTE) != 26:
        raise ValueError("the first three releases must contain exactly twenty-six lessons")
    sql = [
        "BEGIN;",
        "SELECT pg_advisory_xact_lock(hashtext('shanraq-informatics-course'));",
    ]
    slugs = ",".join(literal(slug) for _, slug, _ in ROUTE)
    sql.append(f"""DO $guard$
BEGIN
  IF (SELECT count(*) FROM auth_users WHERE email='baimurza.daulet@gmail.com') <> 1 THEN
    RAISE EXCEPTION 'Expected course author not found';
  END IF;
  IF EXISTS (SELECT 1 FROM articles a JOIN auth_users u ON u.id=a.author_id
             WHERE a.slug IN ({slugs}) AND u.email <> 'baimurza.daulet@gmail.com') THEN
    RAISE EXCEPTION 'Informatics slug belongs to another author';
  END IF;
  IF EXISTS (SELECT 1 FROM article_series_items i
             JOIN article_series s ON s.id=i.series_id JOIN articles a ON a.id=i.article_id
             WHERE a.slug IN ({slugs}) AND s.slug <> 'informatics') THEN
    RAISE EXCEPTION 'Informatics lesson belongs to another course';
  END IF;
END $guard$;""")
    sql.append(f"""INSERT INTO article_series(slug,cover_url,status,code_lang)
VALUES('informatics',{literal(COVER)},'published','informatics')
ON CONFLICT(slug) DO UPDATE SET cover_url=EXCLUDED.cover_url,status='published',
code_lang='informatics',updated_at=now();""")
    for lang, (title, summary) in META.items():
        sql.append(f"""INSERT INTO article_series_i18n(series_id,lang,title,summary)
SELECT id,'{lang}',{literal(title)},{literal(summary)} FROM article_series
WHERE slug='informatics' ON CONFLICT(series_id,lang) DO UPDATE
SET title=EXCLUDED.title,summary=EXCLUDED.summary;""")

    expected = []
    for stem, slug_name, position in ROUTE:
        slug = literal(slug_name)
        lesson_cover = COVER if position <= 100 else REPRESENTATION_COVER if position <= 180 else ALGORITHMS_COVER
        sql.append(f"""INSERT INTO articles(author_id,slug,original_lang,category,subcategory,cover_url,status,published_at)
SELECT id,{slug},'ru','society','education',{literal(lesson_cover)},'published',now()
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
WHERE s.slug='informatics' AND a.slug={slug}
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
