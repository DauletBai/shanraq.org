# Values and memory

_Lead (summary):_ **Lesson 18. Distinguish a value, its owner, and the storage it may use.**

Prerequisites: lessons 1–17. Revisit [scope](/read/rust-09-variables?lang=en) if you are unsure where a name can be used, or [functions](/read/rust-14-functions?lang=en) if a function call is unclear.

## Why look at memory now?

With numeric arrays, passing values to functions seemed effortless. Text that can grow raises a new question: who is responsible for its storage? **Memory** is where a running program keeps data. A **buffer** is an allocated area for data such as text bytes; its size may change. A variable is a name through which code accesses a value. Picture a label on a folder, but remember the limit of the picture: the label need not be a separate box containing a complete copy of every byte. The compiler may choose a different physical layout if the program behaves the same.

**Recall map:** name → value → perhaps a data buffer → owner → end of scope. Cover this line and redraw it. A buffer is needed for some values, not for every value.

## Meet a growing string before seeing its code

`String` is a Rust type that owns text data and can grow. `String::from("Plan")` creates such a value from a text literal. The double colon `::` chooses a function associated with the type `String`; it does not join two strings. `title.push_str(" for today")` adds text to the end of the value named `title`. The dot calls a method on that value, as introduced with `len()` in lesson 15. `let mut` lets this particular binding change; a `String` binding is not automatically mutable. We will compare string types in lesson 22.

In the `organizer` project from lesson 6, save your previous `src/main.rs` elsewhere. Replace the whole file with the next example and run `cargo run` in the folder containing `Cargo.toml`. No other file or dependency changes. Predict the output; the block after the code shows program output only.

```rust
fn main() {
    let fixed: [u32; 2] = [10, 15];
    let mut title = String::from("Plan");
    title.push_str(" for today");
    println!("{} and {} minutes", fixed[0], fixed[1]);
    println!("{title}");
}
```

```text
10 and 15 minutes
Plan for today
```

The array value contains two `u32` numbers and has a known size. A `String` value keeps information about its text and owns a buffer holding its bytes. That buffer can be replaced with a larger allocation as the text grows. This description explains the rule, not an address you can rely on.

## Stack, heap, and scope

The **stack** is normally used for function-call data of known size. The **heap** can supply storage whose size is chosen while the program runs. As a working picture, imagine the array in a call's work area and the string's bytes in a separate growable storage area. The picture's limit matters: optimization can alter actual placement. Stack versus heap alone does not decide whether you may use a name after giving a value to another function. Ownership rules decide that.

A name's **scope** ends at the matching closing brace `}`. Rust automatically releases an owned resource when it is no longer needed; we do not release this string buffer by hand. Now predict which name remains available after the inner braces:

```rust
fn main() {
    let outer = 10;
    {
        let inner = String::from("note");
        println!("{inner}: {outer}");
    }
    println!("Outer value: {outer}");
}
```

```text
note: 10
Outer value: 10
```

The inner `}` ends the scope of `inner`, but not of `outer`. Using `inner` afterward fails to compile with E0425:

<!-- error-code: E0425 -->
```rust,compile_fail
fn main() {
    {
        let inner = String::from("note");
        println!("{inner}");
    }
    println!("{inner}");
}
```

## Recall check

Does `push_str` change the length of the numeric array? Can `outer` be used after the inner block? Is the drawing of stack and heap an exact map of addresses? Answers: no; yes; no. Explain the last answer without the page.

## Exercise

**Required.** Create a `String` containing `Task`; add `: reading` using `push_str`. Inside an inner block, declare `[5, 10]` and print their sum. Outside the block, print the string. Draw both scopes before coding. Expect `Minutes: 15` followed by `Task: reading`. Run `cargo run` and explain why the array's name is unavailable after the inner block.

## Hint

Declare `let mut title = String::from("Task");` before the inner block. Declare the array inside it. The outer `title` remains in scope after that block.

## Reference answer after your attempt

<!-- task-answer -->
```rust
fn main() {
    let mut title = String::from("Task");
    title.push_str(": reading");
    {
        let minutes = [5, 10];
        println!("Minutes: {}", minutes[0] + minutes[1]);
    }
    println!("{title}");
}
```

```text
Minutes: 15
Task: reading
```

Next we will ask what happens when the `String` value itself is passed to a new name. [Rust's ownership chapter](https://doc.rust-lang.org/book/ch04-01-what-is-ownership.html).

[Previous lesson](/read/rust-17-first-checkpoint?lang=en) · [Contents](/course/rust?lang=en)
