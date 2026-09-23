fn main() {
    let title = String::from("é");
    println!("Bytes: {}", title.len());
    println!("char values: {}", title.chars().count());
}
