# Scenarios: test a sequence of actions

_Lead (summary):_ **Lesson 39. Check consequences of several actions and errors that must leave data unchanged.**

Prerequisites: lessons 16, 32–33 and 38. Revisit [lesson 16](/read/rust-16-first-test?lang=en) if `#[test]` is unclear.

## Familiar image and recall map

To check a lock, one turn of the key is not enough: open it, close it, and try the wrong key. A **scenario test** likewise checks a sequence of organizer actions. The analogy has a limit: tests cover selected cases; they cannot prove that all bugs are absent.

**Initial state → action → expected result → new state; error → old state retained.** `#[test]` and `assert_eq!` come from lesson 16. Every test creates its own `Organizer`, so test order must not affect the outcome. `Err(String::from(...))` is an expected error. Comparing `Result` with `assert_eq!` checks both the variant and its contents. When comparing `o.tasks[0].title`, the macro reads the string without taking it out of the list. A successful return does not alone prove the state is correct: after removal, we separately check the remaining ID and `next_id`. For errors, also confirm that tasks and counter did not change.

## Two independent scenarios

Save and completely replace the previous `src/main.rs`. `cargo run` prints a short message. Then run `cargo test`: Cargo runs functions marked `#[test]`. Its progress text may vary; the important result is that every test passes. No new dependencies are needed. To keep this check short, `Task` contains only `id` and `title` here: the `done` field from lesson 38 is temporarily omitted and returns in lesson 40. Each lesson replaces the teaching `main.rs` completely; it does not carry data from an earlier run.

```rust
struct Task {
    id: u32,
    title: String,
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
            return Err(String::from("Empty title"));
        }
        match self.next_id.checked_add(1) {
            Some(next) => {
                let id = self.next_id;
                self.tasks.push(Task {
                    id,
                    title: String::from(title),
                });
                self.next_id = next;
                Ok(id)
            }
            None => Err(String::from("No IDs left")),
        }
    }
    fn remove(&mut self, id: u32) -> Result<(), String> {
        for index in 0..self.tasks.len() {
            if self.tasks[index].id == id {
                self.tasks.remove(index);
                return Ok(());
            }
        }
        Err(String::from("ID not found"))
    }
}
fn main() {
    let mut o = Organizer::new();
    o.add("Read").expect("fixed title");
    o.add("Write").expect("fixed title");
    o.remove(1).expect("fixed ID");
    println!("Remaining: {}", o.tasks[0].title);
}

#[test]
fn ids_survive_removal() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Read"), Ok(1));
    assert_eq!(o.add("Write"), Ok(2));
    assert_eq!(o.remove(1), Ok(()));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 2);
    assert_eq!(o.tasks[0].title, "Write");
    assert_eq!(o.next_id, 3);
}

#[test]
fn overflow_changes_nothing() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Read"), Ok(1));
    o.next_id = u32::MAX;
    assert_eq!(o.add("Write"), Err(String::from("No IDs left")));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 1);
    assert_eq!(o.next_id, u32::MAX);
}
```
```text
Remaining: Write
```

The first check catches a mistaken assumption that ID equals list position. The second deliberately puts the counter at its boundary and checks that `add` refuses before changing data. The example checks `title` too: the right ID attached to the wrong card would be a false success. Data still disappear when the process ends, as in lesson 38.

## Check your understanding

1. Why is checking `Ok(2)` after removal insufficient?
2. What else do we check after an `Err`?
3. Do these tests depend on one another?

Check the retained ID and counter too; state has not changed; each test builds a fresh object. Rebuild the recall map.

## Exercise

**Required.** Add a separate test: after two successful additions, removing unknown ID 99 returns an error, both original IDs remain, and the next successful `add` returns ID 3. Run `cargo test`. Do more than checking the error message.

## Hint

Create a new `Organizer` inside `#[test] fn invalid_remove_keeps_ids()`. Use `assert_eq!(o.remove(99), Err(...))`, then check `o.tasks.len()`, both IDs, and `o.add(...)`.

## Reference after attempting

<!-- task-answer -->
```rust
struct Task {
    id: u32,
    title: String,
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
            return Err(String::from("Empty title"));
        }
        match self.next_id.checked_add(1) {
            Some(next) => {
                let id = self.next_id;
                self.tasks.push(Task {
                    id,
                    title: String::from(title),
                });
                self.next_id = next;
                Ok(id)
            }
            None => Err(String::from("No IDs left")),
        }
    }
    fn remove(&mut self, id: u32) -> Result<(), String> {
        for index in 0..self.tasks.len() {
            if self.tasks[index].id == id {
                self.tasks.remove(index);
                return Ok(());
            }
        }
        Err(String::from("ID not found"))
    }
}
fn main() {
    let mut o = Organizer::new();
    o.add("Read").expect("fixed title");
    o.add("Write").expect("fixed title");
    o.remove(1).expect("fixed ID");
    println!("Remaining: {}", o.tasks[0].title);
}

#[test]
fn ids_survive_removal() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Read"), Ok(1));
    assert_eq!(o.add("Write"), Ok(2));
    assert_eq!(o.remove(1), Ok(()));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 2);
    assert_eq!(o.tasks[0].title, "Write");
    assert_eq!(o.next_id, 3);
}

#[test]
fn overflow_changes_nothing() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Read"), Ok(1));
    o.next_id = u32::MAX;
    assert_eq!(o.add("Write"), Err(String::from("No IDs left")));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 1);
    assert_eq!(o.next_id, u32::MAX);
}

#[test]
fn invalid_remove_keeps_ids() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Read"), Ok(1));
    assert_eq!(o.add("Write"), Ok(2));
    assert_eq!(o.remove(99), Err(String::from("ID not found")));
    assert_eq!(o.tasks.len(), 2);
    assert_eq!(o.tasks[0].id, 1);
    assert_eq!(o.tasks[1].id, 2);
    assert_eq!(o.add("Read"), Ok(3));
}
```
```text
Remaining: Write
```

This output belongs to `cargo run`; `cargo test` separately reports three passing tests. [Official guide to Rust tests](https://doc.rust-lang.org/book/ch11-01-writing-tests.html).

[Previous lesson](/read/rust-38-ids?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-40-memory-checkpoint?lang=en)
