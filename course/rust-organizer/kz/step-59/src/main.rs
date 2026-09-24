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

fn decode(text: &str) -> Result<Organizer, String> {
    let data: Organizer = serde_json::from_str(text).map_err(|e| e.to_string())?;
    if data.version != 1 || data.next_id == 0 {
        return Err(String::from("invalid header"));
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
    Ok(data)
}

fn main() -> Result<(), String> {
    let old_json =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Кітап сатып алу","done":false}]}"#;
    let data = decode(old_json)?;
    println!("Тапсырма саны: {}", data.tasks.len());
    Ok(())
}

#[test]
fn reads_old_task() {
    let text =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Кітап сатып алу","done":false}]}"#;
    assert_eq!(decode(text).unwrap().tasks.len(), 1);
}
