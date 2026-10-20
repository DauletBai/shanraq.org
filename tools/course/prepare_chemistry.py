#!/usr/bin/env python3
"""Prepare an all-or-nothing chemistry block SQL release without remote writes."""
import argparse
import hashlib
import json
from pathlib import Path
import re

try:
    from .build_chemistry_block import LESSONS, LESSON_DIR, ROOT, VISUAL_FILES
except ImportError:
    from build_chemistry_block import LESSONS, LESSON_DIR, ROOT, VISUAL_FILES

COVER = VISUAL_FILES[0]
META = {
    'ru': ('Химия: расследуем мир веществ', 'Понятный курс химии через безопасные опыты, опорные схемы, 3D-визуализации, разбор ошибок и проверку знаний.'),
    'kz': ('Химия: заттар әлемін зерттейміз', 'Қауіпсіз тәжірибе, тірек сызба, 3D-көрініс, қатені талдау және білімді тексеру арқылы химияны үйренеміз.'),
    'en': ('Chemistry: investigate the world of substances', 'Learn chemistry through safe investigations, support maps, 3D visuals, error analysis and retrieval practice.'),
}
LEAD = re.compile(r'_Лид \(summary\):_ \*\*(.+)\*\*$')


def literal(value):
    delimiter = '$chemistry_course$'
    if delimiter in value:
        raise ValueError('SQL delimiter in text')
    return delimiter + value + delimiter


def read_lesson(path):
    lines = path.read_text(encoding='utf-8').splitlines()
    if len(lines) < 6 or not lines[0].startswith('# '):
        raise ValueError(f'missing lesson title: {path}')
    match = LEAD.fullmatch(lines[2])
    if match is None:
        raise ValueError(f'missing summary: {path}')
    return lines[0][2:].strip(), match.group(1), '\n'.join(lines[4:]).strip() + '\n'


def prepare():
    if len(LESSONS) != 5:
        raise ValueError('first release requires five lessons')
    routes = [(i, row['stem'], f'chemistry-{i:02d}-{row["stem"].split("-",1)[1]}')
              for i, row in enumerate(LESSONS, 1)]
    slug_list = ','.join(literal(slug) for _, _, slug in routes)
    sql = ['BEGIN;', "SELECT pg_advisory_xact_lock(hashtext('shanraq-chemistry-course'));",
           f"""DO $guard$
BEGIN
  IF (SELECT count(*) FROM auth_users WHERE email='baimurza.daulet@gmail.com') <> 1 THEN
    RAISE EXCEPTION 'Course author missing';
  END IF;
  IF EXISTS (SELECT 1 FROM articles a JOIN auth_users u ON u.id=a.author_id
             WHERE a.slug IN ({slug_list}) AND u.email <> 'baimurza.daulet@gmail.com') THEN
    RAISE EXCEPTION 'Chemistry slug belongs to another author';
  END IF;
  IF EXISTS (SELECT 1 FROM article_series_items i JOIN article_series s ON s.id=i.series_id
             JOIN articles a ON a.id=i.article_id WHERE a.slug IN ({slug_list}) AND s.slug <> 'chemistry') THEN
    RAISE EXCEPTION 'Chemistry lesson belongs to another course';
  END IF;
END $guard$;""",
           f"""INSERT INTO article_series(slug,cover_url,status,code_lang)
VALUES('chemistry',{literal(COVER)},'published','math')
ON CONFLICT(slug) DO UPDATE SET cover_url=EXCLUDED.cover_url,status='published',code_lang='math',updated_at=now();"""]
    for lang, (title, summary) in META.items():
        sql.append(f"""INSERT INTO article_series_i18n(series_id,lang,title,summary)
SELECT id,'{lang}',{literal(title)},{literal(summary)} FROM article_series WHERE slug='chemistry'
ON CONFLICT(series_id,lang) DO UPDATE SET title=EXCLUDED.title,summary=EXCLUDED.summary;""")
    expected = []
    for number, stem, slug in routes:
        qslug = literal(slug)
        sql.append(f"""INSERT INTO articles(author_id,slug,original_lang,category,subcategory,cover_url,status,published_at)
SELECT id,{qslug},'ru','society','education',{literal(VISUAL_FILES[number-1])},'published',now()
FROM auth_users WHERE email='baimurza.daulet@gmail.com'
ON CONFLICT(slug) DO UPDATE SET category='society',subcategory='education',cover_url=EXCLUDED.cover_url,
status='published',updated_at=now(),published_at=COALESCE(articles.published_at,now());""")
        for lang in ('ru', 'kz', 'en'):
            suffix = '' if lang == 'ru' else '-' + lang
            title, summary, body = read_lesson(LESSON_DIR / f'{stem}{suffix}.md')
            sql.append(f"""INSERT INTO article_translations(article_id,lang,title,summary,body_md,source,status)
SELECT id,'{lang}',{literal(title)},{literal(summary)},{literal(body)},'ai','ready'
FROM articles WHERE slug={qslug}
ON CONFLICT(article_id,lang) DO UPDATE SET title=EXCLUDED.title,summary=EXCLUDED.summary,
body_md=EXCLUDED.body_md,source='ai',status='ready',updated_at=now();""")
            expected.append({'slug': slug, 'lang': lang, 'title': title,
                             'summary_sha256': hashlib.sha256(summary.encode()).hexdigest(),
                             'body_sha256': hashlib.sha256(body.encode()).hexdigest()})
        sql.append(f"""INSERT INTO article_series_items(series_id,article_id,position)
SELECT s.id,a.id,{number*10} FROM article_series s,articles a
WHERE s.slug='chemistry' AND a.slug={qslug}
ON CONFLICT(series_id,article_id) DO UPDATE SET position=EXCLUDED.position;""")
    sql.append('COMMIT;')
    return '\n'.join(sql) + '\n', expected


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--sql', type=Path, required=True)
    parser.add_argument('--expected', type=Path, required=True)
    args = parser.parse_args()
    sql, expected = prepare()
    args.sql.write_text(sql, encoding='utf-8')
    args.expected.write_text(json.dumps(expected, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
    print(f'Prepared {len(expected)} pages without remote changes')


if __name__ == '__main__':
    main()
