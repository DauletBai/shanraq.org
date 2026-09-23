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
                Err(String::from("Пустое название"))
            } else {
                Ok(Command::Add(String::from(title)))
            }
        }
        ("done", Some(text), None) => match text.parse::<u32>() {
            Ok(0) => Err(String::from("Номер должен быть больше нуля")),
            Ok(number) => Ok(Command::Done(number)),
            Err(_) => Err(String::from("Нужен номер")),
        },
        _ => Err(String::from("Неверная команда или число аргументов")),
    }
}

fn show(result: Result<Command, String>) {
    match result {
        Ok(Command::Add(title)) => println!("Добавить: {title}"),
        Ok(Command::List) => println!("Показать список"),
        Ok(Command::Done(number)) => println!("Завершить: {number}"),
        Err(message) => println!("Ошибка: {message}"),
    }
}

fn main() {
    show(parse("add", Some("Читать Rust"), None));
    show(parse("list", None, None));
    show(parse("add", Some("   "), None));
    show(parse("done", Some("2"), None));
    show(parse("done", Some("0"), None));
    show(parse("done", Some("нет"), None));
}
