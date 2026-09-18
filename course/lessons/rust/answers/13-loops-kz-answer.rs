fn main() {
    let mut total_minutes: u32 = 0;
    for day in 1..=5 {
        if day == 3 {
            continue;
        }
        println!("Күн: {day}");
        total_minutes = total_minutes + 15;
    }
    println!("Барлық минут: {total_minutes}");
}
