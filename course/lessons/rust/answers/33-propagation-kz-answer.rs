fn parse_number(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(number) => Ok(number),
        Err(_) => Err(String::from("Сан енгізіңіз")),
    }
}

fn checked_minutes(text: &str) -> Result<u32, String> {
    let minutes = parse_number(text)?;
    if minutes == 0 {
        return Err(String::from("Минут нөл болмауы керек"));
    }
    Ok(minutes)
}

fn main() {
    for text in ["15", "0", "көп"] {
        match checked_minutes(text) {
            Ok(minutes) => println!("Қабылданды: {minutes} мин"),
            Err(message) => println!("Қате: {message}"),
        }
    }
}
