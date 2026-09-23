fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Reading"));
    for index in [0, 2] {
        match titles.get(index) {
            Some(title) => println!("Found: {title}"),
            None => println!("No record at position {index}"),
        }
    }
}
