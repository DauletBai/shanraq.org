# Жадтағы тапсырмалар: қосу, өзгерту, аяқтау, өшіру

_Қысқаша (summary):_ **37-сабақ. Тапсырмаларға төрт әрекет жасап, өшіргенде тізім орындарының жылжуын көру.**

Алдын ала: 21, 24, 27, 31–33 және 36-сабақтар. Өзгермелі қарызға алуды ұмытсаңыз, [21-сабаққа](/read/rust-21-mutable-borrow?lang=kz) оралыңыз.

## Таныс бейне және тірек сызба

Қатарға қойылған тапсырма қағаздарын елестетіңіз. Жаңасын қосуға, атауын түзетуге, орындалды деп белгілеуге немесе алып тастауға болады. `Vec<Task>` — бағдарлама жадындағы осындай қатар. Бағдарлама жабылғанда қағаздар әзірше жоғалады: файлға сақтауды кейін үйренеміз. Ұқсатудың шегі: қатардағы орын қағаздың тұрақты нөмірі емес.

**Тексеру → `Vec<Task>` өзгерту → нәтижені көрсету; бағдарламаны жабу → жадтағы дерек жоғалады.** `Task` құрамында `title: String` және `done: bool` өрістері бар. `done` тапсырма аяқталғанда ақиқат болады. `&mut Vec<Task>` тізімді өзгертуге уақытша дара қолжетім береді, ал `&[Task]` тізім бөлігін тек оқуға береді. `get_mut(index)` `Option<&mut Task>` қайтарады: орын болмауы мүмкін. `tasks.remove(index)` бөлікті өшіріп, кейінгілерді солға жылжытады; оның алдында шекараны тексереміз, өйткені жарамсыз орын берілсе, `remove` бағдарламаны тоқтатады. `usize` — 15-сабақтағы тізім орнының түрі. `Result<(), String>` сәтті әрекет мән қайтармайтынын (`()`), ал қате себебін жеткізетінін білдіреді. `for task in tasks` қағаздарды ретімен оқиды.

## Төрт әрекет

Бұрынғы `src/main.rs` файлын сақтап, толық ауыстырыңыз. `organizer` жобасында `cargo run` орындаңыз. Жаңа тәуелділік жоқ. 0-орын өшірілген соң қай қағаз қалатынын болжаңыз. `[ ]` әлі аяқталмаған тапсырманы, `[дайын]` аяқталғанын білдіреді. Осы оқу мысалындағы `main` ішінде `.expect(...)` сәтті `Ok` мәнін алады, бірақ `Err` болса бағдарламаны тоқтатады; оны тек алдын ала белгілі дұрыс деректерге қолданыңыз.

```rust
struct Task {
    title: String,
    done: bool,
}

fn add(tasks: &mut Vec<Task>, title: &str) -> Result<(), String> {
    if title.trim().is_empty() {
        return Err(String::from("Атау бос"));
    }
    tasks.push(Task {
        title: String::from(title),
        done: false,
    });
    Ok(())
}

fn rename(tasks: &mut Vec<Task>, index: usize, title: &str) -> Result<(), String> {
    match tasks.get_mut(index) {
        Some(task) => {
            task.title = String::from(title);
            Ok(())
        }
        None => Err(String::from("Ондай орын жоқ")),
    }
}

fn complete(tasks: &mut Vec<Task>, index: usize) -> Result<(), String> {
    match tasks.get_mut(index) {
        Some(task) => {
            task.done = true;
            Ok(())
        }
        None => Err(String::from("Ондай орын жоқ")),
    }
}

fn remove(tasks: &mut Vec<Task>, index: usize) -> Result<(), String> {
    if index >= tasks.len() {
        return Err(String::from("Ондай орын жоқ"));
    }
    tasks.remove(index);
    Ok(())
}

fn show(tasks: &[Task]) {
    for task in tasks {
        if task.done {
            println!("[дайын] {}", task.title);
        } else {
            println!("[ ] {}", task.title);
        }
    }
}

fn main() {
    let mut tasks = Vec::new();
    add(&mut tasks, "Нан алу").expect("valid title");
    add(&mut tasks, "Тарау оқу").expect("valid title");
    rename(&mut tasks, 1, "Екі тарау оқу").expect("existing position");
    complete(&mut tasks, 0).expect("existing position");
    remove(&mut tasks, 0).expect("existing position");
    show(&tasks);
}
```
```text
[ ] Екі тарау оқу
```

