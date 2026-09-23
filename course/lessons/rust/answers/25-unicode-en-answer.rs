fn main() {
    let title = String::from("é");
    let combined = "e\u{301}";
    println!("é: bytes {}, char values {}", title.len(), title.chars().count());
    println!(
        "accented e: bytes {}, char values {}",
        combined.len(),
        combined.chars().count()
    );
}
