# Methods: actions beside Task

_Lead (summary):_ **Lesson 28. Create a task, read its state and finish it through methods.**

Before this lesson: lessons 14, 20, 21 and 27. If the difference between `&` and `&mut` is unclear, revisit [lesson 21](/read/rust-21-mutable-borrow?lang=en).

## Familiar image and recall map

A task card has filled-in spaces, but we also need actions: read its state and mark it finished. A **method** is a function associated with a type and called through a value of that type. An `impl Task { ... }` block groups these actions beside `Task`; `impl` means implementation. The image has a limit: a paper card performs no action by itself. The method's code runs only when called.

**`struct Task` → `impl Task` → `Task::new(...)` → `task.is_done()` → `task.finish()`.** Retell the map. A function such as `new` inside `impl` without a `self` parameter belongs to the type and is called as `Task::new(...)`. The `::` path goes from type to function, as in `String::from(...)`. This is an **associated function**, rather than a method on one instance. A method has `self` as its first parameter: `&self` borrows the task for reading; `&mut self` borrows it for changing. In `task.is_done()` and `task.finish()`, Rust passes `task` as that first parameter automatically.

## Read and change a task

Save the old `src/main.rs` in `organizer`, replace the whole file, and run `cargo run` beside `Cargo.toml`. No other file or dependency is needed. Predict the output first.

`new` takes a title and minutes, creates a `Task` with `done: false`, and returns it. In `Task { title, minutes, done: false }`, the two names without `:` are field shorthand: they mean `title: title, minutes: minutes` because the parameter and field names match. We pass known suitable values; this function does not validate user input yet. `is_done` returns a copyable `bool` without changing the task. `finish` changes a field, so it needs `&mut self` and a mutable `task` variable.

```rust
struct Task {
    title: String,
    minutes: u32,
    done: bool,
}

impl Task {
    fn new(title: String, minutes: u32) -> Task {
        Task {
            title,
            minutes,
            done: false,
        }
    }

    fn is_done(&self) -> bool {
        self.done
    }

    fn finish(&mut self) {
        self.done = true;
    }
}

fn main() {
    let mut task = Task::new(String::from("Reading"), 20);
    println!(
        "{}: {} min, done {}",
        task.title,
        task.minutes,
        task.is_done()
    );
    task.finish();
    println!("After: {}", task.is_done());
}
```
```text
Reading: 20 min, done false
After: true
```

Ownership of `title` moves into the new instance. `finish` has no return arrow, so it returns `()` from lesson 14. Its `&mut self` borrow ends after the call, and `is_done` can read the task again.

## When `&mut self` cannot be borrowed

Without `mut` on the variable, a changing method cannot get a mutable reference to it. This **separate** example intentionally fails with E0596:

<!-- error-code: E0596 -->
```rust,compile_fail
struct Task { done: bool }
impl Task {
    fn finish(&mut self) { self.done = true; }
}
fn main() {
    let task = Task { done: false };
    task.finish();
}
```

Adding `mut` to `task` allows the call. Methods still follow ownership rules: `&self` cannot change fields, and `&mut self` needs exclusive access during the call. This `new` does not reject zero minutes. We will add full value checking and an error message after learning `Result`.

## Check your understanding

1. Why do we call `new` through `Task::` but `finish` through `task.`?
2. Which method requires a mutable variable?
3. Why can `is_done()` run after `finish()`?

Check: `new` takes no instance; `finish` borrows `&mut self`; that temporary borrow has ended. Repeat the recall map without the page.

## Exercise

**Required.** Declare `Task { title: String, done: bool }`. In `impl Task`, write `new(title: String) -> Task`, `is_done(&self) -> bool` and `finish(&mut self)`. Create a mutable task named `Walk`, print `Walk: false`, finish it, then print `Walk: true`. Obtain both state values through the method.

**Change.** Remove `mut` from the variable, read the error, and restore `mut`.

## Hint

Have `new` create `Task { title, done: false }`. Have `finish` write `self.done = true;`. Call `task.is_done()` both before and after `task.finish()`.

## Reference answer after trying

<!-- task-answer -->
```rust
struct Task {
    title: String,
    done: bool,
}

impl Task {
    fn new(title: String) -> Task {
        Task { title, done: false }
    }

    fn is_done(&self) -> bool {
        self.done
    }

    fn finish(&mut self) {
        self.done = true;
    }
}

fn main() {
    let mut task = Task::new(String::from("Walk"));
    println!("{}: {}", task.title, task.is_done());
    task.finish();
    println!("{}: {}", task.title, task.is_done());
}
```
```text
Walk: false
Walk: true
```

The data and the actions that belong to it now sit together. [Official methods chapter](https://doc.rust-lang.org/book/ch05-03-method-syntax.html).

[Previous lesson](/read/rust-27-structs?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-29-enums?lang=en)
