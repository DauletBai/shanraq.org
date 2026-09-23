# Бақылау кезеңі: жадтағы жоспарлау құралы

_Қысқаша (summary):_ **40-сабақ. Бір іске қосуда бірнеше пәрмен орындап, жарамсыз енгізуді өңдейтін жоспарлау құралын қолдану.**

Алдын ала: 16, 21, 27–28 және 31–39-сабақтар. Енгізу соңы түсініксіз болса, [34-сабақты](/read/rust-34-input?lang=kz), ID шатасса, [38-сабақты](/read/rust-38-ids?lang=kz) қайталаңыз.

## Таныс бейне және тірек сызба

Қабылдау орнында адам бірнеше өтінішті кезекпен бере алады; қызметкер есік жабылғанша бұрынғы өтініштерді есте сақтайды. Бұл жоспарлау құралы да шығу пәрменіне не енгізудің аяқталуына дейін жолдарды қабылдап, тапсырмаларды бір іске қосудың жадында ұстайды. Ұқсатудың шегі: бағдарлама жабылғанда дерек әзірше жоғалады; файлға сақтауды кейін үйренеміз.

**Енгізу жолы → `Command` талдауы → `Organizer` әрекеті → жауап → келесі жол; `quit` не енгізу соңы → аяқтау.** Мұнда 34-сабақтағы `read_line`, 36-сабақтағы пәрмен санамалауы, 37-сабақтағы тапсырма әрекеттері, 38-сабақтағы тұрақты ID және 39-сабақтағы қате тексеруі бірігеді. Стандартты кітапханадағы `split_once(' ')` әдісі мәтінді **алғашқы бос орыннан** бөліп, `Some((сол бөлік, оң бөлік))` қайтарады; бос орын болмаса `None` болады. Мысалы, `add Rust оқу` әрекетті `add`, ал толық атауды `Rust оқу` деп бөледі. Мұнда тырнақша қажет емес: бүкіл жол әлдеқашан оқылған. Тырнақша жазылса, ол атаудың бір бөлігі болады; бөлгіш ретінде кәдімгі бос орын алынады, табуляция емес. `trim()` шеткі бос орындар мен жол соңын алып тастайды. `self.tasks.is_empty()` тізімде тапсырма бар-жоғын тексереді. `done` мен `delete` пәрмендеріне оң ID керек; `list` пен `quit` қосымша дерек алмайды. Әр қате себебін шығарып, келесі жолды күтеді. `use std::io` енгізу модуліне апаратын жолды қысқартады: `std::io::stdin()` орнына `io::stdin()` жазамыз. Бұл 34-сабақтағы сол енгізу нысаны; модуль құрылымын 44-сабақта талдаймыз. `!execute(...)` «жалғастырмау» дегенді білдіреді, себебі `execute` `quit` үшін `false` қайтарады.

## Бір іске қосуда бірнеше пәрмен

Ескі `src/main.rs` файлын сақтап, толық ауыстырыңыз. `Cargo.toml` жанында `cargo run` орындаңыз. Әр жолға бір пәрмен жазып, Enter басыңыз. Шығу үшін `quit` жазыңыз; енгізу арнасы жабылса да бағдарлама аяқталады. Жаңа файл не тәуелділік керек емес. Төменде енгізу арнасы бірден жабылғандағы нәтиже берілген.

