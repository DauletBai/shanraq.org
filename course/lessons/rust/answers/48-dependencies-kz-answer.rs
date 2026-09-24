fn label(title: &str) -> String {
    format!("[{title}]")
}

fn main() {
    println!("{}", label("Кездесу"));
    println!("{}", label("Жоспар"));
}
