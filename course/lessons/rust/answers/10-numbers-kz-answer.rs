fn main() {
    let total: u32 = 10;
    let done: u32 = 4;
    let remaining = total - done;
    println!("Қалды: {remaining}");
    println!("Толық жұп: {}", remaining / 2);
    println!("Жұпсыз: {}", remaining % 2);
}
