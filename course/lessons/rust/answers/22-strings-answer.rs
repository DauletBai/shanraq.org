fn show_task(title: &str, minutes: u32) {
    println!("{title}: {minutes} мин");
}

fn main() {
    let owned = String::from("Чтение");
    let literal: &str = "Отдых";
    show_task(owned.as_str(), 15);
    show_task(literal, 5);
    println!("Сохранено: {owned}");
}
