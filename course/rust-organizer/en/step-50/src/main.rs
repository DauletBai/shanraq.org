use std::fs::{self, OpenOptions};
use std::io::Write;

fn run() -> std::io::Result<()> {
    let path = std::env::temp_dir().join(format!("rust-lesson-50-{}.txt", std::process::id()));
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(&path)?;
    file.write_all("Buy a book\n".as_bytes())?;
    drop(file);
    let text = fs::read_to_string(&path)?;
    print!("{text}");
    fs::remove_file(&path)?;
    Ok(())
}

fn main() {
    if let Err(error) = run() {
        eprintln!("File error: {error}");
        std::process::exit(1);
    }
}
