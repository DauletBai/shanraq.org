use std::fs::{self, OpenOptions};
use std::io::{ErrorKind, Write};

fn run() -> std::io::Result<()> {
    let path = std::env::temp_dir().join(format!("rust-lesson-50-{}.txt", std::process::id()));
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(&path)?;
    file.write_all("Кітап сатып алу\n".as_bytes())?;
    drop(file);
    let text = fs::read_to_string(&path)?;
    print!("{text}");
    fs::remove_file(&path)?;
    match fs::read_to_string(&path) {
        Err(error) if error.kind() == ErrorKind::NotFound => println!("Файл табылмады"),
        Ok(_) => println!("Күтпеген сәттілік"),
        Err(error) => eprintln!("Файл қатесі: {error}"),
    }
    Ok(())
}

fn main() {
    if let Err(error) = run() {
        eprintln!("Файл қатесі: {error}");
        std::process::exit(1);
    }
}
