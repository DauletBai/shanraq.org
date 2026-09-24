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
        title: String::from("Rust оқу"),
    }];
    if let Some(number) = choose(&numbers, 1) {
        println!("Сан: {number}");
    }
    if let Some(task) = choose(&tasks, 0) {
        println!("Тапсырма: {}", task.title);
    }
    println!("Элемент жоқ: {}", choose(&numbers, 9).is_none());
    if let Some(task) = last(&tasks) {
        println!("Соңғы: {}", task.title);
    }
    let empty: [u32; 0] = [];
    println!("Бос тізім: {}", last(&empty).is_none());
}
