fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Оқу"));
    titles.push(String::from("Серуен"));
    titles.push(String::from("Демалыс"));
    let removed = titles.remove(1);
    println!("Өшірілді: {removed}");
    for title in &titles {
        println!("Тапсырма: {title}");
    }
    println!("Қалды: {}", titles.len());
}
