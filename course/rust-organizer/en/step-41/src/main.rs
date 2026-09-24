struct Task {
    id: u32,
    title: String,
    done: bool,
}

fn main() {
    let mut tasks = vec![
        Task {
            id: 1,
            title: String::from("Read Rust"),
            done: false,
        },
        Task {
            id: 2,
            title: String::from("Write an idea"),
            done: false,
        },
    ];
    let mut reader = tasks.iter();
    if let Some(task) = reader.next() {
        println!("First: {}", task.title);
    }
    for task in tasks.iter_mut() {
        if task.id == 2 {
            task.done = true;
        }
    }
    for task in tasks.iter() {
        println!("Task {}: {}", task.id, task.title);
    }
    let mut done_count = 0;
    for task in tasks.iter() {
        if task.done {
            done_count += 1;
        }
    }
    println!("Done: {}", done_count);
    for task in tasks.into_iter() {
        println!("Owned: {}", task.title);
    }
}
