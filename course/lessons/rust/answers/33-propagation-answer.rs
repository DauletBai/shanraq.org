fn parse_number(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(number) => Ok(number),
        Err(_) => Err(String::from("Введите число")),
    }
}

fn checked_minutes(text: &str) -> Result<u32, String> {
    let minutes = parse_number(text)?;
    if minutes == 0 {
        return Err(String::from("Минуты не могут быть нулём"));
    }
    Ok(minutes)
}

fn main() {
    for text in ["15", "0", "много"] {
        match checked_minutes(text) {
            Ok(minutes) => println!("Принято: {minutes} мин"),
            Err(message) => println!("Ошибка: {message}"),
        }
    }
}
