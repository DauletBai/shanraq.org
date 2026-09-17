#!/usr/bin/env python3
"""Prepare a reviewable, transactional Rust publication. Does not contact a server."""
import argparse
import hashlib
import json
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parents[2]
LEAD = re.compile(r'_[^_]+:_\s*\*\*(.+)\*\*\s*$')


def literal(value):
    if '$rust_course$' in value:
        raise ValueError('SQL quoting delimiter in course text')
    return '$rust_course$' + value + '$rust_course$'


def prepare():
    manifest = json.loads((ROOT / 'tools/course/rust-release.json').read_text(encoding='utf-8'))
    items = manifest['items']
    count = manifest['lessons']
    assert manifest['course'] == 'rust' and manifest['code_lang'] == 'rust'
    assert isinstance(count, int) and 5 <= count <= 60 and count % 5 == 0
    assert manifest['batch'] == count // 5
    assert [i['position'] for i in items] == [1] + list(range(10, count * 10 + 1, 10))
    assert len({i['slug'] for i in items}) == count + 1
    assert items[0]['slug'] == 'rust-before-start'
    assert all(i['slug'] == 'rust-' + i['stem'] for i in items[1:])
    expected = []
    sql = ["BEGIN;", "SELECT pg_advisory_xact_lock(hashtext('shanraq-rust-course'));"]
    sql.append("""DO $guard$
BEGIN
  IF (SELECT count(*) FROM auth_users WHERE email = 'baimurza.daulet@gmail.com') <> 1 THEN
    RAISE EXCEPTION 'Expected course author not found';
  END IF;
  IF EXISTS (SELECT 1 FROM articles a JOIN auth_users u ON u.id = a.author_id
             WHERE a.slug IN (SLUGS) AND u.email <> 'baimurza.daulet@gmail.com') THEN
    RAISE EXCEPTION 'Rust slug belongs to another author';
  END IF;
  IF EXISTS (SELECT 1 FROM article_series_items i
             JOIN article_series s ON s.id = i.series_id
             JOIN articles a ON a.id = i.article_id
             WHERE a.slug IN (SLUGS) AND s.slug <> 'rust') THEN
    RAISE EXCEPTION 'Rust lesson already belongs to another course';
  END IF;
END $guard$;""".replace('SLUGS', ','.join(literal(i['slug']) for i in items)))
    cover = literal(manifest['cover_url'])
    sql.append(f"INSERT INTO article_series(slug,cover_url,status,code_lang) VALUES('rust',{cover},'published','rust') ON CONFLICT(slug) DO UPDATE SET cover_url=EXCLUDED.cover_url,status='published',code_lang='rust',updated_at=now();")
    for lang in ('kz', 'ru', 'en'):
        title, summary = manifest['metadata'][lang]
        sql.append(f"INSERT INTO article_series_i18n(series_id,lang,title,summary) SELECT id,'{lang}',{literal(title)},{literal(summary)} FROM article_series WHERE slug='rust' ON CONFLICT(series_id,lang) DO UPDATE SET title=EXCLUDED.title,summary=EXCLUDED.summary;")
    for item in items:
        slug = literal(item['slug'])
        sql.append(f"INSERT INTO articles(author_id,slug,original_lang,category,subcategory,cover_url,status,published_at) SELECT id,{slug},'ru','it','programming',{cover},'published',now() FROM auth_users WHERE email='baimurza.daulet@gmail.com' ON CONFLICT(slug) DO UPDATE SET status='published',cover_url=EXCLUDED.cover_url,updated_at=now(),published_at=COALESCE(articles.published_at,now());")
        for lang, suffix in [('ru', ''), ('kz', '-kz'), ('en', '-en')]:
            file = ROOT / f"course/lessons/rust/{item['stem']}{suffix}.md"
            lines = file.read_text(encoding='utf-8').splitlines()
            lead = LEAD.fullmatch(lines[2])
            assert lines[0].startswith('# ') and lead, f'Invalid header: {file}'
            title, summary = lines[0][2:], lead[1]
            body = '\n'.join(lines[4:]).strip() + '\n'
            assert not re.search(r'\]\([^)]*\.md\)', body), f'Repository link on public page: {file}'
            assert not any(word in body for word in ('Редакционный черновик', 'Editorial draft.', 'Редакциялық жоба.'))
            for target in re.findall(r'\]\((/read/[^)]+)\)', body):
                assert target.split('?')[0][6:] in {i['slug'] for i in items}, target
                assert target.endswith('?lang='+lang), target
            # Source remains AI-labelled: these are original assistant-written lessons,
            # not imported human-authored texts or an unmarked translation service.
            sql.append(f"INSERT INTO article_translations(article_id,lang,title,summary,body_md,source,status) SELECT id,'{lang}',{literal(title)},{literal(summary)},{literal(body)},'ai','ready' FROM articles WHERE slug={slug} ON CONFLICT(article_id,lang) DO UPDATE SET title=EXCLUDED.title,summary=EXCLUDED.summary,body_md=EXCLUDED.body_md,source='ai',status='ready',updated_at=now();")
            expected.append(dict(slug=item['slug'], lang=lang, title=title,
                                 body_sha256=hashlib.sha256(body.encode()).hexdigest(),
                                 summary_sha256=hashlib.sha256(summary.encode()).hexdigest()))
        sql.append(f"INSERT INTO article_series_items(series_id,article_id,position) SELECT s.id,a.id,{item['position']} FROM article_series s,articles a WHERE s.slug='rust' AND a.slug={slug} ON CONFLICT(series_id,article_id) DO UPDATE SET position=EXCLUDED.position;")
    sql.append("COMMIT;")
    return '\n'.join(sql) + '\n', expected


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--sql', required=True, type=Path)
    parser.add_argument('--expected', required=True, type=Path)
    args = parser.parse_args()
    sql, expected = prepare()
    args.sql.write_text(sql, encoding='utf-8')
    args.expected.write_text(json.dumps(expected, ensure_ascii=False, indent=2)+'\n', encoding='utf-8')
    print(f'Prepared {len(expected)//3-1} lessons + preface, {len(expected)} localized pages; no remote changes')


if __name__ == '__main__':
    main()
