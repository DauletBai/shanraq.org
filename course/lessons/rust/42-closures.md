# Замыкания: выбрать нужные задачи

_Лид (summary):_ **Урок 42. Найти задачи по слову и вывести их в устойчивом порядке.**

Предварительно: [итераторы](/read/rust-41-iterators?lang=ru), [строки](/read/rust-22-strings?lang=ru), [устойчивые ID](/read/rust-38-ids?lang=ru). Сохраните копию предыдущего `main.rs`: пример заменяет файл целиком и пока работает с данными в памяти.

## Образ и опорная карта

Представьте помощника в библиотеке: ему дали короткое правило «отбери карточки, где есть слово Rust». Это правило похоже на **замыкание** — безымянную функцию, которая может помнить значение из окружающего места. В `|task| task.title.contains(query)` вертикальные черты ограничивают входной параметр `task`, а `query` взят из внешней переменной. Образ не обещает, что помощник сделает работу сразу: итератор ленив и выполняет отбор, когда его потребляют.

**Все задачи → `filter` (оставить подходящие) → `collect` (собрать) → `sort_by_key` (упорядочить) → показать.** `filter` принимает замыкание с условием и сохраняет только подходящие элементы. Поскольку `iter()` выдаёт `&Task`, а `filter` берёт этот элемент по ссылке для проверки, параметр внутри `filter` имеет вид `&&Task` — ссылка на ссылку. Rust сам проходит эти два слоя при обращении к `.title`; то же происходит с `.id` в сортировке и `.title` в `map`. `contains` ищет точную последовательность символов: регистр важен, поэтому `rust` не совпадёт с `Rust`. `collect` создаёт здесь `Vec<&Task>`, то есть новый список ссылок; исходные задачи остаются на месте. Двоеточие в `let found: Vec<&Task>` сообщает компилятору, какой результат нужен. `sort_by_key(|task| task.id)` сравнивает ID, не названия. Для вывода `map` превращает каждую ссылку на задачу в `&str` её названия; `as_str()` заимствует текст, не копируя его. `names.len()` считает результат. Здесь `names` создан специально для знакомства с `map`; если нужно только число, достаточно `found.len()` и новый список не требуется. Вертикальные черты не отделяют командные аргументы, как символ `|` в подсказке урока 40: здесь они принадлежат синтаксису Rust.

## Запустите поиск

Замените `src/main.rs` и выполните `cargo run`. Вывод программы показан отдельно от сообщений Cargo. Новых библиотек нет.

```rust
struct Task {
    id: u32,
    title: String,
}

fn main() {
    let tasks = vec![
        Task {
            id: 3,
            title: String::from("Читать Rust"),
        },
        Task {
            id: 2,
            title: String::from("Купить хлеб"),
        },
        Task {
            id: 1,
            title: String::from("Повторить Rust"),
        },
    ];
    let query = "Rust";
    let mut found: Vec<&Task> = tasks
        .iter()
        .filter(|task| task.title.contains(query))
        .collect();
    found.sort_by_key(|task| task.id);
    let names: Vec<&str> = found.iter().map(|task| task.title.as_str()).collect();
    println!("Найдено: {}", names.len());
    for task in found {
        println!("{}: {}", task.id, task.title);
    }
}
```
```text
Найдено: 2
1: Повторить Rust
3: Читать Rust
```

Задача с ID 3 была в исходном списке раньше задачи с ID 1, но вывод отсортирован по ID. Исходный `tasks` не изменён. Если результатов нет, `found` пуст и цикл ничего не печатает; это корректный результат, а не ошибка. Для больших списков поиск каждого слова снова проходит все задачи; ускорение пока не цель урока.

## Вопросы для воспроизведения

1. Откуда замыкание получает `query`?
2. Почему `collect` не удаляет исходные задачи?
3. Изменит ли `sort_by_key` устойчивые ID?

## Задание

**Обязательное.** После основного вывода посчитайте задачи, чьё название содержит слово `хлеб`, и выведите `По слову «хлеб»: 1`. Используйте новый `filter` и `count`, не меняя исходный список. Затем попробуйте слово `Хлеб` и объясните разницу.

## Ответы

Создайте `second_query`. Выражение `tasks.iter().filter(|task| task.title.contains(second_query)).count()` потребляет отбор и возвращает число.

## Эталон после попытки

<!-- task-answer -->
```rust
struct Task {
    id: u32,
    title: String,
}

fn main() {
    let tasks = vec![
        Task {
            id: 3,
            title: String::from("Читать Rust"),
        },
        Task {
            id: 2,
            title: String::from("Купить хлеб"),
        },
        Task {
            id: 1,
            title: String::from("Повторить Rust"),
        },
    ];
    let query = "Rust";
    let mut found: Vec<&Task> = tasks
        .iter()
        .filter(|task| task.title.contains(query))
        .collect();
    found.sort_by_key(|task| task.id);
    let names: Vec<&str> = found.iter().map(|task| task.title.as_str()).collect();
    println!("Найдено: {}", names.len());
    for task in found {
        println!("{}: {}", task.id, task.title);
    }
    let second_query = "хлеб";
    let second_count = tasks
        .iter()
        .filter(|task| task.title.contains(second_query))
        .count();
    println!("По слову «{second_query}»: {second_count}");
}
```
```text
Найдено: 2
1: Повторить Rust
3: Читать Rust
По слову «хлеб»: 1
```

## После проверки

`sort_by_key` сортирует найденные ссылки; ID в самих задачах не меняются. [Официальная глава о замыканиях](https://doc.rust-lang.org/book/ch13-01-closures.html) · [глава об итераторах](https://doc.rust-lang.org/book/ch13-02-iterators.html). Если неясно, почему обход откладывается, повторите [урок 41](/read/rust-41-iterators?lang=ru).

[Предыдущий урок](/read/rust-41-iterators?lang=ru) · [Оглавление](/course/rust?lang=ru)
