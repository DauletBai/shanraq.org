# Tasks in memory: add, rename, complete, remove

_Lead (summary):_ **Lesson 37. Carry out four task actions and see how removal shifts list positions.**

Prerequisites: lessons 21, 24, 27, 31–33 and 36. Review [lesson 21](/read/rust-21-mutable-borrow?lang=en) if mutable borrowing is unclear.

## Familiar image and recall map

Picture task cards laid out in a row. You can add a card, rewrite its title, mark it complete, or take it away. `Vec<Task>` is such a row in program memory. These cards still disappear when the program closes; files come later. The image has a limit: a place in the row is not a permanent card number.

**Check → change `Vec<Task>` → show result; close the program → in-memory data disappear.** `Task` groups `title: String` and `done: bool`. The `done` field becomes true after completion. `&mut Vec<Task>` lends a function exclusive access to change the list; `&[Task]` lends read-only access to a slice of it. `get_mut(index)` returns `Option<&mut Task>` because a position may be absent. `tasks.remove(index)` removes an element and shifts later ones left; we check the boundary first because `remove` itself would stop the program on a bad index. `usize` is the list-position type from lesson 15. `Result<(), String>` means a successful action returns no value (`()`), while an error carries a reason. `for task in tasks` reads the cards in order.

## Four actions

Save and replace the previous `src/main.rs`. Run `cargo run` in `organizer`. No new dependencies are needed. Predict which card remains after position 0 is removed. `[ ]` marks a pending task and `[done]` a completed one. In this teaching `main`, `.expect(...)` extracts a successful `Ok` but stops the program on `Err`; use it only for fixed, known-good data, never raw user input.

```rust
struct Task {
    title: String,
    done: bool,
}

fn add(tasks: &mut Vec<Task>, title: &str) -> Result<(), String> {
    if title.trim().is_empty() {
        return Err(String::from("Title is empty"));
    }
    tasks.push(Task {
        title: String::from(title),
        done: false,
    });
    Ok(())
}

fn rename(tasks: &mut Vec<Task>, index: usize, title: &str) -> Result<(), String> {
    match tasks.get_mut(index) {
        Some(task) => {
            task.title = String::from(title);
            Ok(())
        }
        None => Err(String::from("No such position")),
    }
}

fn complete(tasks: &mut Vec<Task>, index: usize) -> Result<(), String> {
    match tasks.get_mut(index) {
        Some(task) => {
            task.done = true;
            Ok(())
        }
        None => Err(String::from("No such position")),
    }
}

fn remove(tasks: &mut Vec<Task>, index: usize) -> Result<(), String> {
    if index >= tasks.len() {
        return Err(String::from("No such position"));
    }
    tasks.remove(index);
    Ok(())
}

fn show(tasks: &[Task]) {
    for task in tasks {
        if task.done {
            println!("[done] {}", task.title);
        } else {
            println!("[ ] {}", task.title);
        }
    }
}

fn main() {
    let mut tasks = Vec::new();
    add(&mut tasks, "Buy bread").expect("valid title");
    add(&mut tasks, "Read a chapter").expect("valid title");
    rename(&mut tasks, 1, "Read two chapters").expect("existing position");
    complete(&mut tasks, 0).expect("existing position");
    remove(&mut tasks, 0).expect("existing position");
    show(&tasks);
}
```
```text
[ ] Read two chapters
```

The teaching `main` uses `.expect(...)` only because its positions and nonempty titles are fixed in advance. This `Result` method extracts `Ok` but stops the program on `Err`; it is unsuitable for real user input. The exercise handles errors explicitly. After removal, the second element moves to position 0. That motivates a stable ID in the next lesson.

## Check your understanding

1. Why does `get_mut` return `Option` rather than `&mut Task` directly?
2. What happens to the second card position when the first is removed?
3. Will these tasks survive another run?

The position may be absent; the second card moves to position 0; no, there is no saving yet. Rebuild the recall map.

## Exercise

**Required.** Extend `rename` to reject a title made only of spaces without changing the old title. After successful actions, try renaming position 0 to `"   "` and completing position 9. Print both errors, then the list; the old title must remain.

## Hint

Put `title.trim().is_empty()` before `get_mut`. Print each `Result` error with `match`, as in lesson 32.

## Reference after attempting

<!-- task-answer -->
```rust
struct Task {
    title: String,
    done: bool,
}

fn add(tasks: &mut Vec<Task>, title: &str) -> Result<(), String> {
    if title.trim().is_empty() {
        return Err(String::from("Title is empty"));
    }
    tasks.push(Task {
        title: String::from(title),
        done: false,
    });
    Ok(())
}

fn rename(tasks: &mut Vec<Task>, index: usize, title: &str) -> Result<(), String> {
    if title.trim().is_empty() {
        return Err(String::from("Title is empty"));
    }
    match tasks.get_mut(index) {
        Some(task) => {
            task.title = String::from(title);
            Ok(())
        }
        None => Err(String::from("No such position")),
    }
}

fn complete(tasks: &mut Vec<Task>, index: usize) -> Result<(), String> {
    match tasks.get_mut(index) {
        Some(task) => {
            task.done = true;
            Ok(())
        }
        None => Err(String::from("No such position")),
    }
}

fn remove(tasks: &mut Vec<Task>, index: usize) -> Result<(), String> {
    if index >= tasks.len() {
        return Err(String::from("No such position"));
    }
    tasks.remove(index);
    Ok(())
}

fn show(tasks: &[Task]) {
    for task in tasks {
        if task.done {
            println!("[done] {}", task.title);
        } else {
            println!("[ ] {}", task.title);
        }
    }
}

fn main() {
    let mut tasks = Vec::new();
    add(&mut tasks, "Buy bread").expect("valid title");
    add(&mut tasks, "Read a chapter").expect("valid title");
    rename(&mut tasks, 1, "Read two chapters").expect("existing position");
    complete(&mut tasks, 0).expect("existing position");
    remove(&mut tasks, 0).expect("existing position");
    for result in [rename(&mut tasks, 0, "   "), complete(&mut tasks, 9)] {
        if let Err(message) = result {
            println!("Error: {message}");
        }
    }
    show(&tasks);
}
```
```text
Error: Title is empty
Error: No such position
[ ] Read two chapters
```

[Documentation for `Vec::get_mut` and `Vec::remove`](https://doc.rust-lang.org/std/vec/struct.Vec.html#method.remove).

[Previous lesson](/read/rust-36-commands?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-38-ids?lang=en)
