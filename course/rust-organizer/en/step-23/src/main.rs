fn sum(values: &[u32]) -> u32 {
    let mut total: u32 = 0;
    for index in 0..values.len() {
        total = total + values[index];
    }
    total
}

fn main() {
    let minutes: [u32; 4] = [10, 20, 15, 5];
    let middle = &minutes[1..3];
    println!("Middle: {}", sum(middle));
    println!("First: {}", minutes[0]);
}
