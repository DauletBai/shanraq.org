fn normalized_title(text: &str) -> Option<&str> {
    let trimmed = text.trim();
    if trimmed.is_empty() {
        None
    } else {
        Some(trimmed)
    }
}

fn main() {
    let title = normalized_title("  Buy a book  ").unwrap_or("Empty");
    println!("Title: {title}");
}

#[test]
fn trims_outer_spaces() {
    assert_eq!(normalized_title("  Book  "), Some("Book"));
}

#[test]
fn rejects_blank_title() {
    assert_eq!(normalized_title("   "), None);
}

#[test]
fn keeps_unicode_title() {
    assert_eq!(normalized_title("  Café  "), Some("Café"));
}
