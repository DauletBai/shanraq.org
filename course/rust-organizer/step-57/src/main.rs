fn normalized_title(text: &str) -> Option<&str> {
    let trimmed = text.trim();
    if trimmed.is_empty() {
        None
    } else {
        Some(trimmed)
    }
}

fn main() {
    let title = normalized_title("  Купить книгу  ").unwrap_or("Пусто");
    println!("Название: {title}");
}

#[test]
fn trims_outer_spaces() {
    assert_eq!(normalized_title("  Книга  "), Some("Книга"));
}
