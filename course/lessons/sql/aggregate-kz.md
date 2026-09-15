# SUM, COUNT, GROUP BY және HAVING: санат қорытындысы

_Лид (summary):_ **Операцияларды тексерілетін қорытындыға жинап, лимиттен асқан санаттарды табамыз.**

## Бұл не үшін қажет

Түбіртек тізімі әр санатқа қанша кеткенін айтпайды. Агрегат жауапты есептейді және есептеу ережесін ашық қалдырады.

> **Бейне.** GROUP BY түбіртектерді конверттерге бөледі; агрегат әр конвертті есептейді; HAVING қорытындысына қарай конвертті қалдырады.

## Алдымен тұтас көрініс

```sql
SELECT c.name,
       COUNT(*) AS operation_count,
       SUM(t.amount_tiyn) / 100.0 AS spent_kzt
FROM transactions AS t
JOIN categories AS c ON c.id = t.category_id
WHERE c.kind = 'expense'
GROUP BY c.id, c.name
HAVING SUM(t.amount_tiyn) >= 1000000
ORDER BY SUM(t.amount_tiyn) DESC;
```

## Талдау

WHERE топтауға дейін бастапқы жолдарды, HAVING топталған нәтижені сүзеді. COUNT(*) жолдарды, COUNT(column) NULL емес мәндерді санайды, SUM қосады. Агрегат емес таңдалған өрнек топты анықтауы керек. Оқылатын атаумен бірге тұрақты санат id бойынша топтаңыз.

## Сабақ картасы

![SUM, COUNT, GROUP BY және HAVING](/static/course/sql/map-aggregate-kz.svg)

## Өз сөзіңізбен айтыңыз

1. WHERE мен HAVING қашан іске қосылады?
2. COUNT(*) пен COUNT(note) қалай ерекшеленеді?
3. Неліктен санат id бойынша топтаймыз?

## Тапсырма

Қыркүйек шығыны monthly_limit_tiyn мәнінен асқан санаттарды көрсетіңіз.

## Жобадағы орны

Бұл қорытындылар бюджет панелі мен лимит ескертуінің негізі болады.

## Жауаптар

```sql
SELECT c.name,
       c.monthly_limit_tiyn / 100.0 AS limit_kzt,
       SUM(t.amount_tiyn) / 100.0 AS spent_kzt,
       (c.monthly_limit_tiyn - SUM(t.amount_tiyn)) / 100.0 AS left_kzt
FROM categories AS c
JOIN transactions AS t ON t.category_id = c.id
WHERE c.kind = 'expense'
  AND t.happened_on >= '2026-09-01' AND t.happened_on < '2026-10-01'
GROUP BY c.id, c.name, c.monthly_limit_tiyn
HAVING COUNT(*) >= 2;
```

## Дереккөздер

- [SQLite: агрегатные функции](https://sqlite.org/lang_aggfunc.html)
