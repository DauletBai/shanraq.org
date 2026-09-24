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
        title: String::from("Read Rust"),
    }];
    if let Some(number) = choose(&numbers, 1) {
        println!("Number: {number}");
    }
    if let Some(task) = choose(&tasks, 0) {
        println!("Task: {}", task.title);
    }
    println!("No item: {}", choose(&numbers, 9).is_none());
    if let Some(task) = last(&tasks) {
        println!("Last: {}", task.title);
    }
    let empty: [u32; 0] = [];
    println!("Empty list: {}", last(&empty).is_none());
}
