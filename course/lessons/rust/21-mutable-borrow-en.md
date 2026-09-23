# Mutable borrowing

_Lead (summary):_ **Lesson 21. Change a string through a function and explain when other access becomes available again.**

Prerequisites: lessons 19–20. If it is unclear why `&title` keeps the owner, return to [lesson 20](/read/rust-20-borrowing?lang=en).

## A familiar image and recall map

Several people can read the same reference book. Someone correcting a page needs temporary exclusive access so others do not read an unfinished change. **Mutable borrowing** grants temporary permission to change a value without transferring ownership. The image has a limit: Rust checks uses of references in code; it does not issue physical library passes.

**Owner `mut` → `&mut` → change in a function → last use of the reference → owner can act again.** Close this page and reconstruct the chain. Familiar `let mut` permits change through the owner's name. New `&mut title` creates a mutable reference; `&mut String` in a parameter describes its type. Write `&mut` at the call and in the parameter type; these do not make two changes.

## First example

In the `organizer` project, save your previous `src/main.rs` separately and replace the whole file. Run `cargo run` in the folder containing `Cargo.toml`. No other files or dependencies change. Predict the line first; the output block omits Cargo messages.

The `push_str` method from lesson 18 appends text to a `String`. The function below gets temporary permission to make that change; afterward the owner remains in `main`.

```rust
fn mark_done(title: &mut String) {
    title.push_str(" — done");
}

fn main() {
    let mut title = String::from("Reading");
    mark_done(&mut title);
    println!("Task: {title}");
}
```

```text
Task: Reading — done
```

`mark_done` does not own or return the string: its reference permits changing the original value. `mut` on the owner and `&mut` on the reference have different jobs. Having `mut` does not make every reference to the string mutable.

## One changing access at a time

While a mutable reference will still be used, a reading reference to the same value cannot be used at the same time. Otherwise a reader could see a changing value. This **separate** example intentionally fails to compile. E0502 reports conflicting borrows:

<!-- error-code: E0502 -->
```rust,compile_fail
fn main() {
    let mut title = String::from("Plan");
    let reader = &title;
    let writer = &mut title;
    writer.push_str("!");
    println!("{reader}");
}
```

Here `reader` is needed after `writer` is created, so their periods of use overlap. Finish reading **before** mutable borrowing to fix it. Adding `clone()` without thought would create another string instead of changing the original.

```rust
fn main() {
    let mut title = String::from("Plan");
    let reader = &title;
    println!("Before: {reader}");
    let writer = &mut title;
    writer.push_str("!");
    println!("After: {title}");
}
```

```text
Before: Plan
After: Plan!
```

The compiler sees the last use of `reader` before `writer` is created, even though both names appear in one block. Likewise the owner becomes available after the last use of `writer`. Two mutable references used at the same time for one value are also forbidden. The rule concerns **overlapping uses**; it does not forbid changing a string twice during the whole program.

## Check your understanding

1. Why do we need both `let mut title` and `&mut title`?
2. Why does the example with `reader` fail to compile?
3. After what action can the owner read `title` again?

Check: the owner permits change and the reference conveys that permission temporarily; the read overlaps the write; after the mutable reference's last use. Rebuild the recall map from memory.

## Exercise

**Required.** Write `add_label(title: &mut String, label: &str)` to append `label` to the string. You met `&str` in lesson 9: access to text without ownership; the literal `"!"` can be passed as `&str`. Create a mutable `String` containing `Plan`, call the function twice with `"!"`, and print `Plan!!`. Predict, then run. Do not use `clone()` or return a `String` from the function.

**Boundary.** Replace `&mut String` with `&String` in the parameter and explain the compiler's refusal. Restore the working form.

## Hint

Use `let mut` for the owner and `add_label(&mut title, "!")` for each call. Inside the function use `title.push_str(label)`. The first call's borrow ends before the second begins.

## Reference answer after your attempt

<!-- task-answer -->
```rust
fn add_label(title: &mut String, label: &str) {
    title.push_str(label);
}

fn main() {
    let mut title = String::from("Plan");
    add_label(&mut title, "!");
    add_label(&mut title, "!");
    println!("{title}");
}
```

```text
Plan!!
```

Both calls change one string; its owner stays in `main`. [Official explanation of mutable references](https://doc.rust-lang.org/book/ch04-02-references-and-borrowing.html).

[Previous lesson](/read/rust-20-borrowing?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-22-strings?lang=en)
