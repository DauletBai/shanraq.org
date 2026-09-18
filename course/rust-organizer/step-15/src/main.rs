fn main() {
    let minutes: [u32; 3] = [10, 20, 15];
    println!("Первая: {}", minutes[0]);
    println!("Количество: {}", minutes.len());
    let mut total: u32 = 0;
    for value in minutes {
        total = total + value;
    }
    println!("Всего минут: {total}");
}
