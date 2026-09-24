# Settings and documentation

_Summary:_ **Lesson 55. Explain data-path selection and write instructions another person can follow.**

Revisit [paths](/read/rust-49-paths?lang=en), [arguments](/read/rust-35-arguments?lang=en), and [commands](/read/rust-54-cli?lang=en).

## Familiar image and recall map

You can put a mailing address on the letter, keep it in an address book, or use a default address. The one stated on the letter should win. **Recall map:** command path → otherwise environment path → otherwise local default path. The **environment** is a set of settings supplied to a program by its parent process at launch. The analogy does not make a relative path absolute: `tasks.json` depends on the current folder.

`Option<PathBuf>` means a path may be present or absent. `if let Some(path)` extracts it only when present. The first path found is returned, so later branches do not run. `PathBuf::from` constructs a path value and `display()` shows it to a person. `///` is a documentation comment before a function; Cargo puts it in pages made by `cargo doc`. The first example passes values directly for repeatable output. A real launch obtains command arguments with `args` and environment values with `std::env::var_os`; when a variable is absent, `var_os` returns `None`.

## Run and inspect

Create a Cargo project without dependencies and paste the first Rust block into `src/main.rs`. `cargo run` shows a default selection and an environment selection. Then run `cargo doc --no-deps`: it builds HTML documentation for this project in `target/doc` without dependency pages. Open it in a browser from your file manager. A README is a normal `README.md` file at the project root; a newcomer should find requirements, launch commands, data location, and error guidance there.

```rust
use std::path::PathBuf;

/// Choose the data path: command line, environment, then local default.
pub fn choose_path(command: Option<PathBuf>, environment: Option<PathBuf>) -> PathBuf {
    if let Some(path) = command {
        return path;
    }
    if let Some(path) = environment {
        return path;
    }
    PathBuf::from("tasks.json")
}

fn main() {
    println!("Default path: {}", choose_path(None, None).display());
    println!(
        "Environment path: {}",
        choose_path(None, Some(PathBuf::from("family.json"))).display()
    );
}
```
```text
Default path: tasks.json
Environment path: family.json
```

## Recall without looking

1. Which path wins when all three are available?
2. What does a relative `tasks.json` depend on?
3. Who benefits from README and `///` comments?

## Exercise

**Required.** Add a third `choose_path` call with both a command path and an environment path; print the chosen path. Then write `README.md` with installation, `cargo run`, path precedence, and the current-folder warning. Ask someone unfamiliar with the lesson to follow it, or follow it yourself from a fresh empty folder.

## Answers

Call `choose_path(Some(PathBuf::from("work.json")), Some(PathBuf::from("family.json")))`. Do not assume Cargo is installed in the README; use the instructions from lessons 5–6.

<!-- task-answer -->
```rust
use std::path::PathBuf;

/// Choose the data path: command line, environment, then local default.
pub fn choose_path(command: Option<PathBuf>, environment: Option<PathBuf>) -> PathBuf {
    if let Some(path) = command {
        return path;
    }
    if let Some(path) = environment {
        return path;
    }
    PathBuf::from("tasks.json")
}

fn main() {
    println!("Default path: {}", choose_path(None, None).display());
    println!(
        "Environment path: {}",
        choose_path(None, Some(PathBuf::from("family.json"))).display()
    );
    println!(
        "Command path: {}",
        choose_path(
            Some(PathBuf::from("work.json")),
            Some(PathBuf::from("family.json"))
        )
        .display()
    );
}
```
```text
Default path: tasks.json
Environment path: family.json
Command path: work.json
```

An example `README.md` for this practice project:

```markdown
# Organizer: choosing a file

1. Open https://www.rust-lang.org/tools/install, install Rust for your OS, open a new terminal, and check `cargo --version`.
2. Run `cargo run` in the project folder.
3. The command path wins over the environment path; otherwise use `tasks.json` in the current folder. This practice example passes values directly in code and does not create a file yet.
4. If Cargo reports an error, read it and check that you ran the command in the folder containing `Cargo.toml`.
```

Run `cargo doc --no-deps` and find the `choose_path` description in the generated documentation.

## After checking

An environment variable is not necessarily a safe place for secrets; this example contains only a path. For the full organizer, document where data lives, how backups work, what to do with damaged files, and remaining limits. [cargo doc documentation](https://doc.rust-lang.org/cargo/commands/cargo-doc.html).

[Previous lesson](/read/rust-54-cli?lang=en) · [Contents](/course/rust?lang=en)
