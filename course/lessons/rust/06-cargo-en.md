# Your first Cargo project

_Summary:_ **Distinguish checking, building and running.**

Prerequisite: lessons 1–5. Work in organizer unless stated otherwise.

## Why this matters

In lesson 5 you built a disposable installation probe. Now we create the organizer itself. Cargo resembles a workshop coordinator: it knows where the source lives and how to call the tools. A **project** is the collection of files for your work. A **Cargo package** has a Cargo.toml description and one or more buildable targets. Our package currently contains one executable program.

## The whole picture

Enter rust-learning as in lesson 4, not the rust-install-check folder. cargo new creates a new project directory; organizer is its name. --edition selects the language edition; --vcs none avoids creating Git history before we have explained that tool. Run:

```sh
cargo new organizer --edition 2024 --vcs none
cd organizer
cargo check
cargo build
cargo run
```

check checks the code without producing an executable. build creates the executable without running it. run rebuilds when necessary and runs it. The new program prints Hello, world!; Compiling and Finished are Cargo messages, not our program's output.

## Where the program lives

Use your editor's Open Folder command to open organizer. Do not create another organizer inside it. The generated src/main.rs contains:

```rust
fn main() {
    println!("Hello, world!");
}
```

Output:

```text
Hello, world!
```

This is Cargo's supplied example. Lesson 7 explains every symbol; you are not yet expected to write the entry-point syntax independently.

**Cargo.toml** is the package's description in **TOML**, a text configuration format using sections and fields. name identifies the package, version identifies our application version, and edition chooses the language rules. dependencies lists libraries used by the project. Those are reusable parts; there are no third-party dependencies yet. Cargo.lock records selected dependency versions and is maintained by Cargo, not edited as a second settings file.

Rust 1.97.0 and Edition 2024 are different: a tool release and a language edition. target holds generated build output; do not edit those files. Debug is the development build; we will study optimized release builds before distributing the program.

After cargo build you can run the executable directly. Windows PowerShell:

```powershell
.\target\debug\organizer.exe
```

macOS/Linux:

```sh
./target/debug/organizer
```

The leading ./ or .\ starts from the current directory. Direct execution does not rebuild edited source. After a change, use build followed by direct execution, or just run. Perform these commands inside organizer, where Cargo.toml lives.

## If something goes wrong

could not find Cargo.toml usually means you are in the wrong directory. destination already exists means the folder exists: inspect it rather than deleting your work. A linker error returns you to lesson 5's system-tool installation.

## Recall map

new: create; check: check; build: build; run: build if needed and execute. Source in src, configuration in Cargo.toml, generated output in target.

## Warm-up

1. Predict: does saving source update a previously built executable?
2. Complete: ___ checks without producing an executable.
3. Repair a Cargo invocation from the parent directory by entering organizer.

## Exercise

**Required.** Replace only Hello, world! between the quotes with “The organizer is ready”. Save main.rs, run cargo run, then build and run the executable directly for your system. Expect that exact line both times.

**Your own data.** Change the greeting while preserving the quotation marks.

## Hint and reference answer

If the old message remains, check saving and the current directory. The reference below is a complete main.rs, not a command to paste into the shell. Attempt the exercise before reading it.

Warm-up answers: saving does not rebuild; the command is check; Cargo needs the directory containing Cargo.toml. Continue when both execution methods work and you can explain the difference.

<!-- task-answer -->
```rust
fn main() {
    println!("The organizer is ready");
}
```

Expected output:

```text
The organizer is ready
```

[Cargo reference](https://doc.rust-lang.org/cargo/).

[Previous lesson](/read/rust-05-installation?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-07-first-program?lang=en)
