fn mark_done(title: &mut String) {
    title.push_str(" — done");
}

fn main() {
    let mut title = String::from("Reading");
    mark_done(&mut title);
    println!("Task: {title}");
}
