# Самостоятельное расширение: приоритет задачи

_Лид (summary):_ **Урок 59. Добавить приоритет, сохранив чтение старого JSON и отказ от неверных значений.**

Повторите [структуры](/read/rust-27-structs?lang=ru), [JSON и проверку](/read/rust-51-json?lang=ru) и [тесты файлов](/read/rust-56-testing?lang=ru). Сначала сформулируйте поведение, потом пишите код: это самостоятельная работа с эталоном после задания.

## Образ и опорная карта

На карточке задачи появляется цветная отметка срочности. Старая карточка без отметки не становится срочной сама собой: ей дают обычную отметку. **Карта:** записать правило → добавить поле → дать значение старой записи → отвергнуть неверное значение → проверить выборку → объяснить пользователю. Предел образа: цвет на карточке не решает, что задача действительно важна; это решение человека.

**Приоритет** — число от 1 до 3: 1 низкий, 2 обычный, 3 срочный. Старый JSON не содержит поля `priority`. Атрибут `#[serde(default = "ordinary_priority")]` просит Serde вызвать `ordinary_priority`, когда поле отсутствует. Он не исправляет явно записанное неверное число. `1..=3` — включительный диапазон: обе границы входят. `high_open` считает срочные и ещё не завершённые задачи. Номер схемы `version` можно оставить 1 для этого совместимого добавления, если старые файлы читаются, а новое поле всегда имеет документированный смысл; несовместимое изменение потребует отдельного решения о версии и переносе.

## Запустите и разберите

В отдельный проект Cargo добавьте зависимости ниже, поместите первый Rust блок в `src/main.rs`, запустите `cargo run` и `cargo test`. Базовая программа читает старый JSON без приоритета и сообщает число задач. Она проверяет версию, положительные неповторяющиеся ID и непустые названия. Не меняйте этот старый JSON вручную: он нужен для проверки совместимости.

```toml
[dependencies]
serde = { version = "=1.0.228", features = ["derive"] }
serde_json = "=1.0.149"
```

```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;

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

fn decode(text: &str) -> Result<Organizer, String> {
    let data: Organizer = serde_json::from_str(text).map_err(|e| e.to_string())?;
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

fn main() -> Result<(), String> {
    let old_json =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Купить книгу","done":false}]}"#;
    let data = decode(old_json)?;
    println!("Задач: {}", data.tasks.len());
    Ok(())
}

#[test]
fn reads_old_task() {
    let text =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Купить книгу","done":false}]}"#;
    assert_eq!(decode(text).unwrap().tasks.len(), 1);
}
```
```text
Задач: 1
```

## Вопросы для воспроизведения

1. Почему старому файлу нужен определённый приоритет по умолчанию?
2. Должен ли `default` исправлять явно записанный `priority: 0`?
3. Почему подсчёт срочных исключает завершённые задачи?

## Задание

**Обязательное.** До кода запишите четыре правила: допустимы 1, 2, 3; отсутствие поля в старом JSON означает 2; явное 0 или 4 вызывает ошибку; срочными для выборки считаются только незавершённые задачи с 3. Затем добавьте поле и функцию подсчёта, два теста для старого и неверного JSON. Выведите приоритет старой задачи и число срочных после добавления новой незавершённой задачи с 3. Сначала предскажите две строки.

## Ответы

Создайте `fn ordinary_priority() -> u8 { 2 }` до структуры. Над новым полем поставьте `#[serde(default = "ordinary_priority")]`. В проверке используйте `!(1..=3).contains(&task.priority)`; в выборке — `iter().filter(...).count()`.

<!-- task-answer -->
```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;

fn ordinary_priority() -> u8 {
    2
}

#[derive(Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
struct Task {
    id: u64,
    title: String,
    done: bool,
    #[serde(default = "ordinary_priority")]
    priority: u8,
}

#[derive(Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
struct Organizer {
    version: u32,
    next_id: u64,
    tasks: Vec<Task>,
}

fn decode(text: &str) -> Result<Organizer, String> {
    let data: Organizer = serde_json::from_str(text).map_err(|e| e.to_string())?;
    if data.version != 1 || data.next_id == 0 {
        return Err(String::from("invalid header"));
    }
    let mut ids = HashSet::new();
    for task in &data.tasks {
        if task.id == 0
            || task.id >= data.next_id
            || task.title.trim().is_empty()
            || !(1..=3).contains(&task.priority)
            || !ids.insert(task.id)
        {
            return Err(String::from("invalid task"));
        }
    }
    Ok(data)
}

fn high_open(data: &Organizer) -> usize {
    data.tasks
        .iter()
        .filter(|task| task.priority == 3 && !task.done)
        .count()
}

fn main() -> Result<(), String> {
    let old_json =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Купить книгу","done":false}]}"#;
    let mut data = decode(old_json)?;
    println!("Приоритет старой задачи: {}", data.tasks[0].priority);
    data.tasks.push(Task {
        id: 2,
        title: String::from("Позвонить другу"),
        done: false,
        priority: 3,
    });
    data.next_id = 3;
    println!("Срочных незавершённых: {}", high_open(&data));
    Ok(())
}

#[test]
fn old_file_gets_ordinary_priority() {
    let text =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Купить книгу","done":false}]}"#;
    assert_eq!(decode(text).unwrap().tasks[0].priority, 2);
}

#[test]
fn invalid_priority_is_rejected() {
    let text = r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Купить книгу","done":false,"priority":0}]}"#;
    assert!(decode(text).is_err());
}
```
```text
Приоритет старой задачи: 2
Срочных незавершённых: 1
```

## После проверки

Старый файл должен оставаться читаемым, но повреждённый или явно неверный файл нельзя молча заменить пустым списком. По той же схеме можно предложить собственное расширение, однако сначала сформулируйте правило для старых данных и тест отказа. [Документация serde default](https://serde.rs/field-attrs.html#default).

[Предыдущий урок](/read/rust-58-distribution?lang=ru) · [Оглавление](/course/rust?lang=ru)
