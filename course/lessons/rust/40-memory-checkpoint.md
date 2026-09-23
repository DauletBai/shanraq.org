# Контрольная точка: органайзер в памяти

_Лид (summary):_ **Урок 40. Запустить интерактивный органайзер, выполнить несколько команд за один сеанс и обработать неверный ввод.**

Предварительно: уроки 16, 21, 27–28 и 31–39. При трудностях с концом ввода повторите [урок 34](/read/rust-34-input?lang=ru), с ID — [урок 38](/read/rust-38-ids?lang=ru).

## Образ и опорная карта

У стойки регистрации посетитель может подать несколько заявок подряд, а сотрудник помнит уже принятые до закрытия стойки. Так и этот органайзер принимает строки до команды выхода или конца ввода и хранит задачи в памяти одного запуска. Образ ограничен: после закрытия программы данные пока исчезают; сохранение в файл будет позднее.

**Строка ввода → разбор `Command` → действие над `Organizer` → ответ → следующая строка; `quit` или конец ввода → завершение.** Здесь соединяются `read_line` из урока 34, перечисление команд из 36, операции с задачами из 37, устойчивые ID из 38 и проверка ошибок из 39. `split_once(' ')` разделяет текст **по первому пробелу** и возвращает `Some((левая часть, правая часть))` или `None`, если пробела нет. Например, `add Читать Rust` даёт действие `add` и целое название `Читать Rust`. Кавычки здесь не нужны: строка уже прочитана целиком. Если написать кавычки, они войдут в название; разделителем служит обычный пробел, а не табуляция. `trim()` удаляет крайние пробелы и перевод строки. `self.tasks.is_empty()` проверяет, нет ли в списке задач. Действия `done` и `delete` требуют положительный ID; `list` и `quit` не принимают данные. При любой ошибке выводится объяснение, затем ожидается следующая строка. `use std::io` сокращает путь к модулю ввода: вместо `std::io::stdin()` пишем `io::stdin()`. Это тот же объект ввода из урока 34; подробное устройство модулей разберём в уроке 44. `!execute(...)` означает «не продолжать», потому что `execute` возвращает `false` для `quit`.

## Несколько команд за один запуск

Сохраните старый `src/main.rs`, замените его целиком и выполните `cargo run` рядом с `Cargo.toml`. Вводите команды по одной строке, нажимая Enter. Для выхода напишите `quit`; закрытие потока ввода тоже завершит программу. Новых файлов и зависимостей нет. Ниже показан вывод при сразу закрытом потоке ввода, как в автоматической проверке.

```rust
use std::io;

enum Command {
    Add(String),
    List,
    Done(u32),
    Delete(u32),
    Quit,
}

fn parse_id(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(0) => Err(String::from("Нужен положительный ID")),
        Ok(id) => Ok(id),
        Err(_) => Err(String::from("Нужен положительный ID")),
    }
}

fn parse(line: &str) -> Result<Command, String> {
    let line = line.trim();
    if line == "list" {
        return Ok(Command::List);
    }
    if line == "quit" {
        return Ok(Command::Quit);
    }
    match line.split_once(' ') {
        Some(("add", title)) => {
            if title.trim().is_empty() {
                Err(String::from("Название пусто"))
            } else {
                Ok(Command::Add(String::from(title.trim())))
            }
        }
        Some(("done", id)) => Ok(Command::Done(parse_id(id)?)),
        Some(("delete", id)) => Ok(Command::Delete(parse_id(id)?)),
        _ => Err(String::from("Неизвестная команда или неверные аргументы")),
    }
}

struct Task {
    id: u32,
    title: String,
    done: bool,
}
struct Organizer {
    tasks: Vec<Task>,
    next_id: u32,
}
impl Organizer {
    fn new() -> Self {
        Self {
            tasks: Vec::new(),
            next_id: 1,
        }
    }
    fn add(&mut self, title: String) -> Result<u32, String> {
        if title.trim().is_empty() {
            return Err(String::from("Название пусто"));
        }
        match self.next_id.checked_add(1) {
            Some(next) => {
                let id = self.next_id;
                self.tasks.push(Task {
                    id,
                    title,
                    done: false,
                });
                self.next_id = next;
                Ok(id)
            }
            None => Err(String::from("Номера закончились")),
        }
    }
    fn done(&mut self, id: u32) -> Result<(), String> {
        for task in &mut self.tasks {
            if task.id == id {
                task.done = true;
                return Ok(());
            }
        }
        Err(String::from("ID не найден"))
    }
    fn delete(&mut self, id: u32) -> Result<(), String> {
        for index in 0..self.tasks.len() {
            if self.tasks[index].id == id {
                self.tasks.remove(index);
                return Ok(());
            }
        }
        Err(String::from("ID не найден"))
    }
    fn list(&self) {
        if self.tasks.is_empty() {
            println!("Задач нет");
        }
        for task in &self.tasks {
            let mark = if task.done {
                "готово"
            } else {
                "ожидает"
            };
            println!("{}: {} ({mark})", task.id, task.title);
        }
    }
}

fn execute(organizer: &mut Organizer, command: Command) -> bool {
    match command {
        Command::Add(title) => {
            match organizer.add(title) {
                Ok(id) => println!("Добавлено: {id}"),
                Err(message) => println!("Ошибка: {message}"),
            }
            true
        }
        Command::List => {
            organizer.list();
            true
        }
        Command::Done(id) => {
            match organizer.done(id) {
                Ok(()) => println!("Завершено: {id}"),
                Err(message) => println!("Ошибка: {message}"),
            }
            true
        }
        Command::Delete(id) => {
            match organizer.delete(id) {
                Ok(()) => println!("Удалено: {id}"),
                Err(message) => println!("Ошибка: {message}"),
            }
            true
        }
        Command::Quit => false,
    }
}

fn main() {
    let mut organizer = Organizer::new();
    let stdin = io::stdin();
    println!("Команды: add НАЗВАНИЕ | list | done ID | delete ID | quit");
    loop {
        let mut line = String::new();
        match stdin.read_line(&mut line) {
            Ok(0) => break,
            Ok(_) => match parse(&line) {
                Ok(command) => {
                    if !execute(&mut organizer, command) {
                        break;
                    }
                }
                Err(message) => println!("Ошибка: {message}"),
            },
            Err(_) => {
                println!("Ошибка: не удалось прочитать ввод");
                break;
            }
        }
    }
    println!("Сеанс завершён");
}
```
```text
Команды: add НАЗВАНИЕ | list | done ID | delete ID | quit
Сеанс завершён
```

