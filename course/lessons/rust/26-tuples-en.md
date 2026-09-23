# Tuples: two results from one function

_Lead (summary):_ **Lesson 26. Return a task title and its duration together without mixing up the values.**

Before this lesson: lessons 9, 14, 19 and 22. If returning a value from a function feels unfamiliar, revisit [lesson 14](/read/rust-14-functions?lang=en).

## Familiar image and recall map

A task card has a place for a title and another for a number of minutes. It is one card with two different entries. A **tuple** is one value made of a fixed number of ordered elements, which may have different types. Rust knows how many places it has, and each place's type, when it compiles the program. The card image has a limit: tuple places have positions rather than names. The next lesson introduces a struct when named fields would help.

**Two values → parentheses and comma → one tuple → return → unpack into names.** Say the chain without looking at code. `(String, u32)` is a tuple type: the first place holds an owned `String`, the second a `u32` whole number. `(title, minutes)` creates a value in that order. The comma separates places, and the parentheses group them into one value. `(u32, String)` is a different type because the order differs.

## Return one pair

Save your previous `src/main.rs` in `organizer`, replace the whole file with this example, and run `cargo run` beside `Cargo.toml`. No other files or dependencies change. Predict the output without Cargo's own messages.

The type after `->` in `first_task` is the single return type `(String, u32)`. The final line has no `;`, so it returns the tuple. In `let (title, minutes) = first_task();`, the pattern on the left **unpacks** it: the first element gets the name `title`, the second `minutes`. The function still returns one compound value.

```rust
fn first_task() -> (String, u32) {
    let title = String::from("Reading");
    let minutes = 20;
    (title, minutes)
}

fn main() {
    let (title, minutes) = first_task();
    println!("{title}: {minutes} min");
}
```
```text
Reading: 20 min
```

The `String` moves from the function into the tuple, then into `title`; this follows lesson 19's ownership rule. A `u32` can be copied. A tuple works well for a short intermediate result when the meaning of each place is clear near its use.

## Positions and the image's limit

To read one place, put a dot and its **position number** after the tuple name: `.0` is first and `.1` is second. Counting begins at zero, as with an array index. In the separate example below, `_u32` on `25` states its number type explicitly; lesson 10 covered number types. `true` is a Boolean value from lesson 11. The tuple holds minutes and a Boolean answer to whether there was enough time.

```rust
fn main() {
    let report = (25_u32, true);
    println!("Minutes: {}", report.0);
    println!("Finished: {}", report.1);
}
```

```text
Minutes: 25
Finished: true
```

Both elements here copy when read. Taking a `String` through `.0` can instead move ownership to a new variable; it does not copy text for free. A two-element tuple has no `.2`. This **separate** example intentionally fails with E0609:

<!-- error-code: E0609 -->
```rust,compile_fail
fn main() {
    let report = (20_u32, true);
    println!("{}", report.2);
}
```

Do not confuse `(20, true)` with an array `[20, true]`: an array requires all elements to have one type. The empty parentheses `()` from lesson 14 are a tuple with no elements and the unit value. A one-element tuple needs a comma, `(20,)`; `(20)` is just a number in parentheses.

## Recall without the page

1. Why does a function returning `(String, u32)` still return one value?
2. Which place does `.1` select?
3. How does `(20,)` differ from `(20)`?

Check: a tuple is one compound value; `.1` is the second place; the comma makes a one-element tuple. Close the page and say “two values → one tuple → return → unpack” from memory.

## Exercise

**Required.** Write `plan() -> (String, u32)` so it returns `Walk` and `30`. In `main`, unpack those values into `title` and `minutes`, then print `Walk — 30 min`. Both displayed values must come from the function, not from a finished message in `println!`.

**Boundary check.** Swap the order of the returned elements without changing the declared type. Read the compiler error, then restore the order.

## Hint

Make `(String::from("Walk"), 30)` the final expression of the function. Put `(title, minutes)` to the left of `=` in `main`.

## Reference answer after trying

<!-- task-answer -->
```rust
fn plan() -> (String, u32) {
    (String::from("Walk"), 30)
}

fn main() {
    let (title, minutes) = plan();
    println!("{title} — {minutes} min");
}
```
```text
Walk — 30 min
```

The tuple held values of two types together; unpacking gave each result a useful name. [Official tuple explanation](https://doc.rust-lang.org/book/ch03-02-data-types.html#the-tuple-type).

[Previous lesson](/read/rust-25-unicode?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-27-structs?lang=en)
