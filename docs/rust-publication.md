# Rust course publication

Live site: preface and lessons 1–40 published in kz/ru/en. The prepared ninth release manifest includes lessons 41–45, but they remain absent from the live site until the guarded transaction is applied. The full free course remains in progress; later lessons must not be presented as published.

Prepare a cumulative transaction from reviewed sources:

```sh
python3 tools/course/prepare_rust.py --sql /tmp/rust-release.sql --expected /tmp/rust-release.json
python3 -m unittest discover -s tools/course -p test_rust_release.py
python3 tools/coursecheck/rustcheck.py
```

Before publishing, list every push workflow for the exact source commit and wait for every run to succeed. Required checks are CI, Course, docker-smoke and Rust course checks; do not filter the list to only the Rust-specific workflows. Cancelled, skipped or missing required runs do not count as success. Back up production before activation and content publication. The application must include migration 20260917000300 and Rust-aware rendering. The Rust footer link preserves the reader's selected language.

The SQL uses one transaction with an advisory lock. Course membership and all three localized versions become visible together. Author and course collisions stop the transaction; unrelated content is not deleted. The manifest lists all released groups cumulatively, with one preface and numbered lessons at positions 10, 20, and so on. Complete groups of five and all three languages are required. Register every new source in lesson-slugs.json too.

Rust exercises use inline reference answers. There is no server Rust runner; unsupported formatting/check requests are refused instead of reaching the Go checker. Local validation compiles every complete Rust example and answer, checks intentional compiler failures, and builds independent Cargo snapshots for all three languages. Snapshot locations: course/rust-organizer/step-NN (ru), en/step-NN, kz/step-NN. The first five lessons also need manual verification of their installation and navigation steps.

For each future release, compare every body/summary hash with the prepared sources. Visit each page anonymously, verify navigation in all languages, and check main/latest/top feeds for leaks. Check desktop, tablet and both mobile orientations, and verify existing Go book links and readiness.

## Verified eighth publication — 2026-09-23

