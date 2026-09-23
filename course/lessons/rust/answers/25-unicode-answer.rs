fn main() {
    let title = String::from("Я");
    let combined = "e\u{301}";
    println!("Я: байтов {}, значений char {}", title.len(), title.chars().count());
    println!(
        "e с ударением: байтов {}, значений char {}",
        combined.len(),
        combined.chars().count()
    );
}
