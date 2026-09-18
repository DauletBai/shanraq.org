fn main() {
    let total: u32 = 4;
    let done: u32 = 6;
    if done > total {
        println!("Қате: аяқталғаны жалпы саннан көп");
    } else if done == total {
        println!("Барлық тапсырма аяқталды");
    } else {
        println!("Қалғаны: {}", total - done);
    }
}
