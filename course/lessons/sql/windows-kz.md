# Терезелік функциялар: жолдарды жоғалтпайтын жинақталған қалдық

_Лид (summary):_ **Әр операциядан кейінгі қалдық пен шығын рейтингін бастапқы жолдарды сақтай отырып есептейміз.**

## Бұл не үшін қажет

GROUP BY айды бірнеше қорытындыға жинайды. Ақша қай сатып алудан кейін азайғанын табу үшін әр операция және оның жанындағы жинақталған нәтиже керек.

> **Бейне.** GROUP BY карталар бумасын бір сомаға жинайды; терезе ешбір картаны тастамай, әрқайсына ағымдағы қорытынды жазады.

## Алдымен тұтас көрініс

```sql
SELECT t.happened_on, t.id, c.kind, t.amount_tiyn / 100.0 AS amount_kzt,
       205000.0 + SUM(CASE WHEN c.kind='income' THEN t.amount_tiyn ELSE -t.amount_tiyn END)
           OVER (ORDER BY t.happened_on, t.id ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)
           / 100.0 AS balance_kzt
FROM transactions AS t
JOIN categories AS c ON c.id = t.category_id
ORDER BY t.happened_on, t.id;
```

## Талдау

OVER агрегатты терезелік функцияға айналдырады. Терезедегі ORDER BY есеп ретін, сыртқы ORDER BY көрсету ретін береді. ROWS шеңберін анық жазып, тең күнді id арқылы ажыратыңыз. PARTITION BY әр санатқа жеке терезе ашады; ROW_NUMBER, RANK және DENSE_RANK тең мәндерді әртүрлі өңдейді.

## Сабақ картасы

![Терезелік функциялар](/static/course/sql/map-windows-kz.svg)

## Өз сөзіңізбен айтыңыз

1. Терезе GROUP BY-дан қалай ерекшеленеді?
2. Сұрыптау неге екі рет жазылған?
3. PARTITION BY не істейді?

## Тапсырма

ROW_NUMBER() арқылы әр санаттағы шығындарды үлкенінен кішісіне нөмірлеңіз.

## Жобадағы орны

Жинақталған қалдық ай соңын ғана емес, ақша жетіспеген сәтті де көрсетеді.

## Жауаптар

```sql
SELECT c.name, t.amount_tiyn / 100.0 AS amount_kzt,
       ROW_NUMBER() OVER (
           PARTITION BY c.id ORDER BY t.amount_tiyn DESC, t.id
       ) AS expense_no
FROM transactions AS t
JOIN categories AS c ON c.id=t.category_id
WHERE c.kind='expense';
```

## Дереккөздер

- [SQLite: оконные функции](https://sqlite.org/windowfunctions.html)