Оқу мысалындағы `main` ішінде `.expect(...)` қолдануға болады, себебі мұнда бар орындар мен бос емес атаулар алдын ала берілген. Бұл `Result` әдісі `Ok` мәнін алады, бірақ `Err` болса бағдарламаны тоқтатады; нақты енгізуде бұлай істемеу керек. Тапсырмада қатені анық өңдейміз. Өшіруден кейін екінші бөлік 0-орынға жылжиды. Сондықтан келесі сабақта тұрақты ID енгіземіз.

## Түсінгеніңізді тексеріңіз

1. Неге `get_mut` бірден `&mut Task` емес, `Option` қайтарады?
2. Бірінші қағаз өшсе, екіншісінің орны қалай өзгереді?
3. Келесі іске қосқанда осы тапсырмалар қала ма?

Орын болмауы мүмкін; екінші қағаз 0-орынға жылжиды; жоқ, сақтау әлі жасалмаған. Тірек сызбаны қайта айтыңыз.

## Тапсырма

**Міндетті.** `rename` қызметіне тек бос орыннан тұратын атауды қабылдамайтын тексеру қосыңыз. Мұндай жағдайда бұрынғы атау өзгермесін. Сәтті әрекеттерден кейін 0-орынның атауын `"   "` деп өзгертіп, 9-орынды аяқтап көріңіз. Екі қатені, содан соң тізімді көрсетіңіз; бұрынғы атау қалуы керек.

## Ишара

`title.trim().is_empty()` тексеруін `get_mut` алдында қойыңыз. Әр `Result` қатесін 32-сабақтағыдай `match` арқылы көрсетіңіз.

## Талпыныстан кейінгі үлгі жауап

<!-- task-answer -->
```rust
struct Task {
    title: String,
    done: bool,
}

fn add(tasks: &mut Vec<Task>, title: &str) -> Result<(), String> {
    if title.trim().is_empty() {
        return Err(String::from("Атау бос"));
    }
    tasks.push(Task {
        title: String::from(title),
        done: false,
    });
    Ok(())
}

fn rename(tasks: &mut Vec<Task>, index: usize, title: &str) -> Result<(), String> {
    if title.trim().is_empty() {
        return Err(String::from("Атау бос"));
    }
    match tasks.get_mut(index) {
        Some(task) => {
            task.title = String::from(title);
            Ok(())
        }
        None => Err(String::from("Ондай орын жоқ")),
    }
}

fn complete(tasks: &mut Vec<Task>, index: usize) -> Result<(), String> {
    match tasks.get_mut(index) {
        Some(task) => {
            task.done = true;
            Ok(())
        }
        None => Err(String::from("Ондай орын жоқ")),
    }
}

fn remove(tasks: &mut Vec<Task>, index: usize) -> Result<(), String> {
    if index >= tasks.len() {
        return Err(String::from("Ондай орын жоқ"));
    }
    tasks.remove(index);
    Ok(())
}

fn show(tasks: &[Task]) {
    for task in tasks {
        if task.done {
            println!("[дайын] {}", task.title);
        } else {
            println!("[ ] {}", task.title);
        }
    }
}

fn main() {
    let mut tasks = Vec::new();
    add(&mut tasks, "Нан алу").expect("valid title");
    add(&mut tasks, "Тарау оқу").expect("valid title");
    rename(&mut tasks, 1, "Екі тарау оқу").expect("existing position");
    complete(&mut tasks, 0).expect("existing position");
    remove(&mut tasks, 0).expect("existing position");
    for result in [rename(&mut tasks, 0, "   "), complete(&mut tasks, 9)] {
        if let Err(message) = result {
            println!("Қате: {message}");
        }
    }
    show(&tasks);
}
```
```text
Қате: Атау бос
Қате: Ондай орын жоқ
[ ] Екі тарау оқу
```

[`Vec::get_mut` пен `Vec::remove` құжаттамасы](https://doc.rust-lang.org/std/vec/struct.Vec.html#method.remove).

[Алдыңғы сабақ](/read/rust-36-commands?lang=kz) · [Мазмұны](/course/rust?lang=kz) · [Келесі сабақ](/read/rust-38-ids?lang=kz)