Lessons 36–40 are public in kz/ru/en. Release source commit: `d1906344d3c122e2a7f900d4c94132a325d07a33`. All four required push workflows succeeded for that exact commit: [CI](https://github.com/DauletBai/shanraq.org/actions/runs/35880765223), [Course](https://github.com/DauletBai/shanraq.org/actions/runs/35880765168), [docker-smoke](https://github.com/DauletBai/shanraq.org/actions/runs/35880765109), and [Rust course checks](https://github.com/DauletBai/shanraq.org/actions/runs/35880765119). Local checks passed 544 Rust examples, answers and Cargo snapshots, 66 additional boundary cases, 126 Rust language files, 483 internal links and release tests.

Lesson 36 returned 404 in all three languages before activation. A 4,489,458-byte database dump was saved at `/var/backups/shanraq/rust-content-20260923-before-36-40.dump` and its archive listing was checked with `pg_restore` inside the database container. The cumulative SQL completed in one transaction. The course has the preface and 40 numbered lessons. `publish.py --check` found no source differences. After an initial network timeout, a full anonymous HTTP pass succeeded for all 15 new pages, three course pages, three prefaces, navigation including 35→36, main/latest/top feeds, lesson 41 remaining 404, the Go course, `/readyz` and `/healthz`. The application was not restarted. Responsive browser inspection was not rerun in this pass.

## Verified seventh publication — 2026-09-23

Lessons 31–35 are public in kz/ru/en. The exact release source commit is `c2599986b28bd2e87ecb464e4f2927ca1fd1ec4f`. All required push workflows succeeded: [CI](https://github.com/DauletBai/shanraq.org/actions/runs/35871731345), [Course](https://github.com/DauletBai/shanraq.org/actions/runs/35871731478), [docker-smoke](https://github.com/DauletBai/shanraq.org/actions/runs/35871731499), and [Rust course checks](https://github.com/DauletBai/shanraq.org/actions/runs/35871731560). Local release tests, 484 examples, answers and snapshots, and 60 additional boundary cases passed before publication.

Lesson 31 returned 404 before the transaction. A 4,437,448-byte database dump was saved at `/var/backups/shanraq/rust-content-20260923-before-31-35.dump` and its archive listing was checked with `pg_restore` inside the database container. The cumulative SQL completed in one transaction. The course contains the preface and 35 numbered lessons. `publish.py --check` found no source differences. All 15 new pages, three course pages and three prefaces load anonymously. Navigation, including 30→31, works; lesson 36 remains 404. Checked main/latest/top feeds contain no new lesson URLs. The Go course, `/readyz` and `/healthz` return 200. The application was not restarted. Responsive browser inspection was not rerun in this pass.

## Verified sixth publication — 2026-09-23

Lessons 26–30 are public in kz/ru/en, along with an evergreen preface and a clear terminal-program scope in lesson 2. Release source commit: `773990cc58b3fc974cc98fb930d677a0f9ff6732`. All four required push workflows succeeded for that exact commit: [CI](https://github.com/DauletBai/shanraq.org/actions/runs/35867824716), [Course](https://github.com/DauletBai/shanraq.org/actions/runs/35867824533), [docker-smoke](https://github.com/DauletBai/shanraq.org/actions/runs/35867824627), and [Rust course checks](https://github.com/DauletBai/shanraq.org/actions/runs/35867824562). Local release tests, 415 Rust examples, answers and snapshots, language checks and offline links passed before publication.

Lesson 26 returned 404 before the transaction. A 4.4 MB database dump was saved at `/var/backups/shanraq/rust-content-20260923-before-26-30.dump` and its archive listing was checked with `pg_restore` inside the database container. The cumulative SQL completed in one transaction. The course contains the preface and 30 numbered lessons; `publish.py --check` found no source differences. All 15 new pages, three course pages and the three revised prefaces load anonymously. Navigation, including 25→26, works; lesson 31 remains 404. Checked main/latest/top feeds contain no new lesson URLs. The Go course, `/readyz` and `/healthz` return 200. The application was not restarted. Responsive browser inspection was not rerun in this pass.

## Verified fifth publication — 2026-09-23

Lessons 21–25 are public in kz/ru/en. Release source commit: `c75386f3661d9546c8a5d1418b4804def210c7f0`. All four required push workflows succeeded for that exact commit: [CI](https://github.com/DauletBai/shanraq.org/actions/runs/35862383075), [Course](https://github.com/DauletBai/shanraq.org/actions/runs/35862383102), [docker-smoke](https://github.com/DauletBai/shanraq.org/actions/runs/35862383011), and [Rust course checks](https://github.com/DauletBai/shanraq.org/actions/runs/35862382948). Local checks passed 328 Rust examples, answers and snapshots, 36 boundary cases, 81 language files, 781 internal links, and the release tests.

Before the transaction, lesson 21 returned 404 and a 4.4 MB database dump was saved at `/var/backups/shanraq/rust-content-20260923-before-21-25.dump`. The cumulative SQL completed in one transaction. The course now contains the preface and 25 numbered lessons. The publication source checker found no differences between site content and repository lessons. All 15 new pages and three course pages load anonymously; previous/next links, including 20→21, work in all languages. Lesson 26 remains 404. The checked main/latest/top feeds contain no new lesson URLs. The Go course, `/readyz`, and `/healthz` return 200. The application was not restarted. Responsive browser inspection was not rerun in this pass.

## Verified fourth publication — 2026-09-23

Lessons 16–20 are public in kz/ru/en. Release source commit: `78564f95157b7c301d963b11b1883d14151ee799`. All four push workflows succeeded for that exact commit: [CI](https://github.com/DauletBai/shanraq.org/actions/runs/35830931052), [Course](https://github.com/DauletBai/shanraq.org/actions/runs/35830930962), [docker-smoke](https://github.com/DauletBai/shanraq.org/actions/runs/35830931002), and [Rust course checks](https://github.com/DauletBai/shanraq.org/actions/runs/35830931004). Locally, the checker passed 247 examples, answers and Cargo snapshots; 30 additional boundary cases, 66 Rust language files, release preparation tests and internal links also passed.

Before the transaction, the new pages returned 404 and a 4.2 MB database dump was saved at `/var/backups/shanraq/rust-content-20260923-before-16-20.dump`. The cumulative SQL completed in one transaction. The course now contains the preface and 20 numbered lessons. All 63 titles, bodies and summaries match the prepared sources. The 15 new pages, course pages and the 15→16→17 and 19→20 links load anonymously in all languages; lesson 21 returns 404. Main latest/top feeds contain no new lesson URLs. The Go course, `/readyz`, and `/healthz` still return 200. The application was not restarted. Responsive browser inspection was not rerun in this publication pass.

## Verified third publication — 2026-09-23

The preface and lessons 1–15 are public in kz/ru/en. Release source commit: `0d63b93c3303e17df1d52704d3f59e1a3568a23b`. The four required push workflows for that commit succeeded: [CI](https://github.com/DauletBai/shanraq.org/actions/runs/35314882908), [Course](https://github.com/DauletBai/shanraq.org/actions/runs/35314882885), [docker-smoke](https://github.com/DauletBai/shanraq.org/actions/runs/35314882896), and [Rust course checks](https://github.com/DauletBai/shanraq.org/actions/runs/35314882931). Locally, `rustcheck.py` passed 166 released examples, answers and Cargo snapshots, `rust_boundaries.py` passed 30 additional cases, and the three release preparation tests passed.

Before the database transaction, the course had 11 items. A database dump was stored at `/var/backups/shanraq/rust-content-20260923-before-11-15.dump` (4.2 MB). The prepared SQL completed as one transaction. Production then had 16 course items and all 48 localized body/summary pairs matched the prepared SHA-256 values. Anonymous HTTP checks passed for all 15 new pages, their previous/next links, and three general feed variants in each language. All three course and lesson 15 pages returned 200; lesson 16 returned 404 in every language. The application was not restarted.

Responsive browser inspection and existing Go book/readiness probes were not rerun in this publication pass; they remain verification work, not completed checks.

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
