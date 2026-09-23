enum Command {
    Add(String),
    List,
}

fn parse(action: &str, title: Option<&str>, extra: Option<&str>) -> Result<Command, String> {
    match (action, title, extra) {
        ("list", None, None) => Ok(Command::List),
        ("add", Some(title), None) => {
            if title.trim().is_empty() {
                Err(String::from("Пустое название"))
            } else {
                Ok(Command::Add(String::from(title)))
            }
        }
        _ => Err(String::from("Неверная команда или число аргументов")),
    }
}

fn show(result: Result<Command, String>) {
    match result {
        Ok(Command::Add(title)) => println!("Добавить: {title}"),
        Ok(Command::List) => println!("Показать список"),
        Err(message) => println!("Ошибка: {message}"),
    }
}

fn main() {
    show(parse("add", Some("Читать Rust"), None));
    show(parse("list", None, None));
    show(parse("add", Some("   "), None));
}
