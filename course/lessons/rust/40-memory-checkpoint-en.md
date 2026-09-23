# Checkpoint: an organizer in memory

_Lead (summary):_ **Lesson 40. Run an interactive organizer, perform several commands in one session, and handle invalid input.**

Prerequisites: lessons 16, 21, 27–28 and 31–39. Revisit [lesson 34](/read/rust-34-input?lang=en) for end of input and [lesson 38](/read/rust-38-ids?lang=en) for IDs.

## Familiar image and recall map

At a reception desk, a visitor can submit several requests in a row; the clerk remembers them until closing. This organizer similarly accepts lines until a quit command or end of input, keeping tasks in one run of memory. The image has a limit: tasks still disappear when the program closes; file saving comes later.

**Input line → parse `Command` → act on `Organizer` → response → next line; `quit` or end of input → stop.** This combines `read_line` from lesson 34, command enums from 36, task actions from 37, stable IDs from 38, and error checks from 39. The standard-library method `split_once(' ')` splits text **at the first space** and returns `Some((left part, right part))`, or `None` when there is no space. For example, `add Read Rust` yields action `add` and the whole title `Read Rust`. Quotes are unnecessary here because the entire line has already been read. If typed, quotes become part of the title; the separator is an ordinary space, not a tab. `trim()` removes edge spaces and the line ending. `self.tasks.is_empty()` checks whether the task list has no entries. `done` and `delete` require a positive ID; `list` and `quit` take no data. Every error prints a reason and then waits for another line. `use std::io` shortens the path to the input module: we write `io::stdin()` instead of `std::io::stdin()`. It is the same input object from lesson 34; lesson 44 explains modules in depth. `!execute(...)` means “do not continue” because `execute` returns `false` for `quit`.

## Several commands in one run

Save and completely replace the old `src/main.rs`, then run `cargo run` beside `Cargo.toml`. Type one command per line and press Enter. Type `quit` to leave; closing the input stream also ends the program. No new files or dependencies are needed. The output below is for an input stream closed immediately, as in the automated check.

```rust
use std::io;

enum Command {
    Add(String),
    List,
    Done(u32),
    Delete(u32),
    Quit,
}

fn parse_id(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(0) => Err(String::from("A positive ID is required")),
        Ok(id) => Ok(id),
        Err(_) => Err(String::from("A positive ID is required")),
    }
}

fn parse(line: &str) -> Result<Command, String> {
    let line = line.trim();
    if line == "list" {
        return Ok(Command::List);
    }
    if line == "quit" {
        return Ok(Command::Quit);
    }
    match line.split_once(' ') {
        Some(("add", title)) => {
            if title.trim().is_empty() {
                Err(String::from("Empty title"))
            } else {
                Ok(Command::Add(String::from(title.trim())))
            }
        }
        Some(("done", id)) => Ok(Command::Done(parse_id(id)?)),
        Some(("delete", id)) => Ok(Command::Delete(parse_id(id)?)),
        _ => Err(String::from("Unknown command or wrong arguments")),
    }
}

struct Task {
    id: u32,
    title: String,
    done: bool,
}
struct Organizer {
    tasks: Vec<Task>,
    next_id: u32,
}
impl Organizer {
    fn new() -> Self {
        Self {
            tasks: Vec::new(),
            next_id: 1,
        }
    }
    fn add(&mut self, title: String) -> Result<u32, String> {
        if title.trim().is_empty() {
            return Err(String::from("Empty title"));
        }
        match self.next_id.checked_add(1) {
            Some(next) => {
                let id = self.next_id;
                self.tasks.push(Task {
                    id,
                    title,
                    done: false,
                });
                self.next_id = next;
                Ok(id)
            }
            None => Err(String::from("No IDs left")),
        }
    }
    fn done(&mut self, id: u32) -> Result<(), String> {
        for task in &mut self.tasks {
            if task.id == id {
                task.done = true;
                return Ok(());
            }
        }
        Err(String::from("ID not found"))
    }
    fn delete(&mut self, id: u32) -> Result<(), String> {
        for index in 0..self.tasks.len() {
            if self.tasks[index].id == id {
                self.tasks.remove(index);
                return Ok(());
            }
        }
        Err(String::from("ID not found"))
    }
    fn list(&self) {
        if self.tasks.is_empty() {
            println!("No tasks");
        }
        for task in &self.tasks {
            let mark = if task.done { "done" } else { "open" };
            println!("{}: {} ({mark})", task.id, task.title);
        }
    }
}

fn execute(organizer: &mut Organizer, command: Command) -> bool {
    match command {
        Command::Add(title) => {
            match organizer.add(title) {
                Ok(id) => println!("Added: {id}"),
                Err(message) => println!("Error: {message}"),
            }
            true
        }
        Command::List => {
            organizer.list();
            true
        }
        Command::Done(id) => {
            match organizer.done(id) {
                Ok(()) => println!("Completed: {id}"),
                Err(message) => println!("Error: {message}"),
            }
            true
        }
        Command::Delete(id) => {
            match organizer.delete(id) {
                Ok(()) => println!("Deleted: {id}"),
                Err(message) => println!("Error: {message}"),
            }
            true
        }
        Command::Quit => false,
    }
}

fn main() {
    let mut organizer = Organizer::new();
    let stdin = io::stdin();
    println!("Commands: add TITLE | list | done ID | delete ID | quit");
    loop {
        let mut line = String::new();
        match stdin.read_line(&mut line) {
            Ok(0) => break,
            Ok(_) => match parse(&line) {
                Ok(command) => {
                    if !execute(&mut organizer, command) {
                        break;
                    }
                }
                Err(message) => println!("Error: {message}"),
            },
            Err(_) => {
                println!("Error: could not read input");
                break;
            }
        }
    }
    println!("Session ended");
}
```
```text
Commands: add TITLE | list | done ID | delete ID | quit
Session ended
```

