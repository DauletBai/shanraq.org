# Rust course: first publication

Release group: preface plus lessons 1–5, with original kz/ru/en texts. The remaining lessons are not public pages. The Rust course is free and explicitly in progress.

Prepare the SQL and expected hashes from reviewed sources:

```sh
python3 tools/course/prepare_rust.py --sql /tmp/rust-first-batch.sql --expected /tmp/rust-first-batch.json
python3 -m unittest discover -s tools/course -p test_rust_release.py
```

Apply only after deploying migration 20260917000300 and the Rust-aware application. Back up the production database first. The prepared SQL uses one transaction and creates series membership before COMMIT: no lesson becomes visible in the ordinary feed between creation and membership. It requires the existing course author's exact account and refuses collisions with another author/course. It updates only the first batch and does not delete future lessons. The fixed first-batch manifest must be extended with a new release procedure before later batches.

Rust's initial conceptual exercises use the in-lesson reference answers. There is no server Rust runner or parser. Rust formatting/check endpoints refuse unsupported requests; the UI does not offer Go validation for Rust. The existing Go/Python/SQL workflow remains available.

Post-publication: verify 18 translated bodies and summaries against expected SHA-256 hashes; check /course/rust and all /read/rust-* pages in each language; inspect next/previous links; confirm five numbered lessons plus preface; verify home/latest/top feeds omit all six article slugs and /read/rust-06-cargo remains unavailable. Check responsive rendering and the existing Go book links.

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
