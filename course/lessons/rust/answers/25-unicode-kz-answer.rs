fn main() {
    let title = String::from("Ә");
    let combined = "e\u{301}";
    println!("Ә: байт саны {}, char саны {}", title.len(), title.chars().count());
    println!(
        "екпінді e: байт саны {}, char саны {}",
        combined.len(),
        combined.chars().count()
    );
}
