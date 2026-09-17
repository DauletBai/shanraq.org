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
