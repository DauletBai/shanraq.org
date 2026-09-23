fn main() {
    let mut line = String::new();
    match std::io::stdin().read_line(&mut line) {
        Ok(0) => println!("No more input"),
        Ok(_) => {
            let title = line.trim();
            if title == "" {
                println!("A title is needed");
            } else {
                println!("Add: {title}");
            }
        }
        Err(_) => println!("Read error"),
    }
}
