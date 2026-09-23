fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Reading"));
    titles.push(String::from("Rest"));
    println!("Task count: {}", titles.len());
    for title in &titles {
        println!("Task: {title}");
    }
}
