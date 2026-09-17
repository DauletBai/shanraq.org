# The organizer we will build

_Summary:_ **Distinguish the future demonstration from the program available now.**

Editorial draft. Prerequisite: previous lessons; no programming experience for lesson 1.

## Why this matters

A builder needs a plan before the first wall. Our destination is an organizer that keeps tasks and notes on your computer. A **record** describes one item; a **command** asks a program to do something; a **scenario** is a sequence of user actions towards a goal.

## The whole picture

This is a mock-up of the future interface, not commands you can run today. You do not yet have the organizer executable.

```text
organizer add "Study the Rust lesson"
organizer list
organizer done 1
```

The first line adds a task, the second lists tasks and the third completes the item identified by 1. That number will eventually be a stable identifier, not its current row on screen. The final program will also edit descriptions and priorities, search, filter, calculate statistics and save data.

An empty list should produce a helpful message. An invalid command should explain the problem without changing other records. These are part of the result, not optional improvements left until the end.

## How we will get there

First we print a message. Then we calculate, choose actions and repeat them. We introduce our own task type, commands and, later, file storage. Early versions demonstrate individual language features; they do not implement the mock-up yet.

The first interactive version keeps data only in memory. Closing it loses those records. The file lessons remove that limitation. We state this in advance so you do not discover it after entering important notes.

A website, accounts, cloud synchronization and scheduled reminders are outside the required project. You can read on a phone, but the full build-and-file workflow needs a computer. Installation needs an internet connection; normal use of the finished local organizer does not.

## Recall map

Message → data → actions → commands → persistence → checks → release. The finished picture shows our direction, not code you must already understand.

## Warm-up

1. Predict: does a mock-up save anything on your computer?
2. Complete: a scenario is a sequence of ___ towards a ___.
3. Repair the requirement “delete every task whenever an error occurs.”

## Exercise

**Required.** Write four steps: add “Study the Rust lesson”, see the record, mark it complete and see its completed state. Add a fifth step: try to complete a nonexistent task. Expect an explanation and no change to the existing record.

**Your own data.** Invent a task containing spaces and letters from a language you use.

## Hint and reference answer

For each step, state both the action and the observable result. Reference: one incomplete record after adding; the same record completed after finishing; unchanged data after an invalid identifier. A demonstration line in a textbook does not create an actual stored record.

Continue when you can distinguish the final goal from the current stage's capabilities.

[Previous lesson](01-why-rust-en.md) · [Contents](README-en.md) · [Next lesson](03-workspace-en.md)
