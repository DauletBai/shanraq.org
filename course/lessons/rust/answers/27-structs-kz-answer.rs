struct Task {
    title: String,
    minutes: u32,
    done: bool,
}

fn main() {
    let mut task = Task {
        title: String::from("Тазалау"),
        minutes: 15,
        done: false,
    };
    println!(
        "{}: {} мин, орындалды {}",
        task.title, task.minutes, task.done
    );
    task.done = true;
    println!("Кейін: {}", task.done);
}
