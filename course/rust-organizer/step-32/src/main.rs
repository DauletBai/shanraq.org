fn read_minutes(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(minutes) => Ok(minutes),
        Err(_) => Err(String::from("Нужно целое число минут")),
    }
}

fn main() {
    for text in ["25", "скоро"] {
        match read_minutes(text) {
            Ok(minutes) => println!("Минут: {minutes}"),
            Err(message) => println!("Ошибка: {message}"),
        }
    }
}
