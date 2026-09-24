# JSON format

_Summary:_ **Lesson 51. Turn tasks into JSON and validate them after reading.**

Revisit [Task structs](/read/rust-27-structs?lang=en), [external crates](/read/rust-48-dependencies?lang=en), and [file reading](/read/rust-50-files?lang=en). This lesson uses an in-memory string; disk storage comes next.

## Familiar image and recall map

A task list is like cards in a labelled envelope. JSON is a text format with shared rules for field names and values. Another program can read it, but the envelope does not prove the cards are correct. **Recall map:** values → JSON → read JSON → validate meaning → usable tasks. Rebuild this route from memory.

`serde` converts structures to data and back; `serde_json` handles JSON. **Serialization** turns a value into text; **deserialization** rebuilds it. `#[derive(Serialize, Deserialize)]` asks the compiler to create standard implementations; `#[serde(deny_unknown_fields)]` rejects fields the structure does not recognise. An attribute `#[...]` is a compiler instruction attached to the next declaration. `version` records the file shape, so a changed shape cannot silently masquerade as the old one. `next_id` must exceed every existing ID. `HashSet` holds unique values; `insert` returns `false` for a repeated ID. `?` passes an error to the caller. `map_err` converts a library error into text for `run`.

## Run and inspect

Use a separate Cargo project. Put the dependencies below in `Cargo.toml`, replace `src/main.rs` with the first Rust block, and run `cargo run`. Cargo downloads crates on the first build; `Cargo.lock` pins the resolved versions. The example has one task. `to_string` creates JSON text, `from_str` rebuilds a value, and `validate` checks rules JSON cannot know. This output follows the struct field order; another writer may choose another order.

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

fn validate(data: &Organizer) -> Result<(), String> {
    if data.version != 1 || data.next_id == 0 {
        return Err(String::from("invalid version or next ID"));
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
    Ok(())
}

fn run() -> Result<(), String> {
    let source = Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from("Buy a book"),
            done: false,
        }],
    };
    let json = serde_json::to_string(&source).map_err(|error| error.to_string())?;
    let loaded: Organizer = serde_json::from_str(&json).map_err(|error| error.to_string())?;
    validate(&loaded)?;
    println!("{json}");
    println!("Restored tasks: {}", loaded.tasks.len());
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
{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Buy a book","done":false}]}
Restored tasks: 1
```

## Recall without looking

1. How does parsing JSON differ from validating its meaning?
2. Why keep `next_id` when the tasks already have IDs?
3. Why reject unknown fields?

## Exercise

**Required.** Add a second task with ID 2, set `next_id` to 3, and print the number of loaded tasks. Then change the loaded second task ID to 1 and show that validation rejects it. Predict both output lines first.

## Answers

Change the value before `to_string`. After `from_str`, you have an `Organizer` again. For the invalid version, change `loaded.tasks[1].id` and call validation again. Check `validate(&bad).is_err()`.

<!-- task-answer -->
```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;

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

fn validate(data: &Organizer) -> Result<(), String> {
    if data.version != 1 || data.next_id == 0 {
        return Err(String::from("invalid version or next ID"));
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
    Ok(())
}

fn run() -> Result<(), String> {
    let source = Organizer {
        version: 1,
        next_id: 3,
        tasks: vec![
            Task {
                id: 1,
                title: String::from("Buy a book"),
                done: false,
            },
            Task {
                id: 2,
                title: String::from("Call a friend"),
                done: false,
            },
        ],
    };
    let json = serde_json::to_string(&source).map_err(|error| error.to_string())?;
    let mut loaded: Organizer = serde_json::from_str(&json).map_err(|error| error.to_string())?;
    validate(&loaded)?;
    println!("{json}");
    println!("Restored tasks: {}", loaded.tasks.len());
    loaded.tasks[1].id = 1;
    println!("Duplicate ID rejected: {}", validate(&loaded).is_err());
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
{"version":1,"next_id":3,"tasks":[{"id":1,"title":"Buy a book","done":false},{"id":2,"title":"Call a friend","done":false}]}
Restored tasks: 2
Duplicate ID rejected: true
```

## After checking

JSON checks syntax; `validate` checks organizer rules. Never turn damaged data into an empty task list. Return to [lesson 50](/read/rust-50-files?lang=en). [Serde documentation](https://serde.rs/derive.html).

[Previous lesson](/read/rust-50-files?lang=en) · [Contents](/course/rust?lang=en)
