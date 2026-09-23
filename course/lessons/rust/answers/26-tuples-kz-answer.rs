fn plan() -> (String, u32) {
    (String::from("Серуен"), 30)
}

fn main() {
    let (title, minutes) = plan();
    println!("{title} — {minutes} мин");
}
