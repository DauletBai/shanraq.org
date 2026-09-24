struct Task {
    id: u32,
    title: String,
    done: bool,
}

fn main() {
    let mut tasks = vec![
        Task {
            id: 1,
            title: String::from("Rust оқу"),
            done: false,
        },
        Task {
            id: 2,
            title: String::from("Ойды жазу"),
            done: false,
        },
    ];
    let mut reader = tasks.iter();
    if let Some(task) = reader.next() {
        println!("Алғашқы: {}", task.title);
    }
    for task in tasks.iter_mut() {
        if task.id == 2 {
            task.done = true;
        }
    }
    for task in tasks.iter() {
        println!("Тапсырма {}: {}", task.id, task.title);
    }
    let mut done_count = 0;
    for task in tasks.iter() {
        if task.done {
            done_count += 1;
        }
    }
    println!("Дайын: {}", done_count);
    for task in tasks.into_iter() {
        println!("Иелік етеміз: {}", task.title);
    }
}
