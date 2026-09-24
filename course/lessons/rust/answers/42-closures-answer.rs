struct Task {
    id: u32,
    title: String,
}

fn main() {
    let tasks = vec![
        Task {
            id: 3,
            title: String::from("Читать Rust"),
        },
        Task {
            id: 2,
            title: String::from("Купить хлеб"),
        },
        Task {
            id: 1,
            title: String::from("Повторить Rust"),
        },
    ];
    let query = "Rust";
    let mut found: Vec<&Task> = tasks
        .iter()
        .filter(|task| task.title.contains(query))
        .collect();
    found.sort_by_key(|task| task.id);
    let names: Vec<&str> = found.iter().map(|task| task.title.as_str()).collect();
    println!("Найдено: {}", names.len());
    for task in found {
        println!("{}: {}", task.id, task.title);
    }
    let second_query = "хлеб";
    let second_count = tasks
        .iter()
        .filter(|task| task.title.contains(second_query))
        .count();
    println!("По слову «{second_query}»: {second_count}");
}
