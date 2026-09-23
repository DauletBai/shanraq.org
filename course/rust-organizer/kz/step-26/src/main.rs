fn first_task() -> (String, u32) {
    let title = String::from("Оқу");
    let minutes = 20;
    (title, minutes)
}

fn main() {
    let (title, minutes) = first_task();
    println!("{title}: {minutes} мин");
}
