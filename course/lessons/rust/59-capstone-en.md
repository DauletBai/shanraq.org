# Independent extension: task priority

_Summary:_ **Lesson 59. Add priority while still reading old JSON and rejecting invalid values.**

Revisit [structs](/read/rust-27-structs?lang=en), [JSON validation](/read/rust-51-json?lang=en), and [file tests](/read/rust-56-testing?lang=en). Write the behaviour before the code: this is independent work with a reference solution after the exercise.

## Familiar image and recall map

A task card gets an urgency mark. An old card without a mark does not become urgent by itself: it gets the ordinary mark. **Recall map:** state the rule → add a field → assign a value to old records → reject invalid values → check filtering → explain it to the user. The analogy has a limit: a mark does not decide what really matters; a person does.

**Priority** is a number from 1 to 3: 1 low, 2 ordinary, 3 urgent. Old JSON has no `priority` field. `#[serde(default = "ordinary_priority")]` tells Serde to call `ordinary_priority` when that field is absent. It does not repair an explicitly invalid number. `1..=3` is an inclusive range: both endpoints belong. `high_open` counts tasks that are urgent and not yet done. `version` can remain 1 for this compatible addition if old files still load and the new field has a documented meaning; an incompatible change needs a separate version and migration decision.

## Run and inspect

In a separate Cargo project add the dependencies below and put the first Rust block in `src/main.rs`. Run `cargo run` and `cargo test`. The starting program reads old JSON without priority and reports the number of tasks. It validates version, positive unique IDs, and nonblank titles. Do not edit that old JSON by hand: it is the compatibility example.

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

fn decode(text: &str) -> Result<Organizer, String> {
    let data: Organizer = serde_json::from_str(text).map_err(|e| e.to_string())?;
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

fn main() -> Result<(), String> {
    let old_json =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Buy a book","done":false}]}"#;
    let data = decode(old_json)?;
    println!("Tasks: {}", data.tasks.len());
    Ok(())
}

#[test]
fn reads_old_task() {
    let text = r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Buy a book","done":false}]}"#;
    assert_eq!(decode(text).unwrap().tasks.len(), 1);
}
```
```text
Tasks: 1
```

## Recall without looking

1. Why must an old file get a defined default priority?
2. Should `default` repair an explicitly written `priority: 0`?
3. Why does the urgent count exclude completed tasks?

## Exercise

**Required.** Before coding, write four rules: only 1, 2, and 3 are valid; a missing field in old JSON means 2; an explicit 0 or 4 is an error; only open tasks with 3 are counted as urgent. Then add the field and counting function, plus tests for old and invalid JSON. Print the old task priority and the urgent count after adding a new open task with priority 3. Predict both lines first.

## Answers

Define `fn ordinary_priority() -> u8 { 2 }` before the struct. Put `#[serde(default = "ordinary_priority")]` above the new field. Validate with `!(1..=3).contains(&task.priority)` and count with `iter().filter(...).count()`.

<!-- task-answer -->
```rust
use serde::{Deserialize, Serialize};
use std::collections::HashSet;

fn ordinary_priority() -> u8 {
    2
}

#[derive(Serialize, Deserialize)]
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

fn decode(text: &str) -> Result<Organizer, String> {
    let data: Organizer = serde_json::from_str(text).map_err(|e| e.to_string())?;
    if data.version != 1 || data.next_id == 0 {
        return Err(String::from("invalid header"));
    }
    let mut ids = HashSet::new();
    for task in &data.tasks {
        if task.id == 0
            || task.id >= data.next_id
            || task.title.trim().is_empty()
            || !(1..=3).contains(&task.priority)
            || !ids.insert(task.id)
        {
            return Err(String::from("invalid task"));
        }
    }
    Ok(data)
}

fn high_open(data: &Organizer) -> usize {
    data.tasks
        .iter()
        .filter(|task| task.priority == 3 && !task.done)
        .count()
}

fn main() -> Result<(), String> {
    let old_json =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Buy a book","done":false}]}"#;
    let mut data = decode(old_json)?;
    println!("Legacy task priority: {}", data.tasks[0].priority);
    data.tasks.push(Task {
        id: 2,
        title: String::from("Call a friend"),
        done: false,
        priority: 3,
    });
    data.next_id = 3;
    println!("Open high-priority tasks: {}", high_open(&data));
    Ok(())
}

#[test]
fn old_file_gets_ordinary_priority() {
    let text = r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Buy a book","done":false}]}"#;
    assert_eq!(decode(text).unwrap().tasks[0].priority, 2);
}

#[test]
fn invalid_priority_is_rejected() {
    let text = r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Buy a book","done":false,"priority":0}]}"#;
    assert!(decode(text).is_err());
}
```
```text
Legacy task priority: 2
Open high-priority tasks: 1
```

## After checking

Old files must remain readable, but damaged or explicitly invalid files must never silently become an empty list. You can propose another extension by the same method; first state the rule for old data and a rejection test. [Serde default documentation](https://serde.rs/field-attrs.html#default).

[Previous lesson](/read/rust-58-distribution?lang=en) · [Contents](/course/rust?lang=en)
