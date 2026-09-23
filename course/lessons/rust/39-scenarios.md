# Сценарии: проверяем цепочку действий

_Лид (summary):_ **Урок 39. Проверить не одну функцию, а последствия нескольких действий и ошибки без изменения данных.**

Предварительно: уроки 16, 32–33 и 38. Если забыли `#[test]`, вернитесь к [уроку 16](/read/rust-16-first-test?lang=ru).

## Образ и опорная карта

Проверяя замок, недостаточно один раз повернуть ключ: нужно открыть, закрыть и попробовать неподходящий ключ. **Сценарный тест** так же проверяет цепочку действий органайзера. Образ ограничен: тесты охватывают выбранные случаи, но не доказывают отсутствие всех ошибок.

**Начальное состояние → действие → ожидаемый результат → новое состояние; при ошибке → прежнее состояние.** `#[test]` и `assert_eq!` знакомы по уроку 16. Каждый тест создаёт собственный `Organizer`, поэтому порядок запуска тестов не должен влиять на результат. Выражение `Err(String::from(...))` — ожидаемая ошибка. Сравнение `Result` через `assert_eq!` проверяет и вариант, и его содержимое. При сравнении `o.tasks[0].title` макрос читает строку, не вынимая её из списка. Успех операции ещё не гарантирует правильное состояние: после удаления мы отдельно проверяем оставшийся ID и значение `next_id`. Проверка при ошибке должна удостовериться, что ни задачи, ни счётчик не поменялись.

## Два независимых сценария

Сохраните прежний `src/main.rs` и замените полностью. `cargo run` покажет короткое сообщение. Затем выполните `cargo test`: Cargo запустит функции с `#[test]`. Число строк служебного вывода Cargo может отличаться; важен итог: все тесты прошли. Новых зависимостей нет. Для короткой проверки `Task` здесь содержит только `id` и `title`: поле `done` из урока 38 временно опущено и вернётся в уроке 40. Каждый урок заменяет учебный `main.rs` целиком, а не переносит данные прежнего запуска.

```rust
struct Task {
    id: u32,
    title: String,
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
}
fn main() {
    let mut o = Organizer::new();
    o.add("Читать").expect("fixed title");
    o.add("Писать").expect("fixed title");
    o.remove(1).expect("fixed ID");
    println!("Осталась: {}", o.tasks[0].title);
}

#[test]
fn ids_survive_removal() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Читать"), Ok(1));
    assert_eq!(o.add("Писать"), Ok(2));
    assert_eq!(o.remove(1), Ok(()));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 2);
    assert_eq!(o.tasks[0].title, "Писать");
    assert_eq!(o.next_id, 3);
}

#[test]
fn overflow_changes_nothing() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Читать"), Ok(1));
    o.next_id = u32::MAX;
    assert_eq!(o.add("Писать"), Err(String::from("Номера закончились")));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 1);
    assert_eq!(o.next_id, u32::MAX);
}
```
```text
Осталась: Писать
```

Первая проверка ловит ошибку «ID равен месту в списке». Вторая искусственно ставит счётчик на границу и проверяет, что `add` отказывает до изменения данных. В учебном примере поле `title` сверяется тоже: правильный ID у неверной карточки был бы ложным успехом. Данные исчезают после завершения процесса, как и в уроке 38.

## Проверка понимания

1. Почему проверки `Ok(2)` после удаления ID 1 недостаточно?
2. Что дополнительно проверяем после `Err`?
3. Зависят ли эти тесты друг от друга?

Нужно проверить сохранённый ID и счётчик; состояние не изменилось; каждый создаёт свою новую структуру. Восстановите опорную карту.

## Задание

**Обязательное.** Добавьте отдельный тест: после двух успешных добавлений попытка удалить ID 99 возвращает ошибку, оба исходных ID остаются, а следующий успешный `add` выдаёт ID 3. Вызовите `cargo test`. Не проверяйте только сообщение об ошибке.

## Подсказка

Создайте новый `Organizer` внутри `#[test] fn invalid_remove_keeps_ids()`. Используйте `assert_eq!(o.remove(99), Err(...))`, затем проверьте `o.tasks.len()`, оба `id` и `o.add(...)`.

## Эталон после попытки

<!-- task-answer -->
```rust
struct Task {
    id: u32,
    title: String,
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
}
fn main() {
    let mut o = Organizer::new();
    o.add("Читать").expect("fixed title");
    o.add("Писать").expect("fixed title");
    o.remove(1).expect("fixed ID");
    println!("Осталась: {}", o.tasks[0].title);
}

#[test]
fn ids_survive_removal() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Читать"), Ok(1));
    assert_eq!(o.add("Писать"), Ok(2));
    assert_eq!(o.remove(1), Ok(()));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 2);
    assert_eq!(o.tasks[0].title, "Писать");
    assert_eq!(o.next_id, 3);
}

#[test]
fn overflow_changes_nothing() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Читать"), Ok(1));
    o.next_id = u32::MAX;
    assert_eq!(o.add("Писать"), Err(String::from("Номера закончились")));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 1);
    assert_eq!(o.next_id, u32::MAX);
}

#[test]
fn invalid_remove_keeps_ids() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Читать"), Ok(1));
    assert_eq!(o.add("Писать"), Ok(2));
    assert_eq!(o.remove(99), Err(String::from("ID не найден")));
    assert_eq!(o.tasks.len(), 2);
    assert_eq!(o.tasks[0].id, 1);
    assert_eq!(o.tasks[1].id, 2);
    assert_eq!(o.add("Читать"), Ok(3));
}
```
```text
Осталась: Писать
```

Вывод здесь относится к `cargo run`; `cargo test` отдельно сообщает о трёх пройденных тестах. [Официальный разбор тестов Rust](https://doc.rust-lang.org/book/ch11-01-writing-tests.html).

[Предыдущий урок](/read/rust-38-ids?lang=ru) · [Оглавление](/course/rust?lang=ru) · [Следующий урок](/read/rust-40-memory-checkpoint?lang=ru)
