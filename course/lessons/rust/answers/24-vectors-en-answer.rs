fn main() {
    let mut titles: Vec<String> = Vec::new();
    titles.push(String::from("Reading"));
    titles.push(String::from("Walk"));
    titles.push(String::from("Rest"));
    let removed = titles.remove(1);
    println!("Removed: {removed}");
    for title in &titles {
        println!("Task: {title}");
    }
    println!("Remaining: {}", titles.len());
}
