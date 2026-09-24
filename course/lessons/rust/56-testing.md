# Полный набор проверок для файла задач

_Лид (summary):_ **Урок 56. Проверить сохранённые задачи, повреждённый файл, повторное чтение и текст на разных языках.**

Нужны [JSON и проверка данных](/read/rust-51-json?lang=ru), [резервная копия](/read/rust-53-backup?lang=ru) и [первый автоматический тест](/read/rust-16-first-test?lang=ru). Для упражнения используйте только отдельные файлы в системной временной папке.

## Образ и опорная карта

Перед поездкой вы проверяете не только машину на месте, но и повторный запуск после остановки, запасной ключ и дорогу в дождь. Так же тесты проверяют разные сценарии файла. **Карта:** отдельный учебный путь → подготовить данные → выполнить действие → сравнить с ожидаемым → удалить только свой файл. Пройдите эту цепочку по памяти. Образ имеет предел: несколько тестов не доказывают отсутствие всех возможных ошибок.

**Набор тестов** — несколько независимых проверок поведения. **Изоляция** означает, что проверка не портит настоящие данные и не зависит от соседней проверки. ASCII — небольшая таблица основных латинских букв, цифр и знаков; буквы `é` в ней нет. Здесь тест проверяет сохранение буквы за пределами этой таблицы. `practice_path` добавляет к имени файла номер процесса и метку сценария; метки `unicode`, `damaged`, `duplicate`, `again` не совпадают. Это снижает риск пересечения при параллельном запуске тестов. `#[test]` сообщает Cargo, что функция — тест; `cargo test` выполняет их, а `cargo run` запускает `main`. `assert_eq!` сравнивает значения, `assert!` требует истинного условия. Внутри теста `unwrap()` делает неожиданную ошибку провалом теста; в обычном приложении ошибку следует обработать и объяснить пользователю.

## Запустите и разберите

Создайте отдельный проект Cargo. В `Cargo.toml` после `[package]` добавьте зависимости ниже: внутреннее `=` в номере версии требует точную версию, а список `features = ["derive"]` включает создание реализаций Serde. Поместите первый Rust блок в `src/main.rs`. `cargo run` должен вывести число прочитанных задач. Затем `cargo test`: два теста должны пройти. Первый записывает и читает текст с буквой вне ASCII; второй оставляет повреждённый JSON без изменения. `read_checked` сначала читает UTF-8 и JSON, затем проверяет версию, ID и названия. Сообщение `Err` не превращается в пустой органайзер.

```toml
[dependencies]
serde = { version = "=1.0.228", features = ["derive"] }
serde_json = "=1.0.149"
```

```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::fs;
use std::path::{Path, PathBuf};

#[derive(Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
struct Task {
    id: u64,
    title: String,
    done: bool,
}

#[derive(Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
struct Organizer {
    version: u32,
    next_id: u64,
    tasks: Vec<Task>,
}

fn read_checked(path: &Path) -> Result<Organizer, String> {
    let text = fs::read_to_string(path).map_err(|error| error.to_string())?;
    let data: Organizer = serde_json::from_str(&text).map_err(|error| error.to_string())?;
    if data.version != 1 || data.next_id == 0 {
        return Err(String::from("invalid header"));
    }
    let mut ids = HashSet::new();
    for task in &data.tasks {
        if task.id == 0
            || task.id >= data.next_id
            || task.title.trim().is_empty()
            || !ids.insert(task.id)
        {
            return Err(String::from("invalid task"));
        }
    }
    Ok(data)
}

fn practice_path(label: &str) -> PathBuf {
    std::env::temp_dir().join(format!(
        "rust-course-56-{}-{label}.json",
        std::process::id()
    ))
}

fn sample(title: &str) -> Organizer {
    Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from(title),
            done: false,
        }],
    }
}

fn main() -> Result<(), String> {
    let path = practice_path("main");
    let text = serde_json::to_string(&sample("Купить книгу")).map_err(|e| e.to_string())?;
    fs::write(&path, text).map_err(|e| e.to_string())?;
    let loaded = read_checked(&path)?;
    println!("Прочитано задач: {}", loaded.tasks.len());
    fs::remove_file(path).map_err(|e| e.to_string())?;
    Ok(())
}

#[test]
fn unicode_round_trip() {
    let path = practice_path("unicode");
    let text = serde_json::to_string(&sample("Café")).unwrap();
    fs::write(&path, text).unwrap();
    let loaded = read_checked(&path).unwrap();
    assert_eq!(loaded.tasks[0].title, "Café");
    fs::remove_file(path).unwrap();
}

#[test]
fn damaged_json_is_rejected() {
    let path = practice_path("damaged");
    fs::write(&path, "{broken").unwrap();
    assert!(read_checked(&path).is_err());
    assert_eq!(fs::read_to_string(&path).unwrap(), "{broken");
    fs::remove_file(path).unwrap();
}
```
```text
Прочитано задач: 1
```

