# Rust 16–20: sequence and terminology review

Working drafts on kz/ru/en. These lessons are not in the publication manifest and are not live.

| Lesson | First new notation or idea | Explained before use | Familiar image and its limit | Independent result |
|---|---|---|---|---|
| 16 | `#[test]`, `assert_eq!`, `cargo test` | Each symbol and the two macro arguments appear before the first code block | Reweigh a parcel; a few checked parcels do not establish all cases | Test 0, 15, and 16 for `is_short` and observe a deliberate failure |
| 17 | No new Rust syntax | Recalls conditions, loops, arrays, functions and tests from 1–16 | Data → rule → traversal → counters → output → tests | Rebuild a two-result summary and test empty, boundary and ordinary cases |
| 18 | Buffer, `String`, `::`, `push_str`, stack, heap | Buffer precedes the map; code notation precedes the first `String` example; memory regions precede the scope example | Labelled folder and expandable storage; physical placement is not promised | Draw scopes and explain the string's availability after an inner block |
| 19 | Move, `Copy`, `clone()`, `mut` parameter | Move precedes the first example; `clone()` precedes its example; parameter `mut` precedes the exercise | One claim ticket; bytes need not physically move | Pass an owned `String` into a function and return it without cloning |
| 20 | `&`, `&String`, immutable borrowing | Reference creation and type are explained before the first example | Lend a book to read; Rust enforces access and use periods | Read one `String` twice through references and use the owner afterward |

All three languages include the same required tasks, expected output, deliberate compiler failures where applicable, hints, reference answers, recall prompts and routes back to prerequisites. Kazakh examples and visible output are localized. The terminology in [rust-terms.json](rust-terms.json) is working terminology, not a claim of official approval. Kazakh and English wording was reviewed in the authoring pass; automated checks additionally verify structure and code, but cannot measure an individual learner's comprehension.

The first testing example deliberately keeps tests in the same file. This removes `mod`, `use`, and `#[cfg(test)]` from lesson 16; those constructs can be introduced when modules are taught. New symbols are explained before a learner is asked to run their first example.
