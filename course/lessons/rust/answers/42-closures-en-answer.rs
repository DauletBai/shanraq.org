struct Task {
    id: u32,
    title: String,
}

fn main() {
    let tasks = vec![
        Task {
            id: 3,
            title: String::from("Read Rust"),
        },
        Task {
            id: 2,
            title: String::from("Buy bread"),
        },
        Task {
            id: 1,
            title: String::from("Review Rust"),
        },
    ];
    let query = "Rust";
    let mut found: Vec<&Task> = tasks
        .iter()
        .filter(|task| task.title.contains(query))
        .collect();
    found.sort_by_key(|task| task.id);
    let names: Vec<&str> = found.iter().map(|task| task.title.as_str()).collect();
    println!("Found: {}", names.len());
    for task in found {
        println!("{}: {}", task.id, task.title);
    }
    let second_query = "bread";
    let second_count = tasks
        .iter()
        .filter(|task| task.title.contains(second_query))
        .count();
    println!("Found {second_query}: {second_count}");
}
