fn show_task(title: &String, minutes: u32) {
    println!("{title} — {minutes} min");
}

fn main() {
    let title = String::from("Reading");
    show_task(&title, 15);
    show_task(&title, 20);
    println!("Title remains: {title}");
}
