# Ішкі сұрау және WITH: күрделі сұрақты қадамға бөлу

_Лид (summary):_ **Жинақ үлесін есептеуді атауы бар, тексерілетін CTE кезеңдеріне бөлеміз.**

## Бұл не үшін қажет

Бір ұзын формуланы қайталау тексеруді қиындатады. Аралық қадамның атауы қате қорытындының қайда кіргенін көрсетеді.

> **Бейне.** CTE — мөлдір ыдыстар: кіріс пен шығысты бөлек дайындап, тексеріп, содан кейін біріктіреміз.

## Алдымен тұтас көрініс

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

## Талдау

CTE бір команда уақытында ғана бар; ол сақталған кесте емес. Бірнеше CTE есеп кезеңдері ретінде жоғарыдан төмен оқылады. NULLIF(income,0) кіріс нөл болғанда жалған пайыз шығармайды. Корреляцияланған ішкі сұрау әрдайым баяу емес; ішкі сұрау, JOIN және алдын ала топтауды түсініктілік пен жоспарға қарап таңдаңыз.

## Сабақ картасы

![Ішкі сұрау және WITH](/static/course/sql/map-cte-kz.svg)

## Өз сөзіңізбен айтыңыз

1. CTE кестеден несімен ерекшеленеді?
2. Мұнда NULLIF не үшін керек?
3. Барлық ішкі сұрауды неге алдын ала баяу дей алмаймыз?

## Тапсырма

Шығынның кіріске пайызын, ал кіріс жоқта NULL қайтаратын CTE қосыңыз.

## Жобадағы орны

Атауы бар кезеңдер қорытынды есепке кіріп, әр есепті жеке тексеруге мүмкіндік береді.

## Жауаптар

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

## Дереккөздер

- [SQLite: WITH](https://sqlite.org/lang_with.html)