```rust
use std::io;

enum Command {
    Add(String),
    List,
    Done(u32),
    Delete(u32),
    Quit,
}

fn parse_id(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(0) => Err(String::from("Оң ID қажет")),
        Ok(id) => Ok(id),
        Err(_) => Err(String::from("Оң ID қажет")),
    }
}

fn parse(line: &str) -> Result<Command, String> {
    let line = line.trim();
    if line == "list" {
        return Ok(Command::List);
    }
    if line == "quit" {
        return Ok(Command::Quit);
    }
    match line.split_once(' ') {
        Some(("add", title)) => {
            if title.trim().is_empty() {
                Err(String::from("Атау бос"))
            } else {
                Ok(Command::Add(String::from(title.trim())))
            }
        }
        Some(("done", id)) => Ok(Command::Done(parse_id(id)?)),
        Some(("delete", id)) => Ok(Command::Delete(parse_id(id)?)),
        _ => Err(String::from("Белгісіз пәрмен не жарамсыз аргументтер")),
    }
}

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
    fn add(&mut self, title: String) -> Result<u32, String> {
        if title.trim().is_empty() {
            return Err(String::from("Атау бос"));
        }
        match self.next_id.checked_add(1) {
            Some(next) => {
                let id = self.next_id;
                self.tasks.push(Task {
                    id,
                    title,
                    done: false,
                });
                self.next_id = next;
                Ok(id)
            }
            None => Err(String::from("Нөмірлер таусылды")),
        }
    }
    fn done(&mut self, id: u32) -> Result<(), String> {
        for task in &mut self.tasks {
            if task.id == id {
                task.done = true;
                return Ok(());
            }
        }
        Err(String::from("ID табылмады"))
    }
    fn delete(&mut self, id: u32) -> Result<(), String> {
        for index in 0..self.tasks.len() {
            if self.tasks[index].id == id {
                self.tasks.remove(index);
                return Ok(());
            }
        }
        Err(String::from("ID табылмады"))
    }
    fn list(&self) {
        if self.tasks.is_empty() {
            println!("Тапсырма жоқ");
        }
        for task in &self.tasks {
            let mark = if task.done {
                "дайын"
            } else {
                "аяқталмаған"
            };
            println!("{}: {} ({mark})", task.id, task.title);
        }
    }
}

fn execute(organizer: &mut Organizer, command: Command) -> bool {
    match command {
        Command::Add(title) => {
            match organizer.add(title) {
                Ok(id) => println!("Қосылды: {id}"),
                Err(message) => println!("Қате: {message}"),
            }
            true
        }
        Command::List => {
            organizer.list();
            true
        }
        Command::Done(id) => {
            match organizer.done(id) {
                Ok(()) => println!("Аяқталды: {id}"),
                Err(message) => println!("Қате: {message}"),
            }
            true
        }
        Command::Delete(id) => {
            match organizer.delete(id) {
                Ok(()) => println!("Өшірілді: {id}"),
                Err(message) => println!("Қате: {message}"),
            }
            true
        }
        Command::Quit => false,
    }
}

fn main() {
    let mut organizer = Organizer::new();
    let stdin = io::stdin();
    println!("Пәрмендер: add АТАУ | list | done ID | delete ID | quit");
    loop {
        let mut line = String::new();
        match stdin.read_line(&mut line) {
            Ok(0) => break,
            Ok(_) => match parse(&line) {
                Ok(command) => {
                    if !execute(&mut organizer, command) {
                        break;
                    }
                }
                Err(message) => println!("Қате: {message}"),
            },
            Err(_) => {
                println!("Қате: енгізу оқылмады");
                break;
            }
        }
    }
    println!("Жұмыс аяқталды");
}
```
```text
Пәрмендер: add АТАУ | list | done ID | delete ID | quit
Жұмыс аяқталды
```

Қолмен `add Rust оқу`, `add Сынақтарды қайталау`, `list`, `done 1`, `delete 2`, `list`, `done 99`, `quit` енгізіңіз. ID 1 «дайын» күйінде қалады, ал жоқ ID қате береді. Қайта іске қосқанда тізім бос болады: файлсыз кезеңнің осы шегі күтіледі.

Көрсетілген пәрмендер тізбегінің нәтижесі мынадай болады (енгізілген жолдар мұнда қайталанбайды):

```text
Пәрмендер: add АТАУ | list | done ID | delete ID | quit
Қосылды: 1
Қосылды: 2
1: Rust оқу (аяқталмаған)
2: Сынақтарды қайталау (аяқталмаған)
Аяқталды: 1
Өшірілді: 2
1: Rust оқу (дайын)
Қате: ID табылмады
Жұмыс аяқталды
```

## Түсінгеніңізді тексеріңіз

