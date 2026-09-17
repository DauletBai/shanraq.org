# Reading your first program

_Summary:_ **Change the output and explain every line.**

Prerequisite: lessons 1–6. Work in organizer unless stated otherwise.

## Why this matters

The last lesson used a supplied program. Now we read its instructions. Open src/main.rs, replace its entire contents with the following, save it and run cargo run from organizer:

```rust
fn main() {
    println!("My organizer");
    println!("Task: learn Rust");
}
```

Output:

```text
My organizer
Task: learn Rust
```

## Every symbol has a purpose

fn introduces a **function**, a named action. main names the function where execution of our ordinary program starts. The empty parentheses contain no declared parameters. Braces delimit a **block**, a group of actions. Lesson 14 will develop functions with parameters; for now recognise this entry-point form.

println! prints text followed by a line break. The ! indicates a **macro**, a construct expanded into code during compilation. Think of an instruction-generating template. It is not an ordinary function, and writing your own macro is outside this lesson. The parentheses hold what is supplied to it.

Text in double quotes is a **string literal**, a value written directly in the source. The quotes delimit it and are not printed. The semicolon ends this **statement**. Four-space indentation makes nesting readable, but braces determine the actual boundaries.

The two statements run in order and produce two lines. // begins a single-line **comment**, explanatory text for a reader through the end of that line; it is not executed.

## An understandable mistake

Removing a closing quote makes the text boundary unclear to the compiler. Removing ! changes the kind of invocation; it no longer calls the macro we learned. Change one thing at a time and repeat cargo check, introduced in lesson 6.

## Recall map

fn main(): entry; braces: block; println!: output; quoted text: literal; semicolon: statement ending; //: explanation for people.

## Warm-up

1. Predict the order of the two printed lines before running.
2. Complete: println___("Text");
3. Restore a missing closing brace in your copy using the first diagnostic.

## Exercise

**Required.** Print three lines in this order: “Today’s plan”, “Learn Rust”, “Review Cargo commands”. Use three statements. The expected output contains these lines without the quotation marks.

**Your own data.** Replace the second line with your own task.

## Hint and reference answer

Each println! ends a line. Try reconstructing the program using only the recall map. If necessary, practise one line first. Warm-up answers: heading before task; the missing symbol is !; the closing brace ends main's block.

<!-- task-answer -->
```rust
fn main() {
    println!("Today’s plan");
    println!("Learn Rust");
    println!("Review Cargo commands");
}
```

Expected output:

```text
Today’s plan
Learn Rust
Review Cargo commands
```

[Previous lesson](/read/rust-06-cargo?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-08-diagnostics?lang=en)
