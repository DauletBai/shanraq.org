# Оконные функции: накопительный остаток без потери строк

_Лид (summary):_ **Считаем баланс после каждой операции и рейтинг расходов. Оконная функция видит соседние строки, но сохраняет каждое исходное событие.**

## Зачем это нужно

Группировка превращает весь месяц в несколько итогов. Для вопроса «после какой покупки остаток стал минимальным?» нужны и отдельные строки, и накопленный итог рядом с каждой.

> **Образ.** `GROUP BY` сворачивает колоду в сумму. Окно кладёт на каждую карту бегущий итог, не выбрасывая ни одной карты.

## Сразу целиком

```sql
SELECT t.happened_on, t.id, c.kind, t.amount_tiyn / 100.0 AS amount_kzt,
       205000.0 + SUM(CASE WHEN c.kind='income' THEN t.amount_tiyn ELSE -t.amount_tiyn END)
           OVER (ORDER BY t.happened_on, t.id ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)
           / 100.0 AS balance_kzt
FROM transactions AS t
JOIN categories AS c ON c.id = t.category_id
ORDER BY t.happened_on, t.id;
```

## Разбор

`OVER` превращает агрегат в оконную функцию. `ORDER BY` внутри окна задаёт порядок расчёта, внешний `ORDER BY` — порядок показа. Нужны оба.

Рамка `ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW` явно означает «от первой строки до текущей». Без явной рамки значение по умолчанию зависит от вида сортировки, а строки с одинаковым ключом могут считаться одной группой peers. Устойчивый второй ключ `id` снимает неоднозначность.

`PARTITION BY category_id` начал бы отдельное окно для каждой категории. `ROW_NUMBER`, `RANK` и `DENSE_RANK` по-разному обращаются с равными значениями; выбор — часть смысла отчёта.

## Карта урока

![Окно сохраняет строки и добавляет контекст](/static/course/sql/map-windows-ru.svg)

## Скажите своими словами

1. Чем окно отличается от `GROUP BY`?
2. Почему сортировка указана дважды?
3. Что делает `PARTITION BY`?

## Задание

Пронумеруйте расходы внутри каждой категории от самого крупного к меньшему через `ROW_NUMBER()`.

## Куда это встанет в проекте

Накопительный остаток показывает не только итог месяца, но и момент кассового разрыва.

## Ответы

```sql
SELECT c.name, t.amount_tiyn / 100.0 AS amount_kzt,
       ROW_NUMBER() OVER (
           PARTITION BY c.id ORDER BY t.amount_tiyn DESC, t.id
       ) AS expense_no
FROM transactions AS t
JOIN categories AS c ON c.id=t.category_id
WHERE c.kind='expense';
```

## Источники

- [SQLite: оконные функции](https://sqlite.org/windowfunctions.html)

