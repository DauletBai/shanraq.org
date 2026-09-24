# Консольдық органайзерді шығару

_Қысқаша (summary):_ **Сабақ 60. Таныс бөліктерді басқа адам іске қосып, тексере алатын бағдарламаға біріктіру.**

[JSON](/read/rust-51-json?lang=kz), [сақтау](/read/rust-52-safe-save?lang=kz), [көшірме](/read/rust-53-backup?lang=kz), [пәрмендер](/read/rust-54-cli?lang=kz), [сынақтар](/read/rust-56-testing?lang=kz) және [басымдық](/read/rust-59-capstone?lang=kz) керек. Бұл — қорытынды консольдық жоба: әр әрекет жеке іске қосылады, ал тапсырмалар JSON файлында қалады.

## Таныс бейне және тірек сызба

Шеберхананы елестетіңіз: тапсырма қағаздары құлыпты жәшікте жатыр, көрші қағаз алынса да, өз нөмірі өзгермейді; қосалқы жәшіктің кілті алдымен тексеріледі. **Тірек сызба:** аргументтер мен жол → файлды оқу → тексеру → әрекет → уақытша файл арқылы жазу → түсінікті нәтиже не қате. Көшірмеге бөлек тармақ: бастапқыны тексеру → жаңа файл жасау → қайта оқу. Сызбаны мәтінге қарамай айтыңыз. Ұқсастықтың шегі: екі адамның бір жәшіктегі қағазды қатар өзгертуі бөлек үйлестіруді қажет етеді; бұл бағдарлама оны атқармайды.

`--data ЖОЛ` файлды анық көрсетеді; әйтпесе бағдарлама ортадағы `ORGANIZER_DATA` мәнін, содан кейін ағымдағы қалтадағы `tasks.json` файлын алады. `std::env::args()` аргументтерді Unicode мәтіні ретінде алады; UTF-8 болмайтын байттары бар жолдар мұнда қолданылмайды. Файлдың болмауы тек алғашқы іске қосуда рұқсат: онда жадта бос күй жасалады. Қолжетімсіз не бүлінген JSON қате қайтарады, бос тізімге айналмайды. `Result<String, (i32, String)>` қалыпты жауапты немесе «аяқталу коды, қате мәтіні» жұбын сақтайды. 2 коды жарамсыз пәрменді, 1 коды файл не дерек қатесін білдіреді; `eprintln!` stderr-ге жазады. `retain` тізімде шартқа сай элементтерді қалдырады. Басымдық бойынша реттегенде `then` басымдық тең болса, ID-ді салыстырады. `to_owned()` мәтін көрінісінен өз жолын жасайды.

## Қорытынды бағдарлама

Осы сабақтың дайын жоба қалтасын алыңыз немесе төмендегі тәуелділіктері бар бөлек Cargo жобасын жасаңыз. Алғашқы Rust үзіндісін түгелдей `src/main.rs` файлына қойыңыз; оны өзіңіздің шын дерек файлыңызға бөлшектеп көшірмеңіз. Аргументсіз `cargo run` көмек көрсетеді. `cargo run -- --data tasks.json add 3 "Кітап сатып алу"` пәрменіндегі алғашқы `--` Cargo аргументтерін аяқтайды; кейінгі `--data` оқу файлын таңдайды; `3` — шұғыл басымдық; тырнақша арасы бос атауды бір аргумент ретінде сақтайды. Бұл Cargo пәрмені жоба қалтасынан Windows, macOS және Linux жүйелерінде орындалады. Қызметтік `add`, `list`, `done`, `delete`, `rename`, `search`, `filter`, `sort`, `stats`, `backup`, `restore`, `help` сөздері өзгермейді; адамға арналған хабарлар қазақша.

