fn read_minutes(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(minutes) => Ok(minutes),
        Err(_) => Err(String::from("Минут санын енгізіңіз")),
    }
}

fn main() {
    for text in ["25", "жақында"] {
        match read_minutes(text) {
            Ok(minutes) => println!("Минут: {minutes}"),
            Err(message) => println!("Қате: {message}"),
        }
    }
}
