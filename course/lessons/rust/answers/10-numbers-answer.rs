fn main() {
    let total: u32 = 10;
    let done: u32 = 4;
    let remaining = total - done;
    println!("Осталось: {remaining}");
    println!("Полных пар: {}", remaining / 2);
    println!("Без пары: {}", remaining % 2);
}
