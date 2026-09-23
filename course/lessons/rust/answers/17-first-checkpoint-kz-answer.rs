fn is_short(minutes: u32) -> bool {
    minutes > 0 && minutes <= 15
}

fn summarize(minutes: [u32; 4]) -> [u32; 2] {
    let mut count = 0;
    let mut total = 0;
    for value in minutes {
        if is_short(value) {
            count = count + 1;
            total = total + value;
        }
    }
    [count, total]
}

fn main() {
    let result = summarize([0, 8, 20, 15]);
    println!("Қысқа тапсырма: {}", result[0]);
    println!("Минут: {}", result[1]);
}

#[test]
fn no_tasks() {
    assert_eq!(summarize([0, 0, 0, 0]), [0, 0]);
}

#[test]
fn boundaries() {
    assert_eq!(summarize([1, 15, 16, 0]), [2, 16]);
}

#[test]
fn ordinary_case() {
    assert_eq!(summarize([0, 8, 20, 15]), [2, 23]);
}
