# Borrowing for reading

_Lead (summary):_ **Lesson 20. Let a function read a string while its owner keeps using it.**

Prerequisites: lessons 1–19. If a move is unclear, revisit [lesson 19](/read/rust-19-ownership?lang=en). Changing data through a borrowed reference comes in lesson 21.

## A picture and the new symbol

If a friend needs to read your book, you can lend it instead of giving it away or buying another copy. **Borrowing** gives temporary access to a value without transferring ownership. The picture's limit: Rust checks when references are used and which kinds of access may coexist; an ordinary loan agreement does not explain those rules.

The symbol `&` makes a **reference**, a value that points to another value without owning it. In a call, `&title` creates a reference to `title`. In a parameter declaration, `&String` says the function accepts a read-only reference to a `String`. The same `&` appears in two places for two related purposes: create a reference, then describe its type. A reference does not own the string's text buffer.

**Recall map:** owner → `&` → reading reference → function → owner continues. Cover the line and explain it before running code.

## Read without a move

In the `organizer` project from lesson 6, save the former `src/main.rs`, replace the whole file, and run `cargo run` in the folder containing `Cargo.toml`. No other files or dependencies change. Predict both lines; the output block contains program output only.

```rust
fn show_title(title: &String) {
    println!("Inside: {title}");
}

fn main() {
    let title = String::from("Reading");
    show_title(&title);
    println!("After the call: {title}");
}
```

```text
Inside: Reading
After the call: Reading
```

The call creates a reference, and the parameter receives it. The owner in `main` still owns the `String`, so `main` can use `title` afterward. This is **immutable borrowing**: the function can read through its reference but cannot change the text. “Immutable” describes the access through this reference, even if the owner might later change its value.

## Where reading access stops

`push_str` changes a `String`. Trying to call it through a read-only `&String` fails with E0596:

<!-- error-code: E0596 -->
```rust,compile_fail
fn change(title: &String) {
    title.push_str("!");
}

fn main() {
    let title = String::from("Reading");
    change(&title);
}
```

Several read-only references can coexist. While a reference is still needed, changing the value in a way that would invalidate the read is not permitted. Rust can end a borrow after its **last use**, which may be earlier than the next closing brace. Predict whether this separate example compiles:

```rust
fn main() {
    let mut title = String::from("Plan");
    let view = &title;
    println!("Before: {view}");
    title.push_str("!");
    println!("After: {title}");
}
```

```text
Before: Plan
After: Plan!
```

It compiles because the last use of `view` happens before `push_str`. If you tried to print `view` again after the change, Rust would reject the program. We will examine slices and the connection between reference and owner in lessons 23 and 47.

## Recall check

Why does `main` retain `title` here, unlike the `String`-by-value call in lesson 19? Can `&String` call `push_str`? Answers: the reference does not move ownership; no, changing through this reference is forbidden.

## Exercise

**Required.** Write `show_task(title: &String, minutes: u32)` to print `Reading — 15 min`. In `main`, create the string `Reading`, call the function with 15 and 20, and then print `Title remains: Reading`. Predict the three lines first. Do not use `clone()` or pass a `String` by value. If the call fails to compile, check for `&` in both the parameter type and the call.

## Hint

`u32` is passed by value because it is `Copy`; pass the string as `&title`. Inside the function use `println!("{title} — {minutes} min");`.

## Reference answer after your attempt

<!-- task-answer -->
```rust
fn show_task(title: &String, minutes: u32) {
    println!("{title} — {minutes} min");
}

fn main() {
    let title = String::from("Reading");
    show_task(&title, 15);
    show_task(&title, 20);
    println!("Title remains: {title}");
}
```

```text
Reading — 15 min
Reading — 20 min
Title remains: Reading
```

Explain the map from memory once more: the owner stays in `main` and the function receives temporary reading access. [Rust's reference chapter](https://doc.rust-lang.org/book/ch04-02-references-and-borrowing.html).

[Previous lesson](/read/rust-19-ownership?lang=en) · [Contents](/course/rust?lang=en)
