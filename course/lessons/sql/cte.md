# Подзапросы и WITH: сложный вопрос по шагам

_Лид (summary):_ **Разбиваем расчёт нормы сбережений на именованные части с CTE, сравниваем подзапросы и соединения и учимся не повторять одну формулу в нескольких местах.**

## Зачем это нужно

Когда запрос одновременно выбирает месяц, делит доходы и расходы и считает процент, длинная строка скрывает ход мысли. Имя промежуточного результата делает его проверяемым.

> **Образ.** CTE — прозрачные миски на кухне: сначала отдельно подготовили доход и расход, затем смешали. Можно остановиться и проверить содержимое каждой.

## Сразу целиком

```sql
WITH monthly AS (
    SELECT c.kind, SUM(t.amount_tiyn) AS amount_tiyn
    FROM transactions AS t
    JOIN categories AS c ON c.id = t.category_id
    WHERE t.happened_on >= '2026-09-01' AND t.happened_on < '2026-10-01'
    GROUP BY c.kind
), totals AS (
    SELECT COALESCE(SUM(amount_tiyn) FILTER (WHERE kind='income'), 0) AS income,
           COALESCE(SUM(amount_tiyn) FILTER (WHERE kind='expense'), 0) AS expense
    FROM monthly
)
SELECT income / 100.0 AS income_kzt,
       expense / 100.0 AS expense_kzt,
       ROUND(100.0 * (income - expense) / NULLIF(income, 0), 1) AS saving_rate
FROM totals;
```

## Разбор

`WITH name AS (...)` даёт результату запроса имя на время одной команды. Это не сохранённая таблица. Несколько CTE читаются сверху вниз как этапы расчёта.

`NULLIF(income, 0)` защищает деление: при нулевом доходе знаменатель становится `NULL`, и результат честно остаётся неизвестным. Подставить ноль означало бы выдать математически неверный процент.

Коррелированный подзапрос выполняет логически отдельный поиск для строки внешнего запроса. Он бывает самым ясным решением, но для массового расчёта `JOIN` или предварительная агрегация часто проще оптимизируются. Решение выбирают по ясности и плану, а не по лозунгу «подзапросы медленные».

## Карта урока

![CTE делит расчёт на проверяемые этапы](/static/course/sql/map-cte-ru.svg)

## Скажите своими словами

1. Чем CTE отличается от настоящей таблицы?
2. Зачем здесь `NULLIF`?
3. Почему нельзя объявить все подзапросы медленными заранее?

## Задание

Добавьте третий CTE, который вернёт долю расходов в доходе. При отсутствии дохода результат должен быть `NULL`.

## Куда это встанет в проекте

Именованные этапы станут основой итогового отчёта и позволят тестировать каждый расчёт отдельно.

## Ответы

```sql
WITH totals AS (
    SELECT SUM(amount_tiyn) FILTER (WHERE kind='income') AS income,
           SUM(amount_tiyn) FILTER (WHERE kind='expense') AS expense
    FROM transactions JOIN categories ON categories.id=category_id
), ratio AS (
    SELECT 100.0 * expense / NULLIF(income, 0) AS expense_rate FROM totals
)
SELECT expense_rate FROM ratio;
```

## Источники

- [SQLite: WITH](https://sqlite.org/lang_with.html)

