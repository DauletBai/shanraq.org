struct Task {
    id: u32,
    title: String,
    done: bool,
}

fn main() {
    let mut tasks = vec![
        Task {
            id: 1,
            title: String::from("Читать Rust"),
            done: false,
        },
        Task {
            id: 2,
            title: String::from("Записать идею"),
            done: false,
        },
    ];
    let mut reader = tasks.iter();
    if let Some(task) = reader.next() {
        println!("Первая: {}", task.title);
    }
    for task in tasks.iter_mut() {
        if task.id == 2 {
            task.done = true;
        }
    }
    for task in tasks.iter() {
        println!("Задача {}: {}", task.id, task.title);
    }
    let mut done_count = 0;
    for task in tasks.iter() {
        if task.done {
            done_count += 1;
        }
    }
    println!("Готово: {}", done_count);
    let mut pending_count = 0;
    for task in tasks.iter() {
        if !task.done {
            pending_count += 1;
        }
    }
    println!("Ожидают: {}", pending_count);
    for task in tasks.into_iter() {
        println!("Владеем: {}", task.title);
    }
}
