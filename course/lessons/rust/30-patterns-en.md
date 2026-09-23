# match: handle every variant

_Lead (summary):_ **Lesson 30. Read a task state fully and catch a missing case before running the program.**

Before this lesson: lessons 12, 14, 19 and 29. If `Status::Doing(15)` is unclear, revisit [lesson 29](/read/rust-29-enums?lang=en).

## Familiar image and recall map

A clerk receives a card with one of three stamps. To know what to do every time, the desk has an instruction for each stamp. **Pattern matching** checks a value's form and, when it contains data, gives that data names. `match` picks one matching branch. The image has a limit: a clerk can forget an instruction, but Rust checks at compile time that all possible variants are covered.

**Value → `match` → patterns → one branch → result.** Branches are considered from top to bottom. `Status::Doing(minutes)` is a pattern for `Doing` that gives its inner number the new name `minutes`. `=>` separates a pattern from its result or action. Unlike a series of independent `if` conditions, `match` must cover every possible variant or the program will not compile. `_` can stand for “any remaining case”, but it also catches a new state added later.

## A result from complete matching

Save the old `src/main.rs` in `organizer`, replace the whole file, and run `cargo run` beside `Cargo.toml`. No other files or dependencies are needed. Predict the output first.

`match` is an **expression**: its result can be assigned to a variable. Here it answers the narrow question “how many minutes are currently counted as active work?” For `Planned` and `Done`, the answer is 0 because work is not in progress now. For `Doing(minutes)`, it is the number inside that variant. All three branches produce a `u32`, so the entire expression has one type. The `;` after the closing brace completes the `let` statement.

```rust
enum Status {
    Planned,
    Doing(u32),
    Done,
}

fn main() {
    let status = Status::Doing(10);
    let active_minutes: u32 = match status {
        Status::Planned => 0,
        Status::Doing(minutes) => minutes,
        Status::Done => 0,
    };
    println!("Currently active: {active_minutes} min");
}
```
```text
Currently active: 10 min
```

The branches do not all run: one is chosen for the value. Replacing `Doing(10)` with `Done` gives 0. This does not claim that a finished task took no time; it counts only work active **now**.

## A missing branch and `_`

Removing `Done` makes the compiler reject this separate example with E0004. That state cannot be left without a result:

<!-- error-code: E0004 -->
```rust,compile_fail
enum Status { Planned, Doing(u32), Done }
fn main() {
    let status = Status::Done;
    let minutes = match status {
        Status::Planned => 0,
        Status::Doing(value) => value,
    };
    println!("{minutes}");
}
```

`_` is a pattern that matches any value not handled earlier. For example, `Status::Done => ...` and `_ => ...` together cover every state. But if a new variant is added later, `_` hides it from the completeness check. Name variants separately when each organizer state has its own meaning. `_` helps when the remaining cases intentionally share one action.

## When only one variant matters

`if let` combines lesson 12's conditional branch with one pattern. `if let Status::Doing(minutes) = status` does not assign a new value to `status` or compare two numbers. It asks whether the value has the specified form and, if so, names its number `minutes`. `else` handles the other variants. This is a separate program; its result does not depend on the previous one.

```rust
enum Status { Planned, Doing(u32), Done }

fn main() {
    let status = Status::Doing(5);
    if let Status::Doing(minutes) = status {
        println!("Work in progress: {minutes} min");
    } else {
        println!("No work in progress now");
    }
}
```

```text
Work in progress: 5 min
```

Use a full `match` when all variants need different actions. Use `if let` when one variant matters most; add `else` when other cases also require an action.

## Checkpoint for lessons 26–30

Explain without the page: a tuple groups positions in order; a struct names the parts; a method connects an action to a type; an enum limits the variants; `match` checks their complete handling. If a step is missing, return to its lesson from 26–29 and rebuild its recall map.

## Exercise

**Required.** Create `Status` with `Planned`, `Doing(u32)` and `Done`. Write `active_minutes(status: Status) -> u32` using a complete `match`: 0 for `Planned` and `Done`, the inner number for `Doing`. In `main`, call it with `Doing(25)` and `Done`, then print `Active: 25 min` and `After finishing: 0 min`. Do not bypass the `match` by putting finished numbers in `println!`.

**Boundary check.** Remove the `Done` branch, read E0004, and restore it.

## Hint

Make `match status` the function's final expression without `;`. The `Status::Doing(minutes)` branch can simply return `minutes`.

## Reference answer after trying

<!-- task-answer -->
```rust
enum Status {
    Planned,
    Doing(u32),
    Done,
}

fn active_minutes(status: Status) -> u32 {
    match status {
        Status::Planned => 0,
        Status::Doing(minutes) => minutes,
        Status::Done => 0,
    }
}

fn main() {
    let during = active_minutes(Status::Doing(25));
    let after = active_minutes(Status::Done);
    println!("Active: {during} min");
    println!("After finishing: {after} min");
}
```
```text
Active: 25 min
After finishing: 0 min
```

The function now has a defined result for every state. [Official `match` chapter](https://doc.rust-lang.org/book/ch06-02-match.html) and [`if let` explanation](https://doc.rust-lang.org/book/ch06-03-if-let.html).

[Previous lesson](/read/rust-29-enums?lang=en) · [Contents](/course/rust?lang=en)
