fn main() {
    let mut title = String::from("Task");
    title.push_str(": reading");
    {
        let minutes = [5, 10];
        println!("Minutes: {}", minutes[0] + minutes[1]);
    }
    println!("{title}");
}
