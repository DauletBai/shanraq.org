# Сақтық көшірме

_Қысқаша (summary):_ **Сабақ 53. Сақтық көшірмені тексеріп, деректі бөлек файлға қалпына келтіру.**

[JSON](/read/rust-51-json?lang=kz) мен [сақтауды](/read/rust-52-safe-save?lang=kz) еске түсіріңіз. Мұнда уақытша оқу қалтасы қолданылады.

## Таныс бейне және тірек сызба

Қосалқы кілт құлыпқа сай келсе ғана пайдалы. Сол сияқты сақтық көшірме оқылып, тексерілсе ғана көмектеседі. **Тірек сызба:** негізгі файл → тексерілген көшірме → негізгі файлдың бүлінуі → көшірмені тексеру → бөлек қалпына келген файл. Бүлінген бастапқы файл талдау үшін қалады. Бұл ұқсастық бүкіл қалта жойылса, көшірме аман қалады дегенді білдірмейді.

`load` JSON-ды `Organizer` мәніне айналдырып, нұсқасын, нөмірлерін, атауларын тексереді. `io::ErrorKind::InvalidData` пішімі немесе мағынасы жарамсыз деректі білдіреді. `copy_checked` алдымен бастапқы деректі тексереді, содан кейін `create_new(true)` арқылы жаңа нысана файл жасайды, жазады, құрылғыға жеткізуді сұрайды және нәтижені қайта оқиды. `create_new` басқа файлды үстінен жаздырмайды. `read_to_string` дұрыс UTF-8 мәтінін талап етеді. `primary`, `backup`, `recovered` — үш бөлек жол. `fs::write` қалпына келтіруді көрсету үшін тек оқу мысалындағы негізгі файлды бүлдіреді.

## Іске қосып, талдаңыз

51-сабақтағы Serde тәуелділіктерімен жоба ашыңыз. Алғашқы Rust үзіндісін `src/main.rs` файлына қойып, `cargo run` орындаңыз. Ретін бақылаңыз: көшірме бүлінуден бұрын жасалады; бүлінгеннен кейін көшірме оқылып, қалпына келген дерек жаңа файлға жазылады. Бастапқы дерек жарамсыз болса, нысана файл әлі ашылмайды. Жазу қатесі нысанада шала файл қалдыруы мүмкін; оны қайта тексермей жарамды көшірме деуге болмайды.

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
use std::fs::{self, OpenOptions};
use std::io::{self, Write};
use std::path::Path;

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

fn load(text: &str) -> io::Result<Organizer> {
    let data: Organizer = serde_json::from_str(text)
        .map_err(|error| io::Error::new(io::ErrorKind::InvalidData, error))?;
    if data.version != 1 || data.next_id == 0 {
        return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid header"));
    }
    let mut ids = HashSet::new();
    for task in &data.tasks {
        if task.id == 0
            || task.id >= data.next_id
            || task.title.trim().is_empty()
            || !ids.insert(task.id)
        {
            return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid task"));
        }
    }
    Ok(data)
}

fn copy_checked(source: &Path, target: &Path) -> io::Result<()> {
    let text = fs::read_to_string(source)?;
    load(&text)?;
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(target)?;
    file.write_all(text.as_bytes())?;
    file.sync_all()?;
    drop(file);
    load(&fs::read_to_string(target)?)?;
    Ok(())
}

