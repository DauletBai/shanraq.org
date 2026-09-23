enum Command {
    Add(String),
    List,
    Done(u32),
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
        ("done", Some(text), None) => match text.parse::<u32>() {
            Ok(0) => Err(String::from("Number must be positive")),
            Ok(number) => Ok(Command::Done(number)),
            Err(_) => Err(String::from("A number is required")),
        },
        _ => Err(String::from("Unknown command or argument count")),
    }
}

fn show(result: Result<Command, String>) {
    match result {
        Ok(Command::Add(title)) => println!("Add: {title}"),
        Ok(Command::List) => println!("Show list"),
        Ok(Command::Done(number)) => println!("Complete: {number}"),
        Err(message) => println!("Error: {message}"),
    }
}

fn main() {
    show(parse("add", Some("Read Rust"), None));
    show(parse("list", None, None));
    show(parse("add", Some("   "), None));
    show(parse("done", Some("2"), None));
    show(parse("done", Some("0"), None));
    show(parse("done", Some("no"), None));
}
