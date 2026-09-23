fn remaining(total: u32, done: u32) -> u32 {
    total - done
}

fn main() {
    println!("Осталось: {}", remaining(7, 2));
}

#[test]
fn subtracts_done_tasks() {
    assert_eq!(remaining(7, 2), 5);
}
