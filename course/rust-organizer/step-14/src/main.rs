fn remaining(total: u32, done: u32) -> u32 {
    total - done
}

fn print_remaining(count: u32) {
    println!("Осталось: {count}");
}

fn main() {
    let total: u32 = 7;
    let done: u32 = 2;
    if done <= total {
        let count = remaining(total, done);
        print_remaining(count);
    } else {
        println!("Ошибка данных");
    }
}
