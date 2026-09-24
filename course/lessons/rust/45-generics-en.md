# A generic type: one function for different lists

_Summary:_ **Lesson 45. Read and write a function that works with different element types.**

Prerequisites: [slices](/read/rust-23-slices?lang=en), [`Option`](/read/rust-31-option?lang=en), and [modules](/read/rust-44-modules?lang=en). Save the previous project and replace `src/main.rs` with the example. Tasks still exist only in memory.

## Familiar image and recall map

A cloakroom ticket points to a location: it may hold a coat or a bag, but “find the item at this number” is the same rule. A **generic type** lets us write that rule once. In `fn choose<T>(items: &[T], index: usize) -> Option<&T>`, `T` stands for an element type not yet specified. Angle brackets after the function name declare this type parameter; they are not a numeric comparison. When called with numbers, the compiler uses their type; with tasks, it uses `Task`. The image has a limit: the function cannot perform every operation on arbitrary `T`; here it only selects from a slice.

**Slice `&[T]` + position `usize` → `get` → `Some(&T)` or `None`.** `&[T]` borrows a slice of elements of one type; `usize` is an unsigned integer type for indices; its width depends on the target platform. `get` checks bounds and returns `None` for a missing position instead of panicking. `Option<&T>` returns a reference, so the function does not take an element from the array. `&` means reference; `Option` adds the possibility of absence. No explicit lifetime mark is needed here: the only borrowed input, `items`, determines what the returned reference is tied to. The full rule belongs to the later lifetime lesson. You have seen `<>` in `Vec<T>`; now you declare your own `T` rather than only reading an existing type.

## Run one algorithm with two types

Replace `src/main.rs` and run `cargo run`. No dependency is added. Cargo may print tool messages; only program output is shown below.

```rust
fn choose<T>(items: &[T], index: usize) -> Option<&T> {
    items.get(index)
}

struct Task {
    title: String,
}

fn main() {
    let numbers = [10, 20];
    let tasks = [Task {
        title: String::from("Read Rust"),
    }];
    if let Some(number) = choose(&numbers, 1) {
        println!("Number: {number}");
    }
    if let Some(task) = choose(&tasks, 0) {
        println!("Task: {}", task.title);
    }
    println!("No item: {}", choose(&numbers, 9).is_none());
}
```
```text
Number: 20
Task: Read Rust
No item: true
```

`choose(&numbers, 9)` returns `None`, and `is_none()` yields `true`. Writing `items[index]` instead would panic at an invalid position. This alone is not a lookup by a task's stable ID: a slice position can change after removal, as [lesson 38](/read/rust-38-ids?lang=en) showed.

## Recall without looking

1. What does `<T>` after `choose` declare?
2. Why return `Option<&T>` instead of `&T`?
3. Is position `0` the same as a task's stable ID?

## Exercise

**Required.** Write your own `fn last<T>(items: &[T]) -> Option<&T>` that returns the final item or `None` for an empty slice. Use it on `tasks` to print `Last: Read Rust`, then check an empty number array. Avoid indexing with `len() - 1`: that fails for empty arrays.

## Answers

A slice has a `last()` method that already returns `Option<&T>`. Call it inside your function and handle `Some` in `main`.

## Answer after your attempt

<!-- task-answer -->
```rust
fn choose<T>(items: &[T], index: usize) -> Option<&T> {
    items.get(index)
}

fn last<T>(items: &[T]) -> Option<&T> {
    items.last()
}

struct Task {
    title: String,
}

fn main() {
    let numbers = [10, 20];
    let tasks = [Task {
        title: String::from("Read Rust"),
    }];
    if let Some(number) = choose(&numbers, 1) {
        println!("Number: {number}");
    }
    if let Some(task) = choose(&tasks, 0) {
        println!("Task: {}", task.title);
    }
    println!("No item: {}", choose(&numbers, 9).is_none());
    if let Some(task) = last(&tasks) {
        println!("Last: {}", task.title);
    }
    let empty: [u32; 0] = [];
    println!("Empty list: {}", last(&empty).is_none());
}
```
```text
Number: 20
Task: Read Rust
No item: true
Last: Read Rust
Empty list: true
```

## After checking

For an empty array, `last(&[])` yields `None`; give that empty array a type, for example `let empty: [u32; 0] = [];`. [Official chapter on generic types](https://doc.rust-lang.org/book/ch10-01-syntax.html). Revisit [lesson 31](/read/rust-31-option?lang=en) if `Some` and `None` are unclear.

[Previous lesson](/read/rust-44-modules?lang=en) · [Contents](/course/rust?lang=en)
