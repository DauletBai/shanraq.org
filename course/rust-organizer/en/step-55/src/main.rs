use std::path::PathBuf;

/// Choose the data path: command line, environment, then local default.
pub fn choose_path(command: Option<PathBuf>, environment: Option<PathBuf>) -> PathBuf {
    if let Some(path) = command {
        return path;
    }
    if let Some(path) = environment {
        return path;
    }
    PathBuf::from("tasks.json")
}

fn main() {
    println!("Default path: {}", choose_path(None, None).display());
    println!(
        "Environment path: {}",
        choose_path(None, Some(PathBuf::from("family.json"))).display()
    );
}
