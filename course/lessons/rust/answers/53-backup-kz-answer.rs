use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::fs::{self, OpenOptions};
use std::io::{self, Write};
use std::path::Path;

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

fn load(text: &str) -> io::Result<Organizer> {
    let data: Organizer = serde_json::from_str(text)
        .map_err(|error| io::Error::new(io::ErrorKind::InvalidData, error))?;
    if data.version != 1 || data.next_id == 0 {
        return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid header"));
    }
    let mut ids = HashSet::new();
    for task in &data.tasks {
        if task.id == 0
            || task.id >= data.next_id
            || task.title.trim().is_empty()
            || !ids.insert(task.id)
        {
            return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid task"));
        }
    }
    Ok(data)
}

fn copy_checked(source: &Path, target: &Path) -> io::Result<()> {
    let text = fs::read_to_string(source)?;
    load(&text)?;
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(target)?;
    file.write_all(text.as_bytes())?;
    file.sync_all()?;
    drop(file);
    load(&fs::read_to_string(target)?)?;
    Ok(())
}

fn run() -> io::Result<()> {
    let folder = std::env::temp_dir().join(format!("rust-lesson-53-{}", std::process::id()));
    fs::create_dir(&folder)?;
    let primary = folder.join("tasks.json");
    let backup = folder.join("backup.json");
    let recovered = folder.join("recovered.json");
    let data = Organizer {
        version: 1,
        next_id: 2,
        tasks: vec![Task {
            id: 1,
            title: String::from("Кітап сатып алу"),
            done: false,
        }],
    };
    let json = serde_json::to_string(&data)
        .map_err(|error| io::Error::new(io::ErrorKind::InvalidData, error))?;
    fs::write(&primary, json)?;
    copy_checked(&primary, &backup)?;
    fs::write(&primary, "{broken")?;
    println!(
        "Негізгі файл бүлінген: {}",
        load(&fs::read_to_string(&primary)?).is_err()
    );
    copy_checked(&backup, &recovered)?;
    let restored = load(&fs::read_to_string(&recovered)?)?;
    println!("Қалпына келді: {}", restored.tasks[0].title);
    println!(
        "Бастапқы файл өзгермеді: {}",
        fs::read_to_string(&primary)? == "{broken"
    );
    let damaged_backup = folder.join("damaged-backup.json");
    let other = folder.join("other-recovery.json");
    fs::write(&damaged_backup, "{broken")?;
    let rejected = copy_checked(&damaged_backup, &other).is_err() && !other.exists();
    println!("Жарамсыз көшірме қабылданбады: {}", rejected);
    fs::remove_dir_all(folder)?;
    Ok(())
}

fn main() {
    if let Err(error) = run() {
        eprintln!("Қате: {error}");
        std::process::exit(1);
    }
}
