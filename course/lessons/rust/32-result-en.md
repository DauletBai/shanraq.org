# Result: explain a failed operation

_Lead (summary):_ **Lesson 32. Turn text into a number of minutes and explain invalid input without stopping the program.**

Before this lesson: lessons 10, 22 and 29–31. If `Some` and `None` are unclear, revisit [lesson 31](/read/rust-31-option?lang=en).

## Familiar image and recall map

If a card is missing from a drawer, we have an absence. If the card exists but says “soon” where minutes should be, reading a number failed. **`Result<T, E>`** is Rust's standard enum for an operation that may succeed or fail. `Ok(T)` holds a successful value of type `T`; `Err(E)` holds error information of type `E`. The image has a limit: the compiler makes us handle both variants, but cannot write a helpful human message for us.

**Text → attempt to parse → `Ok(number)` or `Err(reason)` → `match` → response.** In `Result<u32, String>`, success holds a `u32` and failure holds an owned `String` message. The comma inside `<u32, String>` separates the two types. `.parse::<u32>()` tries to read a `u32` from text; `::<u32>` tells the method which result type to use. It is a type indication after the method name, not a comparison. The method returns a `Result`, so invalid text alone does not crash the program.

## Parse two inputs

Save the previous `src/main.rs` in `organizer`, replace the whole file, and run `cargo run` beside `Cargo.toml`. No other file or dependency changes. Predict both messages first.

`read_minutes` changes a parsing failure into our own explanation. `Err(_)` recognizes failure; the `_` from lesson 30 discards technical detail we do not show yet. Both `match` branches return `Result<u32, String>`. `main` then handles success and error explicitly.

```rust
fn read_minutes(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(minutes) => Ok(minutes),
        Err(_) => Err(String::from("Enter whole minutes")),
    }
}

fn main() {
    for text in ["25", "soon"] {
        match read_minutes(text) {
            Ok(minutes) => println!("Minutes: {minutes}"),
            Err(message) => println!("Error: {message}"),
        }
    }
}
```
```text
Minutes: 25
Error: Enter whole minutes
```

`0` also parses successfully as a `u32`: it is a valid number spelling. The rule “minutes must be greater than zero” checks meaning, not syntax; we add it in the next lesson. A negative number does not fit `u32`, so the function returns `Err`.

## An error is a value, not a crash

A `Result<u32, String>` cannot be used as an already available `u32`. This **separate** example fails with E0308:

<!-- error-code: E0308 -->
```rust,compile_fail
fn main() {
    let minutes: u32 = "25".parse::<u32>();
    println!("{minutes}");
}
```

`Result` has an `.unwrap()` method that extracts the value from `Ok`, but panics on `Err`. It is unsuitable for user input: an invalid number is an expected situation to explain. Handle both variants with `match` here. `Option` from lesson 31 answers “is there a value?”, while `Result` can also carry a failure reason.

## Check your understanding

1. What does `String` in the second place of `Result<u32, String>` mean?
2. Why is `"0".parse::<u32>()` not an error by itself?
3. Why is `.unwrap()` risky for human input?

Check: the error message type; zero is a valid number spelling; on `Err`, the program panics instead of explaining. Rebuild the recall map from memory.

## Exercise

**Required.** Write `read_minutes(text: &str) -> Result<u32, String>` as in the example, but use `Write minutes as a number` for the error message. Check `"30"` and `"many"`, printing `Minutes: 30` and `Error: Write minutes as a number`. Take the number from `Ok` and the message from `Err`.

**Boundary check.** Check `"0"`. Explain why parsing succeeds even though a task rule may later reject zero.

## Hint

First use `match text.parse::<u32>()`; then use a second `match` in `main` on the function result. Do not call `.unwrap()`.

## Reference answer after trying

<!-- task-answer -->
```rust
fn read_minutes(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(minutes) => Ok(minutes),
        Err(_) => Err(String::from("Write minutes as a number")),
    }
}

fn main() {
    for text in ["30", "many"] {
        match read_minutes(text) {
            Ok(minutes) => println!("Minutes: {minutes}"),
            Err(message) => println!("Error: {message}"),
        }
    }
}
```
```text
Minutes: 30
Error: Write minutes as a number
```

The program reports invalid input and continues running. [Official `Result` chapter](https://doc.rust-lang.org/book/ch09-02-recoverable-errors-with-result.html).

[Previous lesson](/read/rust-31-option?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-33-propagation?lang=en)
