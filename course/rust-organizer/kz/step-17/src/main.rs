fn count_short(minutes: [u32; 4]) -> u32 {
    let mut count = 0;
    for value in minutes {
        if value > 0 && value <= 15 {
            count = count + 1;
        }
    }
    count
}

fn main() {
    let tasks = [0, 8, 20, 15];
    println!("Қысқа тапсырма: {}", count_short(tasks));
}
