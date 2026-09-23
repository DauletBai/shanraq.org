# Enums: one valid task state at a time

_Lead (summary):_ **Lesson 29. Represent a task's state with one named variant that may carry extra data.**

Before this lesson: lessons 11, 12, 14 and 27–28. If struct fields are unfamiliar, revisit [lesson 27](/read/rust-27-structs?lang=en).

## Familiar image and recall map

A sign on a door shows one state: “not started”, “in progress” or “finished”. With three separate `bool` fields, we could accidentally mark a task both in progress and finished. An **enum** defines a type with named **variants**; one value selects exactly one variant at a time. The sign image has a limit: some variants can also carry data, which the simple sign does not show.

**`Status` type → variants → `Status::Doing(15)` value → inspect the variant.** `enum Status { ... }` declares the possible states. `Status::Planned` and `Status::Done` carry no extra data. `Doing(u32)` means the “in progress” variant holds one number of minutes. Parentheses in the declaration give the data type; in `Status::Doing(15)` they supply the value. `::` gives the path from the type to its variant. It is not a separate function.

## One state at a time

Save the old `src/main.rs` in `organizer`, replace the whole file, and run `cargo run` beside `Cargo.toml`. No other file or dependency changes. Predict the output without Cargo's messages.

To show the value inside a variant, choose a branch. `match status { ... }` compares the value with variant **patterns**. `Status::Doing(minutes) => ...` means: if the variant is `Doing`, name its number `minutes` and run the expression after `=>`. The `=>` separates a pattern from its action; commas separate branches. `Planned` and `Done` are also present, so every possible state is handled. The next lesson gives `match` more practice.

```rust
enum Status {
    Planned,
    Doing(u32),
    Done,
}

fn main() {
    let status = Status::Doing(15);
    match status {
        Status::Planned => println!("Not started"),
        Status::Doing(minutes) => println!("In progress: {minutes} min"),
        Status::Done => println!("Finished"),
    }
}
```
```text
In progress: 15 min
```

Here `15` means minutes already spent, not the task's planned duration. `Done` cannot also hold `Doing` data: a value has one variant. Rust does allow `Doing(0)`; checking whether that number makes sense is our job.

## An enum inside Task

A struct field can have our new type. This separate example creates a task and then changes only its state. Its title stays the same.

```rust
enum Status {
    Planned,
    Doing(u32),
    Done,
}

struct Task {
    title: String,
    status: Status,
}

fn main() {
    let mut task = Task {
        title: String::from("Reading"),
        status: Status::Planned,
    };
    task.status = Status::Done;
    println!("{} finished", task.title);
}
```

```text
Reading finished
```

The message is correct after this known assignment. If the state can vary, inspect `task.status` with `match` instead of printing a fixed message. We will do that in lesson 30.

Variant data has a type too. This **separate** example intentionally fails with E0308: `Doing` needs a `u32`, not text.

<!-- error-code: E0308 -->
```rust,compile_fail
enum Status { Doing(u32) }
fn main() {
    let status = Status::Doing("15");
}
```

## Check your understanding

1. Why is `Status` a better fit than three separate `bool` fields for mutually exclusive states?
2. What does `Status::Doing(15)` hold besides the variant name?
3. Does the enum check that its number is positive?

Check: a value selects only one variant; `Doing` holds a `u32` equal to 15; no, that is a separate organizer rule. Close the page and rebuild the recall map.

## Exercise

**Required.** Declare `Status` with `Planned`, `Doing(u32)` and `Done`. Create `Status::Doing(25)` and use `match` to print `In progress: 25 min`. Give the other two variants `Not started` and `Finished` messages so the code handles any state. Take the number in the output from the variant.

**Change.** Set the value to `Status::Done` and check that it prints `Finished`.

## Hint

Inside `match status`, use three branches. The `Status::Doing(minutes)` branch makes `minutes` available to `println!` after `=>`.

## Reference answer after trying

<!-- task-answer -->
```rust
enum Status {
    Planned,
    Doing(u32),
    Done,
}

fn main() {
    let status = Status::Doing(25);
    match status {
        Status::Planned => println!("Not started"),
        Status::Doing(minutes) => println!("In progress: {minutes} min"),
        Status::Done => println!("Finished"),
    }
}
```
```text
In progress: 25 min
```

The state is one value, and the “in progress” data belongs to its specific variant. [Official enum explanation](https://doc.rust-lang.org/book/ch06-01-defining-an-enum.html).

[Previous lesson](/read/rust-28-methods?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-30-patterns?lang=en)
