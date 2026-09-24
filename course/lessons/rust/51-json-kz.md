# JSON пішімі

_Қысқаша (summary):_ **Сабақ 51. Тапсырмаларды JSON мәтініне айналдырып, қайта оқығанда тексеру.**

[Task құрылымын](/read/rust-27-structs?lang=kz), [сыртқы кітапханаларды](/read/rust-48-dependencies?lang=kz) және [файл оқуды](/read/rust-50-files?lang=kz) еске түсіріңіз. Мұнда дерек жадтағы мәтін түрінде болады; дискіге жазу келесі сабақта.

## Таныс бейне және тірек сызба

Тапсырмалар тізімі атауы жазылған хатқалтадағы қағаздарға ұқсайды. JSON — өріс атаулары мен мәндерді ортақ ережемен жазатын мәтіндік пішім. Оны басқа бағдарлама оқи алады, бірақ хатқалтаның өзі ішіндегі деректің дұрыстығын дәлелдемейді. **Тірек сызба:** мәндер → JSON → JSON-ды оқу → мағынасын тексеру → дайын тапсырмалар. Осы жолды мәтінге қарамай айтып шығыңыз.

`serde` құрылымдарды дерекке және кері түрлендіреді; `serde_json` JSON-мен жұмыс істейді. **Сериялау** — мәнді мәтінге айналдыру, **десериялау** — мәтіннен мәнді қайта құрастыру. `#[derive(Serialize, Deserialize)]` құрылымға арналған қалыпты іске асыруды жасатады; `#[serde(deny_unknown_fields)]` танылмаған өрісті қабылдамайды. `#[...]` белгісі кейінгі жариялауға арналған компилятор нұсқауы. `version` — файл құрылымының нұсқа нөмірі; құрылым өзгерсе, ескі файлды жаңасы деп қабылдауға болмайды. `next_id` барлық бар ID-ден үлкен болуы керек. `HashSet` — қайталанбайтын мәндер жиыны; `insert` бұрын кездескен нөмір үшін `false` қайтарады. `?` қатені шақырған жерге өткізеді. `map_err` кітапхана қатесін `run` үшін мәтінге айналдырады.

## Іске қосып, талдаңыз

Жеке Cargo жобасын ашыңыз. `Cargo.toml` файлына төмендегі тәуелділіктерді жазыңыз, `src/main.rs` файлын алғашқы Rust үзіндісімен алмастырып, `cargo run` орындаңыз. Cargo кітапханаларды алғаш жинағанда жүктейді; `Cargo.lock` таңдалған нұсқаларды бекітеді. Үлгіде бір тапсырма бар. `to_string` JSON мәтінін жасайды, `from_str` қайта оқиды, ал `validate` JSON білмейтін мән ережелерін тексереді. Бұл үлгідегі өрістер реті құрылымдағы ретпен шығады; басқа бағдарлама оларды өзгеше орналастыра алады.

`Cargo.toml` файлына тәуелділіктерді қосыңыз (төмендегі жолдар `[package]` бөлімінен кейін):

Тәуелділік жазбасындағы `version = "=1.0.228"` ішкі `=` таңбасы дәл осы нұсқаны талап етеді, ал `features = ["derive"]` өрісі `derive` арқылы іске асыру жасау мүмкіндігін қосады; мұндағы тік жақшалар мүмкіндіктер тізімін білдіреді. Бұл — оқу үлгісіндегі нұсқаларды бекіту, олардың мәңгі өзекті болатынына уәде емес.

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
            title: String::from("Кітап сатып алу"),
            done: false,
        }],
    };
    let json = serde_json::to_string(&source).map_err(|error| error.to_string())?;
    let loaded: Organizer = serde_json::from_str(&json).map_err(|error| error.to_string())?;
    validate(&loaded)?;
    println!("{json}");
    println!("Қалпына келген тапсырма саны: {}", loaded.tasks.len());
    Ok(())
}

fn main() {
    if let Err(error) = run() {
        eprintln!("Қате: {error}");
        std::process::exit(1);
    }
}
```
```text
{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Кітап сатып алу","done":false}]}
Қалпына келген тапсырма саны: 1
```

## Бетке қарамай жауап беріңіз

1. JSON-ды оқу мен оның мағынасын тексерудің айырмасы қандай?
2. Тізім бар болса, `next_id` не үшін қажет?
3. Танылмаған өрісті неге қабылдамаған дұрыс?

## Тапсырма

**Міндетті.** ID-і 2 болатын екінші тапсырманы қосып, `next_id` мәнін 3 етіңіз; қайта оқылған тапсырма санын шығарыңыз. Содан кейін оқылған екінші тапсырманың ID-ін 1-ге ауыстырып, тексеру қатені қайтаратынын көрсетіңіз. Алдымен екі нәтижені болжаңыз.

## Жауаптар

Деректі `to_string` алдында өзгертіңіз. `from_str` соңында қайтадан `Organizer` аласыз. Қате нұсқа үшін `loaded.tasks[1].id` мәнін өзгертіп, тексеруді қайта шақырыңыз. `validate(&bad).is_err()` нәтижесін тексеріңіз.

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
                title: String::from("Кітап сатып алу"),
                done: false,
            },
            Task {
                id: 2,
                title: String::from("Досқа қоңырау шалу"),
                done: false,
            },
        ],
    };
    let json = serde_json::to_string(&source).map_err(|error| error.to_string())?;
    let mut loaded: Organizer = serde_json::from_str(&json).map_err(|error| error.to_string())?;
    validate(&loaded)?;
    println!("{json}");
    println!("Қалпына келген тапсырма саны: {}", loaded.tasks.len());
    loaded.tasks[1].id = 1;
    println!(
        "Қайталанған ID қабылданбады: {}",
        validate(&loaded).is_err()
    );
    Ok(())
}

fn main() {
    if let Err(error) = run() {
        eprintln!("Қате: {error}");
        std::process::exit(1);
    }
}
```
```text
{"version":1,"next_id":3,"tasks":[{"id":1,"title":"Кітап сатып алу","done":false},{"id":2,"title":"Досқа қоңырау шалу","done":false}]}
Қалпына келген тапсырма саны: 2
Қайталанған ID қабылданбады: true
```

## Тексергеннен кейін

JSON жазылу ережесін, `validate` органайзер ережесін тексереді. Бүлінген деректі бос тізіммен алмастырмаңыз. [50-сабаққа оралу](/read/rust-50-files?lang=kz). [Serde құжаттамасы](https://serde.rs/derive.html).

[Алдыңғы сабақ](/read/rust-50-files?lang=kz) · [Мазмұн](/course/rust?lang=kz)
