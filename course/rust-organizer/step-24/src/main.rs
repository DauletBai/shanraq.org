fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Чтение"));
    titles.push(String::from("Отдых"));
    println!("Задач: {}", titles.len());
    for title in &titles {
        println!("Задача: {title}");
    }
}
