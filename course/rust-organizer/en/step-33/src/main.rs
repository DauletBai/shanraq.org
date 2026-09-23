fn parse_number(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(number) => Ok(number),
        Err(_) => Err(String::from("Enter a whole number")),
    }
}

fn checked_minutes(text: &str) -> Result<u32, String> {
    let minutes = parse_number(text)?;
    if minutes == 0 {
        return Err(String::from("Minutes must be greater than zero"));
    }
    Ok(minutes)
}

fn main() {
    for text in ["20", "0", "soon"] {
        match checked_minutes(text) {
            Ok(minutes) => println!("Accepted: {minutes} min"),
            Err(message) => println!("Error: {message}"),
        }
    }
}
