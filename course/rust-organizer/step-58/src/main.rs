fn main() {
    println!("Версия: {}", env!("CARGO_PKG_VERSION"));
}

#[test]
fn package_version_has_three_parts() {
    let parts: Vec<&str> = env!("CARGO_PKG_VERSION").split('.').collect();
    assert_eq!(parts.len(), 3);
}
