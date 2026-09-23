# Repeating actions

_Lead (summary):_ **Lesson 13. Repeat work without missing boundaries or running forever.**

Prerequisites: lessons 1–12. If variables are unclear, revisit [lesson 9](/read/rust-09-variables?lang=en).

## Why this matters

So far, each reached line ran once. A **loop** repeats a block; one pass is an **iteration**. Think of making a daily calendar entry: the rule stays the same while the days change. The analogy has a limit: a calendar has a last page, but a program without an exit can run forever.

Plan ten minutes of tasks on each of three days. The program calculates the plan; it does not wait three actual days. Recall `mut` from lesson 9 and conditions from lesson 12.

## Run the example

Use the `organizer` project from [lesson 6](/read/rust-06-cargo?lang=en), in the folder containing `Cargo.toml`. Save your previous work separately. Replace **all** of `src/main.rs` with the first example below, save it, and run `cargo run` in that folder’s terminal. As explained in lesson 6, this command builds and runs the program. No other files or dependencies change. Each subsequent complete example also replaces the entire file. Output blocks show only program output, without Cargo messages. Predict the output before running.

Read the new form before running it: `for day in 1..=3` visits the numbers from 1 through 3 **including** 3. `for` begins repetition, `day` names the current number, `in` introduces its source, and `..=` includes the right endpoint. Picture a finger visiting three calendar squares in order; the program does not wait for real days. Each execution of the `{ ... }` block is one iteration.

```rust
fn main() {
    let mut total_minutes: u32 = 0;
    for day in 1..=3 {
        println!("Day: {day}");
        total_minutes = total_minutes + 10;
    }
    println!("Total minutes: {total_minutes}");
}
```

```text
Day: 1
Day: 2
Day: 3
Total minutes: 30
```

## Walkthrough

`for` visits values, and `in` introduces their source. A **range** describes values using boundaries. `1..=3` includes 1, 2, 3; `1..3` includes only 1 and 2. The `=` in `..=` includes the right boundary; it is not assignment here. For these ascending integer ranges, `3..3` and `3..1` are empty: this is not countdown notation.

`day` holds the current value within the loop block; you do not increment it manually. `total_minutes` is an **accumulator**, a variable collecting the results of iterations. Declare it before the loop or it would start over every time. `total_minutes = total_minutes + 10` reads the old value, adds ten, and stores the new value. The successive totals are 10, 20, 30. The outer `println!` runs only once.

## When to use while or loop

`while` checks its condition **before** each iteration. If it starts false, the body never runs. `loop` repeats without a condition. `break` immediately leaves the innermost loop; `continue` skips the rest of the current iteration. A `while` then checks its condition again, while a `for` moves to its next value.

The next example is a separate program. Temporarily replace all of `src/main.rs`; do not add a second `main`.

```rust
fn main() {
    let mut remaining: u32 = 3;
    while remaining > 0 {
        println!("Remaining: {remaining}");
        remaining = remaining - 1;
    }
    let mut attempt: u32 = 0;
    loop {
        attempt = attempt + 1;
        if attempt == 2 {
            continue;
        }
        println!("Attempt: {attempt}");
        if attempt == 3 {
            break;
        }
    }
}
```

```text
Remaining: 3
Remaining: 2
Remaining: 1
Attempt: 1
Attempt: 3
```

In `while`, reducing `remaining` to zero is what eventually makes the condition false. In `loop`, `attempt` is incremented **before** a possible `continue`: iteration two prints nothing, and iteration three reaches `break`. A `for` would be simpler for this fixed sequence; `loop` is included to demonstrate control of repetition.

## Avoiding an endless loop

Remove the decrement of `remaining` and the `while` never finishes. Put the increment of `attempt` after `continue` and it may become stuck at one value. Before running a loop, identify the action that brings it closer to an exit. If you accidentally start endless printing, press **Ctrl+C in the terminal** to interrupt the program. This stops execution but does not repair the code.

An **off-by-one error** uses the wrong boundary: `1..3` gives two iterations, not three. Write down the first and last values before editing.

## Reference map

Initial state → next value or condition → body → change → completion.

## Check your understanding

1. How many iterations do `0..3` and `0..=3` give?
2. How do `break` and `continue` differ?
3. What if the accumulator is declared inside the loop?

Check: 3 and 4; leave the loop versus move to its next iteration; the accumulator is recreated and does not retain the overall total.

## Exercise

**Required.** Plan 15 minutes on days 1 through 5, skipping day 3 with `continue`. Print only days 1, 2, 4, 5 and “Total minutes: 60”. Accumulate the total rather than printing a hard-coded answer.

**Boundary check.** Temporarily use `1..1`: no days should print and the total should be zero. Restore `1..=5`. Explain why the day check belongs before printing and adding.

## Hint

Declare the accumulator before `for`. Start the loop body by checking `day == 3`. The reference does not count the skipped day.

## Reference answer after your own attempt

<!-- task-answer -->
```rust
fn main() {
    let mut total_minutes: u32 = 0;
    for day in 1..=5 {
        if day == 3 {
            continue;
        }
        println!("Day: {day}");
        total_minutes = total_minutes + 15;
    }
    println!("Total minutes: {total_minutes}");
}
```

```text
Day: 1
Day: 2
Day: 4
Day: 5
Total minutes: 60
```

[Verification source](https://doc.rust-lang.org/book/ch03-05-control-flow.html)

[Previous lesson](/read/rust-12-branches?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-14-functions?lang=en)
