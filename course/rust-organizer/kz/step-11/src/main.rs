fn main() {
    let done: bool = false;
    let priority: u8 = 3;
    let minutes: u32 = 20;
    let show = !done && (priority >= 3 || minutes <= 15);
    println!("Көрсету: {show}");
    let slots: u32 = 0;
    let tasks: u32 = 6;
    let fits = slots != 0 && tasks / slots <= 2;
    println!("Сыяды: {fits}");
}
