fn main() {
    let fixed: [u32; 2] = [10, 15];
    let mut title = String::from("План");
    title.push_str(" на день");
    println!("{} и {} минут", fixed[0], fixed[1]);
    println!("{title}");
}
