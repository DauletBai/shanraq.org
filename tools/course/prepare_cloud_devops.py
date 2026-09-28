#!/usr/bin/env python3
"""Prepare one atomic SQL publication for the Cloud & DevOps course."""
import argparse
import hashlib
import json
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/cloud-devops"
STEMS = ("preface", "01-why-now", "02-service-map", "03-terminal-git",
         "04-linux-access", "05-processes-logs", "06-networks-dns-http",
         "07-containers", "08-images-lifecycle", "09-dockerfile",
         "10-persistent-data", "11-compose", "12-health-resources",
         "13-registry", "14-cicd", "15-github-actions", "16-release-image",
         "17-secrets", "18-deploy-rollback", "19-cloud", "20-server",
         "21-https", "22-iac", "23-kubernetes", "24-incident")
SLUGS = ("cloud-devops-before-start",) + tuple("cloud-devops-" + s for s in STEMS[1:])
LANGS = {"ru": "", "kz": "-kz", "en": "-en"}
COVER = "/static/covers/it/devops/cloud-devops-course.webp"
META = {
    "ru": ("Cloud & DevOps: от локального приложения до надёжного сервиса",
           "Бесплатный практический курс из 24 уроков: Linux, сети, Docker, CI/CD, облачный сервер, HTTPS, OpenTofu, Ansible, Kubernetes, наблюдаемость, резервное копирование и восстановление на одном проекте."),
    "kz": ("Cloud & DevOps: жергілікті қолданбадан сенімді сервиске дейін",
           "24 сабақтан тұратын тегін тәжірибелік курс: бір жоба арқылы Linux, желі, Docker, CI/CD, бұлт сервері, HTTPS, OpenTofu, Ansible, Kubernetes, бақылау, сақтық көшірме және қалпына келтіру."),
    "en": ("Cloud & DevOps: from a local app to a reliable service",
           "A free 24-lesson practical course covering Linux, networking, Docker, CI/CD, a cloud server, HTTPS, OpenTofu, Ansible, Kubernetes, observability, backup, and recovery through one project."),
}
LEAD = re.compile(r"_[^_]+:_\s*\*\*(.+)\*\*\s*$")


def literal(value: str) -> str:
    delimiter = "$cloud_devops_course$"
    if delimiter in value:
        raise ValueError("SQL delimiter occurs in course content")
    return delimiter + value + delimiter


def lesson(path: Path):
    lines = path.read_text(encoding="utf-8").splitlines()
    if len(lines) < 5 or not lines[0].startswith("# "):
        raise ValueError(f"missing title/body: {path}")
    lead = LEAD.fullmatch(lines[2])
    if not lead:
        raise ValueError(f"missing summary lead: {path}")
    body = "\n".join(lines[4:]).strip() + "\n"
    if any(mark in body for mark in ("Редакционный черновик", "Редакциялық жоба", "Editorial draft")):
        raise ValueError(f"editorial note in public body: {path}")
    return lines[0][2:].strip(), lead.group(1), body


def prepare():
    if len(STEMS) != 25 or len(SLUGS) != 25:
        raise ValueError("expected preface and 24 lessons")
    sql = ["BEGIN;", "SELECT pg_advisory_xact_lock(hashtext('shanraq-cloud-devops-course'));" ]
    slugs = ",".join(literal(s) for s in SLUGS)
    sql.append(f"""DO $guard$
BEGIN
  IF (SELECT count(*) FROM auth_users WHERE email='baimurza.daulet@gmail.com') <> 1 THEN
    RAISE EXCEPTION 'Expected course author not found';
  END IF;
  IF EXISTS (SELECT 1 FROM articles a JOIN auth_users u ON u.id=a.author_id
             WHERE a.slug IN ({slugs}) AND u.email <> 'baimurza.daulet@gmail.com') THEN
    RAISE EXCEPTION 'Cloud DevOps slug belongs to another author';
  END IF;
  IF EXISTS (SELECT 1 FROM article_series_items i
             JOIN article_series s ON s.id=i.series_id JOIN articles a ON a.id=i.article_id
             WHERE a.slug IN ({slugs}) AND s.slug <> 'cloud-devops') THEN
    RAISE EXCEPTION 'Cloud DevOps lesson belongs to another course';
  END IF;
END $guard$;""")
    sql.append(f"""INSERT INTO article_series(slug,cover_url,status,code_lang)
VALUES('cloud-devops',{literal(COVER)},'published','shell')
ON CONFLICT(slug) DO UPDATE SET cover_url=EXCLUDED.cover_url,status='published',
code_lang='shell',updated_at=now();""")
    for lang, (title, summary) in META.items():
        sql.append(f"""INSERT INTO article_series_i18n(series_id,lang,title,summary)
SELECT id,'{lang}',{literal(title)},{literal(summary)} FROM article_series
WHERE slug='cloud-devops' ON CONFLICT(series_id,lang) DO UPDATE
SET title=EXCLUDED.title,summary=EXCLUDED.summary;""")
    expected = []
    for index, (stem, slug_name) in enumerate(zip(STEMS, SLUGS)):
        slug = literal(slug_name)
        sql.append(f"""INSERT INTO articles(author_id,slug,original_lang,category,subcategory,cover_url,status,published_at)
SELECT id,{slug},'ru','it','devops',{literal(COVER)},'published',now()
FROM auth_users WHERE email='baimurza.daulet@gmail.com'
ON CONFLICT(slug) DO UPDATE SET cover_url=EXCLUDED.cover_url,status='published',
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
        position = 1 if index == 0 else index * 10
        sql.append(f"""INSERT INTO article_series_items(series_id,article_id,position)
SELECT s.id,a.id,{position} FROM article_series s,articles a
WHERE s.slug='cloud-devops' AND a.slug={slug}
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
