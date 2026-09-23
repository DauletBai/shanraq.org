fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Оқу"));
    titles.push(String::from("Серуен"));
    for index in [1, 2] {
        match titles.get(index) {
            Some(title) => println!("Табылды: {title}"),
            None => println!("{index}-орын жоқ"),
        }
    }
}
