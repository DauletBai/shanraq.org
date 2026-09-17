fn main() {
    let total: u32 = 7;
    let done: u32 = 2;
    let remaining = total - done;
    println!("Қалды: {remaining}");
    println!("Толық жұп: {}", remaining / 2);
    println!("Жұпсыз: {}", remaining % 2);
}
