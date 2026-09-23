# A command as a checked intention

_Lead (summary):_ **Lesson 36. Separate parsing launch words from acting on the organizer and explain each error.**

Prerequisites: lessons 29–35. If `Option` and `Result` are unclear, revisit [lesson 31](/read/rust-31-option?lang=en).

## Familiar image and recall map

At a service desk, someone checks the request form before carrying it out. The organizer can do the same: turn words into a checked command first. **Command parsing** checks the action name and the data it requires. It does not change the task list yet. The analogy has a limit: this program does not understand arbitrary sentences; it accepts only the forms we define.

**Words → shape check → `Result<Command, String>` → `Ok` with a command or `Err` with a reason → action later.** `Command` is our own enum from lesson 29. `Add(String)` carries a title; `List` carries no data. `Option<&str>` means a borrowed piece of text may be present (`Some`) or absent (`None`), as in lesson 31. A tuple of three values lets `match` check the action, title and extra argument together. `title.trim().is_empty()` tests whether anything remains after removing edge spaces. `String::from(title)` makes an owned string: the original arguments may disappear while the command still needs its title. `_` covers the remaining shapes and returns an explanation. Here `Ok(...)` and `Err(...)` are result values; neither prints.

## Parse, then show

Save the previous `src/main.rs` in your `organizer` project, replace it completely, and run `cargo run` beside `Cargo.toml`. This teaching example parses fixed words; we will connect actual arguments from lesson 35 when assembling the full organizer. No other files or dependencies change. Predict the output before running.

```rust
enum Command {
    Add(String),
    List,
}

fn parse(action: &str, title: Option<&str>, extra: Option<&str>) -> Result<Command, String> {
    match (action, title, extra) {
        ("list", None, None) => Ok(Command::List),
        ("add", Some(title), None) => {
            if title.trim().is_empty() {
                Err(String::from("Empty title"))
            } else {
                Ok(Command::Add(String::from(title)))
            }
        }
        _ => Err(String::from("Unknown command or argument count")),
    }
}

fn show(result: Result<Command, String>) {
    match result {
        Ok(Command::Add(title)) => println!("Add: {title}"),
        Ok(Command::List) => println!("Show list"),
        Err(message) => println!("Error: {message}"),
    }
}

fn main() {
    show(parse("add", Some("Read Rust"), None));
    show(parse("list", None, None));
    show(parse("add", Some("   "), None));
}
```
```text
Add: Read Rust
Show list
Error: Empty title
```

`parse` only creates a checked command. `show` displays it; the task list comes in the next lesson. If you call `parse("list", Some("extra"), None)`, the shape check returns an error instead of ignoring the word.

## Check your understanding

1. Why does `parse` return a command instead of adding a task immediately?
2. How does `None` differ from `Some("")`?
3. Which function prints an error?

Parsing leaves the list alone; one title is absent and the other was supplied empty; `main` prints. Rebuild the recall map without looking.

## Exercise

**Required.** Add `Done(u32)` to the enum for a `done` action with one number. Parse the number with `parse::<u32>()`, reject zero and invalid text, and keep the extra-argument check. Try `done 2`, `done 0`, and `done no`. The parser must still not print.

## Hint

A new `("done", Some(text), None)` arm can call `text.parse::<u32>()`. Handle `Ok(number)` and `Err(_)` with `match`; return `Err` for zero.

## Reference after attempting

<!-- task-answer -->
```rust
enum Command {
    Add(String),
    List,
    Done(u32),
}

fn parse(action: &str, title: Option<&str>, extra: Option<&str>) -> Result<Command, String> {
    match (action, title, extra) {
        ("list", None, None) => Ok(Command::List),
        ("add", Some(title), None) => {
            if title.trim().is_empty() {
                Err(String::from("Empty title"))
            } else {
                Ok(Command::Add(String::from(title)))
            }
        }
        ("done", Some(text), None) => match text.parse::<u32>() {
            Ok(0) => Err(String::from("Number must be positive")),
            Ok(number) => Ok(Command::Done(number)),
            Err(_) => Err(String::from("A number is required")),
        },
        _ => Err(String::from("Unknown command or argument count")),
    }
}

fn show(result: Result<Command, String>) {
    match result {
        Ok(Command::Add(title)) => println!("Add: {title}"),
        Ok(Command::List) => println!("Show list"),
        Ok(Command::Done(number)) => println!("Complete: {number}"),
        Err(message) => println!("Error: {message}"),
    }
}

fn main() {
    show(parse("add", Some("Read Rust"), None));
    show(parse("list", None, None));
    show(parse("add", Some("   "), None));
    show(parse("done", Some("2"), None));
    show(parse("done", Some("0"), None));
    show(parse("done", Some("no"), None));
}
```
```text
Add: Read Rust
Show list
Error: Empty title
Complete: 2
Error: Number must be positive
Error: A number is required
```

[Official explanations of enums and `match`](https://doc.rust-lang.org/book/ch06-02-match.html).

[Previous lesson](/read/rust-35-arguments?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-37-operations?lang=en)
