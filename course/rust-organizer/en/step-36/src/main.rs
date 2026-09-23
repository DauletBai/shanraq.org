enum Command {
    Add(String),
    List,
}

fn parse(action: &str, title: Option<&str>, extra: Option<&str>) -> Result<Command, String> {
    match (action, title, extra) {
        ("list", None, None) => Ok(Command::List),
        ("add", Some(title), None) => {
            if title.trim().is_empty() {
                Err(String::from("Empty title"))
            } else {
                Ok(Command::Add(String::from(title)))
            }
        }
        _ => Err(String::from("Unknown command or argument count")),
    }
}

fn show(result: Result<Command, String>) {
    match result {
        Ok(Command::Add(title)) => println!("Add: {title}"),
        Ok(Command::List) => println!("Show list"),
        Err(message) => println!("Error: {message}"),
    }
}

fn main() {
    show(parse("add", Some("Read Rust"), None));
    show(parse("list", None, None));
    show(parse("add", Some("   "), None));
}
