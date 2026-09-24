use std::path::PathBuf;

/// Выбрать путь: команда, окружение, местный вариант по умолчанию.
pub fn choose_path(command: Option<PathBuf>, environment: Option<PathBuf>) -> PathBuf {
    if let Some(path) = command {
        return path;
    }
    if let Some(path) = environment {
        return path;
    }
    PathBuf::from("дела.json")
}

fn main() {
    println!("По умолчанию: {}", choose_path(None, None).display());
    println!(
        "Из окружения: {}",
        choose_path(None, Some(PathBuf::from("семья.json"))).display()
    );
    println!(
        "Из команды: {}",
        choose_path(
            Some(PathBuf::from("работа.json")),
            Some(PathBuf::from("семья.json"))
        )
        .display()
    );
}
