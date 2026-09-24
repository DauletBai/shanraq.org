# Traits: a shared behavior contract

_Summary (summary):_ **Lesson 46. Give two types a shared behavior and constrain a generic function.**

First review [generics](/read/rust-45-generics?lang=en) and [methods](/read/rust-28-methods?lang=en). Save your earlier project; this example still keeps data only in memory.

## Familiar image and recall map

Two people at a club can introduce themselves even if they store their details differently. An instruction to “introduce yourself” resembles a **trait**: it names required behavior. Each participant supplies the action. The image has a limit: the compiler checks the method signature, not whether the text is useful to a person.

**trait — required method → impl — behavior for a type → T: Summary — condition on the shared function.** `trait Summary` declares `summary(&self) -> String`; the semicolon means there is no body here. `impl Summary for Task` supplies the body for `Task`. In `show<T: Summary>`, the colon sets a trait bound: `show` accepts any `T` that fulfills this contract. `&T` only borrows the value. `clone()` makes an owned string to return; the task keeps its title.

## Use one function for a task

```rust
trait Summary {
    fn summary(&self) -> String;
}

struct Task {
    title: String,
}

impl Summary for Task {
    fn summary(&self) -> String {
        self.title.clone()
    }
}

fn show<T: Summary>(item: &T) {
    println!("{}", item.summary());
}

fn main() {
    let task = Task {
        title: String::from("Read Rust"),
    };
    show(&task);
}
```
```text
Read Rust
```

`show(&task)` uses the implementation for `Task`. Remove `impl Summary for Task` and the compiler rejects the call. For some standard traits, an attribute such as `#[derive(Debug)]` asks the compiler to generate an implementation. It works only with supported traits and suitable fields; it cannot automatically implement an arbitrary contract. Here we write the behavior because the output is our choice.

## How derive works

`Debug` exposes fields for a developer; it is not finished user-facing text. `#[derive(Debug)]` is an attribute before the struct declaration that asks the compiler to implement a standard trait. `{:?}` prints that debug representation.

```rust
#[derive(Debug)]
struct Task {
    title: String,
}

fn main() {
    let task = Task {
        title: String::from("Read Rust"),
    };
    println!("{task:?}");
}
```
```text
Task { title: "Read Rust" }
```

## Recall without looking

1. What does `trait Summary` require?
2. What does `impl Summary for Task` provide?
3. Why does `show<T: Summary>` reject a type without the trait?

## Exercise

**Required.** Add `struct Note { text: String }`, implement `Summary` for it, create a note “Buy a book”, and pass it to the same `show` function. Predict both output lines first.

## Answers

Follow the `impl Summary for Task` pattern, changing `Task` to `Note` and `title` to `text`. Keep the note’s text in the note.

<!-- task-answer -->
```rust
trait Summary {
    fn summary(&self) -> String;
}

struct Task {
    title: String,
}

impl Summary for Task {
    fn summary(&self) -> String {
        self.title.clone()
    }
}

struct Note {
    text: String,
}

impl Summary for Note {
    fn summary(&self) -> String {
        self.text.clone()
    }
}

fn show<T: Summary>(item: &T) {
    println!("{}", item.summary());
}

fn main() {
    let task = Task {
        title: String::from("Read Rust"),
    };
    let note = Note {
        text: String::from("Buy a book"),
    };
    show(&task);
    show(&note);
}
```
```text
Read Rust
Buy a book
```

## After checking

Both lines use the same `show`; the type determines the method body. If `T` is unclear, revisit [lesson 45](/read/rust-45-generics?lang=en). [Official traits chapter](https://doc.rust-lang.org/book/ch10-02-traits.html).

[Previous lesson](/read/rust-45-generics?lang=en) · [Contents](/course/rust?lang=en)
