struct Task {
    title: String,
    done: bool,
}

fn add(tasks: &mut Vec<Task>, title: &str) -> Result<(), String> {
    if title.trim().is_empty() {
        return Err(String::from("Атау бос"));
    }
    tasks.push(Task {
        title: String::from(title),
        done: false,
    });
    Ok(())
}

fn rename(tasks: &mut Vec<Task>, index: usize, title: &str) -> Result<(), String> {
    match tasks.get_mut(index) {
        Some(task) => {
            task.title = String::from(title);
            Ok(())
        }
        None => Err(String::from("Ондай орын жоқ")),
    }
}

fn complete(tasks: &mut Vec<Task>, index: usize) -> Result<(), String> {
    match tasks.get_mut(index) {
        Some(task) => {
            task.done = true;
            Ok(())
        }
        None => Err(String::from("Ондай орын жоқ")),
    }
}

fn remove(tasks: &mut Vec<Task>, index: usize) -> Result<(), String> {
    if index >= tasks.len() {
        return Err(String::from("Ондай орын жоқ"));
    }
    tasks.remove(index);
    Ok(())
}

fn show(tasks: &[Task]) {
    for task in tasks {
        if task.done {
            println!("[дайын] {}", task.title);
        } else {
            println!("[ ] {}", task.title);
        }
    }
}

fn main() {
    let mut tasks = Vec::new();
    add(&mut tasks, "Нан алу").expect("valid title");
    add(&mut tasks, "Тарау оқу").expect("valid title");
    rename(&mut tasks, 1, "Екі тарау оқу").expect("existing position");
    complete(&mut tasks, 0).expect("existing position");
    remove(&mut tasks, 0).expect("existing position");
    show(&tasks);
}
