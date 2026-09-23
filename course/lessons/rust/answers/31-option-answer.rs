fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Чтение"));
    titles.push(String::from("Прогулка"));
    for index in [1, 2] {
        match titles.get(index) {
            Some(title) => println!("Найдено: {title}"),
            None => println!("Позиция {index} отсутствует"),
        }
    }
}
