# Backup copy

_Summary:_ **Lesson 53. Validate a backup and recover its data into a separate file.**

Revisit [JSON](/read/rust-51-json?lang=en) and [saving](/read/rust-52-safe-save?lang=en). This uses a temporary practice folder.

## Familiar image and recall map

A spare key helps only if it fits the lock. Likewise, a backup helps only if it can be read and checked. **Recall map:** primary file → verified copy → damage primary → check copy → separate recovered file. Keep the damaged source for inspection. The analogy does not promise protection when the whole folder disappears.

`load` turns JSON into an `Organizer` and checks version, IDs, and titles. `io::ErrorKind::InvalidData` marks data with invalid shape or meaning. `copy_checked` checks the source first, creates a new target with `create_new(true)`, writes and syncs it, then reads it back. `create_new` refuses to overwrite an existing file. `read_to_string` requires valid UTF-8. In the practice folder, `primary`, `backup`, and `recovered` are three distinct paths. `fs::write` damages only the practice primary file to demonstrate recovery.

## Run and inspect

Create a project with the Serde dependencies from lesson 51. Paste the first Rust block into `src/main.rs` and run `cargo run`. Trace the calls: create a backup before damage; after damage, read the backup and write recovery to a new file. If source validation fails, the target has not been opened. A write failure can leave a partial new target; never treat it as a usable backup without validating it again.

Add these dependencies to `Cargo.toml` (after `[package]`):

In a dependency entry, the inner `=` in `version = "=1.0.228"` requires exactly that version; `features = ["derive"]` enables generated implementations, and the square brackets mark a list of features. These pins make the practice example repeatable; they do not imply the versions will always be current.

```toml
[dependencies]
serde = { version = "=1.0.228", features = ["derive"] }
serde_json = "=1.0.149"
```

```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::fs::{self, OpenOptions};
use std::io::{self, Write};
use std::path::Path;

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

fn load(text: &str) -> io::Result<Organizer> {
    let data: Organizer = serde_json::from_str(text)
        .map_err(|error| io::Error::new(io::ErrorKind::InvalidData, error))?;
    if data.version != 1 || data.next_id == 0 {
        return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid header"));
    }
    let mut ids = HashSet::new();
    for task in &data.tasks {
        if task.id == 0
            || task.id >= data.next_id
            || task.title.trim().is_empty()
            || !ids.insert(task.id)
        {
            return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid task"));
        }
    }
    Ok(data)
}

fn copy_checked(source: &Path, target: &Path) -> io::Result<()> {
    let text = fs::read_to_string(source)?;
    load(&text)?;
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(target)?;
    file.write_all(text.as_bytes())?;
    file.sync_all()?;
    drop(file);
    load(&fs::read_to_string(target)?)?;
    Ok(())
}

fn run() -> io::Result<()> {
    let folder = std::env::temp_dir().join(format!("rust-lesson-53-{}", std::process::id()));
    fs::create_dir(&folder)?;
    let primary = folder.join("tasks.json");
    let backup = folder.join("backup.json");
    let recovered = folder.join("recovered.json");
    let data = Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from("Buy a book"),
            done: false,
        }],
    };
    let json = serde_json::to_string(&data)
        .map_err(|error| io::Error::new(io::ErrorKind::InvalidData, error))?;
    fs::write(&primary, json)?;
    copy_checked(&primary, &backup)?;
    fs::write(&primary, "{broken")?;
    println!(
        "Primary file damaged: {}",
        load(&fs::read_to_string(&primary)?).is_err()
    );
    copy_checked(&backup, &recovered)?;
    let restored = load(&fs::read_to_string(&recovered)?)?;
    println!("Recovered: {}", restored.tasks[0].title);
    println!(
        "Source unchanged: {}",
        fs::read_to_string(&primary)? == "{broken"
    );
    fs::remove_dir_all(folder)?;
    Ok(())
}

fn main() {
    if let Err(error) = run() {
        eprintln!("Error: {error}");
        std::process::exit(1);
    }
}
```
```text
Primary file damaged: true
Recovered: Buy a book
Source unchanged: true
```

## Recall without looking

1. Why validate the source first?
2. Why recover into a third file?
3. What should you do if reading the copy back fails?

## Exercise

**Required.** Create a separate file containing invalid JSON `{broken`, call `copy_checked` with a new target path, and print that it failed and no recovery file appeared. Keep the good backup.

## Answers

Call `copy_checked(&bad_backup, &bad_target)`. Use `is_err()` for the result and `!bad_target.exists()` for the untouched target.

<!-- task-answer -->
```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::fs::{self, OpenOptions};
use std::io::{self, Write};
use std::path::Path;

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

fn load(text: &str) -> io::Result<Organizer> {
    let data: Organizer = serde_json::from_str(text)
        .map_err(|error| io::Error::new(io::ErrorKind::InvalidData, error))?;
    if data.version != 1 || data.next_id == 0 {
        return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid header"));
    }
    let mut ids = HashSet::new();
    for task in &data.tasks {
        if task.id == 0
            || task.id >= data.next_id
            || task.title.trim().is_empty()
            || !ids.insert(task.id)
        {
            return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid task"));
        }
    }
    Ok(data)
}

fn copy_checked(source: &Path, target: &Path) -> io::Result<()> {
    let text = fs::read_to_string(source)?;
    load(&text)?;
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(target)?;
    file.write_all(text.as_bytes())?;
    file.sync_all()?;
    drop(file);
    load(&fs::read_to_string(target)?)?;
    Ok(())
}

fn run() -> io::Result<()> {
    let folder = std::env::temp_dir().join(format!("rust-lesson-53-{}", std::process::id()));
    fs::create_dir(&folder)?;
    let primary = folder.join("tasks.json");
    let backup = folder.join("backup.json");
    let recovered = folder.join("recovered.json");
    let data = Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from("Buy a book"),
            done: false,
        }],
    };
    let json = serde_json::to_string(&data)
        .map_err(|error| io::Error::new(io::ErrorKind::InvalidData, error))?;
    fs::write(&primary, json)?;
    copy_checked(&primary, &backup)?;
    fs::write(&primary, "{broken")?;
    println!(
        "Primary file damaged: {}",
        load(&fs::read_to_string(&primary)?).is_err()
    );
    copy_checked(&backup, &recovered)?;
    let restored = load(&fs::read_to_string(&recovered)?)?;
    println!("Recovered: {}", restored.tasks[0].title);
    println!(
        "Source unchanged: {}",
        fs::read_to_string(&primary)? == "{broken"
    );
    let damaged_backup = folder.join("damaged-backup.json");
    let other = folder.join("other-recovery.json");
    fs::write(&damaged_backup, "{broken")?;
    let rejected = copy_checked(&damaged_backup, &other).is_err() && !other.exists();
    println!("Bad backup rejected: {}", rejected);
    fs::remove_dir_all(folder)?;
    Ok(())
}

fn main() {
    if let Err(error) = run() {
        eprintln!("Error: {error}");
        std::process::exit(1);
    }
}
```
```text
Primary file damaged: true
Recovered: Buy a book
Source unchanged: true
Bad backup rejected: true
```

## After checking

A validated copy is not a complete retention strategy: it may be old or disappear with the primary. Show the recovered version to the user before choosing it. Return to [lesson 51](/read/rust-51-json?lang=en). [create_new documentation](https://doc.rust-lang.org/std/fs/struct.OpenOptions.html#method.create_new).

[Previous lesson](/read/rust-52-safe-save?lang=en) · [Contents](/course/rust?lang=en)
