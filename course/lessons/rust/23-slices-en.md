# Slices: part of data without a copy

_Lead (summary):_ **Lesson 23. Pass part of a numeric collection to a function without building another array.**

Prerequisites: lessons 13, 15, and 20–22. If the boundaries of `1..3` are unclear, revisit [lesson 13](/read/rust-13-loops?lang=en).

## A familiar image and recall map

You can show a colleague two adjacent lines on your list without copying the whole sheet. A **slice** is a view of a contiguous part of a sequence without owning its elements. The image has a limit: Rust ties a slice's usable period to its data and checks boundaries; paper does neither.

**Owner → boundaries → read-only slice → function → owner remains.** `&[u32]` is a reference to a sequence of `u32` numbers with a known length, but the type does not state the exact count. Unlike `[u32; 4]` from lesson 15, its length is not written into the type. `&minutes[1..3]` takes elements at indices 1 and 2: `..` from lesson 13 **excludes** the right boundary. Both boundaries must be valid.

## Sum the middle

In `organizer`, save the previous `src/main.rs` separately and replace the whole file with this example. Run `cargo run` beside `Cargo.toml`. No other files or dependencies change. Predict the result; Cargo messages are omitted below.

Read the new parts before the code. `middle` points to just two numbers in the original array. `values: &[u32]` lets `sum` accept such a slice. `0..values.len()` produces indices from zero up to, **but excluding**, the length, so each `values[index]` is valid. A `u32` number is copied as a value; the whole array is not copied for this call.

```rust
fn sum(values: &[u32]) -> u32 {
    let mut total: u32 = 0;
    for index in 0..values.len() {
        total = total + values[index];
    }
    total
}

fn main() {
    let minutes: [u32; 4] = [10, 20, 15, 5];
    let middle = &minutes[1..3];
    println!("Middle: {}", sum(middle));
    println!("First: {}", minutes[0]);
}
```

```text
Middle: 35
First: 10
```

`middle` covers 20 and 15. It cannot outlive `minutes`. Passing `middle` does not transfer ownership of the array, so `minutes[0]` is still available afterward. `&minutes[0..minutes.len()]` gives a slice of the whole array; the function need not know its length beforehand.

## Bounds and use period

If the right boundary exceeds the length, creating a slice stops the program with a **panic**, as with array bounds in lesson 15. For a four-element array, `&minutes[4..4]` is a valid empty slice, but `&minutes[4..5]` is invalid. Passing the empty slice to `sum` leaves the initial zero as the answer.

A slice is a reference, so it cannot outlive its owner. This **separate** program intentionally fails with E0597:

<!-- error-code: E0597 -->
```rust,compile_fail
fn main() {
    let view: &[u32];
    {
        let values = [10, 20];
        view = &values[0..2];
    }
    println!("{}", view[0]);
}
```

`values` ends in the inner block while `view` is needed outside. Move the owner to the outer block or finish using the slice inside. Renaming the slice alone cannot repair the problem.

`&str` from lesson 22 is also a slice, but of **text**. Text slice boundaries use bytes and must align with UTF-8 character boundaries. Do not transfer array indices straight onto a word; lesson 25 explains this in detail.

## Check your understanding

1. Which elements does `&minutes[1..3]` include?
2. Why does `sum` accept `&[u32]` instead of `[u32; 4]`?
3. Can a slice be used after its array disappears?

Check: indices 1 and 2; the function accepts different slice lengths without copying an array; no. Rebuild the recall map.

## Exercise

**Required.** With `[5, 10, 15, 20]`, call your `sum(values: &[u32])` for the first two and last two elements. Print `First pair: 15`, `Last pair: 35`, then `Element count: 4` from the original array. Do not create two more arrays or print precomputed sums.

**Boundary.** Try an empty slice `&minutes[0..0]`: its sum should be zero. Then restore the required solution.

## Hint

Use `&minutes[0..2]` and `&minutes[2..4]` with the same function. Inside, walk through `0..values.len()`.

## Reference answer after your attempt

<!-- task-answer -->
```rust
fn sum(values: &[u32]) -> u32 {
    let mut total: u32 = 0;
    for index in 0..values.len() {
        total = total + values[index];
    }
    total
}

fn main() {
    let minutes: [u32; 4] = [5, 10, 15, 20];
    println!("First pair: {}", sum(&minutes[0..2]));
    println!("Last pair: {}", sum(&minutes[2..4]));
    println!("Element count: {}", minutes.len());
}
```

```text
First pair: 15
Last pair: 35
Element count: 4
```

One function processes two parts of the array while the owner remains available. [Official guide to slices](https://doc.rust-lang.org/book/ch04-03-slices.html).

[Previous lesson](/read/rust-22-strings?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-24-vectors?lang=en)
