# Organizer: choosing a file

1. Open https://www.rust-lang.org/tools/install, install Rust for your OS, open a new terminal, and check `cargo --version`.
2. Run `cargo run` in the project folder.
3. The command path wins over the environment path; otherwise use `tasks.json` in the current folder. This practice example passes values directly in code and does not create a file yet.
4. If Cargo reports an error, read it and check that you ran the command in the folder containing `Cargo.toml`.
