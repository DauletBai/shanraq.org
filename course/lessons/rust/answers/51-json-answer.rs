use serde::{Deserialize, Serialize};
use std::collections::HashSet;

#[derive(Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
struct Task {
    id: u64,
    title: String,
    done: bool,
}

#[derive(Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
struct Organizer {
    version: u32,
    next_id: u64,
    tasks: Vec<Task>,
}

fn validate(data: &Organizer) -> Result<(), String> {
    if data.version != 1 || data.next_id == 0 {
        return Err(String::from("invalid version or next ID"));
    }
    let mut ids = HashSet::new();
    for task in &data.tasks {
        if task.id == 0
            || task.id >= data.next_id
            || task.title.trim().is_empty()
            || !ids.insert(task.id)
        {
            return Err(String::from("invalid task"));
        }
    }
    Ok(())
}

fn run() -> Result<(), String> {
    let source = Organizer {
        version: 1,
        next_id: 3,
        tasks: vec![
            Task {
                id: 1,
                title: String::from("Купить книгу"),
                done: false,
            },
            Task {
                id: 2,
                title: String::from("Позвонить другу"),
                done: false,
            },
        ],
    };
    let json = serde_json::to_string(&source).map_err(|error| error.to_string())?;
    let mut loaded: Organizer = serde_json::from_str(&json).map_err(|error| error.to_string())?;
    validate(&loaded)?;
    println!("{json}");
    println!("Восстановлено задач: {}", loaded.tasks.len());
    loaded.tasks[1].id = 1;
    println!("Повтор ID отклонён: {}", validate(&loaded).is_err());
    Ok(())
}

fn main() {
    if let Err(error) = run() {
        eprintln!("Ошибка: {error}");
        std::process::exit(1);
    }
}
