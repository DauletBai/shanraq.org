# Әрекеттер тізбегін сынау

_Қысқаша (summary):_ **39-сабақ. Бірнеше әрекеттің салдарын және қате кезінде деректің өзгермеуін тексеру.**

Алдын ала: 16, 32–33 және 38-сабақтар. `#[test]` ұмытылса, [16-сабаққа](/read/rust-16-first-test?lang=kz) оралыңыз.

## Таныс бейне және тірек сызба

Құлыпты тексеру үшін кілтті бір рет бұрау аз: ашып, қайта жауып, басқа кілтті де сынайсыз. **Сценарий сынағы** жоспарлау құралындағы әрекеттер тізбегін осылай тексереді. Ұқсатудың шегі: сынақтар таңдалған жағдайларды қамтиды, барлық қатенің жоқтығын дәлелдемейді.

**Бастапқы күй → әрекет → күтілген нәтиже → жаңа күй; қате → бұрынғы күй сақталады.** `#[test]` пен `assert_eq!` 16-сабақтан таныс. Әр сынақ өз `Organizer` мәнін жасайды, сондықтан сынақтардың орындалу реті нәтижеге әсер етпеуі керек. `Err(String::from(...))` — күтілген қате. `Result` мәндерін `assert_eq!` арқылы салыстыру нұсқаны да, ішіндегі мәтінді де тексереді. `o.tasks[0].title` салыстырылғанда макрос мәтінді тізімнен алып кетпей оқиды. Әрекет сәтті аяқталса да, күйдің дұрыстығы өзінен-өзі дәлелденбейді: өшірген соң қалған ID мен `next_id` мәнін бөлек тексереміз. Қате кезінде тапсырмалар мен есептеуіш өзгермегенін де растау керек.

## Екі дербес сценарий

Бұрынғы `src/main.rs` файлын сақтап, толық ауыстырыңыз. `cargo run` қысқа хабар шығарады. Содан кейін `cargo test` орындаңыз: Cargo `#[test]` белгісі бар қызметтерді іске қосады. Cargo-ның қызметтік жолдары әртүрлі болуы мүмкін; барлық сынақтың өтуі маңызды. Жаңа тәуелділік жоқ. Қысқа тексеру үшін мұндағы `Task` тек `id` мен `title` сақтайды: 38-сабақтағы `done` өрісі уақытша алынған, оны 40-сабақта қайта қосамыз. Әр сабақ оқу үшін `main.rs` файлын толық ауыстырады, бұрынғы іске қосудың дерегін көшірмейді.

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
            return Err(String::from("Атау бос"));
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
            None => Err(String::from("Нөмірлер таусылды")),
        }
    }
    fn remove(&mut self, id: u32) -> Result<(), String> {
        for index in 0..self.tasks.len() {
            if self.tasks[index].id == id {
                self.tasks.remove(index);
                return Ok(());
            }
        }
        Err(String::from("ID табылмады"))
    }
}
fn main() {
    let mut o = Organizer::new();
    o.add("Оқу").expect("fixed title");
    o.add("Жазу").expect("fixed title");
    o.remove(1).expect("fixed ID");
    println!("Қалды: {}", o.tasks[0].title);
}

#[test]
fn ids_survive_removal() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Оқу"), Ok(1));
    assert_eq!(o.add("Жазу"), Ok(2));
    assert_eq!(o.remove(1), Ok(()));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 2);
    assert_eq!(o.tasks[0].title, "Жазу");
    assert_eq!(o.next_id, 3);
}

#[test]
fn overflow_changes_nothing() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Оқу"), Ok(1));
    o.next_id = u32::MAX;
    assert_eq!(o.add("Жазу"), Err(String::from("Нөмірлер таусылды")));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 1);
    assert_eq!(o.next_id, u32::MAX);
}
```
```text
Қалды: Жазу
```

Бірінші тексеру «ID тізім орнына тең» деген қатені табады. Екіншісі есептеуішті әдейі шекке қойып, `add` деректі өзгертпей тұрып бас тартатынын тексереді. Оқу мысалында `title` өрісі де салыстырылады: басқа қағаздағы дұрыс ID жалған сәттілік болар еді. 38-сабақтағыдай, бағдарлама аяқталғанда дерек жойылады.

## Түсінгеніңізді тексеріңіз

1. ID 1 өшкеннен кейін тек `Ok(2)` тексеру неге жеткіліксіз?
2. `Err` жағдайынан кейін тағы нені тексереміз?
3. Бұл сынақтар бір-біріне тәуелді ме?

Қалған ID мен есептеуішті де тексеру керек; күй өзгермеген; әр сынақ жаңа мән жасайды. Тірек сызбаны қайта айтыңыз.

## Тапсырма

**Міндетті.** Жеке сынақ қосыңыз: екі сәтті қосудан кейін жоқ ID 99-ды өшіру қате қайтарады, алғашқы екі ID сақталады, ал келесі сәтті `add` ID 3 береді. `cargo test` орындаңыз. Тек қате хабарын тексерумен шектелмеңіз.

## Ишара

`#[test] fn invalid_remove_keeps_ids()` ішінде жаңа `Organizer` жасаңыз. `assert_eq!(o.remove(99), Err(...))` жазыңыз; содан соң `o.tasks.len()`, екі `id` және `o.add(...)` мәндерін тексеріңіз.

## Талпыныстан кейінгі үлгі жауап

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
            return Err(String::from("Атау бос"));
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
            None => Err(String::from("Нөмірлер таусылды")),
        }
    }
    fn remove(&mut self, id: u32) -> Result<(), String> {
        for index in 0..self.tasks.len() {
            if self.tasks[index].id == id {
                self.tasks.remove(index);
                return Ok(());
            }
        }
        Err(String::from("ID табылмады"))
    }
}
fn main() {
    let mut o = Organizer::new();
    o.add("Оқу").expect("fixed title");
    o.add("Жазу").expect("fixed title");
    o.remove(1).expect("fixed ID");
    println!("Қалды: {}", o.tasks[0].title);
}

#[test]
fn ids_survive_removal() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Оқу"), Ok(1));
    assert_eq!(o.add("Жазу"), Ok(2));
    assert_eq!(o.remove(1), Ok(()));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 2);
    assert_eq!(o.tasks[0].title, "Жазу");
    assert_eq!(o.next_id, 3);
}

#[test]
fn overflow_changes_nothing() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Оқу"), Ok(1));
    o.next_id = u32::MAX;
    assert_eq!(o.add("Жазу"), Err(String::from("Нөмірлер таусылды")));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 1);
    assert_eq!(o.next_id, u32::MAX);
}

#[test]
fn invalid_remove_keeps_ids() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Оқу"), Ok(1));
    assert_eq!(o.add("Жазу"), Ok(2));
    assert_eq!(o.remove(99), Err(String::from("ID табылмады")));
    assert_eq!(o.tasks.len(), 2);
    assert_eq!(o.tasks[0].id, 1);
    assert_eq!(o.tasks[1].id, 2);
    assert_eq!(o.add("Оқу"), Ok(3));
}
```
```text
Қалды: Жазу
```

Бұл нәтиже `cargo run` үшін; `cargo test` үш сынақтың өткенін бөлек көрсетеді. [Rust сынақтары туралы ресми нұсқаулық](https://doc.rust-lang.org/book/ch11-01-writing-tests.html).

[Алдыңғы сабақ](/read/rust-38-ids?lang=kz) · [Мазмұны](/course/rust?lang=kz) · [Келесі сабақ](/read/rust-40-memory-checkpoint?lang=kz)
