# The question mark: return an error to the caller

_Lead (summary):_ **Lesson 33. Parse minutes first, then check their meaning and return a useful error.**

Before this lesson: lessons 12, 14 and 30–32. If `Ok` and `Err` are unclear, revisit [lesson 32](/read/rust-32-result?lang=en).

## Familiar image and recall map

A clerk checks an application in two stages: first whether a number was written, then whether it is acceptable. If the first check fails, there is no reason to continue; the application goes back with an explanation. The **`?`** symbol after a `Result` expression acts similarly: with `Ok`, it extracts the value and continues the function; with `Err`, it immediately returns the error from that function. This is called **error propagation**. The image has a limit: `?` neither invents an explanation nor repairs a failure; it passes a value according to Rust's type rules.

**Parse → `Ok` continues / `Err` returns → check zero → `Ok` or a new `Err`.** The function using `?` here returns `Result<u32, String>`. The expression before `?` also gives `Result<u32, String>`, so the error types match. `return Err(...)` is the early return from lesson 14: the function ends at that line. Here it creates a **new** explanation for our zero-minute rule.

## Separate parsing from validation

Save the previous `src/main.rs` in `organizer`, replace the whole file, and run `cargo run` beside `Cargo.toml`. No other files or dependencies change. Predict the three output lines.

`parse_number` distinguishes number text from other text. `checked_minutes` receives the number through `parse_number(text)?`. On an error, the later lines in that function do not run. If the number is zero, the function creates a different error. `main` decides what to show the user; the checking functions do not print messages themselves.

```rust
fn parse_number(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(number) => Ok(number),
        Err(_) => Err(String::from("Enter a whole number")),
    }
}

fn checked_minutes(text: &str) -> Result<u32, String> {
    let minutes = parse_number(text)?;
    if minutes == 0 {
        return Err(String::from("Minutes must be greater than zero"));
    }
    Ok(minutes)
}

fn main() {
    for text in ["20", "0", "soon"] {
        match checked_minutes(text) {
            Ok(minutes) => println!("Accepted: {minutes} min"),
            Err(message) => println!("Error: {message}"),
        }
    }
}
```
```text
Accepted: 20 min
Error: Minutes must be greater than zero
Error: Enter a whole number
```

The spelling error and the meaning error stay separate. `?` does not print to the terminal automatically: the `Err` branch in `main` does that. Later, this helps keep organizer rules apart from its terminal interface.

## Where `?` is allowed

Use `?` where the current function can return a compatible error. A plain `fn main()` with no return type returns `()` from lesson 14, so it cannot return an `Err(String)` this way. This **separate** example intentionally fails with E0277:

<!-- error-code: E0277 -->
```rust,compile_fail
fn main() {
    let value: Result<u32, String> = Ok(20);
    let minutes = value?;
    println!("{minutes}");
}
```

There are other uses of `?`. For now, stay with functions that have matching `Result<u32, String>` error types. Different error types require a conversion we have explicitly learned; do not insert `.unwrap()` merely to silence a compiler error.

## Check your understanding

1. What happens to `Ok(20)` before `?`?
2. Does the zero check run after `parse_number` returns `Err`?
3. Which function prints the user-facing message?

Check: extract 20 and continue; no, the error returns early; `main`. Rebuild the recall map without the page.

## Exercise

**Required.** Write `parse_number` and `checked_minutes` with `?` as shown, but use `Enter a number` and `Minutes cannot be zero` for the errors. Check `"15"`, `"0"` and `"many"`. Expect `Accepted: 15 min` and two different error messages. Do not print inside the checking functions.

**Boundary check.** Explain why `"0"` reaches the zero check but `"many"` does not.

## Hint

Put `let minutes = parse_number(text)?;` inside a function returning `Result<u32, String>`. Then have `if minutes == 0` return a new `Err`.

## Reference answer after trying

<!-- task-answer -->
```rust
fn parse_number(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(number) => Ok(number),
        Err(_) => Err(String::from("Enter a number")),
    }
}

fn checked_minutes(text: &str) -> Result<u32, String> {
    let minutes = parse_number(text)?;
    if minutes == 0 {
        return Err(String::from("Minutes cannot be zero"));
    }
    Ok(minutes)
}

fn main() {
    for text in ["15", "0", "many"] {
        match checked_minutes(text) {
            Ok(minutes) => println!("Accepted: {minutes} min"),
            Err(message) => println!("Error: {message}"),
        }
    }
}
```
```text
Accepted: 15 min
Error: Minutes cannot be zero
Error: Enter a number
```

Two functions check the values; `main` presents the message. [Official explanation of `?`](https://doc.rust-lang.org/book/ch09-02-recoverable-errors-with-result.html#propagating-errors).

[Previous lesson](/read/rust-32-result?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-34-input?lang=en)
