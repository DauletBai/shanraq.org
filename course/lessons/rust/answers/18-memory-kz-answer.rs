fn main() {
    let mut title = String::from("Тапсырма");
    title.push_str(": оқу");
    {
        let minutes = [5, 10];
        println!("Минут: {}", minutes[0] + minutes[1]);
    }
    println!("{title}");
}
