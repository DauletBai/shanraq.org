# Rust course audit — 2026-09-23

Scope: preface, lessons 1–15 in kz/ru/en, the syllabus, publication manifest, answer files, Cargo snapshots, and working drafts 16–20 in all three languages. The authoring pass includes Kazakh and English prose editing and a lesson-by-lesson search for unexplained notation. Automated checks still measure code behavior rather than learner comprehension.

## Findings

1. **The course is not yet a complete route to the promised organizer.** The public lessons end at arrays (15). Files, JSON persistence, recovery, and the finished program are planned for 49–60. The public description must continue to say that new lessons are being released; it must not promise a finished, currently available project.
2. **The next learning step needs a self-check before ownership.** Lesson 16 introduces executable tests, 17 asks for independent reconstruction, and 18–20 introduce memory, moves, and shared references. Each draft has an image or concrete model, a compact recall map, prediction before execution, a required exercise, a hint, an explained answer, and a return path. This sequence closes the jump from basic arrays to ownership in the syllabus.
3. **Lessons 1–15 had first-use gaps.** The editorial pass now defines the purpose and meaning of new commands, signs and types before their first runnable example. The largest gaps were shell arguments in lesson 4; Cargo flags and generated syntax in 6; `&str` and references in 9; arithmetic and overflow in 10; comparisons and lazy boolean operators in 11; `if` block values in 12; ranges in 13; function parameters and return notation in 14; and array type, indexing and method call notation in 15. Familiar images accompany the less concrete rules, and the limits of those images are stated.
4. **The next batch remains a draft.** Lessons 16–20 have Kazakh, Russian, and English versions, localized answers, and Cargo snapshots. Their text has been edited in this pass, but they are deliberately absent from the live database and release manifest. Release needs the next cumulative manifest and successful CI for its final source commit.
5. **Automated checks have a limited scope.** The checker compiles and runs Rust blocks, checks expected output and errors, and executes `#[test]` functions. This catches code drift and failing test assertions. It cannot establish that a beginner understands the explanation or that installation instructions work on every computer. A later learner walkthrough would add evidence; it is not a substitute for the present authoring review.

## Checks and boundaries

The release block 11–15 passed 166 compiled examples, reference answers and Cargo snapshots, 30 additional boundary cases in three locales, local link/language checks, and three publication preparation tests. Four required push workflows for the source commit succeeded. The published 48 localized bodies and summaries matched source SHA-256 values. All 15 newly public pages and their navigation loaded anonymously; lesson 16 returned 404. See [publication record](rust-publication.md).

All three versions of 16–20, their deliberate failures, reference answers, and independent Cargo snapshots are included in the extended checker. Tests in lessons 16 and 17 run under Rust's test harness. Repeat code and language checks after any editorial changes.
