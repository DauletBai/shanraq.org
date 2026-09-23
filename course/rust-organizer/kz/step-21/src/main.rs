fn mark_done(title: &mut String) {
    title.push_str(" — дайын");
}

fn main() {
    let mut title = String::from("Оқу");
    mark_done(&mut title);
    println!("Тапсырма: {title}");
}
