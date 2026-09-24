# Сохранение без потери данных

_Лид (summary):_ **Урок 52. Записать новые данные через соседний временный файл, сохранив старые при отказе.**

Нужны [пути](/read/rust-49-paths?lang=ru) и [файловые ошибки](/read/rust-50-files?lang=ru). Пока сохраняем учебный текст; JSON из урока 51 можно передать этой функции.

## Образ и опорная карта

Вы не стираете старую страницу блокнота, пока новая ещё пишется. Сначала пишете на отдельном листе, проверяете запись и только затем заменяете страницу. **Карта:** старый файл → новый соседний файл → запись → синхронизация → закрытие → замена. Если шаг записи не удался, старый файл остаётся. Образ имеет предел: внезапное отключение питания и особенности файловой системы требуют более сильных мер.

`with_extension` создаёт путь временного файла рядом с основным; соседство важно для замены в пределах одной файловой системы. `create_new(true)` запрещает перезаписать случайно найденный временный файл. `write_all` передаёт все байты либо ошибку; `sync_all` просит ОС отправить содержимое и метаданные файла на устройство; `drop` закрывает дескриптор. `rename` заменяет целевой путь там, где ОС позволяет, либо возвращает ошибку. `if let Err(error)` выбирает только ветку ошибки; `let _ = remove_file` означает, что неудачу очистки здесь сознательно не превращают в успех сохранения.

## Запустите и разберите

Создайте чистый проект Cargo без зависимостей, скопируйте первый блок в `src/main.rs`, выполните `cargo run`. Учебная папка названа по номеру процесса; если она осталась после аварии, удалите только её после проверки содержимого. В демонстрации старый текст заменяется новым. Ошибка при создании временного файла, записи или синхронизации не вызывает `rename`, поэтому прежний целевой файл не затрагивается. Ошибка `rename` возвращается вызывающему коду.

```rust
use std::fs::{self, OpenOptions};
use std::io::{self, Write};
use std::path::Path;

fn safe_save(path: &Path, text: &str) -> io::Result<()> {
    let temporary = path.with_extension(format!("tmp-{}", std::process::id()));
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(&temporary)?;
    if let Err(error) = file.write_all(text.as_bytes()) {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    if let Err(error) = file.sync_all() {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    drop(file);
    if let Err(error) = fs::rename(&temporary, path) {
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    Ok(())
}

fn run() -> io::Result<()> {
    let folder = std::env::temp_dir().join(format!("rust-lesson-52-{}", std::process::id()));
    fs::create_dir(&folder)?;
    let path = folder.join("tasks.json");
    fs::write(&path, "Старый план")?;
    safe_save(&path, "Новый план")?;
    println!("Сохранено: {}", fs::read_to_string(&path)?);
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
Сохранено: Новый план
```

## Вопросы для воспроизведения

1. Зачем временный файл расположен рядом с основным?
2. Почему `create_new` лучше перезаписи?
3. На каком шаге старый файл впервые может быть заменён?

## Задание

**Обязательное.** До вызова `safe_save` создайте файл с тем же временным именем. Покажите, что сохранение вернуло ошибку и в основном файле остался старый текст. Удалите препятствующий временный файл, повторите сохранение и покажите новый текст.

## Ответы

Вычислите временный путь так же, как функция: `path.with_extension(format!("tmp-{}", std::process::id()))`. Проверьте `is_err()` и сравните прочитанный старый текст.

<!-- task-answer -->
```rust
use std::fs::{self, OpenOptions};
use std::io::{self, Write};
use std::path::Path;

fn safe_save(path: &Path, text: &str) -> io::Result<()> {
    let temporary = path.with_extension(format!("tmp-{}", std::process::id()));
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(&temporary)?;
    if let Err(error) = file.write_all(text.as_bytes()) {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    if let Err(error) = file.sync_all() {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    drop(file);
    if let Err(error) = fs::rename(&temporary, path) {
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    Ok(())
}

fn run() -> io::Result<()> {
    let folder = std::env::temp_dir().join(format!("rust-lesson-52-{}", std::process::id()));
    fs::create_dir(&folder)?;
    let path = folder.join("tasks.json");
    fs::write(&path, "Старый план")?;
    let temporary = path.with_extension(format!("tmp-{}", std::process::id()));
    fs::write(&temporary, "occupied")?;
    let refused = safe_save(&path, "Новый план").is_err();
    let kept = fs::read_to_string(&path)? == "Старый план";
    println!("Старое сохранено: {}", refused && kept);
    fs::remove_file(&temporary)?;
    safe_save(&path, "Новый план")?;
    println!("Сохранено: {}", fs::read_to_string(&path)?);
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
Старое сохранено: true
Сохранено: Новый план
```

## После проверки

Это защита от типичной ошибки записи, но не обещание сохранности при любом сбое. Одного `sync_all` для временного файла недостаточно, чтобы обещать пережить отключение питания после замены; для строгой гарантии нужны правила ОС, файловой системы и синхронизация каталога. Одновременную запись двух процессов этот пример тоже не решает. [Документация rename](https://doc.rust-lang.org/std/fs/fn.rename.html).

[Предыдущий урок](/read/rust-51-json?lang=ru) · [Оглавление](/course/rust?lang=ru)
