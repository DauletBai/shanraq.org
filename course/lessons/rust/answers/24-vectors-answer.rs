fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Чтение"));
    titles.push(String::from("Прогулка"));
    titles.push(String::from("Отдых"));
    let removed = titles.remove(1);
    println!("Удалено: {removed}");
    for title in &titles {
        println!("Задача: {title}");
    }
    println!("Осталось: {}", titles.len());
}