Type `add Read Rust`, `add Review tests`, `list`, `done 1`, `delete 2`, `list`, `done 99`, `quit`. ID 1 remains marked done and unknown ID produces an error. Another run starts with an empty list: that is the expected boundary before file storage.

The listed commands produce this output (the typed lines are not repeated here):

```text
Commands: add TITLE | list | done ID | delete ID | quit
Added: 1
Added: 2
1: Read Rust (open)
2: Review tests (open)
Completed: 1
Deleted: 2
1: Read Rust (done)
Error: ID not found
Session ended
```

## Check your understanding

1. Why does `add Read Rust` keep the space in the title?
2. Why must `done 99` leave the list unchanged?
3. Will tasks survive `quit` and a new run?

`split_once` separates at only the first space; an unknown ID changes no record; without a file, data disappear. Rebuild the recall map without looking.

## Exercise

**Required.** Add `rename ID TITLE`: find the task by stable ID, reject an empty title, and keep the list unchanged for an unknown ID. Try `add Read Rust`, `rename 1 Read two chapters`, `list`, `rename 99 Error`, and `quit`. Run the program again and confirm the list is empty. Do not treat a `Vec` position as an ID.

## Hint

Add `Rename(u32, String)` to `Command`. Call `split_once(' ')` again on the text after `rename `, then check ID and title separately. In `rename`, walk `&mut self.tasks` and change the field only after validation.

## Reference after attempting

