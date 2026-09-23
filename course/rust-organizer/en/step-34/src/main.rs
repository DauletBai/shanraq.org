fn main() {
    let mut line = String::new();
    println!("Enter a title:");
    match std::io::stdin().read_line(&mut line) {
        Ok(0) => println!("Input ended"),
        Ok(_) => {
            let title = line.trim();
            if title == "" {
                println!("Empty title");
            } else {
                println!("Task: {title}");
            }
        }
        Err(_) => println!("Could not read input"),
    }
}
