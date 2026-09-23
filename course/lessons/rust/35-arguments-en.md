# Launch arguments: a command and a title with spaces

_Lead (summary):_ **Lesson 35. Give the organizer an action and title when starting it from the terminal.**

Before this lesson: lessons 4, 14, 22, 26, 30–31 and 34. If `Option` and `Some` are unclear, revisit [lesson 31](/read/rust-31-option?lang=en).

## Familiar image and recall map

At a cafe, you can walk up to the counter and then say your order, or hand over a written order as you enter. Lesson 34 read a line **after launch**; now the program receives words **at launch**. These words are called **command-line arguments**. The image has a limit: the terminal's shell first splits your written command according to its rules; Rust receives the resulting values.

**Command line → `cargo run --` → program arguments → `args()` → `next()` → `Option<String>` → action.** The two dashes `--` after `cargo run` separate Cargo's options from the arguments for our organizer. In `cargo run -- add "Read Rust"`, the shell passes `add` and `Read Rust` as two arguments: quotes keep words with a space together, but the quotes are not part of the title. Without quotes, the title would become several arguments. `std::env::args()` provides `String` values in order; the first is usually the program path. `.next()` takes the next value as `Some(String)`, or returns `None` when no values remain. Each call advances. `_program` names a value that is intentionally unused: the initial `_` silences an unused-name warning.

## One command at launch

Save the previous `src/main.rs` in `organizer`, replace the whole file, and first run `cargo run` beside `Cargo.toml`. No other files or dependencies change. Without arguments, the program shows the usage line below. Then run `cargo run -- add "Read Rust"`: it prints `Add: Read Rust`. This line is a planned action; storing the task comes later.

The tuple `(command, title, extra)` groups three `.next()` results. The pattern `(Some(command), Some(title), None)` matches only when a command and title exist and no extra argument remains. `_` from lesson 30 handles other shapes. `command == "add"` tests the familiar command word, not a translated task title.

```rust
fn main() {
    let mut args = std::env::args();
    let _program = args.next();
    let command = args.next();
    let title = args.next();
    let extra = args.next();
    match (command, title, extra) {
        (Some(command), Some(title), None) => {
            if command == "add" {
                println!("Add: {title}");
            } else {
                println!("Unknown command: {command}");
            }
        }
        _ => println!("Usage: organizer add \"Title\""),
    }
}
```
```text
Usage: organizer add "Title"
```

With `cargo run -- add "Read Rust"`, the program prints `Add: Read Rust`. `add` is the command from lesson 2. The name `organizer` in the usage line shows a later way to run the built program; for now use `cargo run --`. The program does not save data or implement other commands yet.

## Missing or extra words

Without a title after `add`, `title` is `None`, so the usage line appears. If you omit quotes around a spaced title, `extra` becomes `Some(...)`, and the program shows usage instead of silently losing half of the title. `""` on the command line can pass an existing but empty title; this code does not check it yet. We will add that rule while parsing commands. `std::env::args()` handles arguments containing valid Unicode; a different method exists for other bytes, which we do not need here.

## Check your understanding

1. Why use `--` in `cargo run -- add ...`?
2. Why does the first `.next()` go into `_program`?
3. What changes if you remove the quotes around a title with spaces?

Check: separate program arguments from Cargo options; the first argument is the program path; the title splits and creates an extra argument. Rebuild the recall map.

## Exercise

**Required.** Use the same pattern, but print `Planned: TITLE` for `add`, `Command not found` for another command, and `Need: add "Title"` for missing or extra arguments. Try `cargo run -- add "My plan"` and a launch without arguments. The title must come from `args`, not a finished string in the code.

**Boundary check.** Run `cargo run -- add My plan` without quotes and explain why the usage line appears.

## Hint

Keep four `.next()` calls: program path, command, title and extra word. Inside the pattern with two `Some` values and one `None`, compare `command` with `"add"`.

## Reference answer after trying

<!-- task-answer -->
```rust
fn main() {
    let mut args = std::env::args();
    let _program = args.next();
    let command = args.next();
    let title = args.next();
    let extra = args.next();
    match (command, title, extra) {
        (Some(command), Some(title), None) => {
            if command == "add" {
                println!("Planned: {title}");
            } else {
                println!("Command not found");
            }
        }
        _ => println!("Need: add \"Title\""),
    }
}
```
```text
Need: add "Title"
```

The output block shows a launch without arguments; `cargo run -- add "My plan"` prints `Planned: My plan`. [Official command-line argument chapter](https://doc.rust-lang.org/book/ch12-01-accepting-command-line-arguments.html) and [the `cargo run --` rule](https://doc.rust-lang.org/cargo/commands/cargo-run.html).

[Previous lesson](/read/rust-34-input?lang=en) · [Contents](/course/rust?lang=en)
