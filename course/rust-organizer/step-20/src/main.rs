fn show_title(title: &String) {
    println!("В функции: {title}");
}

fn main() {
    let title = String::from("Чтение");
    show_title(&title);
    println!("После вызова: {title}");
}
