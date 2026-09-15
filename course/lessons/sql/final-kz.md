# Индекс, EXPLAIN және қорытынды есеп

_Лид (summary):_ **Отбасы есебін жинап, сұрау жоспарын өлшеп, негізделген индекс құрып, сақтық көшірмені қалпына келтіру арқылы тексереміз.**

## Бұл не үшін қажет

Дұрыс сұрау көп жолда баяулауы, жылдам сұрау қате нәтиже беруі мүмкін. Тексерілмеген көшірме қажет сәтте ашылмауы ықтимал.

> **Бейне.** Индекс — кітаптың әліпбилік көрсеткіші: іздеуді жеделдету үшін орын және жаңарту жұмысы жұмсалады. Сақтық көшірме — шынайы есікте сыналған қосалқы кілт.

## Алдымен тұтас көрініс

```sql
EXPLAIN QUERY PLAN
SELECT category_id, SUM(amount_tiyn)
FROM transactions
WHERE happened_on >= '2026-09-01' AND happened_on < '2026-10-01'
GROUP BY category_id;

CREATE INDEX transactions_month_category
    ON transactions(happened_on, category_id);
```

## Талдау

Индексті өлшенген сұрау жолына сәйкестендіріңіз; әр бағанға индекс жасамаңыз және сегіз жолмен жылдамдық бағаламаңыз. EXPLAIN QUERY PLAN — нұсқаға тәуелді диагностика. Шағын кестеде SCAN дұрыс жоспар болуы мүмкін. VIEW нәтиже көшірмесін емес, сұрауды сақтайды. Ашық SQLite базасын .backup не backup API арқылы көшіріп, көшірмені бөлек ашыңыз, integrity_check орындап, бақылау сомаларын салыстырыңыз. Нақты қаржы дерегін Git-ке қоспаңыз.

## Сабақ картасы

![Индекс, EXPLAIN және қорытынды есеп](/static/course/sql/map-final-kz.svg)

## Өз сөзіңізбен айтыңыз

1. Индекстің құны қандай?
2. SCAN неліктен әрдайым жаман емес?
3. Сақтық көшірме қашан тексерілген саналады?

## Тапсырма

Көшірме жасап, оны бөлек ашыңыз; integrity_check пен айлық есепті орындап, үш қорытындыны салыстырыңыз.

## Жобадағы орны

Дайын жобада фактілер қорғалған, есеп түсіндіріледі, индекс өлшенген, көшірме қалпына келеді.

## Жауаптар

```text
sqlite3 budget.db ".read course/sql-budget/report.sql"
```

```sql
CREATE VIEW expense_by_category AS
SELECT category_id, SUM(amount_tiyn) AS amount_tiyn
FROM transactions
JOIN categories ON categories.id=category_id
WHERE categories.kind='expense'
GROUP BY category_id;
```

```text
sqlite3 budget.db ".backup budget-backup.db"
sqlite3 budget-backup.db "PRAGMA integrity_check;"
sqlite3 budget-backup.db ".read course/sql-budget/report.sql"
```

## Дереккөздер

- [SQLite: план запроса](https://sqlite.org/eqp.html)
- [SQLite: планировщик запросов](https://sqlite.org/queryplanner.html)
- [SQLite: резервное копирование](https://sqlite.org/backup.html)
