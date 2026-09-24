fn choose<T>(items: &[T], index: usize) -> Option<&T> {
    items.get(index)
}

struct Task {
    title: String,
}

fn main() {
    let numbers = [10, 20];
    let tasks = [Task {
        title: String::from("Читать Rust"),
    }];
    if let Some(number) = choose(&numbers, 1) {
        println!("Число: {number}");
    }
    if let Some(task) = choose(&tasks, 0) {
        println!("Задача: {}", task.title);
    }
    println!("Нет элемента: {}", choose(&numbers, 9).is_none());
}
