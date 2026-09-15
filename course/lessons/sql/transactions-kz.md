# Транзакциялар және UPDATE: бәрін өзгерту немесе ештеңені өзгертпеу

_Лид (summary):_ **Деректі қауіпсіз түзетіп, шоттар арасындағы аударымды бір атомдық әрекет ретінде орындаймыз.**

## Бұл не үшін қажет

Шоттар арасындағы аударым кемінде екі байланысты жазбадан тұрады. Ақша алу сақталып, қосу сақталмаса, база жоқ ақшаны көрсетеді.

> **Бейне.** Транзакция — екі қақпалы шлюз: қайық не толық өтеді, не су мен қайық бастапқы орнына қайтады.

## Алдымен тұтас көрініс

```sql
BEGIN IMMEDIATE;

UPDATE transactions
SET note = 'Utilities and rent'
WHERE id = 2;

SELECT changes() AS changed_rows;
COMMIT;
```

## Талдау

ACID атомдық, келісімділік, оқшаулану және тұрақтылықты білдіреді; нақты кепілдік ДҚБЖ-ға байланысты. UPDATE не DELETE алдында дәл сол WHERE бар SELECT орындап, кейін өзгерген жол санын тексеріңіз. COMMIT бекітеді, ROLLBACK аяқталмаған жұмысты қайтарады, SAVEPOINT бір бөлігін қайтарады. Өз шоттарыңыз арасындағы аударым отбасы шығыны емес.

## Сабақ картасы

![Транзакциялар және UPDATE](/static/course/sql/map-transactions-kz.svg)

## Өз сөзіңізбен айтыңыз

1. Атомдық деген не?
2. WHERE шартын өзгеріске дейін қалай тексереміз?
3. Ішкі аударым неге шығын емес?

## Тапсырма

Транзакция ішінде Transport лимитін 40 000 теңгеге өзгертіп, оқыңыз; ROLLBACK жасап, ескі мәннің қайтқанын тексеріңіз.

## Жобадағы орны

Байланысты бюджет өзгерістері бірге орындалып, бекітуге дейін қайтарыла алады.

## Жауаптар

```sql
BEGIN;
UPDATE categories SET monthly_limit_tiyn=4000000 WHERE name='Transport';
SELECT monthly_limit_tiyn FROM categories WHERE name='Transport';
ROLLBACK;
SELECT monthly_limit_tiyn FROM categories WHERE name='Transport';
```

## Дереккөздер

- [SQLite: транзакции](https://sqlite.org/lang_transaction.html)
- [SQLite: изоляция](https://sqlite.org/isolation.html)
