enum Status {
    Planned,
    Doing(u32),
    Done,
}

fn main() {
    let status = Status::Doing(15);
    match status {
        Status::Planned => println!("Не начато"),
        Status::Doing(minutes) => println!("В работе: {minutes} мин"),
        Status::Done => println!("Завершено"),
    }
}
