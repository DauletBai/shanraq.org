# Option: when a record might be absent

_Lead (summary):_ **Lesson 31. Look up a title by position without crashing when that position does not exist.**

Before this lesson: lessons 15, 23–24 and 29–30. If `match` branches are unclear, revisit [lesson 30](/read/rust-30-patterns?lang=en).

## Familiar image and recall map

Ask an archivist for a numbered card that is not in the drawer. They either hand you the card or say that it is absent. **`Option<T>`** is Rust's standard enum for a value that **might be absent**. `Some(T)` holds a value of type `T`; `None` means no value is present. The image has a limit: `None` does not explain why the card is missing. The next lesson introduces a different type when a failure needs a reason.

**Request → `Option<T>` → `Some(value)` or `None` → handle both with `match`.** As with `Vec<String>`, the angle brackets name the type of possible contents. In `Option<&String>`, the contents would be a reference to a `String`, not a copied string. `Some` and `None` need no separate import. An `Option<T>` is not already a `T`; inspect its variant before using the inner value.

## Read a list safely

Save the previous `src/main.rs` in `organizer`, replace the whole file, and run `cargo run` beside `Cargo.toml`. No other files or dependencies change. Predict both output lines.

The `.get(index)` method on `Vec` asks for an element at a position. Unlike `titles[index]`, a missing position does not stop the program: `.get()` returns `Option<&String>`. `Some(title)` in `match` names a reference to the found text; `None` handles absence. We build the list with `Vec::new()` and `push` from lesson 24. Position 0 exists, but position 2 does not exist in a one-element list.

```rust
fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Reading"));
    for index in [0, 2] {
        match titles.get(index) {
            Some(title) => println!("Found: {title}"),
            None => println!("No record at position {index}"),
        }
    }
}
```
```text
Found: Reading
No record at position 2
```

`get` neither copies the title nor transfers ownership. The found value is borrowed from `titles`, so the list remains available.

## Absent is different from empty text

`Some(String::from(""))` means a value exists but its text is empty. `None` means no value exists at all. These are different states. With `None` alone, Rust sometimes cannot infer what could have been inside; `let missing: Option<String> = None;` states the type. For `titles.get(index)`, the list determines it.

Assigning the result of `.get()` to a plain reference without inspecting the variant makes this **separate** example fail with E0308:

<!-- error-code: E0308 -->
```rust,compile_fail
fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Reading"));
    let title: &String = titles.get(0);
    println!("{title}");
}
```

Use `match` with `Some` and `None`, or another absence-handling method after it has been explained. A list position is still not a task's permanent ID: `get` avoids an out-of-bounds access but does not solve the stable identifier problem in lesson 38.

## Check your understanding

1. How does `Option<&String>` differ from `&String`?
2. Does `None` mean that an empty string was found?
3. Why does `.get(2)` avoid a panic for a one-element list?

Check: the reference might be absent; no, empty text is still a value; `get` returns `None`. Close the page and rebuild the recall map.

## Exercise

**Required.** Build a list containing `Reading` and `Walk`. Request positions 1 and 2 through `.get()` and handle both with `match`. Print `Found: Walk` and `Position 2 is absent`. Do not check the length yourself or put a finished title in the message.

**Boundary check.** Create an empty list and explain what `.get(0)` returns.

## Hint

For each `index` in `[1, 2]`, use `match titles.get(index)`. Print the found `title` in `Some(title)` and the `index` in `None`.

## Reference answer after trying

<!-- task-answer -->
```rust
fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Reading"));
    titles.push(String::from("Walk"));
    for index in [1, 2] {
        match titles.get(index) {
            Some(title) => println!("Found: {title}"),
            None => println!("Position {index} is absent"),
        }
    }
}
```
```text
Found: Walk
Position 2 is absent
```

The list remains the owner of its titles, and absence is handled explicitly. [Official `Option` chapter](https://doc.rust-lang.org/book/ch06-01-defining-an-enum.html#the-option-enum).

[Previous lesson](/read/rust-30-patterns?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-32-result?lang=en)
