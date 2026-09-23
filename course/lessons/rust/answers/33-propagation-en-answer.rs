fn parse_number(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(number) => Ok(number),
        Err(_) => Err(String::from("Enter a number")),
    }
}

fn checked_minutes(text: &str) -> Result<u32, String> {
    let minutes = parse_number(text)?;
    if minutes == 0 {
        return Err(String::from("Minutes cannot be zero"));
    }
    Ok(minutes)
}

fn main() {
    for text in ["15", "0", "many"] {
        match checked_minutes(text) {
            Ok(minutes) => println!("Accepted: {minutes} min"),
            Err(message) => println!("Error: {message}"),
        }
    }
}
