fn plan() -> (String, u32) {
    (String::from("Walk"), 30)
}

fn main() {
    let (title, minutes) = plan();
    println!("{title} — {minutes} min");
}
