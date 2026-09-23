struct Task {
    title: String,
    minutes: u32,
    done: bool,
}

fn main() {
    let task = Task {
        title: String::from("Оқу"),
        minutes: 20,
        done: false,
    };
    println!("Тапсырма: {}", task.title);
    println!("Минут: {}", task.minutes);
    println!("Орындалды: {}", task.done);
}
