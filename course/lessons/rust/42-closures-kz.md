# Тұйықталым: керекті тапсырмаларды таңдау

_Қысқаша (summary):_ **42-сабақ. Сөз бойынша тапсырмаларды тауып, тұрақты ретпен көрсету.**

Алдын ала: [итераторлар](/read/rust-41-iterators?lang=kz), [мәтін жолдары](/read/rust-22-strings?lang=kz), [тұрақты ID](/read/rust-38-ids?lang=kz). Алдыңғы `main.rs` файлының көшірмесін сақтаңыз: мысал файлды түгел ауыстырады, дерек әлі жадта ғана тұрады.

## Таныс бейне және тірек сызба

Кітапханадағы көмекшіге «атауында Rust сөзі бар қағаздарды ірікте» деген қысқа ереже берілді делік. **Тұйықталым** (closure) — айналасындағы мәнді есінде ұстай алатын атаусыз қызмет. `|task| task.title.contains(query)` жазылымындағы тік сызықтар кіріс `task` мәнін қоршайды; `query` сырттағы айнымалыдан алынады. Ұқсатудың шегі: ереже берілген сәтте іріктеу бірден жүрмейді. Итератор нәтижені сұрағанда ғана жұмыс бастайды.

**Барлық тапсырма → `filter` (сәйкесін қалдыру) → `collect` (жинау) → `sort_by_key` (реттеу) → көрсету.** `filter` шартты тексеретін тұйықталымды қолданады. `iter()` элементті `&Task` сілтемесімен береді, ал `filter` тексеру үшін сол сілтеменің өзін қарызға алады. Сондықтан оның кірісі `&&Task` — сілтемеге сілтеме. Rust `.title` өрісін оқығанда екі қабатты өзі өтеді; `sort_by_key` ішіндегі `.id` пен `map` ішіндегі `.title` да осылай оқылады. `contains` таңбаларды дәл іздейді: үлкен және кіші әріп бөлек, сондықтан `rust` пен `Rust` бірдей емес. `collect` бұл жерде `Vec<&Task>` — тапсырмаларға сілтемелер тізімін жасайды; бастапқы тапсырмалар орнында қалады. `let found: Vec<&Task>` ішіндегі қос нүкте нәтиженің түрін көрсетеді. `sort_by_key(|task| task.id)` атаумен емес, ID бойынша реттейді. `map` әр тапсырмадан атауының `&str` сілтемесін алады; `as_str()` мәтінді көшірмейді. `names.len()` нәтиже санын береді. Мұнда `names` тізімі `map` жұмысын көрсету үшін жасалған; тек сан керек болса, жаңа тізім құрмай `found.len()` қолдануға болады. Мұндағы тік сызықтар Rust жазылымының бөлігі, 40-сабақтағы пәрмендерді көрсететін ажыратқыш емес.

## Іздеп көріңіз

`src/main.rs` файлын ауыстырып, `cargo run` пәрменін орындаңыз. Бағдарлама нәтижесі Cargo хабарларынан бөлек көрсетілген. Қосымша кітапхана жоқ.

```rust
struct Task {
    id: u32,
    title: String,
}

fn main() {
    let tasks = vec![
        Task {
            id: 3,
            title: String::from("Rust оқу"),
        },
        Task {
            id: 2,
            title: String::from("Нан алу"),
        },
        Task {
            id: 1,
            title: String::from("Rust қайталау"),
        },
    ];
    let query = "Rust";
    let mut found: Vec<&Task> = tasks
        .iter()
        .filter(|task| task.title.contains(query))
        .collect();
    found.sort_by_key(|task| task.id);
    let names: Vec<&str> = found.iter().map(|task| task.title.as_str()).collect();
    println!("Табылды: {}", names.len());
    for task in found {
        println!("{}: {}", task.id, task.title);
    }
}
```
```text
Табылды: 2
1: Rust қайталау
3: Rust оқу
```

Бастапқы тізімде ID 3 алдымен тұрса да, нәтижеде ID 1 алда шығады. Бастапқы `tasks` өзгермейді. Ештеңе табылмаса, `found` бос болады және цикл еш жол шығармайды: бұл қате емес. Әр іздеу тізімді қайта аралайды; жылдамдатуды кейін қарастырамыз.

## Бетке қарамай жауап беріңіз

1. Тұйықталым `query` мәнін қайдан алады?
2. `collect` бастапқы тапсырмаларды неге жоймайды?
3. `sort_by_key` тапсырмалардың тұрақты ID-сін өзгерте ме?

## Тапсырма

**Міндетті.** Негізгі нәтижеден соң атауында `Нан` сөзі бар тапсырмалар санын `filter` және `count` арқылы тауып, `«Нан» сөзі бойынша: 1` шығарыңыз. Бастапқы тізімді өзгертпеңіз. Кейін `нан` деп іздеп, айырмасын түсіндіріңіз.

## Жауаптар

`second_query` жасаңыз. `tasks.iter().filter(|task| task.title.contains(second_query)).count()` іріктеуді орындап, санын қайтарады.

## Талпыныстан кейінгі үлгі жауап

<!-- task-answer -->
```rust
struct Task {
    id: u32,
    title: String,
}

fn main() {
    let tasks = vec![
        Task {
            id: 3,
            title: String::from("Rust оқу"),
        },
        Task {
            id: 2,
            title: String::from("Нан алу"),
        },
        Task {
            id: 1,
            title: String::from("Rust қайталау"),
        },
    ];
    let query = "Rust";
    let mut found: Vec<&Task> = tasks
        .iter()
        .filter(|task| task.title.contains(query))
        .collect();
    found.sort_by_key(|task| task.id);
    let names: Vec<&str> = found.iter().map(|task| task.title.as_str()).collect();
    println!("Табылды: {}", names.len());
    for task in found {
        println!("{}: {}", task.id, task.title);
    }
    let second_query = "Нан";
    let second_count = tasks
        .iter()
        .filter(|task| task.title.contains(second_query))
        .count();
    println!("«{second_query}» сөзі бойынша: {second_count}");
}
```
```text
Табылды: 2
1: Rust қайталау
3: Rust оқу
«Нан» сөзі бойынша: 1
```

## Тексергеннен кейін

`sort_by_key` табылған сілтемелерді реттейді, тапсырмалардың ID-сі өзгермейді. [Тұйықталымдар туралы ресми тарау](https://doc.rust-lang.org/book/ch13-01-closures.html) · [итераторлар туралы тарау](https://doc.rust-lang.org/book/ch13-02-iterators.html). Кейінге қалдырылған аралау түсініксіз болса, [41-сабақты](/read/rust-41-iterators?lang=kz) қайталаңыз.

[Алдыңғы сабақ](/read/rust-41-iterators?lang=kz) · [Мазмұны](/course/rust?lang=kz)
