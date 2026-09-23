# A growable list: Vec

_Lead (summary):_ **Lesson 24. Add and remove task titles without fixing the number of slots in advance.**

Prerequisites: lessons 15 and 19–23. If a slice and its owner blur together, revisit [lesson 23](/read/rust-23-slices?lang=en).

## A familiar image and recall map

The three-slot tray from lesson 15 cannot gain a fourth slot. You can continue a written list with more lines. **`Vec`** is an owned, growable list of values of one type; it tracks its length and obtains more room when needed. The image has a limit: growth may move storage, and removal from the middle shifts later elements. Indices are not permanent task identifiers.

**Empty list → `push` → length → read through `&` → `remove` → new positions.** Recall the map before the exercise. In `Vec<String>`, angle brackets `<...>` name the **type of each element**, here `String`; they are not comparisons. This is a type with a parameter: `Vec` can hold another element type, but one list holds elements of one type. `Vec::new()` creates an empty list. `push` appends an element; `.len()` gives its count.

## From empty to two tasks

In `organizer`, save the previous `src/main.rs` separately, replace the whole file with this example, and run `cargo run` beside `Cargo.toml`. No other files or dependencies change. Predict program output without Cargo messages.

`for title in &titles` borrows the list for reading and supplies a reference to each `String` in turn. `&` here does not give ownership of the list to the loop; the list remains usable afterward. `String::from` makes owned titles; `push` transfers them into the list.

```rust
fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Reading"));
    titles.push(String::from("Rest"));
    println!("Task count: {}", titles.len());
    for title in &titles {
        println!("Task: {title}");
    }
}
```

```text
Task count: 2
Task: Reading
Task: Rest
```

`mut` is needed to change the list. Like an array, `Vec<String>` stores ordered elements of one type and allows repeated titles. Unlike an array, its length can change while the program runs. This `Vec` remains only in memory; closing the program does not save its tasks. Files come later.

## Removal and ownership

`remove(1)` removes the item at index 1 and **returns** its `String` to the caller. Later elements shift left. The index must be less than `.len()`, or the program panics. Check the boundary before removing an index supplied by a user; here we chose the index ourselves. The list still owns its remaining elements.

If a title was already held in a variable, `push(title)` moves its ownership into the list. This separate program intentionally fails with E0382:

<!-- error-code: E0382 -->
```rust,compile_fail
fn main() {
    let title = String::from("Reading");
    let mut titles: Vec<String> = Vec::new();
    titles.push(title);
    println!("{title}");
}
```

The list now owns that `String`; the old name `title` cannot use it. To read a title from the list, traverse `&titles`; avoid copying all strings without reason. You cannot keep using a reference to an element while changing the list in a way that could relocate storage; borrowing rules from lesson 21 prevent this.

## Check your understanding

1. How does `Vec<String>` differ from `[String; 3]`?
2. What does `String` between `<` and `>` mean?
3. Why can a removed task's index not be its permanent ID?

Check: growable versus fixed length; each element's type; removal shifts positions. Do not confuse a list's length with the last record's identifier.

## Exercise

**Required.** Make an empty `Vec<String>` and add `Reading`, `Walk`, `Rest`. Remove the element at index 1 with `remove`; print `Removed: Walk`, then the remaining titles in order and `Remaining: 2`. Use `&titles` for traversal. Do not replace traversal with hard-coded output lines.

**Boundary.** What happens with `remove(3)` on a list of length 3? Why would `<= titles.len()` be the wrong validity check?

## Hint

Store `remove(1)` in `removed`. Traverse with `for title in &titles` after removal. Get the count from `.len()`.

## Reference answer after your attempt

<!-- task-answer -->
```rust
fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Reading"));
    titles.push(String::from("Walk"));
    titles.push(String::from("Rest"));
    let removed = titles.remove(1);
    println!("Removed: {removed}");
    for title in &titles {
        println!("Task: {title}");
    }
    println!("Remaining: {}", titles.len());
}
```

```text
Removed: Walk
Task: Reading
Task: Rest
Remaining: 2
```

`removed` owns the removed value; the list owns the others. [Official guide to Vec](https://doc.rust-lang.org/book/ch08-01-vectors.html).

[Previous lesson](/read/rust-23-slices?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-25-unicode?lang=en)
