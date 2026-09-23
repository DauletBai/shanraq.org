fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Оқу"));
    for index in [0, 2] {
        match titles.get(index) {
            Some(title) => println!("Табылды: {title}"),
            None => println!("{index}-орында жазба жоқ"),
        }
    }
}
