fn is_short(minutes: u32) -> bool {
    minutes > 0 && minutes <= 15
}

fn main() {
    println!("{}", is_short(0));
    println!("{}", is_short(15));
    println!("{}", is_short(16));
}
