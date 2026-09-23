# Variables, values and types

_Summary:_ **Predict values and distinguish reassignment from a new binding.**

Prerequisite: lessons 1–8. Work in organizer unless stated otherwise.

## Why this matters

A **value** is concrete data such as 2. A **variable** associates a name with a value. Think of a labelled place for data. Later we will see the limit of that image: a name does not always mean a separate box with an independent copy.

A **type** determines valid values and operations. Think of rules for the contents: numbers support arithmetic; text supports text operations. A name alone does not determine its type.

## The whole picture

Replace src/main.rs, save and use cargo run from organizer:

A key before the first example: `let` binds a name to a value, `mut` permits later assignment through that name, `: u32` explicitly selects a nonnegative integer type, and `=` assigns a value. `;` ends a statement. Inside `println!`, `{count}` inserts the current value into the text; it is not another program block. Picture a labelled card with a number that may be replaced when `mut` permits it. The picture has a limit: a name does not always mean a separate box containing a copy.

```rust
fn main() {
    let mut count: u32 = 2;
    println!("Before: {count}");
    count = 3;
    println!("After: {count}");
}
```

Output:

```text
Before: 2
After: 3
```

## Reading the declaration

let declares a variable. mut permits assigning a new value. count is our chosen name. : u32 is a **type annotation**, explicitly selecting a 32-bit unsigned integer. Unsigned means no negative values. A **bit** is a binary unit with two possible states; the number of bits limits the range. Lesson 10 gives the exact ranges. = assigns a value; it does not ask whether values are equal. 2 is a numeric literal and ; ends the statement.

Inside the println! text, {count} inserts the variable's value. These braces inside a string are not a Rust block. count = 3 changes the value without changing the name or type. Without mut, reassignment is forbidden.

## Inferred types

An annotation can sometimes be omitted: the compiler determines the type from the value and surrounding uses. This is **type inference**, not the absence of a type. With let count = 2 and no other constraints, the integer type defaults to i32, which permits negative numbers. We explicitly use u32 for this teaching counter rather than asking you to guess the context.

`true` and `false` are `bool` values. `'A'` is a `char`, one Unicode scalar value. Unicode is a shared system for representing characters; a visible character can contain several scalar values. We will return to this in the text lesson. Single quotes surround this `char`; double quotes surround text.

`"Rust"` is a string literal of type `&str`. **Owning data** means being responsible for its storage and eventual release. Here `&` gives a **reference**: reading access to existing text without taking that responsibility; `str` means a piece of text. Picture a card pointing to a page in a book: it makes no second book and is useful only while the book exists. This literal remains available throughout the program; we do not grow it through `&str`. Lessons 18–22 develop the ownership and string rules. Boolean operations arrive in lesson 11. These definitions show that types describe more than numbers.


Try several types in one small program. true is printed as true; Rust does not translate Boolean values. ready, mark and title are our chosen names. Braces inside the println! string insert a named value, just as {count} did above.

```rust
fn main() {
    let ready: bool = true;
    let mark: char = 'A';
    let title: &str = "Rust";
    println!("{ready}");
    println!("{mark}");
    println!("{title}");
}
```

```text
true
A
Rust
```

## A new declaration and scope

Another let with the same name creates a new binding called **shadowing**; it is not reassignment. **Scope** is the part of the program where a name is visible. A nested block may use its own binding, with the outer one visible again afterwards:

```rust
fn main() {
    let count = 2;
    {
        let count = 3;
        println!("Inside: {count}");
    }
    println!("Outside: {count}");
}
```

Output:

```text
Inside: 3
Outside: 2
```

const declares a **constant**, evaluated during compilation. Its type must be written explicitly; mut is not allowed. An immutable variable can receive a value during execution, so these are not the same concept. Constant names conventionally use uppercase letters and underscores. DAILY_LIMIT means a daily limit.

```rust
const DAILY_LIMIT: u32 = 5;

fn main() {
    println!("Plan: {DAILY_LIMIT}");
}
```

Output:

```text
Plan: 5
```

Deliberate error: we forgot permission to change the variable. Expect E0384 rather than output:

<!-- error-code: E0384 -->
```rust,compile_fail
fn main() {
    let count = 2;
    println!("Before: {count}");
    count = 3;
    println!("After: {count}");
}
```

## Recall map

Value → name → type → permitted actions. let creates a binding; mut allows reassignment; another let shadows; blocks delimit scope; const declares a constant.

## Warm-up

1. Predict the outer count after the nested block.
2. Complete let ___ count: u32 = 4; so it may later be assigned 5.
3. Repair E0384 without replacing the variable with a constant.

## Exercise

**Required.** Create a mutable u32 counter set to 4, print “Before: 4”, assign 5 and print “After: 5”. Explain why this is not shadowing.

**Your own data.** Choose two different small nonnegative numbers.

## Hint and reference answer

The reference has only one let followed by reassignment. Warm-up answers: the outer value stays 2; the missing word is mut; permission to reassign fixes the deliberate error. Repeat the two examples separately if the distinction is still unclear before moving to arithmetic.

<!-- task-answer -->
```rust
fn main() {
    let mut count: u32 = 4;
    println!("Before: {count}");
    count = 5;
    println!("After: {count}");
}
```

Expected output:

```text
Before: 4
After: 5
```

[Previous lesson](/read/rust-08-diagnostics?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-10-numbers?lang=en)
