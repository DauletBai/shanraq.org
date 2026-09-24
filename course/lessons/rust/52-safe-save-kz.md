# Деректерді жоғалтпай сақтау

_Қысқаша (summary):_ **Сабақ 52. Жаңа деректі көршілес уақытша файлға жазып, ақау кезінде ескісін сақтау.**

[Жолдарды](/read/rust-49-paths?lang=kz) және [файл қатесін](/read/rust-50-files?lang=kz) еске түсіріңіз. Қазір оқу мәтінін сақтаймыз; 51-сабақтағы JSON-ды осы қызметке беруге болады.

## Таныс бейне және тірек сызба

Дәптердегі ескі бетті жаңасын жазып бітірмей өшірмейсіз. Алдымен бөлек параққа жазып, оны тексересіз, содан соң бетті алмастырасыз. **Тірек сызба:** ескі файл → көршілес уақытша файл → жазу → құрылғыға жеткізуді сұрау → жабу → алмастыру. Жазу сәтсіз болса, ескі файл қалады. Ұқсастықтың шегі: қуаттың кенет үзілуі мен файлдық жүйе ерекшеліктері қосымша шараны қажет етеді.

`with_extension` негізгі файлдың қасында уақытша файл жолын жасайды; бір файлдық жүйедегі алмастыру үшін бұл маңызды. `create_new(true)` бұрыннан бар уақытша файлды үстінен жаздырмайды. `write_all` барлық байтты жазады немесе қате қайтарады; `sync_all` файл дерегі мен қосымша мәліметін құрылғыға жеткізуді операциялық жүйеден сұрайды; `drop` ашық файлды жабады. `rename` жүйе рұқсат етсе, нысана жолды алмастырады, әйтпесе қате қайтарады. `if let Err(error)` тек қате тармағын таңдайды; `let _ = remove_file` тазалау қатесі сақтаудың сәтті нәтижесіне айналмайтынын білдіреді.

## Іске қосып, талдаңыз

Тәуелділігі жоқ жаңа Cargo жобасын ашып, алғашқы Rust үзіндісін `src/main.rs` файлына көшіріңіз де, `cargo run` орындаңыз. Оқу қалтасының атауына үдеріс нөмірі қосылған; бағдарлама тоқтап қалса, ішін қарап барып осы қалтаны ғана өшіріңіз. Үлгіде ескі мәтін жаңасымен алмастырылады. Уақытша файлды жасау, жазу не жеткізу қателері `rename` қадамына жеткізбейді, сондықтан негізгі файл сақталады.

```rust
use std::fs::{self, OpenOptions};
use std::io::{self, Write};
use std::path::Path;

fn safe_save(path: &Path, text: &str) -> io::Result<()> {
    let temporary = path.with_extension(format!("tmp-{}", std::process::id()));
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(&temporary)?;
    if let Err(error) = file.write_all(text.as_bytes()) {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    if let Err(error) = file.sync_all() {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    drop(file);
    if let Err(error) = fs::rename(&temporary, path) {
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    Ok(())
}

fn run() -> io::Result<()> {
    let folder = std::env::temp_dir().join(format!("rust-lesson-52-{}", std::process::id()));
    fs::create_dir(&folder)?;
    let path = folder.join("tasks.json");
    fs::write(&path, "Ескі жоспар")?;
    safe_save(&path, "Жаңа жоспар")?;
    println!("Сақталды: {}", fs::read_to_string(&path)?);
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
Сақталды: Жаңа жоспар
```

## Бетке қарамай жауап беріңіз

1. Уақытша файл неге негізгі файлдың жанында?
2. `create_new` неге пайдалы?
3. Ескі файл алғаш қай қадамда алмастырылуы мүмкін?

## Тапсырма

**Міндетті.** `safe_save` шақырмай тұрып дәл сондай уақытша файл жасаңыз. Сақтау қатемен аяқталғанын және негізгі файлда ескі мәтін қалғанын көрсетіңіз. Кедергі файлды өшіріп, қайта сақтаңыз да жаңа мәтінді шығарыңыз.

## Жауаптар

Уақытша жолды қызметтегідей есептеңіз: `path.with_extension(format!("tmp-{}", std::process::id()))`. `is_err()` нәтижесін және ескі мәтінді тексеріңіз.

<!-- task-answer -->
```rust
use std::fs::{self, OpenOptions};
use std::io::{self, Write};
use std::path::Path;

fn safe_save(path: &Path, text: &str) -> io::Result<()> {
    let temporary = path.with_extension(format!("tmp-{}", std::process::id()));
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(&temporary)?;
    if let Err(error) = file.write_all(text.as_bytes()) {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    if let Err(error) = file.sync_all() {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    drop(file);
    if let Err(error) = fs::rename(&temporary, path) {
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    Ok(())
}

fn run() -> io::Result<()> {
    let folder = std::env::temp_dir().join(format!("rust-lesson-52-{}", std::process::id()));
    fs::create_dir(&folder)?;
    let path = folder.join("tasks.json");
    fs::write(&path, "Ескі жоспар")?;
    let temporary = path.with_extension(format!("tmp-{}", std::process::id()));
    fs::write(&temporary, "occupied")?;
    let refused = safe_save(&path, "Жаңа жоспар").is_err();
    let kept = fs::read_to_string(&path)? == "Ескі жоспар";
    println!("Ескі дерек сақталды: {}", refused && kept);
    fs::remove_file(&temporary)?;
    safe_save(&path, "Жаңа жоспар")?;
    println!("Сақталды: {}", fs::read_to_string(&path)?);
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
Ескі дерек сақталды: true
Сақталды: Жаңа жоспар
```

## Тексергеннен кейін

Бұл кәдімгі жазу қатесінен қорғайды, бірақ кез келген ақауда дерек сақталады деп уәде бермейді. Уақытша файлға жасалған `sync_all` алмастырудан кейінгі қуат үзілуіне толық кепіл емес; қатаң кепіл үшін операциялық жүйе, файлдық жүйе және қалта дерегін құрылғыға жеткізу ережелері қажет. Екі үдерістің қатар жазуы да шешілмеген. [rename құжаттамасы](https://doc.rust-lang.org/std/fs/fn.rename.html).

[Алдыңғы сабақ](/read/rust-51-json?lang=kz) · [Мазмұн](/course/rust?lang=kz)
