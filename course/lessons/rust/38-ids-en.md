# Task IDs: numbers that do not shift

_Lead (summary):_ **Lesson 38. Give each task a stable ID and never reuse it after removal.**

Prerequisites: lessons 10, 13, 27–28, 31–33 and 37. Revisit [lesson 37](/read/rust-37-operations?lang=en) if list positions and IDs blur together.

## Familiar image and recall map

A library book may move to another shelf while keeping its inventory number. A task can likewise move inside a `Vec` while retaining its **ID**, a unique record number. The image has a limit: our program issues the number, so it stays stable only if we maintain the counter correctly. Data still disappear at the next run, and numbering starts again.

**Next ID → check overflow → add → advance counter; remove → never issue the old ID again.** `Organizer` stores `tasks` and `next_id`. The next number starts at 1. `checked_add(1)` returns `Some(next number)` or `None` when a `u32` cannot be increased. We reserve `u32::MAX`: reaching it makes `add` return `Err` **before changing** the list. This simplifies the rule but leaves one possible number unused. Removal never decreases the counter. `for index in 0..self.tasks.len()` visits existing positions only, so reading `self.tasks[index]` inside the loop stays in bounds. We search by `task.id`, not by position. `Self` inside `impl Organizer` names the type `Organizer`; `self` is one organizer, as in lesson 28. `u32::MAX` is the largest `u32`. The teaching `main` uses `.expect(...)` only for fixed valid values; actual input errors need handling. Parentheses in the output show task state.

## Counter and lookup

Save and completely replace the previous `src/main.rs`, then run `cargo run` in `organizer`. No new dependencies are needed. Predict the two printed IDs first.

```rust
struct Task {
    id: u32,
    title: String,
    done: bool,
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
            return Err(String::from("Title is empty"));
        }
        match self.next_id.checked_add(1) {
            Some(next) => {
                let id = self.next_id;
                self.tasks.push(Task {
                    id,
                    title: String::from(title),
                    done: false,
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
    fn complete(&mut self, id: u32) -> Result<(), String> {
        for task in &mut self.tasks {
            if task.id == id {
                task.done = true;
                return Ok(());
            }
        }
        Err(String::from("ID not found"))
    }
    fn show(&self) {
        for task in &self.tasks {
            let mark = if task.done { "done" } else { "open" };
            println!("ID {}: {} ({mark})", task.id, task.title);
        }
    }
}

fn main() {
    let mut organizer = Organizer::new();
    let first = organizer.add("Buy bread").expect("fixed valid title");
    let second = organizer.add("Read a book").expect("fixed valid title");
    organizer.remove(first).expect("fixed existing ID");
    let third = organizer.add("Walk").expect("fixed valid title");
    organizer.complete(second).expect("fixed existing ID");
    println!("ID: {second}, {third}");
    organizer.show();
}
```
```text
ID: 2, 3
ID 2: Read a book (done)
ID 3: Walk (open)
```

Two tasks with the same title still receive different IDs: the title is not a lookup key. Removing ID 1 leaves ID 2 unchanged. New ID 3 shows the issue counter is not the list length. Without a file, another run starts at 1 again; saving and restoring the counter comes in the file stage.

## Check your understanding

1. Why does the next task get ID 3 after ID 1 is removed?
2. Why might looking up position 1 after removal change the wrong task?
3. Can `add` partly change the list when the counter overflows?

The counter does not decrease; positions shift; no, the check precedes `push`. Rebuild the recall map.

## Exercise

**Required.** Add a boundary check: try completing unknown ID 99, set `next_id` to `u32::MAX`, and try adding a task. Print both errors and confirm the task count and existing IDs did not change.

## Hint

Handle each `Err(message)` with `match`. For this teaching check, assign `organizer.next_id = u32::MAX` inside `main`; a real user would not be given that access.

## Reference after attempting

<!-- task-answer -->
```rust
struct Task {
    id: u32,
    title: String,
    done: bool,
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
            return Err(String::from("Title is empty"));
        }
        match self.next_id.checked_add(1) {
            Some(next) => {
                let id = self.next_id;
                self.tasks.push(Task {
                    id,
                    title: String::from(title),
                    done: false,
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
    fn complete(&mut self, id: u32) -> Result<(), String> {
        for task in &mut self.tasks {
            if task.id == id {
                task.done = true;
                return Ok(());
            }
        }
        Err(String::from("ID not found"))
    }
    fn show(&self) {
        for task in &self.tasks {
            let mark = if task.done { "done" } else { "open" };
            println!("ID {}: {} ({mark})", task.id, task.title);
        }
    }
}

fn main() {
    let mut organizer = Organizer::new();
    let first = organizer.add("Buy bread").expect("fixed valid title");
    let second = organizer.add("Read a book").expect("fixed valid title");
    organizer.remove(first).expect("fixed existing ID");
    let third = organizer.add("Walk").expect("fixed valid title");
    organizer.complete(second).expect("fixed existing ID");
    println!("ID: {second}, {third}");
    if let Err(message) = organizer.complete(99) {
        println!("Error: {message}");
    }
    organizer.next_id = u32::MAX;
    if let Err(message) = organizer.add("Walk") {
        println!("Error: {message}");
    }
    println!("Task count: {}", organizer.tasks.len());
    organizer.show();
}
```
```text
ID: 2, 3
Error: ID not found
Error: No IDs left
Task count: 2
ID 2: Read a book (done)
ID 3: Walk (open)
```

[Documentation for `u32::checked_add`](https://doc.rust-lang.org/std/primitive.u32.html#method.checked_add).

[Previous lesson](/read/rust-37-operations?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-39-scenarios?lang=en)
