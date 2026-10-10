#!/usr/bin/env python3
"""Prepare Informatics lessons 1–72 as an atomic SQL publication."""
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
    ("27-python-first-state", "informatics-27-python-first-state", 270),
    ("28-python-input-output-types", "informatics-28-python-input-output-types", 280),
    ("29-python-expressions", "informatics-29-python-expressions", 290),
    ("30-python-conditions", "informatics-30-python-conditions", 300),
    ("31-python-loops", "informatics-31-python-loops", 310),
    ("32-python-collections", "informatics-32-python-collections", 320),
    ("33-python-text-dates", "informatics-33-python-text-dates", 330),
    ("34-python-functions", "informatics-34-python-functions", 340),
    ("35-python-files-json", "informatics-35-python-files-json", 350),
    ("36-python-errors-debugging", "informatics-36-python-errors-debugging", 360),
    ("37-python-tests-modules", "informatics-37-python-tests-modules", 370),
    ("38-python-cli-release", "informatics-38-python-cli-release", 380),
    ("39-network-message-journey", "informatics-39-network-message-journey", 390),
    ("40-ip-router-packets", "informatics-40-ip-router-packets", 400),
    ("41-dns-names", "informatics-41-dns-names", 410),
    ("42-tcp-udp-delivery", "informatics-42-tcp-udp-delivery", 420),
    ("43-tls-trust", "informatics-43-tls-trust", 430),
    ("44-http-browser-server", "informatics-44-http-browser-server", 440),
    ("45-html-css-accessibility", "informatics-45-html-css-accessibility", 450),
    ("46-web-cloud-release", "informatics-46-web-cloud-release", 460),
    ("47-observations-data-schema", "informatics-47-observations-data-schema", 470),
    ("48-spreadsheets-formulas", "informatics-48-spreadsheets-formulas", 480),
    ("49-cleaning-provenance", "informatics-49-cleaning-provenance", 490),
    ("50-charts-honesty", "informatics-50-charts-honesty", 500),
    ("51-relational-keys", "informatics-51-relational-keys", 510),
    ("52-sql-queries", "informatics-52-sql-queries", 520),
    ("53-joins-reports", "informatics-53-joins-reports", 530),
    ("54-data-release", "informatics-54-data-release", 540),
    ("55-threat-model-cia", "informatics-55-threat-model-cia", 550),
    ("56-passwords-hashing-2fa", "informatics-56-passwords-hashing-2fa", 560),
    ("57-authorization-least-privilege", "informatics-57-authorization-least-privilege", 570),
    ("58-phishing-social-deepfakes", "informatics-58-phishing-social-deepfakes", 580),
    ("59-encryption-keys", "informatics-59-encryption-keys", 590),
    ("60-backup-updates-logs", "informatics-60-backup-updates-logs", 600),
    ("61-privacy-rights-wellbeing", "informatics-61-privacy-rights-wellbeing", 610),
    ("62-security-release", "informatics-62-security-release", 620),
    ("63-rules-algorithms-models", "informatics-63-rules-algorithms-models", 630),
    ("64-features-labels-datasets", "informatics-64-features-labels-datasets", 640),
    ("65-train-classifier", "informatics-65-train-classifier", 650),
    ("66-validation-metrics", "informatics-66-validation-metrics", 660),
    ("67-bias-fairness-privacy", "informatics-67-bias-fairness-privacy", 670),
    ("68-generative-ai-llm", "informatics-68-generative-ai-llm", 680),
    ("69-ai-verification-sources", "informatics-69-ai-verification-sources", 690),
    ("70-human-controlled-ai", "informatics-70-human-controlled-ai", 700),
    ("71-release-testing-docs", "informatics-71-release-testing-docs", 710),
    ("72-capstone-defense", "informatics-72-capstone-defense", 720),
)
COVER = "/static/covers/school/informatics/foundations/01-digital-world-computer-project.webp"
REPRESENTATION_COVER = "/static/covers/school/informatics/representation/02-information-data.webp"
ALGORITHMS_COVER = "/static/covers/school/informatics/algorithms/03-algorithmic-thinking.webp"
PYTHON_COVER = "/static/covers/school/informatics/python/04-python-assistant.webp"
WEB_COVER = "/static/covers/school/informatics/web/05-internet-web-cloud.webp"
DATA_COVER = "/static/covers/school/informatics/data/06-data-tables-databases.webp"
SECURITY_COVER = "/static/covers/school/informatics/security/07-digital-security.webp"
AI_COVER = "/static/covers/school/informatics/ai-release/08-ai-final-project.webp"
META = {
    "ru": (
        "Информатика: создаём своего цифрового помощника",
        "72 бесплатных урока: компьютер, Python, веб, данные, безопасность и ИИ. Шаг за шагом создаём и защищаем цифрового помощника.",
    ),
    "kz": (
        "Информатика: өз цифрлық көмекшімізді жасаймыз",
        "72 тегін сабақ: компьютер, Python, веб, деректер, қауіпсіздік және ЖИ. Цифрлық көмекшіні қадамдап құрып, қорғаймыз.",
    ),
    "en": (
        "Informatics: build your own digital assistant",
        "72 free lessons cover computing, Python, the web, data, security and AI through one continuing, tested digital-assistant project.",
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


def prepare(start=1):
    if len(ROUTE) != 72 or start not in (1, 39, 47, 55, 63):
        raise ValueError("expected seventy-two lessons; start must be 1, 39, 47, 55 or 63")
    active_route = ROUTE[38:46] if start == 39 else ROUTE[46:54] if start == 47 else ROUTE[start-1:]
    sql = [
        "BEGIN;",
        "SELECT pg_advisory_xact_lock(hashtext('shanraq-informatics-course'));",
    ]
    slugs = ",".join(literal(slug) for _, slug, _ in active_route)
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
    for stem, slug_name, position in active_route:
        slug = literal(slug_name)
        lesson_cover = COVER if position <= 100 else REPRESENTATION_COVER if position <= 180 else ALGORITHMS_COVER if position <= 260 else PYTHON_COVER if position <= 380 else WEB_COVER if position <= 460 else DATA_COVER if position <= 540 else SECURITY_COVER if position <= 620 else AI_COVER
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
    parser.add_argument("--start", type=int, choices=(1, 39, 47, 55, 63), default=1,
                        help="1 rebuilds all; 55 publishes both final blocks")
    args = parser.parse_args()
    sql, expected = prepare(args.start)
    args.sql.write_text(sql, encoding="utf-8")
    args.expected.write_text(json.dumps(expected, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(f"Prepared {len(expected)} localized pages in one transaction; no remote changes")


if __name__ == "__main__":
    main()
