use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::fs::{self, OpenOptions};
use std::io::Write;
use std::path::{Path, PathBuf};

fn ordinary_priority() -> u8 {
    2
}

#[derive(Clone, Serialize, Deserialize)]
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

impl Organizer {
    fn new() -> Self {
        Self {
            version: 1,
            next_id: 1,
            tasks: Vec::new(),
        }
    }
    fn validate(&self) -> Result<(), String> {
        if self.version != 1 || self.next_id == 0 {
            return Err(String::from("Неверные данные"));
        }
        let mut ids = HashSet::new();
        for task in &self.tasks {
            if task.id == 0
                || task.id >= self.next_id
                || task.title.trim().is_empty()
                || !(1..=3).contains(&task.priority)
                || !ids.insert(task.id)
            {
                return Err(String::from("Неверные данные"));
            }
        }
        Ok(())
    }
    fn add(&mut self, priority: u8, title: String) -> Result<u64, String> {
        if !(1..=3).contains(&priority) {
            return Err(String::from("Приоритет должен быть от 1 до 3"));
        }
        if title.trim().is_empty() {
            return Err(String::from("Название пусто"));
        }
        let next = self.next_id.checked_add(1).ok_or("Номера закончились")?;
        let id = self.next_id;
        self.tasks.push(Task {
            id,
            title: title.trim().to_owned(),
            done: false,
            priority,
        });
        self.next_id = next;
        Ok(id)
    }
    fn task_mut(&mut self, id: u64) -> Result<&mut Task, String> {
        self.tasks
            .iter_mut()
            .find(|task| task.id == id)
            .ok_or(String::from("ID не найден"))
    }
    fn delete(&mut self, id: u64) -> Result<(), String> {
        let before = self.tasks.len();
        self.tasks.retain(|task| task.id != id);
        if self.tasks.len() == before {
            Err(String::from("ID не найден"))
        } else {
            Ok(())
        }
    }
}

fn decode(text: &str) -> Result<Organizer, String> {
    let data: Organizer = serde_json::from_str(text).map_err(|e| e.to_string())?;
    data.validate()?;
    Ok(data)
}

fn load_or_new(path: &Path) -> Result<Organizer, String> {
    match fs::read_to_string(path) {
        Ok(text) => decode(&text),
        Err(error) if error.kind() == std::io::ErrorKind::NotFound => Ok(Organizer::new()),
        Err(error) => Err(error.to_string()),
    }
}

fn save_safe(path: &Path, data: &Organizer) -> Result<(), String> {
    data.validate()?;
    let text = serde_json::to_string(data).map_err(|e| e.to_string())?;
    let temporary = path.with_extension(format!("tmp-{}", std::process::id()));
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(&temporary)
        .map_err(|e| e.to_string())?;
    if let Err(error) = file
        .write_all(text.as_bytes())
        .and_then(|_| file.sync_all())
    {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error.to_string());
    }
    drop(file);
    if let Err(error) = fs::rename(&temporary, path) {
        let _ = fs::remove_file(&temporary);
        return Err(error.to_string());
    }
    Ok(())
}

fn copy_checked(source: &Path, target: &Path) -> Result<(), String> {
    let text = fs::read_to_string(source).map_err(|e| e.to_string())?;
    decode(&text)?;
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(target)
        .map_err(|e| e.to_string())?;
    if let Err(error) = file
        .write_all(text.as_bytes())
        .and_then(|_| file.sync_all())
    {
        drop(file);
        let _ = fs::remove_file(target);
        return Err(error.to_string());
    }
    drop(file);
    if let Err(error) = fs::read_to_string(target)
        .map_err(|e| e.to_string())
        .and_then(|copied| decode(&copied).map(|_| ()))
    {
        let _ = fs::remove_file(target);
        return Err(error);
    }
    Ok(())
}

