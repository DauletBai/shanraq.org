# Your first automated test

_Lead (summary):_ **Lesson 16. Make a calculation check itself against an expected result.**

Prerequisites: lessons 1–15, especially [functions](/read/rust-14-functions?lang=en) and [arrays](/read/rust-15-arrays?lang=en). We will test numbers before adding input or external libraries.

## A familiar picture and a recall map

Suppose you weigh a parcel after packing it. Writing down “this parcel should weigh 5 kg” lets you repeat the check tomorrow. An **automated test** records a similar rule for code: given an input, compare the actual result with the expected one. Its limit is important: checking one parcel does not prove every parcel has the right weight, and one passing test does not prove a function is correct for every input.

**Input → function → actual result ↔ expected result → test report.** Close the page and repeat this chain before reading the code.

## Run the program and its test

Two new notations come first. `#[test]` is an **attribute**, a label immediately before a function with no arguments. `#` begins the label, the square brackets hold its name, and `test` means “run the next function as a test.” `assert_eq!(actual, expected)` compares the first value, calculated by the code, with the second value, predicted by us. Its `!` marks a macro, as with `println!`. Equality passes; inequality fails and shows both values.

Use the `organizer` project from lesson 6. In the folder containing `Cargo.toml`, save your earlier `src/main.rs` elsewhere and replace the whole file with the following program. No other files or dependencies change. `cargo run` builds and runs `main`; `cargo test` builds and runs functions marked `#[test]`. Type both commands in that folder. The output block is the program's output from `cargo run`, without Cargo messages. Predict it first.

```rust
fn remaining(total: u32, done: u32) -> u32 {
    total - done
}

fn main() {
    println!("Remaining: {}", remaining(7, 2));
}

#[test]
fn subtracts_done_tasks() {
    assert_eq!(remaining(7, 2), 5);
}
```

```text
Remaining: 5
```

The first two functions use syntax from lesson 14. `main` does not call the marked test function. Because both are in the same file, the test can call `remaining` directly. We will group tests into modules after modules have been explained in lesson 44.

`assert_eq!` passes silently. `cargo test` should report `1 passed; 0 failed`; surrounding messages depend on your Cargo version. A test reports a problem but cannot repair the function.

## Break it on purpose

Change expected `5` to `6` and run `cargo test`: it must fail. Restore `5`. Change `total - done` to `total + done`, run the test again, then restore subtraction. If a test still passes, check that you saved the file and ran the command in the project folder. The function has a **precondition** from lesson 14: `done <= total`. With `u32`, a larger `done` is invalid input for this particular function. We will return explicit errors after learning `Result`.

## Recall before the exercise

Which command calls `main` and which finds `#[test]` functions? Does one successful example cover all `u32` values? Answer without looking: `cargo run`, then `cargo test`; no.

## Exercise

**Required.** Write `is_short(minutes: u32) -> bool` using the rule from lesson 14: 1–15 minutes inclusive. Add three tests: 0 → `false`, 15 → `true`, 16 → `false`. Write the predictions first, run `cargo test`, then run `cargo run`. For 15, `main` must print `Short: true`. Deliberately make one expected answer wrong, observe the failed test, and restore it. Revisit [lesson 11](/read/rust-11-boolean?lang=en) if the condition is unclear.

## Hint

Put three separate zero-argument functions after `main`, each preceded by `#[test]`. Compare `is_short(...)` with `true` or `false`.

## Reference answer after your attempt

<!-- task-answer -->
```rust
fn is_short(minutes: u32) -> bool {
    minutes > 0 && minutes <= 15
}

fn main() {
    println!("Short: {}", is_short(15));
}

#[test]
fn zero_is_not_a_task() {
    assert_eq!(is_short(0), false);
}

#[test]
fn fifteen_is_short() {
    assert_eq!(is_short(15), true);
}

#[test]
fn sixteen_is_long() {
    assert_eq!(is_short(16), false);
}
```

```text
Short: true
```

The three boundaries matter: a test of 15 alone would miss accepting 0 or 16. [Rust's testing chapter](https://doc.rust-lang.org/book/ch11-01-writing-tests.html).

[Previous lesson](/read/rust-15-arrays?lang=en) · [Contents](/course/rust?lang=en)
