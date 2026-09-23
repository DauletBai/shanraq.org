fn decorate(mut title: String) -> String {
    title.push_str("!");
    title
}

fn main() {
    let original = String::from("Жоспар");
    let decorated = decorate(original);
    println!("{decorated}");
}
