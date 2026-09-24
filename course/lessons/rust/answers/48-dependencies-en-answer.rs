fn label(title: &str) -> String {
    format!("[{title}]")
}

fn main() {
    println!("{}", label("Meeting"));
    println!("{}", label("Plan"));
}
