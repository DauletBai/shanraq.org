fn main() {
    let total: u32 = 7;
    let done: u32 = 2;
    let remaining = total - done;
    println!("Remaining: {remaining}");
    println!("Complete pairs: {}", remaining / 2);
    println!("Unpaired: {}", remaining % 2);
}
