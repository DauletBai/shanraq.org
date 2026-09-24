use serde::{Deserialize, Serialize};
use std::collections::HashSet;

fn ordinary_priority() -> u8 {
    2
}

#[derive(Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
struct Task {
    id: u64,
    title: String,
    done: bool,
    #[serde(default = "ordinary_priority")]
    priority: u8,
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
            || !(1..=3).contains(&task.priority)
            || !ids.insert(task.id)
        {
            return Err(String::from("invalid task"));
        }
    }
    Ok(data)
}

fn high_open(data: &Organizer) -> usize {
    data.tasks
        .iter()
        .filter(|task| task.priority == 3 && !task.done)
        .count()
}

fn main() -> Result<(), String> {
    let old_json =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Кітап сатып алу","done":false}]}"#;
    let mut data = decode(old_json)?;
    println!("Ескі тапсырманың басымдығы: {}", data.tasks[0].priority);
    data.tasks.push(Task {
        id: 2,
        title: String::from("Досқа қоңырау шалу"),
        done: false,
        priority: 3,
    });
    data.next_id = 3;
    println!("Шұғыл аяқталмағандар: {}", high_open(&data));
    Ok(())
}

#[test]
fn old_file_gets_ordinary_priority() {
    let text =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Кітап сатып алу","done":false}]}"#;
    assert_eq!(decode(text).unwrap().tasks[0].priority, 2);
}

#[test]
fn invalid_priority_is_rejected() {
    let text = r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Кітап сатып алу","done":false,"priority":0}]}"#;
    assert!(decode(text).is_err());
}
