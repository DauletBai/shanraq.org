fn is_short(minutes: u32) -> bool {
    minutes > 0 && minutes <= 15
}

fn main() {
    println!("Қысқа: {}", is_short(15));
}

#[test]
fn zero_is_not_a_task() {
    assert_eq!(is_short(0), false);
}

#[test]
fn fifteen_is_short() {
    assert_eq!(is_short(15), true);
}

#[test]
fn sixteen_is_long() {
    assert_eq!(is_short(16), false);
}
