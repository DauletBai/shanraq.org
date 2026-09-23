fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Оқу"));
    titles.push(String::from("Демалыс"));
    println!("Тапсырма саны: {}", titles.len());
    for title in &titles {
        println!("Тапсырма: {title}");
    }
}
