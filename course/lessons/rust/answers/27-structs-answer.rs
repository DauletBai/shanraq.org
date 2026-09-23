struct Task {
    title: String,
    minutes: u32,
    done: bool,
}

fn main() {
    let mut task = Task {
        title: String::from("Уборка"),
        minutes: 15,
        done: false,
    };
    println!(
        "{}: {} мин, сделано {}",
        task.title, task.minutes, task.done
    );
    task.done = true;
    println!("После: {}", task.done);
}
