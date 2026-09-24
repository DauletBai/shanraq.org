# Автоматические проверки на трёх системах

_Лид (summary):_ **Урок 57. Запускать форматирование, Clippy и тесты после каждой отправки кода.**

Повторите [изолированные тесты](/read/rust-56-testing?lang=ru) и [первый проект Cargo](/read/rust-06-cargo?lang=ru). Репозиторий — место, где GitHub хранит файлы проекта и историю их изменений. Работайте в отдельном учебном репозитории: файл задания запускает проверку при отправке кода.

## Образ и опорная карта

Перед выпуском тетрадь читают три корректора: каждый смотрит свой экземпляр, а замечание возвращает автору. **Карта:** изменить код → проверить локально → отправить → три системы собирают и тестируют → прочитать результат → исправить. Это не доказательство отсутствия всех ошибок: проверяются лишь записанные правила и сценарии.

**CI** (непрерывная интеграция) — автоматический запуск проверок для отправленного изменения. **GitHub Actions** читает файл в `.github/workflows` в корне репозитория. В YAML `name:` даёт имя, отступы объединяют поля, `-` начинает элемент списка, `on` задаёт событие, `jobs` — набор работ, `steps` — шаги, `uses` подключает готовое действие, `run` запускает команду. `matrix.os` означает три варианта системы; запись `${{ matrix.os }}` подставляет один из них в `runs-on`. `actions/checkout@v7` берёт исходники. `rustup toolchain install` устанавливает конкретную версию компилятора и компоненты; `rust-toolchain.toml` хранит этот выбор рядом с проектом. `fmt --check` проверяет оформление без изменения файлов, Clippy ищет подозрительные конструкции, `test --locked` выполняет тесты с закреплёнными зависимостями. `-- -D warnings` передаёт Clippy требование считать предупреждения ошибками.

В команде `cargo +1.97.0` знак `+` выбирает установленную версию Rust 1.97.0; `--all-targets` включает тестовые цели, а `--locked` запрещает менять `Cargo.lock`. В YAML квадратные скобки перечисляют события или системы, `@v7` указывает версию действия.

## Запустите и разберите

Скопируйте Rust-код в новый `src/main.rs`; `Cargo.toml` может остаться без зависимостей. Сначала выполните `cargo fmt --check`, `cargo clippy --all-targets -- -D warnings`, `cargo test --locked`, затем `cargo run`. Первый блок содержит один тест. После этого создайте файл `.github/workflows/organizer.yml` с примером ниже и `rust-toolchain.toml` с закреплённой версией. Команды в YAML должны находиться в корне именно этого проекта. Первая установка 1.97.0 требует сети. Изменение YAML проверяется на GitHub после отправки; зелёный знак не заменяет чтение задания человеком.

```rust
fn normalized_title(text: &str) -> Option<&str> {
    let trimmed = text.trim();
    if trimmed.is_empty() {
        None
    } else {
        Some(trimmed)
    }
}

fn main() {
    let title = normalized_title("  Купить книгу  ").unwrap_or("Пусто");
    println!("Название: {title}");
}

#[test]
fn trims_outer_spaces() {
    assert_eq!(normalized_title("  Книга  "), Some("Книга"));
}
```
```text
Название: Купить книгу
```

Файлы для корня учебного проекта:

`rust-toolchain.toml`:

```toml
[toolchain]
channel = "1.97.0"
components = ["rustfmt", "clippy"]
profile = "minimal"
```

`.github/workflows/organizer.yml`:

```yaml
name: Organizer checks
on: [push, pull_request]
jobs:
  verify:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v7
      - run: rustup toolchain install 1.97.0 --profile minimal --component rustfmt,clippy
      - run: cargo +1.97.0 fmt --check
      - run: cargo +1.97.0 clippy --all-targets -- -D warnings
      - run: cargo +1.97.0 test --locked
```

## Вопросы для воспроизведения

1. Почему локальный `cargo test` нужен до отправки?
2. Что означает `matrix.os`?
3. Почему успешный CI не доказывает правильность замысла?

## Задание

**Обязательное.** Добавьте тест пустого названия и тест текста с буквой вне ASCII. Убедитесь, что локально проходят три теста. Создайте оба файла настройки из примера, затем в своём учебном репозитории отправьте изменение и найдите результаты для всех трёх систем. Если GitHub недоступен, сохраните файлы и объясните по карте, какие шаги ещё не выполнены.

## Ответы

`normalized_title("   ")` должен вернуть `None`. Для Unicode проверьте строку с `é`. Не считайте запуск только на одной системе полной матрицей.

<!-- task-answer -->
```rust
fn normalized_title(text: &str) -> Option<&str> {
    let trimmed = text.trim();
    if trimmed.is_empty() {
        None
    } else {
        Some(trimmed)
    }
}

fn main() {
    let title = normalized_title("  Купить книгу  ").unwrap_or("Пусто");
    println!("Название: {title}");
}

#[test]
fn trims_outer_spaces() {
    assert_eq!(normalized_title("  Книга  "), Some("Книга"));
}

#[test]
fn rejects_blank_title() {
    assert_eq!(normalized_title("   "), None);
}

#[test]
fn keeps_unicode_title() {
    assert_eq!(normalized_title("  Café  "), Some("Café"));
}
```
```text
Название: Купить книгу
```

## После проверки

Если одна система показывает ошибку, откройте именно её шаг и воспроизведите команду локально. `fmt`, Clippy и тесты решают разные задачи; не заглушайте ошибки ради зелёного значка. Номер toolchain и версии действий нужно пересматривать при обновлении курса. [Документация Cargo fmt](https://doc.rust-lang.org/cargo/commands/cargo-fmt.html).

[Предыдущий урок](/read/rust-56-testing?lang=ru) · [Оглавление](/course/rust?lang=ru)
