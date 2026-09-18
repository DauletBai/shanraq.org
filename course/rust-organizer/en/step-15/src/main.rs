fn main() {
    let minutes: [u32; 3] = [10, 20, 15];
    println!("First: {}", minutes[0]);
    println!("Count: {}", minutes.len());
    let mut total: u32 = 0;
    for value in minutes {
        total = total + value;
    }
    println!("Total minutes: {total}");
}
