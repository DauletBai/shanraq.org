# Rust lessons 36–40: draft review

Working drafts in kz/ru/en. They are available in the repository for offline review and absent from the public release manifest. Published lesson 35 has no forward link to lesson 36; lesson 40 does not link to an unwritten lesson 41.

| Lesson | New step | Familiar image and its limit | Required independent result |
|---|---|---|---|
| 36 | Parse words into `Result<Command, String>` before acting | A clerk checks an application before fulfilling it; the program accepts defined command forms, not arbitrary sentences | Add `Done(u32)` and reject zero, malformed and extra arguments |
| 37 | Add, rename, complete and remove tasks in `Vec<Task>` | Cards laid in a row; a row position changes after removal | Reject blank renames and missing positions without changing data |
| 38 | Give records stable IDs with a checked counter | A library book keeps its inventory number when moved; IDs are stable only while the counter is maintained | Reject unknown IDs and counter overflow before mutation |
| 39 | Test action sequences and state after errors | Check a lock with several keys and actions; selected tests do not prove absence of every bug | Assert that failed removal preserves records and the next ID |
| 40 | Use the organizer through terminal input in one run | A reception desk handles multiple requests until closing; memory disappears at process exit | Add `rename ID TITLE`, retain spaces, and handle unknown IDs |

Each language has a recall map, explanation before the first runnable example, visible output, required exercise, hint, inline answer and return route. The Kazakh wording uses the working terms in [rust-kz-style.md](rust-kz-style.md). Lesson 40 shows both the empty-input output used by automated checks and an actual multi-command transcript; its answer shows the expected rename transcript. `split_once(' ')` is explained before the code and splits only on the first ASCII space. The parser intentionally does not implement shell quoting, which belongs to the separate launch-argument lesson.

The examples and answers compile and run in all three languages, and each first example matches an independent Cargo snapshot. Lesson 39's tests run under `cargo test`. Additional boundary checks run lesson 40 answers with real input and verify an empty task list. Files are not yet used: every lesson states that task data vanish when the process ends.

Technical behavior was checked against official [enum and match documentation](https://doc.rust-lang.org/book/ch06-02-match.html), [`Vec::remove`](https://doc.rust-lang.org/std/vec/struct.Vec.html#method.remove), [`u32::checked_add`](https://doc.rust-lang.org/std/primitive.u32.html#method.checked_add), [Rust testing guide](https://doc.rust-lang.org/book/ch11-01-writing-tests.html), [`read_line`](https://doc.rust-lang.org/std/io/struct.Stdin.html#method.read_line), and [`split_once`](https://doc.rust-lang.org/std/primitive.str.html#method.split_once). A beginner's manual run of the terminal exercises remains editorial verification work before the next publication.
