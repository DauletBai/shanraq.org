# CREATE TABLE: мүмкін емес деректі қабылдамайтын сұлба

_Лид (summary):_ **Отбасы бюджетінің кестелерін құрып, негізгі тексерулерді дерекқордың өзіне тапсырамыз.**

## Бұл не үшін қажет

Интерфейс тек бір енгізу жолын тексереді. Импорт пен басқа бағдарлама оны айналып өте алады; дерекқор шектеуі барлық жолды қорғайды.

> **Бейне.** Сұлба — мұзға арналған қалып: су кез келген ыдыстан құйылуы мүмкін, бірақ қалып мүмкін емес текше жасамайды.

## Алдымен тұтас көрініс

```sql
PRAGMA foreign_keys = ON;

CREATE TABLE transactions (
    id INTEGER PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(id),
    category_id INTEGER NOT NULL REFERENCES categories(id),
    happened_on TEXT NOT NULL CHECK (
        date(happened_on) IS NOT NULL AND happened_on = date(happened_on)
    ),
    amount_tiyn INTEGER NOT NULL CHECK (amount_tiyn > 0),
    note TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Талдау

NOT NULL мәнді міндеттейді; UNIQUE шот атауының қайталануын тоқтатады; CHECK бір жолдың шартын тексереді; REFERENCES жоқ шотқа не санатқа сілтеме жасатпайды. SQLite-та сыртқы кілттер әр қосылымда PRAGMA foreign_keys=ON арқылы қосылады. INTEGER PRIMARY KEY жол нөмірін өзі бере алады; AUTOINCREMENT көбіне қажет емес. SQLite-та жеке DATE типі болмағандықтан, күн ISO пішіміндегі TEXT ретінде сақталады.

## Сабақ картасы

![CREATE TABLE](/static/course/sql/map-schema-kz.svg)

## Өз сөзіңізбен айтыңыз

1. NOT NULL мен CHECK қандай бөлек міндет атқарады?
2. Жоқ шоттың account_id мәнін берсеңіз не болады?
3. Күн шарты неліктен екі тексеруден тұрады?

## Тапсырма

Кілті, қайталанбайтын бос емес атауы, оң target_tiyn сомасы және міндетті емес дұрыс ISO мерзімі бар goals кестесін құрыңыз.

## Жобадағы орны

Толық сұлба course/sql-budget/schema.sql файлында; кейінгі сұраулар оның кепілдіктеріне сүйенеді.

## Жауаптар

```sql
CREATE TABLE goals (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE CHECK (trim(name) <> ''),
    target_tiyn INTEGER NOT NULL CHECK (target_tiyn > 0),
    due_on TEXT CHECK (due_on IS NULL OR (
        date(due_on) IS NOT NULL AND due_on = date(due_on)
    ))
);
```

## Дереккөздер

- [SQLite: CREATE TABLE](https://sqlite.org/lang_createtable.html)
- [SQLite: AUTOINCREMENT](https://sqlite.org/autoinc.html)
