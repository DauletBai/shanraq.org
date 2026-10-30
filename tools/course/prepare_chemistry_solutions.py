#!/usr/bin/env python3
"""Prepare a guarded, all-or-nothing SQL release of chemistry lessons 50–58."""
import argparse
import hashlib
import json
from pathlib import Path

try:
    from .chemistry_solutions_data import LESSONS
    from .build_chemistry_particles import LESSON_DIR, slug
    from .prepare_chemistry import literal, read_lesson
except ImportError:
    from chemistry_solutions_data import LESSONS
    from build_chemistry_particles import LESSON_DIR, slug
    from prepare_chemistry import literal, read_lesson


def prepare():
    if len(LESSONS) != 9:
        raise ValueError('the solutions block must contain exactly nine lessons')
    routes = [(number, LESSONS[number - 50]['stem'], slug(number, LESSONS[number - 50]['stem']))
              for number in range(50, 59)]
    slug_list = ','.join(literal(item[2]) for item in routes)
    sql = ['BEGIN;', "SELECT pg_advisory_xact_lock(hashtext('shanraq-chemistry-course'));",
           f"""DO $guard$
BEGIN
  IF (SELECT count(*) FROM auth_users WHERE email='baimurza.daulet@gmail.com') <> 1 THEN
    RAISE EXCEPTION 'Course author missing';
  END IF;
  IF (SELECT count(*) FROM article_series_items i JOIN article_series s ON s.id=i.series_id
      JOIN articles a ON a.id=i.article_id WHERE s.slug='chemistry'
      AND a.status='published' AND i.position BETWEEN 10 AND 490 AND i.position % 10=0) <> 49 THEN
    RAISE EXCEPTION 'Previous chemistry blocks incomplete';
  END IF;
  IF (SELECT count(*) FROM article_translations t JOIN articles a ON a.id=t.article_id
      WHERE a.slug='chemistry-49-salt-quantitative-case-review' AND t.lang IN ('ru','kz','en')
      AND t.status='ready') <> 3 THEN
    RAISE EXCEPTION 'Previous chemistry lesson translations incomplete';
  END IF;
  IF (SELECT count(*) FROM article_series WHERE slug='chemistry' AND status='published') <> 1 THEN
    RAISE EXCEPTION 'Published chemistry course missing';
  END IF;
  IF (SELECT count(*) FROM article_translations t JOIN article_series_items i ON i.article_id=t.article_id
      JOIN article_series s ON s.id=i.series_id WHERE s.slug='chemistry' AND i.position BETWEEN 10 AND 490
      AND i.position % 10=0 AND t.lang IN ('ru','kz','en') AND t.status='ready') <> 147 THEN
    RAISE EXCEPTION 'Previous chemistry translations incomplete';
  END IF;
  IF EXISTS (SELECT 1 FROM articles a JOIN auth_users u ON u.id=a.author_id
             WHERE a.slug IN ({slug_list}) AND u.email <> 'baimurza.daulet@gmail.com') THEN
    RAISE EXCEPTION 'Solutions lesson slug belongs to another author';
  END IF;
  IF EXISTS (SELECT 1 FROM article_series_items i JOIN article_series s ON s.id=i.series_id
             JOIN articles a ON a.id=i.article_id WHERE a.slug IN ({slug_list}) AND s.slug <> 'chemistry') THEN
    RAISE EXCEPTION 'Solutions lesson belongs to another course';
  END IF;
  IF EXISTS (SELECT 1 FROM article_series_items i JOIN article_series s ON s.id=i.series_id
             WHERE s.slug='chemistry' AND i.position BETWEEN 500 AND 580
             AND i.position % 10 = 0 AND i.article_id NOT IN
                 (SELECT id FROM articles WHERE slug IN ({slug_list}))) THEN
    RAISE EXCEPTION 'Solutions lesson position occupied';
  END IF;
END $guard$;"""]
    try:
        from .build_chemistry_solutions import bridge_body
    except ImportError:
        from build_chemistry_solutions import bridge_body
    expected = []
    # The previous lesson must point into this block in the same transaction.
    # Only this transaction adds the onward link: block six remains releasable alone.
    for lang in ('ru', 'kz', 'en'):
        suffix = '' if lang == 'ru' else '-' + lang
        title, summary, body = read_lesson(LESSON_DIR / f'49-salt-quantitative-case-review{suffix}.md')
        body = bridge_body(lang)
        sql.append(f"""UPDATE article_translations SET body_md={literal(body)},updated_at=now()
WHERE article_id=(SELECT id FROM articles WHERE slug='chemistry-49-salt-quantitative-case-review')
AND lang='{lang}' AND status='ready';""")
        expected.append({'slug': 'chemistry-49-salt-quantitative-case-review', 'lang': lang, 'title': title,
                         'summary_sha256': hashlib.sha256(summary.encode()).hexdigest(),
                         'body_sha256': hashlib.sha256(body.encode()).hexdigest()})
    for number, stem, article_slug in routes:
        qslug = literal(article_slug)
        cover = f'/static/course/chemistry/visuals/{number:02d}-solutions-ru.svg'
        sql.append(f"""INSERT INTO articles(author_id,slug,original_lang,category,subcategory,cover_url,status,published_at)
SELECT id,{qslug},'ru','society','education',{literal(cover)},'published',now()
FROM auth_users WHERE email='baimurza.daulet@gmail.com'
ON CONFLICT(slug) DO UPDATE SET category='society',subcategory='education',cover_url=EXCLUDED.cover_url,
status='published',updated_at=now(),published_at=COALESCE(articles.published_at,now());""")
        for lang in ('ru', 'kz', 'en'):
            suffix = '' if lang == 'ru' else '-' + lang
            title, summary, body = read_lesson(LESSON_DIR / f'{stem}{suffix}.md')
            sql.append(f"""INSERT INTO article_translations(article_id,lang,title,summary,body_md,source,status,cover_url)
SELECT id,'{lang}',{literal(title)},{literal(summary)},{literal(body)},'ai','ready',{literal(f'/static/course/chemistry/visuals/{number:02d}-solutions-{lang}.svg')}
FROM articles WHERE slug={qslug}
ON CONFLICT(article_id,lang) DO UPDATE SET title=EXCLUDED.title,summary=EXCLUDED.summary,
body_md=EXCLUDED.body_md,source='ai',status='ready',cover_url=EXCLUDED.cover_url,updated_at=now();""")
            expected.append({'slug': article_slug, 'lang': lang, 'title': title,
                             'summary_sha256': hashlib.sha256(summary.encode()).hexdigest(),
                             'body_sha256': hashlib.sha256(body.encode()).hexdigest()})
        sql.append(f"""INSERT INTO article_series_items(series_id,article_id,position)
SELECT s.id,a.id,{number*10} FROM article_series s,articles a
WHERE s.slug='chemistry' AND a.slug={qslug}
ON CONFLICT(series_id,article_id) DO UPDATE SET position=EXCLUDED.position;""")
    sql.append(f"""DO $verify$
BEGIN
  IF (SELECT count(*) FROM article_series_items i JOIN article_series s ON s.id=i.series_id
      JOIN articles a ON a.id=i.article_id WHERE s.slug='chemistry' AND a.slug IN ({slug_list})
      AND a.status='published' AND i.position BETWEEN 500 AND 580 AND i.position % 10=0) <> 9 THEN
    RAISE EXCEPTION 'Solutions block course membership incomplete';
  END IF;
  IF (SELECT count(*) FROM article_translations t JOIN articles a ON a.id=t.article_id
      WHERE a.slug IN ({slug_list}) AND t.lang IN ('ru','kz','en') AND t.status='ready') <> 27 THEN
    RAISE EXCEPTION 'Solutions block translations incomplete';
  END IF;
END $verify$;""")
    sql.append('SET CONSTRAINTS ALL IMMEDIATE;')
    sql.append(f"""DO $queue$
BEGIN
  IF EXISTS (SELECT 1 FROM job_queue WHERE name='syndicate_telegram'
      AND payload->>'article_id' IN (SELECT id::text FROM articles WHERE slug IN ({slug_list}))) THEN
    RAISE EXCEPTION 'Course lessons entered news queue';
  END IF;
END $queue$;""")
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
