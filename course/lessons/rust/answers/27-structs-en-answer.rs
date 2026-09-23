struct Task {
    title: String,
    minutes: u32,
    done: bool,
}

fn main() {
    let mut task = Task {
        title: String::from("Cleaning"),
        minutes: 15,
        done: false,
    };
    println!("{}: {} min, done {}", task.title, task.minutes, task.done);
    task.done = true;
    println!("After: {}", task.done);
}
