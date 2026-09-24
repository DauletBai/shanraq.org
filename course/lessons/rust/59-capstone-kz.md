# Өз бетінше кеңейту: тапсырма басымдығы

_Қысқаша (summary):_ **Сабақ 59. Ескі JSON-ды оқуды сақтай отырып, басымдық қосу және жарамсыз мәнді қабылдамау.**

[Құрылымдарды](/read/rust-27-structs?lang=kz), [JSON мен тексеруді](/read/rust-51-json?lang=kz), [файл сынақтарын](/read/rust-56-testing?lang=kz) еске түсіріңіз. Алдымен әрекет ережесін жазыңыз, содан соң код құрыңыз: бұл жауап үлгісі тапсырмадан кейін берілетін дербес жұмыс.

## Таныс бейне және тірек сызба

Тапсырма қағазына шұғылдық белгісі қойылады. Белгісіз ескі қағаз өздігінен шұғыл болмайды: оған қалыпты белгі беріледі. **Тірек сызба:** ережені жазу → өріс қосу → ескі жазбаға мән беру → жарамсыз мәнді қабылдамау → іріктеуді тексеру → адамға түсіндіру. Ұқсастықтың шегі: қағаздағы белгі тапсырманың шынымен маңызды екенін анықтамайды; оны адам шешеді.

**Басымдық** — 1-ден 3-ке дейінгі сан: 1 төмен, 2 қалыпты, 3 шұғыл. Ескі JSON-да `priority` өрісі жоқ. `#[serde(default = "ordinary_priority")]` белгісі өріс болмаса, Serde `ordinary_priority` қызметін шақыратынын білдіреді. Ол файлда анық жазылған жарамсыз санды түзетпейді. `1..=3` — екі шеті де кіретін аралық. `high_open` шұғыл әрі әлі аяқталмаған тапсырмаларды санайды. Бұл үйлесімді қосымша үшін ескі файл оқылса және жаңа өрістің мағынасы анық болса, `version` мәні 1 күйінде қала алады; үйлеспейтін өзгеріске нұсқа мен деректі көшіру туралы бөлек шешім керек.

## Іске қосып, талдаңыз

Бөлек Cargo жобасына төмендегі тәуелділіктерді қосып, алғашқы Rust үзіндісін `src/main.rs` файлына қойыңыз. `cargo run` және `cargo test` орындаңыз. Бастапқы бағдарлама басымдық өрісі жоқ ескі JSON-ды оқып, тапсырма санын көрсетеді. Ол нұсқаны, оң әрі қайталанбайтын ID-ді және бос емес атауды тексереді. Ескі JSON-ды қолмен өзгертпеңіз: ол үйлесімділікті тексеруге қажет.

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

fn decode(text: &str) -> Result<Organizer, String> {
    let data: Organizer = serde_json::from_str(text).map_err(|e| e.to_string())?;
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

fn main() -> Result<(), String> {
    let old_json =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Кітап сатып алу","done":false}]}"#;
    let data = decode(old_json)?;
    println!("Тапсырма саны: {}", data.tasks.len());
    Ok(())
}

#[test]
fn reads_old_task() {
    let text =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Кітап сатып алу","done":false}]}"#;
    assert_eq!(decode(text).unwrap().tasks.len(), 1);
}
```
```text
Тапсырма саны: 1
```

## Бетке қарамай жауап беріңіз

1. Ескі файлға анықталған әдепкі басымдық не үшін керек?
2. `default` анық жазылған `priority: 0` мәнін түзеуі керек пе?
3. Шұғылдарды санағанда аяқталған тапсырма неге алынбайды?

## Тапсырма

**Міндетті.** Кодтан бұрын төрт ереже жазыңыз: тек 1, 2, 3 жарамды; ескі JSON-да өріс болмаса 2 алынады; анық жазылған 0 не 4 — қате; іріктеуге тек аяқталмаған 3 басымдықтағы тапсырмалар кіреді. Содан кейін өріс пен санау қызметін, ескі және жарамсыз JSON үшін екі сынақ қосыңыз. Ескі тапсырма басымдығын және 3 мәні бар жаңа аяқталмаған тапсырма қосылған соң шұғылдар санын шығарыңыз. Екі жолды алдын ала болжаңыз.

## Жауаптар

Құрылымның алдына `fn ordinary_priority() -> u8 { 2 }` жазыңыз. Жаңа өрістің алдына `#[serde(default = "ordinary_priority")]` қойыңыз. Тексеруде `!(1..=3).contains(&task.priority)`, іріктеуде `iter().filter(...).count()` қолданыңыз.

<!-- task-answer -->
```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;

fn ordinary_priority() -> u8 {
    2
}

#[derive(Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
struct Task {
    id: u64,
    title: String,
    done: bool,
    #[serde(default = "ordinary_priority")]
    priority: u8,
}

#[derive(Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
struct Organizer {
    version: u32,
    next_id: u64,
    tasks: Vec<Task>,
}

fn decode(text: &str) -> Result<Organizer, String> {
    let data: Organizer = serde_json::from_str(text).map_err(|e| e.to_string())?;
    if data.version != 1 || data.next_id == 0 {
        return Err(String::from("invalid header"));
    }
    let mut ids = HashSet::new();
    for task in &data.tasks {
        if task.id == 0
            || task.id >= data.next_id
            || task.title.trim().is_empty()
            || !(1..=3).contains(&task.priority)
            || !ids.insert(task.id)
        {
            return Err(String::from("invalid task"));
        }
    }
    Ok(data)
}

fn high_open(data: &Organizer) -> usize {
    data.tasks
        .iter()
        .filter(|task| task.priority == 3 && !task.done)
        .count()
}

fn main() -> Result<(), String> {
    let old_json =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Кітап сатып алу","done":false}]}"#;
    let mut data = decode(old_json)?;
    println!("Ескі тапсырманың басымдығы: {}", data.tasks[0].priority);
    data.tasks.push(Task {
        id: 2,
        title: String::from("Досқа қоңырау шалу"),
        done: false,
        priority: 3,
    });
    data.next_id = 3;
    println!("Шұғыл аяқталмағандар: {}", high_open(&data));
    Ok(())
}

#[test]
fn old_file_gets_ordinary_priority() {
    let text =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Кітап сатып алу","done":false}]}"#;
    assert_eq!(decode(text).unwrap().tasks[0].priority, 2);
}

#[test]
fn invalid_priority_is_rejected() {
    let text = r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Кітап сатып алу","done":false,"priority":0}]}"#;
    assert!(decode(text).is_err());
}
```
```text
Ескі тапсырманың басымдығы: 2
Шұғыл аяқталмағандар: 1
```

## Тексергеннен кейін

Ескі файл оқылуы тиіс, бірақ бүлінген не анық жарамсыз файлды үнсіз бос тізімге айналдыруға болмайды. Осы тәсілмен өз кеңейтуіңізді де ұсына аласыз; алдымен ескі дерекке арналған ережені және қате сынағын жазыңыз. [serde default құжаттамасы](https://serde.rs/field-attrs.html#default).

[Алдыңғы сабақ](/read/rust-58-distribution?lang=kz) · [Мазмұн](/course/rust?lang=kz)
