# String and &str: storing and reading

_Lead (summary):_ **Lesson 22. Choose an owned string for storage and a text view for reading.**

Prerequisites: lessons 9 and 18–21. If moving a `String` is unclear, revisit [lesson 19](/read/rust-19-ownership?lang=en).

## A familiar image and recall map

Your own notebook stores your notes; an open page lets someone read a note but does not become a second notebook. **`String`** owns its text and can grow when change is allowed. **`&str`** is a read-only view of a contiguous part of existing text with a known length, without owning its storage. The image has a limit: a text literal is also `&str`, although its data lives in the program and need not belong to a `String` variable.

**Store and change → `String`; accept for reading → `&str`; owner stays alive → view remains valid.** Recall this map without looking. `"Rest"` is a `&str` literal; `String::from("Reading")` makes a separate owned string. The `.as_str()` method gives a read-only `&str` view of a whole `String` without copying its text. The dot calls a method, as `.len()` did in lesson 15.

## One parameter for two text sources

In `organizer`, save your previous `src/main.rs` separately, replace the whole file with the following program, and run `cargo run` beside `Cargo.toml`. No other files or dependencies change. Predict three lines; only program output appears below.

Before running, read the calls: `show_title(owned.as_str())` passes a view from an owner, while `show_title(literal)` passes an existing `&str` literal. The parameter `title: &str` reads text in both cases without taking the `String`.

```rust
fn show_title(title: &str) {
    println!("Title: {title}");
}

fn main() {
    let owned = String::from("Reading");
    let literal: &str = "Rest";
    show_title(owned.as_str());
    show_title(literal);
    println!("Stored: {owned}");
}
```

```text
Title: Reading
Title: Rest
Stored: Reading
```

`owned` remains usable: the function did not take ownership. A literal appears directly in source and remains available while the program runs. A view from `owned.as_str()` must not be used after `owned` disappears; the compiler prevents keeping a dangling reference.

## Choosing a parameter type

If a function only reads a title, `&str` usually accepts more text sources than `&String`: literals, whole `String` values through `.as_str()`, and parts of strings from lesson 23. This does not mean `&str` can itself store a new changeable record. If the organizer must keep a title after the call, it needs its own `String`; we will cover passing and storing it when we build a task type.

You cannot call `push_str` through `&str`: it is a view for reading. This separate invalid example yields E0599 because the type has no such method:

<!-- error-code: E0599 -->
```rust,compile_fail
fn main() {
    let owned = String::from("Plan");
    let view = owned.as_str();
    view.push_str("!");
}
```

Changing the original string needs `let mut owned` and the suitable mutable access from lesson 21. Do not change through an old reading view: growth could invalidate it. A `&str` can also describe part of a string rather than all of it.

## Check your understanding

1. Who owns the text bytes after `let owned = String::from("Reading")`?
2. Does `owned.as_str()` return a copy or a view?
3. Why does `show_title(literal)` need no `String::from`?

Check: `owned`; a view without copying text; a literal already has type `&str`. Explain the notebook image and where it stops being exact.

## Exercise

**Required.** Write `show_task(title: &str, minutes: u32)`. In `main`, create an owned `Reading` string and a `Rest` literal; call the function with `owned.as_str(), 15` and `literal, 5`. Then print `Stored: Reading`. Expect three lines in that order. Do not copy the owned string just to display it.

**Your own data.** Change titles and minutes while keeping one call from a `String` and another from a literal.

## Hint

Inside the function use `println!("{title}: {minutes} min");`. The first call needs `.as_str()`; the second does not. `owned` stays in `main`.

## Reference answer after your attempt

<!-- task-answer -->
```rust
fn show_task(title: &str, minutes: u32) {
    println!("{title}: {minutes} min");
}

fn main() {
    let owned = String::from("Reading");
    let literal: &str = "Rest";
    show_task(owned.as_str(), 15);
    show_task(literal, 5);
    println!("Stored: {owned}");
}
```

```text
Reading: 15 min
Rest: 5 min
Stored: Reading
```

The function reads both text sources while `main` keeps ownership of its string. [Official guide to strings](https://doc.rust-lang.org/book/ch08-02-strings.html).

[Previous lesson](/read/rust-21-mutable-borrow?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-23-slices?lang=en)
