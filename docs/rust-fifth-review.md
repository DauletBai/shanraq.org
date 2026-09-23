# Rust lessons 21–25: first-use and language review

Working drafts in kz/ru/en; absent from the publication manifest and the live course. The group follows the published lesson 20. Each new notation is defined before the first runnable code that uses it.

| Lesson | First new idea | Familiar image and its limit | Independent task |
|---|---|---|---|
| 21 | `&mut String`, exclusive changing access and the last use of a reference | Temporarily close a reference book for editing; Rust checks code uses, not physical passes | Append a label twice through one mutable reference at a time |
| 22 | `String` versus `&str`, `.as_str()` | Owned notebook versus an open page; literals are also `&str` | Read a stored title and a literal through one parameter |
| 23 | `&[u32]`, `&array[a..b]`, owner boundary | Show two lines of a list without copying the sheet; Rust enforces ownership and bounds | Sum two separate parts of one array |
| 24 | `Vec<String>`, angle brackets, `push`, `remove`, borrowing during traversal | Add or remove lines on a list; indices move and are not permanent IDs | Add three titles, remove the middle, traverse the remaining two |
| 25 | UTF-8 bytes, `char` values, graphemes, `chars().count()`, Unicode escape | Count a written address by lines or letters; visual units do not predict bytes | Measure a single non-ASCII letter and a combined accent |

All three locales contain localized visible output, predictions, recall maps, deliberately failing examples with a diagnostic code, required exercises, hints, inline answers, and routes to prerequisites. The Kazakh text uses the working terms recorded in [rust-kz-style.md](rust-kz-style.md), with unfamiliar code terms explained where they first appear. The official Rust Book chapters on [borrowing](https://doc.rust-lang.org/book/ch04-02-references-and-borrowing.html), [slices](https://doc.rust-lang.org/book/ch04-03-slices.html), [vectors](https://doc.rust-lang.org/book/ch08-01-vectors.html), and [UTF-8 strings](https://doc.rust-lang.org/book/ch08-02-strings.html) support the technical rules.

The checker compiles each complete example and answer, verifies expected output and intentional failures, and runs each independent Cargo snapshot. It cannot establish that every beginner will understand a metaphor on first reading. The next publication requires the five-lesson manifest, exact-commit CI, backup, and public-page checks.
