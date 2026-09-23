enum Status {
    Planned,
    Doing(u32),
    Done,
}

fn main() {
    let status = Status::Doing(25);
    match status {
        Status::Planned => println!("Not started"),
        Status::Doing(minutes) => println!("In progress: {minutes} min"),
        Status::Done => println!("Finished"),
    }
}
