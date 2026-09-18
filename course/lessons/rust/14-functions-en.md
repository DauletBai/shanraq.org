# Functions

_Lead (summary):_ **Lesson 14. Separate calculation from display and explain how values are passed.**

Prerequisites: lessons 1–13. If variables are unclear, revisit [lesson 9](/read/rust-09-variables?lang=en).

## Why this matters

We already write `fn main()`. Now we will name a calculation of our own. A **function** is a named piece of code that can be called with input values and can return a result. Think of a recipe written once and used for different quantities. The limit of this analogy: a function only works with data passed to it or available in its scope, not an understanding of the whole project.

Separate computing a remaining count from displaying it. That will make calculations easier to check without depending on messages. We reuse conditions and block values from lessons 11–12.

## Run the example

Use the `organizer` project from [lesson 6](/read/rust-06-cargo?lang=en), in the folder containing `Cargo.toml`. Save your previous work separately. Replace **all** of `src/main.rs` with the first example below, save it, and run `cargo run` in that folder’s terminal. As explained in lesson 6, this command builds and runs the program. No other files or dependencies change. Each subsequent complete example also replaces the entire file. Output blocks show only program output, without Cargo messages. Predict the output before running.

```rust
fn remaining(total: u32, done: u32) -> u32 {
    total - done
}

fn print_remaining(count: u32) {
    println!("Remaining: {count}");
}

fn main() {
    let total: u32 = 7;
    let done: u32 = 2;
    if done <= total {
        let count = remaining(total, done);
        print_remaining(count);
    } else {
        println!("Invalid data");
    }
}
```

```text
Remaining: 5
```

## Walkthrough

`fn` declares a function. `remaining` is our chosen name. Inside parentheses, **parameters** name input values and specify their required types, separated by commas. `-> u32` states the return type. Braces enclose the function body. Its last expression, `total - done` without `;`, supplies the result just like a block value in lesson 12.

**Arguments** are the actual values supplied by a call. `remaining(total, done)` passes the caller's variable values. The first argument goes to the first parameter, the second to the second: matching names do not create a name-based connection. `remaining(7, 2)` is a valid call too. `remaining` does not see the local variables in `main`; its parameters `total` and `done` are its own names. Each call supplies parameter values afresh.

Declaring a function does not execute its body. Execution starts in `main`, reaches a call, runs the called function, and resumes the caller's next action. In this file functions can be declared before or after `main`; declaration order does not determine call order. Repeated calls do not automatically retain a local counter.

`print_remaining(count)` only displays a message. With no `-> ...`, its return type is `()` from lesson 12. Printing is an action, not a returned number. Ordinary function calls have no `!`; `println!` remains a macro.

## A contract for inputs

Our `remaining` function requires `done <= total`. This is a **precondition**, a rule the calling code must satisfy. Here `main` checks it **before** the call. The `u32` parameter types alone do not prove this rule. Do not call `remaining(2, 7)` or conceal invalid counts by pretending the result is zero. Lessons on `Option` and `Result` will introduce explicit missing results and errors; for now this contract limits our example.

## Returning early

`return` immediately ends the current function call, not the whole process or merely the enclosing `if` block. It is followed by the value to return. This separate example asks whether a task fits the available time; our rule rejects zero-minute tasks.

```rust
fn fits(minutes: u32, available: u32) -> bool {
    if minutes == 0 {
        return false;
    }
    minutes <= available
}

fn main() {
    println!("{}", fits(0, 30));
    println!("{}", fits(15, 30));
    println!("{}", fits(40, 30));
}
```

```text
false
true
false
```

`available` is the number of available minutes. A zero duration reaches `return false;`; otherwise the final comparison without `;` supplies the result.

## A common mistake

This function promises `u32`, but a semicolon makes its block produce `()`. It intentionally fails with E0308:

<!-- error-code: E0308 -->
```rust,compile_fail
fn remaining(total: u32, done: u32) -> u32 {
    total - done;
}

fn main() {
    println!("{}", remaining(7, 2));
}
```

Remove the semicolon after the final expression, or write `return total - done;` explicitly. The first version is shorter for a final expression.

## Reference map

Arguments → parameters → body → returned value → caller.

## Check your understanding

1. Does declaring a function execute it?
2. What does a function without `-> ...` return?
3. What must `main` check before calling `remaining`?

Check: no, it needs a call; `()`; completion must not exceed the total. A parameter is a declared input name; an argument is a value supplied in a call.

## Exercise

**Required.** Write `is_short(minutes: u32) -> bool`: a task is short if it lasts 1 through 15 minutes inclusive. Compute and return the answer; do not print inside the function. In `main`, print results for arguments 0, 15, 16: `false`, `true`, `false`.

**Self-check.** Identify the parameter, arguments, and return type. Also call the function with 1: expect `true`. Do not write a separate function for each input.

## Hint

Combine `minutes > 0` and `minutes <= 15` with `&&`. Leave the function’s final expression without `;`. In `main`, use empty `{}` in `println!` to display a returned answer.

## Reference answer after your own attempt

<!-- task-answer -->
```rust
fn is_short(minutes: u32) -> bool {
    minutes > 0 && minutes <= 15
}

fn main() {
    println!("{}", is_short(0));
    println!("{}", is_short(15));
    println!("{}", is_short(16));
}
```

```text
false
true
false
```

[Verification source](https://doc.rust-lang.org/book/ch03-03-how-functions-work.html)

[Previous lesson](/read/rust-13-loops?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-15-arrays?lang=en)
