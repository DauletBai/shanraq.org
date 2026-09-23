# The Task struct: data with useful names

_Lead (summary):_ **Lesson 27. Put a task's title, duration and completion state into one type of your own.**

Before this lesson: lessons 9, 11, 19, 22 and 26. If tuples still feel unfamiliar, revisit [lesson 26](/read/rust-26-tuples?lang=en).

## Familiar image and recall map

A task card has labelled spaces for “title”, “minutes” and “done”. With a tuple, you would need to remember what `.0` and `.1` mean. A **struct** (`struct` is short for structure) is a type you define with named **fields**. A field stores one part of the value; an **instance** is one filled-in card. The image has a limit: Rust checks that the named fields exist and have the right types, unlike paper. But the struct alone cannot check a business rule such as minutes being greater than zero.

**`Task` type → fields → instance → read with a dot → change with `mut`.** Say the map before reading code. The keyword `struct` declares a type. Braces after `Task` contain fields written as `name: type`, separated by commas. This describes a form; it does not create a task yet. By Rust convention, type names such as `Task` begin with a capital letter.

## Fill in the card

Save the old `src/main.rs` in `organizer`, replace the whole file, and run `cargo run` beside `Cargo.toml`. No other file or dependency changes. Predict the output first.

`Task { title: ..., minutes: ..., done: ... }` creates an **instance**. Unlike a tuple, every value has a field name. `task.title` reads the field after the dot. `println!` uses the text while formatting, and the task remains available afterwards. `done` has the familiar `bool` type. A comma after the final field is allowed and makes another field easier to add.

```rust
struct Task {
    title: String,
    minutes: u32,
    done: bool,
}

fn main() {
    let task = Task {
        title: String::from("Reading"),
        minutes: 20,
        done: false,
    };
    println!("Task: {}", task.title);
    println!("Minutes: {}", task.minutes);
    println!("Done: {}", task.done);
}
```
```text
Task: Reading
Minutes: 20
Done: false
```

Fields can appear in a different order when creating an instance: their names preserve the meaning.

## Change a field and check its meaning

`let mut task` permits changes to this instance's fields. `task.done = true` writes a new Boolean value. The dot accesses a field here; it does not call a method. You do not put `mut` on the struct definition. `task.minutes > 0` is our organizer's rule, not a guarantee of `struct`: zero has the right `u32` type but may be the wrong duration for this task.

```rust
struct Task {
    title: String,
    minutes: u32,
    done: bool,
}

fn main() {
    let mut task = Task {
        title: String::from("Walk"),
        minutes: 30,
        done: false,
    };
    if task.minutes > 0 {
        task.done = true;
    }
    println!("{}: {}", task.title, task.done);
}
```

```text
Walk: true
```

This condition shows where a check belongs. Later we will introduce a way to prevent an unsuitable task from being created and report why. There is no user input here, so this example is not a complete data validation system.

Every declared field must be filled. This **separate** example omits `done` and intentionally fails with E0063:

<!-- error-code: E0063 -->
```rust,compile_fail
struct Task {
    title: String,
    minutes: u32,
    done: bool,
}

fn main() {
    let task = Task {
        title: String::from("Reading"),
        minutes: 20,
    };
    println!("{}", task.title);
}
```

The compiler checks fields and types, but cannot know how many minutes make sense for a walk. That is where the paper form image ends.

## Check your understanding

1. How does `Task` differ from one `task`?
2. Why does `task.minutes` explain more than `report.1`?
3. Does `minutes: u32` guarantee a positive number?

Check: one is a type definition, the other an instance; a field name states its role; no, zero is also a `u32`. Rebuild the recall map from memory.

## Exercise

**Required.** Declare `Task` with the three fields from this lesson. Create a mutable task named `Cleaning` with `15` minutes and `false`. Print `Cleaning: 15 min, done false`, change `done` to `true`, then print `After: true`. Use field values rather than finished output strings.

**Boundary check.** Try creating an instance without `minutes`. Read the compiler message and restore the field.

## Hint

Start with `let mut task = Task { ... };`, then use `task.done = true;`. Pass the fields to `println!`.

## Reference answer after trying

<!-- task-answer -->
```rust
struct Task {
    title: String,
    minutes: u32,
    done: bool,
}

fn main() {
    let mut task = Task {
        title: String::from("Cleaning"),
        minutes: 15,
        done: false,
    };
    println!("{}: {} min, done {}", task.title, task.minutes, task.done);
    task.done = true;
    println!("After: {}", task.done);
}
```
```text
Cleaning: 15 min, done false
After: true
```

Task data now stays together and can be read by name. The next lesson puts related behavior beside it. [Official struct explanation](https://doc.rust-lang.org/book/ch05-01-defining-structs.html).

[Previous lesson](/read/rust-26-tuples?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-28-methods?lang=en)