fn parse_id(text: &str) -> Result<u64, String> {
    match text.parse::<u64>() {
        Ok(id) if id > 0 => Ok(id),
        _ => Err(String::from("Нужен положительный ID")),
    }
}

fn parse_priority(text: &str) -> Result<u8, String> {
    match text.parse::<u8>() {
        Ok(priority) if (1..=3).contains(&priority) => Ok(priority),
        _ => Err(String::from("Приоритет должен быть от 1 до 3")),
    }
}

fn show(tasks: &[&Task]) -> String {
    if tasks.is_empty() {
        return String::from("Задач нет");
    }
    tasks
        .iter()
        .map(|task| {
            let status = if task.done {
                "готово"
            } else {
                "ожидает"
            };
            format!("{} [{}] {} ({status})", task.id, task.priority, task.title)
        })
        .collect::<Vec<_>>()
        .join("\n")
}

fn execute(args: &[String], path: &Path) -> Result<String, (i32, String)> {
    let usage = || {
        (
            2,
            String::from("Неизвестная команда или неверные аргументы"),
        )
    };
    let data_error = |message: String| (1, message);
    if args.is_empty() || args == ["help"] {
        return Ok(String::from(
            "Команды: add ПРИОРИТЕТ НАЗВАНИЕ | list | done ID | delete ID | rename ID НАЗВАНИЕ | search ТЕКСТ | filter open/done/high | sort id/title/priority | stats | backup ПУТЬ | restore КОПИЯ НОВЫЙ_ПУТЬ | --version",
        ));
    }
    if args == ["--version"] {
        return Ok(env!("CARGO_PKG_VERSION").to_owned());
    }
    match args[0].as_str() {
        "backup" if args.len() == 2 => {
            copy_checked(path, Path::new(&args[1])).map_err(data_error)?;
            Ok(String::from("Копия создана"))
        }
        "restore" if args.len() == 3 => {
            copy_checked(Path::new(&args[1]), Path::new(&args[2])).map_err(data_error)?;
            Ok(String::from("Восстановлено отдельно"))
        }
        "add" if args.len() >= 3 => {
            let priority = parse_priority(&args[1]).map_err(|e| (2, e))?;
            let mut data = load_or_new(path).map_err(data_error)?;
            let id = data
                .add(priority, args[2..].join(" "))
                .map_err(|e| (2, e))?;
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Добавлено: {id}"))
        }
        "list" if args.len() == 1 => {
            let data = load_or_new(path).map_err(data_error)?;
            let mut tasks: Vec<&Task> = data.tasks.iter().collect();
            tasks.sort_by_key(|task| task.id);
            Ok(show(&tasks))
        }
        "done" if args.len() == 2 => {
            let id = parse_id(&args[1]).map_err(|e| (2, e))?;
            let mut data = load_or_new(path).map_err(data_error)?;
            data.task_mut(id).map_err(data_error)?.done = true;
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Завершено: {id}"))
        }
        "delete" if args.len() == 2 => {
            let id = parse_id(&args[1]).map_err(|e| (2, e))?;
            let mut data = load_or_new(path).map_err(data_error)?;
            data.delete(id).map_err(data_error)?;
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Удалено: {id}"))
        }
        "rename" if args.len() >= 3 => {
            let id = parse_id(&args[1]).map_err(|e| (2, e))?;
            let title = args[2..].join(" ");
            if title.trim().is_empty() {
                return Err((2, String::from("Название пусто")));
            }
            let mut data = load_or_new(path).map_err(data_error)?;
            data.task_mut(id).map_err(data_error)?.title = title.trim().to_owned();
            save_safe(path, &data).map_err(data_error)?;
            Ok(format!("Переименовано: {id}"))
        }
        "search" if args.len() >= 2 => {
            let term = args[1..].join(" ");
            if term.trim().is_empty() {
                return Err(usage());
            }
            let data = load_or_new(path).map_err(data_error)?;
            let tasks: Vec<&Task> = data
                .tasks
                .iter()
                .filter(|task| task.title.contains(&term))
                .collect();
            Ok(show(&tasks))
        }
        "filter" if args.len() == 2 => {
            let data = load_or_new(path).map_err(data_error)?;
            let tasks: Vec<&Task> = match args[1].as_str() {
                "open" => data.tasks.iter().filter(|task| !task.done).collect(),
                "done" => data.tasks.iter().filter(|task| task.done).collect(),
                "high" => data
                    .tasks
                    .iter()
                    .filter(|task| task.priority == 3 && !task.done)
                    .collect(),
                _ => return Err(usage()),
            };
            Ok(show(&tasks))
        }
        "sort" if args.len() == 2 => {
            let data = load_or_new(path).map_err(data_error)?;
            let mut tasks: Vec<&Task> = data.tasks.iter().collect();
            match args[1].as_str() {
                "id" => tasks.sort_by_key(|task| task.id),
                "title" => tasks.sort_by(|a, b| a.title.cmp(&b.title)),
                "priority" => {
                    tasks.sort_by(|a, b| b.priority.cmp(&a.priority).then(a.id.cmp(&b.id)))
                }
                _ => return Err(usage()),
            }
            Ok(show(&tasks))
        }
        "stats" if args.len() == 1 => {
            let data = load_or_new(path).map_err(data_error)?;
            let total = data.tasks.len();
            let done = data.tasks.iter().filter(|task| task.done).count();
            let open = total - done;
            let high = data
                .tasks
                .iter()
                .filter(|task| task.priority == 3 && !task.done)
                .count();
            Ok(format!(
                "Всего: {total}, готово: {done}, осталось: {open}, срочных открытых: {high}"
            ))
        }
        _ => Err(usage()),
    }
}

