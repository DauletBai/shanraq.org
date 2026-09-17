# Using compiler diagnostics

_Summary:_ **Separate a language-rule error from a mistake in your intention.**

Prerequisite: lessons 1–7. Work in organizer unless stated otherwise.

## Why this matters

A **diagnostic** is a tool's message about a problem. It points towards a repair, not towards a judgment about your ability. An **error** prevents the current action from succeeding. A **warning** normally permits compilation but identifies something worth checking.

## The whole picture

Start with this correct src/main.rs. Save and run it from organizer using cargo run:

```rust
fn main() {
    println!("Task: learn Rust");
}
```

Output:

```text
Task: learn Rust
```

Now make a deliberate mistake in a temporary copy: change println! to printline! and run cargo check. We expect compilation to fail, not program output:

```rust,compile_fail
fn main() {
    printline!("Task: learn Rust");
}
```

Expected diagnostic fragment: cannot find macro. We have not defined a macro named printline. Paths, line numbers and formatting can differ between tool versions; they need not match someone else's screenshot.

## How to read the message

Find src/main.rs, the line number and the marked location. Read the first useful error. Its location may be the consequence of an earlier missing quote, so inspect the preceding line too. After one repair, check again.

An error such as E0384 has an identifier. rustc --explain E0384 asks for its detailed explanation: --explain is the option, and E0384 is its argument. Not every diagnostic has such a code.

Distinguish three situations. Source that does not compile breaks a language rule. A program that starts and then fails needs its input and runtime conditions investigated. A program that finishes but gives the wrong answer needs its logic checked against the task. Correcting spelling alone cannot resolve every situation.

## Recall map

Read → locate → understand the cause → change one thing → check → run → compare the result.

## Warm-up

1. Predict: can wrong output text compile successfully?
2. Complete: Compiling is printed by ___; our task line by ___.
3. Restore the known macro name println!.

## Exercise

**Required.** Fix the misspelled macro and change the output to “Done”. Run check, then run. Expect a successful check and that one output line.

**Your own data.** Misspell only the text inside quotes. Notice that syntax checking cannot know which word you intended.

## Hint and reference answer

The macro is println!, while you choose the message inside quotes. Warm-up answers: yes, wrong wording can compile; Cargo prints Compiling; our executable prints its message. When asking for help, provide the command, system, version, full diagnostic and a minimal example. Do not send secrets or replace the message with “it doesn't work”.

<!-- task-answer -->
```rust
fn main() {
    println!("Done");
}
```

Expected output:

```text
Done
```

[Previous lesson](/read/rust-07-first-program?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-09-variables?lang=en)
