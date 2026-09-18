fn main() {
    let total: u32 = 7;
    let done: u32 = 2;
    if done > total {
        println!("Ошибка: выполнено больше общего числа");
    } else if done == total {
        println!("Все задачи выполнены");
    } else {
        let remaining = total - done;
        println!("Осталось: {remaining}");
    }
    let priority: u8 = 3;
    let urgent: bool = if priority >= 3 { true } else { false };
    println!("Срочно: {urgent}");
}
