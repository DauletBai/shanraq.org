fn main() {
    let minutes: [u32; 3] = [10, 20, 15];
    println!("Біріншісі: {}", minutes[0]);
    println!("Саны: {}", minutes.len());
    let mut total: u32 = 0;
    for value in minutes {
        total = total + value;
    }
    println!("Барлық минут: {total}");
}
