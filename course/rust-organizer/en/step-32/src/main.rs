fn read_minutes(text: &str) -> Result<u32, String> {
    match text.parse::<u32>() {
        Ok(minutes) => Ok(minutes),
        Err(_) => Err(String::from("Enter whole minutes")),
    }
}

fn main() {
    for text in ["25", "soon"] {
        match read_minutes(text) {
            Ok(minutes) => println!("Minutes: {minutes}"),
            Err(message) => println!("Error: {message}"),
        }
    }
}
