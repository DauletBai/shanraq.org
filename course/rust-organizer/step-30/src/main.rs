enum Status {
    Planned,
    Doing(u32),
    Done,
}

fn main() {
    let status = Status::Doing(10);
    let active_minutes: u32 = match status {
        Status::Planned => 0,
        Status::Doing(minutes) => minutes,
        Status::Done => 0,
    };
    println!("Сейчас в работе: {active_minutes} мин");
}