fn run() -> io::Result<()> {
    let folder = std::env::temp_dir().join(format!("rust-lesson-53-{}", std::process::id()));
    fs::create_dir(&folder)?;
    let primary = folder.join("tasks.json");
    let backup = folder.join("backup.json");
    let recovered = folder.join("recovered.json");
    let data = Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from("Кітап сатып алу"),
            done: false,
        }],
    };
    let json = serde_json::to_string(&data)
        .map_err(|error| io::Error::new(io::ErrorKind::InvalidData, error))?;
    fs::write(&primary, json)?;
    copy_checked(&primary, &backup)?;
    fs::write(&primary, "{broken")?;
    println!(
        "Негізгі файл бүлінген: {}",
        load(&fs::read_to_string(&primary)?).is_err()
    );
    copy_checked(&backup, &recovered)?;
    let restored = load(&fs::read_to_string(&recovered)?)?;
    println!("Қалпына келді: {}", restored.tasks[0].title);
    println!(
        "Бастапқы файл өзгермеді: {}",
        fs::read_to_string(&primary)? == "{broken"
    );
    fs::remove_dir_all(folder)?;
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
Негізгі файл бүлінген: true
Қалпына келді: Кітап сатып алу
Бастапқы файл өзгермеді: true
```

## Бетке қарамай жауап беріңіз

1. Бастапқы дерек неге алдымен тексеріледі?
2. Қалпына келген дерек неге үшінші файлға жазылады?
3. Қайта оқу сәтсіз болса, не істеу керек?

## Тапсырма

**Міндетті.** `{broken` деген жарамсыз JSON бар бөлек файл жасап, `copy_checked` қызметін жаңа жолға қолданыңыз. Қате анықталғанын және қалпына келетін файл жасалмағанын шығарыңыз. Дұрыс көшірмені өшірмеңіз.

## Жауаптар

`copy_checked(&bad_backup, &bad_target)` шақырыңыз. Нәтижені `is_err()`, нысана жоқтығын `!bad_target.exists()` арқылы тексеріңіз.

<!-- task-answer -->
```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::fs::{self, OpenOptions};
use std::io::{self, Write};
use std::path::Path;

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

fn load(text: &str) -> io::Result<Organizer> {
    let data: Organizer = serde_json::from_str(text)
        .map_err(|error| io::Error::new(io::ErrorKind::InvalidData, error))?;
    if data.version != 1 || data.next_id == 0 {
        return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid header"));
    }
    let mut ids = HashSet::new();
    for task in &data.tasks {
        if task.id == 0
            || task.id >= data.next_id
            || task.title.trim().is_empty()
            || !ids.insert(task.id)
        {
            return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid task"));
        }
    }
    Ok(data)
}

fn copy_checked(source: &Path, target: &Path) -> io::Result<()> {
    let text = fs::read_to_string(source)?;
    load(&text)?;
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(target)?;
    file.write_all(text.as_bytes())?;
    file.sync_all()?;
    drop(file);
    load(&fs::read_to_string(target)?)?;
    Ok(())
}

fn run() -> io::Result<()> {
    let folder = std::env::temp_dir().join(format!("rust-lesson-53-{}", std::process::id()));
    fs::create_dir(&folder)?;
    let primary = folder.join("tasks.json");
    let backup = folder.join("backup.json");
    let recovered = folder.join("recovered.json");
    let data = Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from("Кітап сатып алу"),
            done: false,
        }],
    };
    let json = serde_json::to_string(&data)
        .map_err(|error| io::Error::new(io::ErrorKind::InvalidData, error))?;
    fs::write(&primary, json)?;
    copy_checked(&primary, &backup)?;
    fs::write(&primary, "{broken")?;
    println!(
        "Негізгі файл бүлінген: {}",
        load(&fs::read_to_string(&primary)?).is_err()
    );
    copy_checked(&backup, &recovered)?;
    let restored = load(&fs::read_to_string(&recovered)?)?;
    println!("Қалпына келді: {}", restored.tasks[0].title);
    println!(
        "Бастапқы файл өзгермеді: {}",
        fs::read_to_string(&primary)? == "{broken"
    );
    let damaged_backup = folder.join("damaged-backup.json");
    let other = folder.join("other-recovery.json");
    fs::write(&damaged_backup, "{broken")?;
    let rejected = copy_checked(&damaged_backup, &other).is_err() && !other.exists();
    println!("Жарамсыз көшірме қабылданбады: {}", rejected);
    fs::remove_dir_all(folder)?;
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
Негізгі файл бүлінген: true
Қалпына келді: Кітап сатып алу
Бастапқы файл өзгермеді: true
Жарамсыз көшірме қабылданбады: true
```

## Тексергеннен кейін

Тексерілген көшірме сақтау жоспарының бәрін алмастырмайды: ол ескіруі немесе негізгі файлмен бірге жоғалуы мүмкін. Қалпына келтірерде табылған нұсқаны адамға көрсетіңіз. [51-сабаққа оралу](/read/rust-51-json?lang=kz). [create_new құжаттамасы](https://doc.rust-lang.org/std/fs/struct.OpenOptions.html#method.create_new).

[Алдыңғы сабақ](/read/rust-52-safe-save?lang=kz) · [Мазмұн](/course/rust?lang=kz)
