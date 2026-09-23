# Команда как проверенное намерение

_Лид (summary):_ **Урок 36. Отделить разбор слов запуска от действия органайзера и объяснить каждую ошибку.**

Предварительно: уроки 29–35. Если забылись `Option` и `Result`, повторите [уроки 31–32](/read/rust-31-option?lang=ru).

## Образ и опорная карта

В кассе сначала проверяют заполненный бланк, а затем выполняют просьбу. Так и органайзер сначала превращает слова в проверенную команду. **Разбор команды** — проверка названия действия и нужных ему данных. Пока список задач не меняется. Образ ограничен: программа не понимает смысл произвольного русского предложения; она принимает ровно описанные формы.

**Слова → проверка формы → `Result<Command, String>` → `Ok` с командой или `Err` с причиной → выполнение позже.** `Command` — собственное перечисление из урока 29. `Add(String)` хранит название; `List` не хранит данных. В `Option<&str>` ссылка на текст либо есть (`Some`), либо отсутствует (`None`); значение берём из уже знакомого урока 31. Кортеж из трёх значений позволяет `match` проверить действие, название и лишний аргумент вместе. `title.trim().is_empty()` проверяет, осталось ли что-нибудь после удаления пробелов по краям. `String::from(title)` делает собственную строку: исходные аргументы могут перестать существовать, а команда сохранит название. Ветка `_` охватывает все остальные формы и возвращает объяснение. `Ok(...)` и `Err(...)` здесь значения результата, а не печать.

## Разобрать, затем показать

Сохраните прежний `src/main.rs` проекта `organizer`, замените его целиком и выполните `cargo run` рядом с `Cargo.toml`. Это учебный разбор заранее заданных слов: реальные аргументы из урока 35 подключим при сборке полного органайзера. Других файлов и зависимостей нет. Предскажите результат перед запуском.

```rust
enum Command {
    Add(String),
    List,
}

fn parse(action: &str, title: Option<&str>, extra: Option<&str>) -> Result<Command, String> {
    match (action, title, extra) {
        ("list", None, None) => Ok(Command::List),
        ("add", Some(title), None) => {
            if title.trim().is_empty() {
                Err(String::from("Пустое название"))
            } else {
                Ok(Command::Add(String::from(title)))
            }
        }
        _ => Err(String::from("Неверная команда или число аргументов")),
    }
}

fn show(result: Result<Command, String>) {
    match result {
        Ok(Command::Add(title)) => println!("Добавить: {title}"),
        Ok(Command::List) => println!("Показать список"),
        Err(message) => println!("Ошибка: {message}"),
    }
}

fn main() {
    show(parse("add", Some("Читать Rust"), None));
    show(parse("list", None, None));
    show(parse("add", Some("   "), None));
}
```
```text
Добавить: Читать Rust
Показать список
Ошибка: Пустое название
```

`parse` лишь создаёт проверенную команду. `show` сейчас показывает её человеку; список задач появится в следующем уроке. Если вызвать `parse("list", Some("лишнее"), None)`, проверка формы вернёт ошибку, а не проигнорирует слово.

## Проверка понимания

1. Почему `parse` возвращает команду, а не сразу добавляет задачу?
2. Чем `None` отличается от `Some("")`?
3. Кто печатает сообщение об ошибке?

Разбор не меняет список; отсутствующего названия нет, а пустое передано; печатает `main`. Восстановите опорную карту без страницы.

## Задание

**Обязательное.** Добавьте в перечисление `Done(u32)` для действия `done` с одним номером. Разбирайте число через `parse::<u32>()`, отвергайте ноль и неверный текст, оставьте проверку лишнего аргумента. Проверьте `done 2`, `done 0` и `done нет`. Функция разбора по-прежнему не должна печатать.

## Подсказка

Новая ветка `("done", Some(text), None)` может вызвать `text.parse::<u32>()`. Обработайте `Ok(number)` и `Err(_)` через `match`; для нуля верните `Err`.

## Эталон после попытки

<!-- task-answer -->
```rust
enum Command {
    Add(String),
    List,
    Done(u32),
}

fn parse(action: &str, title: Option<&str>, extra: Option<&str>) -> Result<Command, String> {
    match (action, title, extra) {
        ("list", None, None) => Ok(Command::List),
        ("add", Some(title), None) => {
            if title.trim().is_empty() {
                Err(String::from("Пустое название"))
            } else {
                Ok(Command::Add(String::from(title)))
            }
        }
        ("done", Some(text), None) => match text.parse::<u32>() {
            Ok(0) => Err(String::from("Номер должен быть больше нуля")),
            Ok(number) => Ok(Command::Done(number)),
            Err(_) => Err(String::from("Нужен номер")),
        },
        _ => Err(String::from("Неверная команда или число аргументов")),
    }
}

fn show(result: Result<Command, String>) {
    match result {
        Ok(Command::Add(title)) => println!("Добавить: {title}"),
        Ok(Command::List) => println!("Показать список"),
        Ok(Command::Done(number)) => println!("Завершить: {number}"),
        Err(message) => println!("Ошибка: {message}"),
    }
}

fn main() {
    show(parse("add", Some("Читать Rust"), None));
    show(parse("list", None, None));
    show(parse("add", Some("   "), None));
    show(parse("done", Some("2"), None));
    show(parse("done", Some("0"), None));
    show(parse("done", Some("нет"), None));
}
```
```text
Добавить: Читать Rust
Показать список
Ошибка: Пустое название
Завершить: 2
Ошибка: Номер должен быть больше нуля
Ошибка: Нужен номер
```

[Официальные объяснения перечислений и `match`](https://doc.rust-lang.org/book/ch06-02-match.html).

[Предыдущий урок](/read/rust-35-arguments?lang=ru) · [Оглавление](/course/rust?lang=ru) · [Следующий урок](/read/rust-37-operations?lang=ru)
