#!/usr/bin/env python3
"""Generate the localized support maps for the SQL family-budget course."""

from html import escape
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "web/static/course/sql"

MAPS = {
    "why": {
        "ru": ("ФАКТЫ → ВОПРОС → ОТВЕТ", ["FROM · чеки", "WHERE · месяц", "GROUP · конверты", "SUM · итог"]),
        "kz": ("ДЕРЕК → СҰРАҚ → ЖАУАП", ["FROM · түбіртек", "WHERE · ай", "GROUP · санат", "SUM · қорытынды"]),
        "en": ("FACTS → QUESTION → ANSWER", ["FROM · receipts", "WHERE · month", "GROUP · envelopes", "SUM · total"]),
    },
    "model": {
        "ru": ("ОДНА СТРОКА — ОДИН ФАКТ", ["счёт", "категория", "операция", "NULL ≠ 0"]),
        "kz": ("БІР ЖОЛ — БІР ДЕРЕК", ["шот", "санат", "операция", "NULL ≠ 0"]),
        "en": ("ONE ROW — ONE FACT", ["account", "category", "transaction", "NULL ≠ 0"]),
    },
    "schema": {
        "ru": ("СХЕМА — ФОРМА ДЛЯ ДАННЫХ", ["NOT NULL", "CHECK", "UNIQUE", "REFERENCES"]),
        "kz": ("СҰЛБА — ДЕРЕК ҚАЛЫБЫ", ["NOT NULL", "CHECK", "UNIQUE", "REFERENCES"]),
        "en": ("SCHEMA SHAPES DATA", ["NOT NULL", "CHECK", "UNIQUE", "REFERENCES"]),
    },
    "insert": {
        "ru": ("ЗАПРОС ≠ ЗНАЧЕНИЕ", ["SQL · бланк", "? · место", "параметр · значение", "никакой склейки"]),
        "kz": ("СҰРАУ ≠ МӘН", ["SQL · қалып", "? · орын", "параметр · мән", "жолды қоспау"]),
        "en": ("QUERY ≠ VALUE", ["SQL · form", "? · slot", "parameter · value", "never concatenate"]),
    },
    "select": {
        "ru": ("ОТ ВСЕХ СТРОК К НУЖНЫМ", ["FROM", "WHERE", "SELECT", "ORDER · LIMIT"]),
        "kz": ("БАРЛЫҚ ЖОЛДАН КЕРЕГІНЕ", ["FROM", "WHERE", "SELECT", "ORDER · LIMIT"]),
        "en": ("FROM ALL ROWS TO THE NEEDED", ["FROM", "WHERE", "SELECT", "ORDER · LIMIT"]),
    },
    "expressions": {
        "ru": ("ФАКТ ХРАНИМ · ВИД ВЫЧИСЛЯЕМ", ["тиын → тенге", "NULL → COALESCE", "условие → CASE", "дата → границы"]),
        "kz": ("ДЕРЕК САҚТАЛАДЫ · ТҮР ЕСЕПТЕЛЕДІ", ["тиын → теңге", "NULL → COALESCE", "шарт → CASE", "күн → шекара"]),
        "en": ("STORE FACTS · COMPUTE VIEWS", ["tiyn → tenge", "NULL → COALESCE", "condition → CASE", "date → bounds"]),
    },
    "aggregate": {
        "ru": ("ЧЕКИ → КОНВЕРТЫ → ИТОГ", ["WHERE · строки", "GROUP BY · группы", "SUM · сумма", "HAVING · итог"]),
        "kz": ("ТҮБІРТЕК → САНАТ → ҚОРЫТЫНДЫ", ["WHERE · жол", "GROUP BY · топ", "SUM · сома", "HAVING · қорытынды"]),
        "en": ("RECEIPTS → ENVELOPES → TOTAL", ["WHERE · rows", "GROUP BY · groups", "SUM · amount", "HAVING · total"]),
    },
    "join": {
        "ru": ("ЖЕТОН СОЕДИНЯЕТ ТОЧНО", ["операция · много", "category_id", "категория · одна", "проверить дубли"]),
        "kz": ("БЕЛГІ ДӘЛ БАЙЛАНЫСТЫРАДЫ", ["операция · көп", "category_id", "санат · бір", "қайталауды тексеру"]),
        "en": ("A KEY CONNECTS EXACTLY", ["transactions · many", "category_id", "category · one", "check duplicates"]),
    },
    "cte": {
        "ru": ("СЛОЖНОЕ — ПО ПРОЗРАЧНЫМ ШАГАМ", ["monthly", "totals", "ratio", "final SELECT"]),
        "kz": ("КҮРДЕЛІ СҰРАҚ — АШЫҚ ҚАДАММЕН", ["monthly", "totals", "ratio", "final SELECT"]),
        "en": ("COMPLEX QUERY · VISIBLE STEPS", ["monthly", "totals", "ratio", "final SELECT"]),
    },
    "windows": {
        "ru": ("СТРОКА ОСТАЁТСЯ · ИТОГ РАСТЁТ", ["PARTITION", "ORDER", "ROWS", "running SUM"]),
        "kz": ("ЖОЛ ҚАЛАДЫ · ҚОРЫТЫНДЫ ӨСЕДІ", ["PARTITION", "ORDER", "ROWS", "running SUM"]),
        "en": ("KEEP EACH ROW · GROW THE TOTAL", ["PARTITION", "ORDER", "ROWS", "running SUM"]),
    },
    "transactions": {
        "ru": ("ЛИБО ВСЁ · ЛИБО НИЧЕГО", ["BEGIN", "проверить", "COMMIT", "или ROLLBACK"]),
        "kz": ("НЕ БӘРІ · НЕ ЕШТЕҢЕ", ["BEGIN", "тексеру", "COMMIT", "немесе ROLLBACK"]),
        "en": ("ALL OR NOTHING", ["BEGIN", "verify", "COMMIT", "or ROLLBACK"]),
    },
    "audit": {
        "ru": ("ДЕТАЛИ ДОЛЖНЫ СХОДИТЬСЯ", ["строки", "дубли", "связи", "контрольный итог"]),
        "kz": ("БӨЛШЕКТЕР СӘЙКЕС БОЛУЫ КЕРЕК", ["жолдар", "қайталау", "байланыс", "бақылау сомасы"]),
        "en": ("DETAILS MUST RECONCILE", ["rows", "duplicates", "links", "control total"]),
    },
    "final": {
        "ru": ("ИЗМЕРИТЬ · УСКОРИТЬ · ВОССТАНОВИТЬ", ["EXPLAIN", "INDEX", "REPORT", "BACKUP → RESTORE"]),
        "kz": ("ӨЛШЕУ · ЖЕДЕЛДЕТУ · ҚАЛПЫНА КЕЛТІРУ", ["EXPLAIN", "INDEX", "REPORT", "BACKUP → RESTORE"]),
        "en": ("MEASURE · SPEED UP · RESTORE", ["EXPLAIN", "INDEX", "REPORT", "BACKUP → RESTORE"]),
    },
}


