# Кілттер және JOIN: операцияны шот пен санатқа байланыстыру

_Лид (summary):_ **Бюджет моделін қалыпқа келтіріп, кестелерді тұрақты кілтпен, жолдарды көбейтпей байланыстырамыз.**

## Бұл не үшін қажет

Әр түбіртекте санат атауын қайталау емле нұсқаларын және қайшы ережелерді туғызады. Мағынаны бір жерде сақтап, кілтіне сілтеңіз.

> **Бейне.** Операцияда киім ілгіштің жетоны бар; JOIN ұқсас түске емес, жетонның дәл нөміріне қарай киімді қайтарады.

## Алдымен тұтас көрініс

```sql
SELECT t.happened_on, a.name AS account, c.name AS category,
       t.amount_tiyn / 100.0 AS amount_kzt
FROM transactions AS t
JOIN accounts AS a ON a.id = t.account_id
JOIN categories AS c ON c.id = t.category_id
ORDER BY t.happened_on, t.id;
```

## Талдау

INNER JOIN сәйкестікті ғана қалдырады; LEFT JOIN сол жақтағы барлық жолды сақтап, жоқ оң мәнге NULL береді. Ұқсас мәтінмен емес, кілтпен байланыстырыңыз. Қорытындыға дейін кілт бірегейлігін және жол санын тексеріңіз: қайталанған анықтамалық кілті ақшаны көбейте алады.

## Сабақ картасы

![Кілттер және JOIN](/static/course/sql/map-join-kz.svg)

## Өз сөзіңізбен айтыңыз

1. INNER JOIN мен LEFT JOIN қалай ерекшеленеді?
2. Атаумен байланыстыру неге дұрыс емес?
3. JOIN ақшаны қалай көбейтіп жібереді?

## Тапсырма

Операциясы жоқтарын да қамтып, әр шот пен оның операция санын көрсетіңіз.

## Жобадағы орны

JOIN ықшам фактіге түсінікті атау қосады, ал сыртқы кілт байланысты қорғайды.

## Жауаптар

```sql
SELECT c.name, COALESCE(SUM(t.amount_tiyn), 0) / 100.0 AS spent_kzt
FROM categories AS c
LEFT JOIN transactions AS t
  ON t.category_id = c.id
 AND t.happened_on >= '2026-09-01'
 AND t.happened_on <  '2026-10-01'
WHERE c.kind = 'expense'
GROUP BY c.id, c.name
ORDER BY c.name;
```

## Дереккөздер

- [SQLite: JOIN в SELECT](https://sqlite.org/lang_select.html)
- [PostgreSQL: соединения таблиц](https://www.postgresql.org/docs/current/queries-table-expressions.html)
