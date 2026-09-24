# Резервная копия

_Лид (summary):_ **Урок 53. Проверить резервную копию и восстановить данные в отдельный файл.**

Повторите [JSON](/read/rust-51-json?lang=ru) и [сохранение](/read/rust-52-safe-save?lang=ru). Здесь используется учебная временная папка.

## Образ и опорная карта

Запасной ключ полезен, только если он подходит к замку. Так же и копия файла полезна, только если её можно прочитать и проверить. **Карта:** основной файл → проверенная копия → повреждение основного → проверка копии → отдельный восстановленный файл. Исходный повреждённый файл остаётся для разбора. Образ не обещает, что копия защитит от удаления всей папки.

`load` превращает JSON в `Organizer` и проверяет версию, номера и названия. `io::ErrorKind::InvalidData` обозначает данные неверной формы или смысла. `copy_checked` сперва проверяет источник, затем создаёт новый целевой файл через `create_new(true)`, пишет и синхронизирует, снова читает результат. `create_new` не перезапишет чужой файл. `read_to_string` потребует корректный UTF-8. В учебной папке `primary`, `backup` и `recovered` — три разных пути. `fs::write` портит только учебный основной файл, чтобы показать восстановление.

## Запустите и разберите

Создайте проект с теми же зависимостями Serde, что в уроке 51. Замените `src/main.rs` первым Rust блоком и запустите `cargo run`. Проследите порядок вызовов: копию создают до повреждения; после повреждения читают копию, а восстановленный текст пишут в новый файл. Если проверка источника не прошла, целевой файл ещё не открыт. Ошибка записи может оставить неполный новый целевой файл; его нельзя принимать за годную копию без повторной проверки.

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
use std::fs::{self, OpenOptions};
use std::io::{self, Write};
use std::path::Path;

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

fn load(text: &str) -> io::Result<Organizer> {
    let data: Organizer = serde_json::from_str(text)
        .map_err(|error| io::Error::new(io::ErrorKind::InvalidData, error))?;
    if data.version != 1 || data.next_id == 0 {
        return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid header"));
    }
    let mut ids = HashSet::new();
    for task in &data.tasks {
        if task.id == 0
            || task.id >= data.next_id
            || task.title.trim().is_empty()
            || !ids.insert(task.id)
        {
            return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid task"));
        }
    }
    Ok(data)
}

fn copy_checked(source: &Path, target: &Path) -> io::Result<()> {
    let text = fs::read_to_string(source)?;
    load(&text)?;
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(target)?;
    file.write_all(text.as_bytes())?;
    file.sync_all()?;
    drop(file);
    load(&fs::read_to_string(target)?)?;
    Ok(())
}

fn run() -> io::Result<()> {
    let folder = std::env::temp_dir().join(format!("rust-lesson-53-{}", std::process::id()));
    fs::create_dir(&folder)?;
    let primary = folder.join("tasks.json");
    let backup = folder.join("backup.json");
    let recovered = folder.join("recovered.json");
    let data = Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from("Купить книгу"),
            done: false,
        }],
    };
    let json = serde_json::to_string(&data)
        .map_err(|error| io::Error::new(io::ErrorKind::InvalidData, error))?;
    fs::write(&primary, json)?;
    copy_checked(&primary, &backup)?;
    fs::write(&primary, "{broken")?;
    println!(
        "Основной файл повреждён: {}",
        load(&fs::read_to_string(&primary)?).is_err()
    );
    copy_checked(&backup, &recovered)?;
    let restored = load(&fs::read_to_string(&recovered)?)?;
    println!("Восстановлено: {}", restored.tasks[0].title);
    println!(
        "Источник не изменён: {}",
        fs::read_to_string(&primary)? == "{broken"
    );
    fs::remove_dir_all(folder)?;
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
Основной файл повреждён: true
Восстановлено: Купить книгу
Источник не изменён: true
```

## Вопросы для воспроизведения

1. Почему сначала проверяют источник?
2. Почему восстановление идёт в третий файл?
3. Что нужно сделать, если повторное чтение копии не прошло?

## Задание

**Обязательное.** Создайте отдельный файл с неверным JSON `{broken`, попробуйте `copy_checked` в новый путь и выведите, что ошибка обнаружена и восстановленный файл не появился. Не стирайте исправную копию.

## Ответы

Вызовите `copy_checked(&bad_backup, &bad_target)`. Проверьте результат через `is_err()`, а отсутствие цели — через `!bad_target.exists()`.

<!-- task-answer -->
```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::fs::{self, OpenOptions};
use std::io::{self, Write};
use std::path::Path;

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

fn load(text: &str) -> io::Result<Organizer> {
    let data: Organizer = serde_json::from_str(text)
        .map_err(|error| io::Error::new(io::ErrorKind::InvalidData, error))?;
    if data.version != 1 || data.next_id == 0 {
        return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid header"));
    }
    let mut ids = HashSet::new();
    for task in &data.tasks {
        if task.id == 0
            || task.id >= data.next_id
            || task.title.trim().is_empty()
            || !ids.insert(task.id)
        {
            return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid task"));
        }
    }
    Ok(data)
}

fn copy_checked(source: &Path, target: &Path) -> io::Result<()> {
    let text = fs::read_to_string(source)?;
    load(&text)?;
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(target)?;
    file.write_all(text.as_bytes())?;
    file.sync_all()?;
    drop(file);
    load(&fs::read_to_string(target)?)?;
    Ok(())
}

fn run() -> io::Result<()> {
    let folder = std::env::temp_dir().join(format!("rust-lesson-53-{}", std::process::id()));
    fs::create_dir(&folder)?;
    let primary = folder.join("tasks.json");
    let backup = folder.join("backup.json");
    let recovered = folder.join("recovered.json");
    let data = Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from("Купить книгу"),
            done: false,
        }],
    };
    let json = serde_json::to_string(&data)
        .map_err(|error| io::Error::new(io::ErrorKind::InvalidData, error))?;
    fs::write(&primary, json)?;
    copy_checked(&primary, &backup)?;
    fs::write(&primary, "{broken")?;
    println!(
        "Основной файл повреждён: {}",
        load(&fs::read_to_string(&primary)?).is_err()
    );
    copy_checked(&backup, &recovered)?;
    let restored = load(&fs::read_to_string(&recovered)?)?;
    println!("Восстановлено: {}", restored.tasks[0].title);
    println!(
        "Источник не изменён: {}",
        fs::read_to_string(&primary)? == "{broken"
    );
    let damaged_backup = folder.join("damaged-backup.json");
    let other = folder.join("other-recovery.json");
    fs::write(&damaged_backup, "{broken")?;
    let rejected = copy_checked(&damaged_backup, &other).is_err() && !other.exists();
    println!("Плохая копия отклонена: {}", rejected);
    fs::remove_dir_all(folder)?;
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
Основной файл повреждён: true
Восстановлено: Купить книгу
Источник не изменён: true
Плохая копия отклонена: true
```

## После проверки

Проверенная копия не заменяет стратегию хранения: она может устареть или исчезнуть вместе с основным файлом. Перед восстановлением покажите человеку, какую версию данных вы нашли. Вернитесь к [уроку 51](/read/rust-51-json?lang=ru). [Документация create_new](https://doc.rust-lang.org/std/fs/struct.OpenOptions.html#method.create_new).

[Предыдущий урок](/read/rust-52-safe-save?lang=ru) · [Оглавление](/course/rust?lang=ru)
