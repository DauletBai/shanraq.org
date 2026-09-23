enum Command {
    Add(String),
    List,
}

fn parse(action: &str, title: Option<&str>, extra: Option<&str>) -> Result<Command, String> {
    match (action, title, extra) {
        ("list", None, None) => Ok(Command::List),
        ("add", Some(title), None) => {
            if title.trim().is_empty() {
                Err(String::from("Бос атау"))
            } else {
                Ok(Command::Add(String::from(title)))
            }
        }
        _ => Err(String::from("Пәрмен не аргумент саны қате")),
    }
}

fn show(result: Result<Command, String>) {
    match result {
        Ok(Command::Add(title)) => println!("Қосу: {title}"),
        Ok(Command::List) => println!("Тізімді көрсету"),
        Err(message) => println!("Қате: {message}"),
    }
}

fn main() {
    show(parse("add", Some("Rust оқу"), None));
    show(parse("list", None, None));
    show(parse("add", Some("   "), None));
}
