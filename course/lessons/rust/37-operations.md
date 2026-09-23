# Список задач в памяти: добавить, изменить, завершить, удалить

_Лид (summary):_ **Урок 37. Выполнить четыре действия с задачами и увидеть, как удаление сдвигает места в списке.**

Предварительно: уроки 21, 24, 27, 31–33 и 36. Если забыли изменяемое заимствование, вернитесь к [уроку 21](/read/rust-21-mutable-borrow?lang=ru).

## Образ и опорная карта

Представьте карточки задач, разложенные в ряд. Можно положить новую карточку, переписать название, поставить отметку и убрать карточку. `Vec<Task>` — такой ряд в памяти программы. Карточки пока исчезают после закрытия программы: файлы появятся позднее. Образ ограничен: место в ряду не является постоянным номером карточки.

**Проверка → изменение `Vec<Task>` → показ результата; закрытие программы → данные в памяти исчезают.** `Task` объединяет `title: String` и `done: bool`. Поле `done` истинно после завершения. `&mut Vec<Task>` даёт функции временный исключительный доступ к списку для изменения, а `&[Task]` — только чтение среза списка. `get_mut(index)` возвращает `Option<&mut Task>`: место может отсутствовать. `tasks.remove(index)` удаляет элемент и сдвигает следующие влево; перед ним мы проверяем границу, поскольку сам `remove` при неверном месте остановил бы программу. `usize` — тип места в списке из урока 15. `Result<(), String>` означает: успешное действие не возвращает данных (`()`), ошибка возвращает объяснение. `for task in tasks` читает каждую карточку по порядку.

## Четыре действия

Сохраните прежний `src/main.rs` и замените целиком. Запустите `cargo run` в проекте `organizer`. Новых зависимостей нет. Предскажите, какая карточка останется после удаления места 0. `[ ]` означает ожидающую задачу, а `[готово]` — завершённую. В этом учебном `main` `.expect(...)` достаёт успешный `Ok`, но остановит программу на `Err`; он уместен только для заранее известных верных данных, а не для ввода пользователя.

```rust
struct Task {
    title: String,
    done: bool,
}

fn add(tasks: &mut Vec<Task>, title: &str) -> Result<(), String> {
    if title.trim().is_empty() {
        return Err(String::from("Название пусто"));
    }
    tasks.push(Task {
        title: String::from(title),
        done: false,
    });
    Ok(())
}

fn rename(tasks: &mut Vec<Task>, index: usize, title: &str) -> Result<(), String> {
    match tasks.get_mut(index) {
        Some(task) => {
            task.title = String::from(title);
            Ok(())
        }
        None => Err(String::from("Такого места нет")),
    }
}

fn complete(tasks: &mut Vec<Task>, index: usize) -> Result<(), String> {
    match tasks.get_mut(index) {
        Some(task) => {
            task.done = true;
            Ok(())
        }
        None => Err(String::from("Такого места нет")),
    }
}

fn remove(tasks: &mut Vec<Task>, index: usize) -> Result<(), String> {
    if index >= tasks.len() {
        return Err(String::from("Такого места нет"));
    }
    tasks.remove(index);
    Ok(())
}

fn show(tasks: &[Task]) {
    for task in tasks {
        if task.done {
            println!("[готово] {}", task.title);
        } else {
            println!("[ ] {}", task.title);
        }
    }
}

fn main() {
    let mut tasks = Vec::new();
    add(&mut tasks, "Купить хлеб").expect("valid title");
    add(&mut tasks, "Прочитать главу").expect("valid title");
    rename(&mut tasks, 1, "Прочитать две главы").expect("existing position");
    complete(&mut tasks, 0).expect("existing position");
    remove(&mut tasks, 0).expect("existing position");
    show(&tasks);
}
```
```text
[ ] Прочитать две главы
```

В учебном `main` вызовы `.expect(...)` допустимы лишь потому, что здесь мы заранее задали существующие места и непустые названия. Это метод `Result`, который достаёт `Ok`, но при `Err` останавливает программу; в настоящем вводе так поступать нельзя. В задании мы обработаем ошибки явно. После удаления второй элемент переходит на место 0. Это причина в следующем уроке ввести устойчивый ID.

## Проверка понимания

1. Почему `get_mut` возвращает `Option`, а не сразу `&mut Task`?
2. Что случится с местом второй карточки после удаления первой?
3. Останутся ли эти задачи после следующего запуска?

Места может не быть; вторая карточка перейдёт на место 0; нет, пока нет сохранения. Восстановите опорную карту.

## Задание

**Обязательное.** Дополните `rename`: отклоняйте название, состоящее только из пробелов, не меняя прежнее. После успешных действий попробуйте переименовать место 0 в `"   "` и завершить место 9. Выведите обе ошибки и затем список; прежнее название должно сохраниться.

## Подсказка

Проверку `title.trim().is_empty()` поставьте до `get_mut`. Ошибку из `Result` выведите через `match`, как в уроке 32.

## Эталон после попытки

<!-- task-answer -->
```rust
struct Task {
    title: String,
    done: bool,
}

fn add(tasks: &mut Vec<Task>, title: &str) -> Result<(), String> {
    if title.trim().is_empty() {
        return Err(String::from("Название пусто"));
    }
    tasks.push(Task {
        title: String::from(title),
        done: false,
    });
    Ok(())
}

fn rename(tasks: &mut Vec<Task>, index: usize, title: &str) -> Result<(), String> {
    if title.trim().is_empty() {
        return Err(String::from("Название пусто"));
    }
    match tasks.get_mut(index) {
        Some(task) => {
            task.title = String::from(title);
            Ok(())
        }
        None => Err(String::from("Такого места нет")),
    }
}

fn complete(tasks: &mut Vec<Task>, index: usize) -> Result<(), String> {
    match tasks.get_mut(index) {
        Some(task) => {
            task.done = true;
            Ok(())
        }
        None => Err(String::from("Такого места нет")),
    }
}

fn remove(tasks: &mut Vec<Task>, index: usize) -> Result<(), String> {
    if index >= tasks.len() {
        return Err(String::from("Такого места нет"));
    }
    tasks.remove(index);
    Ok(())
}

fn show(tasks: &[Task]) {
    for task in tasks {
        if task.done {
            println!("[готово] {}", task.title);
        } else {
            println!("[ ] {}", task.title);
        }
    }
}

fn main() {
    let mut tasks = Vec::new();
    add(&mut tasks, "Купить хлеб").expect("valid title");
    add(&mut tasks, "Прочитать главу").expect("valid title");
    rename(&mut tasks, 1, "Прочитать две главы").expect("existing position");
    complete(&mut tasks, 0).expect("existing position");
    remove(&mut tasks, 0).expect("existing position");
    for result in [rename(&mut tasks, 0, "   "), complete(&mut tasks, 9)] {
        if let Err(message) = result {
            println!("Ошибка: {message}");
        }
    }
    show(&tasks);
}
```
```text
Ошибка: Название пусто
Ошибка: Такого места нет
[ ] Прочитать две главы
```

[Документация `Vec::get_mut` и `Vec::remove`](https://doc.rust-lang.org/std/vec/struct.Vec.html#method.remove).

[Предыдущий урок](/read/rust-36-commands?lang=ru) · [Оглавление](/course/rust?lang=ru) · [Следующий урок](/read/rust-38-ids?lang=ru)
