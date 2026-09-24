fn normalized_title(text: &str) -> Option<&str> {
    let trimmed = text.trim();
    if trimmed.is_empty() {
        None
    } else {
        Some(trimmed)
    }
}

fn main() {
    let title = normalized_title("  Кітап сатып алу  ").unwrap_or("Бос");
    println!("Атауы: {title}");
}

#[test]
fn trims_outer_spaces() {
    assert_eq!(normalized_title("  Кітап  "), Some("Кітап"));
}
