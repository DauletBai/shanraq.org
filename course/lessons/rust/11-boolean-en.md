# Boolean values

_Lead (summary):_ **Lesson 11. Write and verify a rule for selecting tasks.**

Prerequisites: lessons 1–10. If variables are unclear, revisit [lesson 9](/read/rust-09-variables?lang=en).

## Why this matters

Lesson 10 counted tasks. Now we need a decision: show an unfinished task if it is either important or short. A **Boolean value** has exactly two possibilities: `true` and `false`. Their type is `bool`.

Think of a switch with two positions. The analogy stops there: `false` does not mean “bad task” or “missing information”; it answers a particular question. We still set the data in the source code rather than reading terminal input.

## Run the example

Use the `organizer` project from [lesson 6](/read/rust-06-cargo?lang=en), in the folder containing `Cargo.toml`. Save your previous work separately. Replace **all** of `src/main.rs` with the first example below, save it, and run `cargo run` in that folder’s terminal. As explained in lesson 6, this command builds and runs the program. No other files or dependencies change. Each subsequent complete example also replaces the entire file. Output blocks show only program output, without Cargo messages. Predict the output before running.

```rust
fn main() {
    let done: bool = false;
    let priority: u8 = 3;
    let minutes: u32 = 20;
    let show = !done && (priority >= 3 || minutes <= 15);
    println!("Show: {show}");
    let slots: u32 = 0;
    let tasks: u32 = 6;
    let fits = slots != 0 && tasks / slots <= 2;
    println!("Fits: {fits}");
}
```

```text
Show: true
Fits: false
```

## Walkthrough

`done` records completion, `priority` importance, `minutes` duration, and `show` the selection decision. In this example priority 3 or higher counts as important. `let done: bool = false;` states a type explicitly; Rust infers the type of `show` from the expression.

| Notation | Question |
|---|---|
| `a == b` | Are they equal? |
| `a != b` | Are they different? |
| `a < b`, `a <= b` | Less? Less than or equal? |
| `a > b`, `a >= b` | Greater? Greater than or equal? |

`=` assigns a value; `==` compares values. A comparison produces a `bool`. Before a Boolean, `!` reverses the answer: `!false` is `true`. This use differs from the exclamation mark in the macro name `println!`.

`&&` means “and”: both conditions must hold. `||` is inclusive “or”: at least one must hold, including the case where both hold. Type two `|` characters without a space between them. Our rule checks `!done` and then the parenthesized group. Precedence is `!`, comparisons, `&&`, then `||`; parentheses make the intended grouping clear. A single `&` has a different meaning and is not a substitute for `&&`.

| a | b | a && b | a \|\| b |
|---|---|---|---|
| false | false | false | false |
| false | true | false | true |
| true | false | false | true |
| true | true | true | true |

## Short-circuit evaluation

**Short-circuit evaluation** means evaluating the right operand only when necessary. `&&` skips it when the left operand is `false`; `||` skips it when the left operand is `true`. Think of a guard who does not ask for a room number after finding that you have no pass. Rust does not guess which operation is dangerous: the written order matters.

`slots` counts available places, `tasks` counts tasks, and `fits` stores the decision. With zero places, `slots != 0` is false, so division never happens. Reversing these operands would attempt division by zero. The protection depends on writing the condition correctly.

## A common mistake

Numbers do not automatically become Booleans. This separate, deliberately invalid program produces E0308, a type mismatch.

<!-- error-code: E0308 -->
```rust,compile_fail
fn main() {
    let ready: bool = 1;
    println!("{ready}");
}
```

Choose a correction that matches your intention: use `true`, or compute a comparison such as `minutes > 0` after declaring `minutes`. Do not pick an arbitrary type merely to silence the compiler.

## Reference map

Data → comparisons → bool → combine conditions → check boundaries.

## Check your understanding

1. How do `=` and `==` differ?
2. Predict `false || true` and `!true`.
3. Will division happen when `slots = 0`? Why?

Check: assignment versus comparison; `true` and `false`; the left side of `&&` prevents evaluation of the division.

## Exercise

**Required.** Show a task only when it is unfinished, not archived, and takes at most 15 minutes. Declare `done = false`, `archived = false`, and `minutes: u32 = 15`. Print “Show: true”. Separately try `done = true`, `archived = true`, and `minutes = 16`: each should produce `false`. Reset the other values before each trial.

**Boundary check.** 15 minutes qualifies; 16 does not. Explain why you need `<=` rather than `<`.

## Hint

Translate each part of the rule: “not” becomes `!`; requiring all conditions uses `&&`. The answer below covers the initial case; verify the variations yourself.

## Reference answer after your own attempt

<!-- task-answer -->
```rust
fn main() {
    let done: bool = false;
    let archived: bool = false;
    let minutes: u32 = 15;
    let show = !done && !archived && minutes <= 15;
    println!("Show: {show}");
}
```

```text
Show: true
```

[Verification source](https://doc.rust-lang.org/reference/expressions/operator-expr.html#lazy-boolean-operators)

[Previous lesson](/read/rust-10-numbers?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-12-branches?lang=en)
