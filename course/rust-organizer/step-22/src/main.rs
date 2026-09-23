fn show_title(title: &str) {
    println!("Название: {title}");
}

fn main() {
    let owned = String::from("Чтение");
    let literal: &str = "Отдых";
    show_title(owned.as_str());
    show_title(literal);
    println!("Хранится: {owned}");
}
