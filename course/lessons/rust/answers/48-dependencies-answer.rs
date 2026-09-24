fn label(title: &str) -> String {
    format!("[{title}]")
}

fn main() {
    println!("{}", label("Встреча"));
    println!("{}", label("План"));
}
