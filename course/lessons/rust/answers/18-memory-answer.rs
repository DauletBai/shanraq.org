fn main() {
    let mut title = String::from("Задача");
    title.push_str(": чтение");
    {
        let minutes = [5, 10];
        println!("Минуты: {}", minutes[0] + minutes[1]);
    }
    println!("{title}");
}
