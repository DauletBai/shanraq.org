fn main() {
    let mut line = String::new();
    match std::io::stdin().read_line(&mut line) {
        Ok(0) => println!("Больше ввода нет"),
        Ok(_) => {
            let title = line.trim();
            if title == "" {
                println!("Нужно название");
            } else {
                println!("Добавить: {title}");
            }
        }
        Err(_) => println!("Ошибка чтения"),
    }
}