<!-- task-answer -->
```rust
use std::io;

enum Command {
    Add(String),
    List,
    Done(u32),
    Delete(u32),
    Rename(u32, String),
    Quit,
}

fn parse_id(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(0) => Err(String::from("A positive ID is required")),
        Ok(id) => Ok(id),
        Err(_) => Err(String::from("A positive ID is required")),
    }
}

fn parse(line: &str) -> Result<Command, String> {
    let line = line.trim();
    if line == "list" {
        return Ok(Command::List);
    }
    if line == "quit" {
        return Ok(Command::Quit);
    }
    match line.split_once(' ') {
        Some(("add", title)) => {
            if title.trim().is_empty() {
                Err(String::from("Empty title"))
            } else {
                Ok(Command::Add(String::from(title.trim())))
            }
        }
        Some(("done", id)) => Ok(Command::Done(parse_id(id)?)),
        Some(("delete", id)) => Ok(Command::Delete(parse_id(id)?)),
        Some(("rename", rest)) => match rest.split_once(' ') {
            Some((id, title)) => {
                let id = parse_id(id)?;
                if title.trim().is_empty() {
                    Err(String::from("Empty title"))
                } else {
                    Ok(Command::Rename(id, String::from(title.trim())))
                }
            }
            None => Err(String::from("Unknown command or wrong arguments")),
        },
        _ => Err(String::from("Unknown command or wrong arguments")),
    }
}

struct Task {
    id: u32,
    title: String,
    done: bool,
}
struct Organizer {
    tasks: Vec<Task>,
    next_id: u32,
}
impl Organizer {
    fn new() -> Self {
        Self {
            tasks: Vec::new(),
            next_id: 1,
        }
    }
    fn add(&mut self, title: String) -> Result<u32, String> {
        if title.trim().is_empty() {
            return Err(String::from("Empty title"));
        }
        match self.next_id.checked_add(1) {
            Some(next) => {
                let id = self.next_id;
                self.tasks.push(Task {
                    id,
                    title,
                    done: false,
                });
                self.next_id = next;
                Ok(id)
            }
            None => Err(String::from("No IDs left")),
        }
    }
    fn done(&mut self, id: u32) -> Result<(), String> {
        for task in &mut self.tasks {
            if task.id == id {
                task.done = true;
                return Ok(());
            }
        }
        Err(String::from("ID not found"))
    }
    fn delete(&mut self, id: u32) -> Result<(), String> {
        for index in 0..self.tasks.len() {
            if self.tasks[index].id == id {
                self.tasks.remove(index);
                return Ok(());
            }
        }
        Err(String::from("ID not found"))
    }
    fn rename(&mut self, id: u32, title: String) -> Result<(), String> {
        if title.trim().is_empty() {
            return Err(String::from("Empty title"));
        }
        for task in &mut self.tasks {
            if task.id == id {
                task.title = title;
                return Ok(());
            }
        }
        Err(String::from("ID not found"))
    }
    fn list(&self) {
        if self.tasks.is_empty() {
            println!("No tasks");
        }
        for task in &self.tasks {
            let mark = if task.done { "done" } else { "open" };
            println!("{}: {} ({mark})", task.id, task.title);
        }
    }
}

fn execute(organizer: &mut Organizer, command: Command) -> bool {
    match command {
        Command::Add(title) => {
            match organizer.add(title) {
                Ok(id) => println!("Added: {id}"),
                Err(message) => println!("Error: {message}"),
            }
            true
        }
        Command::List => {
            organizer.list();
            true
        }
        Command::Done(id) => {
            match organizer.done(id) {
                Ok(()) => println!("Completed: {id}"),
                Err(message) => println!("Error: {message}"),
            }
            true
        }
        Command::Delete(id) => {
            match organizer.delete(id) {
                Ok(()) => println!("Deleted: {id}"),
                Err(message) => println!("Error: {message}"),
            }
            true
        }
        Command::Rename(id, title) => {
            match organizer.rename(id, title) {
                Ok(()) => println!("Renamed: {id}"),
                Err(message) => println!("Error: {message}"),
            }
            true
        }
        Command::Quit => false,
    }
}

fn main() {
    let mut organizer = Organizer::new();
    let stdin = io::stdin();
    println!("Commands: add TITLE | list | done ID | delete ID | rename ID TITLE | quit");
    loop {
        let mut line = String::new();
        match stdin.read_line(&mut line) {
            Ok(0) => break,
            Ok(_) => match parse(&line) {
                Ok(command) => {
                    if !execute(&mut organizer, command) {
                        break;
                    }
                }
                Err(message) => println!("Error: {message}"),
            },
            Err(_) => {
                println!("Error: could not read input");
                break;
            }
        }
    }
    println!("Session ended");
}
```
```text
Commands: add TITLE | list | done ID | delete ID | rename ID TITLE | quit
Session ended
```

The reference prints this output for the exercise commands:

```text
Commands: add TITLE | list | done ID | delete ID | rename ID TITLE | quit
Added: 1
Renamed: 1
1: Read two chapters (open)
Error: ID not found
Session ended
```

This output uses an empty input stream. Typing the exercise commands adds organizer responses between the first and last lines. [Official explanations of input reading and `split_once`](https://doc.rust-lang.org/std/io/struct.Stdin.html#method.read_line) · [`split_once`](https://doc.rust-lang.org/std/primitive.str.html#method.split_once).

## Checkpoint after lesson 40

Close this page and explain one command without copying: input line → `parse` → `Command` → action → response. Build the program, add two tasks, complete one, delete the other, and try an unknown ID. Close the program and explain why the list disappears. Revisit [lesson 36](/read/rust-36-commands?lang=en) for parsing, [lesson 37](/read/rust-37-operations?lang=en) for list changes, [lesson 38](/read/rust-38-ids?lang=en) for IDs, and [lesson 39](/read/rust-39-scenarios?lang=en) for error checks.

[Previous lesson](/read/rust-39-scenarios?lang=en) · [Contents](/course/rust?lang=en)
