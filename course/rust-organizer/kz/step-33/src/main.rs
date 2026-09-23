fn parse_number(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(number) => Ok(number),
        Err(_) => Err(String::from("Бүтін сан енгізіңіз")),
    }
}

fn checked_minutes(text: &str) -> Result<u32, String> {
    let minutes = parse_number(text)?;
    if minutes == 0 {
        return Err(String::from("Минут нөлден үлкен болуы керек"));
    }
    Ok(minutes)
}

fn main() {
    for text in ["20", "0", "жақында"] {
        match checked_minutes(text) {
            Ok(minutes) => println!("Қабылданды: {minutes} мин"),
            Err(message) => println!("Қате: {message}"),
        }
    }
}