fn main() {
    let mut args: Vec<String> = std::env::args().skip(1).collect();
    let path = if args.first().map(String::as_str) == Some("--data") {
        if args.len() < 2 {
            eprintln!("Ошибка: Неизвестная команда или неверные аргументы");
            std::process::exit(2);
        }
        args.remove(0);
        PathBuf::from(args.remove(0))
    } else {
        std::env::var_os("ORGANIZER_DATA")
            .map(PathBuf::from)
            .unwrap_or_else(|| PathBuf::from("tasks.json"))
    };
    match execute(&args, &path) {
        Ok(message) => println!("{message}"),
        Err((code, message)) => {
            eprintln!("Ошибка: {message}");
            std::process::exit(code);
        }
    }
}

#[test]
fn damaged_json_is_rejected() {
    assert!(decode("{broken").is_err());
}

#[test]
fn unknown_version_is_rejected() {
    let text = r#"{"version":2,"next_id":1,"tasks":[]}"#;
    assert!(decode(text).is_err());
}

#[test]
fn deleting_one_task_keeps_the_next_id() {
    let mut data = Organizer::new();
    let first = data.add(2, String::from("Купить книгу")).unwrap();
    let second = data.add(3, String::from("Купить книгу")).unwrap();
    data.delete(first).unwrap();
    assert_eq!(data.tasks[0].id, second);
    assert_eq!(data.next_id, second + 1);
}

#[test]
fn old_json_gets_normal_priority() {
    let text =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Купить книгу","done":false}]}"#;
    assert_eq!(decode(text).unwrap().tasks[0].priority, 2);
}

#[test]
fn duplicate_id_is_rejected() {
    let text = r#"{"version":1,"next_id":3,"tasks":[{"id":1,"title":"Купить книгу","done":false,"priority":2},{"id":1,"title":"Позвонить другу","done":false,"priority":3}]}"#;
    assert!(decode(text).is_err());
}

#[test]
fn invalid_priority_is_rejected() {
    let text =
        r#"{"version":1,"next_id":2,"tasks":[{"id":1,"title":"Купить книгу","done":false,"priority":0}]}"#;
    assert!(decode(text).is_err());
}