1. `add Rust оқу` атаудағы бос орынды неге сақтайды?
2. `done 99` кезінде тізім неге өзгермеуі керек?
3. `quit` пен қайта іске қосудан кейін тапсырмалар қала ма?

`split_once` тек алғашқы бос орыннан бөледі; жоқ ID жазбаны өзгертпейді; файлсыз дерек жойылады. Бетті жапқан соң тірек сызбаны қайта айтыңыз.

## Тапсырма

**Міндетті.** `rename ID АТАУ` пәрменін қосыңыз: тапсырманы тұрақты ID бойынша тауып, бос атауды қабылдамасын, ал жоқ ID берілсе тізімді өзгертпесін. `add Rust оқу`, `rename 1 Екі тарау оқу`, `list`, `rename 99 Қате`, `quit` тізбегін сынаңыз. Бағдарламаны қайта іске қосып, тізім бос екенін тексеріңіз. `Vec` орнын ID деп алмаңыз.

## Ишара

`Command` ішіне `Rename(u32, String)` қосыңыз. `rename ` кейінгі мәтінге `split_once(' ')` қайта қолданып, ID мен атауды жеке тексеріңіз. `rename` әдісінде `&mut self.tasks` арқылы өтіп, тексеруден соң ғана өрісті өзгертіңіз.

## Талпыныстан кейінгі үлгі жауап

<!-- task-answer -->
```rust
use std::io;

enum Command {
    Add(String),
    List,
    Done(u32),
    Delete(u32),
    Rename(u32, String),
    Quit,
}

fn parse_id(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(0) => Err(String::from("Оң ID қажет")),
        Ok(id) => Ok(id),
        Err(_) => Err(String::from("Оң ID қажет")),
    }
}

fn parse(line: &str) -> Result<Command, String> {
    let line = line.trim();
    if line == "list" {
        return Ok(Command::List);
    }
    if line == "quit" {
        return Ok(Command::Quit);
    }
    match line.split_once(' ') {
        Some(("add", title)) => {
            if title.trim().is_empty() {
                Err(String::from("Атау бос"))
            } else {
                Ok(Command::Add(String::from(title.trim())))
            }
        }
        Some(("done", id)) => Ok(Command::Done(parse_id(id)?)),
        Some(("delete", id)) => Ok(Command::Delete(parse_id(id)?)),
        Some(("rename", rest)) => match rest.split_once(' ') {
            Some((id, title)) => {
                let id = parse_id(id)?;
                if title.trim().is_empty() {
                    Err(String::from("Атау бос"))
                } else {
                    Ok(Command::Rename(id, String::from(title.trim())))
                }
            }
            None => Err(String::from("Белгісіз пәрмен не жарамсыз аргументтер")),
        },
        _ => Err(String::from("Белгісіз пәрмен не жарамсыз аргументтер")),
    }
}

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
    fn add(&mut self, title: String) -> Result<u32, String> {
        if title.trim().is_empty() {
            return Err(String::from("Атау бос"));
        }
        match self.next_id.checked_add(1) {
            Some(next) => {
                let id = self.next_id;
                self.tasks.push(Task {
                    id,
                    title,
                    done: false,
                });
                self.next_id = next;
                Ok(id)
            }
            None => Err(String::from("Нөмірлер таусылды")),
        }
    }
    fn done(&mut self, id: u32) -> Result<(), String> {
        for task in &mut self.tasks {
            if task.id == id {
                task.done = true;
                return Ok(());
            }
        }
        Err(String::from("ID табылмады"))
    }
    fn delete(&mut self, id: u32) -> Result<(), String> {
        for index in 0..self.tasks.len() {
            if self.tasks[index].id == id {
                self.tasks.remove(index);
                return Ok(());
            }
        }
        Err(String::from("ID табылмады"))
    }
    fn rename(&mut self, id: u32, title: String) -> Result<(), String> {
        if title.trim().is_empty() {
            return Err(String::from("Атау бос"));
        }
        for task in &mut self.tasks {
            if task.id == id {
                task.title = title;
                return Ok(());
            }
        }
        Err(String::from("ID табылмады"))
    }
    fn list(&self) {
        if self.tasks.is_empty() {
            println!("Тапсырма жоқ");
        }
        for task in &self.tasks {
            let mark = if task.done {
                "дайын"
            } else {
                "аяқталмаған"
            };
            println!("{}: {} ({mark})", task.id, task.title);
        }
    }
}

fn execute(organizer: &mut Organizer, command: Command) -> bool {
    match command {
        Command::Add(title) => {
            match organizer.add(title) {
                Ok(id) => println!("Қосылды: {id}"),
                Err(message) => println!("Қате: {message}"),
            }
            true
        }
        Command::List => {
            organizer.list();
            true
        }
        Command::Done(id) => {
            match organizer.done(id) {
                Ok(()) => println!("Аяқталды: {id}"),
                Err(message) => println!("Қате: {message}"),
            }
            true
        }
        Command::Delete(id) => {
            match organizer.delete(id) {
                Ok(()) => println!("Өшірілді: {id}"),
                Err(message) => println!("Қате: {message}"),
            }
            true
        }
        Command::Rename(id, title) => {
            match organizer.rename(id, title) {
                Ok(()) => println!("Атауы өзгерді: {id}"),
                Err(message) => println!("Қате: {message}"),
            }
            true
        }
        Command::Quit => false,
    }
}

fn main() {
    let mut organizer = Organizer::new();
    let stdin = io::stdin();
    println!("Пәрмендер: add АТАУ | list | done ID | delete ID | rename ID АТАУ | quit");
    loop {
        let mut line = String::new();
        match stdin.read_line(&mut line) {
            Ok(0) => break,
            Ok(_) => match parse(&line) {
                Ok(command) => {
                    if !execute(&mut organizer, command) {
                        break;
                    }
                }
                Err(message) => println!("Қате: {message}"),
            },
            Err(_) => {
                println!("Қате: енгізу оқылмады");
                break;
            }
        }
    }
    println!("Жұмыс аяқталды");
}
```
```text
Пәрмендер: add АТАУ | list | done ID | delete ID | rename ID АТАУ | quit
Жұмыс аяқталды
```

