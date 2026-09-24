fn choose<T>(items: &[T], index: usize) -> Option<&T> {
    items.get(index)
}

struct Task {
    title: String,
}

fn main() {
    let numbers = [10, 20];
    let tasks = [Task {
        title: String::from("Rust оқу"),
    }];
    if let Some(number) = choose(&numbers, 1) {
        println!("Сан: {number}");
    }
    if let Some(task) = choose(&tasks, 0) {
        println!("Тапсырма: {}", task.title);
    }
    println!("Элемент жоқ: {}", choose(&numbers, 9).is_none());
}
