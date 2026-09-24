# Тапсырма файлына арналған толық тексерулер

_Қысқаша (summary):_ **Сабақ 56. Сақталған тапсырмаларды, бүлінген файлды, қайта оқуды және көптілді мәтінді тексеру.**

[JSON мен деректі тексеруді](/read/rust-51-json?lang=kz), [сақтық көшірмені](/read/rust-53-backup?lang=kz) және [алғашқы автоматты сынақты](/read/rust-16-first-test?lang=kz) еске түсіріңіз. Жаттығуда жүйенің уақытша қалтасындағы бөлек файлдарды ғана қолданыңыз.

## Таныс бейне және тірек сызба

Жолға шығарда көлікті тұрған жерінде ғана қарамай, тоқтағаннан кейін қайта оталуын, қосалқы кілтті және жаңбырдағы жолды да тексересіз. Файл сынақтары да әртүрлі жағдайды қамтиды. **Тірек сызба:** бөлек оқу жолы → дерек дайындау → әрекетті орындау → күтілетін нәтижемен салыстыру → тек өз файлыңызды өшіру. Осы тізбекті жатқа айтып шығыңыз. Ұқсастықтың шегі: бірнеше сынақ барлық қатенің жоқтығын дәлелдемейді.

**Сынақтар жиыны** — әрекетті тексеретін бірнеше дербес сынақ. **Оқшаулау** сынақ шын деректі бүлдірмейтінін және көрші сынаққа тәуелді еместігін білдіреді. ASCII — латынның негізгі әріптері, цифрлары мен белгілерін қамтитын шағын ескі кесте; `Ә` онда жоқ. Мұндағы сынақ сол кестеден тыс әріпті сақтауды тексереді. `practice_path` файл атына үдеріс нөмірі мен жағдай белгісін қосады; `unicode`, `damaged`, `duplicate`, `again` белгілері бір-бірінен бөлек. Бұл сынақтар қатар орындалғанда жолдардың қақтығысу қаупін азайтады. `#[test]` қызметтің сынақ екенін Cargo-ға білдіреді; `cargo test` соларды, `cargo run` `main` қызметін орындайды. `assert_eq!` мәндерді салыстырады, `assert!` шарттың рас болуын талап етеді. Сынақтағы `unwrap()` күтпеген қатені сынақтың сәтсіздігіне айналдырады; кәдімгі бағдарламада қатені адамға түсіндіру керек.

## Іске қосып, талдаңыз

Жеке Cargo жобасын ашыңыз. `Cargo.toml` ішіндегі `[package]` бөлімінен соң төмендегі тәуелділіктерді қосыңыз: нұсқа ішіндегі `=` дәл осы нұсқаны талап етеді, ал `features = ["derive"]` тізімі Serde іске асыруын жасау мүмкіндігін қосады. Алғашқы Rust үзіндісін `src/main.rs` файлына қойыңыз. `cargo run` оқылған тапсырма санын шығарады. Содан кейін `cargo test` орындаңыз: екі сынақ өтуі тиіс. Біріншісі қазақ әрпі бар мәтінді жазып, қайта оқиды; екіншісі бүлінген JSON-ды өзгертпей қалдырады. `read_checked` алдымен UTF-8 пен JSON-ды оқиды, содан соң нұсқаны, ID мен атауларды тексереді. `Err` бос органайзерге айналмайды.

```toml
[dependencies]
serde = { version = "=1.0.228", features = ["derive"] }
serde_json = "=1.0.149"
```

```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::fs;
use std::path::{Path, PathBuf};

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

fn read_checked(path: &Path) -> Result<Organizer, String> {
    let text = fs::read_to_string(path).map_err(|error| error.to_string())?;
    let data: Organizer = serde_json::from_str(&text).map_err(|error| error.to_string())?;
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

fn practice_path(label: &str) -> PathBuf {
    std::env::temp_dir().join(format!(
        "rust-course-56-{}-{label}.json",
        std::process::id()
    ))
}

fn sample(title: &str) -> Organizer {
    Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from(title),
            done: false,
        }],
    }
}

fn main() -> Result<(), String> {
    let path = practice_path("main");
    let text = serde_json::to_string(&sample("Кітап сатып алу")).map_err(|e| e.to_string())?;
    fs::write(&path, text).map_err(|e| e.to_string())?;
    let loaded = read_checked(&path)?;
    println!("Оқылған тапсырма саны: {}", loaded.tasks.len());
    fs::remove_file(path).map_err(|e| e.to_string())?;
    Ok(())
}

#[test]
fn unicode_round_trip() {
    let path = practice_path("unicode");
    let text = serde_json::to_string(&sample("Әжеге хат жазу")).unwrap();
    fs::write(&path, text).unwrap();
    let loaded = read_checked(&path).unwrap();
    assert_eq!(loaded.tasks[0].title, "Әжеге хат жазу");
    fs::remove_file(path).unwrap();
}

#[test]
fn damaged_json_is_rejected() {
    let path = practice_path("damaged");
    fs::write(&path, "{broken").unwrap();
    assert!(read_checked(&path).is_err());
    assert_eq!(fs::read_to_string(&path).unwrap(), "{broken");
    fs::remove_file(path).unwrap();
}
```
```text
Оқылған тапсырма саны: 1
```

