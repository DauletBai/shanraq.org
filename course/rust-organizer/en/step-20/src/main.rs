fn show_title(title: &String) {
    println!("Inside: {title}");
}

fn main() {
    let title = String::from("Reading");
    show_title(&title);
    println!("After the call: {title}");
}
