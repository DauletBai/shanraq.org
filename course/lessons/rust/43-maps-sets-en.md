# Map and set: count task tags

_Summary:_ **Lesson 43. Build tag statistics without relying on storage order.**

Prerequisites: [tuples](/read/rust-26-tuples?lang=en), [iterators](/read/rust-41-iterators?lang=en), and [`collect`](/read/rust-42-closures?lang=en). Save the previous file and replace `src/main.rs` with this focused example. These are still memory values, not saved organizer data.

## Familiar image and recall map

One warehouse ledger records how many boxes belong to each category; another lists only the distinct categories encountered. **`HashMap<K, V>`** is a key → value map, here tag → count. **`HashSet<T>`** is a set of distinct values, here tags. Angle brackets name the key and value types, as they did in `Vec<T>`. A hash is a code computed from a key to help choose a storage location, like a warehouse section number. Different keys can land in the same section; the collection handles that. The image does not define reading order: iteration order of `HashMap` and `HashSet` is not guaranteed and may change between runs.

**Tag → read old count → add one → store; tag → insert into set; display → sort first.** `use std::collections::{HashMap, HashSet}` brings two names from one module into this file; the braces after `::` list those names. `get(tag)` yields `Some(&number)` or `None`; `*value` reads the number through that reference. `insert` replaces the value for an existing key or adds a new key. Inserting a tag again into a set does not make a duplicate. `into_iter().collect()` moves the distinct tags into `Vec<&str>`, and `sort()` gives predictable output. This is Rust string order, not language-aware dictionary order. `counts["work"]` is safe only because this example already inserted the key; an unknown key would panic. Use `get` and handle `None` for user input.

## Run the count

Replace `src/main.rs` and run `cargo run`. `HashMap` and `HashSet` are in Rust's standard library; no dependency is added. The lines below are program output, separate from Cargo messages.

```rust
use std::collections::{HashMap, HashSet};

fn main() {
    let tags = ["work", "home", "work"];
    let mut counts: HashMap<&str, usize> = HashMap::new();
    let mut unique: HashSet<&str> = HashSet::new();
    for tag in tags {
        let old = match counts.get(tag) {
            Some(value) => *value,
            None => 0,
        };
        counts.insert(tag, old + 1);
        unique.insert(tag);
    }
    let mut sorted: Vec<&str> = unique.into_iter().collect();
    sorted.sort();
    println!("Tags: {:?}", sorted);
    println!("Work: {}", counts["work"]);
}
```
```text
Tags: ["home", "work"]
Work: 2
```

Two occurrences of `work` produce count 2 but just one distinct tag. Sorting changes only the list used for display, not the counts. This will support later organizer summaries, but tasks are still lost when the process ends.

## Recall without looking

1. Why is there only one copy of a tag in the set?
2. What does `None` from `counts.get(tag)` mean?
3. Why sort before printing tags?

## Exercise

**Required.** Print the number of tasks with tag `home` as `Home: 1`. Do not rely on `HashMap` iteration order. Then add one more `home` to the input array and predict both the new count and the set's contents.

## Answers

The example inserted this key, so `counts["home"]` is safe after the existing output. For a new user-supplied key, call `get` first.

## Answer after your attempt

<!-- task-answer -->
```rust
use std::collections::{HashMap, HashSet};

fn main() {
    let tags = ["work", "home", "work"];
    let mut counts: HashMap<&str, usize> = HashMap::new();
    let mut unique: HashSet<&str> = HashSet::new();
    for tag in tags {
        let old = match counts.get(tag) {
            Some(value) => *value,
            None => 0,
        };
        counts.insert(tag, old + 1);
        unique.insert(tag);
    }
    let mut sorted: Vec<&str> = unique.into_iter().collect();
    sorted.sort();
    println!("Tags: {:?}", sorted);
    println!("Work: {}", counts["work"]);
    println!("Home: {}", counts["home"]);
}
```
```text
Tags: ["home", "work"]
Work: 2
Home: 1
```

## After checking

With a second `home`, its count becomes 2 while the distinct tag list stays the same. [Official `HashMap` chapter](https://doc.rust-lang.org/book/ch08-03-hash-maps.html) · [collections overview](https://doc.rust-lang.org/std/collections/). Revisit [lesson 42](/read/rust-42-closures?lang=en) if `collect` is unclear.

[Previous lesson](/read/rust-42-closures?lang=en) · [Contents](/course/rust?lang=en)
