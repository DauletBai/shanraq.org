# Terminal input: a line, an empty answer and end of input

_Lead (summary):_ **Lesson 34. Read a task title and distinguish a blank line, end of input and a read error.**

Before this lesson: lessons 4, 21–22, 25 and 31–33. If `Result` is unclear, revisit [lesson 32](/read/rust-32-result?lang=en).

## Familiar image and recall map

At a reception desk, someone can write a title, leave the line blank, or close the channel. These are three different events. **Standard input** is a stream of text a program usually receives from the terminal; another process can also send text into it. **End of file** (EOF) means no more data will arrive. Pressing Enter on a blank line provides a line; it is not EOF. The image has a limit: `read_line` reports the number of **bytes**, not visible letters, and reading can fail.

**Empty `String` → `read_line(&mut ...)` → `Ok(0)` / `Ok(number)` / `Err` → `trim()` → response.** `std::io::stdin()` accesses standard input. `std::io` is the path to input and output tools in Rust's standard library; `::` joins path parts as in lesson 28. `String::new()` creates an empty string. `read_line(&mut line)` adds a line to it, so it needs a mutable reference. It returns `Result<usize, std::io::Error>`: `Ok(0)` means EOF, a positive number in `Ok` gives bytes read, and `Err` means reading failed. `usize` is the size and index type from the array lessons. `trim()` gives a `&str` without whitespace or a line ending **at the edges**; spaces inside the title remain. If the result equals `""`, the title is empty.

## Read one answer

Save the previous `src/main.rs` in `organizer`, replace the whole file, and run `cargo run` beside `Cargo.toml`. No other files or dependencies change. After the prompt, type `Reading` and press Enter: the program prints `Task: Reading`. Your typed characters might also appear in the terminal; they are not output from `println!`.

For automated checking, the example receives an already ended empty input stream. The output block below therefore shows the EOF branch. In normal interactive use, enter a title to try the other branch.

```rust
fn main() {
    let mut line = String::new();
    println!("Enter a title:");
    match std::io::stdin().read_line(&mut line) {
        Ok(0) => println!("Input ended"),
        Ok(_) => {
            let title = line.trim();
            if title == "" {
                println!("Empty title");
            } else {
                println!("Task: {title}");
            }
        }
        Err(_) => println!("Could not read input"),
    }
}
```
```text
Enter a title:
Input ended
```

Pressing Enter without text enters `Ok(_)`; `trim()` leaves an empty title, so the program says `Empty title`. An ended stream enters `Ok(0)`. `trim()` also removes spaces at both ends of a title, which may fit task names but must be a deliberate choice for other data. This reads one line; it is not yet the organizer's complete command loop.

## EOF is not an empty title

Try three observations separately: the title `Reading`, an empty line followed by Enter, and an input stream with no data. Confusing `Ok(0)` with `Ok(_)` can leave a program asking for a title after its input source has closed. Do not silently turn `Err` into an empty task, and do not use `.unwrap()` on human input.

## Check your understanding

1. Why is `line` declared with `mut`?
2. How does `Ok(0)` differ from Enter on a blank line?
3. Does `trim()` remove a space between words?

Check: `read_line` changes the string; `Ok(0)` is EOF, while Enter provides a line ending; an internal space remains. Rebuild the recall map without the page.

## Exercise

**Required.** Read one title from standard input. For EOF, print `No more input`; for a title empty after `trim()`, print `A title is needed`; for a nonempty title, print `Add: TITLE` using the title from input. For `Err`, print `Read error`. Try all three ordinary cases; to test EOF, end the input stream using your shell's method or provide an empty file.

**Boundary check.** Enter `My plan` and confirm the internal space remains.

## Hint

Inside `match`, put `Ok(0)` first, then `Ok(_)`, then `Err(_)`. In `Ok(_)`, use `let title = line.trim();` and compare with empty text.

## Reference answer after trying

<!-- task-answer -->
```rust
fn main() {
    let mut line = String::new();
    match std::io::stdin().read_line(&mut line) {
        Ok(0) => println!("No more input"),
        Ok(_) => {
            let title = line.trim();
            if title == "" {
                println!("A title is needed");
            } else {
                println!("Add: {title}");
            }
        }
        Err(_) => println!("Read error"),
    }
}
```
```text
No more input
```

The reference output is for an empty stream; entering `My plan` produces `Add: My plan`. [Official `read_line` documentation](https://doc.rust-lang.org/std/io/struct.Stdin.html#method.read_line).

[Previous lesson](/read/rust-33-propagation?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-35-arguments?lang=en)
