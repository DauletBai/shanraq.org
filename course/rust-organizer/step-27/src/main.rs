struct Task {
    title: String,
    minutes: u32,
    done: bool,
}

fn main() {
    let task = Task {
        title: String::from("Чтение"),
        minutes: 20,
        done: false,
    };
    println!("Задача: {}", task.title);
    println!("Минут: {}", task.minutes);
    println!("Сделано: {}", task.done);
}
