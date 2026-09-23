# Rust lessons 26–30: publication review

Published in kz/ru/en on 2026-09-23. Lesson 25 now links to lesson 26, and the group is in the release manifest. Publication evidence is in [rust-publication.md](rust-publication.md).

| Lesson | First new idea | Familiar image and its limit | Independent task |
|---|---|---|---|
| 26 | Tuple positions, `.0`, unpacking and a tuple return type | Two labelled places on one card; tuple places are positional | Return a title and duration together, then unpack them |
| 27 | `struct Task`, named fields and instances | A filled task form; type checks do not enforce sensible minutes | Build and change a task through its fields |
| 28 | `impl`, associated `new`, `&self`, `&mut self`, field shorthand | Actions attached to a card; code only runs when called | Implement construction, reading and finishing methods |
| 29 | `enum`, variants with data, introductory `match` | One state sign at a time; a variant may carry data | Handle `Doing(u32)` and the two other states |
| 30 | Exhaustive `match`, `_`, `if let`, a `match` result | An instruction for every state stamp; wildcard can hide later states | Return active minutes for every status |

Each locale explains new notation before the first runnable example, then asks for a prediction, provides visible output, an intentional diagnostic, recall prompts, a required exercise, a hint and an inline reference answer. The Kazakh version uses the working terms in [rust-kz-style.md](rust-kz-style.md), with familiar images and their limits. The group ends with a 26–30 recall checkpoint. A tuple or struct does not itself validate business meaning; the lessons state that `Result` will supply a later error-reporting path.

Technical rules were checked against the official Rust Book: [tuples](https://doc.rust-lang.org/book/ch03-02-data-types.html#the-tuple-type), [structs](https://doc.rust-lang.org/book/ch05-01-defining-structs.html), [methods](https://doc.rust-lang.org/book/ch05-03-method-syntax.html), [enums](https://doc.rust-lang.org/book/ch06-01-defining-an-enum.html), and [`match`](https://doc.rust-lang.org/book/ch06-02-match.html). Automated checks compile complete examples and answers, verify visible output and intended compiler errors, and run Cargo snapshots. They do not prove that every learner will understand a metaphor without trying the exercises.
