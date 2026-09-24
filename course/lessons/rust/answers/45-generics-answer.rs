fn choose<T>(items: &[T], index: usize) -> Option<&T> {
    items.get(index)
}

fn last<T>(items: &[T]) -> Option<&T> {
    items.last()
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
    if let Some(task) = last(&tasks) {
        println!("Последняя: {}", task.title);
    }
    let empty: [u32; 0] = [];
    println!("Пустой список: {}", last(&empty).is_none());
}
