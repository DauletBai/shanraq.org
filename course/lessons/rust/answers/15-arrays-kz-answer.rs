fn is_short(minutes: u32) -> bool {
    minutes > 0 && minutes <= 15
}

fn main() {
    let minutes: [u32; 4] = [0, 10, 25, 15];
    let mut count: u32 = 0;
    let mut total: u32 = 0;
    for value in minutes {
        if is_short(value) {
            count = count + 1;
            total = total + value;
        }
    }
    println!("Қысқа тапсырмалар: {count}");
    println!("Барлық минут: {total}");
}
