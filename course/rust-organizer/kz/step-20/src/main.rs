fn show_title(title: &String) {
    println!("Функцияда: {title}");
}

fn main() {
    let title = String::from("Оқу");
    show_title(&title);
    println!("Шақырудан кейін: {title}");
}
