# Үш жүйедегі автоматты тексерулер

_Қысқаша (summary):_ **Сабақ 57. Код жіберілген сайын пішімдеуді, Clippy-ді және сынақтарды орындау.**

[Оқшауланған сынақтарды](/read/rust-56-testing?lang=kz) және [алғашқы Cargo жобасын](/read/rust-06-cargo?lang=kz) еске түсіріңіз. Код қоры (repository) — GitHub-та жоба файлдары мен олардың өзгеріс тарихы сақталатын орын. Бөлек оқу қорын пайдаланыңыз: тапсырма файлы код жіберілгенде тексеруді бастайды.

## Таныс бейне және тірек сызба

Баспаға дейін дәптерді үш түзетуші қарайды: әрқайсысы өз данасын оқиды да, ескертуді авторға қайтарады. **Тірек сызба:** кодты өзгерту → жергілікті тексеру → жіберу → үш жүйеде жинау және сынау → нәтижені оқу → түзету. Бұл барлық қатенің жоқтығын дәлелдемейді: тек жазылған ережелер мен жағдайлар тексеріледі.

**CI** (үздіксіз біріктіру) — жіберілген өзгерісті автоматты тексеру. **GitHub Actions** қордың түбіндегі `.github/workflows` файлдарын оқиды. YAML ішіндегі `name:` атау береді, шегініс өрістерді топтайды, `-` тізім элементін бастайды, `on` іске қосу жағдайын, `jobs` жұмыстарды, `steps` қадамдарды көрсетеді, `uses` дайын әрекетті қосады, `run` пәрменді орындайды. `matrix.os` үш жүйе нұсқасын білдіреді; `${{ matrix.os }}` жазуы соның бірін `runs-on` орнына қояды. `actions/checkout@v7` бастапқы кодты алады. `rustup toolchain install` компилятордың нақты нұсқасы мен қосымша құралдарын орнатады; `rust-toolchain.toml` таңдауды жобамен бірге сақтайды. `fmt --check` файлды өзгертпей пішімді тексереді, Clippy күмәнді жазылымдарды іздейді, `test --locked` бекітілген тәуелділіктермен сынақтарды орындайды. `-- -D warnings` Clippy ескертулерін қате деп есептеуді талап етеді.

`cargo +1.97.0` пәрменіндегі `+` таңбасы орнатылған Rust 1.97.0 нұсқасын таңдайды; `--all-targets` сынақ нысандарын да қамтиды, `--locked` `Cargo.lock` өзгерісіне жол бермейді. YAML ішіндегі тік жақшалар жағдайлар мен жүйелер тізімін, `@v7` әрекет нұсқасын көрсетеді.

## Іске қосып, талдаңыз

Rust кодын жаңа жобаның `src/main.rs` файлына көшіріңіз; `Cargo.toml` тәуелділіксіз бола алады. Алдымен `cargo fmt --check`, `cargo clippy --all-targets -- -D warnings`, `cargo test --locked`, содан соң `cargo run` орындаңыз. Алғашқы үзіндіде бір сынақ бар. Одан кейін төмендегі `.github/workflows/organizer.yml` және `rust-toolchain.toml` файлдарын жоба түбінде жасаңыз. YAML ішіндегі пәрмендер осы жобаның түбінде орындалуы тиіс. 1.97.0 нұсқасын алғаш орнатуға желі керек. GitHub-қа жіберген соң YAML-дың нәтижесін сонда көресіз; жасыл белгі тапсырманы адам оқуының орнын баспайды.

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
    let title = normalized_title("  Кітап сатып алу  ").unwrap_or("Бос");
    println!("Атауы: {title}");
}

#[test]
fn trims_outer_spaces() {
    assert_eq!(normalized_title("  Кітап  "), Some("Кітап"));
}
```
```text
Атауы: Кітап сатып алу
```

Оқу жобасының түбіне арналған файлдар:

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

## Бетке қарамай жауап беріңіз

1. Жіберерден бұрын жергілікті `cargo test` неге керек?
2. `matrix.os` нені білдіреді?
3. CI сәтті өтсе де, ойдың дұрыстығы неге дәлелденбейді?

## Тапсырма

**Міндетті.** Бос атау мен ASCII-ден тыс әрпі бар мәтінге екі сынақ қосыңыз. Жергілікті үш сынақтың өткенін тексеріңіз. Үлгідегі екі баптау файлын жасап, өз оқу қорыңызға жіберіңіз де үш жүйенің нәтижесін табыңыз. GitHub қолжетімсіз болса, файлдарды сақтап, тірек сызба бойынша қай қадамдар орындалмағанын түсіндіріңіз.

## Жауаптар

`normalized_title("   ")` `None` қайтаруы тиіс. Көптілді мәтін үшін `Ә` әрпін тексеріңіз. Бір жүйедегі іске қосуды толық үш жүйелік тексеру деп атамаңыз.

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
    let title = normalized_title("  Кітап сатып алу  ").unwrap_or("Бос");
    println!("Атауы: {title}");
}

#[test]
fn trims_outer_spaces() {
    assert_eq!(normalized_title("  Кітап  "), Some("Кітап"));
}

#[test]
fn rejects_blank_title() {
    assert_eq!(normalized_title("   "), None);
}

#[test]
fn keeps_unicode_title() {
    assert_eq!(normalized_title("  Әже  "), Some("Әже"));
}
```
```text
Атауы: Кітап сатып алу
```

## Тексергеннен кейін

Бір жүйеде қате шықса, дәл сол жүйедегі қадамды ашып, пәрменді жергілікті қайталаңыз. `fmt`, Clippy және сынақтар әртүрлі міндет атқарады; жасыл белгі үшін қатені жасырып қоймаңыз. Құралдар мен әрекеттердің нұсқасын курс жаңарғанда қайта қарау керек. [Cargo fmt құжаттамасы](https://doc.rust-lang.org/cargo/commands/cargo-fmt.html).

[Алдыңғы сабақ](/read/rust-56-testing?lang=kz) · [Мазмұн](/course/rust?lang=kz)
