use std::path::{Path, PathBuf};

fn main() {
    let folder = Path::new("data");
    let file: PathBuf = folder.join("tasks.txt");
    println!("{}", file.is_absolute());
    println!("{}", file.file_name().unwrap().to_string_lossy());
    let backup = folder.join("backup.txt");
    println!("{}", backup.file_name().unwrap().to_string_lossy());
}
