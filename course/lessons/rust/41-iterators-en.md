# Iterators: read, change, or transfer ownership

_Summary:_ **Lesson 41. Choose how to visit tasks without losing access to them too early.**

Prerequisites: [borrowing in lesson 20](/read/rust-20-borrowing?lang=en), [`Vec` in lesson 24](/read/rust-24-vectors?lang=en), and the [organizer in lesson 40](/read/rust-40-memory-checkpoint?lang=en). Save a copy of that project. Here you replace `main.rs` with a small focused program; data from an earlier run do not carry over.

## Familiar image and recall map

A librarian can walk along a shelf to read spines, update cards, or take books away. An **iterator** is an object that yields elements one at a time and remembers its position. The image stops at Rust ownership: `iter()` yields shared `&Task` references, `iter_mut()` yields exclusive `&mut Task` references for changes, and `into_iter()` takes the elements from the vector. After the last operation, the old `Vec` cannot be used. Lessons 19–21 explain these references.

**List → `iter` (read) / `iter_mut` (change) / `into_iter` (take) → next item or end.** One call to `next()` returns `Some` with an item, or `None` at the end. `reader` is mutable because each step changes its position. Merely creating an iterator does not visit items; it waits until asked. A `for` loop keeps asking. The three passes in the example do not overlap: mutation starts after the reading pass ends. `if let Some(task)` runs only when an item exists. `!task.done` means the task is still pending.

## Run and observe

Replace `src/main.rs` and use `cargo run` from [lesson 6](/read/rust-06-cargo?lang=en). No extra dependency is needed. Cargo's own messages are separate from this program output. On a narrow screen, swipe sideways inside the light code box to read a long line; the page itself stays in place.

```rust
struct Task {
    id: u32,
    title: String,
    done: bool,
}

fn main() {
    let mut tasks = vec![
        Task {
            id: 1,
            title: String::from("Read Rust"),
            done: false,
        },
        Task {
            id: 2,
            title: String::from("Write an idea"),
            done: false,
        },
    ];
    let mut reader = tasks.iter();
    if let Some(task) = reader.next() {
        println!("First: {}", task.title);
    }
    for task in tasks.iter_mut() {
        if task.id == 2 {
            task.done = true;
        }
    }
    for task in tasks.iter() {
        println!("Task {}: {}", task.id, task.title);
    }
    let mut done_count = 0;
    for task in tasks.iter() {
        if task.done {
            done_count += 1;
        }
    }
    println!("Done: {}", done_count);
    for task in tasks.into_iter() {
        println!("Owned: {}", task.title);
    }
}
```
```text
First: Read Rust
Task 1: Read Rust
Task 2: Write an idea
Done: 1
Owned: Read Rust
Owned: Write an idea
```

The program first reads a task, then marks task 2 done. Its final loop takes ownership of the titles. Writing `tasks.len()` after `into_iter()` would fail to compile because the vector was moved. For displaying tasks, `iter()` is usually enough. The data still live only for this run.

## Recall without looking

1. Which pass can mark a task done?
2. Why can you not read `tasks` after `into_iter()`?
3. Does `let reader = tasks.iter()` visit any items by itself?

## Exercise

**Required.** Before `into_iter()`, make another pass to count pending tasks and print `Pending: 1`. Preserve the other output lines. Explain why counting does not need ownership of the tasks.

## Answers

Start `pending_count` at zero, walk through `tasks.iter()`, and add one when `!task.done`.

## Answer after your attempt

<!-- task-answer -->
```rust
struct Task {
    id: u32,
    title: String,
    done: bool,
}

fn main() {
    let mut tasks = vec![
        Task {
            id: 1,
            title: String::from("Read Rust"),
            done: false,
        },
        Task {
            id: 2,
            title: String::from("Write an idea"),
            done: false,
        },
    ];
    let mut reader = tasks.iter();
    if let Some(task) = reader.next() {
        println!("First: {}", task.title);
    }
    for task in tasks.iter_mut() {
        if task.id == 2 {
            task.done = true;
        }
    }
    for task in tasks.iter() {
        println!("Task {}: {}", task.id, task.title);
    }
    let mut done_count = 0;
    for task in tasks.iter() {
        if task.done {
            done_count += 1;
        }
    }
    println!("Done: {}", done_count);
    let mut pending_count = 0;
    for task in tasks.iter() {
        if !task.done {
            pending_count += 1;
        }
    }
    println!("Pending: {}", pending_count);
    for task in tasks.into_iter() {
        println!("Owned: {}", task.title);
    }
}
```
```text
First: Read Rust
Task 1: Read Rust
Task 2: Write an idea
Done: 1
Pending: 1
Owned: Read Rust
Owned: Write an idea
```

## After checking

Check that counting happens before ownership moves and leaves IDs and status unchanged. [Official iterator chapter](https://doc.rust-lang.org/book/ch13-02-iterators.html). If references are unclear, revisit [lesson 20](/read/rust-20-borrowing?lang=en).

[Previous lesson](/read/rust-40-memory-checkpoint?lang=en) · [Contents](/course/rust?lang=en)
