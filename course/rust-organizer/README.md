# Organizer learning snapshots

Independent Cargo projects for lessons 06–15; these demonstrate language basics, not the completed organizer. Each `src/main.rs` matches the first runnable Rust block in its language version. The root `step-NN` projects use Russian output; `en/step-NN` uses English and `kz/step-NN` uses Kazakh. All three sets are independently built and checked. `expected.txt` records stdout. Exercise answers are in `course/lessons/rust/answers`.

For a snapshot, enter its directory and run `cargo run --locked --offline`. For all checks, run `python3 tools/coursecheck/rustcheck.py` from the repository root. The checker builds in a temporary directory and leaves learner snapshots untouched. Rust 1.97.0 and Edition 2024 are the checked baseline.
