# CREATE TABLE: схема, которая не принимает невозможное

_Лид (summary):_ **Создаём таблицы семейного бюджета и переносим часть проверок к данным: обязательные поля, уникальные имена, положительные суммы, допустимые виды категорий, даты и внешние ключи.**

## Зачем это нужно

Проверка только в интерфейсе защищает один путь ввода. Импорт, другой скрипт или ручная команда могут его обойти. Ограничение базы действует для каждого пути.

> **Образ.** Схема — форма для льда. Вода может прийти из любого кувшина, но невозможный кубик форма не создаст.

## Сразу целиком

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

## Разбор

`CREATE TABLE` задаёт контракт. `NOT NULL` требует значение. `UNIQUE` запрещает два одинаковых имени счёта. `CHECK` проверяет правило строки. `REFERENCES` не позволяет сослаться на несуществующий счёт или категорию.

В SQLite внешние ключи включаются отдельно для каждого соединения командой `PRAGMA foreign_keys = ON`. Это особенность SQLite, а не универсальная команда SQL.

`INTEGER PRIMARY KEY` в SQLite связан с внутренним идентификатором строки и получает значение автоматически. Не добавляйте `AUTOINCREMENT` без конкретной необходимости: он запрещает повторное использование старых номеров ценой дополнительной работы, но не делает обычный ключ «надёжнее».

Дата хранится как ISO-текст, потому что у SQLite нет отдельного типа `DATE`. В PostgreSQL для неё выбрали бы `date`, а для денег — целое число минимальных единиц или подходящий `numeric`, но не двоичный `float`.

## Карта урока

![Ограничения защищают таблицу на входе](/static/course/sql/map-schema-ru.svg)

## Скажите своими словами

1. Почему `NOT NULL` и `CHECK` решают разные задачи?
2. Что произойдёт при `account_id = 99`, если такого счёта нет?
3. Почему проверки даты две?

## Задание

Создайте таблицу `goals`: ключ, уникальное непустое название и положительная целевая сумма в тиынах. Срок может отсутствовать, а если указан — должен быть корректной ISO-датой.

## Куда это встанет в проекте

Полная схема находится в `course/sql-budget/schema.sql`. Все дальнейшие запросы полагаются на её гарантии.

## Ответы

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

## Источники

- [SQLite: CREATE TABLE](https://sqlite.org/lang_createtable.html)
- [SQLite: AUTOINCREMENT](https://sqlite.org/autoinc.html)

