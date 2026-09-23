# Practice: a small summary

_Lead (summary):_ **Lesson 17. Build and test a task summary independently using familiar syntax.**

Prerequisites: lessons 1–16. This lesson introduces no new Rust syntax. Revisit [lesson 16](/read/rust-16-first-test?lang=en) if a test will not run, or [lessons 13](/read/rust-13-loops?lang=en) and [15](/read/rust-15-arrays?lang=en) if you confuse array values with positions.

## Reconstruct the map first

Close earlier lessons and draw: **data → selection rule → visit each value → counters → display → tests**. Explain `u32`, `[u32; 4]`, `for`, `if`, `fn`, and `#[test]` without looking. Then check your account against the earlier pages. Recognizing code is easier than recreating it; this exercise asks you to recreate it.

The numbers represent minutes for four tasks. Zero means no task; a short task lasts 1–15 minutes inclusive. At this stage names and user input are absent. The values are written in the source file, so closing the program does not save a changed task list.

## Warm-up: read a program

Use the `organizer` project from lesson 6, in the folder with `Cargo.toml`. Save your old `src/main.rs` separately and replace the whole file with this example. Run `cargo run`, which builds and starts `main`. The output block excludes Cargo messages. Predict the count before running.

```rust
fn count_short(minutes: [u32; 4]) -> u32 {
    let mut count = 0;
    for value in minutes {
        if value > 0 && value <= 15 {
            count = count + 1;
        }
    }
    count
}

fn main() {
    let tasks = [0, 8, 20, 15];
    println!("Short tasks: {}", count_short(tasks));
}
```

```text
Short tasks: 2
```

The function returns a number; `main` handles the message. Why are 0 and 20 excluded? What happens if 15 becomes 16? Predict, then edit and run. The numeric array can be copied when passed to this function; do not extend that rule to strings without studying ownership.

## Exercise: independent checkpoint

**Required.** Without copying the warm-up, write `is_short(minutes: u32) -> bool` and `summarize(minutes: [u32; 4]) -> [u32; 2]`. The first number in the returned two-element array is the count of short tasks; the second is their total minutes. For `[0, 8, 20, 15]`, `main` must print `Short tasks: 2` and `Minutes: 23`. Add tests for `[0, 0, 0, 0]` → `[0, 0]`, `[1, 15, 16, 0]` → `[2, 16]`, and the ordinary case. Run `cargo test`, then `cargo run`.

Follow this sequence: write expected answers on paper, recreate the functions, check the first compiler or test error in your own words, and rerun. Change only the input array to `[15, 15, 15, 15]`; expect 4 and 60. Tomorrow, try again from an empty file. If the next step is unclear, return to the map and identify the missing link.

## Hint

Start two counters at zero. Visit each array value. Inside `if is_short(value)`, increase the count and add the minutes. The final expression of the function is `[count, total]`. In `main`, use positions 0 and 1 of the result, as explained in lesson 15.

## Reference answer after your attempt

<!-- task-answer -->
```rust
fn is_short(minutes: u32) -> bool {
    minutes > 0 && minutes <= 15
}

fn summarize(minutes: [u32; 4]) -> [u32; 2] {
    let mut count = 0;
    let mut total = 0;
    for value in minutes {
        if is_short(value) {
            count = count + 1;
            total = total + value;
        }
    }
    [count, total]
}

fn main() {
    let result = summarize([0, 8, 20, 15]);
    println!("Short tasks: {}", result[0]);
    println!("Minutes: {}", result[1]);
}

#[test]
fn no_tasks() {
    assert_eq!(summarize([0, 0, 0, 0]), [0, 0]);
}

#[test]
fn boundaries() {
    assert_eq!(summarize([1, 15, 16, 0]), [2, 16]);
}

#[test]
fn ordinary_case() {
    assert_eq!(summarize([0, 8, 20, 15]), [2, 23]);
}
```

```text
Short tasks: 2
Minutes: 23
```

The result array is convenient but asks us to remember what positions 0 and 1 mean. Later a structure will give these values names. If you needed the answer to finish, repeat the task without it before moving on.

[Previous lesson](/read/rust-16-first-test?lang=en) · [Contents](/course/rust?lang=en)
