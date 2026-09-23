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
                Err(String::from("Бос атау"))
            } else {
                Ok(Command::Add(String::from(title)))
            }
        }
        ("done", Some(text), None) => match text.parse::<u32>() {
            Ok(0) => Err(String::from("Нөмір нөлден үлкен болуы керек")),
            Ok(number) => Ok(Command::Done(number)),
            Err(_) => Err(String::from("Нөмір қажет")),
        },
        _ => Err(String::from("Пәрмен не аргумент саны қате")),
    }
}

fn show(result: Result<Command, String>) {
    match result {
        Ok(Command::Add(title)) => println!("Қосу: {title}"),
        Ok(Command::List) => println!("Тізімді көрсету"),
        Ok(Command::Done(number)) => println!("Аяқтау: {number}"),
        Err(message) => println!("Қате: {message}"),
    }
}

fn main() {
    show(parse("add", Some("Rust оқу"), None));
    show(parse("list", None, None));
    show(parse("add", Some("   "), None));
    show(parse("done", Some("2"), None));
    show(parse("done", Some("0"), None));
    show(parse("done", Some("жоқ"), None));
}
