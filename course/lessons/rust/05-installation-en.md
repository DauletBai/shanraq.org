# Installing Rust and checking the tools

_Summary:_ **Check the installation and distinguish the tools involved.**


## Why this matters

A workshop has tools with different jobs. **rustup** installs and selects their set; **rustc** compiles Rust; **Cargo** coordinates projects, builds, runs and dependencies. Cargo calls rustc rather than replacing it. A **dependency** is a library used by a project; a **library** is reusable code. Think of a ready-made part whose purpose and limitations still need to be understood. A **toolchain** is a compatible set of tools at a particular version. **rust-analyzer** provides editor assistance; suggestions do not replace checking the program.

## Installation for your system

Open the [official installation page](https://rust-lang.org/tools/install/) and follow the branch for your operating system. Do not download installers from advertisements. On Unix the page offers a command that downloads and runs an installation script. This is a program from the internet: verify the address and command on the official page, and do not disable certificate checks to make installation succeed.

**Windows.** Use rustup-init for your architecture. If prompted for Visual Studio C++ tools, follow the provided installation route; C++ components and the Windows SDK are needed. An SDK is a kit for building programs for a platform. The full Visual Studio editor is not required. For this beginner course, choose ordinary Windows installation instead of mixing it with WSL, a separate Linux environment inside Windows.

**macOS.** If Apple command-line tools are missing, xcode-select --install opens their official installer. Complete it, then follow rustup's instructions. Finish installation dialogs before checking Rust.

**Linux.** You may need system build tools in addition to rustup. On Ubuntu/Debian these are commonly supplied by build-essential. sudo apt install build-essential requires administrator permission: sudo requests it, apt manages packages, install selects installation and build-essential names the package. Use your own distribution's instructions on other systems. On a managed computer, ask the administrator; this course does not require bypassing restrictions.

A **linker** combines prepared parts into an executable, like an assembler of a finished product. An error mentioning linker or link.exe may mean missing system tools rather than incorrect Rust code.

## Check the result

Close and reopen the terminal. Run these individually; --version requests a version number:

```sh
rustup --version
rustc --version
cargo --version
```

Each should report a tool name and version. Their numbers need not match. Our examples are checked with Rust 1.97.0 and Edition 2024. The next lesson explains the difference. Start with stable, the stable release channel, rather than nightly, which includes experimental capabilities.

If a command is missing, first reopen the terminal and confirm installation completed. **PATH** is the list of folders searched for executable commands, like a list of places to look. It is not your project folder. Do not edit it blindly; consult [rustup's installation guidance](https://rust-lang.github.io/rustup/installation/).

## Editor and build check

Install VS Code from code.visualstudio.com and the rust-analyzer extension published by rust-lang. An **extension** adds editor features. Another editor is acceptable, but our illustrated instructions will use VS Code.

Version output does not prove that linking works. Check it now from rust-learning. cargo new creates a project; rust-install-check names its folder; --edition 2024 selects the language edition; --vcs none avoids creating Git history. cd enters the folder. cargo run builds when necessary and runs the program. This is an installation probe using a supplied example, not an exercise in writing unexplained source.

```sh
cargo new rust-install-check --edition 2024 --vcs none
cd rust-install-check
cargo run
```

Expect Hello, world! from the program. Compiling and Finished are tool messages. If the folder exists, enter it and repeat only cargo run; do not delete your work. Do not name the probe organizer: that will be our separate learning project. Return to rust-learning using cd .. . Lessons 6–7 will explain the generated files and every source-code symbol. This first group now ends with an actual build check.

## Recall map

rustup: tool set; Cargo: coordination; rustc: compilation; linker: combining; rust-analyzer: editor help. You do not manually run all five each time.

## Warm-up

1. Predict: does cargo --version prove a project has been built?
2. Complete: editor suggestions come from ___; compilation is performed by ___.
3. Repair “edit main.rs when the cargo command is not found.”

## Exercise

**Required.** Save the three version lines in notes.txt, along with your operating system and shell. Expect all commands to be found and be able to distinguish the tools. Complete the build probe above and record its result.

**Your own data.** Note where your editor is installed.

## Hint and reference answer

The answer contains rustup, rustc and cargo version lines, not necessarily identical numbers. For help, include the system, command and complete error text. Remove private details where needed; passwords and access tokens are unnecessary.

[Previous lesson](/read/rust-04-terminal?lang=en) · [Contents](/course/rust?lang=en)

[Next lesson: your first Cargo project](/read/rust-06-cargo?lang=en).
