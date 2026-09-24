# Rust lessons 46–50: draft review

Status: draft in kz/ru/en. The publication manifest remains at 45 lessons. This block prepares the organizer for file storage; it does not yet claim that the organizer persists its tasks. JSON data, safe replacement and recovery are scheduled for lessons 51–53.

| Lesson | New step | Boundary explained | Independent result |
|---|---|---|---|
| 46 | `trait`, `impl`, `T: Summary`, `#[derive(Debug)]` | A trait states available behavior, while the implementation defines it; debug output is for developers | Implement the same summary behavior for a note |
| 47 | Named lifetime `'a` | An annotation relates references; it cannot keep local data alive | Choose between two strings safely; reject a dangling return |
| 48 | Local path dependency, `Cargo.toml`, `Cargo.lock`, crates.io | A version requirement differs from a locked version; adding outside code calls for documentation and license review | Move `label` into a separate library crate and use it offline |
| 49 | `Path`, `PathBuf`, relative paths | A valid path does not guarantee a file exists or is accessible | Build a second path without assuming an OS separator |
| 50 | `OpenOptions`, `Write`, `read_to_string`, I/O errors | Existing data must not be overwritten during practice; a missing or unreadable file is an error, not an empty task list | Read a practice file, remove it, and identify `NotFound` |

Each locale has a familiar-life image, its limit, a text recall map, a runnable first example with expected output, recall questions, a required exercise, a hint, a folded answer and a return route. Each first example has a matching Cargo snapshot. The answer source and output files match the inline answers.

Local verification passed 670 Rust examples, answer files and Cargo snapshots across the course, including the checker’s negative fixtures. All 15 new snapshots passed `cargo fmt --check`; language checks found no prohibited prose fragments in the new pages. Offline link checking found 78 internal and five external links in the 15 pages, with no broken internal links.

The deliberately invalid example in lesson 47 is marked `rust,compile_fail` and expects compiler error E0515. Lesson 48 additionally gives all files needed for the two-crate solution in the folded answer; the local crate exercise was run in fresh Cargo projects in all three languages with `--offline`. Lesson 50 uses a process-specific name in the system temporary directory and `create_new(true)` so its demonstration does not overwrite an existing file. On an unhandled I/O error, it writes to stderr and exits with status 1.

Technical rules were checked against the official [trait chapter](https://doc.rust-lang.org/book/ch10-02-traits.html), [lifetime chapter](https://doc.rust-lang.org/book/ch10-03-lifetime-syntax.html), [Cargo dependency reference](https://doc.rust-lang.org/cargo/reference/specifying-dependencies.html), [PathBuf documentation](https://doc.rust-lang.org/std/path/struct.PathBuf.html) and [file reading documentation](https://doc.rust-lang.org/std/fs/fn.read_to_string.html).

The site Markdown renderer and current article template were used to preview all 15 pages at 1366×900, 768×1024, 390×844 and 844×390. [All 60 responsive checks](rust-tenth-browser-checks.json) found no page overflow and confirmed code blocks and folded answers. Representative Kazakh, Russian and English screenshots were visually inspected, including scrolled code on phone and desktop. The preview template retains the previously published lesson count; it will be populated from the course during a future release.

The required exercises were reconstructed from the lesson instructions in fresh Cargo projects, without copying the answer files: 15 walkthroughs across ru/kz/en passed. The two-crate task was built offline in each language. Additional manual checks confirmed that `cargo run --manifest-path` retains the caller's working directory and that an existing practice file is preserved when `create_new(true)` fails, with process status 1. The four push workflows for the draft source commit `9faeba7c32399cc88d6f1c9e3a2d33e6bedb68be` all succeeded. Before activation, `publish.py --check` found no differences in previously published lessons and confirmed the new 15 translations are not live.

The release manifest now includes lessons 46–50. Publication still requires green CI for the exact release commit, a production backup, and post-publication checks of all localized pages, navigation, feeds and responsive layout.