def svg(title, steps):
    cards = []
    for i, step in enumerate(steps):
        x = 55 + i * 245
        cards.append(f'<rect x="{x}" y="245" width="205" height="120" rx="20" fill="#fff" stroke="#d83232" stroke-width="4"/>')
        cards.append(f'<text x="{x+102}" y="312" text-anchor="middle" font-size="21" font-weight="700" fill="#303030">{escape(step)}</text>')
        if i < 3:
            cards.append(f'<path d="M{x+207} 305 H{x+238}" stroke="#2f7d45" stroke-width="7" marker-end="url(#a)"/>')
    return f'''<svg xmlns="http://www.w3.org/2000/svg" width="1100" height="620" viewBox="0 0 1100 620" role="img">
<defs><marker id="a" markerWidth="10" markerHeight="10" refX="7" refY="3" orient="auto"><path d="M0 0L0 6L8 3Z" fill="#2f7d45"/></marker></defs>
<rect width="1100" height="620" rx="32" fill="#f7f5f1"/>
<image href="/static/brand/shanraq.svg" x="492" y="57" width="116" height="116" preserveAspectRatio="xMidYMid meet"/>
<text x="550" y="205" text-anchor="middle" font-family="Arial,sans-serif" font-size="31" font-weight="800" fill="#303030">{escape(title)}</text>
{''.join(cards)}
<path d="M80 455 H1020" stroke="#e1ddd6" stroke-width="3"/>
<text x="550" y="515" text-anchor="middle" font-family="Arial,sans-serif" font-size="23" fill="#555">SQL</text>
</svg>'''


OUT.mkdir(parents=True, exist_ok=True)
for name, translations in MAPS.items():
    for lang, (title, steps) in translations.items():
        (OUT / f"map-{name}-{lang}.svg").write_text(svg(title, steps), encoding="utf-8")
print(f"generated {len(MAPS) * 3} maps in {OUT}")
