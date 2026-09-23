enum Status {
    Planned,
    Doing(u32),
    Done,
}

fn active_minutes(status: Status) -> u32 {
    match status {
        Status::Planned => 0,
        Status::Doing(minutes) => minutes,
        Status::Done => 0,
    }
}

fn main() {
    let during = active_minutes(Status::Doing(25));
    let after = active_minutes(Status::Done);
    println!("Работа: {during} мин");
    println!("После завершения: {after} мин");
}
