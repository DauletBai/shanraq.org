fn print_title(title: String) {
    println!("Задача: {title}");
}

fn main() {
    let title = String::from("Чтение");
    print_title(title);
}
