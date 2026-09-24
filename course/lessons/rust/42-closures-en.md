# Closures: select the tasks you need

_Summary:_ **Lesson 42. Search task titles and show matches in a stable order.**

Prerequisites: [iterators](/read/rust-41-iterators?lang=en), [strings](/read/rust-22-strings?lang=en), and [stable IDs](/read/rust-38-ids?lang=en). Save the previous `main.rs`: this focused example replaces it and still keeps data only in memory.

## Familiar image and recall map

Imagine giving a library assistant the rule “keep cards whose title contains Rust.” A **closure** is an unnamed function that can remember a value from its surroundings. In `|task| task.title.contains(query)`, the vertical bars enclose the input `task`, while `query` comes from outside. The image has a limit: giving the rule does not do all the work at once. The lazy iterator runs it when results are requested.

**All tasks → `filter` (keep matches) → `collect` (gather) → `sort_by_key` (order) → display.** `filter` uses a closure that returns a condition. `iter()` yields `&Task`; `filter` borrows that item to test it, so its closure receives `&&Task`, a reference to a reference. Rust follows both layers automatically for `.title`; the same happens for `.id` in sorting and `.title` in `map`. `contains` matches exact characters, including letter case: `rust` differs from `Rust`. Here `collect` builds `Vec<&Task>`, a new list of references; the original tasks stay where they are. The colon in `let found: Vec<&Task>` tells Rust the desired result type. `sort_by_key(|task| task.id)` orders by ID rather than title. `map` turns each found task into a `&str` view of its title; `as_str()` borrows text without copying. `names.len()` counts matches. Here `names` exists to show `map`; if you only need the count, `found.len()` avoids building another list. The vertical bars here are Rust closure syntax, not the separators used in the lesson 40 command prompt.

## Run the search

Replace `src/main.rs` and run `cargo run`. Program output is shown apart from Cargo messages. No new dependency is needed.

```rust
struct Task {
    id: u32,
    title: String,
}

fn main() {
    let tasks = vec![
        Task {
            id: 3,
            title: String::from("Read Rust"),
        },
        Task {
            id: 2,
            title: String::from("Buy bread"),
        },
        Task {
            id: 1,
            title: String::from("Review Rust"),
        },
    ];
    let query = "Rust";
    let mut found: Vec<&Task> = tasks
        .iter()
        .filter(|task| task.title.contains(query))
        .collect();
    found.sort_by_key(|task| task.id);
    let names: Vec<&str> = found.iter().map(|task| task.title.as_str()).collect();
    println!("Found: {}", names.len());
    for task in found {
        println!("{}: {}", task.id, task.title);
    }
}
```
```text
Found: 2
1: Review Rust
3: Read Rust
```

Task 3 appears before task 1 in the original list, but the matches are displayed by ID. The original `tasks` stays unchanged. If there is no match, `found` is empty and the loop prints nothing; that is a valid result. Each search scans the list again; speed is not the point of this lesson.

## Recall without looking

1. Where does the closure get `query`?
2. Why does `collect` leave the original tasks in place?
3. Does `sort_by_key` change the tasks' stable IDs?

## Exercise

**Required.** After the first output, count titles containing `bread` with a new `filter` and `count`; print `Found bread: 1`. Keep the original list untouched. Then try `Bread` and explain the different result.

## Answers

Make `second_query`. `tasks.iter().filter(|task| task.title.contains(second_query)).count()` consumes that search and returns a number.

## Answer after your attempt

<!-- task-answer -->
```rust
struct Task {
    id: u32,
    title: String,
}

fn main() {
    let tasks = vec![
        Task {
            id: 3,
            title: String::from("Read Rust"),
        },
        Task {
            id: 2,
            title: String::from("Buy bread"),
        },
        Task {
            id: 1,
            title: String::from("Review Rust"),
        },
    ];
    let query = "Rust";
    let mut found: Vec<&Task> = tasks
        .iter()
        .filter(|task| task.title.contains(query))
        .collect();
    found.sort_by_key(|task| task.id);
    let names: Vec<&str> = found.iter().map(|task| task.title.as_str()).collect();
    println!("Found: {}", names.len());
    for task in found {
        println!("{}: {}", task.id, task.title);
    }
    let second_query = "bread";
    let second_count = tasks
        .iter()
        .filter(|task| task.title.contains(second_query))
        .count();
    println!("Found {second_query}: {second_count}");
}
```
```text
Found: 2
1: Review Rust
3: Read Rust
Found bread: 1
```

## After checking

`sort_by_key` orders the gathered references, not the IDs inside tasks. [Official closures chapter](https://doc.rust-lang.org/book/ch13-01-closures.html) · [iterator chapter](https://doc.rust-lang.org/book/ch13-02-iterators.html). Revisit [lesson 41](/read/rust-41-iterators?lang=en) if delayed iteration is unclear.

[Previous lesson](/read/rust-41-iterators?lang=en) · [Contents](/course/rust?lang=en)
