use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::fs;
use std::path::{Path, PathBuf};

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

fn read_checked(path: &Path) -> Result<Organizer, String> {
    let text = fs::read_to_string(path).map_err(|error| error.to_string())?;
    let data: Organizer = serde_json::from_str(&text).map_err(|error| error.to_string())?;
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

fn practice_path(label: &str) -> PathBuf {
    std::env::temp_dir().join(format!(
        "rust-course-56-{}-{label}.json",
        std::process::id()
    ))
}

fn sample(title: &str) -> Organizer {
    Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from(title),
            done: false,
        }],
    }
}

fn main() -> Result<(), String> {
    let path = practice_path("main");
    let text = serde_json::to_string(&sample("Купить книгу")).map_err(|e| e.to_string())?;
    fs::write(&path, text).map_err(|e| e.to_string())?;
    let loaded = read_checked(&path)?;
    println!("Прочитано задач: {}", loaded.tasks.len());
    fs::remove_file(path).map_err(|e| e.to_string())?;
    Ok(())
}

#[test]
fn unicode_round_trip() {
    let path = practice_path("unicode");
    let text = serde_json::to_string(&sample("Café")).unwrap();
    fs::write(&path, text).unwrap();
    let loaded = read_checked(&path).unwrap();
    assert_eq!(loaded.tasks[0].title, "Café");
    fs::remove_file(path).unwrap();
}

#[test]
fn damaged_json_is_rejected() {
    let path = practice_path("damaged");
    fs::write(&path, "{broken").unwrap();
    assert!(read_checked(&path).is_err());
    assert_eq!(fs::read_to_string(&path).unwrap(), "{broken");
    fs::remove_file(path).unwrap();
}
