#!/usr/bin/env python3
"""Idempotently install or update the complete SQL course on production."""
import io,json,os,re,subprocess,sys
from pathlib import Path
ROOT=Path(__file__).resolve().parents[2]; LESSONS=ROOT/"course/lessons"; HOST=os.getenv("SHANRAQ_HOST","root@85.202.192.61")
PSQL="cd /opt/shanraq && docker compose -f docker-compose.prod.yml exec -T db psql -U shanraq -d shanraq -v ON_ERROR_STOP=1 -q"
LEAD=re.compile(r"_[^_]+:_\s*\*\*(.+)\*\*\s*$")
ORDER=["why","model","schema","insert","select","expressions","aggregate","join","cte","windows","transactions","audit","final"]
SLUGS=json.loads((ROOT/"tools/course/lesson-slugs.json").read_text())

def parsed(path):
 lines=path.read_text(encoding="utf-8").split("\n"); m=LEAD.match(lines[2].strip())
 if not lines[0].startswith("# ") or not m: raise SystemExit(f"bad lesson header: {path}")
 return lines[0][2:].strip(),m.group(1),"\n".join(lines[4:]).strip()+"\n"
def lit(value,tag):
 if f"${tag}$" in value: raise SystemExit(f"quoting marker in {tag}")
 return f"${tag}${value}${tag}$"
def remote(script=None,query=None):
 cmd=["ssh",HOST,PSQL+(f' -tAc "{query}"' if query else "")]
 r=subprocess.run(cmd,input=script,text=True,capture_output=True)
 if r.returncode: raise SystemExit(r.stderr.strip() or "remote psql failed")
 return r.stdout.strip()

def main():
 if "--check" in sys.argv:
  got=remote(query="select count(*) from article_series_items i join article_series s on s.id=i.series_id where s.slug='sql'")
  print(f"production SQL lessons: {got or 0}/13"); return 0 if got=="13" else 1
 sql=["BEGIN;", "INSERT INTO article_series(slug,cover_url,status,code_lang) VALUES('sql','/static/covers/it/sql-family-budget.webp','published','sql') ON CONFLICT(slug) DO UPDATE SET cover_url=EXCLUDED.cover_url,status='published',code_lang='sql',updated_at=now();"]
 meta={"ru":("SQL: семейный бюджет без догадок","13 практических уроков: от первой таблицы до проверяемого отчёта и резервной копии. Работает на компьютере, планшете и смартфоне."),"kz":("SQL: отбасы бюджетін болжамсыз есептеу","Алғашқы кестеден тексерілетін есеп пен сақтық көшірмеге дейінгі 13 тәжірибелік сабақ. Компьютерде, планшетте және смартфонда жұмыс істейді."),"en":("SQL: a family budget without guesswork","13 practical lessons from a first table to an auditable report and verified backup, on computers, tablets, and phones.")}
 for lang,(title,summary) in meta.items(): sql.append(f"INSERT INTO article_series_i18n(series_id,lang,title,summary) SELECT id,'{lang}',{lit(title,'st')},{lit(summary,'ss')} FROM article_series WHERE slug='sql' ON CONFLICT(series_id,lang) DO UPDATE SET title=EXCLUDED.title,summary=EXCLUDED.summary;")
 for pos,stem in enumerate(ORDER,1):
  slug=SLUGS[f"sql/{stem}"][0]
  sql.append(f"INSERT INTO articles(author_id,slug,original_lang,category,subcategory,cover_url,status,published_at) SELECT id,'{slug}','ru','it','programming','/static/covers/it/sql-family-budget.webp','published',now() FROM auth_users WHERE email='baimurza.daulet@gmail.com' ON CONFLICT(slug) DO UPDATE SET status='published',cover_url=EXCLUDED.cover_url,updated_at=now(),published_at=COALESCE(articles.published_at,now());")
  for suffix,lang in (("","ru"),("-kz","kz"),("-en","en")):
   title,summary,body=parsed(LESSONS/f"sql/{stem}{suffix}.md")
   sql.append(f"INSERT INTO article_translations(article_id,lang,title,summary,body_md,source,status) SELECT id,'{lang}',{lit(title,'title')},{lit(summary,'summary')},{lit(body,'body')},'human','ready' FROM articles WHERE slug='{slug}' ON CONFLICT(article_id,lang) DO UPDATE SET title=EXCLUDED.title,summary=EXCLUDED.summary,body_md=EXCLUDED.body_md,source='human',status='ready',updated_at=now();")
  sql.append(f"INSERT INTO article_series_items(series_id,article_id,position) SELECT s.id,a.id,{pos} FROM article_series s,articles a WHERE s.slug='sql' AND a.slug='{slug}' ON CONFLICT(series_id,article_id) DO UPDATE SET position=EXCLUDED.position;")
 sql.append("COMMIT;"); remote("\n".join(sql)+"\n"); print("SQL course installed: 13 lessons × 3 languages")
 return 0
if __name__=="__main__": raise SystemExit(main())
