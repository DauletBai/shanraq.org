# External libraries and dependencies

_Summary (summary):_ **Lesson 48. Add a justified library and inspect what the project includes.**

First review [Cargo](/read/rust-06-cargo?lang=en) and [modules](/read/rust-44-modules?lang=en). Keep your organizer project separate: this is a small practice copy. The first run needs no network.

## Familiar image and recall map

A ready-made bicycle part saves work, but you check its size, maker, and fitting instructions. A **dependency** is a library included in our project. A **crate** is Rust’s unit of code and compilation; `crates.io` is a public catalog of crates. The image has a limit: library code runs with our program’s capabilities, while a catalog listing and a version number do not establish quality or safety.

**need → choose library and license → record in Cargo.toml → exact selection in Cargo.lock → build and verify.** `Cargo.toml` declares allowed dependencies. `Cargo.lock` records the exact versions selected for a reproducible application build. `version = "1.2"` is a condition for compatible versions, not necessarily the exact version downloaded; Cargo defines the compatibility rules. The `license` field on the crate page and in its source describes usage terms: read them before distributing your program. `cargo add` writes the dependency entry; `cargo tree` shows the dependency chain. A first download from crates.io needs a network. We will start with a tiny local library, which demonstrates a real dependency without outside code or a network.

## First run the function in one file

```rust
fn label(title: &str) -> String {
    format!("[{title}]")
}

fn main() {
    println!("{}", label("Meeting"));
}
```
```text
[Meeting]
```

The first example has no dependency: `label` adds brackets to the title. Next we will move that very function into a separate crate. The application’s `src/main.rs` then becomes a consumer. `path` in Cargo.toml points to a local crate; it does not publish the library on crates.io.

## Recall without looking

1. How does a version requirement in `Cargo.toml` differ from the selection in `Cargo.lock`?
2. Why inspect license and documentation before adding outside code?
3. Why can a `path` dependency work offline?

## Exercise

**Required.** First add a second call to `label("Plan")` in this example. Then create a practice library `organizer_label` beside the application directory: `cargo new --lib organizer_label --vcs none` makes a library project without a separate Git repository. Move `label` into `organizer_label/src/lib.rs` and put `pub` before `fn` to make it callable outside. In the application Cargo.toml add `organizer_label = { path = "../organizer_label" }` under `[dependencies]`. Call `organizer_label::label` from `main.rs`. Run `cargo run --offline` and check two lines. Run `cargo new` from the parent directory of the two projects; run `cargo run --offline` inside the application directory.

## Answers

If `cargo new` reports that the directory already exists, choose another name for this practice library. `lib.rs` has no `main`; write `pub fn label(title: &str) -> String`. In the main file, start the call with the dependency name.

<!-- task-answer -->
```rust
fn label(title: &str) -> String {
    format!("[{title}]")
}

fn main() {
    println!("{}", label("Meeting"));
    println!("{}", label("Plan"));
}
```
```text
[Meeting]
[Plan]
```

### Complete result after moving into the library

The `organizer` and `organizer_label` directories are siblings. `cargo new --lib` creates `organizer_label/Cargo.toml`; its package name should be `organizer_label`. Put this in `organizer_label/src/lib.rs`:

```text
pub fn label(title: &str) -> String {
    format!("[{title}]")
}
```

Keep the existing `[package]` fields in `organizer/Cargo.toml` and add this line under `[dependencies]`:

```text
organizer_label = { path = "../organizer_label" }
```

`organizer/src/main.rs`:

```text
fn main() {
    println!("{}", organizer_label::label("Meeting"));
    println!("{}", organizer_label::label("Plan"));
}
```

```text
[Meeting]
[Plan]
```

## After checking

After the local crate exercise, search crates.io for a library you actually need. Read its docs, version, license, update date, and other dependencies. For a real need, `cargo add NAME` writes the manifest entry and updates the lockfile; `cargo tree` shows the chain. Do not add a library just for this exercise. Commit the application `Cargo.lock` to Git. [Official Cargo dependencies guide](https://doc.rust-lang.org/cargo/guide/dependencies.html).

[Previous lesson](/read/rust-47-lifetimes?lang=en) · [Contents](/course/rust?lang=en)
