# Lifetimes: how references relate

_Summary (summary):_ **Lesson 47. Tell when a returned reference may point into an input.**

First review [references](/read/rust-20-borrowing?lang=en), [slices](/read/rust-23-slices?lang=en), and [generics](/read/rust-45-generics?lang=en). Data still lasts only for one run.

## Familiar image and recall map

A library card lets you borrow a book while the library still has it. A reference similarly gives access to another value, but does not keep that value alive. A **reference lifetime** is the part of a program in which that borrow can be used; the compiler determines it from where values are created and used. The image has a limit: these are code regions, not clock times or manual timers.

**input references → shared 'a relationship → returned reference valid only while the relevant inputs remain usable.** In `fn longer<'a>(left: &'a str, right: &'a str) -> &'a str`, the apostrophe and letter name a relationship between references, not a separate value. The letter is arbitrary; repeating it connects the output to either possible input. The annotation does not extend either value’s life. In lesson 45 there was only one borrowed input, so the compiler inferred this relationship. Here two inputs can supply the result, so we write it.

## Choose text within a safe region

```rust
fn longer<'a>(left: &'a str, right: &'a str) -> &'a str {
    if left.len() >= right.len() {
        left
    } else {
        right
    }
}

fn main() {
    let short = String::from("task");
    let result;
    {
        let long = String::from("shopping");
        result = longer(&short, &long);
        println!("{result}");
    }
    println!("{short}");
}
```
```text
shopping
task
```

`result` is used inside the inner braces while both strings exist. The closing brace destroys `long`; using `result` afterward is invalid, even if another call chose `short`: the function signature permits the other choice. `.len()` counts bytes, not visible letters.

A function cannot return a reference to its local string. The deliberately invalid example below produces E0515: the string is destroyed when the function returns. Naming `'a` does not keep it alive.

error-code: E0515
```rust,compile_fail
fn dangling<'a>() -> &'a str {
    let text = String::from("temporary");
    text.as_str()
}

fn main() {}
```

## Recall without looking

1. What does the repeated `'a` connect?
2. Does the annotation keep `long` alive longer?
3. Why can we not print `result` after the inner brace?

## Exercise

**Required.** After the inner block, create `chosen = longer(short.as_str(), "schedule")` and print it. A string literal stays available for the whole program, and `short` still exists. Predict the output. Do not move the earlier `result` outside.

## Answers

`as_str()` borrows the text in `short` without making another string. Use `chosen` before `main` ends.

<!-- task-answer -->
```rust
fn longer<'a>(left: &'a str, right: &'a str) -> &'a str {
    if left.len() >= right.len() {
        left
    } else {
        right
    }
}

fn main() {
    let short = String::from("task");
    let result;
    {
        let long = String::from("shopping");
        result = longer(&short, &long);
        println!("{result}");
    }
    println!("{short}");
    let chosen = longer(short.as_str(), "schedule");
    println!("{chosen}");
}
```
```text
shopping
task
schedule
```

## After checking

Returning a reference to a `String` created inside a function fails to compile: that value disappears when the function returns. If scopes are unclear, revisit [lesson 19](/read/rust-19-ownership?lang=en). [Official lifetimes chapter](https://doc.rust-lang.org/book/ch10-03-lifetime-syntax.html).

[Previous lesson](/read/rust-46-traits?lang=en) · [Contents](/course/rust?lang=en)