Вручную попробуйте `add Читать Rust`, `add Повторить тесты`, `list`, `done 1`, `delete 2`, `list`, `done 99`, `quit`. ID 1 останется с отметкой «готово», а неизвестный ID даст ошибку. Повторный запуск начнётся с пустого списка: это ожидаемая граница этапа без файлов.

При указанной последовательности команд вывод будет таким (введённые строки здесь не повторяются):

```text
Команды: add НАЗВАНИЕ | list | done ID | delete ID | quit
Добавлено: 1
Добавлено: 2
1: Читать Rust (ожидает)
2: Повторить тесты (ожидает)
Завершено: 1
Удалено: 2
1: Читать Rust (готово)
Ошибка: ID не найден
Сеанс завершён
```

## Проверка понимания

1. Почему `add Читать Rust` сохраняет название с пробелом?
2. Почему при `done 99` список должен остаться прежним?
3. Сохранятся ли задачи после `quit` и следующего запуска?

`split_once` отделяет только первый пробел; неизвестный ID не меняет записи; без файла данные исчезают. Восстановите опорную карту без страницы.

## Задание

**Обязательное.** Добавьте команду `rename ID НАЗВАНИЕ`: она должна искать задачу по устойчивому ID, отклонять пустое название и не менять список при неизвестном ID. Проверьте последовательность `add Читать Rust`, `rename 1 Читать две главы`, `list`, `rename 99 Ошибка`, `quit`. Затем заново запустите программу и убедитесь, что список пуст. Не используйте позицию в `Vec` как ID.

## Подсказка

В `Command` добавьте `Rename(u32, String)`. Для строки после `rename ` ещё раз вызовите `split_once(' ')`, отдельно проверьте ID и название. В методе `rename` проходите по `&mut self.tasks` и меняйте поле лишь после проверки.

## Эталон после попытки

