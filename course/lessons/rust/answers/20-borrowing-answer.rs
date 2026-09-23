fn show_task(title: &String, minutes: u32) {
    println!("{title} — {minutes} мин");
}

fn main() {
    let title = String::from("Чтение");
    show_task(&title, 15);
    show_task(&title, 20);
    println!("Название осталось: {title}");
}
