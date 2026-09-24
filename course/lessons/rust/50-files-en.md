# Reading and writing files

_Summary (summary):_ **Lesson 50. Write text, read it back, and explain an input/output failure.**

First review [Result](/read/rust-32-result?lang=en), [error propagation](/read/rust-33-propagation?lang=en), and [paths](/read/rust-49-paths?lang=en). This is a separate practice file in a temporary directory, not real organizer data.

## Familiar image and recall map

Writing in a notebook can fail: it may be missing, inaccessible, or damaged. **Input/output** (I/O) is a program’s interaction with the outside world, including files. `Result` distinguishes success from a specific error. The image has a limit: a file can change between a check and a read, so finding a path first does not guarantee that a later operation succeeds.

**choose a separate temporary path → create without overwriting → write bytes → close → read UTF-8 → remove → handle failure.** `temp_dir()` gives the system temporary directory; `process::id()` adds this process’s number to the practice name. `create_new(true)` requires a new file and will not overwrite an existing one. `write_all` is a method of the `Write` trait, so we import it; it attempts to write all text bytes. `as_bytes()` exposes the UTF-8 string as bytes. `drop(file)` closes the file before reading, including on Windows. `read_to_string` reads the whole file as UTF-8: a missing file, denied access, or invalid UTF-8 produces `Err`. `?` passes an error out of `run`, while `std::io::Result<()>` means success without an extra value. `eprintln!` writes errors separately from normal output. `ErrorKind` names the kind of error; `NotFound` means the file could not be found.

`std::process::exit(1)` ends the program with status 1 for an unresolved error; status 0 would misleadingly report success to a caller.

## Save and read back a practice line

On a phone, swipe a long code line sideways inside its light box; the page itself stays within the screen.

```rust
use std::fs::{self, OpenOptions};
use std::io::Write;

fn run() -> std::io::Result<()> {
    let path = std::env::temp_dir().join(format!("rust-lesson-50-{}.txt", std::process::id()));
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(&path)?;
    file.write_all("Buy a book\n".as_bytes())?;
    drop(file);
    let text = fs::read_to_string(&path)?;
    print!("{text}");
    fs::remove_file(&path)?;
    Ok(())
}

fn main() {
    if let Err(error) = run() {
        eprintln!("File error: {error}");
        std::process::exit(1);
    }
}
```
```text
Buy a book
```

After reading, `remove_file` deletes only this practice file. If an earlier action fails, the file may remain in the temporary directory; a real application needs to account for this. If the name happens to repeat, `create_new` fails and preserves the earlier file. This is not persistent storage for the organizer: safe saving and recovery come later.

## Recall without looking

1. How does `create_new(true)` protect existing data?
2. Why can `read_to_string` return `Err` for a file that exists?
3. Where does `eprintln!` send an error?

## Exercise

**Required.** After removing the file, try `fs::read_to_string(&path)` again. Use `match` and `ErrorKind::NotFound` to print `File missing`; include separate branches for an unexpected success and another error. Predict the result first.

## Answers

Add `use std::io::{ErrorKind, Write};`. After `remove_file`, compare `error.kind()` with `ErrorKind::NotFound`; do not compare the operating system’s message text, which varies by language and system.

<!-- task-answer -->
```rust
use std::fs::{self, OpenOptions};
use std::io::{ErrorKind, Write};

fn run() -> std::io::Result<()> {
    let path = std::env::temp_dir().join(format!("rust-lesson-50-{}.txt", std::process::id()));
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(&path)?;
    file.write_all("Buy a book\n".as_bytes())?;
    drop(file);
    let text = fs::read_to_string(&path)?;
    print!("{text}");
    fs::remove_file(&path)?;
    match fs::read_to_string(&path) {
        Err(error) if error.kind() == ErrorKind::NotFound => println!("File missing"),
        Ok(_) => println!("Unexpected success"),
        Err(error) => eprintln!("File error: {error}"),
    }
    Ok(())
}

fn main() {
    if let Err(error) = run() {
        eprintln!("File error: {error}");
        std::process::exit(1);
    }
}
```
```text
Buy a book
File missing
```

## After checking

An inaccessible file may yield `PermissionDenied`; invalid UTF-8 also makes `read_to_string` fail. Do not present either case as an empty task list. Real data also needs protection against interrupted writes and a backup; later lessons cover these. [Official read_to_string documentation](https://doc.rust-lang.org/std/fs/fn.read_to_string.html).

[Previous lesson](/read/rust-49-paths?lang=en) · [Contents](/course/rust?lang=en)