<!-- task-answer -->
```rust
use std::io;

enum Command {
    Add(String),
    List,
    Done(u32),
    Delete(u32),
    Rename(u32, String),
    Quit,
}

fn parse_id(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(0) => Err(String::from("Нужен положительный ID")),
        Ok(id) => Ok(id),
        Err(_) => Err(String::from("Нужен положительный ID")),
    }
}

fn parse(line: &str) -> Result<Command, String> {
    let line = line.trim();
    if line == "list" {
        return Ok(Command::List);
    }
    if line == "quit" {
        return Ok(Command::Quit);
    }
    match line.split_once(' ') {
        Some(("add", title)) => {
            if title.trim().is_empty() {
                Err(String::from("Название пусто"))
            } else {
                Ok(Command::Add(String::from(title.trim())))
            }
        }
        Some(("done", id)) => Ok(Command::Done(parse_id(id)?)),
        Some(("delete", id)) => Ok(Command::Delete(parse_id(id)?)),
        Some(("rename", rest)) => match rest.split_once(' ') {
            Some((id, title)) => {
                let id = parse_id(id)?;
                if title.trim().is_empty() {
                    Err(String::from("Название пусто"))
                } else {
                    Ok(Command::Rename(id, String::from(title.trim())))
                }
            }
            None => Err(String::from("Неизвестная команда или неверные аргументы")),
        },
        _ => Err(String::from("Неизвестная команда или неверные аргументы")),
    }
}

struct Task {
    id: u32,
    title: String,
    done: bool,
}
struct Organizer {
    tasks: Vec<Task>,
    next_id: u32,
}
impl Organizer {
    fn new() -> Self {
        Self {
            tasks: Vec::new(),
            next_id: 1,
        }
    }
    fn add(&mut self, title: String) -> Result<u32, String> {
        if title.trim().is_empty() {
            return Err(String::from("Название пусто"));
        }
        match self.next_id.checked_add(1) {
            Some(next) => {
                let id = self.next_id;
                self.tasks.push(Task {
                    id,
                    title,
                    done: false,
                });
                self.next_id = next;
                Ok(id)
            }
            None => Err(String::from("Номера закончились")),
        }
    }
    fn done(&mut self, id: u32) -> Result<(), String> {
        for task in &mut self.tasks {
            if task.id == id {
                task.done = true;
                return Ok(());
            }
        }
        Err(String::from("ID не найден"))
    }
    fn delete(&mut self, id: u32) -> Result<(), String> {
        for index in 0..self.tasks.len() {
            if self.tasks[index].id == id {
                self.tasks.remove(index);
                return Ok(());
            }
        }
        Err(String::from("ID не найден"))
    }
    fn rename(&mut self, id: u32, title: String) -> Result<(), String> {
        if title.trim().is_empty() {
            return Err(String::from("Название пусто"));
        }
        for task in &mut self.tasks {
            if task.id == id {
                task.title = title;
                return Ok(());
            }
        }
        Err(String::from("ID не найден"))
    }
    fn list(&self) {
        if self.tasks.is_empty() {
            println!("Задач нет");
        }
        for task in &self.tasks {
            let mark = if task.done {
                "готово"
            } else {
                "ожидает"
            };
            println!("{}: {} ({mark})", task.id, task.title);
        }
    }
}

fn execute(organizer: &mut Organizer, command: Command) -> bool {
    match command {
        Command::Add(title) => {
            match organizer.add(title) {
                Ok(id) => println!("Добавлено: {id}"),
                Err(message) => println!("Ошибка: {message}"),
            }
            true
        }
        Command::List => {
            organizer.list();
            true
        }
        Command::Done(id) => {
            match organizer.done(id) {
                Ok(()) => println!("Завершено: {id}"),
                Err(message) => println!("Ошибка: {message}"),
            }
            true
        }
        Command::Delete(id) => {
            match organizer.delete(id) {
                Ok(()) => println!("Удалено: {id}"),
                Err(message) => println!("Ошибка: {message}"),
            }
            true
        }
        Command::Rename(id, title) => {
            match organizer.rename(id, title) {
                Ok(()) => println!("Изменено: {id}"),
                Err(message) => println!("Ошибка: {message}"),
            }
            true
        }
        Command::Quit => false,
    }
}

fn main() {
    let mut organizer = Organizer::new();
    let stdin = io::stdin();
    println!("Команды: add НАЗВАНИЕ | list | done ID | delete ID | rename ID НАЗВАНИЕ | quit");
    loop {
        let mut line = String::new();
        match stdin.read_line(&mut line) {
            Ok(0) => break,
            Ok(_) => match parse(&line) {
                Ok(command) => {
                    if !execute(&mut organizer, command) {
                        break;
                    }
                }
                Err(message) => println!("Ошибка: {message}"),
            },
            Err(_) => {
                println!("Ошибка: не удалось прочитать ввод");
                break;
            }
        }
    }
    println!("Сеанс завершён");
}
```
```text
Команды: add НАЗВАНИЕ | list | done ID | delete ID | rename ID НАЗВАНИЕ | quit
Сеанс завершён
```

При выполнении задания эталон выводит:

```text
Команды: add НАЗВАНИЕ | list | done ID | delete ID | rename ID НАЗВАНИЕ | quit
Добавлено: 1
Изменено: 1
1: Читать две главы (ожидает)
Ошибка: ID не найден
Сеанс завершён
```

Этот вывод относится к пустому потоку. При вводе команд из задания между первой и последней строками появятся ответы органайзера. [Официальные объяснения чтения ввода и `split_once`](https://doc.rust-lang.org/std/io/struct.Stdin.html#method.read_line) · [`split_once`](https://doc.rust-lang.org/std/primitive.str.html#method.split_once).

## Контрольная точка после урока 40

Закройте страницу и без копирования объясните путь одной команды: строка → `parse` → `Command` → действие → ответ. Соберите программу, добавьте две задачи, завершите одну, удалите другую и проверьте неверный ID. Затем закройте программу и объясните, почему список исчез. Если теряется разбор команды, повторите [урок 36](/read/rust-36-commands?lang=ru); если изменение списка — [урок 37](/read/rust-37-operations?lang=ru); если ID — [урок 38](/read/rust-38-ids?lang=ru); если проверка ошибок — [урок 39](/read/rust-39-scenarios?lang=ru).

[Предыдущий урок](/read/rust-39-scenarios?lang=ru) · [Оглавление](/course/rust?lang=ru)
