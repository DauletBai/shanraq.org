fn show_task(title: &str, minutes: u32) {
    println!("{title}: {minutes} min");
}

fn main() {
    let owned = String::from("Reading");
    let literal: &str = "Rest";
    show_task(owned.as_str(), 15);
    show_task(literal, 5);
    println!("Stored: {owned}");
}
