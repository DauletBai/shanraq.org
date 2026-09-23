fn mark_done(title: &mut String) {
    title.push_str(" — готово");
}

fn main() {
    let mut title = String::from("Чтение");
    mark_done(&mut title);
    println!("Задача: {title}");
}
