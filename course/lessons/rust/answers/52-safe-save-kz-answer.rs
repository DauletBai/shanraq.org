use std::fs::{self, OpenOptions};
use std::io::{self, Write};
use std::path::Path;

fn safe_save(path: &Path, text: &str) -> io::Result<()> {
    let temporary = path.with_extension(format!("tmp-{}", std::process::id()));
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(&temporary)?;
    if let Err(error) = file.write_all(text.as_bytes()) {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    if let Err(error) = file.sync_all() {
        drop(file);
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    drop(file);
    if let Err(error) = fs::rename(&temporary, path) {
        let _ = fs::remove_file(&temporary);
        return Err(error);
    }
    Ok(())
}

fn run() -> io::Result<()> {
    let folder = std::env::temp_dir().join(format!("rust-lesson-52-{}", std::process::id()));
    fs::create_dir(&folder)?;
    let path = folder.join("tasks.json");
    fs::write(&path, "Ескі жоспар")?;
    let temporary = path.with_extension(format!("tmp-{}", std::process::id()));
    fs::write(&temporary, "occupied")?;
    let refused = safe_save(&path, "Жаңа жоспар").is_err();
    let kept = fs::read_to_string(&path)? == "Ескі жоспар";
    println!("Ескі дерек сақталды: {}", refused && kept);
    fs::remove_file(&temporary)?;
    safe_save(&path, "Жаңа жоспар")?;
    println!("Сақталды: {}", fs::read_to_string(&path)?);
    fs::remove_dir_all(folder)?;
    Ok(())
}

fn main() {
    if let Err(error) = run() {
        eprintln!("Қате: {error}");
        std::process::exit(1);
    }
}
