# ID задачи: номер, который не сдвигается

_Лид (summary):_ **Урок 38. Выдать каждой задаче устойчивый ID и не переиспользовать его после удаления.**

Предварительно: уроки 10, 13, 27–28, 31–33 и 37. Если место в `Vec` и ID смешиваются, повторите [урок 37](/read/rust-37-operations?lang=ru).

## Образ и опорная карта

В библиотеке книга может сменить полку, но её инвентарный номер остаётся прежним. Так и задача может перейти на другое место в `Vec`, сохранив свой **ID** — уникальный номер записи. Образ ограничен: номер выдаёт наша программа; он устойчив только пока мы правильно ведём счётчик. После нового запуска данные пока исчезают, и нумерация начинается заново.

**Следующий ID → проверка переполнения → добавление → счётчик растёт; удаление → старый ID не выдаётся снова.** `Organizer` хранит `tasks` и `next_id`. Изначально следующий номер равен 1. `checked_add(1)` возвращает `Some(новый номер)` или `None`, если число типа `u32` уже нельзя увеличить. Мы резервируем `u32::MAX`: при достижении этого значения добавление возвращает `Err` **до изменения** списка. Это упрощает правило, но один допустимый номер остаётся неиспользованным. Счётчик никогда не уменьшается при удалении. `for index in 0..self.tasks.len()` перебирает только существующие места; внутри цикла чтение `self.tasks[index]` поэтому не выходит за границу. Поиск идёт по `task.id`, а не по месту. `Self` внутри `impl Organizer` означает сам тип `Organizer`; `self` — конкретный органайзер, как в уроке 28. `u32::MAX` — наибольшее значение `u32`. Учебный `main` применяет `.expect(...)` лишь к заранее заданным верным данным; при настоящем вводе ошибки надо обрабатывать. Круглые скобки в выводе обозначают состояние задачи.

## Счётчик и поиск

Сохраните прежний `src/main.rs`, замените его целиком и выполните `cargo run` в проекте `organizer`. Новых зависимостей нет. До запуска предскажите два напечатанных ID.

```rust
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
    fn add(&mut self, title: &str) -> Result<u32, String> {
        if title.trim().is_empty() {
            return Err(String::from("Пустое название"));
        }
        match self.next_id.checked_add(1) {
            Some(next) => {
                let id = self.next_id;
                self.tasks.push(Task {
                    id,
                    title: String::from(title),
                    done: false,
                });
                self.next_id = next;
                Ok(id)
            }
            None => Err(String::from("Номера закончились")),
        }
    }
    fn remove(&mut self, id: u32) -> Result<(), String> {
        for index in 0..self.tasks.len() {
            if self.tasks[index].id == id {
                self.tasks.remove(index);
                return Ok(());
            }
        }
        Err(String::from("ID не найден"))
    }
    fn complete(&mut self, id: u32) -> Result<(), String> {
        for task in &mut self.tasks {
            if task.id == id {
                task.done = true;
                return Ok(());
            }
        }
        Err(String::from("ID не найден"))
    }
    fn show(&self) {
        for task in &self.tasks {
            let mark = if task.done {
                "готово"
            } else {
                "ожидает"
            };
            println!("ID {}: {} ({mark})", task.id, task.title);
        }
    }
}

fn main() {
    let mut organizer = Organizer::new();
    let first = organizer.add("Купить хлеб").expect("fixed valid title");
    let second = organizer.add("Прочитать книгу").expect("fixed valid title");
    organizer.remove(first).expect("fixed existing ID");
    let third = organizer.add("Прогуляться").expect("fixed valid title");
    organizer.complete(second).expect("fixed existing ID");
    println!("ID: {second}, {third}");
    organizer.show();
}
```
```text
ID: 2, 3
ID 2: Прочитать книгу (готово)
ID 3: Прогуляться (ожидает)
```

Две задачи с одинаковым названием тоже получат разные ID: название не служит ключом. Удаление ID 1 не меняет ID 2. Новый ID 3 показывает, что счётчик выдачи не равен длине списка. После нового запуска без файла программа снова начинает с 1; сохранение и восстановление счётчика появятся на файловом этапе.

## Проверка понимания

1. Почему после удаления ID 1 следующая задача получает ID 3?
2. Почему поиск по месту 1 после удаления мог бы изменить не ту задачу?
3. Может ли `add` частично изменить список при переполнении?

Счётчик не уменьшается; места сдвигаются; нет, проверка стоит перед `push`. Восстановите опорную карту.

## Задание

**Обязательное.** Добавьте проверку границы: попытайтесь завершить несуществующий ID 99, затем установите `next_id` в `u32::MAX` и попытайтесь добавить задачу. Покажите обе ошибки и подтвердите, что число задач и прежние ID не изменились.

## Подсказка

Обработайте `Err(message)` через `match`. Для учебной проверки внутри того же `main` можно присвоить `organizer.next_id = u32::MAX`; обычному пользователю такой доступ давать не нужно.

## Эталон после попытки

<!-- task-answer -->
```rust
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
    fn add(&mut self, title: &str) -> Result<u32, String> {
        if title.trim().is_empty() {
            return Err(String::from("Пустое название"));
        }
        match self.next_id.checked_add(1) {
            Some(next) => {
                let id = self.next_id;
                self.tasks.push(Task {
                    id,
                    title: String::from(title),
                    done: false,
                });
                self.next_id = next;
                Ok(id)
            }
            None => Err(String::from("Номера закончились")),
        }
    }
    fn remove(&mut self, id: u32) -> Result<(), String> {
        for index in 0..self.tasks.len() {
            if self.tasks[index].id == id {
                self.tasks.remove(index);
                return Ok(());
            }
        }
        Err(String::from("ID не найден"))
    }
    fn complete(&mut self, id: u32) -> Result<(), String> {
        for task in &mut self.tasks {
            if task.id == id {
                task.done = true;
                return Ok(());
            }
        }
        Err(String::from("ID не найден"))
    }
    fn show(&self) {
        for task in &self.tasks {
            let mark = if task.done {
                "готово"
            } else {
                "ожидает"
            };
            println!("ID {}: {} ({mark})", task.id, task.title);
        }
    }
}

fn main() {
    let mut organizer = Organizer::new();
    let first = organizer.add("Купить хлеб").expect("fixed valid title");
    let second = organizer.add("Прочитать книгу").expect("fixed valid title");
    organizer.remove(first).expect("fixed existing ID");
    let third = organizer.add("Прогуляться").expect("fixed valid title");
    organizer.complete(second).expect("fixed existing ID");
    println!("ID: {second}, {third}");
    if let Err(message) = organizer.complete(99) {
        println!("Ошибка: {message}");
    }
    organizer.next_id = u32::MAX;
    if let Err(message) = organizer.add("Прогуляться") {
        println!("Ошибка: {message}");
    }
    println!("Всего задач: {}", organizer.tasks.len());
    organizer.show();
}
```
```text
ID: 2, 3
Ошибка: ID не найден
Ошибка: Номера закончились
Всего задач: 2
ID 2: Прочитать книгу (готово)
ID 3: Прогуляться (ожидает)
```

[Документация `u32::checked_add`](https://doc.rust-lang.org/std/primitive.u32.html#method.checked_add).

[Предыдущий урок](/read/rust-37-operations?lang=ru) · [Оглавление](/course/rust?lang=ru) · [Следующий урок](/read/rust-39-scenarios?lang=ru)
