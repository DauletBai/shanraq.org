enum Status {
    Planned,
    Doing(u32),
    Done,
}

fn main() {
    let status = Status::Doing(15);
    match status {
        Status::Planned => println!("Басталмаған"),
        Status::Doing(minutes) => println!("Орындалуда: {minutes} мин"),
        Status::Done => println!("Аяқталған"),
    }
}
