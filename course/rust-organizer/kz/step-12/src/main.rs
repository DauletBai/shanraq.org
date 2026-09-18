fn main() {
    let total: u32 = 7;
    let done: u32 = 2;
    if done > total {
        println!("Қате: аяқталғаны жалпы саннан көп");
    } else if done == total {
        println!("Барлық тапсырма аяқталды");
    } else {
        let remaining = total - done;
        println!("Қалғаны: {remaining}");
    }
    let priority: u8 = 3;
    let urgent: bool = if priority >= 3 { true } else { false };
    println!("Шұғыл: {urgent}");
}
