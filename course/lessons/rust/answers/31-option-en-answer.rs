fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Reading"));
    titles.push(String::from("Walk"));
    for index in [1, 2] {
        match titles.get(index) {
            Some(title) => println!("Found: {title}"),
            None => println!("Position {index} is absent"),
        }
    }
}
