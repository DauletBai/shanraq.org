# Arrays

_Lead (summary):_ **Lesson 15. Process a fixed set of durations and check its boundaries.**

Prerequisites: lessons 1–14. If variables are unclear, revisit [lesson 9](/read/rust-09-variables?lang=en).

## Why this matters

Until now we used a separate variable for each task. An **array** holds values of one type with a length known at compile time. Picture a tray with three compartments: contents can change, but you cannot attach a fourth compartment to the same array. Every element must have the same type; similar meanings do not make different types interchangeable.

Store three task durations and compute their total with a loop from lesson 13. This is still a numeric summary, not a complete task with a title and status; structures will provide that later.

## Run the example

Use the `organizer` project from [lesson 6](/read/rust-06-cargo?lang=en), in the folder containing `Cargo.toml`. Save your previous work separately. Replace **all** of `src/main.rs` with the first example below, save it, and run `cargo run` in that folder’s terminal. As explained in lesson 6, this command builds and runs the program. No other files or dependencies change. Each subsequent complete example also replaces the entire file. Output blocks show only program output, without Cargo messages. Predict the output before running.

```rust
fn main() {
    let minutes: [u32; 3] = [10, 20, 15];
    println!("First: {}", minutes[0]);
    println!("Count: {}", minutes.len());
    let mut total: u32 = 0;
    for value in minutes {
        total = total + value;
    }
    println!("Total minutes: {total}");
}
```

```text
First: 10
Count: 3
Total minutes: 45
```

## Walkthrough

In the type **`[u32; 3]`**, `u32` is the element type and `3` is the element count. Square brackets enclose the array description; the inner semicolon separates type and length rather than ending a statement. **`[10, 20, 15]`** constructs an array value, separating elements with commas. Length is part of the type: `[u32; 3]` and `[u32; 4]` are different types.

An **index** is a position starting at zero. `minutes[0]` is the first element, `minutes[1]` the second, and `minutes[2]` the third. Do not confuse everyday ordinal numbers with code indices. An array of length 3 has no element at index 3.

`minutes.len()` reports the length. This is our first **method call**: the dot associates the operation `len` with the value `minutes`, and parentheses call it without extra arguments. We will write our own methods in lesson 28. `len` returns `usize`, a platform-dependent unsigned integer type used for sizes and indices. For example, `let index: usize = 1;` explicitly declares an index. An index is not a duration in minutes of type `u32`; do not add them without a meaningful reason.

`for value in minutes` visits the array values in order. `value` is the next number, not an index. Direct traversal is simpler here than `0..minutes.len()`. In this example `u32` numbers and the array containing them are copied as values. This does not explain ownership rules for strings; lessons 18–22 will address those.

## Changing elements and checking bounds

`let mut` allows element changes, not a change of array length. This separate program replaces the second element and validates a requested index **before** accessing it:

```rust
fn print_at(minutes: [u32; 3], index: usize) {
    if index < minutes.len() {
        println!("{}", minutes[index]);
    } else {
        println!("No such position");
    }
}

fn main() {
    let mut minutes: [u32; 3] = [10, 20, 15];
    minutes[1] = 25;
    print_at(minutes, 3);
}
```

```text
No such position
```

`print_at` is an ordinary function using lesson 14. Its first parameter is the entire `[u32; 3]` array, and its second is a `usize` index. `print_at(minutes, 3)` supplies them in order; the function returns `()` and checks the index itself. This can later validate an externally supplied index; for now the argument is fixed in the source.

With `index = 3`, access is skipped. Change the index to 1 and expect `25`. The condition must be strict: `index < minutes.len()`, not `<=`. An empty array has no valid index at all. Computing `minutes[index]` before checking bounds is too late.

An out-of-bounds access causes a **panic**: the program stops normal execution with an error message. If the invalid index is already obvious at compile time, the compiler may reject the code beforehand. Rust does not simply allow reading neighboring memory. An unexpected stop is still inconvenient, so validate the index first. Lesson 31 will introduce safer element access with `get` and `Option`.

Repeated elements have a shorter notation: `let minutes: [u32; 3] = [0; 3];` creates three zeros. In the **value** `[0; 3]`, the left part is the initial value and the right part is the count. In the **type** `[u32; 3]`, the left part is a type. Similar notation serves two different purposes.

## A common mistake

The declared length is three but the value contains two elements. This deliberately invalid example produces E0308:

<!-- error-code: E0308 -->
```rust,compile_fail
fn main() {
    let minutes: [u32; 3] = [10, 20];
    println!("{}", minutes.len());
}
```

Correct the length or the contents according to the actual requirement. If a user must add an arbitrary number of tasks, a fixed array is unsuitable; lesson 24 introduces the growable `Vec` list.

## Reference map

One type and length → values → traversal → function-based selection → accumulate results.

## Check your understanding

1. What is the last index of an array of length 4?
2. Does `mut` change the length?
3. How does the value in `for value in minutes` differ from an index?

Check: 3; no; `value` is a duration, an index is a position. You are ready for the next group if you can explain selection, repetition, and function calls in your solution. Lesson 16 on tests is still planned.

## Exercise

**Required.** Use `[0, 10, 25, 15]` of type `[u32; 4]`. Reuse `is_short` from lesson 14: durations 1 through 15 minutes count as short. With `for`, compute how many short tasks there are and their total duration. Expect “Short tasks: 2” and “Total minutes: 25”.

**Try your own data.** `[0, 16, 25, 30]` should give 0 and 0; `[1, 15, 2, 3]` should give 4 and 21. Our rule excludes zero-minute tasks. Change only the data, not the formulas.

## Hint

Declare two accumulators before the loop. Inside `if is_short(value)`, add 1 to the count and `value` to the total. Keep the previous lesson’s function outside `main`.

## Reference answer after your own attempt

<!-- task-answer -->
```rust
fn is_short(minutes: u32) -> bool {
    minutes > 0 && minutes <= 15
}

fn main() {
    let minutes: [u32; 4] = [0, 10, 25, 15];
    let mut count: u32 = 0;
    let mut total: u32 = 0;
    for value in minutes {
        if is_short(value) {
            count = count + 1;
            total = total + value;
        }
    }
    println!("Short tasks: {count}");
    println!("Total minutes: {total}");
}
```

```text
Short tasks: 2
Total minutes: 25
```

[Verification source](https://doc.rust-lang.org/book/ch03-02-data-types.html)

[Previous lesson](/read/rust-14-functions?lang=en) · [Contents](/course/rust?lang=en)
