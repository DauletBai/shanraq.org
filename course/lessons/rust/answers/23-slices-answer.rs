fn sum(values: &[u32]) -> u32 {
    let mut total: u32 = 0;
    for index in 0..values.len() {
        total = total + values[index];
    }
    total
}

fn main() {
    let minutes: [u32; 4] = [5, 10, 15, 20];
    println!("Первые: {}", sum(&minutes[0..2]));
    println!("Последние: {}", sum(&minutes[2..4]));
    println!("Всего элементов: {}", minutes.len());
}
