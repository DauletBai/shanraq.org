# Building and sharing the program

_Summary:_ **Lesson 58. Build for your OS, show a version, and hand over a clear archive.**

Revisit [settings and README](/read/rust-55-documentation?lang=en) and [automatic checks](/read/rust-57-quality?lang=en). Distribute only a checked practice project from a separate folder.

## Familiar image and recall map

A recipe and baked bread are different things. Source code resembles the recipe; an **executable file** resembles bread baked for a particular oven. **Recall map:** check source → build a release → run the file → include README and licence notice → hand over an archive → test after extraction. The analogy has a limit: a file built for one OS usually cannot run directly on another.

`cargo build --release` builds with release settings. It may take longer than a practice `cargo build`, but the result is prepared for use. Cargo puts it under `target/release`: `organizer.exe` on Windows, `organizer` on macOS/Linux. `version = "0.1.0"` in `Cargo.toml` is the project version. `env!("CARGO_PKG_VERSION")` is a macro: Cargo inserts that version at compile time. It is embedded in the program and does not change when a runtime environment variable changes. `#[test]` checks that the version has three parts; three parts alone promise no compatibility. An **archive** bundles a program, instructions, and required notices into one file.

## Run and inspect

The manifest field `license = "MIT"` declares the terms for this course example code; choose terms deliberately for your own new code. In a separate Cargo project set `version = "0.1.0"` and put the first Rust block in `src/main.rs`. Run `cargo test --locked`, then `cargo build --release`. Run the built file from `target/release`; the output below is for this manifest. Rebuild after changing `version`, or the old executable will still show the old number. Create `README.md` in the practice folder: state the required OS, launch command, and that organizer data is separate from the executable. The course examples use the MIT section of [LICENSE-COURSE](https://github.com/DauletBai/shanraq.org/blob/main/LICENSE-COURSE); retain its copyright and permission notice when sharing a substantial part of that code.

A manually tested ZIP archive made by one tool lost executable permission after extraction on macOS. On macOS/Linux, `tar.gz` is convenient because it preserves that permission. If the extracted file cannot run because of permissions, `chmod +x organizer` in its folder adds execute permission (`chmod` changes permissions and `+x` adds execution). Windows handles this differently; test `organizer.exe` there.

For a repeatable macOS/Linux step, run `tar -czf organizer.tar.gz organizer README.md LICENSE` in the folder holding those three files, then create an empty `extracted` subfolder, enter it, and run `tar -xzf ../organizer.tar.gz`. `tar` packs or extracts files; `-c` creates, `-x` extracts, `-z` uses gzip compression, and `-f` gives the archive name. On Windows, make a ZIP in File Explorer and test the extracted `.exe`.

```toml
[package]
name = "organizer"
version = "0.1.0"
edition = "2024"
license = "MIT"

[dependencies]
```

```rust
fn main() {
    println!("Version: {}", env!("CARGO_PKG_VERSION"));
}

#[test]
fn package_version_has_three_parts() {
    let parts: Vec<&str> = env!("CARGO_PKG_VERSION").split('.').collect();
    assert_eq!(parts.len(), 3);
}
```
```text
Version: 0.1.0
```

## Recall without looking

1. Why does editing `Cargo.toml` not change an executable already built?
2. Why test an archive after extracting it?
3. What else must you give the recipient besides the executable?

## Exercise

**Required.** Add `--version` handling and make an unknown argument exit with status 2. Build a release, put the executable, `README.md`, and licence notice in a separate folder, and archive that folder with your OS tools. Extract it into another folder on the same OS and run `--version`. Predict the output before running.

## Answers

Read words after the program name with `std::env::args().skip(1)`. For errors, use `eprintln!` and `std::process::exit(2)`. Take the course-code permission text from `LICENSE-COURSE`.

<!-- task-answer -->
```rust
fn version_text() -> String {
    format!("Version: {}", env!("CARGO_PKG_VERSION"))
}

fn run(args: &[String]) -> Result<String, String> {
    match args {
        [] => Ok(version_text()),
        [flag] if flag == "--version" => Ok(version_text()),
        _ => Err(String::from("Unknown argument")),
    }
}

fn main() {
    let args: Vec<String> = std::env::args().skip(1).collect();
    match run(&args) {
        Ok(message) => println!("{message}"),
        Err(message) => {
            eprintln!("{message}");
            std::process::exit(2);
        }
    }
}

#[test]
fn prints_version_for_flag() {
    assert!(run(&[String::from("--version")]).unwrap().contains("0.1.0"));
}

#[test]
fn rejects_unknown_argument() {
    assert!(run(&[String::from("--unknown")]).is_err());
}
```
```text
Version: 0.1.0
```

`README.md`:

```markdown
# Organizer 0.1.0

Built for the OS on which this archive was made. Extract all files, then run `organizer --version` (Windows: `organizer.exe --version`). This practice program only shows its version and does not yet open a task file. Keep the future organizer task file separate from the program and back it up. On macOS/Linux, if the extracted file lacks execute permission, run `chmod +x organizer` in this folder. On failure, read stderr and never replace data with an empty file.
```

## After checking

Do not promise that a Windows executable works on macOS or Linux. Build a separate file for another OS and check it there; lesson 57 CI checks code but does not automatically hand over a release archive. Keep secrets and personal data out of the practice archive. [cargo build documentation](https://doc.rust-lang.org/cargo/commands/cargo-build.html).

[Previous lesson](/read/rust-57-quality?lang=en) · [Contents](/course/rust?lang=en)