Тапсырмадағы әрекеттер үшін үлгі жауаптың нәтижесі:

```text
Пәрмендер: add АТАУ | list | done ID | delete ID | rename ID АТАУ | quit
Қосылды: 1
Атауы өзгерді: 1
1: Екі тарау оқу (аяқталмаған)
Қате: ID табылмады
Жұмыс аяқталды
```

Бұл нәтиже бос енгізу арнасына қатысты. Тапсырмадағы пәрмендер енгізілсе, бірінші және соңғы жолдардың арасында құрал жауаптары пайда болады. [Енгізуді оқу мен `split_once` туралы ресми түсіндірме](https://doc.rust-lang.org/std/io/struct.Stdin.html#method.read_line) · [`split_once`](https://doc.rust-lang.org/std/primitive.str.html#method.split_once).

## 40-сабақтан кейінгі бақылау

Бетті жауып, бір пәрменнің жолын көшірмей түсіндіріңіз: енгізу жолы → `parse` → `Command` → әрекет → жауап. Бағдарламаны құрастырып, екі тапсырма қосыңыз, біреуін аяқтап, екіншісін өшіріңіз және жоқ ID-ді сынаңыз. Содан кейін бағдарламаны жауып, тізімнің неге жоғалғанын түсіндіріңіз. Пәрменді талдау қиын болса, [36-сабақты](/read/rust-36-commands?lang=kz); тізімді өзгерту қиын болса, [37-сабақты](/read/rust-37-operations?lang=kz); ID шатасса, [38-сабақты](/read/rust-38-ids?lang=kz); қате тексеруі түсініксіз болса, [39-сабақты](/read/rust-39-scenarios?lang=kz) қайталаңыз.

[Алдыңғы сабақ](/read/rust-39-scenarios?lang=kz) · [Мазмұны](/course/rust?lang=kz)
