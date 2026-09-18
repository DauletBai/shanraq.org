fn main() {
    let total: u32 = 4;
    let done: u32 = 6;
    if done > total {
        println!("Error: completed exceeds total");
    } else if done == total {
        println!("All tasks completed");
    } else {
        println!("Remaining: {}", total - done);
    }
}
