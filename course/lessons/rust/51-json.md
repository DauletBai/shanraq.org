# Формат JSON

_Лид (summary):_ **Урок 51. Превратить задачи в JSON и проверить их после чтения.**

Повторите [структуру Task](/read/rust-27-structs?lang=ru), [внешние библиотеки](/read/rust-48-dependencies?lang=ru) и [чтение файлов](/read/rust-50-files?lang=ru). Здесь строка в памяти; запись на диск будет в следующем уроке.

## Образ и опорная карта

Список задач похож на карточки в подписанном конверте. JSON — текстовый формат, где имена полей и значения записаны по общим правилам. Текст удобно переслать и прочесть другой программой, но сам конверт не подтверждает правдивость записей. **Карта:** значения → JSON → чтение JSON → проверка смысла → готовые задачи. Воспроизведите путь без подсказки.

`serde` превращает структуры в данные и обратно; `serde_json` работает именно с JSON. **Сериализация** — перевод значения в текст, **десериализация** — обратное действие. `#[derive(Serialize, Deserialize)]` просит компилятор создать стандартные реализации для структуры; `#[serde(deny_unknown_fields)]` отклоняет незнакомые поля, чтобы опечатка не пропала молча. Атрибут `#[...]` — инструкция компилятору над следующим объявлением. `version` — номер формы файла: при изменении формы старый файл нельзя молча считать новым. `next_id` должен быть больше всех существующих ID. `HashSet` — множество без повторов; `insert` возвращает `false`, если номер уже был. `?` передаёт ошибку вызывающему коду. `map_err` здесь переводит ошибку библиотеки в текст для `run`.

## Запустите и разберите

Откройте отдельный проект Cargo. В `Cargo.toml` используйте зависимости ниже, замените `src/main.rs` первым блоком Rust и запустите `cargo run`. Cargo при первой сборке загружает библиотеки; `Cargo.lock` закрепляет разрешённые версии. В образце есть одна задача. `to_string` создаёт строку JSON, `from_str` читает её обратно, а `validate` проверяет ограничения, которые сам JSON не знает. Порядок полей в этом конкретном выводе зависит от порядка полей структуры; другие программы могут записать их иначе.

В `Cargo.toml` добавьте зависимости (строки ниже стоят после `[package]`):

В записи зависимости `version = "=1.0.228"` внутренний знак `=` требует именно эту версию, а `features = ["derive"]` включает возможность создавать реализации через `derive`; квадратные скобки здесь обозначают список возможностей. Это учебная фиксация версий, а не обещание вечной актуальности.

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

fn validate(data: &Organizer) -> Result<(), String> {
    if data.version != 1 || data.next_id == 0 {
        return Err(String::from("invalid version or next ID"));
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
    Ok(())
}

fn run() -> Result<(), String> {
    let source = Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from("Купить книгу"),
            done: false,
        }],
    };
    let json = serde_json::to_string(&source).map_err(|error| error.to_string())?;
    let loaded: Organizer = serde_json::from_str(&json).map_err(|error| error.to_string())?;
    validate(&loaded)?;
    println!("{json}");
    println!("Восстановлено задач: {}", loaded.tasks.len());
    Ok(())
}

fn main() {
    if let Err(error) = run() {
        eprintln!("Ошибка: {error}");
        std::process::exit(1);
    }
}
```
```text
{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Купить книгу","done":false}]}
Восстановлено задач: 1
```

## Вопросы для воспроизведения

1. Чем чтение JSON отличается от проверки его смысла?
2. Зачем нужен `next_id`, если уже есть список задач?
3. Почему незнакомое поле следует отклонить?

## Задание

**Обязательное.** Добавьте вторую задачу с ID 2, измените `next_id` на 3 и выведите число восстановленных задач. Затем после чтения замените у второй задачи ID на 1 и покажите, что проверка возвращает ошибку. Сначала предскажите обе строки вывода.

## Ответы

Изменяйте данные до `to_string`; после `from_str` у вас снова `Organizer`. Для ошибочного варианта измените `loaded.tasks[1].id`, затем снова вызовите проверку. Проверяйте `validate(&bad).is_err()`, а не аварийное завершение.

<!-- task-answer -->
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

fn validate(data: &Organizer) -> Result<(), String> {
    if data.version != 1 || data.next_id == 0 {
        return Err(String::from("invalid version or next ID"));
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
    Ok(())
}

fn run() -> Result<(), String> {
    let source = Organizer {
        version: 1,
        next_id: 3,
        tasks: vec![
            Task {
                id: 1,
                title: String::from("Купить книгу"),
                done: false,
            },
            Task {
                id: 2,
                title: String::from("Позвонить другу"),
                done: false,
            },
        ],
    };
    let json = serde_json::to_string(&source).map_err(|error| error.to_string())?;
    let mut loaded: Organizer = serde_json::from_str(&json).map_err(|error| error.to_string())?;
    validate(&loaded)?;
    println!("{json}");
    println!("Восстановлено задач: {}", loaded.tasks.len());
    loaded.tasks[1].id = 1;
    println!("Повтор ID отклонён: {}", validate(&loaded).is_err());
    Ok(())
}

fn main() {
    if let Err(error) = run() {
        eprintln!("Ошибка: {error}");
        std::process::exit(1);
    }
}
```
```text
{"version":1,"next_id":3,"tasks":[{"id":1,"title":"Купить книгу","done":false},{"id":2,"title":"Позвонить другу","done":false}]}
Восстановлено задач: 2
Повтор ID отклонён: true
```

## После проверки

JSON проверяет синтаксис, `validate` — правила органайзера. Не подменяйте повреждённые данные пустым списком. Возврат к [уроку 50](/read/rust-50-files?lang=ru). [Документация Serde](https://serde.rs/derive.html).

[Предыдущий урок](/read/rust-50-files?lang=ru) · [Оглавление](/course/rust?lang=ru)
