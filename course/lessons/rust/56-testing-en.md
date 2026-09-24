# A complete test set for the task file

_Summary:_ **Lesson 56. Test saved tasks, damaged data, repeated reading, and text in different languages.**

Revisit [JSON validation](/read/rust-51-json?lang=en), [backups](/read/rust-53-backup?lang=en), and [your first automated test](/read/rust-16-first-test?lang=en). Use separate files in the system temporary folder for this exercise.

## Familiar image and recall map

Before a trip, you check more than a parked car: you test restarting after a stop, the spare key, and driving in rain. File tests likewise cover different situations. **Recall map:** separate practice path → prepare data → perform action → compare with expectation → remove only your file. Rebuild that sequence from memory. The analogy has a limit: a few tests cannot prove that every possible defect is absent.

A **test set** is several independent checks of behaviour. **Isolation** means a test does not damage real data or depend on another test. ASCII is a small older table of basic Latin letters, digits, and signs; it does not contain `é`. This test checks a character outside that table. `practice_path` adds the process ID and a scenario label to a filename; `unicode`, `damaged`, `duplicate`, and `again` are distinct labels. This reduces collisions when tests run concurrently. `#[test]` tells Cargo that a function is a test; `cargo test` runs tests while `cargo run` executes `main`. `assert_eq!` compares values, and `assert!` requires a true condition. In a test, `unwrap()` turns an unexpected error into a failed test; a normal application should explain the error to its user.

## Run and inspect

Create a separate Cargo project. Add the dependencies below after `[package]` in `Cargo.toml`: the inner `=` requires an exact version, while `features = ["derive"]` enables Serde-generated implementations. Put the first Rust block in `src/main.rs`. `cargo run` prints the number of tasks read. Then run `cargo test`: two tests should pass. The first writes and reads text containing a character outside ASCII; the second leaves damaged JSON unchanged. `read_checked` reads UTF-8 and JSON, then validates version, IDs, and titles. An `Err` does not become an empty organizer.

```toml
[dependencies]
serde = { version = "=1.0.228", features = ["derive"] }
serde_json = "=1.0.149"
```

```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::fs;
use std::path::{Path, PathBuf};

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

fn read_checked(path: &Path) -> Result<Organizer, String> {
    let text = fs::read_to_string(path).map_err(|error| error.to_string())?;
    let data: Organizer = serde_json::from_str(&text).map_err(|error| error.to_string())?;
    if data.version != 1 || data.next_id == 0 {
        return Err(String::from("invalid header"));
    }
    let mut ids = HashSet::new();
    for task in &data.tasks {
        if task.id == 0
            || task.id >= data.next_id
            || task.title.trim().is_empty()
            || !ids.insert(task.id)
        {
            return Err(String::from("invalid task"));
        }
    }
    Ok(data)
}

fn practice_path(label: &str) -> PathBuf {
    std::env::temp_dir().join(format!(
        "rust-course-56-{}-{label}.json",
        std::process::id()
    ))
}

fn sample(title: &str) -> Organizer {
    Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from(title),
            done: false,
        }],
    }
}

fn main() -> Result<(), String> {
    let path = practice_path("main");
    let text = serde_json::to_string(&sample("Buy a book")).map_err(|e| e.to_string())?;
    fs::write(&path, text).map_err(|e| e.to_string())?;
    let loaded = read_checked(&path)?;
    println!("Tasks read: {}", loaded.tasks.len());
    fs::remove_file(path).map_err(|e| e.to_string())?;
    Ok(())
}

#[test]
fn unicode_round_trip() {
    let path = practice_path("unicode");
    let text = serde_json::to_string(&sample("Café plan")).unwrap();
    fs::write(&path, text).unwrap();
    let loaded = read_checked(&path).unwrap();
    assert_eq!(loaded.tasks[0].title, "Café plan");
    fs::remove_file(path).unwrap();
}

#[test]
fn damaged_json_is_rejected() {
    let path = practice_path("damaged");
    fs::write(&path, "{broken").unwrap();
    assert!(read_checked(&path).is_err());
    assert_eq!(fs::read_to_string(&path).unwrap(), "{broken");
    fs::remove_file(path).unwrap();
}
```
```text
Tasks read: 1
```

