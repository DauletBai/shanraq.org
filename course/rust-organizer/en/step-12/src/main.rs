fn main() {
    let total: u32 = 7;
    let done: u32 = 2;
    if done > total {
        println!("Error: completed exceeds total");
    } else if done == total {
        println!("All tasks completed");
    } else {
        let remaining = total - done;
        println!("Remaining: {remaining}");
    }
    let priority: u8 = 3;
    let urgent: bool = if priority >= 3 { true } else { false };
    println!("Urgent: {urgent}");
}
