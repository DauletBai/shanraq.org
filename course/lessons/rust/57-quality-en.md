# Automatic checks on three systems

_Summary:_ **Lesson 57. Run formatting, Clippy, and tests whenever code is sent.**

Revisit [isolated tests](/read/rust-56-testing?lang=en) and [your first Cargo project](/read/rust-06-cargo?lang=en). A repository stores project files and their change history on GitHub. Use a separate practice repository: its workflow file starts checks when you send code.

## Familiar image and recall map

Before printing a notebook, three proofreaders inspect their own copies and report problems to the author. **Recall map:** change code → check locally → send it → build and test on three systems → read results → fix. The analogy has a limit: only recorded rules and scenarios are checked, so success does not prove the absence of every bug.

**CI** (continuous integration) automatically checks a submitted change. **GitHub Actions** reads files under `.github/workflows` at the repository root. In YAML, `name:` gives a name, indentation groups fields, `-` starts a list item, `on` selects events, `jobs` defines work, `steps` lists actions, `uses` invokes an existing action, and `run` runs a command. `matrix.os` names three operating-system variants; `${{ matrix.os }}` inserts one into `runs-on`. `actions/checkout@v7` fetches source. `rustup toolchain install` installs a specific compiler version and components; `rust-toolchain.toml` stores that choice with the project. `fmt --check` checks formatting without changing files, Clippy looks for suspicious code, and `test --locked` runs tests with pinned dependencies. `-- -D warnings` asks Clippy to treat warnings as errors.

In `cargo +1.97.0`, `+` selects the installed Rust 1.97.0 toolchain; `--all-targets` includes test targets, and `--locked` forbids changes to `Cargo.lock`. YAML square brackets list events or systems, and `@v7` selects an action version.

## Run and inspect

Paste the Rust code into a new `src/main.rs`; `Cargo.toml` needs no dependencies. Run `cargo fmt --check`, `cargo clippy --all-targets -- -D warnings`, `cargo test --locked`, then `cargo run`. The first block has one test. Create `.github/workflows/organizer.yml` and `rust-toolchain.toml` at the project root using the examples below. Commands in the YAML must run at this project root. The first installation of 1.97.0 needs network access. GitHub shows workflow results after a push; a green result does not replace human review of the exercise.

```rust
fn normalized_title(text: &str) -> Option<&str> {
    let trimmed = text.trim();
    if trimmed.is_empty() {
        None
    } else {
        Some(trimmed)
    }
}

fn main() {
    let title = normalized_title("  Buy a book  ").unwrap_or("Empty");
    println!("Title: {title}");
}

#[test]
fn trims_outer_spaces() {
    assert_eq!(normalized_title("  Book  "), Some("Book"));
}
```
```text
Title: Buy a book
```

Files for the practice project root:

`rust-toolchain.toml`:

```toml
[toolchain]
channel = "1.97.0"
components = ["rustfmt", "clippy"]
profile = "minimal"
```

`.github/workflows/organizer.yml`:

```yaml
name: Organizer checks
on: [push, pull_request]
jobs:
  verify:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v7
      - run: rustup toolchain install 1.97.0 --profile minimal --component rustfmt,clippy
      - run: cargo +1.97.0 fmt --check
      - run: cargo +1.97.0 clippy --all-targets -- -D warnings
      - run: cargo +1.97.0 test --locked
```

## Recall without looking

1. Why run `cargo test` locally before sending code?
2. What does `matrix.os` mean?
3. Why does successful CI not prove that the intended behaviour is correct?

## Exercise

**Required.** Add tests for a blank title and a title containing a character outside ASCII. Confirm that all three tests pass locally. Create both configuration files from the example, push to your practice repository, and find the result on all three systems. If GitHub is unavailable, keep the files and explain which map steps are still pending.

## Answers

`normalized_title("   ")` should return `None`. Use `é` for the Unicode case. One operating system is not a complete matrix.

<!-- task-answer -->
```rust
fn normalized_title(text: &str) -> Option<&str> {
    let trimmed = text.trim();
    if trimmed.is_empty() {
        None
    } else {
        Some(trimmed)
    }
}

fn main() {
    let title = normalized_title("  Buy a book  ").unwrap_or("Empty");
    println!("Title: {title}");
}

#[test]
fn trims_outer_spaces() {
    assert_eq!(normalized_title("  Book  "), Some("Book"));
}

#[test]
fn rejects_blank_title() {
    assert_eq!(normalized_title("   "), None);
}

#[test]
fn keeps_unicode_title() {
    assert_eq!(normalized_title("  Café  "), Some("Café"));
}
```
```text
Title: Buy a book
```

## After checking

If one system fails, open that system’s failing step and reproduce the command locally. `fmt`, Clippy, and tests check different things; do not hide errors to obtain a green mark. Revisit the toolchain and action versions when the course is updated. [Cargo fmt documentation](https://doc.rust-lang.org/cargo/commands/cargo-fmt.html).

[Previous lesson](/read/rust-56-testing?lang=en) · [Contents](/course/rust?lang=en)
