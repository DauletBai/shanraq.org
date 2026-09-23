fn main() {
    let mut line = String::new();
    println!("Атау енгізіңіз:");
    match std::io::stdin().read_line(&mut line) {
        Ok(0) => println!("Енгізу аяқталды"),
        Ok(_) => {
            let title = line.trim();
            if title == "" {
                println!("Атау бос");
            } else {
                println!("Тапсырма: {title}");
            }
        }
        Err(_) => println!("Енгізуді оқу мүмкін болмады"),
    }
}
