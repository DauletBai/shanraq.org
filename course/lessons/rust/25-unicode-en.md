# Text, UTF-8, and visible letters

_Lead (summary):_ **Lesson 25. Distinguish bytes, `char` values, and visible letters in a task title.**

Prerequisites: lessons 3, 9, and 22–24. If `String` and `&str` blur together, return to [lesson 22](/read/rust-22-strings?lang=en).

## A familiar image and recall map

You can count an address on paper by lines, words, or letters; these answer different questions. Computers also store and count text in different units. **UTF-8** encodes text as bytes; one visible mark can occupy several bytes. A Rust **`char`** is one Unicode scalar value, a permitted code point, and does not always match a letter as a person sees it. A **grapheme** is a visible writing unit that can contain several scalar values. The image has a limit: your eyes may recognize a letter, but cannot infer its byte count without UTF-8 rules.

**Text → bytes with `.len()` → scalar values with `.chars()` → visible graphemes.** Say what question each number answers. `.chars()` yields `char` values one after another; this sequence provider is an **iterator**. `.count()` walks that sequence and counts its items. It does not count graphemes. We will study iterators in detail later; for now these two calls count `char` values without making another `String`.

## One letter, two bytes

In `organizer`, save your previous `src/main.rs` separately, replace it completely, and run `cargo run` next to `Cargo.toml`. No other files or dependencies change. Predict two numbers; output below excludes Cargo messages.

For a `String`, `.len()` returns the number of **bytes**, while `.chars().count()` returns the number of Unicode scalar values. `é` occupies two bytes in UTF-8 but is one `char` value and one visible letter in this form.

```rust
fn main() {
    let title = String::from("é");
    println!("Bytes: {}", title.len());
    println!("char values: {}", title.chars().count());
}
```

```text
Bytes: 2
char values: 1
```

`len()` is not wrong: it answers a storage question, not how long the word looks to a reader. For Latin `A`, both counts are 1, but one such example cannot establish a rule for all text.

## When one visible letter has two values

The next spelling, `"e\u{301}"`, is a separate example. Inside a string literal, `\u{301}` means Unicode code point 301 in hexadecimal notation: an accent mark added to the preceding `e`. The backslash begins special notation inside quotes, `u` names Unicode, and `{301}` gives the code. This is **not** a variable insertion like `{title}` in `println!`. Together `e` and the mark look like one letter, but form two `char` values and three UTF-8 bytes.

```rust
fn main() {
    let title = "e\u{301}";
    println!("Bytes: {}", title.len());
    println!("char values: {}", title.chars().count());
}
```

```text
Bytes: 3
char values: 2
```

Thus `.chars().count()` does not always answer “how many visible letters?” either. Rust's standard library does not provide a ready grapheme counter; that needs separate text processing. We will not impose a “number of letters” limit on organizer titles until we choose an exact rule.

## Why `title[0]` is unavailable

Index zero does not state whether we want the first byte, `char` value, or grapheme. A non-ASCII letter occupies multiple bytes; its first byte alone is not a complete character. This **separate** example therefore fails with E0277:

<!-- error-code: E0277 -->
```rust,compile_fail
fn main() {
    let title = String::from("é");
    println!("{}", title[0]);
}
```

Use `.chars()` to read scalar values in order. Text slices from lesson 23 have **byte** boundaries: slicing in the middle of a UTF-8 character panics. Do not guess such boundaries. `String` still suits storing a title, and `&str` suits reading it.

## Check your understanding

1. What does `"é".len()` count?
2. Does `.chars().count()` count graphemes?
3. Why is `title[0]` not a safe way to get the first letter?

Check: bytes; no, Unicode `char` values; the intended text unit is unspecified. Rebuild the recall map without looking.

## Exercise

**Required.** Create a `String` with `é` and a `"e\u{301}"` literal. For each, print byte and `char` counts in this order: `é: bytes 2, char values 1` and `accented e: bytes 3, char values 2`. Compute both numbers with methods; do not hard-code the answers in `println!`.

**Your observation.** Replace `é` with `A`: both numbers become 1. Restore the first example.

## Hint

For each text, pass `.len()` and `.chars().count()` as arguments to `println!`. The literal has type `&str`, but both methods work on it too.

## Reference answer after your attempt

<!-- task-answer -->
```rust
fn main() {
    let title = String::from("é");
    let combined = "e\u{301}";
    println!("é: bytes {}, char values {}", title.len(), title.chars().count());
    println!(
        "accented e: bytes {}, char values {}",
        combined.len(),
        combined.chars().count()
    );
}
```

```text
é: bytes 2, char values 1
accented e: bytes 3, char values 2
```

The result separates storage size from scalar count and calls neither a count of visible letters. [Official UTF-8 and string guide](https://doc.rust-lang.org/book/ch08-02-strings.html).

[Previous lesson](/read/rust-24-vectors?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-26-tuples?lang=en)
