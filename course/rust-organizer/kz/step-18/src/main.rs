fn main() {
    let fixed: [u32; 2] = [10, 15];
    let mut title = String::from("Жоспар");
    title.push_str(" бүгінге");
    println!("{} және {} минут", fixed[0], fixed[1]);
    println!("{title}");
}
