# Rust course publication

Current manifest: preface and lessons 1–15 in kz/ru/en (third group: 11–15, prepared for publication). The full free course remains in progress; later lessons must not be presented as published.

Prepare a cumulative transaction from reviewed sources:

```sh
python3 tools/course/prepare_rust.py --sql /tmp/rust-release.sql --expected /tmp/rust-release.json
python3 -m unittest discover -s tools/course -p test_rust_release.py
python3 tools/coursecheck/rustcheck.py
```

Before publishing, list every push workflow for the exact source commit and wait for every run to succeed. Required checks are CI, Course, docker-smoke and Rust course checks; do not filter the list to only the Rust-specific workflows. Cancelled, skipped or missing required runs do not count as success. Back up production before activation and content publication. The application must include migration 20260917000300 and Rust-aware rendering. The Rust footer link preserves the reader's selected language.

The SQL uses one transaction with an advisory lock. Course membership and all three localized versions become visible together. Author and course collisions stop the transaction; unrelated content is not deleted. The manifest lists all released groups cumulatively, with one preface and numbered lessons at positions 10, 20, and so on. Complete groups of five and all three languages are required. Register every new source in lesson-slugs.json too.

Rust exercises use inline reference answers. There is no server Rust runner; unsupported formatting/check requests are refused instead of reaching the Go checker. Local validation compiles every complete Rust example and answer, checks intentional compiler failures, and builds independent Cargo snapshots for all three languages. Snapshot locations: course/rust-organizer/step-NN (ru), en/step-NN, kz/step-NN. The first five lessons also need manual verification of their installation and navigation steps.

After this release, compare all 48 body/summary hashes with the prepared sources. Visit each page anonymously, verify fifteen numbered lessons plus the preface, forward/back links and the footer in all languages. Check main/latest/top feeds for leaks; lesson 16 must remain unavailable. Check desktop, tablet and both mobile orientations, and verify existing Go book links and readiness.

## Second group — published

Lessons 6–10 explain Cargo commands, program structure, diagnostics, variables and types, then numeric calculations. Each includes recall, a required exercise, a hint, an inline reference and expected output. Kazakh and English have original localized explanations, outputs, answers and Cargo snapshots. Public activation and final checks are recorded below.

## Verified publication — 2026-09-17

Application commit: `158223cd3024e07c200ff34ed4959ccffd09f71f`. CI, docker-smoke and the three-platform Rust course checks succeeded before activation. The final terminology clarification was checked locally and published as database content; it does not change the deployed application binary.

The initial deployment probe used a wrong sample URL and stopped before writing course data. After checking the actual readiness and Go-sample URLs, the prepared Rust content transaction completed successfully. No lesson was exposed outside its course during publication.

Verified results:

- Six course items: one unnumbered preface and exactly five numbered lessons.
- Eighteen bodies and summaries match reviewed source SHA-256 hashes.
- All three localized course pages and all lesson pages load without login.
- Main/latest/top feed checks in all three languages contain none of the lesson URLs.
- Lesson six remains unavailable; no public links lead to repository Markdown files.
- Twenty-four browser checks: course and installation lesson, three languages, four viewports (1366×900, 768×1024, 390×844, 844×390). No horizontal page overflow; desktop Kazakh contents and mobile Kazakh lesson were also visually inspected.
- Existing Go book product and site readiness remain available.

Production backup: `/var/backups/shanraq/rust-course-20260917T162153Z` (database plus rollback container/image material). Runtime activation restarted only the application service.

Public course: https://shanraq.org/course/rust?lang=kz | https://shanraq.org/course/rust?lang=ru | https://shanraq.org/course/rust?lang=en

Checks: https://github.com/DauletBai/shanraq.org/actions/runs/35245165737 (site), https://github.com/DauletBai/shanraq.org/actions/runs/35245165785 (container), https://github.com/DauletBai/shanraq.org/actions/runs/35245166023 (Rust).

## Verified second publication — 2026-09-17

Deployed application commit: `d386641499f50345adacd2ad51d0e94f3143507a`. Final content and checker fix: `575fd8cc3552c17e39c8b0f8498d010db9700815`. The subsequent documentation-only commit records results and publication statuses without changing application code or lesson bodies.

The first activation missed a failing fourth workflow, Course. Its Kazakh prose check rejected a word in lesson 6. The phrase was corrected, the common language checks were included in rustcheck, and Rust fences (including compile_fail) were added to the common checker with regression tests. The final content correction was published after the complete workflow list succeeded. Future release verification must enumerate every push workflow for the exact commit.

- All four workflows succeeded for the final content commit: CI (format, vet, vulnerability scan, PostgreSQL integration tests and race detector), Course (shared course checks), docker-smoke and Rust course checks on Windows/macOS/Linux. Shared language validation passes for all 375 course Markdown files.
- All 82 runnable examples, deliberate compiler failures, answers and localized Cargo snapshots passed. Inline answers are checked against their source files and expected output.
- The Rust footer is now an active language-preserving link. Twelve browser navigation checks followed it successfully, across three languages and four viewports.
- The course contains one preface and exactly ten numbered lessons. All 33 published bodies and summaries match the committed sources: see [expected hashes](rust-second-release-hashes.json).
- Anonymous access, previous/next navigation and the course catalog work in all three languages. Lesson 11 remains unavailable; no links point to repository Markdown files.
- Main/latest/top and second-page feeds contain no Rust lesson URLs. Existing Go book product and readiness checks pass.
- Twenty-four responsive page checks (course and lesson 9) have no horizontal overflow: [browser report](rust-second-browser-checks.json). Footer views were also captured, and representative mobile/desktop pages inspected.

Application-release backup: `/var/backups/shanraq/rust-course-20260917T170737Z`. Content-correction backup: `/var/backups/shanraq/rust-content-20260917T172928Z`. Only the application service was restarted for the footer change. The later wording correction updated one translation, guarded by its previous SHA-256, without restarting services.

Checks: [site](https://github.com/DauletBai/shanraq.org/actions/runs/35251890391), [container](https://github.com/DauletBai/shanraq.org/actions/runs/35251890312), [Rust](https://github.com/DauletBai/shanraq.org/actions/runs/35251890321), [Course](https://github.com/DauletBai/shanraq.org/actions/runs/35251890327).
