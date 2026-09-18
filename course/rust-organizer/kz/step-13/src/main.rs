fn main() {
    let mut total_minutes: u32 = 0;
    for day in 1..=3 {
        println!("Күн: {day}");
        total_minutes = total_minutes + 10;
    }
    println!("Барлық минут: {total_minutes}");
}
