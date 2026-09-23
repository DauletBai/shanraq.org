# Numbers and arithmetic

_Summary:_ **Check calculations and the limits of numeric types.**

Prerequisite: lessons 1–9. Work in organizer unless stated otherwise.

## Why this matters

Our organizer now has numbers. **Integers** have no fractional part. **Floating-point numbers** approximately represent a wide range of numbers. A **bit** is a binary unit; a **byte** consists of eight bits. The available bits limit the values a type can represent, like the capacity of a measuring cup.

## The whole picture

Replace src/main.rs and run cargo run from organizer:

Before running, read the new signs: `-` subtracts, `/` divides, and `%` gives the remainder of integer division. Empty `{}` in `println!` is replaced with the following value after the comma; for example, `println!("{}", 5)` prints 5. Picture grouping five tasks into pairs: two full pairs and one task left over. Here `%` means remainder, not a percentage.

```rust
fn main() {
    let total: u32 = 7;
    let done: u32 = 2;
    let remaining = total - done;
    println!("Remaining: {remaining}");
    println!("Complete pairs: {}", remaining / 2);
    println!("Unpaired: {}", remaining % 2);
}
```

Output:

```text
Remaining: 5
Complete pairs: 2
Unpaired: 1
```

`-` subtracts, `/` divides, and `%` gives the remainder. Empty {} in println! inserts the following argument after the comma. This lets us print a calculation without naming an extra variable.

Seven tasks minus two completed leaves five. Integer division 5 / 2 gives 2, with remainder 5 % 2 equal to 1. Integer division truncates towards zero. Integer division by zero is invalid.

## Choosing a type

u8 stores values from 0 through 255; i8 stores −128 through 127. u32 stores 0 through 4,294,967,295; i32 stores −2,147,483,648 through 2,147,483,647. The u prefix means unsigned, i means signed, and the number is the bit count. There are also 16-, 64- and 128-bit integers. usize depends on the platform and is used for sizes and indices, which we will meet with collections. The type has a limit even when today's value is small.

f32 and f64 are 32- and 64-bit floating-point types; f64 generally offers greater precision and range. An unconstrained floating-point literal defaults to f64. Not every decimal fraction has an exact binary representation. Task counts need integers; exact monetary accounting should not casually use floating-point numbers. Money is outside this project.

`+`, `-`, `*`, `/` and `%` mean addition, subtraction, multiplication, division and remainder. Multiplication and division precede addition; parentheses specify the intended order explicitly.


Compare floating-point division and operation order. We specify f64 explicitly; a decimal point is used in code. These particular numbers have exact binary representations. average is our chosen name for a division result, not a task count. Predict all three lines first.

```rust
fn main() {
    let average: f64 = 5.0 / 2.0;
    println!("{average}");
    println!("{}", 2 + 3 * 4);
    println!("{}", (2 + 3) * 4);
}
```

```text
2.5
14
20
```

## Overflow and conversions

**Overflow** occurs when a result cannot fit its type. In our fixed example, done never exceeds total. If that condition is broken, unsigned subtraction does not produce an ordinary negative value. With overflow checks enabled, as they normally are in a debug build, the program stops with an error when executed. If checks are disabled, the value may wrap around the range like a counter after its last digit; that is not a correct task count. Conditions and input-validation lessons will enforce this rule explicitly.

Applying a minus sign to an unsigned u32 is invalid. This deliberately wrong example should produce E0600:

<!-- error-code: E0600 -->
```rust,compile_fail
fn main() {
    let count: u32 = -1;
    println!("{count}");
}
```

A type conversion is an operation, not merely changing how you spell a number. Later we will study checked conversions that can fail. `as` is an explicit type-conversion operator. Do not treat it as a universal repair: narrowing may lose data. Here our types agree and our fixed inputs satisfy the task's conditions.

## Recall map

Meaning → type → range → operation → expected answer → check. Compiling does not prove that total and done make sense.

## Warm-up

1. Predict 6 / 2 and 6 % 2 before running.
2. Complete remaining = total ___ done.
3. Before changing the type to accept a negative task count, decide whether that value makes sense.

## Exercise

**Required.** With total = 10 and done = 4, produce “Remaining: 6”, “Complete pairs: 3” and “Unpaired: 0”. Calculate the results rather than printing fixed answers.

**Your own data.** Use total = done. All three results should be zero.

## Hint and reference answer

Only the inputs change; the formulas remain the same. Warm-up answers: 3 and 0; subtraction; negative counts are invalid for our task. The zero case helps detect a mistaken addition.

<!-- task-answer -->
```rust
fn main() {
    let total: u32 = 10;
    let done: u32 = 4;
    let remaining = total - done;
    println!("Remaining: {remaining}");
    println!("Complete pairs: {}", remaining / 2);
    println!("Unpaired: {}", remaining % 2);
}
```

Expected output:

```text
Remaining: 6
Complete pairs: 3
Unpaired: 0
```

Continue when you can declare a variable, explain its type, change its value and calculate remaining tasks independently. Lessons 11–15 cover logic, branches, loops, functions and arrays. Continue with lesson 11.

[Previous lesson](/read/rust-09-variables?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-11-boolean?lang=en)
