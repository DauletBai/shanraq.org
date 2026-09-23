# Пәрмен — тексерілген ниет

_Қысқаша (summary):_ **36-сабақ. Іске қосқанда берілген сөздерді талдауды жоспарлау құралының әрекетінен бөліп, әр қатені түсіндіру.**

Алдын ала: 29–35-сабақтар. `Option` пен `Result` ұмытылса, [31-сабақты](/read/rust-31-option?lang=kz) қайталаңыз.

## Таныс бейне және тірек сызба

Қабылдау орнында өтінішті алдымен тексеріп, содан кейін ғана орындайды. Жоспарлау құралы да сөздерді әуелі тексерілген пәрменге айналдырады. **Пәрменді талдау** — әрекет атауы мен оған керекті деректерді тексеру. Бұл кезеңде тапсырмалар тізімі өзгермейді. Ұқсатудың шегі: бағдарлама кез келген сөйлемді түсінбейді; біз анықтаған қалыптарды ғана қабылдайды.

**Сөздер → қалпын тексеру → `Result<Command, String>` → пәрмені бар `Ok` не себебі бар `Err` → әрекет кейін орындалады.** `Command` — 29-сабақтағы өзіміз жасаған санамалау. `Add(String)` атауды сақтайды, ал `List` дерек сақтамайды. `Option<&str>` мәтінге сілтеме бар (`Some`) не жоқ (`None`) дегенді білдіреді; бұл 31-сабақтан таныс. Үш мәнді кортеж арқылы `match` әрекетті, атауды және артық аргументті бірге тексереді. `title.trim().is_empty()` шеткі бос орындарды алып тастаған соң мәтін қалды ма, соны анықтайды. `String::from(title)` дербес мәтін жасайды: бастапқы аргументтер жойылса да, пәрмен атауды сақтайды. `_` қалған қалыптарды қамтып, себебін қайтарады. Мұндағы `Ok(...)` пен `Err(...)` нәтиже мәндері, олар ештеңе басып шығармайды.

## Әуелі талдау, кейін көрсету

`organizer` жобасындағы бұрынғы `src/main.rs` файлын сақтап, оны толық ауыстырыңыз. `Cargo.toml` тұрған бумада `cargo run` орындаңыз. Бұл мысал алдын ала берілген сөздерді талдайды; нақты іске қосу аргументтерін толық құралды құрастырғанда жалғаймыз. Өзге файлдар мен тәуелділіктер өзгермейді. Іске қоспай тұрып нәтижені болжаңыз.

```rust
enum Command {
    Add(String),
    List,
}

fn parse(action: &str, title: Option<&str>, extra: Option<&str>) -> Result<Command, String> {
    match (action, title, extra) {
        ("list", None, None) => Ok(Command::List),
        ("add", Some(title), None) => {
            if title.trim().is_empty() {
                Err(String::from("Бос атау"))
            } else {
                Ok(Command::Add(String::from(title)))
            }
        }
        _ => Err(String::from("Пәрмен не аргумент саны қате")),
    }
}

fn show(result: Result<Command, String>) {
    match result {
        Ok(Command::Add(title)) => println!("Қосу: {title}"),
        Ok(Command::List) => println!("Тізімді көрсету"),
        Err(message) => println!("Қате: {message}"),
    }
}

fn main() {
    show(parse("add", Some("Rust оқу"), None));
    show(parse("list", None, None));
    show(parse("add", Some("   "), None));
}
```
```text
Қосу: Rust оқу
Тізімді көрсету
Қате: Бос атау
```

`parse` тек тексерілген пәрмен жасайды. `show` оны көрсетеді; тапсырмалар тізімін келесі сабақта қосамыз. `parse("list", Some("артық"), None)` шақырылса, қалыпты тексеру сөзді елемей қоймай, қате қайтарады.

## Түсінгеніңізді тексеріңіз

1. Неге `parse` тапсырманы бірден қоспай, пәрмен қайтарады?
2. `None` пен `Some("")` несімен өзгеше?
3. Қате хабарын қай қызмет басып шығарады?

Талдау тізімді өзгертпейді; бірінде атау жоқ, екіншісінде бос атау берілген; хабарды `main` шығарады. Бетті жапқан соң тірек сызбаны қайта айтыңыз.

## Тапсырма

**Міндетті.** Санамалауға бір нөмірі бар `done` әрекеті үшін `Done(u32)` қосыңыз. Нөмірді `parse::<u32>()` арқылы талдап, нөл мен жарамсыз мәтінді қабылдамаңыз; артық аргументті де тексеріңіз. `done 2`, `done 0`, `done жоқ` жағдайларын сынаңыз. Талдау қызметі ештеңе басып шығармауы керек.

## Ишара

Жаңа `("done", Some(text), None)` тармағында `text.parse::<u32>()` қолданыңыз. `Ok(number)` мен `Err(_)` нұсқаларын `match` арқылы өңдеп, нөлге `Err` қайтарыңыз.

## Талпыныстан кейінгі үлгі жауап

<!-- task-answer -->
```rust
enum Command {
    Add(String),
    List,
    Done(u32),
}

fn parse(action: &str, title: Option<&str>, extra: Option<&str>) -> Result<Command, String> {
    match (action, title, extra) {
        ("list", None, None) => Ok(Command::List),
        ("add", Some(title), None) => {
            if title.trim().is_empty() {
                Err(String::from("Бос атау"))
            } else {
                Ok(Command::Add(String::from(title)))
            }
        }
        ("done", Some(text), None) => match text.parse::<u32>() {
            Ok(0) => Err(String::from("Нөмір нөлден үлкен болуы керек")),
            Ok(number) => Ok(Command::Done(number)),
            Err(_) => Err(String::from("Нөмір қажет")),
        },
        _ => Err(String::from("Пәрмен не аргумент саны қате")),
    }
}

fn show(result: Result<Command, String>) {
    match result {
        Ok(Command::Add(title)) => println!("Қосу: {title}"),
        Ok(Command::List) => println!("Тізімді көрсету"),
        Ok(Command::Done(number)) => println!("Аяқтау: {number}"),
        Err(message) => println!("Қате: {message}"),
    }
}

fn main() {
    show(parse("add", Some("Rust оқу"), None));
    show(parse("list", None, None));
    show(parse("add", Some("   "), None));
    show(parse("done", Some("2"), None));
    show(parse("done", Some("0"), None));
    show(parse("done", Some("жоқ"), None));
}
```
```text
Қосу: Rust оқу
Тізімді көрсету
Қате: Бос атау
Аяқтау: 2
Қате: Нөмір нөлден үлкен болуы керек
Қате: Нөмір қажет
```

[Санамалау мен `match` туралы ресми түсіндірме](https://doc.rust-lang.org/book/ch06-02-match.html).

[Алдыңғы сабақ](/read/rust-35-arguments?lang=kz) · [Мазмұны](/course/rust?lang=kz) · [Келесі сабақ](/read/rust-37-operations?lang=kz)
