# Useful commands

_Summary:_ **Lesson 54. Provide clear help and distinguish normal output from errors.**

Revisit [launch arguments](/read/rust-35-arguments?lang=en) and [command parsing](/read/rust-36-commands?lang=en). This is a small organizer shell without persistent data.

## Familiar image and recall map

A service desk has one counter for routine requests and another for problems. In a terminal, **stdout** carries ordinary output and **stderr** carries error messages. An **exit status** tells another program whether the action succeeded: 0 means success; here 2 means a bad command. **Recall map:** arguments → select command → stdout or stderr → exit status. The analogy does not guarantee a fixed display order between the two streams.

`std::env::args()` yields launch words; its first word is the program path. `skip(1)` omits it and `collect()` gathers the rest into `Vec<String>`. `args.first()` can find no word, so `map` transforms an optional value. `as_str()` supplies a string view for matching. In `match`, `|` means any listed alternative and `if args.len() <= 1` adds a condition. `println!` writes to stdout, `eprintln!` to stderr. `std::process::exit(code)` reports a nonzero status to the system.

## Run and inspect

Create an empty Cargo project and paste the first Rust block into `src/main.rs`. `cargo run -- help` passes `help` to the program: the two dashes end Cargo arguments. `cargo run -- list` prints an empty list. `cargo run -- mystery` reports an error; status 2 belongs to the program. The first output below comes from running without arguments. To inspect the program status without Cargo in the middle, run the built file at `target/debug/organizer` with an unknown word and inspect your shell’s status.

```rust
fn run(args: &[String]) -> i32 {
    match args.first().map(|word| word.as_str()) {
        None | Some("help") | Some("--help") if args.len() <= 1 => {
            println!("Commands: help, list");
            0
        }
        Some("list") if args.len() == 1 => {
            println!("No tasks yet");
            0
        }
        _ => {
            eprintln!("Unknown command");
            2
        }
    }
}

fn main() {
    let args: Vec<String> = std::env::args().skip(1).collect();
    let code = run(&args);
    if code != 0 {
        std::process::exit(code);
    }
}
```
```text
Commands: help, list
```

## Recall without looking

1. Why skip the first `args()` word?
2. Which stream gets an unknown-command message?
3. Why does a script need an exit status?

## Exercise

**Required.** Add a `status` command with no extra words. It should print task count 0 and return status 0. Add it to help. Ensure `status extra` still fails with status 2.

## Answers

Add a `Some("status") if args.len() == 1` arm before `_`. Update the help text and leave the error arm intact.

<!-- task-answer -->
```rust
fn run(args: &[String]) -> i32 {
    match args.first().map(|word| word.as_str()) {
        None | Some("help") | Some("--help") if args.len() <= 1 => {
            println!("Commands: help, list, status");
            0
        }
        Some("list") if args.len() == 1 => {
            println!("No tasks yet");
            0
        }
        Some("status") if args.len() == 1 => {
            println!("Tasks: 0");
            0
        }
        _ => {
            eprintln!("Unknown command");
            2
        }
    }
}

fn main() {
    let args: Vec<String> = std::env::args().skip(1).collect();
    let code = run(&args);
    if code != 0 {
        std::process::exit(code);
    }
}
```
```text
Commands: help, list, status
```

## After checking

Streams and status matter when scripts invoke the organizer. A real file-save failure should go to stderr with a nonzero status. Return to [lesson 36](/read/rust-36-commands?lang=en). [args documentation](https://doc.rust-lang.org/std/env/fn.args.html).

[Previous lesson](/read/rust-53-backup?lang=en) · [Contents](/course/rust?lang=en)
