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
    ("17-school-map", "kazakh-language-17-school-map", 170),
    ("18-days-subjects", "kazakh-language-18-days-subjects", 180),
    ("19-clock-time", "kazakh-language-19-clock-time", 190),
    ("20-directions-imperative", "kazakh-language-20-directions-imperative", 200),
    ("21-from-to-route", "kazakh-language-21-from-to-route", 210),
    ("22-city-transport", "kazakh-language-22-city-transport", 220),
    ("23-rules-permission", "kazakh-language-23-rules-permission", 230),
    ("24-school-city-mastery", "kazakh-language-24-school-city-mastery", 240),
    ("25-daily-routine", "kazakh-language-25-daily-routine", 250),
    ("26-sequence-connectors", "kazakh-language-26-sequence-connectors", 260),
    ("27-shop-quantity-price", "kazakh-language-27-shop-quantity-price", 270),
    ("28-cafe-order", "kazakh-language-28-cafe-order", 280),
    ("29-appointment-time", "kazakh-language-29-appointment-time", 290),
    ("30-health-pharmacy", "kazakh-language-30-health-pharmacy", 300),
    ("31-service-error-documents", "kazakh-language-31-service-error-documents", 310),
    ("32-day-services-mastery", "kazakh-language-32-day-services-mastery", 320),
    ("33-case-role-map", "kazakh-language-33-case-role-map", 330),
    ("34-genitive-relationship", "kazakh-language-34-genitive-relationship", 340),
    ("35-accusative-specific-object", "kazakh-language-35-accusative-specific-object", 350),
    ("36-dative-goal-recipient", "kazakh-language-36-dative-goal-recipient", 360),
    ("37-locative-place-time", "kazakh-language-37-locative-place-time", 370),
    ("38-ablative-source", "kazakh-language-38-ablative-source", 380),
    ("39-instrumental-companion-means", "kazakh-language-39-instrumental-companion-means", 390),
    ("40-space-cases-mastery", "kazakh-language-40-space-cases-mastery", 400),
    ("41-time-scene-map", "kazakh-language-41-time-scene-map", 410),
    ("42-habit-present-future", "kazakh-language-42-habit-present-future", 420),
    ("43-action-now-progressive", "kazakh-language-43-action-now-progressive", 430),
    ("44-completed-past", "kazakh-language-44-completed-past", 440),
    ("45-future-plan-intention", "kazakh-language-45-future-plan-intention", 450),
    ("46-time-questions-negation", "kazakh-language-46-time-questions-negation", 460),
    ("47-linked-time-story", "kazakh-language-47-linked-time-story", 470),
    ("48-action-time-mastery", "kazakh-language-48-action-time-mastery", 480),
    ("49-request-help", "kazakh-language-49-request-help", 490),
    ("50-agree-time", "kazakh-language-50-agree-time", 500),
    ("51-possible-necessary", "kazakh-language-51-possible-necessary", 510),
    ("52-cause-consequence", "kazakh-language-52-cause-consequence", 520),
    ("53-if-plan-changes", "kazakh-language-53-if-plan-changes", 530),
    ("54-explain-problem", "kazakh-language-54-explain-problem", 540),
    ("55-choose-together", "kazakh-language-55-choose-together", 550),
    ("56-plans-problems-mastery", "kazakh-language-56-plans-problems-mastery", 560),
    ("57-speech-chunks", "kazakh-language-57-speech-chunks", 570),
    ("58-question-intonation", "kazakh-language-58-question-intonation", 580),
    ("59-listen-for-main-point", "kazakh-language-59-listen-for-main-point", 590),
    ("60-repair-a-missed-word", "kazakh-language-60-repair-a-missed-word", 600),
    ("61-short-story-chain", "kazakh-language-61-short-story-chain", 610),
    ("62-unknown-word-context", "kazakh-language-62-unknown-word-context", 620),
    ("63-retell-and-follow-up", "kazakh-language-63-retell-and-follow-up", 630),
    ("64-connected-speech-mastery", "kazakh-language-64-connected-speech-mastery", 640),
)
SERIES_COVER = "/static/covers/school/kazakh-language/foundations/01-first-conversation.webp"
BLOCK_COVERS = {
    1: SERIES_COVER,
    2: "/static/covers/school/kazakh-language/people-home/02-people-home.webp",
    3: "/static/covers/school/kazakh-language/school-city/03-school-city.webp",
    4: "/static/covers/school/kazakh-language/day-services/04-day-services.webp",
    5: "/static/covers/school/kazakh-language/space-cases/05-space-cases.webp",
    6: "/static/covers/school/kazakh-language/action-time/06-action-time.webp",
    7: "/static/covers/school/kazakh-language/plans-problems/07-plans-problems.webp",
    8: "/static/covers/school/kazakh-language/connected-speech/08-connected-speech.webp",
}
META = {
    "ru": (
        "Казахский язык: начинаем говорить с первой встречи",
        "Бесплатный практический курс казахского языка от первых звуков до самостоятельного общения. 64 урока учат знакомиться, ориентироваться в городе, решать бытовые задачи, рассказывать о времени, договариваться, понимать связную речь и точно пересказывать услышанное в проекте «Моя среда».",
    ),
    "kz": (
        "Қазақ тілі: алғашқы кездесуден бастап сөйлейміз",
        "Алғашқы дыбыстардан дербес қарым-қатынасқа дейінгі тегін тәжірибелік қазақ тілі курсы. 64 сабақ танысуды, қалада бағдарлауды, тұрмыстық міндеттерді шешуді, уақытты баяндауды, келісуді, байланысты сөзді түсінуді және естігенді дәл жеткізуді «Менің ортам» жобасында дамытады.",
    ),
    "en": (
        "Kazakh: start speaking from the first meeting",
        "A free practical Kazakh course from first sounds to independent communication. Sixty-four lessons cover introductions, city routes, everyday services, past and planned actions, negotiation, connected speech, and accurate retelling in the continuing My World project.",
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
    if len(ROUTE) != 64:
        raise ValueError("the current release must contain exactly sixty-four lessons")
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
        cover = BLOCK_COVERS[(position - 1) // 80 + 1]
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
