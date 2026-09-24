# Releasing the console organizer

_Summary:_ **Lesson 60. Combine familiar parts into a program another person can run and verify.**

You need [JSON](/read/rust-51-json?lang=en), [saving](/read/rust-52-safe-save?lang=en), [backups](/read/rust-53-backup?lang=en), [commands](/read/rust-54-cli?lang=en), [tests](/read/rust-56-testing?lang=en), and [priority](/read/rust-59-capstone?lang=en). This is the final console project: each command is a separate launch, while tasks remain in a JSON file.

## Familiar image and recall map

Imagine a workshop: task cards sit in a cabinet, a card’s ID survives removal of its neighbour, and a spare cabinet is opened only after checking its key. **Recall map:** arguments and path → read file → validate → act → save through a temporary file → clear output or error. Backups follow another branch: validate source → create a new file → read it back. Rebuild the map without looking. The analogy has a limit: two people writing to the same cabinet at once need separate coordination; this program does not provide it.

`--data PATH` selects a file explicitly; otherwise the program uses the `ORGANIZER_DATA` environment setting, then `tasks.json` in the current folder. `std::env::args()` reads arguments as Unicode strings; paths with non-UTF-8 bytes are unsupported here. A missing file means a first launch only: the program then creates an empty state in memory. An access error or damaged JSON returns an error rather than an empty list. `Result<String, (i32, String)>` carries either normal output or a pair of exit status and error text. Status 2 means a bad command; 1 means a data or file problem. `eprintln!` writes to stderr. `retain` keeps list elements that satisfy a condition. When sorting by priority, `then` compares IDs if priorities tie. `to_owned()` creates an owned string from a text view.

## Final program

Use the ready project folder for this lesson, or create a separate Cargo project with the dependencies below. Put the entire first Rust block in `src/main.rs`; do not paste pieces into your real organizer data. `cargo run` without arguments shows help. In `cargo run -- --data tasks.json add 3 "Buy a book"`, the first `--` ends Cargo’s options, `--data` selects a practice file, `3` means urgent, and quotes keep a title with spaces together. This Cargo command works from the project folder on Windows, macOS, and Linux. The service words `add`, `list`, `done`, `delete`, `rename`, `search`, `filter`, `sort`, `stats`, `backup`, `restore`, and `help` stay in English; messages for the person are localised.

The [complete practice project](https://github.com/DauletBai/shanraq.org/tree/main/course/rust-organizer/en/step-60) will be available in the repository; it contains `Cargo.toml`, `Cargo.lock`, `README.md`, `LICENSE`, and `src/main.rs`.

Read help as alternatives: `|` separates commands; do not type them all together. In a task line, `[3]` means priority 3; the square brackets are part of the printed message. Before running, locate five parts of the code: data model and validation, reading, careful saving and copying, command execution, and path selection in `main`.

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
            return Err(String::from("Invalid data"));
        }
        let mut ids = HashSet::new();
        for task in &self.tasks {
            if task.id == 0
                || task.id >= self.next_id
                || task.title.trim().is_empty()
                || !(1..=3).contains(&task.priority)
                || !ids.insert(task.id)
            {
                return Err(String::from("Invalid data"));
            }
        }
        Ok(())
    }
    fn add(&mut self, priority: u8, title: String) -> Result<u64, String> {
        if !(1..=3).contains(&priority) {
            return Err(String::from("Priority must be 1 to 3"));
        }
        if title.trim().is_empty() {
            return Err(String::from("Title is empty"));
        }
        let next = self.next_id.checked_add(1).ok_or("No IDs remain")?;
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
            .ok_or(String::from("ID not found"))
    }
    fn delete(&mut self, id: u64) -> Result<(), String> {
        let before = self.tasks.len();
        self.tasks.retain(|task| task.id != id);
        if self.tasks.len() == before {
            Err(String::from("ID not found"))
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
        _ => Err(String::from("A positive ID is required")),
    }
}

fn parse_priority(text: &str) -> Result<u8, String> {
    match text.parse::<u8>() {
        Ok(priority) if (1..=3).contains(&priority) => Ok(priority),
        _ => Err(String::from("Priority must be 1 to 3")),
    }
}

