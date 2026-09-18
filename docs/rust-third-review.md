# Rust 11–15: editorial and publication review

Prepared 2026-09-18. This report distinguishes prepared material from verified production publication. Publication results belong in rust-publication.md.

## Sequence and prerequisites

| Lesson | New learning | Reused learning | Required independent result |
|---|---|---|---|
| 11 | Comparisons, &&, logical OR, !, short-circuit order | bool introduced in 9; arithmetic in 10 | Select only unfinished, non-archived tasks lasting at most 15 minutes |
| 12 | if/else if/else, block values, unit | bool and comparison | Reject done > total before subtraction; handle equality and a positive remainder |
| 13 | for, ranges, while, loop, break, continue, accumulation | mut, conditions, arithmetic | Plan days 1–5 excluding day 3, total 60 minutes |
| 14 | Typed parameters, positional arguments, return, preconditions | block values, unit, comparisons | Return whether a task lasts 1–15 minutes, without printing inside the function |
| 15 | Fixed arrays, index zero, usize, len, bounds, panic | functions, selection, iteration | Count short tasks and total their duration without counting zero |

Every locale explains the same symbols, starting directory, changed file, run command, expected output, intentional errors, boundary cases, required exercise, hint and inline reference. All reference answers also exist as independent source and output files. All 15 primary examples have independent Cargo snapshots. No external crates, input parsing, ownership prerequisites, unwrap, or unintroduced formatting traits are required.

Potential gaps resolved: = versus ==; inclusive OR; short-circuit guard before division; unit introduced with block values before function returns; declaration order versus call order; parameter names versus argument positions; loop exits and Ctrl+C; inclusive versus exclusive range ends; len introduced as a method before its first use; array length versus maximum index; invalid inputs never silently become a valid zero remainder. Function preconditions and deferred Result/Option handling are explicitly stated. Copying numeric arrays is not generalized to strings.

Kazakh text was authored and reviewed with the same lesson structure, without a translation service. Output labels, exercises and reference outputs are localized. Existing terms компилятор and жиым are retained. Working terminology is not represented as officially approved. Automated language checks supplement this review; they do not establish native-reader validation.

## Verification sources

Consulted 2026-09-18 for semantics; lesson explanations and examples are original:
- https://doc.rust-lang.org/reference/expressions/operator-expr.html#lazy-boolean-operators
- https://doc.rust-lang.org/book/ch03-05-control-flow.html
- https://doc.rust-lang.org/book/ch03-03-how-functions-work.html
- https://doc.rust-lang.org/book/ch03-02-data-types.html

## Publication boundary

Cumulative release: preface + 15 numbered lessons, 48 localized pages. Lesson 16 remains unpublished. Membership and content are inserted in one transaction so lessons cannot temporarily appear in article feeds. Required workflows: CI, Course, docker-smoke and Rust course checks for the same commit. Production verification must compare all source hashes and navigation, inspect responsive pages and check main/top/latest feeds anonymously.

## Local verification

- `python3 tools/coursecheck/rustcheck.py`: 166 examples, intentional failures, answers and Cargo snapshots passed.
- `python3 tools/coursecheck/rust_boundaries.py`: exercise variations check archived/completed/overlong tasks, equal and invalid counts, empty iteration, positive duration boundary, no/all qualifying array values and valid array access in all three languages.
- Language checks passed for all 51 Rust Markdown files; offline checks found no broken internal links.
- `go test ./pkg/modules/articles ./pkg/modules/ai` passed for the public-description and assistant-guidance wording changes.
