fn show_title(title: &str) {
    println!("Title: {title}");
}

fn main() {
    let owned = String::from("Reading");
    let literal: &str = "Rest";
    show_title(owned.as_str());
    show_title(literal);
    println!("Stored: {owned}");
}