fn show(tasks: &[&Task]) -> String {
    if tasks.is_empty() {
        return String::from("No tasks");
    }
    tasks
        .iter()
        .map(|task| {
            let status = if task.done { "done" } else { "open" };
            format!("{} [{}] {} ({status})", task.id, task.priority, task.title)
        })
        .collect::<Vec<_>>()
        .join("\n")
}

fn execute(args: &[String], path: &Path) -> Result<String, (i32, String)> {
    let usage = || (2, String::from("Unknown command or invalid arguments"));
    let data_error = |message: String| (1, message);
    if args.is_empty() || args == ["help"] {
        return Ok(String::from(
            "Commands: add PRIORITY TITLE | list | done ID | delete ID | rename ID TITLE | search TEXT | filter open/done/high | sort id/title/priority | stats | backup PATH | restore BACKUP NEW_PATH | --version",
        ));
    }
    if args == ["--version"] {
        return Ok(env!("CARGO_PKG_VERSION").to_owned());
    }
    match args[0].as_str() {
        "backup" if args.len() == 2 => {
            copy_checked(path, Path::new(&args[1])).map_err(data_error)?;
            Ok(String::from("Backup created"))
        }
        "restore" if args.len() == 3 => {
            copy_checked(Path::new(&args[1]), Path::new(&args[2])).map_err(data_error)?;
            Ok(String::from("Recovered separately"))
        }
        "add" if args.len() >= 3 => {
            let priority = parse_priority(&args[1]).map_err(|e| (2, e))?;
            let mut data = load_or_new(path).map_err(data_error)?;
            let id = data
                .add(priority, args[2..].join(" "))
                .map_err(|e| (2, e))?;
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Added: {id}"))
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
            Ok(format!("Done: {id}"))
        }
        "delete" if args.len() == 2 => {
            let id = parse_id(&args[1]).map_err(|e| (2, e))?;
            let mut data = load_or_new(path).map_err(data_error)?;
            data.delete(id).map_err(data_error)?;
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Deleted: {id}"))
        }
        "rename" if args.len() >= 3 => {
            let id = parse_id(&args[1]).map_err(|e| (2, e))?;
            let title = args[2..].join(" ");
            if title.trim().is_empty() {
                return Err((2, String::from("Title is empty")));
            }
            let mut data = load_or_new(path).map_err(data_error)?;
            data.task_mut(id).map_err(data_error)?.title = title.trim().to_owned();
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Renamed: {id}"))
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
                "Total: {total}, done: {done}, open: {open}, high open: {high}"
            ))
        }
        _ => Err(usage()),
    }
}

