fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Чтение"));
    for index in [0, 2] {
        match titles.get(index) {
            Some(title) => println!("Найдено: {title}"),
            None => println!("Нет записи на позиции {index}"),
        }
    }
}
