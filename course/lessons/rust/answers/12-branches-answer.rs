fn main() {
    let total: u32 = 4;
    let done: u32 = 6;
    if done > total {
        println!("Ошибка: выполнено больше общего числа");
    } else if done == total {
        println!("Все задачи выполнены");
    } else {
        println!("Осталось: {}", total - done);
    }
}
