fn main() {
    let mut args = std::env::args();
    let _program = args.next();
    let command = args.next();
    let title = args.next();
    let extra = args.next();
    match (command, title, extra) {
        (Some(command), Some(title), None) => {
            if command == "add" {
                println!("Запланировано: {title}");
            } else {
                println!("Команда не найдена");
            }
        }
        _ => println!("Нужно: add \"Название\""),
    }
}