## Вопросы для воспроизведения

1. Чем изолированный тест отличается от работы с вашим настоящим файлом?
2. Почему проверка повреждённого файла сравнивает его содержимое после ошибки?
3. Что случится, если тестовая функция с `unwrap()` неожиданно получит `Err`?

## Задание

**Обязательное.** Добавьте тест, который записывает две задачи с одинаковым ID и ожидает отказа `read_checked`. Добавьте второй тест, который читает один и тот же файл дважды и убеждается, что название осталось прежним. У каждого теста должна быть своя метка пути; после проверки удалите только созданный им файл. Запустите `cargo test` и убедитесь, что прошли четыре теста.

## Ответы

Для дубликата создайте `sample`, поставьте `next_id = 3`, добавьте ещё одну задачу с `id: 1` и запишите JSON. Для повторного чтения сравните `first.tasks[0].title` и `second.tasks[0].title`.

<!-- task-answer -->
```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::fs;
use std::path::{Path, PathBuf};

#[derive(Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
struct Task {
    id: u64,
    title: String,
    done: bool,
}

#[derive(Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
struct Organizer {
    version: u32,
    next_id: u64,
    tasks: Vec<Task>,
}

fn read_checked(path: &Path) -> Result<Organizer, String> {
    let text = fs::read_to_string(path).map_err(|error| error.to_string())?;
    let data: Organizer = serde_json::from_str(&text).map_err(|error| error.to_string())?;
    if data.version != 1 || data.next_id == 0 {
        return Err(String::from("invalid header"));
    }
    let mut ids = HashSet::new();
    for task in &data.tasks {
        if task.id == 0
            || task.id >= data.next_id
            || task.title.trim().is_empty()
            || !ids.insert(task.id)
        {
            return Err(String::from("invalid task"));
        }
    }
    Ok(data)
}

fn practice_path(label: &str) -> PathBuf {
    std::env::temp_dir().join(format!(
        "rust-course-56-{}-{label}.json",
        std::process::id()
    ))
}

fn sample(title: &str) -> Organizer {
    Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from(title),
            done: false,
        }],
    }
}

fn main() -> Result<(), String> {
    let path = practice_path("main");
    let text = serde_json::to_string(&sample("Купить книгу")).map_err(|e| e.to_string())?;
    fs::write(&path, text).map_err(|e| e.to_string())?;
    let loaded = read_checked(&path)?;
    println!("Прочитано задач: {}", loaded.tasks.len());
    fs::remove_file(path).map_err(|e| e.to_string())?;
    Ok(())
}

#[test]
fn unicode_round_trip() {
    let path = practice_path("unicode");
    let text = serde_json::to_string(&sample("Café")).unwrap();
    fs::write(&path, text).unwrap();
    let loaded = read_checked(&path).unwrap();
    assert_eq!(loaded.tasks[0].title, "Café");
    fs::remove_file(path).unwrap();
}

#[test]
fn damaged_json_is_rejected() {
    let path = practice_path("damaged");
    fs::write(&path, "{broken").unwrap();
    assert!(read_checked(&path).is_err());
    assert_eq!(fs::read_to_string(&path).unwrap(), "{broken");
    fs::remove_file(path).unwrap();
}

#[test]
fn duplicate_ids_are_rejected() {
    let path = practice_path("duplicate");
    let mut data = sample("Купить книгу");
    data.next_id = 3;
    data.tasks.push(Task {
        id: 1,
        title: String::from("second"),
        done: false,
    });
    fs::write(&path, serde_json::to_string(&data).unwrap()).unwrap();
    assert!(read_checked(&path).is_err());
    fs::remove_file(path).unwrap();
}

#[test]
fn second_read_keeps_data() {
    let path = practice_path("again");
    fs::write(
        &path,
        serde_json::to_string(&sample("Купить книгу")).unwrap(),
    )
    .unwrap();
    let first = read_checked(&path).unwrap();
    let second = read_checked(&path).unwrap();
    assert_eq!(first.tasks[0].title, second.tasks[0].title);
    fs::remove_file(path).unwrap();
}
```
```text
Прочитано задач: 1
```

## После проверки

Если тест оставил файл после аварии, удалите только файл с учебным именем после просмотра. Параллельные тесты могут выполняться в другом порядке: не заставляйте их пользоваться одним путём. Проверка отказа записи и защита от отключения питания требуют отдельных сценариев и среды; этот урок подтверждает перечисленные свойства, а не все возможные. [Документация cargo test](https://doc.rust-lang.org/cargo/commands/cargo-test.html).

[Предыдущий урок](/read/rust-55-documentation?lang=ru) · [Оглавление](/course/rust?lang=ru)
