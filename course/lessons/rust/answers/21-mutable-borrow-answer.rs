fn add_label(title: &mut String, label: &str) {
    title.push_str(label);
}

fn main() {
    let mut title = String::from("План");
    add_label(&mut title, "!");
    add_label(&mut title, "!");
    println!("{title}");
}
