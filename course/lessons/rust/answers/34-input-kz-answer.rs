fn main() {
    let mut line = String::new();
    match std::io::stdin().read_line(&mut line) {
        Ok(0) => println!("Басқа енгізу жоқ"),
        Ok(_) => {
            let title = line.trim();
            if title == "" {
                println!("Атау қажет");
            } else {
                println!("Қосу: {title}");
            }
        }
        Err(_) => println!("Оқу қатесі"),
    }
}
