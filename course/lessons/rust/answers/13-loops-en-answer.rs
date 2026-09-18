fn main() {
    let mut total_minutes: u32 = 0;
    for day in 1..=5 {
        if day == 3 {
            continue;
        }
        println!("Day: {day}");
        total_minutes = total_minutes + 15;
    }
    println!("Total minutes: {total_minutes}");
}
