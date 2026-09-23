# Тапсырма ID-і: орны ауысса да өзгермейтін нөмір

_Қысқаша (summary):_ **38-сабақ. Әр тапсырмаға тұрақты ID беріп, өшірілгеннен кейін оны қайта пайдаланбау.**

Алдын ала: 10, 13, 27–28, 31–33 және 37-сабақтар. `Vec` орны мен ID шатасса, [37-сабақты](/read/rust-37-operations?lang=kz) қайталаңыз.

## Таныс бейне және тірек сызба

Кітапханада кітап басқа сөреге қойылса да, тіркеу нөмірі өзгермейді. Тапсырма да `Vec` ішіндегі орнын ауыстырып, өзінің **ID-ін**, яғни бірегей жазба нөмірін сақтайды. Ұқсатудың шегі: нөмірді өз бағдарламамыз береді, сондықтан есептеуішті дұрыс жүргізсек қана ол тұрақты болады. Келесі іске қосқанда деректер әлі жоғалады, нөмірлеу қайта басталады.

**Келесі ID → асып кетуді тексеру → қосу → есептеуішті өсіру; өшіру → ескі ID қайта берілмейді.** `Organizer` ішінде `tasks` және `next_id` бар. Келесі нөмір алғашында 1. `checked_add(1)` санды өсіруге болса `Some(келесі нөмір)`, `u32` шегінен асса `None` қайтарады. Біз `u32::MAX` мәнін пайдаланбай қалдырамыз: оған жеткенде `add` тізімді **өзгертпей тұрып** `Err` қайтарады. Бұл ережені жеңілдетеді, бірақ мүмкін нөмірдің бірі қолданылмайды. Өшіргенде есептеуіш кемімейді. `for index in 0..self.tasks.len()` бар орындарды ғана аралайды, сондықтан цикл ішіндегі `self.tasks[index]` шекарадан шықпайды. Іздеу орынмен емес, `task.id` бойынша жүреді. `impl Organizer` ішіндегі `Self` — `Organizer` түрінің өзі, ал `self` — нақты жоспарлау құралы; бұл 28-сабақтан таныс. `u32::MAX` — `u32` түрінің ең үлкен мәні. Оқу мысалындағы `main` `.expect(...)` әдісін тек алдын ала белгілі дұрыс деректерге қолданады; нақты енгізудегі қатені өңдеу керек. Нәтижедегі дөңгелек жақша тапсырма күйін көрсетеді.

## Есептеуіш пен іздеу

Бұрынғы `src/main.rs` файлын сақтап, толық ауыстырыңыз. `organizer` жобасында `cargo run` орындаңыз. Жаңа тәуелділік жоқ. Шығатын екі ID-ді алдын ала болжаңыз.

```rust
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
                    done: false,
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
    fn complete(&mut self, id: u32) -> Result<(), String> {
        for task in &mut self.tasks {
            if task.id == id {
                task.done = true;
                return Ok(());
            }
        }
        Err(String::from("ID табылмады"))
    }
    fn show(&self) {
        for task in &self.tasks {
            let mark = if task.done {
                "дайын"
            } else {
                "аяқталмаған"
            };
            println!("ID {}: {} ({mark})", task.id, task.title);
        }
    }
}

fn main() {
    let mut organizer = Organizer::new();
    let first = organizer.add("Нан алу").expect("fixed valid title");
    let second = organizer.add("Кітап оқу").expect("fixed valid title");
    organizer.remove(first).expect("fixed existing ID");
    let third = organizer.add("Серуен").expect("fixed valid title");
    organizer.complete(second).expect("fixed existing ID");
    println!("ID: {second}, {third}");
    organizer.show();
}
```
```text
ID: 2, 3
ID 2: Кітап оқу (дайын)
ID 3: Серуен (аяқталмаған)
```

Атауы бірдей екі тапсырма да бөлек ID алады: атау іздеу кілті емес. ID 1 өшкенде ID 2 өзгермейді. Жаңа ID 3 есептеуіштің тізім ұзындығына тең емесін көрсетеді. Файлсыз келесі іске қосуда бағдарлама қайтадан 1-ден бастайды; есептеуішті сақтау мен қалпына келтіруді файл кезеңінде қосамыз.

## Түсінгеніңізді тексеріңіз

1. ID 1 өшкеннен кейін жаңа тапсырма неге ID 3 алады?
2. Өшіргеннен кейін 1-орынды іздеу неге басқа тапсырманы өзгертуі мүмкін?
3. Есептеуіш асып кеткенде `add` тізімді жартылай өзгерте ме?

Есептеуіш кемімейді; орындар жылжиды; жоқ, тексеру `push` алдында тұр. Тірек сызбаны қайта айтыңыз.

## Тапсырма

**Міндетті.** Шекараны тексеріңіз: жоқ ID 99-ды аяқтап көріңіз, содан кейін `next_id` мәнін `u32::MAX` етіп, тапсырма қосып көріңіз. Екі қатені көрсетіп, тапсырма саны мен бұрынғы ID-лер өзгермегенін растаңыз.

## Ишара

Әр `Err(message)` мәнін `match` арқылы өңдеңіз. Оқу сынағы үшін сол `main` ішінде `organizer.next_id = u32::MAX` деп жазуға болады; нақты пайдаланушыға мұндай қолжетім берілмейді.

## Талпыныстан кейінгі үлгі жауап

<!-- task-answer -->
```rust
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
                    done: false,
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
    fn complete(&mut self, id: u32) -> Result<(), String> {
        for task in &mut self.tasks {
            if task.id == id {
                task.done = true;
                return Ok(());
            }
        }
        Err(String::from("ID табылмады"))
    }
    fn show(&self) {
        for task in &self.tasks {
            let mark = if task.done {
                "дайын"
            } else {
                "аяқталмаған"
            };
            println!("ID {}: {} ({mark})", task.id, task.title);
        }
    }
}

fn main() {
    let mut organizer = Organizer::new();
    let first = organizer.add("Нан алу").expect("fixed valid title");
    let second = organizer.add("Кітап оқу").expect("fixed valid title");
    organizer.remove(first).expect("fixed existing ID");
    let third = organizer.add("Серуен").expect("fixed valid title");
    organizer.complete(second).expect("fixed existing ID");
    println!("ID: {second}, {third}");
    if let Err(message) = organizer.complete(99) {
        println!("Қате: {message}");
    }
    organizer.next_id = u32::MAX;
    if let Err(message) = organizer.add("Серуен") {
        println!("Қате: {message}");
    }
    println!("Тапсырма саны: {}", organizer.tasks.len());
    organizer.show();
}
```
```text
ID: 2, 3
Қате: ID табылмады
Қате: Нөмірлер таусылды
Тапсырма саны: 2
ID 2: Кітап оқу (дайын)
ID 3: Серуен (аяқталмаған)
```

[`u32::checked_add` құжаттамасы](https://doc.rust-lang.org/std/primitive.u32.html#method.checked_add).

[Алдыңғы сабақ](/read/rust-37-operations?lang=kz) · [Мазмұны](/course/rust?lang=kz) · [Келесі сабақ](/read/rust-39-scenarios?lang=kz)
