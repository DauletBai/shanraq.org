# What Rust is and why it exists

_Summary:_ **Explain the choice of language without promising absolute safety.**

Editorial draft. Prerequisite: previous lessons; no programming experience for lesson 1.

## Why this matters

Imagine a list of things to do on paper. You count the unfinished items by hand. A **program** can follow instructions to do that for you. An **algorithm** is a step-by-step method for solving a problem. A recipe is a useful comparison, but a person may fill in a missing step; a computer needs precise instructions.

**Rust** is a programming language. You write **source code**, text following that language's rules. A **compiler** checks the code and transforms it into a form the computer can execute. Think of a translator who also checks rules. The comparison has a limit: the compiler does not know whether you intended to count your tasks differently.

## The whole picture

Source text → checking and building → running → a result. Editing the source does not automatically update an executable you built earlier. With three tasks and one completed task, we expect two remaining tasks. A program that adds instead of subtracting may faithfully carry out the wrong instruction. We therefore need checks of the result, not only successful compilation.

## Why Rust was created

Graydon Hoare started Rust; Mozilla later supported its development. Browsers and other large programs need efficiency, control over resources and reliability. **Memory** holds data while a program runs. Reading data after its storage has been released, or changing shared data without appropriate coordination, can cause serious problems. Rust's designers sought ways to prevent classes of these errors while retaining control over resources.

Later we will study ownership: who is responsible for data and when it may be used. You do not need to memorise those rules yet. Rust checks many memory-access rules before execution. It does not guarantee correct calculations, successful file writes or correct handling of every input.

Rust is used for command-line tools, servers, system components and devices. It is not only for operating systems and is not automatically the best tool for every job. Our organizer lets us learn the language without simultaneously building a web server or learning processor architecture.

## Recall map

Problem → precise actions → source code → compiler → executable → check the result. The compiler checks language rules; people and tests check the intended behaviour.

## Warm-up

1. Predict: can the compiler know that you intended subtraction rather than addition?
2. Complete: source code is ___; an executable is produced by ___.
3. Repair the claim: “If Rust compiled it, it has no bugs.”

## Exercise

**Required.** Explain program, source code and compiler in your own words. Describe how to find the remaining tasks when there are three tasks and one is complete. Expected result: two remain, with an explanation of why compilation alone does not prove this answer.

**Your own data.** Choose another total and a completed count no greater than the total.

## Hint and reference answer

Name the inputs, action and expected output separately. Reference: 3 − 1 = 2. Source code is the written instructions; the compiler checks rules and transforms them; the program performs the actions as written. Compilation is a valuable check, not proof of every design decision.

Continue when you can distinguish the text of a program from its result without copying these definitions.

## Sources

[History from Rust's original author](https://rustfoundation.org/media/10-years-of-stable-rust-an-infrastructure-story/); [Rust 1.0 goals](https://blog.rust-lang.org/2015/05/15/Rust-1.0/).

[Contents](README-en.md) · [Next lesson](02-organizer-en.md)
