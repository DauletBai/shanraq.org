# Трейты: общий договор о поведении

_Лид (summary):_ **Урок 46. Задать общее поведение двум типам и ограничить обобщённую функцию.**

Предварительно: [обобщённые типы](/read/rust-45-generics?lang=ru) и [методы](/read/rust-28-methods?lang=ru). Сохраните прежний проект; пример снова хранит данные только в памяти.

## Образ и опорная карта

На кружке два участника умеют представляться, хотя имена у них хранятся по-разному. Правило «назови себя» похоже на **трейт** (trait): он перечисляет доступные действия. Участник обязан выполнить договор; сам договор не говорит, как он это сделает. Предел образа: компилятор проверяет сигнатуру метода, но не обещает, что полученный текст полезен человеку.

**trait — обещание метода → impl — реализация для типа → T: Summary — допуск к общей функции.** `trait Summary` задаёт метод `summary(&self) -> String`; точка с запятой означает, что тела здесь нет. `impl Summary for Task` даёт тело именно для `Task`. В `show<T: Summary>` двоеточие означает ограничение типа: `show` принимает любой `T`, который выполняет этот договор. `&T` лишь одалживает значение. `clone()` создаёт свою строку для возвращаемого `String`; исходный заголовок остаётся у задачи.

## Одна функция для задачи

```rust
trait Summary {
    fn summary(&self) -> String;
}

struct Task {
    title: String,
}

impl Summary for Task {
    fn summary(&self) -> String {
        self.title.clone()
    }
}

fn show<T: Summary>(item: &T) {
    println!("{}", item.summary());
}

fn main() {
    let task = Task {
        title: String::from("Читать Rust"),
    };
    show(&task);
}
```
```text
Читать Rust
```

`show(&task)` вызывает реализацию для `Task`. Если убрать `impl Summary for Task`, компилятор отвергнет вызов: функция вправе рассчитывать на метод только у типов с этим трейтом. У готовых трейтов вроде `Debug` реализацию иногда создаёт атрибут `#[derive(Debug)]`: это запрос компилятору сгенерировать подходящий код. Он работает лишь для поддерживаемых трейтов и полей; не делает произвольный трейт автоматически. Здесь реализация написана вручную, потому что итоговый текст — наше решение.

## Как работает derive

`Debug` показывает поля для разработчика; это не готовый текст для пользователя. `#[derive(Debug)]` — атрибут перед объявлением структуры; он просит компилятор реализовать стандартный трейт. Формат `{:?}` вызывает это отладочное представление.

```rust
#[derive(Debug)]
struct Task {
    title: String,
}

fn main() {
    let task = Task {
        title: String::from("Читать Rust"),
    };
    println!("{task:?}");
}
```
```text
Task { title: "Читать Rust" }
```

## Вопросы для воспроизведения

1. Что обещает `trait Summary`?
2. За что отвечает `impl Summary for Task`?
3. Почему `show<T: Summary>` не примет тип без этого трейта?

## Задание

**Обязательное.** Добавьте `struct Note { text: String }`, реализуйте для него `Summary`, создайте заметку «Купить книгу» и передайте её в ту же `show`. Сначала предскажите две строки вывода.

## Ответы

Скопируйте форму `impl Summary for Task`, замените `Task` на `Note` и `title` на `text`. Новый тип должен вернуть свой `String`, не отнимая текст у заметки.

<!-- task-answer -->
```rust
trait Summary {
    fn summary(&self) -> String;
}

struct Task {
    title: String,
}

impl Summary for Task {
    fn summary(&self) -> String {
        self.title.clone()
    }
}

struct Note {
    text: String,
}

impl Summary for Note {
    fn summary(&self) -> String {
        self.text.clone()
    }
}

fn show<T: Summary>(item: &T) {
    println!("{}", item.summary());
}

fn main() {
    let task = Task {
        title: String::from("Читать Rust"),
    };
    let note = Note {
        text: String::from("Купить книгу"),
    };
    show(&task);
    show(&note);
}
```
```text
Читать Rust
Купить книгу
```

## После проверки

Обе строки проходят через одну `show`, а конкретный метод выбирается по типу. Если неясна буква `T`, вернитесь к [уроку 45](/read/rust-45-generics?lang=ru). [Официальная глава о трейтах](https://doc.rust-lang.org/book/ch10-02-traits.html).

[Предыдущий урок](/read/rust-45-generics?lang=ru) · [Оглавление](/course/rust?lang=ru)