[Оқу жобасының қалтасын](https://github.com/DauletBai/shanraq.org/tree/main/course/rust-organizer/kz/step-60) қорда сақталғаннан кейін алуға болады; ішінде `Cargo.toml`, `Cargo.lock`, `README.md`, `LICENSE` және `src/main.rs` бар.

Көмек жолы бірнеше таңдауды көрсетеді: `|` таңбасы пәрмендерді бөледі, бәрін бірден термейсіз. Тапсырма жолындағы `[3]` басымдықтың 3 екенін білдіреді; тік жақшалар мұнда экранға шығатын мәтіннің бөлігі. Іске қоспас бұрын кодтан бес бөлікті табыңыз: дерек үлгісі мен тексеру, оқу, сақтықпен жазу мен көшіру, пәрменді орындау, `main` ішіндегі жолды таңдау.

```toml
[dependencies]
serde = { version = "=1.0.228", features = ["derive"] }
serde_json = "=1.0.149"
```

```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::fs::{self, OpenOptions};
use std::io::Write;
use std::path::{Path, PathBuf};

fn ordinary_priority() -> u8 {
    2
}

#[derive(Clone, Serialize, Deserialize)]
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

impl Organizer {
    fn new() -> Self {
        Self {
            version: 1,
            next_id: 1,
            tasks: Vec::new(),
        }
    }
    fn validate(&self) -> Result<(), String> {
        if self.version != 1 || self.next_id == 0 {
            return Err(String::from("Жарамсыз дерек"));
        }
        let mut ids = HashSet::new();
        for task in &self.tasks {
            if task.id == 0
                || task.id >= self.next_id
                || task.title.trim().is_empty()
                || !(1..=3).contains(&task.priority)
                || !ids.insert(task.id)
            {
                return Err(String::from("Жарамсыз дерек"));
            }
        }
        Ok(())
    }
    fn add(&mut self, priority: u8, title: String) -> Result<u64, String> {
        if !(1..=3).contains(&priority) {
            return Err(String::from("Басымдық 1 мен 3 арасында болуы керек"));
        }
        if title.trim().is_empty() {
            return Err(String::from("Атау бос"));
        }
        let next = self.next_id.checked_add(1).ok_or("Нөмірлер таусылды")?;
        let id = self.next_id;
        self.tasks.push(Task {
            id,
            title: title.trim().to_owned(),
            done: false,
            priority,
        });
        self.next_id = next;
        Ok(id)
    }
    fn task_mut(&mut self, id: u64) -> Result<&mut Task, String> {
        self.tasks
            .iter_mut()
            .find(|task| task.id == id)
            .ok_or(String::from("ID табылмады"))
    }
    fn delete(&mut self, id: u64) -> Result<(), String> {
        let before = self.tasks.len();
        self.tasks.retain(|task| task.id != id);
        if self.tasks.len() == before {
            Err(String::from("ID табылмады"))
        } else {
            Ok(())
        }
    }
}

fn decode(text: &str) -> Result<Organizer, String> {
    let data: Organizer = serde_json::from_str(text).map_err(|e| e.to_string())?;
    data.validate()?;
    Ok(data)
}

fn load_or_new(path: &Path) -> Result<Organizer, String> {
    match fs::read_to_string(path) {
        Ok(text) => decode(&text),
        Err(error) if error.kind() == std::io::ErrorKind::NotFound => Ok(Organizer::new()),
        Err(error) => Err(error.to_string()),
    }
}

fn save_safe(path: &Path, data: &Organizer) -> Result<(), String> {
    data.validate()?;
    let text = serde_json::to_string(data).map_err(|e| e.to_string())?;
    let temporary = path.with_extension(format!("tmp-{}", std::process::id()));
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(&temporary)
        .map_err(|e| e.to_string())?;
    if let Err(error) = file
        .write_all(text.as_bytes())
        .and_then(|_| file.sync_all())
    {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error.to_string());
    }
    drop(file);
    if let Err(error) = fs::rename(&temporary, path) {
        let _ = fs::remove_file(&temporary);
        return Err(error.to_string());
    }
    Ok(())
}

fn copy_checked(source: &Path, target: &Path) -> Result<(), String> {
    let text = fs::read_to_string(source).map_err(|e| e.to_string())?;
    decode(&text)?;
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(target)
        .map_err(|e| e.to_string())?;
    if let Err(error) = file
        .write_all(text.as_bytes())
        .and_then(|_| file.sync_all())
    {
        drop(file);
        let _ = fs::remove_file(target);
        return Err(error.to_string());
    }
    drop(file);
    if let Err(error) = fs::read_to_string(target)
        .map_err(|e| e.to_string())
        .and_then(|copied| decode(&copied).map(|_| ()))
    {
        let _ = fs::remove_file(target);
        return Err(error);
    }
    Ok(())
}

fn parse_id(text: &str) -> Result<u64, String> {
    match text.parse::<u64>() {
        Ok(id) if id > 0 => Ok(id),
        _ => Err(String::from("Оң ID керек")),
    }
}

fn parse_priority(text: &str) -> Result<u8, String> {
    match text.parse::<u8>() {
        Ok(priority) if (1..=3).contains(&priority) => Ok(priority),
        _ => Err(String::from("Басымдық 1 мен 3 арасында болуы керек")),
    }
}

fn show(tasks: &[&Task]) -> String {
    if tasks.is_empty() {
        return String::from("Тапсырма жоқ");
    }
    tasks
        .iter()
        .map(|task| {
            let status = if task.done {
                "аяқталды"
            } else {
                "күтіп тұр"
            };
            format!("{} [{}] {} ({status})", task.id, task.priority, task.title)
        })
        .collect::<Vec<_>>()
        .join("\n")
}

fn execute(args: &[String], path: &Path) -> Result<String, (i32, String)> {
    let usage = || (2, String::from("Белгісіз пәрмен не жарамсыз аргументтер"));
    let data_error = |message: String| (1, message);
    if args.is_empty() || args == ["help"] {
        return Ok(String::from(
            "Пәрмендер: add БАСЫМДЫҚ АТАУ | list | done ID | delete ID | rename ID АТАУ | search МӘТІН | filter open/done/high | sort id/title/priority | stats | backup ЖОЛ | restore КӨШІРМЕ ЖАҢА_ЖОЛ | --version",
        ));
    }
    if args == ["--version"] {
        return Ok(env!("CARGO_PKG_VERSION").to_owned());
    }
    match args[0].as_str() {
        "backup" if args.len() == 2 => {
            copy_checked(path, Path::new(&args[1])).map_err(data_error)?;
            Ok(String::from("Көшірме жасалды"))
        }
        "restore" if args.len() == 3 => {
            copy_checked(Path::new(&args[1]), Path::new(&args[2])).map_err(data_error)?;
            Ok(String::from("Бөлек қалпына келді"))
        }
        "add" if args.len() >= 3 => {
            let priority = parse_priority(&args[1]).map_err(|e| (2, e))?;
            let mut data = load_or_new(path).map_err(data_error)?;
            let id = data
                .add(priority, args[2..].join(" "))
                .map_err(|e| (2, e))?;
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Қосылды: {id}"))
        }
        "list" if args.len() == 1 => {
            let data = load_or_new(path).map_err(data_error)?;
            let mut tasks: Vec<&Task> = data.tasks.iter().collect();
            tasks.sort_by_key(|task| task.id);
            Ok(show(&tasks))
        }
        "done" if args.len() == 2 => {
            let id = parse_id(&args[1]).map_err(|e| (2, e))?;
            let mut data = load_or_new(path).map_err(data_error)?;
            data.task_mut(id).map_err(data_error)?.done = true;
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Аяқталды: {id}"))
        }
        "delete" if args.len() == 2 => {
            let id = parse_id(&args[1]).map_err(|e| (2, e))?;
            let mut data = load_or_new(path).map_err(data_error)?;
            data.delete(id).map_err(data_error)?;
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Өшірілді: {id}"))
        }
        "rename" if args.len() >= 3 => {
            let id = parse_id(&args[1]).map_err(|e| (2, e))?;
            let title = args[2..].join(" ");
            if title.trim().is_empty() {
                return Err((2, String::from("Атау бос")));
            }
            let mut data = load_or_new(path).map_err(data_error)?;
            data.task_mut(id).map_err(data_error)?.title = title.trim().to_owned();
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Атауы өзгерді: {id}"))
        }
        "search" if args.len() >= 2 => {
            let term = args[1..].join(" ");
            if term.trim().is_empty() {
                return Err(usage());
            }
            let data = load_or_new(path).map_err(data_error)?;
            let tasks: Vec<&Task> = data
                .tasks
                .iter()
                .filter(|task| task.title.contains(&term))
                .collect();
            Ok(show(&tasks))
        }
        "filter" if args.len() == 2 => {
            let data = load_or_new(path).map_err(data_error)?;
            let tasks: Vec<&Task> = match args[1].as_str() {
                "open" => data.tasks.iter().filter(|task| !task.done).collect(),
                "done" => data.tasks.iter().filter(|task| task.done).collect(),
                "high" => data
                    .tasks
                    .iter()
                    .filter(|task| task.priority == 3 && !task.done)
                    .collect(),
                _ => return Err(usage()),
            };
            Ok(show(&tasks))
        }
        "sort" if args.len() == 2 => {
            let data = load_or_new(path).map_err(data_error)?;
            let mut tasks: Vec<&Task> = data.tasks.iter().collect();
            match args[1].as_str() {
                "id" => tasks.sort_by_key(|task| task.id),
                "title" => tasks.sort_by(|a, b| a.title.cmp(&b.title)),
                "priority" => {
                    tasks.sort_by(|a, b| b.priority.cmp(&a.priority).then(a.id.cmp(&b.id)))
                }
                _ => return Err(usage()),
            }
            Ok(show(&tasks))
        }
        "stats" if args.len() == 1 => {
            let data = load_or_new(path).map_err(data_error)?;
            let total = data.tasks.len();
            let done = data.tasks.iter().filter(|task| task.done).count();
            let open = total - done;
            let high = data
                .tasks
                .iter()
                .filter(|task| task.priority == 3 && !task.done)
                .count();
            Ok(format!(
                "Барлығы: {total}, аяқталды: {done}, қалды: {open}, шұғыл ашық: {high}"
            ))
        }
        _ => Err(usage()),
    }
}

fn main() {
    let mut args: Vec<String> = std::env::args().skip(1).collect();
    let path = if args.first().map(String::as_str) == Some("--data") {
        if args.len() < 2 {
            eprintln!("Қате: Белгісіз пәрмен не жарамсыз аргументтер");
            std::process::exit(2);
        }
        args.remove(0);
        PathBuf::from(args.remove(0))
    } else {
        std::env::var_os("ORGANIZER_DATA")
            .map(PathBuf::from)
            .unwrap_or_else(|| PathBuf::from("tasks.json"))
    };
    match execute(&args, &path) {
        Ok(message) => println!("{message}"),
        Err((code, message)) => {
            eprintln!("Қате: {message}");
            std::process::exit(code);
        }
    }
}

#[test]
fn damaged_json_is_rejected() {
    assert!(decode("{broken").is_err());
}

#[test]
fn unknown_version_is_rejected() {
    let text = r#"{"version":2,"next_id":1,"tasks":[]}"#;
    assert!(decode(text).is_err());
}

#[test]
fn deleting_one_task_keeps_the_next_id() {
    let mut data = Organizer::new();
    let first = data.add(2, String::from("Кітап сатып алу")).unwrap();
    let second = data.add(3, String::from("Кітап сатып алу")).unwrap();
    data.delete(first).unwrap();
    assert_eq!(data.tasks[0].id, second);
    assert_eq!(data.next_id, second + 1);
}
```
```text
Пәрмендер: add БАСЫМДЫҚ АТАУ | list | done ID | delete ID | rename ID АТАУ | search МӘТІН | filter open/done/high | sort id/title/priority | stats | backup ЖОЛ | restore КӨШІРМЕ ЖАҢА_ЖОЛ | --version
```

| Әрекет | Оқу файлының күйі | Көрінетін нәтиже |
|---|---|---|
| `help` | Файл жасалмайды | Пәрмендер тізімі |
| `add 3 "Кітап сатып алу"` | 1-тапсырма, келесі ID 2 | `Қосылды: 1` |
| Жаңа іске қосудағы `list` | Сол файл оқылады | `1 [3] Кітап сатып алу (күтіп тұр)` |
| `done 1` | 1-тапсырма аяқталған | `Аяқталды: 1` |
| `restore backup.json recovered.json` | Негізгі файл өзгермейді | Жаңа `recovered.json` жасалады |



`list` ID-ді, тік жақшадағы басымдықты, атауды және күйді көрсетеді. `search` әріптің үлкен-кішілігін ескере отырып дәл үзіндіні іздейді; `sort title` жолдарды Rust ережесімен салыстырады, бұл қазақша сөздік рет емес. `filter high` тек басымдығы 3 болатын аяқталмаған тапсырмаларды көрсетеді. `backup ЖОЛ` тексерілген көшірме жасайды; `restore КӨШІРМЕ ЖАҢА_ЖОЛ` оны тек жаңа файлға қалпына келтіреді. Содан кейін бүлінген бастапқы файлды өшірмей, жаңа файлды `--data` арқылы оқи аласыз. `save_safe` көршілес уақытша файл арқылы сақтайды: бұл алмастыруға дейінгі кәдімгі жазу қатесінен қорғайды, бірақ қуат үзілуі не қатар жазу кезінде толық кепіл бермейді. Төмендегі бірнеше пәрмен бөлек оқу қалтасында орындалады. Файл бұрыннан бар болса, басқа оқу атауын таңдаңыз.

## Бетке қарамай жауап беріңіз

1. Бүлінген файл неге алғашқы іске қосу деп есептелмейді?
2. `restore` неге жаңа жолды талап етеді?
3. `search` пен `sort title` қазақша сөздік іздеу және реттеуден несімен бөлек?
4. `rename` алдындағы жазу қатесінде не болады?

## Тапсырма

**Міндетті.** Жаңа оқу қалтасында `cargo test` орындаңыз; одан кейін басымдығы 3, атауында бос орны бар тапсырма қосыңыз, бөлек іске қосып тізімді көрсетіңіз, тапсырманы аяқтап, тізімді тағы оқыңыз. Көшірме жасаңыз, тек оқу негізгі JSON-ын бүлдіріп, `list` қатемен аяқталатынын және файл өшірілмегенін көріңіз; көшірмені жаңа жолға қалпына келтіріп, оны `--data` арқылы оқыңыз. Белгісіз пәрмен мен жарамсыз ID-ді тексеріңіз. `priority` өрісі жоқ ескі JSON, қайталанған ID және жарамсыз басымдық үшін сынақтар қосыңыз. README-ді досыңызға беріп, автор түсіндірмесінсіз қадамдарды қайталауын сұраңыз.

## Жауаптар

Әр `cargo run` кезінде `-- --data ЖОЛ` беріңіз: жеке іске қосулар бір оқу файлын қолданады. Бүлдіруге тек алдын ала жасалған оқу файлын алыңыз. Сынақта үйдегі шын файлдың орнына `decode` шақырыңыз. Пәрмен үлгісі жауап кодынан кейін тұр.

<!-- task-answer -->
```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::fs::{self, OpenOptions};
use std::io::Write;
use std::path::{Path, PathBuf};

fn ordinary_priority() -> u8 {
    2
}

#[derive(Clone, Serialize, Deserialize)]
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

impl Organizer {
    fn new() -> Self {
        Self {
            version: 1,
            next_id: 1,
            tasks: Vec::new(),
        }
    }
    fn validate(&self) -> Result<(), String> {
        if self.version != 1 || self.next_id == 0 {
            return Err(String::from("Жарамсыз дерек"));
        }
        let mut ids = HashSet::new();
        for task in &self.tasks {
            if task.id == 0
                || task.id >= self.next_id
                || task.title.trim().is_empty()
                || !(1..=3).contains(&task.priority)
                || !ids.insert(task.id)
            {
                return Err(String::from("Жарамсыз дерек"));
            }
        }
        Ok(())
    }
    fn add(&mut self, priority: u8, title: String) -> Result<u64, String> {
        if !(1..=3).contains(&priority) {
            return Err(String::from("Басымдық 1 мен 3 арасында болуы керек"));
        }
        if title.trim().is_empty() {
            return Err(String::from("Атау бос"));
        }
        let next = self.next_id.checked_add(1).ok_or("Нөмірлер таусылды")?;
        let id = self.next_id;
        self.tasks.push(Task {
            id,
            title: title.trim().to_owned(),
            done: false,
            priority,
        });
        self.next_id = next;
        Ok(id)
    }
    fn task_mut(&mut self, id: u64) -> Result<&mut Task, String> {
        self.tasks
            .iter_mut()
            .find(|task| task.id == id)
            .ok_or(String::from("ID табылмады"))
    }
    fn delete(&mut self, id: u64) -> Result<(), String> {
        let before = self.tasks.len();
        self.tasks.retain(|task| task.id != id);
        if self.tasks.len() == before {
            Err(String::from("ID табылмады"))
        } else {
            Ok(())
        }
    }
}

fn decode(text: &str) -> Result<Organizer, String> {
    let data: Organizer = serde_json::from_str(text).map_err(|e| e.to_string())?;
    data.validate()?;
    Ok(data)
}

fn load_or_new(path: &Path) -> Result<Organizer, String> {
    match fs::read_to_string(path) {
        Ok(text) => decode(&text),
        Err(error) if error.kind() == std::io::ErrorKind::NotFound => Ok(Organizer::new()),
        Err(error) => Err(error.to_string()),
    }
}

fn save_safe(path: &Path, data: &Organizer) -> Result<(), String> {
    data.validate()?;
    let text = serde_json::to_string(data).map_err(|e| e.to_string())?;
    let temporary = path.with_extension(format!("tmp-{}", std::process::id()));
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(&temporary)
        .map_err(|e| e.to_string())?;
    if let Err(error) = file
        .write_all(text.as_bytes())
        .and_then(|_| file.sync_all())
    {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error.to_string());
    }
    drop(file);
    if let Err(error) = fs::rename(&temporary, path) {
        let _ = fs::remove_file(&temporary);
        return Err(error.to_string());
    }
    Ok(())
}

fn copy_checked(source: &Path, target: &Path) -> Result<(), String> {
    let text = fs::read_to_string(source).map_err(|e| e.to_string())?;
    decode(&text)?;
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(target)
        .map_err(|e| e.to_string())?;
    if let Err(error) = file
        .write_all(text.as_bytes())
        .and_then(|_| file.sync_all())
    {
        drop(file);
        let _ = fs::remove_file(target);
        return Err(error.to_string());
    }
    drop(file);
    if let Err(error) = fs::read_to_string(target)
        .map_err(|e| e.to_string())
        .and_then(|copied| decode(&copied).map(|_| ()))
    {
        let _ = fs::remove_file(target);
        return Err(error);
    }
    Ok(())
}

fn parse_id(text: &str) -> Result<u64, String> {
    match text.parse::<u64>() {
        Ok(id) if id > 0 => Ok(id),
        _ => Err(String::from("Оң ID керек")),
    }
}

fn parse_priority(text: &str) -> Result<u8, String> {
    match text.parse::<u8>() {
        Ok(priority) if (1..=3).contains(&priority) => Ok(priority),
        _ => Err(String::from("Басымдық 1 мен 3 арасында болуы керек")),
    }
}

fn show(tasks: &[&Task]) -> String {
    if tasks.is_empty() {
        return String::from("Тапсырма жоқ");
    }
    tasks
        .iter()
        .map(|task| {
            let status = if task.done {
                "аяқталды"
            } else {
                "күтіп тұр"
            };
            format!("{} [{}] {} ({status})", task.id, task.priority, task.title)
        })
        .collect::<Vec<_>>()
        .join("\n")
}

fn execute(args: &[String], path: &Path) -> Result<String, (i32, String)> {
    let usage = || (2, String::from("Белгісіз пәрмен не жарамсыз аргументтер"));
    let data_error = |message: String| (1, message);
    if args.is_empty() || args == ["help"] {
        return Ok(String::from(
            "Пәрмендер: add БАСЫМДЫҚ АТАУ | list | done ID | delete ID | rename ID АТАУ | search МӘТІН | filter open/done/high | sort id/title/priority | stats | backup ЖОЛ | restore КӨШІРМЕ ЖАҢА_ЖОЛ | --version",
        ));
    }
    if args == ["--version"] {
        return Ok(env!("CARGO_PKG_VERSION").to_owned());
    }
    match args[0].as_str() {
        "backup" if args.len() == 2 => {
            copy_checked(path, Path::new(&args[1])).map_err(data_error)?;
            Ok(String::from("Көшірме жасалды"))
        }
        "restore" if args.len() == 3 => {
            copy_checked(Path::new(&args[1]), Path::new(&args[2])).map_err(data_error)?;
            Ok(String::from("Бөлек қалпына келді"))
        }
        "add" if args.len() >= 3 => {
            let priority = parse_priority(&args[1]).map_err(|e| (2, e))?;
            let mut data = load_or_new(path).map_err(data_error)?;
            let id = data
                .add(priority, args[2..].join(" "))
                .map_err(|e| (2, e))?;
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Қосылды: {id}"))
        }
        "list" if args.len() == 1 => {
            let data = load_or_new(path).map_err(data_error)?;
            let mut tasks: Vec<&Task> = data.tasks.iter().collect();
            tasks.sort_by_key(|task| task.id);
            Ok(show(&tasks))
        }
        "done" if args.len() == 2 => {
            let id = parse_id(&args[1]).map_err(|e| (2, e))?;
            let mut data = load_or_new(path).map_err(data_error)?;
            data.task_mut(id).map_err(data_error)?.done = true;
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Аяқталды: {id}"))
        }
        "delete" if args.len() == 2 => {
            let id = parse_id(&args[1]).map_err(|e| (2, e))?;
            let mut data = load_or_new(path).map_err(data_error)?;
            data.delete(id).map_err(data_error)?;
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Өшірілді: {id}"))
        }
        "rename" if args.len() >= 3 => {
            let id = parse_id(&args[1]).map_err(|e| (2, e))?;
            let title = args[2..].join(" ");
            if title.trim().is_empty() {
                return Err((2, String::from("Атау бос")));
            }
            let mut data = load_or_new(path).map_err(data_error)?;
            data.task_mut(id).map_err(data_error)?.title = title.trim().to_owned();
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Атауы өзгерді: {id}"))
        }
        "search" if args.len() >= 2 => {
            let term = args[1..].join(" ");
            if term.trim().is_empty() {
                return Err(usage());
            }
            let data = load_or_new(path).map_err(data_error)?;
            let tasks: Vec<&Task> = data
                .tasks
                .iter()
                .filter(|task| task.title.contains(&term))
                .collect();
            Ok(show(&tasks))
        }
        "filter" if args.len() == 2 => {
            let data = load_or_new(path).map_err(data_error)?;
            let tasks: Vec<&Task> = match args[1].as_str() {
                "open" => data.tasks.iter().filter(|task| !task.done).collect(),
                "done" => data.tasks.iter().filter(|task| task.done).collect(),
                "high" => data
                    .tasks
                    .iter()
                    .filter(|task| task.priority == 3 && !task.done)
                    .collect(),
                _ => return Err(usage()),
            };
            Ok(show(&tasks))
        }
        "sort" if args.len() == 2 => {
            let data = load_or_new(path).map_err(data_error)?;
            let mut tasks: Vec<&Task> = data.tasks.iter().collect();
            match args[1].as_str() {
                "id" => tasks.sort_by_key(|task| task.id),
                "title" => tasks.sort_by(|a, b| a.title.cmp(&b.title)),
                "priority" => {
                    tasks.sort_by(|a, b| b.priority.cmp(&a.priority).then(a.id.cmp(&b.id)))
                }
                _ => return Err(usage()),
            }
            Ok(show(&tasks))
        }
        "stats" if args.len() == 1 => {
            let data = load_or_new(path).map_err(data_error)?;
            let total = data.tasks.len();
            let done = data.tasks.iter().filter(|task| task.done).count();
            let open = total - done;
            let high = data
                .tasks
                .iter()
                .filter(|task| task.priority == 3 && !task.done)
                .count();
            Ok(format!(
                "Барлығы: {total}, аяқталды: {done}, қалды: {open}, шұғыл ашық: {high}"
            ))
        }
        _ => Err(usage()),
    }
}

fn main() {
    let mut args: Vec<String> = std::env::args().skip(1).collect();
    let path = if args.first().map(String::as_str) == Some("--data") {
        if args.len() < 2 {
            eprintln!("Қате: Белгісіз пәрмен не жарамсыз аргументтер");
            std::process::exit(2);
        }
        args.remove(0);
        PathBuf::from(args.remove(0))
    } else {
        std::env::var_os("ORGANIZER_DATA")
            .map(PathBuf::from)
            .unwrap_or_else(|| PathBuf::from("tasks.json"))
    };
    match execute(&args, &path) {
        Ok(message) => println!("{message}"),
        Err((code, message)) => {
            eprintln!("Қате: {message}");
            std::process::exit(code);
        }
    }
}

#[test]
fn damaged_json_is_rejected() {
    assert!(decode("{broken").is_err());
}

#[test]
fn unknown_version_is_rejected() {
    let text = r#"{"version":2,"next_id":1,"tasks":[]}"#;
    assert!(decode(text).is_err());
}

#[test]
fn deleting_one_task_keeps_the_next_id() {
    let mut data = Organizer::new();
    let first = data.add(2, String::from("Кітап сатып алу")).unwrap();
    let second = data.add(3, String::from("Кітап сатып алу")).unwrap();
    data.delete(first).unwrap();
    assert_eq!(data.tasks[0].id, second);
    assert_eq!(data.next_id, second + 1);
}

#[test]
fn old_json_gets_normal_priority() {
    let text =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Кітап сатып алу","done":false}]}"#;
    assert_eq!(decode(text).unwrap().tasks[0].priority, 2);
}

#[test]
fn duplicate_id_is_rejected() {
    let text = r#"{"version":1,"next_id":3,"tasks":[{"id":1,"title":"Кітап сатып алу","done":false,"priority":2},{"id":1,"title":"Досқа қоңырау шалу","done":false,"priority":3}]}"#;
    assert!(decode(text).is_err());
}

#[test]
fn invalid_priority_is_rejected() {
    let text =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Кітап сатып алу","done":false,"priority":0}]}"#;
    assert!(decode(text).is_err());
}
```
```text
Пәрмендер: add БАСЫМДЫҚ АТАУ | list | done ID | delete ID | rename ID АТАУ | search МӘТІН | filter open/done/high | sort id/title/priority | stats | backup ЖОЛ | restore КӨШІРМЕ ЖАҢА_ЖОЛ | --version
```

```text
cargo run -- --data tasks.json add 3 "Кітап сатып алу"
cargo run -- --data tasks.json list
cargo run -- --data tasks.json done 1
cargo run -- --data tasks.json backup backup.json
cargo run -- --data tasks.json restore backup.json recovered.json
cargo run -- --data recovered.json list
```

`README.md`:

```markdown
# Консольдық органайзер

Rust 1.97.0 мен Cargo-ны https://www.rust-lang.org/tools/install бетіндегі нұсқаумен орнатып, жаңа терминал ашыңыз да `cargo --version` арқылы тексеріңіз. Жоба қалтасында `cargo test`, содан соң `cargo run -- --data tasks.json add 3 "Кітап сатып алу"` және бөлек іске қосып `cargo run -- --data tasks.json list` орындаңыз. `cargo run -- help` басқа пәрмендерді көрсетеді. Дерек бағдарлама ішінде емес, таңдалған JSON файлында қалады. Тәуекелі бар әрекет алдында `backup backup.json` арқылы көшірме жасаңыз. `restore backup.json recovered.json` жаңа жолға жазады; бастапқы файл өшірілмейді. JSON бүлінсе, бағдарлама қатемен аяқталады: оны бос тізіммен алмастырмаңыз. `--data` болмаса, `ORGANIZER_DATA` алынады; ол да болмаса, ағымдағы қалтадағы `tasks.json` қолданылады. Басқа жүйеге бағдарламаны сол жүйеде бөлек құрастырыңыз.
```

## Тексергеннен кейін

Автоматтандыруға кіріспей, қате мәтінімен бірге аяқталу кодын да тексеріңіз. Күтпеген қуат үзілуінен және бір уақытта бірнеше жазушыдан қатаң қорғау үшін файлдық жүйе мен қалтаны құрылғыға жеткізуге қатысты қосымша ереже керек. Түймелері бар жеке терезені кейін осы тексерілген әрекеттердің үстіне құруға болады; курстың міндетті нәтижесі — консольдық бағдарлама. [Rust ресми кітабы](https://doc.rust-lang.org/book/).

[Алдыңғы сабақ](/read/rust-59-capstone?lang=kz) · [Мазмұн](/course/rust?lang=kz)