## Recall without looking

1. How does an isolated test differ from using your real file?
2. Why compare the damaged file contents after the error?
3. What happens if `unwrap()` encounters an unexpected `Err` in a test?

## Exercise

**Required.** Add a test that writes two tasks with the same ID and expects `read_checked` to reject them. Add another test that reads one file twice and confirms the title stayed the same. Give each test a different path label and remove only its own file. Run `cargo test` and confirm that four tests pass.

## Answers

For the repeated ID, start with `sample`, set `next_id = 3`, add a second task with `id: 1`, and write JSON. For repeated reading, compare `first.tasks[0].title` and `second.tasks[0].title`.

<!-- task-answer -->
```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::fs;
use std::path::{Path, PathBuf};

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

fn read_checked(path: &Path) -> Result<Organizer, String> {
    let text = fs::read_to_string(path).map_err(|error| error.to_string())?;
    let data: Organizer = serde_json::from_str(&text).map_err(|error| error.to_string())?;
    if data.version != 1 || data.next_id == 0 {
        return Err(String::from("invalid header"));
    }
    let mut ids = HashSet::new();
    for task in &data.tasks {
        if task.id == 0
            || task.id >= data.next_id
            || task.title.trim().is_empty()
            || !ids.insert(task.id)
        {
            return Err(String::from("invalid task"));
        }
    }
    Ok(data)
}

fn practice_path(label: &str) -> PathBuf {
    std::env::temp_dir().join(format!(
        "rust-course-56-{}-{label}.json",
        std::process::id()
    ))
}

fn sample(title: &str) -> Organizer {
    Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from(title),
            done: false,
        }],
    }
}

fn main() -> Result<(), String> {
    let path = practice_path("main");
    let text = serde_json::to_string(&sample("Buy a book")).map_err(|e| e.to_string())?;
    fs::write(&path, text).map_err(|e| e.to_string())?;
    let loaded = read_checked(&path)?;
    println!("Tasks read: {}", loaded.tasks.len());
    fs::remove_file(path).map_err(|e| e.to_string())?;
    Ok(())
}

#[test]
fn unicode_round_trip() {
    let path = practice_path("unicode");
    let text = serde_json::to_string(&sample("Café plan")).unwrap();
    fs::write(&path, text).unwrap();
    let loaded = read_checked(&path).unwrap();
    assert_eq!(loaded.tasks[0].title, "Café plan");
    fs::remove_file(path).unwrap();
}

#[test]
fn damaged_json_is_rejected() {
    let path = practice_path("damaged");
    fs::write(&path, "{broken").unwrap();
    assert!(read_checked(&path).is_err());
    assert_eq!(fs::read_to_string(&path).unwrap(), "{broken");
    fs::remove_file(path).unwrap();
}

#[test]
fn duplicate_ids_are_rejected() {
    let path = practice_path("duplicate");
    let mut data = sample("Buy a book");
    data.next_id = 3;
    data.tasks.push(Task {
        id: 1,
        title: String::from("second"),
        done: false,
    });
    fs::write(&path, serde_json::to_string(&data).unwrap()).unwrap();
    assert!(read_checked(&path).is_err());
    fs::remove_file(path).unwrap();
}

#[test]
fn second_read_keeps_data() {
    let path = practice_path("again");
    fs::write(&path, serde_json::to_string(&sample("Buy a book")).unwrap()).unwrap();
    let first = read_checked(&path).unwrap();
    let second = read_checked(&path).unwrap();
    assert_eq!(first.tasks[0].title, second.tasks[0].title);
    fs::remove_file(path).unwrap();
}
```
```text
Tasks read: 1
```

## After checking

If a failed test leaves a file behind, inspect it and remove only the practice file. Concurrent tests can run in any order, so do not share one path. Testing write failures and power cuts needs separate scenarios and facilities; this lesson confirms the listed properties, not every possibility. [cargo test documentation](https://doc.rust-lang.org/cargo/commands/cargo-test.html).

[Previous lesson](/read/rust-55-documentation?lang=en) · [Contents](/course/rust?lang=en)
