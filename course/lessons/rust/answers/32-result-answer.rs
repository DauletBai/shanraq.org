fn read_minutes(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(minutes) => Ok(minutes),
        Err(_) => Err(String::from("Введите число минут")),
    }
}

fn main() {
    for text in ["30", "много"] {
        match read_minutes(text) {
            Ok(minutes) => println!("Минут: {minutes}"),
            Err(message) => println!("Ошибка: {message}"),
        }
    }
}
