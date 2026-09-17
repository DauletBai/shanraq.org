fn main() {
    let total: u32 = 7;
    let done: u32 = 2;
    let remaining = total - done;
    println!("Осталось: {remaining}");
    println!("Полных пар: {}", remaining / 2);
    println!("Без пары: {}", remaining % 2);
}
