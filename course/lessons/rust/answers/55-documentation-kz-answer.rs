use std::path::PathBuf;

/// Дерек жолын таңдау: пәрмен, орта, содан соң әдепкі жергілікті жол.
pub fn choose_path(command: Option<PathBuf>, environment: Option<PathBuf>) -> PathBuf {
    if let Some(path) = command {
        return path;
    }
    if let Some(path) = environment {
        return path;
    }
    PathBuf::from("тапсырмалар.json")
}

fn main() {
    println!("Әдепкі жол: {}", choose_path(None, None).display());
    println!(
        "Ортадағы жол: {}",
        choose_path(None, Some(PathBuf::from("отбасы.json"))).display()
    );
    println!(
        "Пәрмендегі жол: {}",
        choose_path(
            Some(PathBuf::from("жұмыс.json")),
            Some(PathBuf::from("отбасы.json"))
        )
        .display()
    );
}
