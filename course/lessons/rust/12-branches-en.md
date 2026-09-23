# Branches and expressions

_Lead (summary):_ **Lesson 12. Choose a safe action and obtain a value from a block.**

Prerequisites: lessons 1–11. If variables are unclear, revisit [lesson 9](/read/rust-09-variables?lang=en).

## Why this matters

Lesson 11 produced `true` or `false`. **Branching** uses that answer to select which part of a program to run. Picture a fork in a path: take another route if one is closed. Unlike a person, the computer does not choose a route for itself; it checks the conditions you wrote in order.

Return to lesson 10. Subtracting completed tasks from a `u32` total is invalid when completion exceeds the total. We will check the data before performing that subtraction.

## Run the example

Use the `organizer` project from [lesson 6](/read/rust-06-cargo?lang=en), in the folder containing `Cargo.toml`. Save your previous work separately. Replace **all** of `src/main.rs` with the first example below, save it, and run `cargo run` in that folder’s terminal. As explained in lesson 6, this command builds and runs the program. No other files or dependencies change. Each subsequent complete example also replaces the entire file. Output blocks show only program output, without Cargo messages. Predict the output before running.

Before the example, read the new form. `if` checks a `bool` condition and runs the following `{ ... }` block when it is true. `else if` checks another possibility; `else`, without a condition, handles what remains. At most one branch in this chain runs. In `let urgent: bool = if ... { true } else { false };`, the chain itself chooses a value, and both blocks must yield `bool`. We will examine that rule after running the example.

```rust
fn main() {
    let total: u32 = 7;
    let done: u32 = 2;
    if done > total {
        println!("Error: completed exceeds total");
    } else if done == total {
        println!("All tasks completed");
    } else {
        let remaining = total - done;
        println!("Remaining: {remaining}");
    }
    let priority: u8 = 3;
    let urgent: bool = if priority >= 3 { true } else { false };
    println!("Urgent: {urgent}");
}
```

```text
Remaining: 5
Urgent: true
```

## Walkthrough

`if` introduces a `bool` condition followed by a block in `{ ... }`. If the condition is false, `else if` checks the next condition. `else` provides the fallback without a condition of its own. Only the first matching branch in this chain runs; later conditions are not checked.

With `done = 2`, both conditions are false and subtraction runs. With `done = 7`, the middle branch runs. With `done = 8`, the first branch prevents the unsafe subtraction. Two separate `if` expressions are different: both might run. Parentheses around a condition are optional; braces around the block are required. The number `1` cannot stand in for `true`.

## Expressions produce values

An **expression** computes a value: subtraction, comparison, and even `if`. A **statement** performs an action; a `let` declaration is a statement. In `let urgent: bool = if ...;`, the selected branch supplies either `true` or `false`.

A final expression **without `;`** supplies the block's value. Branches that supply a value need compatible types: a Boolean in one and a number in the other do not work here. Both our branches supply `bool`. The semicolon after the complete `let` declaration is still necessary. This simple example could be shortened to `let urgent = priority >= 3;`; the longer form demonstrates how `if` works.

A block without a final value expression produces **`()`**, called the **unit value**. Its type is also written `()`. It carries no useful numeric or Boolean answer and is neither zero nor `false`. `println!` displays text but does not return that text as its result. We will revisit unit with functions in lesson 14.

An ordinary block can produce a value too:

```rust
fn main() {
    let minutes = {
        let preparation = 5;
        preparation + 10
    };
    println!("{minutes}");
}
```

```text
15
```

`preparation` exists only inside its braces; the computed number receives the outer name `minutes`. This uses scope from lesson 9 to separate intermediate work from its result.

## A common mistake

An extra semicolon removes the first branch's Boolean result. This separate, deliberately invalid example produces E0308:

<!-- error-code: E0308 -->
```rust,compile_fail
fn main() {
    let priority: u8 = 3;
    let urgent: bool = if priority >= 3 { true; } else { false };
    println!("{urgent}");
}
```

Remove the semicolon after `true` so both branches supply `bool`. Keep the semicolon after the complete declaration.

## Reference map

Validate inputs → first matching branch → action or value.

## Check your understanding

1. How many branches of one chain run?
2. What do `{ 5 + 2 }` and `{ 5 + 2; }` produce?
3. Where is `remaining` available?

Check: at most one; `7` and `()`; only within the `else` block. Using `>=` in the first condition would wrongly reject equal counts.

## Exercise

**Required.** Validate task counts with `total = 4` and `done = 6`. Print only “Error: completed exceeds total”. Then try `done = 4` → “All tasks completed” and `done = 1` → “Remaining: 3”. Never compute the subtraction before the check.

**Further thought.** If the first condition becomes `done >= total`, which case can no longer reach the second branch?

## Hint

Reject invalid data first, then check equality. Subtract only inside `else`. The answer demonstrates invalid input handled by a correctly working program.

## Reference answer after your own attempt

<!-- task-answer -->
```rust
fn main() {
    let total: u32 = 4;
    let done: u32 = 6;
    if done > total {
        println!("Error: completed exceeds total");
    } else if done == total {
        println!("All tasks completed");
    } else {
        println!("Remaining: {}", total - done);
    }
}
```

```text
Error: completed exceeds total
```

[Verification source](https://doc.rust-lang.org/book/ch03-05-control-flow.html)

[Previous lesson](/read/rust-11-boolean?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-13-loops?lang=en)
