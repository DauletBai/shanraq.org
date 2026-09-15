# SUM, COUNT, GROUP BY и HAVING: итог по категориям

_Лид (summary):_ **Сжимаем множество операций в проверяемые итоги, считаем строки отдельно от значений и находим категории, вышедшие за лимит.**

## Зачем это нужно

Список чеков отвечает «что произошло». Бюджет требует другого уровня: сколько ушло на еду, какова доля обязательных расходов, где превышен план.

> **Образ.** `GROUP BY` раскладывает чеки по конвертам, агрегат считает содержимое каждого конверта, `HAVING` оставляет только конверты с нужным итогом.

## Сразу целиком

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

## Разбор

`COUNT(*)` считает строки. `COUNT(note)` считает только строки, где `note` не `NULL`. `SUM`, `AVG`, `MIN`, `MAX` игнорируют `NULL`; на пустом наборе `SUM` возвращает `NULL`, поэтому итог часто оборачивают в `COALESCE`.

`WHERE` отбирает строки до группировки. `HAVING` отбирает уже собранные группы. Условие по виду категории принадлежит `WHERE`, условие по сумме — `HAVING`.

В строгом переносимом SQL каждый неагрегированный столбец результата должен участвовать в `GROUP BY`. Некоторые СУБД допускают сокращения, но они могут вернуть произвольное значение из группы.

## Карта урока

![Строки собираются в группы и проходят HAVING](/static/course/sql/map-aggregate-ru.svg)

## Скажите своими словами

1. Чем `COUNT(*)` отличается от `COUNT(note)`?
2. Где поставить условие `kind = 'expense'`?
3. Что вернёт `SUM` для пустого набора?

## Задание

Для каждой расходной категории покажите лимит, сумму сентября и остаток лимита. Оставьте только категории, у которых было хотя бы две операции.

## Куда это встанет в проекте

Это центральная таблица отчёта: факт расходов сопоставляется с планом.

## Ответы

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

## Источники

- [SQLite: агрегатные функции](https://sqlite.org/lang_aggfunc.html)

