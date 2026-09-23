fn main() {
    let fixed: [u32; 2] = [10, 15];
    let mut title = String::from("Plan");
    title.push_str(" for today");
    println!("{} and {} minutes", fixed[0], fixed[1]);
    println!("{title}");
}
