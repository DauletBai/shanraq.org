# Ownership and moves

_Lead (summary):_ **Lesson 19. Explain why an old name becomes unavailable after passing a String value.**

Prerequisites: lessons 1–18. Revisit [scope and memory](/read/rust-18-memory?lang=en) if you cannot tell where a name may be used.

## The rule with a picture

A `String` value has an **owner**, the part of the program responsible for the string's managed data. Giving that value to another name normally **moves ownership**. Imagine one claim ticket for an item in storage: after handing the ticket over, the old holder cannot claim the item. The picture has a limit: Rust does not promise to move the text bytes to a new physical address. “Move” describes which name may use the value.

**Recall map:** create → own → move → old name unavailable → new owner eventually releases the resource. Repeat it before reading the program.

## Observe a move

Use `organizer` from lesson 6. Save the old `src/main.rs`, replace the whole file, and run `cargo run` in the folder containing `Cargo.toml`. No dependencies change. The output block is program output, without Cargo messages.

```rust
fn print_title(title: String) {
    println!("Task: {title}");
}

fn main() {
    let title = String::from("Reading");
    print_title(title);
}
```

```text
Task: Reading
```

The parameter `title: String` receives ownership at the call. The function can use it. When the function finishes, its owned string resource is released automatically. The old `title` name in `main` cannot be used after the call. This attempt fails with E0382:

<!-- error-code: E0382 -->
```rust,compile_fail
fn print_title(title: String) {
    println!("{title}");
}

fn main() {
    let title = String::from("Reading");
    print_title(title);
    println!("{title}");
}
```

The compiler is not saying the string became empty. It prevents using a name whose value has moved and may already have been released. Predict which line the compiler objects to before trying the example yourself.

## Why numbers behave differently

`u32` has the **`Copy` trait**. A trait describes a property or behavior of a type; we will define our own in lesson 46. For `Copy`, passing a value produces an implicit independent copy, so the earlier name remains usable. The next program uses only syntax already seen:

```rust
fn print_minutes(minutes: u32) {
    println!("Inside: {minutes}");
}

fn main() {
    let minutes = 15;
    print_minutes(minutes);
    println!("Outside: {minutes}");
}
```

```text
Inside: 15
Outside: 15
```

Arrays whose elements are `Copy` can also be `Copy`. This does not mean that strings or every collection are copied. Do not generalize from a small numeric example.

If two independent owners of the text really are needed, `clone()` explicitly creates another `String` with copied text data. The dot calls a method, as in lesson 15; `clone` means “make an independent copy.” This can allocate memory, so it is a deliberate choice:

```rust
fn main() {
    let first = String::from("Reading");
    let second = first.clone();
    println!("First: {first}");
    println!("Second: {second}");
}
```

```text
First: Reading
Second: Reading
```

Cloning is useful for two owners. When a function only needs to read the text, the next lesson will pass a reference without copying the string.

## Recall check

After `let b = a;`, which names remain usable if `a` is a `String`? Which if it is a `u32`? Answer: only `b` for `String`; both names for `u32`, with separate numeric values. Explain this with the claim-ticket picture and its limit.

## Exercise

Before the exercise, one more use of familiar `mut`: in `fn decorate(mut title: String)`, it belongs to the parameter name inside the function. The function owns the value and may add text to it. It does not restore access through the old name in `main` after the move.

**Required.** Write `decorate(mut title: String) -> String`. It adds `!` using `push_str` and returns the string. In `main`, pass `String::from("Plan")`, store the returned value, and print `Plan!`. Explain ownership before the call, inside the function, and after the return. Do not use `clone()`: passing the same value there and back is enough.

## Hint

Write the parameter as `mut title: String`: the `mut` from lesson 9 now lets the function change its own parameter binding. Make `title` the final expression without `;`, then assign the returned value to a new name in `main`.

## Reference answer after your attempt

<!-- task-answer -->
```rust
fn decorate(mut title: String) -> String {
    title.push_str("!");
    title
}

fn main() {
    let original = String::from("Plan");
    let decorated = decorate(original);
    println!("{decorated}");
}
```

```text
Plan!
```

Ownership goes from `original` to the parameter and then to `decorated`. [Rust's ownership chapter](https://doc.rust-lang.org/book/ch04-01-what-is-ownership.html).

[Previous lesson](/read/rust-18-memory?lang=en) · [Contents](/course/rust?lang=en)
