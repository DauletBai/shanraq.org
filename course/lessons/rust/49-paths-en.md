# Paths and the file system

_Summary (summary):_ **Lesson 49. Choose an understandable data file location across systems.**

First review the [working directory](/read/rust-03-workspace?lang=en), [terminal](/read/rust-04-terminal?lang=en), and [errors](/read/rust-32-result?lang=en). We only build a path here; no file is created.

## Familiar image and recall map

An address on an envelope says where to look for a house; it does not promise that the house exists or its door is open. A **path** is an address for a file or folder. The **current working directory** is the location from which a program resolves a relative path. The image has a limit: the same relative path can lead to different places when run from different directories.

**Path — borrowed path view → join — add a name → PathBuf — owned path → check before acting.** `Path::new("data")` makes a borrowed path view; `PathBuf` owns the joined path. `join("tasks.txt")` adds a component according to the operating system’s rules. `std::path::{Path, PathBuf}` imports two names from one module; the braces group the names. `is_absolute()` tells whether the address starts at a root; `data/tasks.txt` is relative. `file_name()` returns `Option` because a root or unusual path can lack a final name. This teaching example certainly has one, so `unwrap()` is appropriate; handle `None` for a user supplied path. `to_string_lossy()` displays an OS file name as text, potentially replacing characters that are not valid UTF-8; keep `Path` or `PathBuf` for actual file operations.

## Build an address without touching the disk

```rust
use std::path::{Path, PathBuf};

fn main() {
    let folder = Path::new("data");
    let file: PathBuf = folder.join("tasks.txt");
    println!("{}", file.is_absolute());
    println!("{}", file.file_name().unwrap().to_string_lossy());
}
```
```text
false
tasks.txt
```

`false` means relative; the final line is the file name. We avoid printing the whole path because directory separators differ between Windows and Unix. Paths can contain spaces and Kazakh letters; `join` avoids writing a separator yourself.

The file operation checks access permissions: an existing path can still yield `PermissionDenied`. A correct address does not open a locked door.

## Recall without looking

1. From where is `data/tasks.txt` resolved?
2. How do `PathBuf` and `Path` differ?
3. Does a well-formed path prove that a file exists?

## Exercise

**Required.** Add a second path `folder.join("backup.txt")` and print only its name, `backup.txt`. Then run the project from another working directory with `cargo run --manifest-path PATH_TO_PROJECT/Cargo.toml`: `--manifest-path` locates the Cargo project; `cargo run` keeps the terminal’s current working directory for the program. Explain where the relative path would now resolve. Create no files and use no real personal data path.

## Answers

Repeat the `file_name()` line for `backup`. The `unwrap()` here is only for our known practice name.

<!-- task-answer -->
```rust
use std::path::{Path, PathBuf};

fn main() {
    let folder = Path::new("data");
    let file: PathBuf = folder.join("tasks.txt");
    println!("{}", file.is_absolute());
    println!("{}", file.file_name().unwrap().to_string_lossy());
    let backup = folder.join("backup.txt");
    println!("{}", backup.file_name().unwrap().to_string_lossy());
}
```
```text
false
tasks.txt
backup.txt
```

## After checking

Choose a real data location explicitly and show it to the user. The next lesson treats missing and inaccessible files as input/output results. If working directories remain unclear, revisit [lesson 4](/read/rust-04-terminal?lang=en). [Official PathBuf documentation](https://doc.rust-lang.org/std/path/struct.PathBuf.html).

[Previous lesson](/read/rust-48-dependencies?lang=en) · [Contents](/course/rust?lang=en)
