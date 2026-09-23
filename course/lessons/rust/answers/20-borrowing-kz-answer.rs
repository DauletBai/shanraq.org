fn show_task(title: &String, minutes: u32) {
    println!("{title} — {minutes} мин");
}

fn main() {
    let title = String::from("Оқу");
    show_task(&title, 15);
    show_task(&title, 20);
    println!("Атау сақталды: {title}");
}
