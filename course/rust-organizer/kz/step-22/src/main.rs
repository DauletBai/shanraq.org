fn show_title(title: &str) {
    println!("Атауы: {title}");
}

fn main() {
    let owned = String::from("Оқу");
    let literal: &str = "Демалыс";
    show_title(owned.as_str());
    show_title(literal);
    println!("Сақталды: {owned}");
}