fn main() {
    let mut args: Vec<String> = std::env::args().skip(1).collect();
    let path = if args.first().map(String::as_str) == Some("--data") {
        if args.len() < 2 {
            eprintln!("Error: Unknown command or invalid arguments");
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
            eprintln!("Error: {message}");
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
    let first = data.add(2, String::from("Buy a book")).unwrap();
    let second = data.add(3, String::from("Buy a book")).unwrap();
    data.delete(first).unwrap();
    assert_eq!(data.tasks[0].id, second);
    assert_eq!(data.next_id, second + 1);
}
```
```text
Commands: add PRIORITY TITLE | list | done ID | delete ID | rename ID TITLE | search TEXT | filter open/done/high | sort id/title/priority | stats | backup PATH | restore BACKUP NEW_PATH | --version
```

| Action | Practice file state | Visible result |
|---|---|---|
| `help` | No file created | Command list |
| `add 3 "Buy a book"` | Task 1, next ID 2 | `Added: 1` |
| `list` in a new launch | Reads the same file | `1 [3] Buy a book (open)` |
| `done 1` | Task 1 is done | `Done: 1` |
| `restore backup.json recovered.json` | Primary is unchanged | New `recovered.json` is created |



`list` shows ID, priority in square brackets, title, and state. `search` looks for an exact case-sensitive substring; `sort title` follows Rust string ordering, not a local dictionary order. `filter high` shows only unfinished tasks with priority 3. `backup PATH` makes a checked copy; `restore BACKUP NEW_PATH` recovers only into a new file. You can then open that new file with `--data` without erasing a damaged source. `save_safe` writes through a neighbouring temporary file. It handles ordinary write errors before replacement but does not guarantee survival of a power cut or simultaneous writers. Run the command sequence below in a separate practice folder. Choose another practice name if a file already exists.

## Recall without looking

1. Why is damaged data not treated as a first launch?
2. Why does `restore` require a new path?
3. How do `search` and `sort title` differ from language-aware search and sorting?
4. What happens after a write failure before `rename`?

## Exercise

**Required.** In a fresh practice folder run `cargo test`, then add a priority-3 task with a title containing spaces, list it in a separate launch, mark it done, and list again. Make a backup, damage only the practice primary JSON, confirm `list` fails without erasing it, recover the backup to a new path, and read that path with `--data`. Try an unknown command and invalid ID. Add tests for old JSON without `priority`, a repeated ID, and an invalid priority. Hand the README to a friend and ask them to repeat the route without your explanation.

## Answers

Pass `-- --data PATH` on each `cargo run` so separate launches share one practice file. Damage only a practice file you created. In tests call `decode` rather than your real home file. A command reference follows the answer code.

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
            return Err(String::from("Invalid data"));
        }
        let mut ids = HashSet::new();
        for task in &self.tasks {
            if task.id == 0
                || task.id >= self.next_id
                || task.title.trim().is_empty()
                || !(1..=3).contains(&task.priority)
                || !ids.insert(task.id)
            {
                return Err(String::from("Invalid data"));
            }
        }
        Ok(())
    }
    fn add(&mut self, priority: u8, title: String) -> Result<u64, String> {
        if !(1..=3).contains(&priority) {
            return Err(String::from("Priority must be 1 to 3"));
        }
        if title.trim().is_empty() {
            return Err(String::from("Title is empty"));
        }
        let next = self.next_id.checked_add(1).ok_or("No IDs remain")?;
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
            .ok_or(String::from("ID not found"))
    }
    fn delete(&mut self, id: u64) -> Result<(), String> {
        let before = self.tasks.len();
        self.tasks.retain(|task| task.id != id);
        if self.tasks.len() == before {
            Err(String::from("ID not found"))
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
        _ => Err(String::from("A positive ID is required")),
    }
}

fn parse_priority(text: &str) -> Result<u8, String> {
    match text.parse::<u8>() {
        Ok(priority) if (1..=3).contains(&priority) => Ok(priority),
        _ => Err(String::from("Priority must be 1 to 3")),
    }
}

fn show(tasks: &[&Task]) -> String {
    if tasks.is_empty() {
        return String::from("No tasks");
    }
    tasks
        .iter()
        .map(|task| {
            let status = if task.done { "done" } else { "open" };
            format!("{} [{}] {} ({status})", task.id, task.priority, task.title)
        })
        .collect::<Vec<_>>()
        .join("\n")
}

fn execute(args: &[String], path: &Path) -> Result<String, (i32, String)> {
    let usage = || (2, String::from("Unknown command or invalid arguments"));
    let data_error = |message: String| (1, message);
    if args.is_empty() || args == ["help"] {
        return Ok(String::from(
            "Commands: add PRIORITY TITLE | list | done ID | delete ID | rename ID TITLE | search TEXT | filter open/done/high | sort id/title/priority | stats | backup PATH | restore BACKUP NEW_PATH | --version",
        ));
    }
    if args == ["--version"] {
        return Ok(env!("CARGO_PKG_VERSION").to_owned());
    }
    match args[0].as_str() {
        "backup" if args.len() == 2 => {
            copy_checked(path, Path::new(&args[1])).map_err(data_error)?;
            Ok(String::from("Backup created"))
        }
        "restore" if args.len() == 3 => {
            copy_checked(Path::new(&args[1]), Path::new(&args[2])).map_err(data_error)?;
            Ok(String::from("Recovered separately"))
        }
        "add" if args.len() >= 3 => {
            let priority = parse_priority(&args[1]).map_err(|e| (2, e))?;
            let mut data = load_or_new(path).map_err(data_error)?;
            let id = data
                .add(priority, args[2..].join(" "))
                .map_err(|e| (2, e))?;
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Added: {id}"))
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
            Ok(format!("Done: {id}"))
        }
        "delete" if args.len() == 2 => {
            let id = parse_id(&args[1]).map_err(|e| (2, e))?;
            let mut data = load_or_new(path).map_err(data_error)?;
            data.delete(id).map_err(data_error)?;
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Deleted: {id}"))
        }
        "rename" if args.len() >= 3 => {
            let id = parse_id(&args[1]).map_err(|e| (2, e))?;
            let title = args[2..].join(" ");
            if title.trim().is_empty() {
                return Err((2, String::from("Title is empty")));
            }
            let mut data = load_or_new(path).map_err(data_error)?;
            data.task_mut(id).map_err(data_error)?.title = title.trim().to_owned();
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Renamed: {id}"))
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
                "Total: {total}, done: {done}, open: {open}, high open: {high}"
            ))
        }
        _ => Err(usage()),
    }
}

fn main() {
    let mut args: Vec<String> = std::env::args().skip(1).collect();
    let path = if args.first().map(String::as_str) == Some("--data") {
        if args.len() < 2 {
            eprintln!("Error: Unknown command or invalid arguments");
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
            eprintln!("Error: {message}");
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
    let first = data.add(2, String::from("Buy a book")).unwrap();
    let second = data.add(3, String::from("Buy a book")).unwrap();
    data.delete(first).unwrap();
    assert_eq!(data.tasks[0].id, second);
    assert_eq!(data.next_id, second + 1);
}

#[test]
fn old_json_gets_normal_priority() {
    let text = r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Buy a book","done":false}]}"#;
    assert_eq!(decode(text).unwrap().tasks[0].priority, 2);
}

#[test]
fn duplicate_id_is_rejected() {
    let text = r#"{"version":1,"next_id":3,"tasks":[{"id":1,"title":"Buy a book","done":false,"priority":2},{"id":1,"title":"Call a friend","done":false,"priority":3}]}"#;
    assert!(decode(text).is_err());
}

#[test]
fn invalid_priority_is_rejected() {
    let text =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Buy a book","done":false,"priority":0}]}"#;
    assert!(decode(text).is_err());
}
```
```text
Commands: add PRIORITY TITLE | list | done ID | delete ID | rename ID TITLE | search TEXT | filter open/done/high | sort id/title/priority | stats | backup PATH | restore BACKUP NEW_PATH | --version
```

```text
cargo run -- --data tasks.json add 3 "Buy a book"
cargo run -- --data tasks.json list
cargo run -- --data tasks.json done 1
cargo run -- --data tasks.json backup backup.json
cargo run -- --data tasks.json restore backup.json recovered.json
cargo run -- --data recovered.json list
```

`README.md`:

```markdown
# Console organizer

Install Rust 1.97.0 and Cargo using https://www.rust-lang.org/tools/install, open a new terminal, and check `cargo --version`. In the project folder run `cargo test`, then `cargo run -- --data tasks.json add 3 "Buy a book"`, and separately `cargo run -- --data tasks.json list`. `cargo run -- help` shows other commands. Data stays in the selected JSON file, not in the executable. Before risky changes, make a copy with `backup backup.json`. `restore backup.json recovered.json` writes a new path and does not erase the source. Damaged JSON causes an error: do not replace it with an empty list. Without `--data`, the program uses `ORGANIZER_DATA`, then `tasks.json` in the current folder. Build separately for another OS.
```

## After checking

Check the exit status along with error text before automating anything. Strong protection against sudden power loss and concurrent writers needs additional file-system and directory-sync policy. A separate window with buttons can later use these validated operations; the required course result is a console program. [Official Rust book](https://doc.rust-lang.org/book/).

[Previous lesson](/read/rust-59-capstone?lang=en) · [Contents](/course/rust?lang=en)
