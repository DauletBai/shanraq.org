# Saving without losing data

_Summary:_ **Lesson 52. Write new data through a neighbouring temporary file and retain the old file on failure.**

Revisit [paths](/read/rust-49-paths?lang=en) and [file errors](/read/rust-50-files?lang=en). We save training text here; the JSON from lesson 51 can be passed to this function.

## Familiar image and recall map

You do not erase an old notebook page before writing its replacement. Write on a separate page, check it, then swap pages. **Recall map:** old file → neighbouring temporary file → write → sync → close → replace. A write failure leaves the old file in place. The analogy has a limit: power loss and file-system behaviour need stronger safeguards.

`with_extension` makes a temporary path beside the main file; sharing a file system matters for replacement. `create_new(true)` refuses to overwrite an existing temporary file. `write_all` writes all bytes or reports an error; `sync_all` asks the OS to flush file contents and metadata to storage; `drop` closes the file. `rename` replaces the target where the OS permits, or returns an error. `if let Err(error)` selects the error case; `let _ = remove_file` deliberately ignores a cleanup error without reporting a successful save.

## Run and inspect

Create a new Cargo project without dependencies, paste the first Rust block into `src/main.rs`, and run `cargo run`. The practice folder includes the process ID; if a crash leaves it behind, inspect its contents before removing only that folder. The example replaces old text with new text. Failure to create, write, or sync the temporary file prevents `rename`, so the old target stays in place. A failed rename is returned to the caller.

```rust
use std::fs::{self, OpenOptions};
use std::io::{self, Write};
use std::path::Path;

fn safe_save(path: &Path, text: &str) -> io::Result<()> {
    let temporary = path.with_extension(format!("tmp-{}", std::process::id()));
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(&temporary)?;
    if let Err(error) = file.write_all(text.as_bytes()) {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    if let Err(error) = file.sync_all() {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    drop(file);
    if let Err(error) = fs::rename(&temporary, path) {
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    Ok(())
}

fn run() -> io::Result<()> {
    let folder = std::env::temp_dir().join(format!("rust-lesson-52-{}", std::process::id()));
    fs::create_dir(&folder)?;
    let path = folder.join("tasks.json");
    fs::write(&path, "Old plan")?;
    safe_save(&path, "New plan")?;
    println!("Saved: {}", fs::read_to_string(&path)?);
    fs::remove_dir_all(folder)?;
    Ok(())
}

fn main() {
    if let Err(error) = run() {
        eprintln!("Error: {error}");
        std::process::exit(1);
    }
}
```
```text
Saved: New plan
```

## Recall without looking

1. Why place the temporary file beside the target?
2. Why use `create_new`?
3. At which step can the old target first be replaced?

## Exercise

**Required.** Before calling `safe_save`, create a file at the exact temporary path. Show that saving fails and the main file still contains the old text. Remove the blocking file, save again, and show the new text.

## Answers

Compute the path as the function does: `path.with_extension(format!("tmp-{}", std::process::id()))`. Check `is_err()` and read the old text back.

<!-- task-answer -->
```rust
use std::fs::{self, OpenOptions};
use std::io::{self, Write};
use std::path::Path;

fn safe_save(path: &Path, text: &str) -> io::Result<()> {
    let temporary = path.with_extension(format!("tmp-{}", std::process::id()));
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(&temporary)?;
    if let Err(error) = file.write_all(text.as_bytes()) {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    if let Err(error) = file.sync_all() {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    drop(file);
    if let Err(error) = fs::rename(&temporary, path) {
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    Ok(())
}

fn run() -> io::Result<()> {
    let folder = std::env::temp_dir().join(format!("rust-lesson-52-{}", std::process::id()));
    fs::create_dir(&folder)?;
    let path = folder.join("tasks.json");
    fs::write(&path, "Old plan")?;
    let temporary = path.with_extension(format!("tmp-{}", std::process::id()));
    fs::write(&temporary, "occupied")?;
    let refused = safe_save(&path, "New plan").is_err();
    let kept = fs::read_to_string(&path)? == "Old plan";
    println!("Old data kept: {}", refused && kept);
    fs::remove_file(&temporary)?;
    safe_save(&path, "New plan")?;
    println!("Saved: {}", fs::read_to_string(&path)?);
    fs::remove_dir_all(folder)?;
    Ok(())
}

fn main() {
    if let Err(error) = run() {
        eprintln!("Error: {error}");
        std::process::exit(1);
    }
}
```
```text
Old data kept: true
Saved: New plan
```

## After checking

This handles an ordinary write failure; it does not promise survival of every failure. Syncing the temporary file alone does not guarantee survival of a power cut after replacement; strong guarantees depend on OS and file-system rules and syncing the containing directory. Concurrent writers also need a separate policy. [rename documentation](https://doc.rust-lang.org/std/fs/fn.rename.html).

[Previous lesson](/read/rust-51-json?lang=en) · [Contents](/course/rust?lang=en)