## Бетке қарамай жауап беріңіз

1. Оқшауланған сынақтың өзіңіздің шын файлыңызбен жұмыс істеуден айырмасы қандай?
2. Бүлінген файл сынағы қате шыққаннан кейін мазмұнды неге салыстырады?
3. `unwrap()` күтпеген `Err` алса, сынаққа не болады?

## Тапсырма

**Міндетті.** Бірдей ID-і бар екі тапсырманы жазып, `read_checked` қабылдамайтынын тексеретін сынақ қосыңыз. Бір файлды екі рет оқып, атаудың өзгермегенін тексеретін екінші сынақ қосыңыз. Әр сынақтың жол белгісі бөлек болсын; соңында тек өзі жасаған файлды өшірсін. `cargo test` орындап, төрт сынақтың өткенін көріңіз.

## Жауаптар

Қайталанған ID үшін `sample` мәнін алып, `next_id = 3` етіп, `id: 1` болатын екінші тапсырманы қосып, JSON жазыңыз. Қайта оқу үшін `first.tasks[0].title` мен `second.tasks[0].title` мәндерін салыстырыңыз.

<!-- task-answer -->
```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::fs;
use std::path::{Path, PathBuf};

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

fn read_checked(path: &Path) -> Result<Organizer, String> {
    let text = fs::read_to_string(path).map_err(|error| error.to_string())?;
    let data: Organizer = serde_json::from_str(&text).map_err(|error| error.to_string())?;
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

fn practice_path(label: &str) -> PathBuf {
    std::env::temp_dir().join(format!(
        "rust-course-56-{}-{label}.json",
        std::process::id()
    ))
}

fn sample(title: &str) -> Organizer {
    Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from(title),
            done: false,
        }],
    }
}

fn main() -> Result<(), String> {
    let path = practice_path("main");
    let text = serde_json::to_string(&sample("Кітап сатып алу")).map_err(|e| e.to_string())?;
    fs::write(&path, text).map_err(|e| e.to_string())?;
    let loaded = read_checked(&path)?;
    println!("Оқылған тапсырма саны: {}", loaded.tasks.len());
    fs::remove_file(path).map_err(|e| e.to_string())?;
    Ok(())
}

#[test]
fn unicode_round_trip() {
    let path = practice_path("unicode");
    let text = serde_json::to_string(&sample("Әжеге хат жазу")).unwrap();
    fs::write(&path, text).unwrap();
    let loaded = read_checked(&path).unwrap();
    assert_eq!(loaded.tasks[0].title, "Әжеге хат жазу");
    fs::remove_file(path).unwrap();
}

#[test]
fn damaged_json_is_rejected() {
    let path = practice_path("damaged");
    fs::write(&path, "{broken").unwrap();
    assert!(read_checked(&path).is_err());
    assert_eq!(fs::read_to_string(&path).unwrap(), "{broken");
    fs::remove_file(path).unwrap();
}

#[test]
fn duplicate_ids_are_rejected() {
    let path = practice_path("duplicate");
    let mut data = sample("Кітап сатып алу");
    data.next_id = 3;
    data.tasks.push(Task {
        id: 1,
        title: String::from("second"),
        done: false,
    });
    fs::write(&path, serde_json::to_string(&data).unwrap()).unwrap();
    assert!(read_checked(&path).is_err());
    fs::remove_file(path).unwrap();
}

#[test]
fn second_read_keeps_data() {
    let path = practice_path("again");
    fs::write(
        &path,
        serde_json::to_string(&sample("Кітап сатып алу")).unwrap(),
    )
    .unwrap();
    let first = read_checked(&path).unwrap();
    let second = read_checked(&path).unwrap();
    assert_eq!(first.tasks[0].title, second.tasks[0].title);
    fs::remove_file(path).unwrap();
}
```
```text
Оқылған тапсырма саны: 1
```

## Тексергеннен кейін

Сынақ тоқтап қалып, файл қалдырса, ішін көргеннен кейін тек оқу атауы бар файлды өшіріңіз. Қатар орындалатын сынақтардың реті өзгеруі мүмкін: бәріне бір жолды бермеңіз. Жазу қатесі мен қуат үзілуін тексеру үшін бөлек жағдайлар мен орта қажет; бұл сабақ аталған қасиеттерді ғана растайды. [cargo test құжаттамасы](https://doc.rust-lang.org/cargo/commands/cargo-test.html).

[Алдыңғы сабақ](/read/rust-55-documentation?lang=kz) · [Мазмұн](/course/rust?lang=kz)
