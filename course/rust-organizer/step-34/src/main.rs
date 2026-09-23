fn main() {
    let mut line = String::new();
    println!("Введите название:");
    match std::io::stdin().read_line(&mut line) {
        Ok(0) => println!("Ввод закончился"),
        Ok(_) => {
            let title = line.trim();
            if title == "" {
                println!("Пустое название");
            } else {
                println!("Задача: {title}");
            }
        }
        Err(_) => println!("Не удалось прочитать ввод"),
    }
}
