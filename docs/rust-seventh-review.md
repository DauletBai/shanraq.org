# Rust lessons 31–35: draft review

Working drafts in kz/ru/en. They are linked for offline checking but absent from the public release manifest. Published lesson 30 has no forward link to lesson 31, and the live lesson 31 returns 404.

| Lesson | First new idea | Familiar image and its limit | Independent task |
|---|---|---|---|
| 31 | `Option<T>`, `Some`, `None`, `Vec::get` | An archivist either finds a card or reports absence; absence has no reason | Read two positions and handle the missing one |
| 32 | `Result<T, E>`, parsing and the danger of `unwrap()` | A card may contain invalid minutes; a useful error message is authored, not generated | Parse a valid and invalid duration |
| 33 | `?`, compatible errors and early return | A clerk sends a failed application back before later checks; `?` does not write an explanation | Separate parsing from positive-minute validation |
| 34 | Standard input, `read_line`, EOF and `trim` | A reception channel can carry text, a blank line, or close; bytes are not letters | Read a title and distinguish empty, ended and failed input |
| 35 | `std::env::args`, `.next()`, `cargo run --`, quotes | An order handed over at entry; the shell splits the written line first | Accept `add` and a title containing spaces |

Each locale introduces notation before runnable code and provides a prediction, visible result, recall map, required exercise, hint and inline reference answer. The first Rust example in lesson 34 is checked with an ended input stream; its exercise asks readers to try interactive input as well. Lesson 35's first example is checked without arguments; the exercise supplies arguments. Additional boundary checks execute answer programs with actual input and arguments in all three languages. The Kazakh working terms appear in [rust-kz-style.md](rust-kz-style.md).

Technical rules were checked against the official Rust Book and library documentation: [`Option`](https://doc.rust-lang.org/book/ch06-01-defining-an-enum.html#the-option-enum), [`Result` and `?`](https://doc.rust-lang.org/book/ch09-02-recoverable-errors-with-result.html), [`read_line`](https://doc.rust-lang.org/std/io/struct.Stdin.html#method.read_line), [command-line arguments](https://doc.rust-lang.org/book/ch12-01-accepting-command-line-arguments.html), and [`cargo run --`](https://doc.rust-lang.org/cargo/commands/cargo-run.html). Automated checks cannot replace a beginner's hands-on run of the terminal exercises.
