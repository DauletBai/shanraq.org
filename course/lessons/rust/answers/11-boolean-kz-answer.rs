fn main() {
    let done: bool = false;
    let archived: bool = false;
    let minutes: u32 = 15;
    let show = !done && !archived && minutes <= 15;
    println!("Көрсету: {show}");
}
