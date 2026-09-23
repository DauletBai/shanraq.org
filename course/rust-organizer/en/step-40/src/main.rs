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
