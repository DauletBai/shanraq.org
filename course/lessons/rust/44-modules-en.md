# Modules: separate logic from startup

_Summary:_ **Lesson 44. Split code into files and expose only the names callers need.**

Prerequisites: [structs](/read/rust-27-structs?lang=en), [methods](/read/rust-28-methods?lang=en), and the [organizer](/read/rust-40-memory-checkpoint?lang=en). Start from a copy: this short example replaces `src/main.rs` so the visibility boundary is easy to see. File persistence for tasks has not started yet.

## Familiar image and recall map

A workshop has a shared counter and a closed tool cabinet. Only the tools visitors need go on the counter. A **module** (`mod`) groups names and sets a visibility boundary. `pub` exposes a name across that boundary. The image stops there: code privacy helps define a clear interface; it does not encrypt data or secure a computer against its user.

**`main` → `use crate::organizer::Task` → public `Task::new` and `task.print` → private fields.** `crate` names the root of the current compilation unit. `::` joins segments of a path to a name. `use` lets us write short `Task` rather than the full path. `pub struct Task` exposes the type, but fields `id` and `title` remain private without their own `pub`. Thus `main` can call the public constructor and method but cannot write `task.title = ...`. The struct can check changes through its methods. `Self` inside `impl Task` means `Task`, as in lesson 28.

## First, one file

Replace `src/main.rs` and run `cargo run`. Cargo messages are separate from the program output. No dependency is added.

```rust
mod organizer {
    pub struct Task {
        id: u32,
        title: String,
    }
    impl Task {
        pub fn new(id: u32, title: &str) -> Self {
            Self {
                id,
                title: String::from(title),
            }
        }
        pub fn print(&self) {
            println!("Task {}: {}", self.id, self.title);
        }
    }
}

use crate::organizer::Task;

fn main() {
    let task = Task::new(1, "Read Rust");
    task.print();
}
```
```text
Task 1: Read Rust
```

## Then, two files

Move the code **inside** `mod organizer { ... }` into a new `src/organizer.rs`, without the outer braces. In `src/main.rs`, replace the whole module block with `mod organizer;` and keep `use crate::organizer::Task;` and `fn main()`. The semicolon after `mod organizer` tells Rust to find the body in a separate file. `cargo run` should print the same line. A `private field` error means `main` tried to access a field directly; use the public method. For `file not found`, check the name and location `src/organizer.rs`.

`main.rs` is the root of a Cargo package's executable program. `lib.rs` can be the root of a separate library part in the same package; other programs can use its functions and it can be tested independently. You do not need `lib.rs` in this exercise. The pair `main.rs` and `organizer.rs` already separates startup from logic. We can move shared logic into a library when we are ready to define its public behavior.

## Recall without looking

1. Why does `pub struct Task` not automatically expose `title`?
2. What does the semicolon in `mod organizer;` mean?
3. How do `main.rs` and `lib.rs` differ?

## Exercise

**Required.** Add a public `id(&self) -> u32` method and use it to print `Number: 1`; keep the field `id` private. Make this change in the two-file version too and run `cargo run`. Do not solve this by adding `pub` to the field.

## Answers

A method inside `impl Task` may read `self.id`; call `task.id()` outside. After splitting files, keep that method in `src/organizer.rs`.

## Answer after your attempt

The complete one-file answer follows. For two files, apply the move described above.

<!-- task-answer -->
```rust
mod organizer {
    pub struct Task {
        id: u32,
        title: String,
    }
    impl Task {
        pub fn new(id: u32, title: &str) -> Self {
            Self {
                id,
                title: String::from(title),
            }
        }
        pub fn print(&self) {
            println!("Task {}: {}", self.id, self.title);
        }
        pub fn id(&self) -> u32 {
            self.id
        }
    }
}

use crate::organizer::Task;

fn main() {
    let task = Task::new(1, "Read Rust");
    task.print();
    println!("Number: {}", task.id());
}
```
```text
Task 1: Read Rust
Number: 1
```

## After checking

[Official explanation of modules](https://doc.rust-lang.org/book/ch07-02-defining-modules-to-control-scope-and-privacy.html) · [separating modules into files](https://doc.rust-lang.org/book/ch07-05-separating-modules-into-different-files.html). Revisit [lesson 28](/read/rust-28-methods?lang=en) if `impl` is unclear.

[Previous lesson](/read/rust-43-maps-sets?lang=en) · [Contents](/course/rust?lang=en)
