fn parse_number(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(number) => Ok(number),
        Err(_) => Err(String::from("Нужно целое число")),
    }
}

fn checked_minutes(text: &str) -> Result<u32, String> {
    let minutes = parse_number(text)?;
    if minutes == 0 {
        return Err(String::from("Минуты должны быть больше нуля"));
    }
    Ok(minutes)
}

fn main() {
    for text in ["20", "0", "скоро"] {
        match checked_minutes(text) {
            Ok(minutes) => println!("Принято: {minutes} мин"),
            Err(message) => println!("Ошибка: {message}"),
        }
    }
}
