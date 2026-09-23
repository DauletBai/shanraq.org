struct Task {
    title: String,
    minutes: u32,
    done: bool,
}

fn main() {
    let task = Task {
        title: String::from("Reading"),
        minutes: 20,
        done: false,
    };
    println!("Task: {}", task.title);
    println!("Minutes: {}", task.minutes);
    println!("Done: {}", task.done);
}
